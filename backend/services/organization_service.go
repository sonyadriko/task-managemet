package services

import (
	"task-management/models"
	"task-management/repositories"
)

type OrganizationService struct {
	orgRepo *repositories.OrganizationRepository
}

func NewOrganizationService(orgRepo *repositories.OrganizationRepository) *OrganizationService {
	return &OrganizationService{orgRepo: orgRepo}
}

func (s *OrganizationService) Create(org *models.Organization) error {
	return s.orgRepo.Create(org)
}

func (s *OrganizationService) GetByID(id uint) (*models.Organization, error) {
	return s.orgRepo.FindByID(id)
}

func (s *OrganizationService) GetAll() ([]models.Organization, error) {
	return s.orgRepo.FindAll()
}

func (s *OrganizationService) Update(org *models.Organization) error {
	return s.orgRepo.Update(org)
}

func (s *OrganizationService) Delete(id uint) error {
	return s.orgRepo.Delete(id)
}
