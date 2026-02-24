package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s-agent-new/internal/infrastructure/auth"

	mcp "github.com/metoro-io/mcp-golang"
)

type AuthArgs struct {
	Token string `json:"token" jsonschema:"required,description:JWT Bearer token for authentication and RBAC impersonation"`
}

type NamespacedArgs struct {
	AuthArgs
	Namespace string `json:"namespace" jsonschema:"required,description:Kubernetes namespace"`
}

type NamespacedResourceArgs struct {
	NamespacedArgs
	Name string `json:"name" jsonschema:"required,description:Resource name"`
}

func authenticate(ctx context.Context, token string) (context.Context, error) {
	return auth.ValidateTokenAndBuildContext(ctx, token)
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

// Session-mode arg types — no token field; auth is resolved once at session init.

type EmptyArgs struct{}

type SessionNamespacedArgs struct {
	Namespace string `json:"namespace" jsonschema:"required,description:Kubernetes namespace"`
}

type SessionNamespacedResourceArgs struct {
	SessionNamespacedArgs
	Name string `json:"name" jsonschema:"required,description:Resource name"`
}
