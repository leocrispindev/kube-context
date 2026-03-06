package auth

import (
	"context"
	"net/http"
)

type contextKey string

const ImpersonationContextKey contextKey = "k8s_impersonation_info"

type ImpersonationInfo struct {
	User   string
	Groups []string
	Extras map[string][]string
}

func contextWithImpersonation(ctx context.Context, info ImpersonationInfo) context.Context {
	return context.WithValue(ctx, ImpersonationContextKey, info)
}

func impersonationFromContext(ctx context.Context) (ImpersonationInfo, bool) {
	info, ok := ctx.Value(ImpersonationContextKey).(ImpersonationInfo)
	return info, ok
}

type ImpersonationRoundTripper struct {
	Base http.RoundTripper
}

func (rt *ImpersonationRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	info, ok := impersonationFromContext(req.Context())
	if !ok || info.User == "" {
		return rt.Base.RoundTrip(req)
	}

	cloned := req.Clone(req.Context())
	cloned.Header.Set("Impersonate-User", info.User)
	for _, group := range info.Groups {
		if group != "" {
			cloned.Header.Add("Impersonate-Group", group)
		}
	}

	for key, values := range info.Extras {
		headerKey := "Impersonate-Extra-" + key
		for _, value := range values {
			if value != "" {
				cloned.Header.Add(headerKey, value)
			}
		}
	}

	return rt.Base.RoundTrip(cloned)
}
