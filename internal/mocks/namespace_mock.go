package mocks

import (
	"context"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	v1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/mock"
)

type MockNamespaceService struct {
	mock.Mock
}

func (m *MockNamespaceService) ListNamespaces(ctx context.Context) (*dto.NamespaceList, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.NamespaceList), args.Error(1)
}

func (m *MockNamespaceService) CreateNamespace(ctx context.Context, namespace string, labels map[string]string, annotations map[string]string) error {
	args := m.Called(ctx, namespace, labels, annotations)
	return args.Error(0)
}

func (m *MockNamespaceService) DeleteNamespace(ctx context.Context, namespace string) error {
	args := m.Called(ctx, namespace)
	return args.Error(0)
}

func (m *MockNamespaceService) GetNamespace(ctx context.Context, namespace string) (*dto.Namespace, error) {
	args := m.Called(ctx, namespace)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Namespace), args.Error(1)
}

func (m *MockNamespaceService) UpdateNamespace(ctx context.Context, namespace string, labels map[string]string, annotations map[string]string) error {
	args := m.Called(ctx, namespace, labels, annotations)
	return args.Error(0)
}

func (m *MockNamespaceService) UpdateResources(ctx context.Context, namespace string, resources v1.ResourceQuotaSpec) error {
	args := m.Called(ctx, namespace, resources)
	return args.Error(0)
}
