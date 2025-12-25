package repositories

import (
	"task-management/models"

	"gorm.io/gorm"
)

type AssignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

func (r *AssignmentRepository) Create(assignment *models.IssueAssignment) error {
	return r.db.Create(assignment).Error
}

func (r *AssignmentRepository) FindByIssue(issueID uint) ([]models.IssueAssignment, error) {
	var assignments []models.IssueAssignment
	err := r.db.Preload("User").Preload("AssignedByUser").
		Where("issue_id = ?", issueID).
		Order("assigned_at DESC").Find(&assignments).Error
	return assignments, err
}

func (r *AssignmentRepository) FindActiveByUser(userID uint) ([]models.IssueAssignment, error) {
	var assignments []models.IssueAssignment
	err := r.db.Preload("Issue.Status").Preload("Issue.Team").
		Where("user_id = ? AND is_active = true", userID).
		Order("start_date ASC").Find(&assignments).Error
	return assignments, err
}

func (r *AssignmentRepository) DeactivateOldAssignments(issueID uint) error {
	return r.db.Model(&models.IssueAssignment{}).
		Where("issue_id = ? AND is_active = true", issueID).
		Update("is_active", false).Error
}
