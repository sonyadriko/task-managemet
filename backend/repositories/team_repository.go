package repositories

import (
	"task-management/models"

	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *TeamRepository) FindByID(id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.Preload("Members.User").Where("id = ? AND deleted_at IS NULL", id).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) FindByOrganization(orgID uint) ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Where("organization_id = ? AND deleted_at IS NULL", orgID).Find(&teams).Error
	return teams, err
}

func (r *TeamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

func (r *TeamRepository) Delete(id uint) error {
	return r.db.Model(&models.Team{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *TeamRepository) AddMember(member *models.TeamMember) error {
	return r.db.Create(member).Error
}

func (r *TeamRepository) RemoveMember(teamID, userID uint) error {
	return r.db.Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&models.TeamMember{}).Error
}

func (r *TeamRepository) GetMemberRole(teamID, userID uint) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.Where("team_id = ? AND user_id = ?", teamID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *TeamRepository) GetMembers(teamID uint) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.db.Preload("User").Where("team_id = ?", teamID).Find(&members).Error
	return members, err
}
