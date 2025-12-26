package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IssueID   uint      `gorm:"not null" json:"issue_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Comment) TableName() string {
	return "issue_comments"
}
