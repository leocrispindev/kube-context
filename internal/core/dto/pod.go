package dto

type PodList struct {
	Pods []string `json:"pods"`
}

type PodDetails struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Status    string            `json:"status"`
	PodIP     string            `json:"pod_ip"`
	NodeName  string            `json:"node_name"`
	CreatedAt string            `json:"created_at"`
	Labels    map[string]string `json:"labels"`
	Images    []string          `json:"images"`
	Restarts  int32             `json:"restarts"`
}

type PodCreate struct {
	Name      string            `json:"name" binding:"required"`
	Namespace string            `json:"namespace" binding:"required"`
	Image     string            `json:"image" binding:"required"`
	Labels    map[string]string `json:"labels,omitempty"`
	Port      int32             `json:"port,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	Command   []string          `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
}

type PodUpdate struct {
	Labels map[string]string `json:"labels,omitempty"`
	Image  string            `json:"image,omitempty"`
}
