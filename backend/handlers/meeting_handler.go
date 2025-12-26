package handlers

import (
	"net/http"
	"strconv"
	"task-management/middleware"
	"task-management/models"
	"task-management/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type MeetingHandler struct {
	meetingRepo *repositories.MeetingRepository
}

func NewMeetingHandler(meetingRepo *repositories.MeetingRepository) *MeetingHandler {
	return &MeetingHandler{meetingRepo: meetingRepo}
}

type CreateMeetingRequest struct {
	TeamID           uint                     `json:"team_id" binding:"required"`
	Title            string                   `json:"title" binding:"required"`
	Description      string                   `json:"description"`
	MeetingDate      string                   `json:"meeting_date" binding:"required"`
	StartTime        string                   `json:"start_time" binding:"required"`
	EndTime          string                   `json:"end_time" binding:"required"`
	Location         string                   `json:"location"`
	IsRecurring      bool                     `json:"is_recurring"`
	RecurringPattern *models.RecurringPattern `json:"recurring_pattern"`
	AttendeeIDs      []uint                   `json:"attendee_ids"`
}

func (h *MeetingHandler) Create(c *gin.Context) {
	var req CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meetingDate, err := time.Parse("2006-01-02", req.MeetingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting date format"})
		return
	}

	userID := middleware.GetUserID(c)

	meeting := &models.Meeting{
		TeamID:           req.TeamID,
		Title:            req.Title,
		Description:      req.Description,
		MeetingDate:      meetingDate,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Location:         req.Location,
		IsRecurring:      req.IsRecurring,
		RecurringPattern: req.RecurringPattern,
		CreatedBy:        userID,
	}

	if err := h.meetingRepo.Create(meeting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meeting"})
		return
	}

	// Add attendees
	for _, attendeeID := range req.AttendeeIDs {
		attendee := &models.MeetingAttendee{
			MeetingID: meeting.ID,
			UserID:    attendeeID,
			Status:    models.AttendeeStatusPending,
		}
		h.meetingRepo.AddAttendee(attendee)
	}

	// Fetch complete meeting with attendees
	created, _ := h.meetingRepo.FindByID(meeting.ID)
	c.JSON(http.StatusCreated, created)
}

func (h *MeetingHandler) List(c *gin.Context) {
	teamIDStr := c.Query("team_id")
	if teamIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_id is required"})
		return
	}
	teamID, _ := strconv.ParseUint(teamIDStr, 10, 32)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var meetings []models.Meeting
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		endDate, _ := time.Parse("2006-01-02", endDateStr)
		meetings, err = h.meetingRepo.FindByDateRange(uint(teamID), startDate, endDate)
	} else {
		meetings, err = h.meetingRepo.FindByTeam(uint(teamID))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meetings"})
		return
	}

	c.JSON(http.StatusOK, meetings)
}

func (h *MeetingHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	meeting, err := h.meetingRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}
	c.JSON(http.StatusOK, meeting)
}

func (h *MeetingHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	existing, err := h.meetingRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	var req CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meetingDate, _ := time.Parse("2006-01-02", req.MeetingDate)

	existing.Title = req.Title
	existing.Description = req.Description
	existing.MeetingDate = meetingDate
	existing.StartTime = req.StartTime
	existing.EndTime = req.EndTime
	existing.Location = req.Location
	existing.IsRecurring = req.IsRecurring
	existing.RecurringPattern = req.RecurringPattern

	if err := h.meetingRepo.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meeting"})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *MeetingHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.meetingRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meeting"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Meeting deleted"})
}

// Attendee actions
func (h *MeetingHandler) AddAttendee(c *gin.Context) {
	meetingID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendee := &models.MeetingAttendee{
		MeetingID: uint(meetingID),
		UserID:    req.UserID,
		Status:    models.AttendeeStatusPending,
	}

	if err := h.meetingRepo.AddAttendee(attendee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attendee added"})
}

func (h *MeetingHandler) RespondToMeeting(c *gin.Context) {
	meetingID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := middleware.GetUserID(c)

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var status models.AttendeeStatus
	switch req.Status {
	case "ACCEPTED":
		status = models.AttendeeStatusAccepted
	case "DECLINED":
		status = models.AttendeeStatusDeclined
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if err := h.meetingRepo.UpdateAttendeeStatus(uint(meetingID), userID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Response recorded"})
}
