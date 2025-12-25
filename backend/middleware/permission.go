package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionChecker interface {
	HasTeamAccess(userID, teamID uint, requiredRole string) (bool, error)
}

func RequireTeamRole(permissionService PermissionChecker, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get team ID from request (could be from path, query, or body)
		teamIDParam := c.Param("teamId")
		if teamIDParam == "" {
			teamIDParam = c.Query("team_id")
		}

		if teamIDParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID required"})
			c.Abort()
			return
		}

		var teamID uint
		if _, err := fmt.Sscanf(teamIDParam, "%d", &teamID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
			c.Abort()
			return
		}

		hasAccess, err := permissionService.HasTeamAccess(userID, teamID, requiredRole)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
