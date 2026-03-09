package tools

import (
	"context"

	"github.com/leocrispindev/kube-context/internal/core/dto"
	podService "github.com/leocrispindev/kube-context/internal/core/service/pod"

	mcp "github.com/metoro-io/mcp-golang"
)

type CreatePodArgs struct {
	dto.PodCreate
}

type UpdatePodArgs struct {
	NamespacedResourceArgs
	dto.PodUpdate
}

func RegisterPodTools(server *mcp.Server, svc *podService.Service) {
	server.RegisterTool("list_pods",
		"List all pods in a Kubernetes namespace.",
		func(ctx context.Context, args NamespacedArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.ListPods(ctx, args.Namespace)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("get_pod",
		"Get details of a specific pod by name and namespace.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.GetPod(ctx, args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("create_pod",
		"Create a new pod in Kubernetes.",
		func(ctx context.Context, args CreatePodArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.CreatePod(ctx, &args.PodCreate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("update_pod",
		"Update an existing pod (labels and/or image).",
		func(ctx context.Context, args UpdatePodArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.UpdatePod(ctx, args.Namespace, args.Name, &args.PodUpdate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("delete_pod",
		"Delete a pod from a Kubernetes namespace.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			if err := svc.DeletePod(ctx, args.Namespace, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("pod deleted")
		})

	server.RegisterTool("restart_pod",
		"Restart a pod by deleting it (the controller will recreate it).",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			if err := svc.RestartPod(ctx, args.Namespace, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("pod restarted")
		})
}
