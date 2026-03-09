package ingress

import (
	"context"
	"fmt"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceInterface interface {
	ListIngresses(ctx context.Context, namespace string) (*dto.IngressList, error)
	GetIngress(ctx context.Context, namespace string, ingressName string) (*dto.Ingress, error)
	CreateIngress(ctx context.Context, ingressCreate *dto.IngressCreate) (*dto.Ingress, error)
	UpdateIngress(ctx context.Context, namespace string, ingressName string, updates *dto.IngressUpdate) (*dto.Ingress, error)
	DeleteIngress(ctx context.Context, namespace string, ingressName string) error
}

type Service struct {
	client kubernetes.Interface
}

func NewService(client kubernetes.Interface) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ListIngresses(ctx context.Context, namespace string) (*dto.IngressList, error) {
	ingresses, err := s.client.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %w", err)
	}

	var dtos []dto.Ingress
	for _, ing := range ingresses.Items {
		dtos = append(dtos, s.convertToDTO(&ing))
	}

	return &dto.IngressList{Ingresses: dtos}, nil
}

func (s *Service) GetIngress(ctx context.Context, namespace string, ingressName string) (*dto.Ingress, error) {
	ing, err := s.client.NetworkingV1().Ingresses(namespace).Get(ctx, ingressName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ingress: %w", err)
	}

	result := s.convertToDTO(ing)
	return &result, nil
}

func (s *Service) CreateIngress(ctx context.Context, ingressCreate *dto.IngressCreate) (*dto.Ingress, error) {
	var rules []networkingv1.IngressRule
	for _, r := range ingressCreate.Rules {
		var paths []networkingv1.HTTPIngressPath
		for _, p := range r.Paths {
			pathType := networkingv1.PathType(p.PathType)
			if p.PathType == "" {
				prefix := networkingv1.PathTypePrefix
				pathType = prefix
			}

			paths = append(paths, networkingv1.HTTPIngressPath{
				Path:     p.Path,
				PathType: &pathType,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: p.ServiceName,
						Port: networkingv1.ServiceBackendPort{
							Number: p.ServicePort,
						},
					},
				},
			})
		}

		rules = append(rules, networkingv1.IngressRule{
			Host: r.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: paths,
				},
			},
		})
	}

	var tls []networkingv1.IngressTLS
	for _, t := range ingressCreate.TLS {
		tls = append(tls, networkingv1.IngressTLS{
			Hosts:      t.Hosts,
			SecretName: t.SecretName,
		})
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ingressCreate.Name,
			Namespace:   ingressCreate.Namespace,
			Labels:      ingressCreate.Labels,
			Annotations: ingressCreate.Annotations,
		},
		Spec: networkingv1.IngressSpec{
			Rules: rules,
			TLS:   tls,
		},
	}

	if ingressCreate.IngressClass != "" {
		ingress.Spec.IngressClassName = &ingressCreate.IngressClass
	}

	created, err := s.client.NetworkingV1().Ingresses(ingressCreate.Namespace).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create ingress: %w", err)
	}

	return s.GetIngress(ctx, created.Namespace, created.Name)
}

func (s *Service) UpdateIngress(ctx context.Context, namespace string, ingressName string, updates *dto.IngressUpdate) (*dto.Ingress, error) {
	ingress, err := s.client.NetworkingV1().Ingresses(namespace).Get(ctx, ingressName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ingress: %w", err)
	}

	if updates.IngressClass != "" {
		ingress.Spec.IngressClassName = &updates.IngressClass
	}

	if len(updates.Labels) > 0 {
		if ingress.Labels == nil {
			ingress.Labels = make(map[string]string)
		}
		for k, v := range updates.Labels {
			ingress.Labels[k] = v
		}
	}

	if len(updates.Annotations) > 0 {
		if ingress.Annotations == nil {
			ingress.Annotations = make(map[string]string)
		}
		for k, v := range updates.Annotations {
			ingress.Annotations[k] = v
		}
	}

	_, err = s.client.NetworkingV1().Ingresses(namespace).Update(ctx, ingress, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update ingress: %w", err)
	}

	return s.GetIngress(ctx, namespace, ingressName)
}

func (s *Service) DeleteIngress(ctx context.Context, namespace string, ingressName string) error {
	err := s.client.NetworkingV1().Ingresses(namespace).Delete(ctx, ingressName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete ingress: %w", err)
	}
	return nil
}

func (s *Service) convertToDTO(ing *networkingv1.Ingress) dto.Ingress {
	var rules []dto.IngressRule
	for _, r := range ing.Spec.Rules {
		var paths []dto.IngressPath
		if r.HTTP != nil {
			for _, p := range r.HTTP.Paths {
				pathType := string(*p.PathType)
				paths = append(paths, dto.IngressPath{
					Path:        p.Path,
					PathType:    pathType,
					ServiceName: p.Backend.Service.Name,
					ServicePort: p.Backend.Service.Port.Number,
				})
			}
		}

		rules = append(rules, dto.IngressRule{
			Host:  r.Host,
			Paths: paths,
		})
	}

	var tls []dto.IngressTLS
	for _, t := range ing.Spec.TLS {
		tls = append(tls, dto.IngressTLS{
			Hosts:      t.Hosts,
			SecretName: t.SecretName,
		})
	}

	ingressClass := ""
	if ing.Spec.IngressClassName != nil {
		ingressClass = *ing.Spec.IngressClassName
	}

	return dto.Ingress{
		Name:         ing.Name,
		Namespace:    ing.Namespace,
		IngressClass: ingressClass,
		Rules:        rules,
		TLS:          tls,
		Labels:       ing.Labels,
		Annotations:  ing.Annotations,
		CreationDate: ing.CreationTimestamp.String(),
	}
}
