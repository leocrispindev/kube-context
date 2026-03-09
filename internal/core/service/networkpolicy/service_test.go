package networkpolicy

import (
	"context"
	"testing"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type NetworkPolicyServiceTestSuite struct {
	suite.Suite
	client *fake.Clientset
	svc    *Service
}

func (s *NetworkPolicyServiceTestSuite) SetupTest() {
	s.client = fake.NewSimpleClientset()
	s.svc = NewService(s.client)
}

func TestNetworkPolicyServiceSuite(t *testing.T) {
	suite.Run(t, new(NetworkPolicyServiceTestSuite))
}

func (s *NetworkPolicyServiceTestSuite) seedNetworkPolicy(name, namespace string) {
	protocol := v1.ProtocolTCP
	port := intstr.FromInt32(5432)
	np := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": "web"},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{"role": "db"},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
				networkingv1.PolicyTypeEgress,
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"role": "frontend"},
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{Protocol: &protocol, Port: &port},
					},
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					To: []networkingv1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"env": "prod"},
							},
						},
					},
					Ports: []networkingv1.NetworkPolicyPort{
						{Protocol: &protocol, Port: &port},
					},
				},
			},
		},
	}
	_, err := s.client.NetworkingV1().NetworkPolicies(namespace).Create(context.Background(), np, metav1.CreateOptions{})
	s.Require().NoError(err)
}

// --- ListNetworkPolicies ---

func (s *NetworkPolicyServiceTestSuite) TestListNetworkPolicies_Success() {
	s.seedNetworkPolicy("np-1", "default")
	s.seedNetworkPolicy("np-2", "default")

	result, err := s.svc.ListNetworkPolicies(context.Background(), "default")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.NetworkPolicies, 2)
	assert.Equal(s.T(), "np-1", result.NetworkPolicies[0].Name)
	assert.Equal(s.T(), "np-2", result.NetworkPolicies[1].Name)
}

func (s *NetworkPolicyServiceTestSuite) TestListNetworkPolicies_EmptyNamespace() {
	result, err := s.svc.ListNetworkPolicies(context.Background(), "empty-ns")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.NetworkPolicies)
}

func (s *NetworkPolicyServiceTestSuite) TestListNetworkPolicies_Error() {
	s.client.PrependReactor("list", "networkpolicies", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, assert.AnError
	})

	result, err := s.svc.ListNetworkPolicies(context.Background(), "default")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to list network policies")
}

// --- GetNetworkPolicy ---

func (s *NetworkPolicyServiceTestSuite) TestGetNetworkPolicy_Success() {
	s.seedNetworkPolicy("my-np", "default")

	result, err := s.svc.GetNetworkPolicy(context.Background(), "default", "my-np")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "my-np", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), map[string]string{"role": "db"}, result.PodSelector)
	assert.Equal(s.T(), []string{"Ingress", "Egress"}, result.PolicyTypes)
	assert.Len(s.T(), result.IngressRules, 1)
	assert.Len(s.T(), result.IngressRules[0].From, 1)
	assert.Equal(s.T(), map[string]string{"role": "frontend"}, result.IngressRules[0].From[0].PodSelector)
	assert.Len(s.T(), result.IngressRules[0].Ports, 1)
	assert.Equal(s.T(), "TCP", result.IngressRules[0].Ports[0].Protocol)
	assert.Equal(s.T(), int32(5432), result.IngressRules[0].Ports[0].Port)
	assert.Len(s.T(), result.EgressRules, 1)
	assert.Len(s.T(), result.EgressRules[0].To, 1)
	assert.Equal(s.T(), map[string]string{"env": "prod"}, result.EgressRules[0].To[0].NamespaceSelector)
	assert.Equal(s.T(), map[string]string{"app": "web"}, result.Labels)
}

func (s *NetworkPolicyServiceTestSuite) TestGetNetworkPolicy_NotFound() {
	result, err := s.svc.GetNetworkPolicy(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get network policy")
}

// --- CreateNetworkPolicy ---

func (s *NetworkPolicyServiceTestSuite) TestCreateNetworkPolicy_Success() {
	create := &dto.NetworkPolicyCreate{
		Name:        "new-np",
		Namespace:   "default",
		PodSelector: map[string]string{"app": "api"},
		PolicyTypes: []string{"Ingress"},
		IngressRules: []dto.NetworkPolicyIngress{
			{
				From: []dto.NetworkPolicyPeer{
					{PodSelector: map[string]string{"role": "web"}},
				},
			},
		},
		Labels: map[string]string{"env": "test"},
	}

	result, err := s.svc.CreateNetworkPolicy(context.Background(), create)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "new-np", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), map[string]string{"app": "api"}, result.PodSelector)
	assert.Equal(s.T(), []string{"Ingress"}, result.PolicyTypes)
	assert.Len(s.T(), result.IngressRules, 1)
	assert.Len(s.T(), result.IngressRules[0].From, 1)
	assert.Equal(s.T(), map[string]string{"role": "web"}, result.IngressRules[0].From[0].PodSelector)
	assert.Equal(s.T(), map[string]string{"env": "test"}, result.Labels)
}

func (s *NetworkPolicyServiceTestSuite) TestCreateNetworkPolicy_Error() {
	s.client.PrependReactor("create", "networkpolicies", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, assert.AnError
	})

	create := &dto.NetworkPolicyCreate{
		Name:      "fail-np",
		Namespace: "default",
	}

	result, err := s.svc.CreateNetworkPolicy(context.Background(), create)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to create network policy")
}

// --- UpdateNetworkPolicy ---

func (s *NetworkPolicyServiceTestSuite) TestUpdateNetworkPolicy_Success() {
	s.seedNetworkPolicy("update-np", "default")

	updates := &dto.NetworkPolicyUpdate{
		Labels: map[string]string{"version": "v2"},
	}

	result, err := s.svc.UpdateNetworkPolicy(context.Background(), "default", "update-np", updates)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "v2", result.Labels["version"])
	assert.Equal(s.T(), "web", result.Labels["app"])
}

func (s *NetworkPolicyServiceTestSuite) TestUpdateNetworkPolicy_NotFound() {
	updates := &dto.NetworkPolicyUpdate{
		Labels: map[string]string{"version": "v2"},
	}

	result, err := s.svc.UpdateNetworkPolicy(context.Background(), "default", "nonexistent", updates)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get network policy")
}

// --- DeleteNetworkPolicy ---

func (s *NetworkPolicyServiceTestSuite) TestDeleteNetworkPolicy_Success() {
	s.seedNetworkPolicy("delete-np", "default")

	err := s.svc.DeleteNetworkPolicy(context.Background(), "default", "delete-np")

	assert.NoError(s.T(), err)

	_, getErr := s.svc.GetNetworkPolicy(context.Background(), "default", "delete-np")
	assert.Error(s.T(), getErr)
}

func (s *NetworkPolicyServiceTestSuite) TestDeleteNetworkPolicy_NotFound() {
	err := s.svc.DeleteNetworkPolicy(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to delete network policy")
}
