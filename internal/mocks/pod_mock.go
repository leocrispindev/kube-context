package mocks

import (
	"context"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	"github.com/stretchr/testify/mock"
)

type MockPodService struct {
	mock.Mock
}

func (m *MockPodService) ListPods(ctx context.Context, namespace string) (*dto.PodList, error) {
	args := m.Called(ctx, namespace)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PodList), args.Error(1)
}

func (m *MockPodService) ScalePod(ctx context.Context, namespace string, podName string, replicas int32) (*dto.PodList, error) {
	args := m.Called(ctx, namespace, podName, replicas)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PodList), args.Error(1)
}

func (m *MockPodService) GetPod(ctx context.Context, namespace string, podName string) (*dto.PodDetails, error) {
	args := m.Called(ctx, namespace, podName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PodDetails), args.Error(1)
}

func (m *MockPodService) DeletePod(ctx context.Context, namespace string, podName string) error {
	args := m.Called(ctx, namespace, podName)
	return args.Error(0)
}

func (m *MockPodService) RestartPod(ctx context.Context, namespace string, podName string) error {
	args := m.Called(ctx, namespace, podName)
	return args.Error(0)
}

func (m *MockPodService) CreatePod(ctx context.Context, podCreate *dto.PodCreate) (*dto.PodDetails, error) {
	args := m.Called(ctx, podCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PodDetails), args.Error(1)
}

func (m *MockPodService) UpdatePod(ctx context.Context, namespace string, podName string, updates *dto.PodUpdate) (*dto.PodDetails, error) {
	args := m.Called(ctx, namespace, podName, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PodDetails), args.Error(1)
}
