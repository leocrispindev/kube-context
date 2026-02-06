package gateway

import "k8s.io/client-go/kubernetes"

type KubernetesCluster interface {
	Client() (*kubernetes.Clientset, error)
}
