package ingress

import (
	"net/http"

	"k8s-agent-new/internal/core/dto"
	ingressService "k8s-agent-new/internal/core/service/ingress"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *ingressService.Service
}

func NewHandler(service *ingressService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetIngresses(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	ingresses, err := h.service.ListIngresses(ctx, namespace)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingresses)
}

func (h *Handler) GetIngress(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	ingressName := c.Param("name")

	ingress, err := h.service.GetIngress(ctx, namespace, ingressName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingress)
}

func (h *Handler) CreateIngress(c *gin.Context) {
	ctx := c.Request.Context()
	var ingressCreate dto.IngressCreate
	if err := c.ShouldBindJSON(&ingressCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingress, err := h.service.CreateIngress(ctx, &ingressCreate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ingress)
}

func (h *Handler) UpdateIngress(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	ingressName := c.Param("name")

	var updates dto.IngressUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingress, err := h.service.UpdateIngress(ctx, namespace, ingressName, &updates)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingress)
}

func (h *Handler) DeleteIngress(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	ingressName := c.Param("name")

	err := h.service.DeleteIngress(ctx, namespace, ingressName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ingress deleted"})
}
