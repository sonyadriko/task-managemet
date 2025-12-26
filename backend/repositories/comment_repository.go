package repositories

import (
	"task-management/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindByIssue(issueID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Preload("User").Where("issue_id = ?", issueID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
