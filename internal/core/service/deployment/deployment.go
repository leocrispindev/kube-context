package deployment

import (
	"context"
	"fmt"

	"k8s-agent-new/internal/core/dto"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceInterface interface {
	ListDeployments(ctx context.Context, namespace string) (*dto.DeploymentList, error)
	UpdateDeployment(ctx context.Context, namespace string, deploymentName string, updates dto.DeploymentUpdate) error
	GetRolloutStatus(ctx context.Context, namespace string, deploymentName string) (*dto.RolloutStatus, error)
	TogglePauseDeployment(ctx context.Context, namespace string, deploymentName string, pause bool) error
	GetDeployment(ctx context.Context, namespace string, deploymentName string) (*dto.Deployment, error)
	CreateDeployment(ctx context.Context, deployCreate *dto.DeploymentCreate) (*dto.Deployment, error)
	DeleteDeployment(ctx context.Context, namespace string, deploymentName string) error
}

type Service struct {
	client kubernetes.Interface
}

func NewService(client kubernetes.Interface) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ListDeployments(ctx context.Context, namespace string) (*dto.DeploymentList, error) {
	deployments, err := s.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	var dtos []dto.Deployment
	for _, d := range deployments.Items {
		var images []string
		for _, c := range d.Spec.Template.Spec.Containers {
			images = append(images, c.Image)
		}

		dtos = append(dtos, dto.Deployment{
			Name:            d.Name,
			Namespace:       d.Namespace,
			Replicas:        *d.Spec.Replicas,
			ReadyReplicas:   d.Status.ReadyReplicas,
			UpdatedReplicas: d.Status.UpdatedReplicas,
			Strategy:        string(d.Spec.Strategy.Type),
			Images:          images,
			Labels:          d.Labels,
			CreationDate:    d.CreationTimestamp.String(),
		})
	}

	return &dto.DeploymentList{Deployments: dtos}, nil
}

func (s *Service) UpdateDeployment(ctx context.Context, namespace string, deploymentName string, updates dto.DeploymentUpdate) error {
	deployment, err := s.client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	if updates.Replicas != nil {
		deployment.Spec.Replicas = updates.Replicas
	}

	if len(updates.Images) > 0 {
		for i, c := range deployment.Spec.Template.Spec.Containers {
			if newImage, ok := updates.Images[c.Name]; ok {
				deployment.Spec.Template.Spec.Containers[i].Image = newImage
			}
		}
	}

	if len(updates.Labels) > 0 {
		if deployment.Labels == nil {
			deployment.Labels = make(map[string]string)
		}
		for k, v := range updates.Labels {
			deployment.Labels[k] = v
		}
	}

	_, err = s.client.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *Service) GetRolloutStatus(ctx context.Context, namespace string, deploymentName string) (*dto.RolloutStatus, error) {
	d, err := s.client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	message := "Deployment is healthy"
	if d.Generation > d.Status.ObservedGeneration {
		message = "Waiting for deployment spec update to be observed..."
	} else if d.Spec.Replicas != nil && d.Status.UpdatedReplicas < *d.Spec.Replicas {
		message = fmt.Sprintf("Waiting for rollout to finish: %d out of %d new replicas have been updated...", d.Status.UpdatedReplicas, *d.Spec.Replicas)
	} else if d.Status.Replicas > d.Status.UpdatedReplicas {
		message = fmt.Sprintf("Waiting for rollout to finish: %d old replicas are pending termination...", d.Status.Replicas-d.Status.UpdatedReplicas)
	} else if d.Status.AvailableReplicas < d.Status.UpdatedReplicas {
		message = fmt.Sprintf("Waiting for rollout to finish: %d of %d updated replicas are available...", d.Status.AvailableReplicas, d.Status.UpdatedReplicas)
	}

	return &dto.RolloutStatus{
		Replicas:        d.Status.Replicas,
		UpdatedReplicas: d.Status.UpdatedReplicas,
		ReadyReplicas:   d.Status.ReadyReplicas,
		Unavailable:     d.Status.UnavailableReplicas,
		Message:         message,
	}, nil
}

func (s *Service) TogglePauseDeployment(ctx context.Context, namespace string, deploymentName string, pause bool) error {
	deployment, err := s.client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	deployment.Spec.Paused = pause

	_, err = s.client.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment pause status: %w", err)
	}

	return nil
}

func (s *Service) GetDeployment(ctx context.Context, namespace string, deploymentName string) (*dto.Deployment, error) {
	d, err := s.client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	var images []string
	for _, c := range d.Spec.Template.Spec.Containers {
		images = append(images, c.Image)
	}

	return &dto.Deployment{
		Name:            d.Name,
		Namespace:       d.Namespace,
		Replicas:        *d.Spec.Replicas,
		ReadyReplicas:   d.Status.ReadyReplicas,
		UpdatedReplicas: d.Status.UpdatedReplicas,
		Strategy:        string(d.Spec.Strategy.Type),
		Images:          images,
		Labels:          d.Labels,
		CreationDate:    d.CreationTimestamp.String(),
	}, nil
}

func (s *Service) CreateDeployment(ctx context.Context, deployCreate *dto.DeploymentCreate) (*dto.Deployment, error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployCreate.Name,
			Namespace: deployCreate.Namespace,
			Labels:    deployCreate.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deployCreate.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deployCreate.Name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deployCreate.Name,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  deployCreate.Name,
							Image: deployCreate.Image,
						},
					},
				},
			},
		},
	}

	if deployCreate.Port > 0 {
		deployment.Spec.Template.Spec.Containers[0].Ports = []v1.ContainerPort{
			{
				ContainerPort: deployCreate.Port,
			},
		}
	}

	if len(deployCreate.Env) > 0 {
		var envVars []v1.EnvVar
		for key, value := range deployCreate.Env {
			envVars = append(envVars, v1.EnvVar{
				Name:  key,
				Value: value,
			})
		}
		deployment.Spec.Template.Spec.Containers[0].Env = envVars
	}

	if len(deployCreate.Labels) > 0 {
		for k, v := range deployCreate.Labels {
			deployment.Spec.Template.Labels[k] = v
		}
	}

	created, err := s.client.AppsV1().Deployments(deployCreate.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	return s.GetDeployment(ctx, created.Namespace, created.Name)
}

func (s *Service) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) error {
	err := s.client.AppsV1().Deployments(namespace).Delete(ctx, deploymentName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}
	return nil
}
