package networkpolicy

import (
	"net/http"

	"k8s-agent-new/internal/core/dto"
	networkpolicyService "k8s-agent-new/internal/core/service/networkpolicy"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *networkpolicyService.Service
}

func NewHandler(service *networkpolicyService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetNetworkPolicies(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	networkPolicies, err := h.service.ListNetworkPolicies(ctx, namespace)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, networkPolicies)
}

func (h *Handler) GetNetworkPolicy(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	networkPolicyName := c.Param("name")

	networkPolicy, err := h.service.GetNetworkPolicy(ctx, namespace, networkPolicyName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, networkPolicy)
}

func (h *Handler) CreateNetworkPolicy(c *gin.Context) {
	ctx := c.Request.Context()
	var networkPolicyCreate dto.NetworkPolicyCreate
	if err := c.ShouldBindJSON(&networkPolicyCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	networkPolicy, err := h.service.CreateNetworkPolicy(ctx, &networkPolicyCreate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, networkPolicy)
}

func (h *Handler) UpdateNetworkPolicy(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	networkPolicyName := c.Param("name")

	var updates dto.NetworkPolicyUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	networkPolicy, err := h.service.UpdateNetworkPolicy(ctx, namespace, networkPolicyName, &updates)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, networkPolicy)
}

func (h *Handler) DeleteNetworkPolicy(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	networkPolicyName := c.Param("name")

	err := h.service.DeleteNetworkPolicy(ctx, namespace, networkPolicyName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "network policy deleted"})
}
