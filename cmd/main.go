package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/leocrispindev/kube-context/internal/core/service/configmap"
	"github.com/leocrispindev/kube-context/internal/core/service/deployment"
	"github.com/leocrispindev/kube-context/internal/core/service/ingress"
	"github.com/leocrispindev/kube-context/internal/core/service/namespace"
	"github.com/leocrispindev/kube-context/internal/core/service/networkpolicy"
	"github.com/leocrispindev/kube-context/internal/core/service/pod"
	"github.com/leocrispindev/kube-context/internal/core/service/service"
	"github.com/leocrispindev/kube-context/internal/infrastructure/adapter"
	mcpSetup "github.com/leocrispindev/kube-context/internal/infrastructure/mcp"

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
