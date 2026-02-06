package pod

import (
	"net/http"

	"k8s-agent-new/internal/core/dto"
	podService "k8s-agent-new/internal/core/service/pod"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *podService.Service
}

func NewHandler(service *podService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetPods(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	pods, err := h.service.ListPods(ctx, namespace)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pods)
}

func (h *Handler) GetPod(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	podName := c.Param("name")
	pod, err := h.service.GetPod(ctx, namespace, podName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pod)
}

func (h *Handler) DeletePod(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	podName := c.Param("name")
	err := h.service.DeletePod(ctx, namespace, podName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pod deleted"})
}

func (h *Handler) RestartPod(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	podName := c.Param("name")
	err := h.service.RestartPod(ctx, namespace, podName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pod restarted"})
}

func (h *Handler) CreatePod(c *gin.Context) {
	ctx := c.Request.Context()
	var podCreate dto.PodCreate
	if err := c.ShouldBindJSON(&podCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pod, err := h.service.CreatePod(ctx, &podCreate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pod)
}

func (h *Handler) UpdatePod(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	podName := c.Param("name")

	var updates dto.PodUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pod, err := h.service.UpdatePod(ctx, namespace, podName, &updates)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pod)
}
