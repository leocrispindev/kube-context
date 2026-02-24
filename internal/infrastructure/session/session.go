package session

import (
	"context"
	"fmt"

	"k8s-agent-new/internal/infrastructure/adapter"
	"k8s-agent-new/internal/infrastructure/auth"
)

// Session holds a pre-authenticated identity for the lifetime of a STDIO connection.
// The JWT is validated once at process startup; every subsequent tool call reuses
// this identity without requiring the token as a parameter.
type Session struct {
	identity adapter.ImpersonationInfo
}

func New(token string) (*Session, error) {
	if token == "" {
		return nil, fmt.Errorf("session token is required (set MCP_AUTH_TOKEN)")
	}

	ctx, err := auth.ValidateTokenAndBuildContext(context.Background(), token)
	if err != nil {
		return nil, fmt.Errorf("session auth failed: %w", err)
	}

	info, ok := adapter.ImpersonationFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract identity from token")
	}

	return &Session{identity: info}, nil
}

// Context returns a fresh context carrying the session's impersonation identity.
func (s *Session) Context() context.Context {
	return adapter.ContextWithImpersonation(context.Background(), s.identity)
}

func (s *Session) User() string {
	return s.identity.User
}
