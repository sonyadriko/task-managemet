package handlers

import (
	"net/http"
	"strconv"
	"task-management/middleware"
	"task-management/models"
	"task-management/services"

	"github.com/gin-gonic/gin"
)

type IssueHandler struct {
	issueService      *services.IssueService
	assignmentService *services.AssignmentService
	permissionService *services.PermissionService
}

func NewIssueHandler(
	issueService *services.IssueService,
	assignmentService *services.AssignmentService,
	permissionService *services.PermissionService,
) *IssueHandler {
	return &IssueHandler{
		issueService:      issueService,
		assignmentService: assignmentService,
		permissionService: permissionService,
	}
}

func (h *IssueHandler) List(c *gin.Context) {
	teamIDStr := c.Query("team_id")
	if teamIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_id required"})
		return
	}

	teamID, _ := strconv.ParseUint(teamIDStr, 10, 32)
	issues, err := h.issueService.GetByTeam(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issues)
}

func (h *IssueHandler) Create(c *gin.Context) {
	var issue models.Issue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.issueService.Create(&issue, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, issue)
}

func (h *IssueHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	issue, err := h.issueService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}
	c.JSON(http.StatusOK, issue)
}

func (h *IssueHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var issue models.Issue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issue.ID = uint(id)
	if err := h.issueService.Update(&issue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

func (h *IssueHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.issueService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Issue deleted"})
}

func (h *IssueHandler) Assign(c *gin.Context) {
	issueID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req services.AssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.IssueID = uint(issueID)
	userID := middleware.GetUserID(c)

	if err := h.assignmentService.Assign(&req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Issue assigned"})
}

func (h *IssueHandler) UpdateStatus(c *gin.Context) {
	issueID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		StatusID uint `json:"status_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.issueService.UpdateStatus(uint(issueID), req.StatusID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

func (h *IssueHandler) Hold(c *gin.Context) {
	issueID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.issueService.Hold(uint(issueID), userID, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue put on hold"})
}

func (h *IssueHandler) Resume(c *gin.Context) {
	issueID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := middleware.GetUserID(c)

	if err := h.issueService.Resume(uint(issueID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue resumed"})
}

func (h *IssueHandler) GetActivities(c *gin.Context) {
	issueID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	activities, err := h.issueService.GetActivities(uint(issueID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func (h *IssueHandler) LogWork(c *gin.Context) {
	issueID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var log models.IssueWorkLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.issueService.LogWork(uint(issueID), userID, &log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Work logged"})
}
