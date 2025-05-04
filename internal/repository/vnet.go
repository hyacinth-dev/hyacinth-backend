package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/model"
)

type VNetRepository interface {
	GetAllByUser(ctx context.Context, userId string) ([]model.VNet, error)
	GetAll(ctx context.Context) ([]model.VNet, error)
	GetByID(ctx context.Context, vnetID string) (*model.VNet, error)
	Update(ctx context.Context, vnet *model.VNet) error
	Create(ctx context.Context, vnet *model.VNet) error
	Delete(ctx context.Context, vnetID string) error
}

type vnetRepository struct {
	*Repository
}

func NewVNetRepository(r *Repository) VNetRepository {
	return &vnetRepository{Repository: r}
}

func (r *vnetRepository) GetAllByUser(ctx context.Context, userId string) ([]model.VNet, error) {
	var vnets []model.VNet
	if err := r.DB(ctx).Where("user_id = ?", userId).Find(&vnets).Error; err != nil {
		return nil, err
	}
	return vnets, nil
}

func (r *vnetRepository) GetAll(ctx context.Context) ([]model.VNet, error) {
	var vnets []model.VNet
	if err := r.DB(ctx).Find(&vnets).Error; err != nil {
		return nil, err
	}
	return vnets, nil
}

func (r *vnetRepository) GetByID(ctx context.Context, vnetID string) (*model.VNet, error) {
	var vnet model.VNet
	if err := r.DB(ctx).Where("id = ?", vnetID).First(&vnet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &vnet, nil
}

func (r *vnetRepository) Update(ctx context.Context, vnet *model.VNet) error {
	if err := r.DB(ctx).Save(vnet).Error; err != nil {
		return err
	}
	return nil
}

func (r *vnetRepository) Create(ctx context.Context, vnet *model.VNet) error {
	if err := r.DB(ctx).Create(vnet).Error; err != nil {
		return err
	}
	return nil
}

func (r *vnetRepository) Delete(ctx context.Context, vnetID string) error {
	if err := r.DB(ctx).Where("id = ?", vnetID).Delete(&model.VNet{}).Error; err != nil {
		return err
	}
	return nil
}
