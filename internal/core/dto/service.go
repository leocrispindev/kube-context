package dto

type ServiceList struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name                string            `json:"name"`
	Namespace           string            `json:"namespace"`
	Type                string            `json:"type"`
	ClusterIP           string            `json:"cluster_ip"`
	ExternalIPs         []string          `json:"external_ips,omitempty"`
	LoadBalancerIngress []LBIngress       `json:"load_balancer_ingress,omitempty"`
	Ports               []ServicePort     `json:"ports"`
	Selector            map[string]string `json:"selector"`
	Labels              map[string]string `json:"labels"`
	CreationDate        string            `json:"creation_date"`
}

type LBIngress struct {
	IP       string `json:"ip,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

type ServicePort struct {
	Name       string `json:"name,omitempty" jsonschema:"description:Port name"`
	Port       int32  `json:"port" jsonschema:"required,description:Service port number"`
	TargetPort int32  `json:"target_port" jsonschema:"required,description:Target port on the pod"`
	NodePort   int32  `json:"node_port,omitempty" jsonschema:"description:Node port (for NodePort/LoadBalancer types)"`
	Protocol   string `json:"protocol" jsonschema:"required,description:Protocol (TCP or UDP)"`
}

type ServiceCreate struct {
	Name      string            `json:"name" jsonschema:"required,description:Name of the service"`
	Namespace string            `json:"namespace" jsonschema:"required,description:Kubernetes namespace"`
	Type      string            `json:"type" jsonschema:"required,description:Service type (ClusterIP/NodePort/LoadBalancer)"`
	Selector  map[string]string `json:"selector" jsonschema:"required,description:Pod selector labels"`
	Ports     []ServicePort     `json:"ports" jsonschema:"required,description:List of service ports"`
	Labels    map[string]string `json:"labels,omitempty" jsonschema:"description:Labels for the service"`
}

type ServiceUpdate struct {
	Type     string            `json:"type,omitempty" jsonschema:"description:Updated service type"`
	Selector map[string]string `json:"selector,omitempty" jsonschema:"description:Updated pod selector labels"`
	Labels   map[string]string `json:"labels,omitempty" jsonschema:"description:Updated labels"`
}
