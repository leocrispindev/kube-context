package namespace

import (
	"net/http"

	namespaceService "k8s-agent-new/internal/core/service/namespace"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *namespaceService.Service
}

func NewHandler(service *namespaceService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetNamespaces(c *gin.Context) {
	ctx := c.Request.Context()
	namespaces, err := h.service.ListNamespaces(ctx)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespaces)
}

func (h *Handler) GetNamespace(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	namespace, err := h.service.GetNamespace(ctx, name)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespace)
}

func (h *Handler) CreateNamespace(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Name        string            `json:"name"`
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateNamespace(ctx, req.Name, req.Labels, req.Annotations)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) DeleteNamespace(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	err := h.service.DeleteNamespace(ctx, name)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateNamespace(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	var req struct {
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateNamespace(ctx, name, req.Labels, req.Annotations)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
