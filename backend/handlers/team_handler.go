package handlers

import (
	"net/http"
	"strconv"
	"task-management/middleware"
	"task-management/models"
	"task-management/services"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService       *services.TeamService
	permissionService *services.PermissionService
}

func NewTeamHandler(teamService *services.TeamService, permissionService *services.PermissionService) *TeamHandler {
	return &TeamHandler{
		teamService:       teamService,
		permissionService: permissionService,
	}
}

func (h *TeamHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	teams, err := h.teamService.GetByOrganization(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) Create(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team.OrganizationID = middleware.GetOrganizationID(c)
	if err := h.teamService.Create(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	team, err := h.teamService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var team models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team.ID = uint(id)
	if err := h.teamService.Update(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.teamService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted"})
}

func (h *TeamHandler) GetMembers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	members, err := h.teamService.GetMembers(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *TeamHandler) AddMember(c *gin.Context) {
	teamID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Email    string          `json:"email" binding:"required"`
		FullName string          `json:"full_name"`
		Role     models.TeamRole `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Try to find existing user
	user, err := h.teamService.FindUserByEmail(req.Email)
	if err != nil {
		// User doesn't exist - create new user if fullName provided
		if req.FullName == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "User not found. Provide full_name to create new user.",
				"new_user": true,
			})
			return
		}

		// Get organization ID from the team
		team, err := h.teamService.GetByID(uint(teamID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Team not found"})
			return
		}

		// Create new user with default password
		newUser, err := h.teamService.CreateUser(req.Email, req.FullName, team.OrganizationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}
		user = newUser
	}

	// Add user to team
	if err := h.teamService.AddMember(uint(teamID), user.ID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":          "Member added successfully",
		"user_id":          user.ID,
		"default_password": "password123",
	})
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID, _ := strconv.ParseUint(c.Param("userId"), 10, 32)

	if err := h.teamService.RemoveMember(uint(teamID), uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}
