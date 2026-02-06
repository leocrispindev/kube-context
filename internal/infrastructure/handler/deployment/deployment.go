package deployment

import (
	"net/http"

	"k8s-agent-new/internal/core/dto"
	deploymentService "k8s-agent-new/internal/core/service/deployment"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *deploymentService.Service
}

func NewHandler(service *deploymentService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetDeployments(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	deployments, err := h.service.ListDeployments(ctx, namespace)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployments)
}

func (h *Handler) UpdateDeployment(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	deploymentName := c.Param("name")

	var updates dto.DeploymentUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateDeployment(ctx, namespace, deploymentName, updates)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deployment updated"})
}

func (h *Handler) GetRolloutStatus(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	deploymentName := c.Param("name")

	status, err := h.service.GetRolloutStatus(ctx, namespace, deploymentName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *Handler) TogglePauseDeployment(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	deploymentName := c.Param("name")

	var body struct {
		Pause bool `json:"pause"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.TogglePauseDeployment(ctx, namespace, deploymentName, body.Pause)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := "deployment resumed"
	if body.Pause {
		message = "deployment paused"
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *Handler) GetDeployment(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	deploymentName := c.Param("name")

	deployment, err := h.service.GetDeployment(ctx, namespace, deploymentName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployment)
}

func (h *Handler) CreateDeployment(c *gin.Context) {
	ctx := c.Request.Context()
	var deployCreate dto.DeploymentCreate
	if err := c.ShouldBindJSON(&deployCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deployment, err := h.service.CreateDeployment(ctx, &deployCreate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deployment)
}

func (h *Handler) DeleteDeployment(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	deploymentName := c.Param("name")

	err := h.service.DeleteDeployment(ctx, namespace, deploymentName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deployment deleted"})
}
