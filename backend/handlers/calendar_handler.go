package handlers

import (
	"net/http"
	"strconv"
	"task-management/services"
	"time"

	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	calendarService   *services.CalendarService
	permissionService *services.PermissionService
}

func NewCalendarHandler(calendarService *services.CalendarService, permissionService *services.PermissionService) *CalendarHandler {
	return &CalendarHandler{
		calendarService:   calendarService,
		permissionService: permissionService,
	}
}

func (h *CalendarHandler) GetCalendar(c *gin.Context) {
	// Parse query parameters
	var teamID *uint
	if teamIDStr := c.Query("team_id"); teamIDStr != "" {
		id, _ := strconv.ParseUint(teamIDStr, 10, 32)
		teamIDVal := uint(id)
		teamID = &teamIDVal
	}

	var userID *uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, _ := strconv.ParseUint(userIDStr, 10, 32)
		userIDVal := uint(id)
		userID = &userIDVal
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format (use YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format (use YYYY-MM-DD)"})
		return
	}

	// Get calendar events
	events, err := h.calendarService.GetCalendarEvents(teamID, userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
