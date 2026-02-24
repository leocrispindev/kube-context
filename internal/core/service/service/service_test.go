package service

import (
	"context"
	"testing"

	"k8s-agent-new/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type ServiceServiceTestSuite struct {
	suite.Suite
	client *fake.Clientset
	svc    *Service
}

func (s *ServiceServiceTestSuite) SetupTest() {
	s.client = fake.NewSimpleClientset()
	s.svc = NewService(s.client)
}

func TestServiceServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceServiceTestSuite))
}

func (s *ServiceServiceTestSuite) seedService(name, namespace string) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": "web"},
		},
		Spec: v1.ServiceSpec{
			Type:      v1.ServiceTypeClusterIP,
			ClusterIP: "10.0.0.1",
			Selector:  map[string]string{"app": "web"},
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(8080),
					Protocol:   v1.ProtocolTCP,
				},
			},
		},
	}
	_, err := s.client.CoreV1().Services(namespace).Create(context.Background(), svc, metav1.CreateOptions{})
	s.Require().NoError(err)
}

// --- ListServices ---

func (s *ServiceServiceTestSuite) TestListServices_Success() {
	s.seedService("svc-1", "default")
	s.seedService("svc-2", "default")

	result, err := s.svc.ListServices(context.Background(), "default")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Services, 2)
	assert.Equal(s.T(), "svc-1", result.Services[0].Name)
	assert.Equal(s.T(), "svc-2", result.Services[1].Name)
	assert.Equal(s.T(), "ClusterIP", result.Services[0].Type)
}

func (s *ServiceServiceTestSuite) TestListServices_EmptyNamespace() {
	result, err := s.svc.ListServices(context.Background(), "empty-ns")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Services)
}

func (s *ServiceServiceTestSuite) TestListServices_Error() {
	s.client.PrependReactor("list", "services", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, assert.AnError
	})

	result, err := s.svc.ListServices(context.Background(), "default")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to list services")
}

// --- GetService ---

func (s *ServiceServiceTestSuite) TestGetService_Success() {
	s.seedService("my-svc", "default")

	result, err := s.svc.GetService(context.Background(), "default", "my-svc")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "my-svc", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), "ClusterIP", result.Type)
	assert.Equal(s.T(), "10.0.0.1", result.ClusterIP)
	assert.Equal(s.T(), map[string]string{"app": "web"}, result.Selector)
	assert.Equal(s.T(), map[string]string{"app": "web"}, result.Labels)
	assert.Len(s.T(), result.Ports, 1)
	assert.Equal(s.T(), "http", result.Ports[0].Name)
	assert.Equal(s.T(), int32(80), result.Ports[0].Port)
	assert.Equal(s.T(), int32(8080), result.Ports[0].TargetPort)
	assert.Equal(s.T(), "TCP", result.Ports[0].Protocol)
}

func (s *ServiceServiceTestSuite) TestGetService_NotFound() {
	result, err := s.svc.GetService(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get service")
}

// --- CreateService ---

func (s *ServiceServiceTestSuite) TestCreateService_Success() {
	create := &dto.ServiceCreate{
		Name:      "new-svc",
		Namespace: "default",
		Type:      "ClusterIP",
		Selector:  map[string]string{"app": "api"},
		Ports: []dto.ServicePort{
			{
				Name:       "http",
				Port:       80,
				TargetPort: 3000,
				Protocol:   "TCP",
			},
		},
		Labels: map[string]string{"env": "prod"},
	}

	result, err := s.svc.CreateService(context.Background(), create)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "new-svc", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), "ClusterIP", result.Type)
	assert.Equal(s.T(), map[string]string{"app": "api"}, result.Selector)
	assert.Equal(s.T(), map[string]string{"env": "prod"}, result.Labels)
	assert.Len(s.T(), result.Ports, 1)
	assert.Equal(s.T(), "http", result.Ports[0].Name)
	assert.Equal(s.T(), int32(80), result.Ports[0].Port)
	assert.Equal(s.T(), int32(3000), result.Ports[0].TargetPort)
	assert.Equal(s.T(), "TCP", result.Ports[0].Protocol)
}

func (s *ServiceServiceTestSuite) TestCreateService_Error() {
	s.client.PrependReactor("create", "services", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, assert.AnError
	})

	create := &dto.ServiceCreate{
		Name:      "fail-svc",
		Namespace: "default",
		Type:      "ClusterIP",
	}

	result, err := s.svc.CreateService(context.Background(), create)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to create service")
}

// --- UpdateService ---

func (s *ServiceServiceTestSuite) TestUpdateService_Success() {
	s.seedService("update-svc", "default")

	updates := &dto.ServiceUpdate{
		Type:     "NodePort",
		Selector: map[string]string{"version": "v2"},
		Labels:   map[string]string{"tier": "frontend"},
	}

	result, err := s.svc.UpdateService(context.Background(), "default", "update-svc", updates)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "NodePort", result.Type)
	assert.Equal(s.T(), "v2", result.Selector["version"])
	assert.Equal(s.T(), "web", result.Selector["app"])
	assert.Equal(s.T(), "frontend", result.Labels["tier"])
	assert.Equal(s.T(), "web", result.Labels["app"])
}

func (s *ServiceServiceTestSuite) TestUpdateService_NotFound() {
	updates := &dto.ServiceUpdate{
		Type: "NodePort",
	}

	result, err := s.svc.UpdateService(context.Background(), "default", "nonexistent", updates)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get service")
}

// --- DeleteService ---

func (s *ServiceServiceTestSuite) TestDeleteService_Success() {
	s.seedService("delete-svc", "default")

	err := s.svc.DeleteService(context.Background(), "default", "delete-svc")

	assert.NoError(s.T(), err)

	_, getErr := s.svc.GetService(context.Background(), "default", "delete-svc")
	assert.Error(s.T(), getErr)
}

func (s *ServiceServiceTestSuite) TestDeleteService_NotFound() {
	err := s.svc.DeleteService(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to delete service")
}
