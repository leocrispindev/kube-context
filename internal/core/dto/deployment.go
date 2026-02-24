package dto

type DeploymentList struct {
	Deployments []Deployment `json:"deployments"`
}

type Deployment struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	Replicas        int32             `json:"replicas"`
	ReadyReplicas   int32             `json:"ready_replicas"`
	UpdatedReplicas int32             `json:"updated_replicas"`
	Strategy        string            `json:"strategy"`
	Images          []string          `json:"images"`
	Labels          map[string]string `json:"labels"`
	CreationDate    string            `json:"creation_date"`
}

type DeploymentUpdate struct {
	Replicas *int32            `json:"replicas,omitempty" jsonschema:"description:Desired number of replicas"`
	Images   map[string]string `json:"images,omitempty" jsonschema:"description:Map of container name to new image"`
	Labels   map[string]string `json:"labels,omitempty" jsonschema:"description:Updated labels"`
}

type RolloutStatus struct {
	Replicas        int32  `json:"replicas"`
	UpdatedReplicas int32  `json:"updated_replicas"`
	ReadyReplicas   int32  `json:"ready_replicas"`
	Unavailable     int32  `json:"unavailable_replicas"`
	Message         string `json:"message"`
}

type DeploymentCreate struct {
	Name      string            `json:"name" jsonschema:"required,description:Name of the deployment"`
	Namespace string            `json:"namespace" jsonschema:"required,description:Kubernetes namespace"`
	Image     string            `json:"image" jsonschema:"required,description:Container image (e.g. nginx:latest)"`
	Labels    map[string]string `json:"labels,omitempty" jsonschema:"description:Labels for the deployment"`
	Replicas  int32             `json:"replicas,omitempty" jsonschema:"description:Number of replicas (default 1)"`
	Port      int32             `json:"port,omitempty" jsonschema:"description:Container port to expose"`
	Env       map[string]string `json:"env,omitempty" jsonschema:"description:Environment variables as key-value pairs"`
}
