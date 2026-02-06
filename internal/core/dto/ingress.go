package dto

type IngressList struct {
	Ingresses []Ingress `json:"ingresses"`
}

type Ingress struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	IngressClass string            `json:"ingress_class,omitempty"`
	Rules        []IngressRule     `json:"rules"`
	TLS          []IngressTLS      `json:"tls,omitempty"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	CreationDate string            `json:"creation_date"`
}

type IngressRule struct {
	Host  string        `json:"host"`
	Paths []IngressPath `json:"paths"`
}

type IngressPath struct {
	Path        string `json:"path"`
	PathType    string `json:"path_type"`
	ServiceName string `json:"service_name"`
	ServicePort int32  `json:"service_port"`
}

type IngressTLS struct {
	Hosts      []string `json:"hosts"`
	SecretName string   `json:"secret_name"`
}

type IngressCreate struct {
	Name         string            `json:"name" binding:"required"`
	Namespace    string            `json:"namespace" binding:"required"`
	IngressClass string            `json:"ingress_class,omitempty"`
	Rules        []IngressRule     `json:"rules" binding:"required"`
	TLS          []IngressTLS      `json:"tls,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
}

type IngressUpdate struct {
	IngressClass string            `json:"ingress_class,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
}
