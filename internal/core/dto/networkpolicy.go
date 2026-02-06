package dto

type NetworkPolicyList struct {
	NetworkPolicies []NetworkPolicy `json:"networkpolicies"`
}

type NetworkPolicy struct {
	Name         string                 `json:"name"`
	Namespace    string                 `json:"namespace"`
	PodSelector  map[string]string      `json:"pod_selector"`
	PolicyTypes  []string               `json:"policy_types"`
	IngressRules []NetworkPolicyIngress `json:"ingress_rules,omitempty"`
	EgressRules  []NetworkPolicyEgress  `json:"egress_rules,omitempty"`
	Labels       map[string]string      `json:"labels"`
	CreationDate string                 `json:"creation_date"`
}

type NetworkPolicyIngress struct {
	From  []NetworkPolicyPeer `json:"from,omitempty"`
	Ports []NetworkPolicyPort `json:"ports,omitempty"`
}

type NetworkPolicyEgress struct {
	To    []NetworkPolicyPeer `json:"to,omitempty"`
	Ports []NetworkPolicyPort `json:"ports,omitempty"`
}

type NetworkPolicyPeer struct {
	PodSelector       map[string]string `json:"pod_selector,omitempty"`
	NamespaceSelector map[string]string `json:"namespace_selector,omitempty"`
	IPBlock           *IPBlock          `json:"ip_block,omitempty"`
}

type IPBlock struct {
	CIDR   string   `json:"cidr"`
	Except []string `json:"except,omitempty"`
}

type NetworkPolicyPort struct {
	Protocol string `json:"protocol,omitempty"`
	Port     int32  `json:"port,omitempty"`
}

type NetworkPolicyCreate struct {
	Name         string                 `json:"name" binding:"required"`
	Namespace    string                 `json:"namespace" binding:"required"`
	PodSelector  map[string]string      `json:"pod_selector" binding:"required"`
	PolicyTypes  []string               `json:"policy_types" binding:"required"`
	IngressRules []NetworkPolicyIngress `json:"ingress_rules,omitempty"`
	EgressRules  []NetworkPolicyEgress  `json:"egress_rules,omitempty"`
	Labels       map[string]string      `json:"labels,omitempty"`
}

type NetworkPolicyUpdate struct {
	Labels map[string]string `json:"labels,omitempty"`
}
