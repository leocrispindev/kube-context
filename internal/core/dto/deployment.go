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
	Replicas *int32            `json:"replicas,omitempty"`
	Images   map[string]string `json:"images,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
}

type RolloutStatus struct {
	Replicas        int32  `json:"replicas"`
	UpdatedReplicas int32  `json:"updated_replicas"`
	ReadyReplicas   int32  `json:"ready_replicas"`
	Unavailable     int32  `json:"unavailable_replicas"`
	Message         string `json:"message"`
}

type DeploymentCreate struct {
	Name      string            `json:"name" binding:"required"`
	Namespace string            `json:"namespace" binding:"required"`
	Image     string            `json:"image" binding:"required"`
	Labels    map[string]string `json:"labels,omitempty"`
	Replicas  int32             `json:"replicas,omitempty"`
	Port      int32             `json:"port,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
}
