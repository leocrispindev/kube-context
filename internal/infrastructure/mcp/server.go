package mcp

import (
	"k8s-agent-new/internal/core/service/configmap"
	"k8s-agent-new/internal/core/service/deployment"
	"k8s-agent-new/internal/core/service/ingress"
	"k8s-agent-new/internal/core/service/namespace"
	"k8s-agent-new/internal/core/service/networkpolicy"
	"k8s-agent-new/internal/core/service/pod"
	"k8s-agent-new/internal/core/service/service"
	"k8s-agent-new/internal/infrastructure/mcp/tools"

	mcplib "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport"
)

type Services struct {
	Pod           *pod.Service
	Namespace     *namespace.Service
	Deployment    *deployment.Service
	ConfigMap     *configmap.Service
	Service       *service.Service
	Ingress       *ingress.Service
	NetworkPolicy *networkpolicy.Service
}

func NewMCPServer(t transport.Transport, svcs Services) *mcplib.Server {
	server := mcplib.NewServer(t)

	tools.RegisterNamespaceTools(server, svcs.Namespace)
	tools.RegisterPodTools(server, svcs.Pod)
	tools.RegisterDeploymentTools(server, svcs.Deployment)
	tools.RegisterConfigMapTools(server, svcs.ConfigMap)
	tools.RegisterServiceTools(server, svcs.Service)
	tools.RegisterIngressTools(server, svcs.Ingress)
	tools.RegisterNetworkPolicyTools(server, svcs.NetworkPolicy)

	return server
}
