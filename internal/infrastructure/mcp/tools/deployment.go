package tools

import (
	"context"
	"fmt"

	"k8s-agent-new/internal/core/dto"
	deploymentService "k8s-agent-new/internal/core/service/deployment"
	"k8s-agent-new/internal/infrastructure/session"

	mcp "github.com/metoro-io/mcp-golang"
)

type CreateDeploymentArgs struct {
	AuthArgs
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
		func(args NamespacedArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
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
		func(args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
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
		func(args CreateDeploymentArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
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
		func(args UpdateDeploymentArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
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
		func(args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
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
		func(args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
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
		func(args TogglePauseDeploymentArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(context.Background(), args.Token)
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
			return msgResponse(fmt.Sprintf(msg))
		})
}

// --- Session-mode (STDIO) ---

type SessionUpdateDeploymentArgs struct {
	SessionNamespacedResourceArgs
	dto.DeploymentUpdate
}

type SessionTogglePauseDeploymentArgs struct {
	SessionNamespacedResourceArgs
	Pause bool `json:"pause" jsonschema:"required,description:True to pause the deployment or false to resume it"`
}

func RegisterDeploymentSessionTools(server *mcp.Server, svc *deploymentService.Service, sess *session.Session) {
	server.RegisterTool("list_deployments",
		"List all deployments in a Kubernetes namespace.",
		func(args SessionNamespacedArgs) (*mcp.ToolResponse, error) {
			result, err := svc.ListDeployments(sess.Context(), args.Namespace)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("get_deployment",
		"Get details of a specific deployment by name and namespace.",
		func(args SessionNamespacedResourceArgs) (*mcp.ToolResponse, error) {
			result, err := svc.GetDeployment(sess.Context(), args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("create_deployment",
		"Create a new deployment in Kubernetes.",
		func(args dto.DeploymentCreate) (*mcp.ToolResponse, error) {
			result, err := svc.CreateDeployment(sess.Context(), &args)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("update_deployment",
		"Update an existing deployment (replicas, images, labels).",
		func(args SessionUpdateDeploymentArgs) (*mcp.ToolResponse, error) {
			if err := svc.UpdateDeployment(sess.Context(), args.Namespace, args.Name, args.DeploymentUpdate); err != nil {
				return nil, err
			}
			return msgResponse("deployment updated")
		})

	server.RegisterTool("delete_deployment",
		"Delete a deployment from a Kubernetes namespace.",
		func(args SessionNamespacedResourceArgs) (*mcp.ToolResponse, error) {
			if err := svc.DeleteDeployment(sess.Context(), args.Namespace, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("deployment deleted")
		})

	server.RegisterTool("get_rollout_status",
		"Get the rollout status of a deployment.",
		func(args SessionNamespacedResourceArgs) (*mcp.ToolResponse, error) {
			result, err := svc.GetRolloutStatus(sess.Context(), args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("toggle_pause_deployment",
		"Pause or resume a deployment rollout.",
		func(args SessionTogglePauseDeploymentArgs) (*mcp.ToolResponse, error) {
			if err := svc.TogglePauseDeployment(sess.Context(), args.Namespace, args.Name, args.Pause); err != nil {
				return nil, err
			}
			msg := "deployment resumed"
			if args.Pause {
				msg = "deployment paused"
			}
			return msgResponse(fmt.Sprintf(msg))
		})
}
