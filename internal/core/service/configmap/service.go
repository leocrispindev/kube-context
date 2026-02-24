package configmap

import (
	"context"
	"fmt"

	"k8s-agent-new/internal/core/dto"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceInterface interface {
	ListConfigMaps(ctx context.Context, namespace string) (*dto.ConfigMapList, error)
	GetConfigMap(ctx context.Context, namespace string, configMapName string) (*dto.ConfigMap, error)
	CreateConfigMap(ctx context.Context, configMapCreate *dto.ConfigMapCreate) (*dto.ConfigMap, error)
	UpdateConfigMap(ctx context.Context, namespace string, configMapName string, updates *dto.ConfigMapUpdate) (*dto.ConfigMap, error)
	DeleteConfigMap(ctx context.Context, namespace string, configMapName string) error
}

type Service struct {
	client kubernetes.Interface
}

func NewService(client kubernetes.Interface) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ListConfigMaps(ctx context.Context, namespace string) (*dto.ConfigMapList, error) {
	configMaps, err := s.client.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list configmaps: %w", err)
	}

	var dtos []dto.ConfigMap
	for _, cm := range configMaps.Items {
		dtos = append(dtos, dto.ConfigMap{
			Name:         cm.Name,
			Namespace:    cm.Namespace,
			Data:         cm.Data,
			Labels:       cm.Labels,
			CreationDate: cm.CreationTimestamp.String(),
		})
	}

	return &dto.ConfigMapList{ConfigMaps: dtos}, nil
}

func (s *Service) GetConfigMap(ctx context.Context, namespace string, configMapName string) (*dto.ConfigMap, error) {
	cm, err := s.client.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configmap: %w", err)
	}

	return &dto.ConfigMap{
		Name:         cm.Name,
		Namespace:    cm.Namespace,
		Data:         cm.Data,
		Labels:       cm.Labels,
		CreationDate: cm.CreationTimestamp.String(),
	}, nil
}

func (s *Service) CreateConfigMap(ctx context.Context, configMapCreate *dto.ConfigMapCreate) (*dto.ConfigMap, error) {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapCreate.Name,
			Namespace: configMapCreate.Namespace,
			Labels:    configMapCreate.Labels,
		},
		Data: configMapCreate.Data,
	}

	created, err := s.client.CoreV1().ConfigMaps(configMapCreate.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create configmap: %w", err)
	}

	return s.GetConfigMap(ctx, created.Namespace, created.Name)
}

func (s *Service) UpdateConfigMap(ctx context.Context, namespace string, configMapName string, updates *dto.ConfigMapUpdate) (*dto.ConfigMap, error) {
	configMap, err := s.client.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configmap: %w", err)
	}

	if len(updates.Data) > 0 {
		if configMap.Data == nil {
			configMap.Data = make(map[string]string)
		}
		for k, v := range updates.Data {
			configMap.Data[k] = v
		}
	}

	if len(updates.Labels) > 0 {
		if configMap.Labels == nil {
			configMap.Labels = make(map[string]string)
		}
		for k, v := range updates.Labels {
			configMap.Labels[k] = v
		}
	}

	_, err = s.client.CoreV1().ConfigMaps(namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update configmap: %w", err)
	}

	return s.GetConfigMap(ctx, namespace, configMapName)
}

func (s *Service) DeleteConfigMap(ctx context.Context, namespace string, configMapName string) error {
	err := s.client.CoreV1().ConfigMaps(namespace).Delete(ctx, configMapName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete configmap: %w", err)
	}
	return nil
}
