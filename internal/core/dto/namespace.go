package dto

import v1 "k8s.io/api/core/v1"

type NamespaceList struct {
	Namespaces []Namespace `json:"namespaces"`
}

type Namespace struct {
	Name        string             `json:"name,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Status      v1.NamespaceStatus `json:"status,omitempty"`
}
