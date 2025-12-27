package handlers

import (
	"net/http"
	"time"

	"task-management/middleware"
	"task-management/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsHandler struct {
	db *gorm.DB
}

func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	return &AnalyticsHandler{db: db}
}

type DashboardAnalytics struct {
	TotalTasks       int64            `json:"total_tasks"`
	CompletedTasks   int64            `json:"completed_tasks"`
	InProgressTasks  int64            `json:"in_progress_tasks"`
	OnHoldTasks      int64            `json:"on_hold_tasks"`
	OverdueTasks     int64            `json:"overdue_tasks"`
	TasksByStatus    []StatusCount    `json:"tasks_by_status"`
	TasksByPriority  []PriorityCount  `json:"tasks_by_priority"`
	WeeklyActivity   []DailyCount     `json:"weekly_activity"`
	RecentActivities []RecentActivity `json:"recent_activities"`
	TeamStats        []TeamStat       `json:"team_stats"`
}

type StatusCount struct {
	StatusID   uint   `json:"status_id"`
	StatusName string `json:"status_name"`
	Color      string `json:"color"`
	Count      int64  `json:"count"`
}

type PriorityCount struct {
	Priority string `json:"priority"`
	Count    int64  `json:"count"`
}

type DailyCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type RecentActivity struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

type TeamStat struct {
	TeamID      uint   `json:"team_id"`
	TeamName    string `json:"team_name"`
	TaskCount   int64  `json:"task_count"`
	MemberCount int64  `json:"member_count"`
}

func (h *AnalyticsHandler) GetDashboardAnalytics(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orgID := middleware.GetOrganizationID(c)

	analytics := DashboardAnalytics{}

	// Get all team IDs user belongs to
	var teamIDs []uint
	h.db.Model(&models.TeamMember{}).Where("user_id = ?", userID).Pluck("team_id", &teamIDs)

	if len(teamIDs) == 0 {
		c.JSON(http.StatusOK, analytics)
		return
	}

	// Total tasks
	h.db.Model(&models.Issue{}).Where("team_id IN ? AND deleted_at IS NULL", teamIDs).Count(&analytics.TotalTasks)

	// Completed tasks (assuming status with name containing 'done' or 'complete')
	var doneStatusIDs []uint
	h.db.Model(&models.IssueStatus{}).Where("organization_id = ? AND (LOWER(name) LIKE '%done%' OR LOWER(name) LIKE '%complete%')", orgID).Pluck("id", &doneStatusIDs)
	if len(doneStatusIDs) > 0 {
		h.db.Model(&models.Issue{}).Where("team_id IN ? AND status_id IN ? AND deleted_at IS NULL", teamIDs, doneStatusIDs).Count(&analytics.CompletedTasks)
	}

	// In progress tasks
	var progressStatusIDs []uint
	h.db.Model(&models.IssueStatus{}).Where("organization_id = ? AND LOWER(name) LIKE '%progress%'", orgID).Pluck("id", &progressStatusIDs)
	if len(progressStatusIDs) > 0 {
		h.db.Model(&models.Issue{}).Where("team_id IN ? AND status_id IN ? AND deleted_at IS NULL", teamIDs, progressStatusIDs).Count(&analytics.InProgressTasks)
	}

	// On hold tasks
	h.db.Model(&models.Issue{}).Where("team_id IN ? AND is_on_hold = true AND deleted_at IS NULL", teamIDs).Count(&analytics.OnHoldTasks)

	// Overdue tasks
	h.db.Model(&models.Issue{}).Where("team_id IN ? AND deadline < ? AND deleted_at IS NULL AND status_id NOT IN ?", teamIDs, time.Now(), doneStatusIDs).Count(&analytics.OverdueTasks)

	// Tasks by status
	var statusCounts []struct {
		StatusID uint
		Count    int64
	}
	h.db.Model(&models.Issue{}).
		Select("status_id, COUNT(*) as count").
		Where("team_id IN ? AND deleted_at IS NULL", teamIDs).
		Group("status_id").
		Find(&statusCounts)

	for _, sc := range statusCounts {
		var status models.IssueStatus
		h.db.First(&status, sc.StatusID)
		analytics.TasksByStatus = append(analytics.TasksByStatus, StatusCount{
			StatusID:   sc.StatusID,
			StatusName: status.Name,
			Color:      status.Color,
			Count:      sc.Count,
		})
	}

	// Tasks by priority
	var priorityCounts []struct {
		Priority string
		Count    int64
	}
	h.db.Model(&models.Issue{}).
		Select("priority, COUNT(*) as count").
		Where("team_id IN ? AND deleted_at IS NULL", teamIDs).
		Group("priority").
		Find(&priorityCounts)

	for _, pc := range priorityCounts {
		analytics.TasksByPriority = append(analytics.TasksByPriority, PriorityCount{
			Priority: pc.Priority,
			Count:    pc.Count,
		})
	}

	// Weekly activity (tasks created in last 7 days)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.AddDate(0, 0, 1)

		var count int64
		h.db.Model(&models.Issue{}).
			Where("team_id IN ? AND created_at >= ? AND created_at < ? AND deleted_at IS NULL", teamIDs, startOfDay, endOfDay).
			Count(&count)

		analytics.WeeklyActivity = append(analytics.WeeklyActivity, DailyCount{
			Date:  startOfDay.Format("Mon"),
			Count: count,
		})
	}

	// Team stats
	for _, teamID := range teamIDs {
		var team models.Team
		h.db.First(&team, teamID)

		var taskCount int64
		h.db.Model(&models.Issue{}).Where("team_id = ? AND deleted_at IS NULL", teamID).Count(&taskCount)

		var memberCount int64
		h.db.Model(&models.TeamMember{}).Where("team_id = ?", teamID).Count(&memberCount)

		analytics.TeamStats = append(analytics.TeamStats, TeamStat{
			TeamID:      teamID,
			TeamName:    team.Name,
			TaskCount:   taskCount,
			MemberCount: memberCount,
		})
	}

	c.JSON(http.StatusOK, analytics)
}
