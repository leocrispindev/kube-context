package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"k8s-agent-new/internal/infrastructure/adapter"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT_SECRET is required"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Expected: Bearer <token>"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token cannot be empty"})
			c.Abort()
			return
		}

		claims, err := parseAndValidateJWT(token, []byte(secret))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		info := buildImpersonationInfo(claims)
		if info.User == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT does not contain a valid user identity"})
			c.Abort()
			return
		}

		ctx := adapter.ContextWithImpersonation(c.Request.Context(), info)
		c.Request = c.Request.WithContext(ctx)
		c.Set(string(adapter.ImpersonationContextKey), info)

		c.Next()
	}
}

func parseAndValidateJWT(token string, secret []byte) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	headerBytes, err := decodeBase64URL(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid JWT header encoding")
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("invalid JWT header")
	}

	if alg, ok := header["alg"].(string); ok && alg != "HS256" {
		return nil, fmt.Errorf("unsupported JWT alg: %s", alg)
	}

	payloadBytes, err := decodeBase64URL(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid JWT payload encoding")
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("invalid JWT payload")
	}

	if err := validateSignature(parts[0], parts[1], parts[2], secret); err != nil {
		return nil, err
	}

	if exp, ok := claims["exp"]; ok {
		if err := validateExpiration(exp); err != nil {
			return nil, err
		}
	}

	return claims, nil
}

func validateSignature(headerPart string, payloadPart string, signaturePart string, secret []byte) error {
	message := headerPart + "." + payloadPart
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(message))
	expected := mac.Sum(nil)

	signature, err := decodeBase64URL(signaturePart)
	if err != nil {
		return fmt.Errorf("invalid JWT signature encoding")
	}

	if !hmac.Equal(signature, expected) {
		return fmt.Errorf("invalid JWT signature")
	}

	return nil
}

func validateExpiration(exp interface{}) error {
	var expUnix int64

	switch value := exp.(type) {
	case float64:
		expUnix = int64(value)
	case int64:
		expUnix = value
	case json.Number:
		parsed, err := value.Int64()
		if err != nil {
			return fmt.Errorf("invalid JWT exp claim")
		}
		expUnix = parsed
	default:
		return fmt.Errorf("invalid JWT exp claim")
	}

	if time.Now().Unix() > expUnix {
		return fmt.Errorf("JWT token has expired")
	}

	return nil
}

func decodeBase64URL(value string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(value)
}

func buildImpersonationInfo(claims map[string]interface{}) adapter.ImpersonationInfo {
	info := adapter.ImpersonationInfo{
		Extras: make(map[string][]string),
	}

	info.User = firstClaimString(claims, []string{"sub", "user", "username", "preferred_username", "email"})
	info.Groups = extractStringSlice(claims["groups"])

	extras := claims["extra"]
	if extras == nil {
		extras = claims["extras"]
	}
	info.Extras = extractExtras(extras)

	return info
}

func firstClaimString(claims map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if value, ok := claims[key]; ok {
			if str, ok := value.(string); ok && str != "" {
				return str
			}
		}
	}
	return ""
}

func extractStringSlice(value interface{}) []string {
	switch typed := value.(type) {
	case []interface{}:
		var result []string
		for _, item := range typed {
			if str, ok := item.(string); ok && str != "" {
				result = append(result, str)
			}
		}
		return result
	case []string:
		return typed
	case string:
		if typed == "" {
			return nil
		}
		return []string{typed}
	default:
		return nil
	}
}

func extractExtras(value interface{}) map[string][]string {
	result := make(map[string][]string)
	if value == nil {
		return result
	}

	switch typed := value.(type) {
	case map[string]interface{}:
		for key, raw := range typed {
			result[key] = extractStringSlice(raw)
		}
	case map[string][]string:
		return typed
	}

	return result
}
