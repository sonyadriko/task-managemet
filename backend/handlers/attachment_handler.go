package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"task-management/middleware"
	"task-management/models"
	"task-management/repositories"
	"task-management/services"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	attachmentRepo *repositories.AttachmentRepository
	storageService *services.StorageService
}

func NewAttachmentHandler(attachmentRepo *repositories.AttachmentRepository, storageService *services.StorageService) *AttachmentHandler {
	return &AttachmentHandler{
		attachmentRepo: attachmentRepo,
		storageService: storageService,
	}
}

// Upload handles file upload
func (h *AttachmentHandler) Upload(c *gin.Context) {
	issueID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file size (max 10MB)
	const maxSize = 10 << 20 // 10MB
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 10MB"})
		return
	}

	// Get MIME type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Validate file type (only images, documents, text)
	allowedTypes := map[string]bool{
		// Images
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/svg+xml": true,
		// Documents
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
		"application/vnd.ms-powerpoint":                                             true,
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
		// Text
		"text/plain":       true,
		"text/csv":         true,
		"text/html":        true,
		"text/markdown":    true,
		"application/json": true,
		"application/xml":  true,
	}

	if !allowedTypes[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed. Only images, documents, and text files are accepted."})
		return
	}

	// Upload to R2
	ctx := context.Background()
	storageKey, err := h.storageService.Upload(ctx, file, header.Filename, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// Save attachment record
	userID := middleware.GetUserID(c)
	attachment := &models.Attachment{
		IssueID:          uint(issueID),
		Filename:         storageKey,
		OriginalFilename: header.Filename,
		FileSize:         header.Size,
		MimeType:         mimeType,
		StorageKey:       storageKey,
		UploadedBy:       &userID,
	}

	if err := h.attachmentRepo.Create(attachment); err != nil {
		// Try to delete uploaded file if DB save fails
		h.storageService.Delete(ctx, storageKey)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save attachment record"})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

// List returns all attachments for an issue
func (h *AttachmentHandler) List(c *gin.Context) {
	issueID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	attachments, err := h.attachmentRepo.FindByIssue(uint(issueID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attachments"})
		return
	}

	c.JSON(http.StatusOK, attachments)
}

// Download generates a presigned URL for downloading
func (h *AttachmentHandler) Download(c *gin.Context) {
	attachmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
		return
	}

	attachment, err := h.attachmentRepo.FindByID(uint(attachmentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	ctx := context.Background()
	url, err := h.storageService.GetPresignedURL(ctx, attachment.StorageKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":      url,
		"filename": attachment.OriginalFilename,
	})
}

// Delete removes an attachment
func (h *AttachmentHandler) Delete(c *gin.Context) {
	attachmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
		return
	}

	attachment, err := h.attachmentRepo.FindByID(uint(attachmentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	// Delete from R2
	ctx := context.Background()
	if err := h.storageService.Delete(ctx, attachment.StorageKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from storage"})
		return
	}

	// Delete from database
	if err := h.attachmentRepo.Delete(uint(attachmentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attachment record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attachment deleted"})
}
