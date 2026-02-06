package dto

type ServiceList struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	Type         string            `json:"type"`
	ClusterIP    string            `json:"cluster_ip"`
	ExternalIPs  []string          `json:"external_ips,omitempty"`
	Ports        []ServicePort     `json:"ports"`
	Selector     map[string]string `json:"selector"`
	Labels       map[string]string `json:"labels"`
	CreationDate string            `json:"creation_date"`
}

type ServicePort struct {
	Name       string `json:"name,omitempty"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"target_port"`
	NodePort   int32  `json:"node_port,omitempty"`
	Protocol   string `json:"protocol"`
}

type ServiceCreate struct {
	Name      string            `json:"name" binding:"required"`
	Namespace string            `json:"namespace" binding:"required"`
	Type      string            `json:"type" binding:"required"`
	Selector  map[string]string `json:"selector" binding:"required"`
	Ports     []ServicePort     `json:"ports" binding:"required"`
	Labels    map[string]string `json:"labels,omitempty"`
}

type ServiceUpdate struct {
	Type     string            `json:"type,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
}
