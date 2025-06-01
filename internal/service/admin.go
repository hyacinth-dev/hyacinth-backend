package service

import (
	// "context"
	// "hyacinth-backend/internal/model"
	// v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/repository"
)

type AdminService interface {
}

func NewAdminService(
	service *Service,
	userRepo repository.UserRepository,
	usageRepo repository.UsageRepository,
) AdminService {
	return &adminService{
		Service:   service,
		userRepo:  userRepo,
		usageRepo: usageRepo,
	}
}

type adminService struct {
	*Service
	userRepo  repository.UserRepository
	usageRepo repository.UsageRepository
}


