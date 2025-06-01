package repository

import (
	"context"
	"hyacinth-backend/internal/model"
)

type VnetRepository interface {
	GetVnetByVnetId(ctx context.Context, vnetId string) (*model.Vnet, error)
	GetVnetByUserId(ctx context.Context, userId string) (*[]model.Vnet, error)
	UpdateVnet(ctx context.Context, vnet *model.Vnet) error
	CreateVnet(ctx context.Context, vnet *model.Vnet) error
	DeleteVnet(ctx context.Context, vnetId string) error
	CheckVnetTokenExists(ctx context.Context, token string, excludeVnetId string) (bool, error)
	GetOnlineTunnels(ctx context.Context, userId string) (int, error)
	GetOnlineDevicesCount(ctx context.Context, userId string) (int, error)
	GetRunningVnetCount(ctx context.Context, userId string) (int, error)
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

func (r *vnetRepository) GetVnetByVnetId(ctx context.Context, vnetId string) (*model.Vnet, error) {
	var vnet model.Vnet

	r.DB(ctx).Where("vnet_id = ?", vnetId).First(&vnet)

	return &vnet, nil
}

func (r *vnetRepository) GetVnetByUserId(ctx context.Context, userId string) (*[]model.Vnet, error) {
	var vnets []model.Vnet
	err := r.DB(ctx).Where("user_id = ?", userId).Find(&vnets).Error
	if err != nil {
		return nil, err
	}
	return &vnets, nil
}

func (r *vnetRepository) UpdateVnet(ctx context.Context, vnet *model.Vnet) error {
	return r.DB(ctx).Save(vnet).Error
}

func (r *vnetRepository) CreateVnet(ctx context.Context, vnet *model.Vnet) error {
	return r.DB(ctx).Create(vnet).Error
}

func (r *vnetRepository) DeleteVnet(ctx context.Context, vnetId string) error {
	return r.DB(ctx).Where("vnet_id = ?", vnetId).Delete(&model.Vnet{}).Error
}

func (r *vnetRepository) CheckVnetTokenExists(ctx context.Context, token string, excludeVnetId string) (bool, error) {
	var count int64
	query := r.DB(ctx).Model(&model.Vnet{}).Where("token = ? AND deleted_at IS NULL", token)
	if excludeVnetId != "" {
		query = query.Where("vnet_id != ?", excludeVnetId)
	}
	err := query.Count(&count).Error
	// print("count=", count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *vnetRepository) GetOnlineTunnels(ctx context.Context, userId string) (int, error) {
	var count int64
	query := r.DB(ctx).Model(&model.Vnet{}).Where("enabled = ? AND deleted_at IS NULL", true)
	if userId != "0" && userId != "" {
		query = query.Where("user_id = ?", userId)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *vnetRepository) GetOnlineDevicesCount(ctx context.Context, userId string) (int, error) {
	var totalOnlineDevices int64
	query := r.DB(ctx).Model(&model.Vnet{}).Where("enabled = ? AND deleted_at IS NULL", true)
	if userId != "0" && userId != "" {
		query = query.Where("user_id = ?", userId)
	}

	// 使用SUM函数计算所有虚拟网络的clients_online总和
	err := query.Select("COALESCE(SUM(clients_online), 0)").Scan(&totalOnlineDevices).Error
	if err != nil {
		return 0, err
	}
	return int(totalOnlineDevices), nil
}

func (r *vnetRepository) GetRunningVnetCount(ctx context.Context, userId string) (int, error) {
	var count int64
	err := r.DB(ctx).Model(&model.Vnet{}).Where("user_id = ? AND enabled = ? AND deleted_at IS NULL", userId, true).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
