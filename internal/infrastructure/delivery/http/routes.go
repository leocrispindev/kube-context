package http

import (
	"k8s-agent-new/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Pod           PodHandler
	Namespace     NamespaceHandler
	Deployment    DeploymentHandler
	ConfigMap     ConfigMapHandler
	Service       ServiceHandler
	Ingress       IngressHandler
	NetworkPolicy NetworkPolicyHandler
}

type PodHandler interface {
	GetPods(*gin.Context)
	GetPod(*gin.Context)
	CreatePod(*gin.Context)
	UpdatePod(*gin.Context)
	DeletePod(*gin.Context)
	RestartPod(*gin.Context)
}

type NamespaceHandler interface {
	GetNamespaces(*gin.Context)
	GetNamespace(*gin.Context)
	CreateNamespace(*gin.Context)
	DeleteNamespace(*gin.Context)
	UpdateNamespace(*gin.Context)
}

type DeploymentHandler interface {
	GetDeployments(*gin.Context)
	GetDeployment(*gin.Context)
	CreateDeployment(*gin.Context)
	UpdateDeployment(*gin.Context)
	DeleteDeployment(*gin.Context)
	GetRolloutStatus(*gin.Context)
	TogglePauseDeployment(*gin.Context)
}

type ConfigMapHandler interface {
	GetConfigMaps(*gin.Context)
	GetConfigMap(*gin.Context)
	CreateConfigMap(*gin.Context)
	UpdateConfigMap(*gin.Context)
	DeleteConfigMap(*gin.Context)
}

type ServiceHandler interface {
	GetServices(*gin.Context)
	GetService(*gin.Context)
	CreateService(*gin.Context)
	UpdateService(*gin.Context)
	DeleteService(*gin.Context)
}

type IngressHandler interface {
	GetIngresses(*gin.Context)
	GetIngress(*gin.Context)
	CreateIngress(*gin.Context)
	UpdateIngress(*gin.Context)
	DeleteIngress(*gin.Context)
}

type NetworkPolicyHandler interface {
	GetNetworkPolicies(*gin.Context)
	GetNetworkPolicy(*gin.Context)
	CreateNetworkPolicy(*gin.Context)
	UpdateNetworkPolicy(*gin.Context)
	DeleteNetworkPolicy(*gin.Context)
}

func AddRoutes(r *gin.Engine, handlers Handlers) {
	healthGroup := r.Group("/health")
	healthGroup.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	k8sGroup := r.Group("/kubernetes")
	k8sGroup.Use(middleware.JWTAuthMiddleware())
	{
		namespaceGroup := k8sGroup.Group("/namespaces")
		namespaceGroup.GET("", handlers.Namespace.GetNamespaces)
		namespaceGroup.GET("/:name", handlers.Namespace.GetNamespace)
		namespaceGroup.POST("", handlers.Namespace.CreateNamespace)
		namespaceGroup.DELETE("/:name", handlers.Namespace.DeleteNamespace)
		namespaceGroup.PUT("/:name", handlers.Namespace.UpdateNamespace)

		podGroup := k8sGroup.Group("/:namespace/pod")
		podGroup.GET("", handlers.Pod.GetPods)
		podGroup.GET("/:name", handlers.Pod.GetPod)
		podGroup.POST("", handlers.Pod.CreatePod)
		podGroup.PUT("/:name", handlers.Pod.UpdatePod)
		podGroup.DELETE("/:name", handlers.Pod.DeletePod)
		podGroup.POST("/:name/restart", handlers.Pod.RestartPod)

		deploymentGroup := k8sGroup.Group("/:namespace/deployments")
		deploymentGroup.GET("", handlers.Deployment.GetDeployments)
		deploymentGroup.GET("/:name", handlers.Deployment.GetDeployment)
		deploymentGroup.POST("", handlers.Deployment.CreateDeployment)
		deploymentGroup.PATCH("/:name", handlers.Deployment.UpdateDeployment)
		deploymentGroup.DELETE("/:name", handlers.Deployment.DeleteDeployment)
		deploymentGroup.GET("/:name/rollout/status", handlers.Deployment.GetRolloutStatus)
		deploymentGroup.POST("/:name/pause", handlers.Deployment.TogglePauseDeployment)

		configMapGroup := k8sGroup.Group("/:namespace/configmaps")
		configMapGroup.GET("", handlers.ConfigMap.GetConfigMaps)
		configMapGroup.GET("/:name", handlers.ConfigMap.GetConfigMap)
		configMapGroup.POST("", handlers.ConfigMap.CreateConfigMap)
		configMapGroup.PUT("/:name", handlers.ConfigMap.UpdateConfigMap)
		configMapGroup.DELETE("/:name", handlers.ConfigMap.DeleteConfigMap)

		serviceGroup := k8sGroup.Group("/:namespace/services")
		serviceGroup.GET("", handlers.Service.GetServices)
		serviceGroup.GET("/:name", handlers.Service.GetService)
		serviceGroup.POST("", handlers.Service.CreateService)
		serviceGroup.PUT("/:name", handlers.Service.UpdateService)
		serviceGroup.DELETE("/:name", handlers.Service.DeleteService)

		ingressGroup := k8sGroup.Group("/:namespace/ingresses")
		ingressGroup.GET("", handlers.Ingress.GetIngresses)
		ingressGroup.GET("/:name", handlers.Ingress.GetIngress)
		ingressGroup.POST("", handlers.Ingress.CreateIngress)
		ingressGroup.PUT("/:name", handlers.Ingress.UpdateIngress)
		ingressGroup.DELETE("/:name", handlers.Ingress.DeleteIngress)

		networkPolicyGroup := k8sGroup.Group("/:namespace/networkpolicies")
		networkPolicyGroup.GET("", handlers.NetworkPolicy.GetNetworkPolicies)
		networkPolicyGroup.GET("/:name", handlers.NetworkPolicy.GetNetworkPolicy)
		networkPolicyGroup.POST("", handlers.NetworkPolicy.CreateNetworkPolicy)
		networkPolicyGroup.PUT("/:name", handlers.NetworkPolicy.UpdateNetworkPolicy)
		networkPolicyGroup.DELETE("/:name", handlers.NetworkPolicy.DeleteNetworkPolicy)
	}
}
