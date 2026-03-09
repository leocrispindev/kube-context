package tools

import (
	"context"

	"github.com/leocrispindev/kube-context/internal/core/dto"
	deploymentService "github.com/leocrispindev/kube-context/internal/core/service/deployment"

	mcp "github.com/metoro-io/mcp-golang"
)

type CreateDeploymentArgs struct {
	dto.DeploymentCreate
}

type UpdateDeploymentArgs struct {
	NamespacedResourceArgs
	dto.DeploymentUpdate
}

type TogglePauseDeploymentArgs struct {
	NamespacedResourceArgs
	Pause bool `json:"pause" jsonschema:"required,description:True to pause the deployment or false to resume it"`
}

func RegisterDeploymentTools(server *mcp.Server, svc *deploymentService.Service) {
	server.RegisterTool("list_deployments",
		"List all deployments in a Kubernetes namespace.",
		func(ctx context.Context, args NamespacedArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.ListDeployments(ctx, args.Namespace)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("get_deployment",
		"Get details of a specific deployment by name and namespace.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.GetDeployment(ctx, args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("create_deployment",
		"Create a new deployment in Kubernetes.",
		func(ctx context.Context, args CreateDeploymentArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.CreateDeployment(ctx, &args.DeploymentCreate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("update_deployment",
		"Update an existing deployment (replicas, images, labels).",
		func(ctx context.Context, args UpdateDeploymentArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			if err := svc.UpdateDeployment(ctx, args.Namespace, args.Name, args.DeploymentUpdate); err != nil {
				return nil, err
			}
			return msgResponse("deployment updated")
		})

	server.RegisterTool("delete_deployment",
		"Delete a deployment from a Kubernetes namespace.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			if err := svc.DeleteDeployment(ctx, args.Namespace, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("deployment deleted")
		})

	server.RegisterTool("get_rollout_status",
		"Get the rollout status of a deployment.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.GetRolloutStatus(ctx, args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("toggle_pause_deployment",
		"Pause or resume a deployment rollout.",
		func(ctx context.Context, args TogglePauseDeploymentArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			if err := svc.TogglePauseDeployment(ctx, args.Namespace, args.Name, args.Pause); err != nil {
				return nil, err
			}
			msg := "deployment resumed"
			if args.Pause {
				msg = "deployment paused"
			}
			return msgResponse(msg)
		})
}
