package services

import (
	"errors"
	"task-management/models"
	"task-management/repositories"
)

type PermissionService struct {
	teamRepo *repositories.TeamRepository
}

func NewPermissionService(teamRepo *repositories.TeamRepository) *PermissionService {
	return &PermissionService{teamRepo: teamRepo}
}

// HasTeamAccess checks if a user has the required role or higher in a team
func (s *PermissionService) HasTeamAccess(userID, teamID uint, requiredRole string) (bool, error) {
	member, err := s.teamRepo.GetMemberRole(teamID, userID)
	if err != nil {
		return false, nil // User is not a member
	}

	// Role hierarchy: manager > assistant > member > stakeholder
	roleLevel := map[models.TeamRole]int{
		models.RoleManager:     4,
		models.RoleAssistant:   3,
		models.RoleMember:      2,
		models.RoleStakeholder: 1,
	}

	userLevel := roleLevel[member.Role]
	requiredLevel := roleLevel[models.TeamRole(requiredRole)]

	return userLevel >= requiredLevel, nil
}

// CanEditIssue checks if a user can edit a specific issue
func (s *PermissionService) CanEditIssue(userID, teamID uint, issue *models.Issue) (bool, error) {
	member, err := s.teamRepo.GetMemberRole(teamID, userID)
	if err != nil {
		return false, errors.New("user is not a team member")
	}

	// Managers and assistants can edit any issue
	if member.Role == models.RoleManager || member.Role == models.RoleAssistant {
		return true, nil
	}

	// Members can only edit issues assigned to them
	if member.Role == models.RoleMember {
		for _, assignment := range issue.Assignments {
			if assignment.UserID == userID && assignment.IsActive {
				return true, nil
			}
		}
	}

	// Stakeholders cannot edit
	return false, nil
}
