package models

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Users []User `gorm:"foreignKey:OrganizationID" json:"users,omitempty"`
	Teams []Team `gorm:"foreignKey:OrganizationID" json:"teams,omitempty"`
}
