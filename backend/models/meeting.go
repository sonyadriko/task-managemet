package models

import "time"

type RecurringPattern string

const (
	RecurringDaily   RecurringPattern = "DAILY"
	RecurringWeekly  RecurringPattern = "WEEKLY"
	RecurringMonthly RecurringPattern = "MONTHLY"
)

type AttendeeStatus string

const (
	AttendeeStatusPending  AttendeeStatus = "PENDING"
	AttendeeStatusAccepted AttendeeStatus = "ACCEPTED"
	AttendeeStatusDeclined AttendeeStatus = "DECLINED"
)

type Meeting struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	TeamID           uint              `gorm:"not null" json:"team_id"`
	Title            string            `gorm:"size:500;not null" json:"title"`
	Description      string            `gorm:"type:text" json:"description"`
	MeetingDate      time.Time         `gorm:"type:date;not null" json:"meeting_date"`
	StartTime        string            `gorm:"type:time;not null" json:"start_time"`
	EndTime          string            `gorm:"type:time;not null" json:"end_time"`
	Location         string            `gorm:"type:text" json:"location"`
	IsRecurring      bool              `gorm:"default:false" json:"is_recurring"`
	RecurringPattern *RecurringPattern `gorm:"type:recurring_pattern" json:"recurring_pattern,omitempty"`
	CreatedBy        uint              `gorm:"not null" json:"created_by"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`

	// Relationships
	Team      Team              `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Creator   User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Attendees []MeetingAttendee `gorm:"foreignKey:MeetingID" json:"attendees,omitempty"`
}

type MeetingAttendee struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MeetingID uint           `gorm:"not null" json:"meeting_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Status    AttendeeStatus `gorm:"type:attendee_status;default:'PENDING'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
