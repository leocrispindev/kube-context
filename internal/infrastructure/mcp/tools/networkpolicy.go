package tools

import (
	"context"

	"k8s-agent-new/internal/core/dto"
	networkpolicyService "k8s-agent-new/internal/core/service/networkpolicy"

	mcp "github.com/metoro-io/mcp-golang"
)

type CreateNetworkPolicyArgs struct {
	dto.NetworkPolicyCreate
}

type UpdateNetworkPolicyArgs struct {
	NamespacedResourceArgs
	dto.NetworkPolicyUpdate
}

func RegisterNetworkPolicyTools(server *mcp.Server, svc *networkpolicyService.Service) {
	server.RegisterTool("list_network_policies",
		"List all network policies in a Kubernetes namespace.",
		func(ctx context.Context, args NamespacedArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.ListNetworkPolicies(ctx, args.Namespace)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("get_network_policy",
		"Get details of a specific network policy by name and namespace.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.GetNetworkPolicy(ctx, args.Namespace, args.Name)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("create_network_policy",
		"Create a new network policy in Kubernetes.",
		func(ctx context.Context, args CreateNetworkPolicyArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.CreateNetworkPolicy(ctx, &args.NetworkPolicyCreate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("update_network_policy",
		"Update an existing network policy (labels).",
		func(ctx context.Context, args UpdateNetworkPolicyArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			result, err := svc.UpdateNetworkPolicy(ctx, args.Namespace, args.Name, &args.NetworkPolicyUpdate)
			if err != nil {
				return nil, err
			}
			return jsonResponse(result)
		})

	server.RegisterTool("delete_network_policy",
		"Delete a network policy from a Kubernetes namespace.",
		func(ctx context.Context, args NamespacedResourceArgs) (*mcp.ToolResponse, error) {
			ctx, err := authenticate(ctx)
			if err != nil {
				return nil, err
			}
			if err := svc.DeleteNetworkPolicy(ctx, args.Namespace, args.Name); err != nil {
				return nil, err
			}
			return msgResponse("network policy deleted")
		})
}
