package mocks

import (
	"context"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	"github.com/stretchr/testify/mock"
)

type MockServiceService struct {
	mock.Mock
}

func (m *MockServiceService) ListServices(ctx context.Context, namespace string) (*dto.ServiceList, error) {
	args := m.Called(ctx, namespace)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ServiceList), args.Error(1)
}

func (m *MockServiceService) GetService(ctx context.Context, namespace string, serviceName string) (*dto.Service, error) {
	args := m.Called(ctx, namespace, serviceName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Service), args.Error(1)
}

func (m *MockServiceService) CreateService(ctx context.Context, serviceCreate *dto.ServiceCreate) (*dto.Service, error) {
	args := m.Called(ctx, serviceCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Service), args.Error(1)
}

func (m *MockServiceService) UpdateService(ctx context.Context, namespace string, serviceName string, updates *dto.ServiceUpdate) (*dto.Service, error) {
	args := m.Called(ctx, namespace, serviceName, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Service), args.Error(1)
}

func (m *MockServiceService) DeleteService(ctx context.Context, namespace string, serviceName string) error {
	args := m.Called(ctx, namespace, serviceName)
	return args.Error(0)
}
