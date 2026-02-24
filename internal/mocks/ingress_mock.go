package mocks

import (
	"context"

	"k8s-agent-new/internal/core/dto"

	"github.com/stretchr/testify/mock"
)

type MockIngressService struct {
	mock.Mock
}

func (m *MockIngressService) ListIngresses(ctx context.Context, namespace string) (*dto.IngressList, error) {
	args := m.Called(ctx, namespace)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.IngressList), args.Error(1)
}

func (m *MockIngressService) GetIngress(ctx context.Context, namespace string, ingressName string) (*dto.Ingress, error) {
	args := m.Called(ctx, namespace, ingressName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Ingress), args.Error(1)
}

func (m *MockIngressService) CreateIngress(ctx context.Context, ingressCreate *dto.IngressCreate) (*dto.Ingress, error) {
	args := m.Called(ctx, ingressCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Ingress), args.Error(1)
}

func (m *MockIngressService) UpdateIngress(ctx context.Context, namespace string, ingressName string, updates *dto.IngressUpdate) (*dto.Ingress, error) {
	args := m.Called(ctx, namespace, ingressName, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Ingress), args.Error(1)
}

func (m *MockIngressService) DeleteIngress(ctx context.Context, namespace string, ingressName string) error {
	args := m.Called(ctx, namespace, ingressName)
	return args.Error(0)
}
