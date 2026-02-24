package tools

import (
	"context"

	namespaceService "k8s-agent-new/internal/core/service/namespace"

	mcp "github.com/metoro-io/mcp-golang"
)

type GetNamespaceArgs struct {
	AuthArgs
	Name string `json:"name" jsonschema:"required,description:Namespace name"`
}

type CreateNamespaceArgs struct {
	AuthArgs
	Name        string            `json:"name" jsonschema:"required,description:Namespace name to create"`
	Labels      map[string]string `json:"labels,omitempty" jsonschema:"description:Labels for the namespace"`
	Annotations map[string]string `json:"annotations,omitempty" jsonschema:"description:Annotations for the namespace"`
}

type UpdateNamespaceArgs struct {
	AuthArgs
	Name        string            `json:"name" jsonschema:"required,description:Namespace name to update"`
	Labels      map[string]string `json:"labels,omitempty" jsonschema:"description:Updated labels"`
	Annotations map[string]string `json:"annotations,omitempty" jsonschema:"description:Updated annotations"`
}

type DeleteNamespaceArgs struct {
	AuthArgs
	Name string `json:"name" jsonschema:"required,description:Namespace name to delete"`
}

func RegisterNamespaceTools(server *mcp.Server, svc *namespaceService.Service) {
	server.RegisterTool("list_namespaces",
		"List all Kubernetes namespaces accessible by the authenticated user.",
		func(args AuthArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.ListNamespaces(ctx)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("get_namespace",
		"Get details of a specific Kubernetes namespace.",
		func(args GetNamespaceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.GetNamespace(ctx, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("create_namespace",
		"Create a new Kubernetes namespace.",
		func(args CreateNamespaceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			if err := svc.CreateNamespace(ctx, args.Name, args.Labels, args.Annotations); err != nil {
				return nil, err
			}
			return msgResponse("namespace created")
		})

	server.RegisterTool("delete_namespace",
		"Delete a Kubernetes namespace.",
		func(args DeleteNamespaceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			if err := svc.DeleteNamespace(ctx, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("namespace deleted")
		})

	server.RegisterTool("update_namespace",
		"Update labels and annotations of a Kubernetes namespace.",
		func(args UpdateNamespaceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			if err := svc.UpdateNamespace(ctx, args.Name, args.Labels, args.Annotations); err != nil {
				return nil, err
			}
			return msgResponse("namespace updated")
		})
}
