package services

import (
	"errors"
	"task-management/models"
	"task-management/repositories"
)

type TeamService struct {
	teamRepo *repositories.TeamRepository
	userRepo *repositories.UserRepository
}

func NewTeamService(teamRepo *repositories.TeamRepository, userRepo *repositories.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamService) Create(team *models.Team) error {
	return s.teamRepo.Create(team)
}

func (s *TeamService) GetByID(id uint) (*models.Team, error) {
	return s.teamRepo.FindByID(id)
}

func (s *TeamService) GetByOrganization(orgID uint) ([]models.Team, error) {
	return s.teamRepo.FindByOrganization(orgID)
}

func (s *TeamService) Update(team *models.Team) error {
	return s.teamRepo.Update(team)
}

func (s *TeamService) Delete(id uint) error {
	return s.teamRepo.Delete(id)
}

func (s *TeamService) AddMember(teamID, userID uint, role models.TeamRole) error {
	// Verify user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	member := &models.TeamMember{
		TeamID: teamID,
		UserID: userID,
		Role:   role,
	}

	return s.teamRepo.AddMember(member)
}

func (s *TeamService) RemoveMember(teamID, userID uint) error {
	return s.teamRepo.RemoveMember(teamID, userID)
}

func (s *TeamService) GetMembers(teamID uint) ([]models.TeamMember, error) {
	return s.teamRepo.GetMembers(teamID)
}

func (s *TeamService) GetMemberRole(teamID, userID uint) (*models.TeamMember, error) {
	return s.teamRepo.GetMemberRole(teamID, userID)
}
