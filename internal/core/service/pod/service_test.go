package pod

import (
	"context"
	"fmt"
	"testing"

	"k8s-agent-new/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type PodServiceTestSuite struct {
	suite.Suite
	client *fake.Clientset
	svc    *Service
}

func (s *PodServiceTestSuite) SetupTest() {
	s.client = fake.NewSimpleClientset()
	s.svc = NewService(s.client)
}

func TestPodServiceSuite(t *testing.T) {
	suite.Run(t, new(PodServiceTestSuite))
}

func (s *PodServiceTestSuite) seedPod(name, namespace string) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": name},
		},
		Spec: v1.PodSpec{
			NodeName: "node-1",
			Containers: []v1.Container{
				{Name: name, Image: "nginx:latest"},
			},
		},
		Status: v1.PodStatus{
			Phase: v1.PodRunning,
			PodIP: "10.0.0.1",
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:         name,
					RestartCount: 2,
				},
			},
		},
	}
	_, err := s.client.CoreV1().Pods(namespace).Create(context.Background(), pod, metav1.CreateOptions{})
	s.Require().NoError(err)
}

// ---------- ListPods ----------

func (s *PodServiceTestSuite) TestListPods_Success() {
	s.seedPod("pod-a", "default")
	s.seedPod("pod-b", "default")

	result, err := s.svc.ListPods(context.Background(), "default")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Pods, 2)
	assert.Contains(s.T(), result.Pods, "pod-a")
	assert.Contains(s.T(), result.Pods, "pod-b")
}

func (s *PodServiceTestSuite) TestListPods_EmptyNamespace() {
	result, err := s.svc.ListPods(context.Background(), "empty-ns")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Pods)
}

func (s *PodServiceTestSuite) TestListPods_Error() {
	s.client.Fake.PrependReactor("list", "pods", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, fmt.Errorf("internal server error")
	})

	result, err := s.svc.ListPods(context.Background(), "default")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

// ---------- GetPod ----------

func (s *PodServiceTestSuite) TestGetPod_Success() {
	s.seedPod("my-pod", "default")

	result, err := s.svc.GetPod(context.Background(), "default", "my-pod")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "my-pod", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), "Running", result.Status)
	assert.Equal(s.T(), "10.0.0.1", result.PodIP)
	assert.Equal(s.T(), "node-1", result.NodeName)
	assert.Equal(s.T(), []string{"nginx:latest"}, result.Images)
	assert.Equal(s.T(), int32(2), result.Restarts)
	assert.Equal(s.T(), map[string]string{"app": "my-pod"}, result.Labels)
	assert.NotEmpty(s.T(), result.CreatedAt)
}

func (s *PodServiceTestSuite) TestGetPod_WaitingStatus() {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "waiting-pod",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{Name: "app", Image: "bad-image:latest"}},
		},
		Status: v1.PodStatus{
			Phase: v1.PodPending,
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name: "app",
					State: v1.ContainerState{
						Waiting: &v1.ContainerStateWaiting{Reason: "ImagePullBackOff"},
					},
				},
			},
		},
	}
	_, err := s.client.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	s.Require().NoError(err)

	result, err := s.svc.GetPod(context.Background(), "default", "waiting-pod")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "ImagePullBackOff", result.Status)
}

func (s *PodServiceTestSuite) TestGetPod_NotFound() {
	result, err := s.svc.GetPod(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get pod")
}

// ---------- DeletePod ----------

func (s *PodServiceTestSuite) TestDeletePod_Success() {
	s.seedPod("to-delete", "default")

	err := s.svc.DeletePod(context.Background(), "default", "to-delete")

	assert.NoError(s.T(), err)

	_, err = s.client.CoreV1().Pods("default").Get(context.Background(), "to-delete", metav1.GetOptions{})
	assert.Error(s.T(), err)
}

func (s *PodServiceTestSuite) TestDeletePod_NotFound() {
	err := s.svc.DeletePod(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to delete pod")
}

// ---------- RestartPod ----------

func (s *PodServiceTestSuite) TestRestartPod_Success() {
	s.seedPod("restart-me", "default")

	err := s.svc.RestartPod(context.Background(), "default", "restart-me")

	assert.NoError(s.T(), err)

	_, err = s.client.CoreV1().Pods("default").Get(context.Background(), "restart-me", metav1.GetOptions{})
	assert.Error(s.T(), err)
}

func (s *PodServiceTestSuite) TestRestartPod_NotFound() {
	err := s.svc.RestartPod(context.Background(), "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to delete pod")
}

// ---------- CreatePod ----------

func (s *PodServiceTestSuite) TestCreatePod_Success() {
	create := &dto.PodCreate{
		Name:      "new-pod",
		Namespace: "default",
		Image:     "redis:7",
		Labels:    map[string]string{"app": "redis"},
		Port:      6379,
		Env:       map[string]string{"REDIS_PASSWORD": "secret"},
		Command:   []string{"redis-server"},
		Args:      []string{"--requirepass", "secret"},
	}

	result, err := s.svc.CreatePod(context.Background(), create)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "new-pod", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), []string{"redis:7"}, result.Images)
	assert.Equal(s.T(), "redis", result.Labels["app"])

	got, err := s.client.CoreV1().Pods("default").Get(context.Background(), "new-pod", metav1.GetOptions{})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int32(6379), got.Spec.Containers[0].Ports[0].ContainerPort)
	assert.Equal(s.T(), []string{"redis-server"}, got.Spec.Containers[0].Command)
	assert.Equal(s.T(), []string{"--requirepass", "secret"}, got.Spec.Containers[0].Args)
}

func (s *PodServiceTestSuite) TestCreatePod_MinimalFields() {
	create := &dto.PodCreate{
		Name:      "minimal",
		Namespace: "default",
		Image:     "busybox",
	}

	result, err := s.svc.CreatePod(context.Background(), create)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "minimal", result.Name)
	assert.Equal(s.T(), []string{"busybox"}, result.Images)
}

func (s *PodServiceTestSuite) TestCreatePod_AlreadyExists() {
	s.seedPod("existing", "default")

	create := &dto.PodCreate{
		Name:      "existing",
		Namespace: "default",
		Image:     "nginx:latest",
	}

	result, err := s.svc.CreatePod(context.Background(), create)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to create pod")
}

// ---------- UpdatePod ----------

func (s *PodServiceTestSuite) TestUpdatePod_Labels() {
	s.seedPod("update-me", "default")

	result, err := s.svc.UpdatePod(context.Background(), "default", "update-me", &dto.PodUpdate{
		Labels: map[string]string{"version": "v2", "tier": "backend"},
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "v2", result.Labels["version"])
	assert.Equal(s.T(), "backend", result.Labels["tier"])
	assert.Equal(s.T(), "update-me", result.Labels["app"])
}

func (s *PodServiceTestSuite) TestUpdatePod_Image() {
	s.seedPod("update-me", "default")

	result, err := s.svc.UpdatePod(context.Background(), "default", "update-me", &dto.PodUpdate{
		Image: "nginx:1.25",
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), []string{"nginx:1.25"}, result.Images)
}

func (s *PodServiceTestSuite) TestUpdatePod_LabelsAndImage() {
	s.seedPod("update-me", "default")

	result, err := s.svc.UpdatePod(context.Background(), "default", "update-me", &dto.PodUpdate{
		Labels: map[string]string{"env": "staging"},
		Image:  "nginx:1.25",
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "staging", result.Labels["env"])
	assert.Equal(s.T(), []string{"nginx:1.25"}, result.Images)
}

func (s *PodServiceTestSuite) TestUpdatePod_NotFound() {
	result, err := s.svc.UpdatePod(context.Background(), "default", "nonexistent", &dto.PodUpdate{
		Image: "nginx:latest",
	})

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "failed to get pod")
}
