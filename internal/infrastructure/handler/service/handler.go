package service

import (
	"net/http"

	"k8s-agent-new/internal/core/dto"
	serviceService "k8s-agent-new/internal/core/service/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *serviceService.Service
}

func NewHandler(service *serviceService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetServices(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	services, err := h.service.ListServices(ctx, namespace)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

func (h *Handler) GetService(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	serviceName := c.Param("name")

	svc, err := h.service.GetService(ctx, namespace, serviceName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, svc)
}

func (h *Handler) CreateService(c *gin.Context) {
	ctx := c.Request.Context()
	var serviceCreate dto.ServiceCreate
	if err := c.ShouldBindJSON(&serviceCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc, err := h.service.CreateService(ctx, &serviceCreate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, svc)
}

func (h *Handler) UpdateService(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	serviceName := c.Param("name")

	var updates dto.ServiceUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc, err := h.service.UpdateService(ctx, namespace, serviceName, &updates)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, svc)
}

func (h *Handler) DeleteService(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	serviceName := c.Param("name")

	err := h.service.DeleteService(ctx, namespace, serviceName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service deleted"})
}
