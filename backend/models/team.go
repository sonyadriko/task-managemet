package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrganizationID uint           `gorm:"not null" json:"organization_id"`
	ParentTeamID   *uint          `json:"parent_team_id,omitempty"`
	Name           string         `gorm:"size:255;not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Organization Organization  `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	ParentTeam   *Team         `gorm:"foreignKey:ParentTeamID" json:"parent_team,omitempty"`
	SubTeams     []Team        `gorm:"foreignKey:ParentTeamID" json:"sub_teams,omitempty"`
	Members      []TeamMember  `gorm:"foreignKey:TeamID" json:"members,omitempty"`
	Issues       []Issue       `gorm:"foreignKey:TeamID" json:"issues,omitempty"`
	Statuses     []IssueStatus `gorm:"foreignKey:TeamID" json:"statuses,omitempty"`
}

type TeamRole string

const (
	RoleStakeholder TeamRole = "stakeholder"
	RoleManager     TeamRole = "manager"
	RoleAssistant   TeamRole = "assistant"
	RoleMember      TeamRole = "member"
)

type TeamMember struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	TeamID   uint      `gorm:"not null" json:"team_id"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	Role     TeamRole  `gorm:"type:team_role;default:member;not null" json:"role"`
	JoinedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`

	// Relationships
	Team Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
