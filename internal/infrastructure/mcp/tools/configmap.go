package tools

import (
	"context"

	"k8s-agent-new/internal/core/dto"
	configmapService "k8s-agent-new/internal/core/service/configmap"

	mcp "github.com/metoro-io/mcp-golang"
)

type CreateConfigMapArgs struct {
	AuthArgs
	dto.ConfigMapCreate
}

type UpdateConfigMapArgs struct {
	NamespacedResourceArgs
	dto.ConfigMapUpdate
}

func RegisterConfigMapTools(server *mcp.Server, svc *configmapService.Service) {
	server.RegisterTool("list_configmaps",
		"List all configmaps in a Kubernetes namespace.",
		func(args NamespacedArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.ListConfigMaps(ctx, args.Namespace)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("get_configmap",
		"Get details of a specific configmap by name and namespace.",
		func(args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.GetConfigMap(ctx, args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("create_configmap",
		"Create a new configmap in Kubernetes.",
		func(args CreateConfigMapArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.CreateConfigMap(ctx, &args.ConfigMapCreate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("update_configmap",
		"Update an existing configmap (data and/or labels).",
		func(args UpdateConfigMapArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			result, err := svc.UpdateConfigMap(ctx, args.Namespace, args.Name, &args.ConfigMapUpdate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("delete_configmap",
		"Delete a configmap from a Kubernetes namespace.",
		func(args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
			if err != nil {
				return nil, err
			}
			if err := svc.DeleteConfigMap(ctx, args.Namespace, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("configmap deleted")
		})
}
