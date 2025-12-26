package models

import (
	"time"

	"gorm.io/gorm"
)

type IssuePriority string

const (
	PriorityLow    IssuePriority = "LOW"
	PriorityNormal IssuePriority = "NORMAL"
	PriorityHigh   IssuePriority = "HIGH"
	PriorityUrgent IssuePriority = "URGENT"
)

type Issue struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TeamID      uint           `gorm:"not null" json:"team_id"`
	StatusID    *uint          `json:"status_id,omitempty"`
	Title       string         `gorm:"size:500;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Priority    IssuePriority  `gorm:"type:issue_priority;default:NORMAL" json:"priority"`
	Deadline    *time.Time     `gorm:"type:date" json:"deadline,omitempty"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Team        Team              `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Status      *IssueStatus      `gorm:"foreignKey:StatusID" json:"status,omitempty"`
	Creator     User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Assignments []IssueAssignment `gorm:"foreignKey:IssueID" json:"assignments,omitempty"`
	WorkLogs    []IssueWorkLog    `gorm:"foreignKey:IssueID" json:"work_logs,omitempty"`
	Activities  []IssueActivity   `gorm:"foreignKey:IssueID" json:"activities,omitempty"`
	HoldReasons []IssueHoldReason `gorm:"foreignKey:IssueID" json:"hold_reasons,omitempty"`
}

type IssueStatus struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrganizationID uint      `gorm:"not null" json:"organization_id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Position       int       `gorm:"not null" json:"position"`
	IsFinal        bool      `gorm:"default:false" json:"is_final"`
	Color          string    `gorm:"size:7;default:#6B7280" json:"color"`
	CreatedAt      time.Time `json:"created_at"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Issues       []Issue      `gorm:"foreignKey:StatusID" json:"issues,omitempty"`
}
