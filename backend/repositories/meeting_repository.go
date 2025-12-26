package repositories

import (
	"task-management/models"
	"time"

	"gorm.io/gorm"
)

type MeetingRepository struct {
	db *gorm.DB
}

func NewMeetingRepository(db *gorm.DB) *MeetingRepository {
	return &MeetingRepository{db: db}
}

func (r *MeetingRepository) Create(meeting *models.Meeting) error {
	return r.db.Create(meeting).Error
}

func (r *MeetingRepository) FindByID(id uint) (*models.Meeting, error) {
	var meeting models.Meeting
	err := r.db.Preload("Team").Preload("Creator").Preload("Attendees.User").First(&meeting, id).Error
	if err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (r *MeetingRepository) FindByTeam(teamID uint) ([]models.Meeting, error) {
	var meetings []models.Meeting
	err := r.db.Preload("Creator").Preload("Attendees.User").
		Where("team_id = ?", teamID).
		Order("meeting_date ASC, start_time ASC").
		Find(&meetings).Error
	return meetings, err
}

func (r *MeetingRepository) FindByDateRange(teamID uint, startDate, endDate time.Time) ([]models.Meeting, error) {
	var meetings []models.Meeting
	err := r.db.Preload("Creator").Preload("Attendees.User").
		Where("team_id = ? AND meeting_date >= ? AND meeting_date <= ?", teamID, startDate, endDate).
		Order("meeting_date ASC, start_time ASC").
		Find(&meetings).Error
	return meetings, err
}

func (r *MeetingRepository) Update(meeting *models.Meeting) error {
	return r.db.Save(meeting).Error
}

func (r *MeetingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Meeting{}, id).Error
}

// Attendee operations
func (r *MeetingRepository) AddAttendee(attendee *models.MeetingAttendee) error {
	return r.db.Create(attendee).Error
}

func (r *MeetingRepository) UpdateAttendeeStatus(meetingID, userID uint, status models.AttendeeStatus) error {
	return r.db.Model(&models.MeetingAttendee{}).
		Where("meeting_id = ? AND user_id = ?", meetingID, userID).
		Update("status", status).Error
}

func (r *MeetingRepository) RemoveAttendee(meetingID, userID uint) error {
	return r.db.Where("meeting_id = ? AND user_id = ?", meetingID, userID).
		Delete(&models.MeetingAttendee{}).Error
}

func (r *MeetingRepository) GetUserMeetings(userID uint, startDate, endDate time.Time) ([]models.Meeting, error) {
	var meetings []models.Meeting
	err := r.db.Preload("Creator").Preload("Attendees.User").
		Joins("JOIN meeting_attendees ON meeting_attendees.meeting_id = meetings.id").
		Where("meeting_attendees.user_id = ? AND meetings.meeting_date >= ? AND meetings.meeting_date <= ?",
			userID, startDate, endDate).
		Order("meeting_date ASC, start_time ASC").
		Find(&meetings).Error
	return meetings, err
}
