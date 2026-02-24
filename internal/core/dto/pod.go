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
	Name      string            `json:"name" jsonschema:"required,description:Name of the pod to create"`
	Namespace string            `json:"namespace" jsonschema:"required,description:Kubernetes namespace for the pod"`
	Image     string            `json:"image" jsonschema:"required,description:Container image (e.g. nginx:latest)"`
	Labels    map[string]string `json:"labels,omitempty" jsonschema:"description:Key-value labels for the pod"`
	Port      int32             `json:"port,omitempty" jsonschema:"description:Container port to expose"`
	Env       map[string]string `json:"env,omitempty" jsonschema:"description:Environment variables as key-value pairs"`
	Command   []string          `json:"command,omitempty" jsonschema:"description:Container command override"`
	Args      []string          `json:"args,omitempty" jsonschema:"description:Arguments for the container command"`
}

type PodUpdate struct {
	Labels map[string]string `json:"labels,omitempty" jsonschema:"description:Updated labels for the pod"`
	Image  string            `json:"image,omitempty" jsonschema:"description:New container image"`
}
