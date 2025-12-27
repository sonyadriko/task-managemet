package handlers

import (
	"net/http"

	"task-management/middleware"
	"task-management/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

type TeamPermission struct {
	TeamID    uint   `json:"team_id"`
	TeamName  string `json:"team_name"`
	Role      string `json:"role"`
	CanEdit   bool   `json:"can_edit"`
	CanDelete bool   `json:"can_delete"`
	CanManage bool   `json:"can_manage"`
}

type UserPermissions struct {
	UserID     uint             `json:"user_id"`
	Email      string           `json:"email"`
	FullName   string           `json:"full_name"`
	Teams      []TeamPermission `json:"teams"`
	IsOrgAdmin bool             `json:"is_org_admin"`
}

func (h *UserHandler) GetMyPermissions(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	permissions := UserPermissions{
		UserID:     user.ID,
		Email:      user.Email,
		FullName:   user.FullName,
		Teams:      []TeamPermission{},
		IsOrgAdmin: false, // Can be extended for org-level admins
	}

	// Get all team memberships
	var memberships []models.TeamMember
	h.db.Preload("Team").Where("user_id = ?", userID).Find(&memberships)

	for _, m := range memberships {
		teamPerm := TeamPermission{
			TeamID:   m.TeamID,
			TeamName: m.Team.Name,
			Role:     string(m.Role),
		}

		// Set permissions based on role
		switch m.Role {
		case models.RoleManager:
			teamPerm.CanEdit = true
			teamPerm.CanDelete = true
			teamPerm.CanManage = true
		case models.RoleAssistant:
			teamPerm.CanEdit = true
			teamPerm.CanDelete = true
			teamPerm.CanManage = false
		case models.RoleMember:
			teamPerm.CanEdit = true // Can edit assigned only (checked per-issue)
			teamPerm.CanDelete = false
			teamPerm.CanManage = false
		case models.RoleStakeholder:
			teamPerm.CanEdit = false
			teamPerm.CanDelete = false
			teamPerm.CanManage = false
		}

		permissions.Teams = append(permissions.Teams, teamPerm)
	}

	c.JSON(http.StatusOK, permissions)
}

func (h *UserHandler) GetTeamRole(c *gin.Context, teamID uint) *models.TeamMember {
	userID := middleware.GetUserID(c)

	var member models.TeamMember
	if err := h.db.Where("user_id = ? AND team_id = ?", userID, teamID).First(&member).Error; err != nil {
		return nil
	}
	return &member
}
