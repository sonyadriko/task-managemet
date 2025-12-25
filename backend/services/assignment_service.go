package services

import (
	"errors"
	"task-management/models"
	"task-management/repositories"
	"time"
)

type AssignmentService struct {
	assignmentRepo *repositories.AssignmentRepository
	issueRepo      *repositories.IssueRepository
	userRepo       *repositories.UserRepository
}

func NewAssignmentService(
	assignmentRepo *repositories.AssignmentRepository,
	issueRepo *repositories.IssueRepository,
	userRepo *repositories.UserRepository,
) *AssignmentService {
	return &AssignmentService{
		assignmentRepo: assignmentRepo,
		issueRepo:      issueRepo,
		userRepo:       userRepo,
	}
}

type AssignmentRequest struct {
	IssueID   uint      `json:"issue_id" binding:"required"`
	UserID    uint      `json:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

func (s *AssignmentService) Assign(req *AssignmentRequest, assignedBy uint) error {
	// Verify issue exists
	_, err := s.issueRepo.FindByID(req.IssueID)
	if err != nil {
		return errors.New("issue not found")
	}

	// Verify user exists
	_, err = s.userRepo.FindByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		return errors.New("end date must be after start date")
	}

	// Deactivate old assignments (optional: keep history)
	// s.assignmentRepo.DeactivateOldAssignments(req.IssueID)

	// Create new assignment
	assignment := &models.IssueAssignment{
		IssueID:    req.IssueID,
		UserID:     req.UserID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		AssignedBy: &assignedBy,
		IsActive:   true,
	}

	if err := s.assignmentRepo.Create(assignment); err != nil {
		return err
	}

	// Log activity
	activity := &models.IssueActivity{
		IssueID:      req.IssueID,
		UserID:       &assignedBy,
		ActivityType: models.ActivityAssigned,
		Description:  "Issue assigned",
	}
	return s.issueRepo.CreateActivity(activity)
}

func (s *AssignmentService) GetByIssue(issueID uint) ([]models.IssueAssignment, error) {
	return s.assignmentRepo.FindByIssue(issueID)
}

func (s *AssignmentService) GetActiveByUser(userID uint) ([]models.IssueAssignment, error) {
	return s.assignmentRepo.FindActiveByUser(userID)
}
