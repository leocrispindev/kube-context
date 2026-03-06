package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"k8s-agent-new/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
	mcp "github.com/metoro-io/mcp-golang"
)

type NamespacedArgs struct {
	Namespace string `json:"namespace" jsonschema:"required,description:Kubernetes namespace"`
}

type NamespacedResourceArgs struct {
	NamespacedArgs
	Name string `json:"name" jsonschema:"required,description:Resource name"`
}

func authenticate(ctx context.Context) (context.Context, error) {
	token, err := tokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return auth.ValidateTokenAndBuildContext(ctx, token)
}

func tokenFromContext(ctx context.Context) (string, error) {
	if ginCtx, ok := ctx.Value("ginContext").(*gin.Context); ok && ginCtx != nil {
		headerToken := strings.TrimSpace(ginCtx.GetHeader("Authorization"))
		if headerToken == "" {
			headerToken = strings.TrimSpace(ginCtx.GetHeader("X-Auth-Token"))
		}
		if headerToken == "" {
			return "", fmt.Errorf("missing authentication token in HTTP header (use Authorization: Bearer <token>)")
		}
		return normalizeToken(headerToken)
	}

	envToken := strings.TrimSpace(os.Getenv("K8S_TOKEN"))
	if envToken == "" {
		return "", fmt.Errorf("missing authentication token for stdio transport (set K8S_TOKEN)")
	}
	return envToken, nil
}

func normalizeToken(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("authentication token cannot be empty")
	}

	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		token := strings.TrimSpace(parts[1])
		if token == "" {
			return "", fmt.Errorf("bearer token cannot be empty")
		}
		return token, nil
	}

	return raw, nil
}

func jsonResponse(data any) (*mcp.ToolResponse, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	return mcp.NewToolResponse(mcp.NewTextContent(string(b))), nil
}

func msgResponse(msg string) (*mcp.ToolResponse, error) {
	return mcp.NewToolResponse(mcp.NewTextContent(msg)), nil
}

type EmptyArgs struct{}
