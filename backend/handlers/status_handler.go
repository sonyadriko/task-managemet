package handlers

import (
	"net/http"
	"strconv"
	"task-management/models"
	"task-management/repositories"
	"task-management/services"

	"github.com/gin-gonic/gin"
)

type StatusHandler struct {
	statusRepo        *repositories.StatusRepository
	permissionService *services.PermissionService
}

func NewStatusHandler(statusRepo *repositories.StatusRepository, permissionService *services.PermissionService) *StatusHandler {
	return &StatusHandler{
		statusRepo:        statusRepo,
		permissionService: permissionService,
	}
}

func (h *StatusHandler) GetByTeam(c *gin.Context) {
	teamID, _ := strconv.ParseUint(c.Param("teamId"), 10, 32)
	statuses, err := h.statusRepo.FindByTeam(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, statuses)
}

func (h *StatusHandler) Create(c *gin.Context) {
	var status models.IssueStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.statusRepo.Create(&status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, status)
}

func (h *StatusHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var status models.IssueStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status.ID = uint(id)
	if err := h.statusRepo.Update(&status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *StatusHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.statusRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status deleted"})
}
