package handlers

import (
	"net/http"
	"strconv"
	"task-management/middleware"
	"task-management/models"
	"task-management/repositories"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentRepo *repositories.CommentRepository
}

func NewCommentHandler(commentRepo *repositories.CommentRepository) *CommentHandler {
	return &CommentHandler{commentRepo: commentRepo}
}

// Create adds a new comment to an issue
func (h *CommentHandler) Create(c *gin.Context) {
	issueID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	userID := middleware.GetUserID(c)

	comment := &models.Comment{
		IssueID: uint(issueID),
		UserID:  userID,
		Content: input.Content,
	}

	if err := h.commentRepo.Create(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Fetch with user info
	created, _ := h.commentRepo.FindByID(comment.ID)
	if created != nil {
		c.JSON(http.StatusCreated, created)
	} else {
		c.JSON(http.StatusCreated, comment)
	}
}

// List returns all comments for an issue
func (h *CommentHandler) List(c *gin.Context) {
	issueID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	comments, err := h.commentRepo.FindByIssue(uint(issueID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// Update modifies a comment (only by owner)
func (h *CommentHandler) Update(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.commentRepo.FindByID(uint(commentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check ownership
	userID := middleware.GetUserID(c)
	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own comments"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	comment.Content = input.Content
	if err := h.commentRepo.Update(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// Delete removes a comment (only by owner)
func (h *CommentHandler) Delete(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.commentRepo.FindByID(uint(commentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check ownership
	userID := middleware.GetUserID(c)
	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
		return
	}

	if err := h.commentRepo.Delete(uint(commentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
