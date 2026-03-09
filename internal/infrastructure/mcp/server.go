package mcp

import (
	"github.com/leocrispindev/kube-context/internal/core/service/configmap"
	"github.com/leocrispindev/kube-context/internal/core/service/deployment"
	"github.com/leocrispindev/kube-context/internal/core/service/ingress"
	"github.com/leocrispindev/kube-context/internal/core/service/namespace"
	"github.com/leocrispindev/kube-context/internal/core/service/networkpolicy"
	"github.com/leocrispindev/kube-context/internal/core/service/pod"
	"github.com/leocrispindev/kube-context/internal/core/service/service"
	"github.com/leocrispindev/kube-context/internal/infrastructure/mcp/tools"

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
