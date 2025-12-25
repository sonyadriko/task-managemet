package services

import (
	"task-management/models"
	"task-management/repositories"
)

type IssueService struct {
	issueRepo  *repositories.IssueRepository
	statusRepo *repositories.StatusRepository
}

func NewIssueService(issueRepo *repositories.IssueRepository, statusRepo *repositories.StatusRepository) *IssueService {
	return &IssueService{
		issueRepo:  issueRepo,
		statusRepo: statusRepo,
	}
}

func (s *IssueService) Create(issue *models.Issue, createdBy uint) error {
	issue.CreatedBy = createdBy

	if err := s.issueRepo.Create(issue); err != nil {
		return err
	}

	// Log activity
	activity := &models.IssueActivity{
		IssueID:      issue.ID,
		UserID:       &createdBy,
		ActivityType: models.ActivityCreated,
		Description:  "Issue created",
	}
	return s.issueRepo.CreateActivity(activity)
}

func (s *IssueService) GetByID(id uint) (*models.Issue, error) {
	return s.issueRepo.FindByID(id)
}

func (s *IssueService) GetByTeam(teamID uint) ([]models.Issue, error) {
	return s.issueRepo.FindByTeam(teamID)
}

func (s *IssueService) Update(issue *models.Issue) error {
	return s.issueRepo.Update(issue)
}

func (s *IssueService) Delete(id uint) error {
	return s.issueRepo.Delete(id)
}

func (s *IssueService) UpdateStatus(issueID, newStatusID, userID uint) error {
	issue, err := s.issueRepo.FindByID(issueID)
	if err != nil {
		return err
	}

	oldStatusID := issue.StatusID

	// Update status
	issue.StatusID = &newStatusID
	if err := s.issueRepo.Update(issue); err != nil {
		return err
	}

	// Log status change
	statusLog := &models.IssueStatusLog{
		IssueID:      issueID,
		FromStatusID: oldStatusID,
		ToStatusID:   &newStatusID,
		ChangedBy:    &userID,
	}
	if err := s.issueRepo.CreateStatusLog(statusLog); err != nil {
		return err
	}

	// Log activity
	activity := &models.IssueActivity{
		IssueID:      issueID,
		UserID:       &userID,
		ActivityType: models.ActivityStatusChanged,
		Description:  "Status changed",
	}
	return s.issueRepo.CreateActivity(activity)
}

func (s *IssueService) Hold(issueID, userID uint, reason string) error {
	// Create hold reason
	holdReason := &models.IssueHoldReason{
		IssueID:   issueID,
		Reason:    reason,
		CreatedBy: &userID,
	}
	if err := s.issueRepo.CreateHoldReason(holdReason); err != nil {
		return err
	}

	// Log activity
	activity := &models.IssueActivity{
		IssueID:      issueID,
		UserID:       &userID,
		ActivityType: models.ActivityHold,
		Description:  "Issue put on hold: " + reason,
	}
	return s.issueRepo.CreateActivity(activity)
}

func (s *IssueService) Resume(issueID, userID uint) error {
	// Resolve hold reason
	if err := s.issueRepo.ResolveHoldReason(issueID, userID); err != nil {
		return err
	}

	// Log activity
	activity := &models.IssueActivity{
		IssueID:      issueID,
		UserID:       &userID,
		ActivityType: models.ActivityResumed,
		Description:  "Issue resumed",
	}
	return s.issueRepo.CreateActivity(activity)
}

func (s *IssueService) GetActivities(issueID uint) ([]models.IssueActivity, error) {
	return s.issueRepo.GetActivities(issueID)
}

func (s *IssueService) LogWork(issueID, userID uint, log *models.IssueWorkLog) error {
	log.IssueID = issueID
	log.UserID = userID
	return s.issueRepo.CreateWorkLog(log)
}
