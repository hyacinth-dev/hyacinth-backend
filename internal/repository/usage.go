package repository

import (
	"context"
	"hyacinth-backend/internal/model"
)

type UsageRepository interface {
	GetAllUsagesByUserId(ctx context.Context, userId int64) ([]*model.Usage, error)
}

func NewUsageRepository(
	repository *Repository,
) UsageRepository {
	return &usageRepository{
		Repository: repository,
	}
}

type usageRepository struct {
	*Repository
}

func (r *usageRepository) GetAllUsagesByUserId(ctx context.Context, userId int64) ([]*model.Usage, error) {
	var usages []*model.Usage
	if err := r.DB(ctx).Where("user_id = ?", userId).Find(&usages).Error; err != nil {
		return nil, err
	}
	return usages, nil
}
