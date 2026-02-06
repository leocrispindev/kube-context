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
	Name      string            `json:"name" binding:"required"`
	Namespace string            `json:"namespace" binding:"required"`
	Data      map[string]string `json:"data" binding:"required"`
	Labels    map[string]string `json:"labels,omitempty"`
}

type ConfigMapUpdate struct {
	Data   map[string]string `json:"data,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}
