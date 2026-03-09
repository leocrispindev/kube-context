package namespace

import (
	"context"
	"testing"

	"github.com/leocrispindev/kube-context/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type NamespaceServiceTestSuite struct {
	suite.Suite
	client *fake.Clientset
	svc    *Service
}

func (s *NamespaceServiceTestSuite) SetupTest() {
	s.client = fake.NewSimpleClientset()
	s.svc = NewService(s.client)
}

func TestNamespaceServiceSuite(t *testing.T) {
	suite.Run(t, new(NamespaceServiceTestSuite))
}

// --- ListNamespaces ---

func (s *NamespaceServiceTestSuite) TestListNamespaces_Success() {
	ctx := context.Background()

	namespaces := []*v1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "default",
				Labels:      map[string]string{"env": "production"},
				Annotations: map[string]string{"note": "primary"},
			},
			Status: v1.NamespaceStatus{Phase: v1.NamespaceActive},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "kube-system",
				Labels:      map[string]string{"env": "system"},
				Annotations: map[string]string{"note": "system-ns"},
			},
			Status: v1.NamespaceStatus{Phase: v1.NamespaceActive},
		},
	}

	for _, ns := range namespaces {
		_, err := s.client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
		s.Require().NoError(err)
	}

	result, err := s.svc.ListNamespaces(ctx)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Namespaces, 2)

	assert.Equal(s.T(), "default", result.Namespaces[0].Name)
	assert.Equal(s.T(), map[string]string{"env": "production"}, result.Namespaces[0].Labels)
	assert.Equal(s.T(), map[string]string{"note": "primary"}, result.Namespaces[0].Annotations)
	assert.Equal(s.T(), string(v1.NamespaceActive), result.Namespaces[0].Status)

	assert.Equal(s.T(), "kube-system", result.Namespaces[1].Name)
}

func (s *NamespaceServiceTestSuite) TestListNamespaces_Empty() {
	ctx := context.Background()

	result, err := s.svc.ListNamespaces(ctx)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Namespaces)
}

// --- CreateNamespace ---

func (s *NamespaceServiceTestSuite) TestCreateNamespace_Success() {
	ctx := context.Background()
	labels := map[string]string{"team": "backend"}
	annotations := map[string]string{"owner": "dev"}

	err := s.svc.CreateNamespace(ctx, "test-ns", labels, annotations)

	assert.NoError(s.T(), err)

	ns, getErr := s.client.CoreV1().Namespaces().Get(ctx, "test-ns", metav1.GetOptions{})
	assert.NoError(s.T(), getErr)
	assert.Equal(s.T(), "test-ns", ns.Name)
	assert.Equal(s.T(), labels, ns.Labels)
	assert.Equal(s.T(), annotations, ns.Annotations)
}

func (s *NamespaceServiceTestSuite) TestCreateNamespace_AlreadyExists() {
	ctx := context.Background()

	err := s.svc.CreateNamespace(ctx, "dup-ns", nil, nil)
	s.Require().NoError(err)

	err = s.svc.CreateNamespace(ctx, "dup-ns", nil, nil)

	assert.Error(s.T(), err)
}

// --- DeleteNamespace ---

func (s *NamespaceServiceTestSuite) TestDeleteNamespace_Success() {
	ctx := context.Background()

	err := s.svc.CreateNamespace(ctx, "to-delete", nil, nil)
	s.Require().NoError(err)

	err = s.svc.DeleteNamespace(ctx, "to-delete")

	assert.NoError(s.T(), err)

	_, getErr := s.client.CoreV1().Namespaces().Get(ctx, "to-delete", metav1.GetOptions{})
	assert.Error(s.T(), getErr)
}

func (s *NamespaceServiceTestSuite) TestDeleteNamespace_NotFound() {
	ctx := context.Background()

	err := s.svc.DeleteNamespace(ctx, "nonexistent")

	assert.Error(s.T(), err)
}

// --- GetNamespace ---

func (s *NamespaceServiceTestSuite) TestGetNamespace_Success() {
	ctx := context.Background()

	_, err := s.client.CoreV1().Namespaces().Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "my-ns",
			Labels:      map[string]string{"app": "web"},
			Annotations: map[string]string{"description": "web namespace"},
		},
		Status: v1.NamespaceStatus{Phase: v1.NamespaceActive},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	result, err := s.svc.GetNamespace(ctx, "my-ns")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), &dto.Namespace{
		Name:        "my-ns",
		Labels:      map[string]string{"app": "web"},
		Annotations: map[string]string{"description": "web namespace"},
		Status:      string(v1.NamespaceActive),
	}, result)
}

func (s *NamespaceServiceTestSuite) TestGetNamespace_NotFound() {
	ctx := context.Background()

	result, err := s.svc.GetNamespace(ctx, "missing-ns")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

// --- UpdateNamespace ---

func (s *NamespaceServiceTestSuite) TestUpdateNamespace_Success() {
	ctx := context.Background()

	_, err := s.client.CoreV1().Namespaces().Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "update-ns",
			Labels: map[string]string{"version": "v1"},
		},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	newLabels := map[string]string{"version": "v2"}
	newAnnotations := map[string]string{"updated": "true"}

	err = s.svc.UpdateNamespace(ctx, "update-ns", newLabels, newAnnotations)

	assert.NoError(s.T(), err)

	ns, getErr := s.client.CoreV1().Namespaces().Get(ctx, "update-ns", metav1.GetOptions{})
	assert.NoError(s.T(), getErr)
	assert.Equal(s.T(), newLabels, ns.Labels)
	assert.Equal(s.T(), newAnnotations, ns.Annotations)
}

func (s *NamespaceServiceTestSuite) TestUpdateNamespace_NotFound() {
	ctx := context.Background()

	err := s.svc.UpdateNamespace(ctx, "ghost-ns", nil, nil)

	assert.Error(s.T(), err)
}

// --- UpdateResources ---

func (s *NamespaceServiceTestSuite) TestUpdateResources_Success() {
	ctx := context.Background()

	_, err := s.client.CoreV1().Namespaces().Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "resource-ns"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	resources := v1.ResourceQuotaSpec{
		Hard: v1.ResourceList{
			v1.ResourceCPU: resource.MustParse("4"),
		},
	}

	err = s.svc.UpdateResources(ctx, "resource-ns", resources)

	assert.NoError(s.T(), err)
}

func (s *NamespaceServiceTestSuite) TestUpdateResources_NotFound() {
	ctx := context.Background()

	resources := v1.ResourceQuotaSpec{
		Hard: v1.ResourceList{
			v1.ResourceCPU: resource.MustParse("2"),
		},
	}

	err := s.svc.UpdateResources(ctx, "no-ns", resources)

	assert.Error(s.T(), err)
}
