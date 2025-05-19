package repository

import (
    "context"
	"hyacinth-backend/internal/model"
)

type VnetRepository interface {
	GetVnet(ctx context.Context, id int64) (*model.Vnet, error)
}

func NewVnetRepository(
	repository *Repository,
) VnetRepository {
	return &vnetRepository{
		Repository: repository,
	}
}

type vnetRepository struct {
	*Repository
}

func (r *vnetRepository) GetVnet(ctx context.Context, id int64) (*model.Vnet, error) {
	var vnet model.Vnet

	return &vnet, nil
}
