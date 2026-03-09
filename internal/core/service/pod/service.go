package pod

import (
	"context"
	"fmt"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceInterface interface {
	ListPods(ctx context.Context, namespace string) (*dto.PodList, error)
	ScalePod(ctx context.Context, namespace string, podName string, replicas int32) (*dto.PodList, error)
	GetPod(ctx context.Context, namespace string, podName string) (*dto.PodDetails, error)
	DeletePod(ctx context.Context, namespace string, podName string) error
	RestartPod(ctx context.Context, namespace string, podName string) error
	CreatePod(ctx context.Context, podCreate *dto.PodCreate) (*dto.PodDetails, error)
	UpdatePod(ctx context.Context, namespace string, podName string, updates *dto.PodUpdate) (*dto.PodDetails, error)
}

type Service struct {
	client kubernetes.Interface
}

func NewService(client kubernetes.Interface) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ListPods(ctx context.Context, namespace string) (*dto.PodList, error) {
	pods, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return &dto.PodList{
		Pods: podNames,
	}, nil
}

func (s *Service) ScalePod(ctx context.Context, namespace string, podName string, replicas int32) (*dto.PodList, error) {
	pod, err := s.client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	if len(pod.OwnerReferences) == 0 {
		return nil, fmt.Errorf("pod %s has no owner references, cannot scale", podName)
	}

	for _, owner := range pod.OwnerReferences {
		switch owner.Kind {
		case "ReplicaSet":
			if err := s.scaleReplicaSet(ctx, namespace, podName, owner.Name, replicas); err != nil {
				return nil, err
			}
		case "StatefulSet":
			if err := s.scaleStatefulSet(ctx, namespace, owner.Name, replicas); err != nil {
				return nil, err
			}
		}
	}

	return s.ListPods(ctx, namespace)
}

func (s *Service) scaleStatefulSet(ctx context.Context, namespace string, statefulSetName string, replicas int32) error {
	sts, err := s.client.AppsV1().StatefulSets(namespace).Get(ctx, statefulSetName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get statefulset %s: %w", statefulSetName, err)
	}
	sts.Spec.Replicas = &replicas
	_, err = s.client.AppsV1().StatefulSets(namespace).Update(ctx, sts, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update statefulset: %w", err)
	}
	return nil
}

func (s *Service) scaleReplicaSet(ctx context.Context, namespace string, podName string, ownerName string, replicas int32) error {
	deploymentName := ""
	rs, err := s.client.AppsV1().ReplicaSets(namespace).Get(ctx, ownerName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get replicaset %s: %w", ownerName, err)
	}
	for _, rsOwner := range rs.OwnerReferences {
		if rsOwner.Kind == "Deployment" {
			deploymentName = rsOwner.Name
			break
		}
	}

	if deploymentName == "" {
		return fmt.Errorf("could not find a deployment managing pod %s", podName)
	}

	deployment, err := s.client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment %s: %w", deploymentName, err)
	}

	deployment.Spec.Replicas = &replicas
	_, err = s.client.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *Service) GetPod(ctx context.Context, namespace string, podName string) (*dto.PodDetails, error) {
	pod, err := s.client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	status := string(pod.Status.Phase)
	restarts := int32(0)
	var images []string

	for _, container := range pod.Spec.Containers {
		images = append(images, container.Image)
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		restarts += containerStatus.RestartCount
		if containerStatus.State.Waiting != nil {
			status = containerStatus.State.Waiting.Reason
		}
	}

	if pod.DeletionTimestamp != nil {
		status = "Terminating"
	}

	return &dto.PodDetails{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		Status:    status,
		PodIP:     pod.Status.PodIP,
		NodeName:  pod.Spec.NodeName,
		CreatedAt: pod.CreationTimestamp.String(),
		Labels:    pod.Labels,
		Images:    images,
		Restarts:  restarts,
	}, nil
}

func (s *Service) DeletePod(ctx context.Context, namespace string, podName string) error {
	err := s.client.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod: %w", err)
	}
	return nil
}

func (s *Service) RestartPod(ctx context.Context, namespace string, podName string) error {
	return s.DeletePod(ctx, namespace, podName)
}

func (s *Service) CreatePod(ctx context.Context, podCreate *dto.PodCreate) (*dto.PodDetails, error) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podCreate.Name,
			Namespace: podCreate.Namespace,
			Labels:    podCreate.Labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    podCreate.Name,
					Image:   podCreate.Image,
					Command: podCreate.Command,
					Args:    podCreate.Args,
				},
			},
		},
	}

	if podCreate.Port > 0 {
		pod.Spec.Containers[0].Ports = []v1.ContainerPort{
			{
				ContainerPort: podCreate.Port,
			},
		}
	}

	if len(podCreate.Env) > 0 {
		var envVars []v1.EnvVar
		for key, value := range podCreate.Env {
			envVars = append(envVars, v1.EnvVar{
				Name:  key,
				Value: value,
			})
		}
		pod.Spec.Containers[0].Env = envVars
	}

	createdPod, err := s.client.CoreV1().Pods(podCreate.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %w", err)
	}

	return s.GetPod(ctx, createdPod.Namespace, createdPod.Name)
}

func (s *Service) UpdatePod(ctx context.Context, namespace string, podName string, updates *dto.PodUpdate) (*dto.PodDetails, error) {
	pod, err := s.client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	if len(updates.Labels) > 0 {
		if pod.Labels == nil {
			pod.Labels = make(map[string]string)
		}
		for k, v := range updates.Labels {
			pod.Labels[k] = v
		}
	}

	if updates.Image != "" {
		if len(pod.Spec.Containers) > 0 {
			pod.Spec.Containers[0].Image = updates.Image
		}
	}

	_, err = s.client.CoreV1().Pods(namespace).Update(ctx, pod, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update pod: %w", err)
	}

	return s.GetPod(ctx, namespace, podName)
}
