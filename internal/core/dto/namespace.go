package dto

type NamespaceList struct {
	Namespaces []Namespace `json:"namespaces"`
}

type Namespace struct {
	Name        string            `json:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Status      string            `json:"status,omitempty"`
}
