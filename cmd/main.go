package main

import (
	"log"
	"os"

	"k8s-agent-new/internal/infrastructure/delivery/http"
	configmapHandler "k8s-agent-new/internal/infrastructure/handler/configmap"
	deploymentHandler "k8s-agent-new/internal/infrastructure/handler/deployment"
	ingressHandler "k8s-agent-new/internal/infrastructure/handler/ingress"
	namespaceHandler "k8s-agent-new/internal/infrastructure/handler/namespace"
	networkpolicyHandler "k8s-agent-new/internal/infrastructure/handler/networkpolicy"
	podHandler "k8s-agent-new/internal/infrastructure/handler/pod"
	serviceHandler "k8s-agent-new/internal/infrastructure/handler/service"

	"k8s-agent-new/internal/core/service/configmap"
	"k8s-agent-new/internal/core/service/deployment"
	"k8s-agent-new/internal/core/service/ingress"
	"k8s-agent-new/internal/core/service/namespace"
	"k8s-agent-new/internal/core/service/networkpolicy"
	"k8s-agent-new/internal/core/service/pod"
	"k8s-agent-new/internal/core/service/service"
	"k8s-agent-new/internal/infrastructure/adapter"

	"github.com/gin-gonic/gin"
)

func startup() (http.Handlers, error) {
	clientProvider, err := adapter.InitClientProvider()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	clientset := clientProvider.Clientset()

	podSvc := pod.NewService(clientset)
	nsSvc := namespace.NewService(clientset)
	deploySvc := deployment.NewService(clientset)
	cmSvc := configmap.NewService(clientset)
	svcSvc := service.NewService(clientset)
	ingSvc := ingress.NewService(clientset)
	netpolSvc := networkpolicy.NewService(clientset)

	return http.Handlers{
		Pod:           podHandler.NewHandler(podSvc),
		Namespace:     namespaceHandler.NewHandler(nsSvc),
		Deployment:    deploymentHandler.NewHandler(deploySvc),
		ConfigMap:     configmapHandler.NewHandler(cmSvc),
		Service:       serviceHandler.NewHandler(svcSvc),
		Ingress:       ingressHandler.NewHandler(ingSvc),
		NetworkPolicy: networkpolicyHandler.NewHandler(netpolSvc),
	}, nil
}

func main() {
	handlers, err := startup()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	r := gin.Default()

	http.AddRoutes(r, handlers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("==============================================")
	log.Println("  Kubernetes Agent Server")
	log.Println("==============================================")
	log.Printf("Server starting on port %s", port)
	log.Println("==============================================")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
