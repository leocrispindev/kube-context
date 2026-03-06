package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func ValidateTokenAndBuildContext(ctx context.Context, token string) (context.Context, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	claims, err := parseAndValidateJWT(token, []byte(secret))
	if err != nil {
		return nil, err
	}

	info := buildImpersonationInfo(claims)
	if info.User == "" {
		return nil, fmt.Errorf("JWT does not contain a valid user identity")
	}

	return contextWithImpersonation(ctx, info), nil
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

func validateSignature(headerPart, payloadPart, signaturePart string, secret []byte) error {
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

func buildImpersonationInfo(claims map[string]interface{}) ImpersonationInfo {
	info := ImpersonationInfo{
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
