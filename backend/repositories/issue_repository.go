package repositories

import (
	"task-management/models"

	"gorm.io/gorm"
)

type IssueRepository struct {
	db *gorm.DB
}

func NewIssueRepository(db *gorm.DB) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) Create(issue *models.Issue) error {
	return r.db.Create(issue).Error
}

func (r *IssueRepository) FindByID(id uint) (*models.Issue, error) {
	var issue models.Issue
	err := r.db.Preload("Status").Preload("Creator").Preload("Team").
		Preload("Assignments.User").Preload("Assignments.AssignedByUser").
		Where("id = ? AND deleted_at IS NULL", id).First(&issue).Error
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (r *IssueRepository) FindByTeam(teamID uint) ([]models.Issue, error) {
	var issues []models.Issue
	err := r.db.Preload("Status").Preload("Creator").
		Preload("Assignments.User").
		Where("team_id = ? AND deleted_at IS NULL", teamID).
		Order("created_at DESC").Find(&issues).Error
	return issues, err
}

func (r *IssueRepository) Update(issue *models.Issue) error {
	return r.db.Save(issue).Error
}

func (r *IssueRepository) Delete(id uint) error {
	return r.db.Model(&models.Issue{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *IssueRepository) CreateActivity(activity *models.IssueActivity) error {
	return r.db.Create(activity).Error
}

func (r *IssueRepository) GetActivities(issueID uint) ([]models.IssueActivity, error) {
	var activities []models.IssueActivity
	err := r.db.Preload("User").Where("issue_id = ?", issueID).Order("created_at DESC").Find(&activities).Error
	return activities, err
}

func (r *IssueRepository) CreateStatusLog(log *models.IssueStatusLog) error {
	return r.db.Create(log).Error
}

func (r *IssueRepository) CreateHoldReason(reason *models.IssueHoldReason) error {
	return r.db.Create(reason).Error
}

func (r *IssueRepository) ResolveHoldReason(issueID, userID uint) error {
	return r.db.Model(&models.IssueHoldReason{}).
		Where("issue_id = ? AND resolved_at IS NULL", issueID).
		Updates(map[string]interface{}{
			"resolved_at": gorm.Expr("CURRENT_TIMESTAMP"),
			"resolved_by": userID,
		}).Error
}

func (r *IssueRepository) CreateWorkLog(log *models.IssueWorkLog) error {
	return r.db.Create(log).Error
}

func (r *IssueRepository) GetWorkLogs(issueID uint) ([]models.IssueWorkLog, error) {
	var logs []models.IssueWorkLog
	err := r.db.Preload("User").Where("issue_id = ?", issueID).Order("work_date DESC").Find(&logs).Error
	return logs, err
}
