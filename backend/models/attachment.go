package models

import "time"

type Attachment struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	IssueID          uint      `gorm:"not null" json:"issue_id"`
	Filename         string    `gorm:"size:255;not null" json:"filename"`
	OriginalFilename string    `gorm:"size:255;not null" json:"original_filename"`
	FileSize         int64     `gorm:"not null" json:"file_size"`
	MimeType         string    `gorm:"size:100;not null" json:"mime_type"`
	StorageKey       string    `gorm:"size:500;not null" json:"storage_key"`
	UploadedBy       *uint     `json:"uploaded_by,omitempty"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Issue Issue `gorm:"foreignKey:IssueID" json:"issue,omitempty"`
	User  *User `gorm:"foreignKey:UploadedBy" json:"user,omitempty"`
}

func (Attachment) TableName() string {
	return "issue_attachments"
}
