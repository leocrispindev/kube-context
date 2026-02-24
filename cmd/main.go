package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"k8s-agent-new/internal/core/service/configmap"
	"k8s-agent-new/internal/core/service/deployment"
	"k8s-agent-new/internal/core/service/ingress"
	"k8s-agent-new/internal/core/service/namespace"
	"k8s-agent-new/internal/core/service/networkpolicy"
	"k8s-agent-new/internal/core/service/pod"
	"k8s-agent-new/internal/core/service/service"
	"k8s-agent-new/internal/infrastructure/adapter"
	mcpSetup "k8s-agent-new/internal/infrastructure/mcp"

	"github.com/gin-gonic/gin"
	mcphttp "github.com/metoro-io/mcp-golang/transport/http"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	done := make(chan struct{})
	useStdio := flag.Bool("stdio", false, "use stdio transport (stdin/stdout) instead of HTTP")
	flag.Parse()

	clientProvider, err := adapter.InitClientProvider()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	clientset := clientProvider.Clientset()

	svcs := mcpSetup.Services{
		Pod:           pod.NewService(clientset),
		Namespace:     namespace.NewService(clientset),
		Deployment:    deployment.NewService(clientset),
		ConfigMap:     configmap.NewService(clientset),
		Service:       service.NewService(clientset),
		Ingress:       ingress.NewService(clientset),
		NetworkPolicy: networkpolicy.NewService(clientset),
	}

	if *useStdio {
		fmt.Fprintln(os.Stderr, "MCP server (stdio) starting...")
		transport := stdio.NewStdioServerTransport()
		server := mcpSetup.NewMCPServer(transport, svcs)
		fmt.Fprintln(os.Stderr, "MCP server (stdio) running. Waiting for messages.")
		if err := server.Serve(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		<-done
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	transport := mcphttp.NewGinTransport()
	server := mcpSetup.NewMCPServer(transport, svcs)
	go server.Serve()

	r := gin.Default()
	r.POST("/mcp", transport.Handler())
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	log.Println("==============================================")
	log.Println("  Kubernetes Agent MCP Server")
	log.Println("==============================================")
	log.Printf("MCP endpoint: POST http://localhost:%s/mcp", port)
	log.Printf("Health check: GET  http://localhost:%s/health", port)
	log.Println("==============================================")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
