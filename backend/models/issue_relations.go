package models

import (
	"time"
)

type IssueAssignment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	IssueID    uint      `gorm:"not null" json:"issue_id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	StartDate  time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate    time.Time `gorm:"type:date;not null" json:"end_date"`
	AssignedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
	AssignedBy *uint     `json:"assigned_by,omitempty"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`

	// Relationships
	Issue          Issue `gorm:"foreignKey:IssueID" json:"issue,omitempty"`
	User           User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AssignedByUser *User `gorm:"foreignKey:AssignedBy" json:"assigned_by_user,omitempty"`
}

type IssueWorkLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	IssueID      uint      `gorm:"not null" json:"issue_id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	WorkDate     time.Time `gorm:"type:date;not null" json:"work_date"`
	MinutesSpent int       `gorm:"not null" json:"minutes_spent"`
	Notes        string    `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Issue Issue `gorm:"foreignKey:IssueID" json:"issue,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type IssueStatusLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	IssueID      uint      `gorm:"not null" json:"issue_id"`
	FromStatusID *uint     `json:"from_status_id,omitempty"`
	ToStatusID   *uint     `json:"to_status_id,omitempty"`
	ChangedBy    *uint     `json:"changed_by,omitempty"`
	ChangedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"changed_at"`

	// Relationships
	Issue         Issue        `gorm:"foreignKey:IssueID" json:"issue,omitempty"`
	FromStatus    *IssueStatus `gorm:"foreignKey:FromStatusID" json:"from_status,omitempty"`
	ToStatus      *IssueStatus `gorm:"foreignKey:ToStatusID" json:"to_status,omitempty"`
	ChangedByUser *User        `gorm:"foreignKey:ChangedBy" json:"changed_by_user,omitempty"`
}

type ActivityType string

const (
	ActivityCreated         ActivityType = "created"
	ActivityAssigned        ActivityType = "assigned"
	ActivityStatusChanged   ActivityType = "status_changed"
	ActivityPriorityChanged ActivityType = "priority_changed"
	ActivityCommented       ActivityType = "commented"
	ActivityHold            ActivityType = "hold"
	ActivityResumed         ActivityType = "resumed"
)

type IssueActivity struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	IssueID      uint         `gorm:"not null" json:"issue_id"`
	UserID       *uint        `json:"user_id,omitempty"`
	ActivityType ActivityType `gorm:"type:activity_type;not null" json:"activity_type"`
	Description  string       `gorm:"type:text" json:"description"`
	Metadata     *string      `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	CreatedAt    time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Issue Issue `gorm:"foreignKey:IssueID" json:"issue,omitempty"`
	User  *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type IssueHoldReason struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	IssueID    uint       `gorm:"not null" json:"issue_id"`
	Reason     string     `gorm:"type:text;not null" json:"reason"`
	CreatedBy  *uint      `json:"created_by,omitempty"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy *uint      `json:"resolved_by,omitempty"`

	// Relationships
	Issue          Issue `gorm:"foreignKey:IssueID" json:"issue,omitempty"`
	CreatedByUser  *User `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	ResolvedByUser *User `gorm:"foreignKey:ResolvedBy" json:"resolved_by_user,omitempty"`
}
