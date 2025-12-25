package repositories

import (
	"time"

	"gorm.io/gorm"
)

type CalendarRepository struct {
	db *gorm.DB
}

func NewCalendarRepository(db *gorm.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

type CalendarEvent struct {
	IssueID     uint      `json:"issue_id"`
	IssueTitle  string    `json:"issue_title"`
	Priority    string    `json:"priority"`
	StatusID    *uint     `json:"status_id"`
	StatusName  string    `json:"status_name"`
	StatusColor string    `json:"status_color"`
	UserID      uint      `json:"user_id"`
	UserName    string    `json:"user_name"`
	TeamID      uint      `json:"team_id"`
	TeamName    string    `json:"team_name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

func (r *CalendarRepository) GetCalendarEvents(teamID *uint, userID *uint, startDate, endDate time.Time) ([]CalendarEvent, error) {
	var events []CalendarEvent

	query := r.db.Table("issue_assignments").
		Select(`
			issues.id as issue_id,
			issues.title as issue_title,
			issues.priority,
			issue_statuses.id as status_id,
			issue_statuses.name as status_name,
			issue_statuses.color as status_color,
			users.id as user_id,
			users.full_name as user_name,
			teams.id as team_id,
			teams.name as team_name,
			issue_assignments.start_date,
			issue_assignments.end_date
		`).
		Joins("INNER JOIN issues ON issue_assignments.issue_id = issues.id").
		Joins("LEFT JOIN issue_statuses ON issues.status_id = issue_statuses.id").
		Joins("INNER JOIN users ON issue_assignments.user_id = users.id").
		Joins("INNER JOIN teams ON issues.team_id = teams.id").
		Where("issue_assignments.is_active = true").
		Where("issues.deleted_at IS NULL").
		Where("issue_assignments.end_date >= ? AND issue_assignments.start_date <= ?", startDate, endDate)

	if teamID != nil {
		query = query.Where("issues.team_id = ?", *teamID)
	}

	if userID != nil {
		query = query.Where("issue_assignments.user_id = ?", *userID)
	}

	err := query.Order("issue_assignments.start_date ASC").Scan(&events).Error
	return events, err
}
