package service

import (
	"context"
	v1 "hyacinth-backend/api/v1"

	// "hyacinth-backend/internal/model"
	"hyacinth-backend/internal/repository"
)

type UsageService interface {
	GetUsage(ctx context.Context, req *v1.GetUsageRequest) (*v1.GetUsageResponseData, error)
}

func NewUsageService(
	service *Service,
	usageRepository repository.UsageRepository,
) UsageService {
	return &usageService{
		Service:         service,
		usageRepository: usageRepository,
	}
}

type usageService struct {
	*Service
	usageRepository repository.UsageRepository
}

func (s *usageService) GetUsage(ctx context.Context, req *v1.GetUsageRequest) (*v1.GetUsageResponseData, error) {
	usages, err := s.usageRepository.GetUsage(ctx, req.UserId, req.VnetId, req.Range)
	if err != nil {
		return nil, err
	}
	return &v1.GetUsageResponseData{
		Usages: *usages,
	}, nil
}
