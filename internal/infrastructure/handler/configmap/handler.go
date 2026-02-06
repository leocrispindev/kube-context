package configmap

import (
	"net/http"

	"k8s-agent-new/internal/core/dto"
	configmapService "k8s-agent-new/internal/core/service/configmap"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *configmapService.Service
}

func NewHandler(service *configmapService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetConfigMaps(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	configMaps, err := h.service.ListConfigMaps(ctx, namespace)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configMaps)
}

func (h *Handler) GetConfigMap(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	configMapName := c.Param("name")

	configMap, err := h.service.GetConfigMap(ctx, namespace, configMapName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configMap)
}

func (h *Handler) CreateConfigMap(c *gin.Context) {
	ctx := c.Request.Context()
	var configMapCreate dto.ConfigMapCreate
	if err := c.ShouldBindJSON(&configMapCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configMap, err := h.service.CreateConfigMap(ctx, &configMapCreate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, configMap)
}

func (h *Handler) UpdateConfigMap(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	configMapName := c.Param("name")

	var updates dto.ConfigMapUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configMap, err := h.service.UpdateConfigMap(ctx, namespace, configMapName, &updates)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configMap)
}

func (h *Handler) DeleteConfigMap(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	configMapName := c.Param("name")

	err := h.service.DeleteConfigMap(ctx, namespace, configMapName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "configmap deleted"})
}
