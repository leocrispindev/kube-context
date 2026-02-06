package service

import (
	"context"
	"fmt"

	"k8s-agent-new/internal/core/dto"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	client *kubernetes.Clientset
}

func NewService(client *kubernetes.Clientset) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ListServices(ctx context.Context, namespace string) (*dto.ServiceList, error) {
	services, err := s.client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var dtos []dto.Service
	for _, svc := range services.Items {
		var ports []dto.ServicePort
		for _, p := range svc.Spec.Ports {
			ports = append(ports, dto.ServicePort{
				Name:       p.Name,
				Port:       p.Port,
				TargetPort: p.TargetPort.IntVal,
				NodePort:   p.NodePort,
				Protocol:   string(p.Protocol),
			})
		}

		dtos = append(dtos, dto.Service{
			Name:         svc.Name,
			Namespace:    svc.Namespace,
			Type:         string(svc.Spec.Type),
			ClusterIP:    svc.Spec.ClusterIP,
			ExternalIPs:  svc.Spec.ExternalIPs,
			Ports:        ports,
			Selector:     svc.Spec.Selector,
			Labels:       svc.Labels,
			CreationDate: svc.CreationTimestamp.String(),
		})
	}

	return &dto.ServiceList{Services: dtos}, nil
}

func (s *Service) GetService(ctx context.Context, namespace string, serviceName string) (*dto.Service, error) {
	svc, err := s.client.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	var ports []dto.ServicePort
	for _, p := range svc.Spec.Ports {
		ports = append(ports, dto.ServicePort{
			Name:       p.Name,
			Port:       p.Port,
			TargetPort: p.TargetPort.IntVal,
			NodePort:   p.NodePort,
			Protocol:   string(p.Protocol),
		})
	}

	return &dto.Service{
		Name:         svc.Name,
		Namespace:    svc.Namespace,
		Type:         string(svc.Spec.Type),
		ClusterIP:    svc.Spec.ClusterIP,
		ExternalIPs:  svc.Spec.ExternalIPs,
		Ports:        ports,
		Selector:     svc.Spec.Selector,
		Labels:       svc.Labels,
		CreationDate: svc.CreationTimestamp.String(),
	}, nil
}

func (s *Service) CreateService(ctx context.Context, serviceCreate *dto.ServiceCreate) (*dto.Service, error) {
	var ports []v1.ServicePort
	for _, p := range serviceCreate.Ports {
		protocol := v1.ProtocolTCP
		if p.Protocol != "" {
			protocol = v1.Protocol(p.Protocol)
		}

		ports = append(ports, v1.ServicePort{
			Name:       p.Name,
			Port:       p.Port,
			TargetPort: intstr.FromInt(int(p.TargetPort)),
			Protocol:   protocol,
		})
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceCreate.Name,
			Namespace: serviceCreate.Namespace,
			Labels:    serviceCreate.Labels,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceType(serviceCreate.Type),
			Selector: serviceCreate.Selector,
			Ports:    ports,
		},
	}

	created, err := s.client.CoreV1().Services(serviceCreate.Namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return s.GetService(ctx, created.Namespace, created.Name)
}

func (s *Service) UpdateService(ctx context.Context, namespace string, serviceName string, updates *dto.ServiceUpdate) (*dto.Service, error) {
	service, err := s.client.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	if updates.Type != "" {
		service.Spec.Type = v1.ServiceType(updates.Type)
	}

	if len(updates.Selector) > 0 {
		if service.Spec.Selector == nil {
			service.Spec.Selector = make(map[string]string)
		}
		for k, v := range updates.Selector {
			service.Spec.Selector[k] = v
		}
	}

	if len(updates.Labels) > 0 {
		if service.Labels == nil {
			service.Labels = make(map[string]string)
		}
		for k, v := range updates.Labels {
			service.Labels[k] = v
		}
	}

	_, err = s.client.CoreV1().Services(namespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	return s.GetService(ctx, namespace, serviceName)
}

func (s *Service) DeleteService(ctx context.Context, namespace string, serviceName string) error {
	err := s.client.CoreV1().Services(namespace).Delete(ctx, serviceName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}
