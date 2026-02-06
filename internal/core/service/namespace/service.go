package namespace

import (
	"context"

	"k8s-agent-new/internal/core/dto"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	client *kubernetes.Clientset
}

func NewService(client *kubernetes.Clientset) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ListNamespaces(ctx context.Context) (*dto.NamespaceList, error) {
	namespaces, err := s.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceNames []dto.Namespace
	for _, ns := range namespaces.Items {
		namespaceNames = append(namespaceNames, dto.Namespace{
			Name:        ns.Name,
			Labels:      ns.Labels,
			Annotations: ns.Annotations,
			Status:      ns.Status,
		})
	}

	return &dto.NamespaceList{
		Namespaces: namespaceNames,
	}, nil
}

func (s *Service) CreateNamespace(ctx context.Context, namespace string, labels map[string]string, annotations map[string]string) error {
	_, err := s.client.CoreV1().Namespaces().Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespace,
			Labels:      labels,
			Annotations: annotations,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteNamespace(ctx context.Context, namespace string) error {
	return s.client.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
}

func (s *Service) GetNamespace(ctx context.Context, namespace string) (*dto.Namespace, error) {
	ns, err := s.client.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &dto.Namespace{
		Name:        ns.Name,
		Labels:      ns.Labels,
		Annotations: ns.Annotations,
		Status:      ns.Status,
	}, nil
}

func (s *Service) UpdateNamespace(ctx context.Context, namespace string, labels map[string]string, annotations map[string]string) error {
	_, err := s.client.CoreV1().Namespaces().Update(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        namespace,
			Labels:      labels,
			Annotations: annotations,
		},
	}, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateResources(ctx context.Context, namespace string, resources v1.ResourceQuotaSpec) error {
	_, err := s.client.CoreV1().Namespaces().Update(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Spec: v1.NamespaceSpec{
			Finalizers: []v1.FinalizerName{
				v1.FinalizerKubernetes,
			},
		},
	}, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
