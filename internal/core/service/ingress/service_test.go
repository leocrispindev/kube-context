package ingress

import (
	"context"
	"testing"

	"k8s-agent-new/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type IngressServiceTestSuite struct {
	suite.Suite
	client *fake.Clientset
	svc    *Service
}

func (s *IngressServiceTestSuite) SetupTest() {
	s.client = fake.NewSimpleClientset()
	s.svc = NewService(s.client)
}

func TestIngressServiceSuite(t *testing.T) {
	suite.Run(t, new(IngressServiceTestSuite))
}

func (s *IngressServiceTestSuite) seedIngress(name, namespace string) {
	pathType := networkingv1.PathTypePrefix
	ingressClass := "nginx"
	ing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      map[string]string{"app": "web"},
			Annotations: map[string]string{"nginx.ingress.kubernetes.io/rewrite-target": "/"},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClass,
			Rules: []networkingv1.IngressRule{
				{
					Host: "example.com",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/api",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "backend-svc",
											Port: networkingv1.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
			},
			TLS: []networkingv1.IngressTLS{
				{Hosts: []string{"example.com"}, SecretName: "tls-secret"},
			},
		},
	}
	_, err := s.client.NetworkingV1().Ingresses(namespace).Create(context.Background(), ing, metav1.CreateOptions{})
	s.Require().NoError(err)
}

// --- ListIngresses ---

func (s *IngressServiceTestSuite) TestListIngresses_Success() {
	s.seedIngress("ing-1", "default")
	s.seedIngress("ing-2", "default")

	result, err := s.svc.ListIngresses(context.Background(), "default")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Ingresses, 2)
	assert.Equal(s.T(), "ing-1", result.Ingresses[0].Name)
	assert.Equal(s.T(), "ing-2", result.Ingresses[1].Name)
	assert.Equal(s.T(), "nginx", result.Ingresses[0].IngressClass)
}

func (s *IngressServiceTestSuite) TestListIngresses_EmptyNamespace() {
	result, err := s.svc.ListIngresses(context.Background(), "empty-ns")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Ingresses)
}

func (s *IngressServiceTestSuite) TestListIngresses_Error() {
	s.client.PrependReactor("list", "ingresses", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, assert.AnError
	})

	result, err := s.svc.ListIngresses(context.Background(), "default")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to list ingresses")
}

// --- GetIngress ---

func (s *IngressServiceTestSuite) TestGetIngress_Success() {
	s.seedIngress("my-ingress", "default")

	result, err := s.svc.GetIngress(context.Background(), "default", "my-ingress")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "my-ingress", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), "nginx", result.IngressClass)
	assert.Len(s.T(), result.Rules, 1)
	assert.Equal(s.T(), "example.com", result.Rules[0].Host)
	assert.Len(s.T(), result.Rules[0].Paths, 1)
	assert.Equal(s.T(), "/api", result.Rules[0].Paths[0].Path)
	assert.Equal(s.T(), "Prefix", result.Rules[0].Paths[0].PathType)
	assert.Equal(s.T(), "backend-svc", result.Rules[0].Paths[0].ServiceName)
	assert.Equal(s.T(), int32(8080), result.Rules[0].Paths[0].ServicePort)
	assert.Len(s.T(), result.TLS, 1)
	assert.Equal(s.T(), "tls-secret", result.TLS[0].SecretName)
	assert.Equal(s.T(), []string{"example.com"}, result.TLS[0].Hosts)
	assert.Equal(s.T(), map[string]string{"app": "web"}, result.Labels)
}

func (s *IngressServiceTestSuite) TestGetIngress_NotFound() {
	result, err := s.svc.GetIngress(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get ingress")
}

// --- CreateIngress ---

func (s *IngressServiceTestSuite) TestCreateIngress_Success() {
	create := &dto.IngressCreate{
		Name:         "new-ingress",
		Namespace:    "default",
		IngressClass: "nginx",
		Rules: []dto.IngressRule{
			{
				Host: "new.example.com",
				Paths: []dto.IngressPath{
					{
						Path:        "/",
						PathType:    "Prefix",
						ServiceName: "frontend-svc",
						ServicePort: 80,
					},
				},
			},
		},
		TLS: []dto.IngressTLS{
			{Hosts: []string{"new.example.com"}, SecretName: "new-tls"},
		},
		Labels:      map[string]string{"env": "staging"},
		Annotations: map[string]string{"key": "value"},
	}

	result, err := s.svc.CreateIngress(context.Background(), create)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "new-ingress", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), "nginx", result.IngressClass)
	assert.Len(s.T(), result.Rules, 1)
	assert.Equal(s.T(), "new.example.com", result.Rules[0].Host)
	assert.Equal(s.T(), "/", result.Rules[0].Paths[0].Path)
	assert.Equal(s.T(), "Prefix", result.Rules[0].Paths[0].PathType)
	assert.Equal(s.T(), "frontend-svc", result.Rules[0].Paths[0].ServiceName)
	assert.Equal(s.T(), int32(80), result.Rules[0].Paths[0].ServicePort)
	assert.Len(s.T(), result.TLS, 1)
	assert.Equal(s.T(), "new-tls", result.TLS[0].SecretName)
	assert.Equal(s.T(), map[string]string{"env": "staging"}, result.Labels)
	assert.Equal(s.T(), map[string]string{"key": "value"}, result.Annotations)
}

func (s *IngressServiceTestSuite) TestCreateIngress_Error() {
	s.client.PrependReactor("create", "ingresses", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, assert.AnError
	})

	create := &dto.IngressCreate{
		Name:      "fail-ingress",
		Namespace: "default",
	}

	result, err := s.svc.CreateIngress(context.Background(), create)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to create ingress")
}

// --- UpdateIngress ---

func (s *IngressServiceTestSuite) TestUpdateIngress_Success() {
	s.seedIngress("update-ingress", "default")

	updates := &dto.IngressUpdate{
		IngressClass: "traefik",
		Labels:       map[string]string{"version": "v2"},
		Annotations:  map[string]string{"custom": "annotation"},
	}

	result, err := s.svc.UpdateIngress(context.Background(), "default", "update-ingress", updates)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "traefik", result.IngressClass)
	assert.Equal(s.T(), "v2", result.Labels["version"])
	assert.Equal(s.T(), "web", result.Labels["app"])
	assert.Equal(s.T(), "annotation", result.Annotations["custom"])
	assert.Equal(s.T(), "/", result.Annotations["nginx.ingress.kubernetes.io/rewrite-target"])
}

func (s *IngressServiceTestSuite) TestUpdateIngress_NotFound() {
	updates := &dto.IngressUpdate{
		Labels: map[string]string{"version": "v2"},
	}

	result, err := s.svc.UpdateIngress(context.Background(), "default", "nonexistent", updates)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get ingress")
}

// --- DeleteIngress ---

func (s *IngressServiceTestSuite) TestDeleteIngress_Success() {
	s.seedIngress("delete-ingress", "default")

	err := s.svc.DeleteIngress(context.Background(), "default", "delete-ingress")

	assert.NoError(s.T(), err)

	_, getErr := s.svc.GetIngress(context.Background(), "default", "delete-ingress")
	assert.Error(s.T(), getErr)
}

func (s *IngressServiceTestSuite) TestDeleteIngress_NotFound() {
	err := s.svc.DeleteIngress(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to delete ingress")
}
