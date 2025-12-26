package repositories

import (
	"task-management/models"

	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(attachment *models.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *AttachmentRepository) FindByIssue(issueID uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.Preload("User").Where("issue_id = ?", issueID).Order("created_at DESC").Find(&attachments).Error
	return attachments, err
}

func (r *AttachmentRepository) FindByID(id uint) (*models.Attachment, error) {
	var attachment models.Attachment
	err := r.db.First(&attachment, id).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *AttachmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Attachment{}, id).Error
}
