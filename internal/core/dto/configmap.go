package dto

type ConfigMapList struct {
	ConfigMaps []ConfigMap `json:"configmaps"`
}

type ConfigMap struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	Data         map[string]string `json:"data"`
	Labels       map[string]string `json:"labels"`
	CreationDate string            `json:"creation_date"`
}

type ConfigMapCreate struct {
	Name      string            `json:"name" jsonschema:"required,description:Name of the configmap"`
	Namespace string            `json:"namespace" jsonschema:"required,description:Kubernetes namespace"`
	Data      map[string]string `json:"data" jsonschema:"required,description:Key-value data entries for the configmap"`
	Labels    map[string]string `json:"labels,omitempty" jsonschema:"description:Labels for the configmap"`
}

type ConfigMapUpdate struct {
	Data   map[string]string `json:"data,omitempty" jsonschema:"description:Updated data entries"`
	Labels map[string]string `json:"labels,omitempty" jsonschema:"description:Updated labels"`
}
