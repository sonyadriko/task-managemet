package services

import (
	"task-management/repositories"
	"time"
)

type CalendarService struct {
	calendarRepo *repositories.CalendarRepository
}

func NewCalendarService(calendarRepo *repositories.CalendarRepository) *CalendarService {
	return &CalendarService{calendarRepo: calendarRepo}
}

func (s *CalendarService) GetCalendarEvents(teamID *uint, userID *uint, startDate, endDate time.Time) ([]repositories.CalendarEvent, error) {
	return s.calendarRepo.GetCalendarEvents(teamID, userID, startDate, endDate)
}
