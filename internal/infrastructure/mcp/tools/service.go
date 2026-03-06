package tools

import (
	"context"

	"k8s-agent-new/internal/core/dto"
	svcService "k8s-agent-new/internal/core/service/service"

	mcp "github.com/metoro-io/mcp-golang"
)

type CreateServiceArgs struct {
	dto.ServiceCreate
}

type UpdateServiceArgs struct {
	NamespacedResourceArgs
	dto.ServiceUpdate
}

func RegisterServiceTools(server *mcp.Server, svc *svcService.Service) {
	server.RegisterTool("list_services",
		"List all services in a Kubernetes namespace.",
		func(ctx context.Context, args NamespacedArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.ListServices(ctx, args.Namespace)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("get_service",
		"Get details of a specific Kubernetes service by name and namespace.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.GetService(ctx, args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("create_service",
		"Create a new Kubernetes service.",
		func(ctx context.Context, args CreateServiceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.CreateService(ctx, &args.ServiceCreate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("update_service",
		"Update an existing Kubernetes service (type, selector, labels).",
		func(ctx context.Context, args UpdateServiceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.UpdateService(ctx, args.Namespace, args.Name, &args.ServiceUpdate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("delete_service",
		"Delete a Kubernetes service from a namespace.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			if err := svc.DeleteService(ctx, args.Namespace, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("service deleted")
		})
}
