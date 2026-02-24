package mocks

import (
	"context"

	"k8s-agent-new/internal/core/dto"

	"github.com/stretchr/testify/mock"
)

type MockDeploymentService struct {
	mock.Mock
}

func (m *MockDeploymentService) ListDeployments(ctx context.Context, namespace string) (*dto.DeploymentList, error) {
	args := m.Called(ctx, namespace)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DeploymentList), args.Error(1)
}

func (m *MockDeploymentService) UpdateDeployment(ctx context.Context, namespace string, deploymentName string, updates dto.DeploymentUpdate) error {
	args := m.Called(ctx, namespace, deploymentName, updates)
	return args.Error(0)
}

func (m *MockDeploymentService) GetRolloutStatus(ctx context.Context, namespace string, deploymentName string) (*dto.RolloutStatus, error) {
	args := m.Called(ctx, namespace, deploymentName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RolloutStatus), args.Error(1)
}

func (m *MockDeploymentService) TogglePauseDeployment(ctx context.Context, namespace string, deploymentName string, pause bool) error {
	args := m.Called(ctx, namespace, deploymentName, pause)
	return args.Error(0)
}

func (m *MockDeploymentService) GetDeployment(ctx context.Context, namespace string, deploymentName string) (*dto.Deployment, error) {
	args := m.Called(ctx, namespace, deploymentName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Deployment), args.Error(1)
}

func (m *MockDeploymentService) CreateDeployment(ctx context.Context, deployCreate *dto.DeploymentCreate) (*dto.Deployment, error) {
	args := m.Called(ctx, deployCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Deployment), args.Error(1)
}

func (m *MockDeploymentService) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) error {
	args := m.Called(ctx, namespace, deploymentName)
	return args.Error(0)
}
