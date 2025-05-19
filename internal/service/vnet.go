package service

import (
    "context"
	"hyacinth-backend/internal/model"
	"hyacinth-backend/internal/repository"
)

type VnetService interface {
	GetVnet(ctx context.Context, id int64) (*model.Vnet, error)
}
func NewVnetService(
    service *Service,
    vnetRepository repository.VnetRepository,
) VnetService {
	return &vnetService{
		Service:        service,
		vnetRepository: vnetRepository,
	}
}

type vnetService struct {
	*Service
	vnetRepository repository.VnetRepository
}

func (s *vnetService) GetVnet(ctx context.Context, id int64) (*model.Vnet, error) {
	return s.vnetRepository.GetVnet(ctx, id)
}
