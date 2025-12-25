package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrganizationID uint           `gorm:"not null" json:"organization_id"`
	Email          string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash   string         `gorm:"size:255;not null" json:"-"`
	FullName       string         `gorm:"size:255;not null" json:"full_name"`
	Timezone       string         `gorm:"size:50;default:Asia/Jakarta" json:"timezone"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Organization  Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	TeamMembers   []TeamMember `gorm:"foreignKey:UserID" json:"team_members,omitempty"`
	CreatedIssues []Issue      `gorm:"foreignKey:CreatedBy" json:"created_issues,omitempty"`
}
