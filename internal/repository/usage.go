package repository

import (
	"context"
	"hyacinth-backend/internal/model"
)

type UsageRepository interface {
	GetUsageByUserAndRange(ctx context.Context, userId string, startDate string, endDate string) ([]model.Usages, error)
	GetTotalUsage(ctx context.Context) (int, error)
	GetByUser(ctx context.Context, userId string) (int, error)
}

func NewUsageRepository(
	r *Repository,
) UsageRepository {
	return &usageRepository{
		Repository: r,
	}
}

type usageRepository struct {
	*Repository
}

func (r *usageRepository) GetUsageByUserAndRange(ctx context.Context, userId string, startDate string, endDate string) ([]model.Usages, error) {
	var usages []model.Usages
	err := r.DB(ctx).
		Where("user_id = ? AND date >= ? AND date <= ?", userId, startDate, endDate).
		Find(&usages).Error

	if err != nil {
		return nil, err
	}

	return usages, nil
}

func (r *usageRepository) GetTotalUsage(ctx context.Context) (int, error) {
	var total int64
	err := r.DB(ctx).Model(&model.Usages{}).Select("SUM(usage)").Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return int(total), nil
}

func (r *usageRepository) GetByUser(ctx context.Context, userId string) (int, error) {
	var total int64
	err := r.DB(ctx).Model(&model.Usages{}).
		Where("user_id = ?", userId).
		Select("SUM(usage)").Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return int(total), nil
}
