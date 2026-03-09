package mocks

import (
	"context"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	"github.com/stretchr/testify/mock"
)

type MockNetworkPolicyService struct {
	mock.Mock
}

func (m *MockNetworkPolicyService) ListNetworkPolicies(ctx context.Context, namespace string) (*dto.NetworkPolicyList, error) {
	args := m.Called(ctx, namespace)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.NetworkPolicyList), args.Error(1)
}

func (m *MockNetworkPolicyService) GetNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string) (*dto.NetworkPolicy, error) {
	args := m.Called(ctx, namespace, networkPolicyName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.NetworkPolicy), args.Error(1)
}

func (m *MockNetworkPolicyService) CreateNetworkPolicy(ctx context.Context, networkPolicyCreate *dto.NetworkPolicyCreate) (*dto.NetworkPolicy, error) {
	args := m.Called(ctx, networkPolicyCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.NetworkPolicy), args.Error(1)
}

func (m *MockNetworkPolicyService) UpdateNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string, updates *dto.NetworkPolicyUpdate) (*dto.NetworkPolicy, error) {
	args := m.Called(ctx, namespace, networkPolicyName, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.NetworkPolicy), args.Error(1)
}

func (m *MockNetworkPolicyService) DeleteNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string) error {
	args := m.Called(ctx, namespace, networkPolicyName)
	return args.Error(0)
}
