package auth

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey string

const bearerTokenContextKey contextKey = "k8s_bearer_token"

// ValidateTokenAndBuildContext keeps the existing call site contract, but now
// just validates presence and stores the raw bearer token in context.
func ValidateTokenAndBuildContext(ctx context.Context, token string) (context.Context, error) {
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}
	return context.WithValue(ctx, bearerTokenContextKey, token), nil
}

func bearerTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(bearerTokenContextKey).(string)
	return token, ok
}

// BearerTokenRoundTripper injects the request-scoped token into Kubernetes API calls.
type BearerTokenRoundTripper struct {
	Base http.RoundTripper
}

func (rt *BearerTokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token, ok := bearerTokenFromContext(req.Context())
	if !ok || token == "" {
		return rt.Base.RoundTrip(req)
	}

	cloned := req.Clone(req.Context())
	cloned.Header.Set("Authorization", "Bearer "+token)
	return rt.Base.RoundTrip(cloned)
}
