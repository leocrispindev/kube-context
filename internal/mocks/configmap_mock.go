package mocks

import (
	"context"

	"k8s-agent-new/internal/core/dto"

	"github.com/stretchr/testify/mock"
)

type MockConfigMapService struct {
	mock.Mock
}

func (m *MockConfigMapService) ListConfigMaps(ctx context.Context, namespace string) (*dto.ConfigMapList, error) {
	args := m.Called(ctx, namespace)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ConfigMapList), args.Error(1)
}

func (m *MockConfigMapService) GetConfigMap(ctx context.Context, namespace string, configMapName string) (*dto.ConfigMap, error) {
	args := m.Called(ctx, namespace, configMapName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ConfigMap), args.Error(1)
}

func (m *MockConfigMapService) CreateConfigMap(ctx context.Context, configMapCreate *dto.ConfigMapCreate) (*dto.ConfigMap, error) {
	args := m.Called(ctx, configMapCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ConfigMap), args.Error(1)
}

func (m *MockConfigMapService) UpdateConfigMap(ctx context.Context, namespace string, configMapName string, updates *dto.ConfigMapUpdate) (*dto.ConfigMap, error) {
	args := m.Called(ctx, namespace, configMapName, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ConfigMap), args.Error(1)
}

func (m *MockConfigMapService) DeleteConfigMap(ctx context.Context, namespace string, configMapName string) error {
	args := m.Called(ctx, namespace, configMapName)
	return args.Error(0)
}
