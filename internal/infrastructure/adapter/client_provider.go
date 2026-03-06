package adapter

import (
	"fmt"
	"k8s-agent-new/internal/infrastructure/auth"
	"net/http"
	"os"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientProvider struct {
	clientset *kubernetes.Clientset
}

var (
	clientProvider     *ClientProvider
	clientProviderOnce sync.Once
	clientProviderErr  error
)

func InitClientProvider() (*ClientProvider, error) {
	clientProviderOnce.Do(func() {
		config, err := buildConfig()
		if err != nil {
			clientProviderErr = err
			return
		}
		//Implement `Impersonnation` in every request
		// this allows to impersonate a user in the request
		// this allows different users to have different permissions in the request
		// this is a security feature that allows to restrict the access to the resources
		config.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
			return &auth.ImpersonationRoundTripper{Base: rt}
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			clientProviderErr = fmt.Errorf("failed to create clientset: %w", err)
			return
		}

		clientProvider = &ClientProvider{
			clientset: clientset,
		}
	})

	return clientProvider, clientProviderErr
}

func (p *ClientProvider) Clientset() *kubernetes.Clientset {
	return p.clientset
}

func buildConfig() (*rest.Config, error) {
	apiServerHost := os.Getenv("K8S_API_SERVER")
	insecure := os.Getenv("K8S_INSECURE") == "true"

	if apiServerHost != "" {
		config := &rest.Config{
			Host: apiServerHost,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: insecure,
			},
		}

		if token := os.Getenv("K8S_BASE_TOKEN"); token != "" {
			config.BearerToken = token
		}

		caCertPath := os.Getenv("K8S_CA_CERT_PATH")
		if caCertPath != "" {
			caCertData, err := os.ReadFile(caCertPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA certificate from %s: %w", caCertPath, err)
			}
			config.TLSClientConfig.CAData = caCertData
		} else if !insecure {
			return nil, fmt.Errorf("K8S_CA_CERT_PATH is required when K8S_INSECURE is false")
		}

		return config, nil
	}

	if os.Getenv("K8S_IN_CLUSTER") == "true" {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load in-cluster config: %w", err)
		}
		return config, nil
	}

	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve home directory: %w", err)
		}
		kubeconfigPath = homeDir + "/.kube/config"
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig from %s: %w", kubeconfigPath, err)
	}

	return config, nil
}
