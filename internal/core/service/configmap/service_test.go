package configmap

import (
	"context"
	"testing"

	"k8s-agent-new/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type ConfigMapServiceTestSuite struct {
	suite.Suite
	client *fake.Clientset
	svc    *Service
}

func (s *ConfigMapServiceTestSuite) SetupTest() {
	s.client = fake.NewSimpleClientset()
	s.svc = NewService(s.client)
}

func TestConfigMapServiceSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapServiceTestSuite))
}

// --- ListConfigMaps ---

func (s *ConfigMapServiceTestSuite) TestListConfigMaps_Success() {
	ctx := context.Background()
	namespace := "default"

	configMaps := []*v1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-one",
				Namespace: namespace,
				Labels:    map[string]string{"tier": "frontend"},
			},
			Data: map[string]string{"key1": "val1"},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-two",
				Namespace: namespace,
				Labels:    map[string]string{"tier": "backend"},
			},
			Data: map[string]string{"key2": "val2"},
		},
	}

	for _, cm := range configMaps {
		_, err := s.client.CoreV1().ConfigMaps(namespace).Create(ctx, cm, metav1.CreateOptions{})
		s.Require().NoError(err)
	}

	result, err := s.svc.ListConfigMaps(ctx, namespace)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.ConfigMaps, 2)
	assert.Equal(s.T(), "cm-one", result.ConfigMaps[0].Name)
	assert.Equal(s.T(), namespace, result.ConfigMaps[0].Namespace)
	assert.Equal(s.T(), map[string]string{"key1": "val1"}, result.ConfigMaps[0].Data)
	assert.Equal(s.T(), "cm-two", result.ConfigMaps[1].Name)
}

func (s *ConfigMapServiceTestSuite) TestListConfigMaps_Empty() {
	ctx := context.Background()

	result, err := s.svc.ListConfigMaps(ctx, "empty-ns")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.ConfigMaps)
}

func (s *ConfigMapServiceTestSuite) TestListConfigMaps_FiltersByNamespace() {
	ctx := context.Background()

	_, err := s.client.CoreV1().ConfigMaps("ns-a").Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "cm-a", Namespace: "ns-a"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	_, err = s.client.CoreV1().ConfigMaps("ns-b").Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "cm-b", Namespace: "ns-b"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	result, err := s.svc.ListConfigMaps(ctx, "ns-a")

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result.ConfigMaps, 1)
	assert.Equal(s.T(), "cm-a", result.ConfigMaps[0].Name)
}

// --- GetConfigMap ---

func (s *ConfigMapServiceTestSuite) TestGetConfigMap_Success() {
	ctx := context.Background()
	namespace := "default"

	_, err := s.client.CoreV1().ConfigMaps(namespace).Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cm",
			Namespace: namespace,
			Labels:    map[string]string{"app": "api"},
		},
		Data: map[string]string{"config.yaml": "key: value"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	result, err := s.svc.GetConfigMap(ctx, namespace, "my-cm")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "my-cm", result.Name)
	assert.Equal(s.T(), namespace, result.Namespace)
	assert.Equal(s.T(), map[string]string{"config.yaml": "key: value"}, result.Data)
	assert.Equal(s.T(), map[string]string{"app": "api"}, result.Labels)
	assert.NotEmpty(s.T(), result.CreationDate)
}

func (s *ConfigMapServiceTestSuite) TestGetConfigMap_NotFound() {
	ctx := context.Background()

	result, err := s.svc.GetConfigMap(ctx, "default", "nonexistent")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

// --- CreateConfigMap ---

func (s *ConfigMapServiceTestSuite) TestCreateConfigMap_Success() {
	ctx := context.Background()

	input := &dto.ConfigMapCreate{
		Name:      "new-cm",
		Namespace: "default",
		Data:      map[string]string{"db_host": "localhost"},
		Labels:    map[string]string{"env": "dev"},
	}

	result, err := s.svc.CreateConfigMap(ctx, input)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "new-cm", result.Name)
	assert.Equal(s.T(), "default", result.Namespace)
	assert.Equal(s.T(), map[string]string{"db_host": "localhost"}, result.Data)
	assert.Equal(s.T(), map[string]string{"env": "dev"}, result.Labels)

	raw, getErr := s.client.CoreV1().ConfigMaps("default").Get(ctx, "new-cm", metav1.GetOptions{})
	assert.NoError(s.T(), getErr)
	assert.Equal(s.T(), "new-cm", raw.Name)
}

func (s *ConfigMapServiceTestSuite) TestCreateConfigMap_AlreadyExists() {
	ctx := context.Background()

	input := &dto.ConfigMapCreate{
		Name:      "dup-cm",
		Namespace: "default",
		Data:      map[string]string{"k": "v"},
	}

	_, err := s.svc.CreateConfigMap(ctx, input)
	s.Require().NoError(err)

	_, err = s.svc.CreateConfigMap(ctx, input)

	assert.Error(s.T(), err)
}

// --- UpdateConfigMap ---

func (s *ConfigMapServiceTestSuite) TestUpdateConfigMap_Success() {
	ctx := context.Background()
	namespace := "default"

	_, err := s.client.CoreV1().ConfigMaps(namespace).Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "upd-cm",
			Namespace: namespace,
			Labels:    map[string]string{"version": "v1", "keep": "yes"},
		},
		Data: map[string]string{"original_key": "original_val", "keep_key": "keep_val"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	updates := &dto.ConfigMapUpdate{
		Data:   map[string]string{"original_key": "new_val", "added_key": "added_val"},
		Labels: map[string]string{"version": "v2", "new_label": "true"},
	}

	result, err := s.svc.UpdateConfigMap(ctx, namespace, "upd-cm", updates)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)

	assert.Equal(s.T(), "new_val", result.Data["original_key"])
	assert.Equal(s.T(), "keep_val", result.Data["keep_key"])
	assert.Equal(s.T(), "added_val", result.Data["added_key"])

	assert.Equal(s.T(), "v2", result.Labels["version"])
	assert.Equal(s.T(), "yes", result.Labels["keep"])
	assert.Equal(s.T(), "true", result.Labels["new_label"])
}

func (s *ConfigMapServiceTestSuite) TestUpdateConfigMap_NotFound() {
	ctx := context.Background()

	updates := &dto.ConfigMapUpdate{
		Data: map[string]string{"k": "v"},
	}

	result, err := s.svc.UpdateConfigMap(ctx, "default", "ghost-cm", updates)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *ConfigMapServiceTestSuite) TestUpdateConfigMap_EmptyUpdates() {
	ctx := context.Background()
	namespace := "default"

	_, err := s.client.CoreV1().ConfigMaps(namespace).Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "noop-cm",
			Namespace: namespace,
			Labels:    map[string]string{"original": "label"},
		},
		Data: map[string]string{"original": "data"},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	updates := &dto.ConfigMapUpdate{}

	result, err := s.svc.UpdateConfigMap(ctx, namespace, "noop-cm", updates)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), map[string]string{"original": "data"}, result.Data)
	assert.Equal(s.T(), map[string]string{"original": "label"}, result.Labels)
}

// --- DeleteConfigMap ---

func (s *ConfigMapServiceTestSuite) TestDeleteConfigMap_Success() {
	ctx := context.Background()
	namespace := "default"

	_, err := s.client.CoreV1().ConfigMaps(namespace).Create(ctx, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "del-cm", Namespace: namespace},
	}, metav1.CreateOptions{})
	s.Require().NoError(err)

	err = s.svc.DeleteConfigMap(ctx, namespace, "del-cm")

	assert.NoError(s.T(), err)

	_, getErr := s.client.CoreV1().ConfigMaps(namespace).Get(ctx, "del-cm", metav1.GetOptions{})
	assert.Error(s.T(), getErr)
}

func (s *ConfigMapServiceTestSuite) TestDeleteConfigMap_NotFound() {
	ctx := context.Background()

	err := s.svc.DeleteConfigMap(ctx, "default", "nonexistent")

	assert.Error(s.T(), err)
}
