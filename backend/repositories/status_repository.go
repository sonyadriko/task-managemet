package repositories

import (
	"task-management/models"

	"gorm.io/gorm"
)

type StatusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

func (r *StatusRepository) Create(status *models.IssueStatus) error {
	return r.db.Create(status).Error
}

func (r *StatusRepository) FindByID(id uint) (*models.IssueStatus, error) {
	var status models.IssueStatus
	err := r.db.Where("id = ?", id).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *StatusRepository) FindByTeam(teamID uint) ([]models.IssueStatus, error) {
	var statuses []models.IssueStatus
	err := r.db.Where("team_id = ?", teamID).Order("position ASC").Find(&statuses).Error
	return statuses, err
}

func (r *StatusRepository) Update(status *models.IssueStatus) error {
	return r.db.Save(status).Error
}

func (r *StatusRepository) Delete(id uint) error {
	return r.db.Delete(&models.IssueStatus{}, id).Error
}
