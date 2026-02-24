package dto

type NetworkPolicyList struct {
	NetworkPolicies []NetworkPolicy `json:"networkpolicies"`
}

type NetworkPolicy struct {
	Name         string                `json:"name"`
	Namespace    string                `json:"namespace"`
	PodSelector  map[string]string     `json:"pod_selector"`
	PolicyTypes  []string              `json:"policy_types"`
	IngressRules []NetworkPolicyIngress `json:"ingress_rules,omitempty"`
	EgressRules  []NetworkPolicyEgress  `json:"egress_rules,omitempty"`
	Labels       map[string]string     `json:"labels"`
	CreationDate string                `json:"creation_date"`
}

type NetworkPolicyIngress struct {
	From  []NetworkPolicyPeer `json:"from,omitempty" jsonschema:"description:Source peers allowed"`
	Ports []NetworkPolicyPort `json:"ports,omitempty" jsonschema:"description:Allowed ports"`
}

type NetworkPolicyEgress struct {
	To    []NetworkPolicyPeer `json:"to,omitempty" jsonschema:"description:Destination peers allowed"`
	Ports []NetworkPolicyPort `json:"ports,omitempty" jsonschema:"description:Allowed ports"`
}

type NetworkPolicyPeer struct {
	PodSelector       map[string]string `json:"pod_selector,omitempty" jsonschema:"description:Pod label selector"`
	NamespaceSelector map[string]string `json:"namespace_selector,omitempty" jsonschema:"description:Namespace label selector"`
	IPBlock           *IPBlock          `json:"ip_block,omitempty" jsonschema:"description:IP CIDR block"`
}

type IPBlock struct {
	CIDR   string   `json:"cidr" jsonschema:"required,description:CIDR range (e.g. 10.0.0.0/8)"`
	Except []string `json:"except,omitempty" jsonschema:"description:Exception CIDRs"`
}

type NetworkPolicyPort struct {
	Protocol string `json:"protocol,omitempty" jsonschema:"description:Protocol (TCP/UDP/SCTP)"`
	Port     int32  `json:"port,omitempty" jsonschema:"description:Port number"`
}

type NetworkPolicyCreate struct {
	Name         string                `json:"name" jsonschema:"required,description:Name of the network policy"`
	Namespace    string                `json:"namespace" jsonschema:"required,description:Kubernetes namespace"`
	PodSelector  map[string]string     `json:"pod_selector" jsonschema:"required,description:Pod selector labels"`
	PolicyTypes  []string              `json:"policy_types" jsonschema:"required,description:Policy types (Ingress and/or Egress)"`
	IngressRules []NetworkPolicyIngress `json:"ingress_rules,omitempty" jsonschema:"description:Ingress rules"`
	EgressRules  []NetworkPolicyEgress  `json:"egress_rules,omitempty" jsonschema:"description:Egress rules"`
	Labels       map[string]string     `json:"labels,omitempty" jsonschema:"description:Labels for the network policy"`
}

type NetworkPolicyUpdate struct {
	Labels map[string]string `json:"labels,omitempty" jsonschema:"description:Updated labels"`
}
