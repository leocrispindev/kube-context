package networkpolicy

import (
	"context"
	"fmt"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceInterface interface {
	ListNetworkPolicies(ctx context.Context, namespace string) (*dto.NetworkPolicyList, error)
	GetNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string) (*dto.NetworkPolicy, error)
	CreateNetworkPolicy(ctx context.Context, networkPolicyCreate *dto.NetworkPolicyCreate) (*dto.NetworkPolicy, error)
	UpdateNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string, updates *dto.NetworkPolicyUpdate) (*dto.NetworkPolicy, error)
	DeleteNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string) error
}

type Service struct {
	client kubernetes.Interface
}

func NewService(client kubernetes.Interface) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ListNetworkPolicies(ctx context.Context, namespace string) (*dto.NetworkPolicyList, error) {
	netpols, err := s.client.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list network policies: %w", err)
	}

	var dtos []dto.NetworkPolicy
	for _, np := range netpols.Items {
		dtos = append(dtos, s.convertToDTO(&np))
	}

	return &dto.NetworkPolicyList{NetworkPolicies: dtos}, nil
}

func (s *Service) GetNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string) (*dto.NetworkPolicy, error) {
	np, err := s.client.NetworkingV1().NetworkPolicies(namespace).Get(ctx, networkPolicyName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get network policy: %w", err)
	}

	result := s.convertToDTO(np)
	return &result, nil
}

func (s *Service) CreateNetworkPolicy(ctx context.Context, networkPolicyCreate *dto.NetworkPolicyCreate) (*dto.NetworkPolicy, error) {
	var policyTypes []networkingv1.PolicyType
	for _, pt := range networkPolicyCreate.PolicyTypes {
		policyTypes = append(policyTypes, networkingv1.PolicyType(pt))
	}

	var ingressRules []networkingv1.NetworkPolicyIngressRule
	for _, rule := range networkPolicyCreate.IngressRules {
		ingressRules = append(ingressRules, s.convertIngressRule(rule))
	}

	var egressRules []networkingv1.NetworkPolicyEgressRule
	for _, rule := range networkPolicyCreate.EgressRules {
		egressRules = append(egressRules, s.convertEgressRule(rule))
	}

	networkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      networkPolicyCreate.Name,
			Namespace: networkPolicyCreate.Namespace,
			Labels:    networkPolicyCreate.Labels,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: networkPolicyCreate.PodSelector,
			},
			PolicyTypes: policyTypes,
			Ingress:     ingressRules,
			Egress:      egressRules,
		},
	}

	created, err := s.client.NetworkingV1().NetworkPolicies(networkPolicyCreate.Namespace).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create network policy: %w", err)
	}

	return s.GetNetworkPolicy(ctx, created.Namespace, created.Name)
}

func (s *Service) UpdateNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string, updates *dto.NetworkPolicyUpdate) (*dto.NetworkPolicy, error) {
	networkPolicy, err := s.client.NetworkingV1().NetworkPolicies(namespace).Get(ctx, networkPolicyName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get network policy: %w", err)
	}

	if len(updates.Labels) > 0 {
		if networkPolicy.Labels == nil {
			networkPolicy.Labels = make(map[string]string)
		}
		for k, v := range updates.Labels {
			networkPolicy.Labels[k] = v
		}
	}

	_, err = s.client.NetworkingV1().NetworkPolicies(namespace).Update(ctx, networkPolicy, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update network policy: %w", err)
	}

	return s.GetNetworkPolicy(ctx, namespace, networkPolicyName)
}

func (s *Service) DeleteNetworkPolicy(ctx context.Context, namespace string, networkPolicyName string) error {
	err := s.client.NetworkingV1().NetworkPolicies(namespace).Delete(ctx, networkPolicyName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete network policy: %w", err)
	}
	return nil
}

func (s *Service) convertToDTO(np *networkingv1.NetworkPolicy) dto.NetworkPolicy {
	var policyTypes []string
	for _, pt := range np.Spec.PolicyTypes {
		policyTypes = append(policyTypes, string(pt))
	}

	var ingressRules []dto.NetworkPolicyIngress
	for _, rule := range np.Spec.Ingress {
		var from []dto.NetworkPolicyPeer
		for _, peer := range rule.From {
			from = append(from, s.convertPeerToDTO(peer))
		}

		var ports []dto.NetworkPolicyPort
		for _, port := range rule.Ports {
			ports = append(ports, s.convertPortToDTO(port))
		}

		ingressRules = append(ingressRules, dto.NetworkPolicyIngress{
			From:  from,
			Ports: ports,
		})
	}

	var egressRules []dto.NetworkPolicyEgress
	for _, rule := range np.Spec.Egress {
		var to []dto.NetworkPolicyPeer
		for _, peer := range rule.To {
			to = append(to, s.convertPeerToDTO(peer))
		}

		var ports []dto.NetworkPolicyPort
		for _, port := range rule.Ports {
			ports = append(ports, s.convertPortToDTO(port))
		}

		egressRules = append(egressRules, dto.NetworkPolicyEgress{
			To:    to,
			Ports: ports,
		})
	}

	return dto.NetworkPolicy{
		Name:         np.Name,
		Namespace:    np.Namespace,
		PodSelector:  np.Spec.PodSelector.MatchLabels,
		PolicyTypes:  policyTypes,
		IngressRules: ingressRules,
		EgressRules:  egressRules,
		Labels:       np.Labels,
		CreationDate: np.CreationTimestamp.String(),
	}
}

func (s *Service) convertPeerToDTO(peer networkingv1.NetworkPolicyPeer) dto.NetworkPolicyPeer {
	result := dto.NetworkPolicyPeer{}

	if peer.PodSelector != nil {
		result.PodSelector = peer.PodSelector.MatchLabels
	}

	if peer.NamespaceSelector != nil {
		result.NamespaceSelector = peer.NamespaceSelector.MatchLabels
	}

	if peer.IPBlock != nil {
		result.IPBlock = &dto.IPBlock{
			CIDR:   peer.IPBlock.CIDR,
			Except: peer.IPBlock.Except,
		}
	}

	return result
}

func (s *Service) convertPortToDTO(port networkingv1.NetworkPolicyPort) dto.NetworkPolicyPort {
	result := dto.NetworkPolicyPort{}

	if port.Protocol != nil {
		result.Protocol = string(*port.Protocol)
	}

	if port.Port != nil {
		result.Port = port.Port.IntVal
	}

	return result
}

func (s *Service) convertIngressRule(rule dto.NetworkPolicyIngress) networkingv1.NetworkPolicyIngressRule {
	var from []networkingv1.NetworkPolicyPeer
	for _, peer := range rule.From {
		from = append(from, s.convertPeerFromDTO(peer))
	}

	var ports []networkingv1.NetworkPolicyPort
	for _, port := range rule.Ports {
		ports = append(ports, s.convertPortFromDTO(port))
	}

	return networkingv1.NetworkPolicyIngressRule{
		From:  from,
		Ports: ports,
	}
}

func (s *Service) convertEgressRule(rule dto.NetworkPolicyEgress) networkingv1.NetworkPolicyEgressRule {
	var to []networkingv1.NetworkPolicyPeer
	for _, peer := range rule.To {
		to = append(to, s.convertPeerFromDTO(peer))
	}

	var ports []networkingv1.NetworkPolicyPort
	for _, port := range rule.Ports {
		ports = append(ports, s.convertPortFromDTO(port))
	}

	return networkingv1.NetworkPolicyEgressRule{
		To:    to,
		Ports: ports,
	}
}

func (s *Service) convertPeerFromDTO(peer dto.NetworkPolicyPeer) networkingv1.NetworkPolicyPeer {
	result := networkingv1.NetworkPolicyPeer{}

	if len(peer.PodSelector) > 0 {
		result.PodSelector = &metav1.LabelSelector{
			MatchLabels: peer.PodSelector,
		}
	}

	if len(peer.NamespaceSelector) > 0 {
		result.NamespaceSelector = &metav1.LabelSelector{
			MatchLabels: peer.NamespaceSelector,
		}
	}

	if peer.IPBlock != nil {
		result.IPBlock = &networkingv1.IPBlock{
			CIDR:   peer.IPBlock.CIDR,
			Except: peer.IPBlock.Except,
		}
	}

	return result
}

func (s *Service) convertPortFromDTO(port dto.NetworkPolicyPort) networkingv1.NetworkPolicyPort {
	return networkingv1.NetworkPolicyPort{}
}
