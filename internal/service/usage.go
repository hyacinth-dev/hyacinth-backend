package service

import (
	"context"
	v1 "hyacinth-backend/api/v1"
	// "hyacinth-backend/internal/model"
	"hyacinth-backend/internal/repository"
)

type UsageService interface {
	GetAllUsagesByUserId(ctx context.Context, req *v1.GetUsageRequest) (*v1.GetUsageResponseData, error)
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

func (s *usageService) GetAllUsagesByUserId(ctx context.Context, req *v1.GetUsageRequest) (*v1.GetUsageResponseData, error) {
	usages, err := s.usageRepository.GetAllUsagesByUserId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	usageData := make([]v1.UsageData, 0)
	for _, usage := range usages {
		usageData = append(usageData, v1.UsageData{
			Date:  usage.Date,
			Usage: usage.Usage,
		})
	}

	return &v1.GetUsageResponseData{
		Usages: usageData,
	}, nil
}
