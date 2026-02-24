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
	Host  string        `json:"host" jsonschema:"required,description:Hostname for the rule"`
	Paths []IngressPath `json:"paths" jsonschema:"required,description:Paths for this host"`
}

type IngressPath struct {
	Path        string `json:"path" jsonschema:"required,description:URL path (e.g. /)"`
	PathType    string `json:"path_type" jsonschema:"required,description:Path type (Prefix/Exact/ImplementationSpecific)"`
	ServiceName string `json:"service_name" jsonschema:"required,description:Backend service name"`
	ServicePort int32  `json:"service_port" jsonschema:"required,description:Backend service port"`
}

type IngressTLS struct {
	Hosts      []string `json:"hosts" jsonschema:"required,description:TLS hostnames"`
	SecretName string   `json:"secret_name" jsonschema:"required,description:Name of the TLS secret"`
}

type IngressCreate struct {
	Name         string            `json:"name" jsonschema:"required,description:Name of the ingress"`
	Namespace    string            `json:"namespace" jsonschema:"required,description:Kubernetes namespace"`
	IngressClass string            `json:"ingress_class,omitempty" jsonschema:"description:Ingress class name"`
	Rules        []IngressRule     `json:"rules" jsonschema:"required,description:Ingress routing rules"`
	TLS          []IngressTLS      `json:"tls,omitempty" jsonschema:"description:TLS configuration"`
	Labels       map[string]string `json:"labels,omitempty" jsonschema:"description:Labels for the ingress"`
	Annotations  map[string]string `json:"annotations,omitempty" jsonschema:"description:Annotations for the ingress"`
}

type IngressUpdate struct {
	IngressClass string            `json:"ingress_class,omitempty" jsonschema:"description:Updated ingress class"`
	Labels       map[string]string `json:"labels,omitempty" jsonschema:"description:Updated labels"`
	Annotations  map[string]string `json:"annotations,omitempty" jsonschema:"description:Updated annotations"`
}
