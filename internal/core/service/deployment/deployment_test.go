package deployment

import (
	"context"
	"fmt"
	"testing"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func int32Ptr(i int32) *int32 { return &i }

type DeploymentServiceTestSuite struct {
	suite.Suite
	client *fake.Clientset
	svc    *Service
}

func (s *DeploymentServiceTestSuite) SetupTest() {
	s.client = fake.NewSimpleClientset()
	s.svc = NewService(s.client)
}

func TestDeploymentServiceSuite(t *testing.T) {
	suite.Run(t, new(DeploymentServiceTestSuite))
}

func (s *DeploymentServiceTestSuite) seedDeployment(name, namespace string, replicas int32) {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": name},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": name},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{Name: name, Image: "nginx:latest"},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{
			Replicas:          replicas,
			ReadyReplicas:     replicas,
			UpdatedReplicas:   replicas,
			AvailableReplicas: replicas,
		},
	}
	_, err := s.client.AppsV1().Deployments(namespace).Create(context.Background(), dep, metav1.CreateOptions{})
	s.Require().NoError(err)
}

// ---------- ListDeployments ----------

func (s *DeploymentServiceTestSuite) TestListDeployments_Success() {
	s.seedDeployment("web", "default", 3)
	s.seedDeployment("api", "default", 2)

	result, err := s.svc.ListDeployments(context.Background(), "default")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Deployments, 2)

	names := []string{result.Deployments[0].Name, result.Deployments[1].Name}
	assert.Contains(s.T(), names, "web")
	assert.Contains(s.T(), names, "api")

	for _, d := range result.Deployments {
		assert.Equal(s.T(), "default", d.Namespace)
		assert.Equal(s.T(), "RollingUpdate", d.Strategy)
		assert.Equal(s.T(), []string{"nginx:latest"}, d.Images)
		assert.NotEmpty(s.T(), d.Labels)
	}
}

func (s *DeploymentServiceTestSuite) TestListDeployments_EmptyNamespace() {
	result, err := s.svc.ListDeployments(context.Background(), "empty-ns")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Deployments)
}

func (s *DeploymentServiceTestSuite) TestListDeployments_Error() {
	s.client.Fake.PrependReactor("list", "deployments", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, fmt.Errorf("internal server error")
	})

	result, err := s.svc.ListDeployments(context.Background(), "default")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to list deployments")
}

// ---------- GetDeployment ----------

func (s *DeploymentServiceTestSuite) TestGetDeployment_Success() {
	s.seedDeployment("my-deploy", "default", 3)

	result, err := s.svc.GetDeployment(context.Background(), "default", "my-deploy")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "my-deploy", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), int32(3), result.Replicas)
	assert.Equal(s.T(), int32(3), result.ReadyReplicas)
	assert.Equal(s.T(), int32(3), result.UpdatedReplicas)
	assert.Equal(s.T(), "RollingUpdate", result.Strategy)
	assert.Equal(s.T(), []string{"nginx:latest"}, result.Images)
	assert.Equal(s.T(), map[string]string{"app": "my-deploy"}, result.Labels)
	assert.NotEmpty(s.T(), result.CreationDate)
}

func (s *DeploymentServiceTestSuite) TestGetDeployment_NotFound() {
	result, err := s.svc.GetDeployment(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get deployment")
}

// ---------- CreateDeployment ----------

func (s *DeploymentServiceTestSuite) TestCreateDeployment_Success() {
	create := &dto.DeploymentCreate{
		Name:      "new-deploy",
		Namespace: "default",
		Image:     "nginx:1.21",
		Labels:    map[string]string{"env": "test"},
		Replicas:  2,
		Port:      80,
		Env:       map[string]string{"APP_ENV": "production"},
	}

	result, err := s.svc.CreateDeployment(context.Background(), create)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "new-deploy", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), int32(2), result.Replicas)
	assert.Contains(s.T(), result.Images, "nginx:1.21")
	assert.Equal(s.T(), "test", result.Labels["env"])

	got, err := s.client.AppsV1().Deployments("default").Get(context.Background(), "new-deploy", metav1.GetOptions{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int32(80), got.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort)
}

func (s *DeploymentServiceTestSuite) TestCreateDeployment_AlreadyExists() {
	s.seedDeployment("existing", "default", 1)

	create := &dto.DeploymentCreate{
		Name:      "existing",
		Namespace: "default",
		Image:     "nginx:latest",
		Replicas:  1,
	}

	result, err := s.svc.CreateDeployment(context.Background(), create)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to create deployment")
}

// ---------- UpdateDeployment ----------

func (s *DeploymentServiceTestSuite) TestUpdateDeployment_Replicas() {
	s.seedDeployment("my-deploy", "default", 2)

	err := s.svc.UpdateDeployment(context.Background(), "default", "my-deploy", dto.DeploymentUpdate{
		Replicas: int32Ptr(5),
	})

	assert.NoError(s.T(), err)

	got, err := s.client.AppsV1().Deployments("default").Get(context.Background(), "my-deploy", metav1.GetOptions{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int32(5), *got.Spec.Replicas)
}

func (s *DeploymentServiceTestSuite) TestUpdateDeployment_Images() {
	s.seedDeployment("my-deploy", "default", 1)

	err := s.svc.UpdateDeployment(context.Background(), "default", "my-deploy", dto.DeploymentUpdate{
		Images: map[string]string{"my-deploy": "nginx:1.25"},
	})

	assert.NoError(s.T(), err)

	got, err := s.client.AppsV1().Deployments("default").Get(context.Background(), "my-deploy", metav1.GetOptions{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "nginx:1.25", got.Spec.Template.Spec.Containers[0].Image)
}

func (s *DeploymentServiceTestSuite) TestUpdateDeployment_Labels() {
	s.seedDeployment("my-deploy", "default", 1)

	err := s.svc.UpdateDeployment(context.Background(), "default", "my-deploy", dto.DeploymentUpdate{
		Labels: map[string]string{"version": "v2"},
	})

	assert.NoError(s.T(), err)

	got, err := s.client.AppsV1().Deployments("default").Get(context.Background(), "my-deploy", metav1.GetOptions{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "v2", got.Labels["version"])
	assert.Equal(s.T(), "my-deploy", got.Labels["app"])
}

func (s *DeploymentServiceTestSuite) TestUpdateDeployment_NotFound() {
	err := s.svc.UpdateDeployment(context.Background(), "default", "nonexistent", dto.DeploymentUpdate{
		Replicas: int32Ptr(3),
	})

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to get deployment")
}

// ---------- GetRolloutStatus ----------

func (s *DeploymentServiceTestSuite) TestGetRolloutStatus_Healthy() {
	s.seedDeployment("healthy", "default", 3)

	result, err := s.svc.GetRolloutStatus(context.Background(), "default", "healthy")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int32(3), result.Replicas)
	assert.Equal(s.T(), int32(3), result.ReadyReplicas)
	assert.Equal(s.T(), int32(3), result.UpdatedReplicas)
	assert.Equal(s.T(), int32(0), result.Unavailable)
	assert.Equal(s.T(), "Deployment is healthy", result.Message)
}

func (s *DeploymentServiceTestSuite) TestGetRolloutStatus_WaitingForUpdate() {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rolling",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(3),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "rolling"}},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "rolling"}},
				Spec:       v1.PodSpec{Containers: []v1.Container{{Name: "rolling", Image: "nginx"}}},
			},
		},
		Status: appsv1.DeploymentStatus{
			Replicas:        3,
			UpdatedReplicas: 1,
			ReadyReplicas:   2,
		},
	}
	_, err := s.client.AppsV1().Deployments("default").Create(context.Background(), dep, metav1.CreateOptions{})
	s.Require().NoError(err)

	result, err := s.svc.GetRolloutStatus(context.Background(), "default", "rolling")

	assert.NoError(s.T(), err)
	assert.Contains(s.T(), result.Message, "Waiting for rollout to finish")
	assert.Contains(s.T(), result.Message, "1 out of 3 new replicas have been updated")
}

func (s *DeploymentServiceTestSuite) TestGetRolloutStatus_WaitingAvailability() {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pending",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(3),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "pending"}},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "pending"}},
				Spec:       v1.PodSpec{Containers: []v1.Container{{Name: "pending", Image: "nginx"}}},
			},
		},
		Status: appsv1.DeploymentStatus{
			Replicas:          3,
			UpdatedReplicas:   3,
			ReadyReplicas:     1,
			AvailableReplicas: 1,
		},
	}
	_, err := s.client.AppsV1().Deployments("default").Create(context.Background(), dep, metav1.CreateOptions{})
	s.Require().NoError(err)

	result, err := s.svc.GetRolloutStatus(context.Background(), "default", "pending")

	assert.NoError(s.T(), err)
	assert.Contains(s.T(), result.Message, "1 of 3 updated replicas are available")
}

func (s *DeploymentServiceTestSuite) TestGetRolloutStatus_NotFound() {
	result, err := s.svc.GetRolloutStatus(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get deployment")
}

// ---------- TogglePauseDeployment ----------

func (s *DeploymentServiceTestSuite) TestTogglePauseDeployment_Pause() {
	s.seedDeployment("pausable", "default", 1)

	err := s.svc.TogglePauseDeployment(context.Background(), "default", "pausable", true)

	assert.NoError(s.T(), err)

	got, err := s.client.AppsV1().Deployments("default").Get(context.Background(), "pausable", metav1.GetOptions{})
	assert.NoError(s.T(), err)
	assert.True(s.T(), got.Spec.Paused)
}

func (s *DeploymentServiceTestSuite) TestTogglePauseDeployment_Unpause() {
	s.seedDeployment("pausable", "default", 1)
	_ = s.svc.TogglePauseDeployment(context.Background(), "default", "pausable", true)

	err := s.svc.TogglePauseDeployment(context.Background(), "default", "pausable", false)

	assert.NoError(s.T(), err)

	got, err := s.client.AppsV1().Deployments("default").Get(context.Background(), "pausable", metav1.GetOptions{})
	assert.NoError(s.T(), err)
	assert.False(s.T(), got.Spec.Paused)
}

func (s *DeploymentServiceTestSuite) TestTogglePauseDeployment_NotFound() {
	err := s.svc.TogglePauseDeployment(context.Background(), "default", "nonexistent", true)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to get deployment")
}

// ---------- DeleteDeployment ----------

func (s *DeploymentServiceTestSuite) TestDeleteDeployment_Success() {
	s.seedDeployment("to-delete", "default", 1)

	err := s.svc.DeleteDeployment(context.Background(), "default", "to-delete")

	assert.NoError(s.T(), err)

	_, err = s.client.AppsV1().Deployments("default").Get(context.Background(), "to-delete", metav1.GetOptions{})
	assert.Error(s.T(), err)
}

func (s *DeploymentServiceTestSuite) TestDeleteDeployment_NotFound() {
	err := s.svc.DeleteDeployment(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to delete deployment")
}
