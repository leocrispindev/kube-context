package tools

import (
	"context"

	"k8s-agent-new/internal/core/dto"
	ingressService "k8s-agent-new/internal/core/service/ingress"

	mcp "github.com/metoro-io/mcp-golang"
)

type CreateIngressArgs struct {
	AuthArgs
	dto.IngressCreate
}

type UpdateIngressArgs struct {
	NamespacedResourceArgs
	dto.IngressUpdate
}

func RegisterIngressTools(server *mcp.Server, svc *ingressService.Service) {
	server.RegisterTool("list_ingresses",
		"List all ingresses in a Kubernetes namespace.",
		func(args NamespacedArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.ListIngresses(ctx, args.Namespace)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("get_ingress",
		"Get details of a specific ingress by name and namespace.",
		func(args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.GetIngress(ctx, args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("create_ingress",
		"Create a new ingress in Kubernetes.",
		func(args CreateIngressArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.CreateIngress(ctx, &args.IngressCreate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("update_ingress",
		"Update an existing ingress (class, labels, annotations).",
		func(args UpdateIngressArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.UpdateIngress(ctx, args.Namespace, args.Name, &args.IngressUpdate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("delete_ingress",
		"Delete an ingress from a Kubernetes namespace.",
		func(args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			if err := svc.DeleteIngress(ctx, args.Namespace, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("ingress deleted")
		})
}
