package service

import (
	"context"
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/model"
	"hyacinth-backend/internal/repository"
	"sync"
)

type VnetService interface {
	GetVnetByUserId(ctx context.Context, id string) (*[]model.Vnet, error)
	GetVnetByVnetId(ctx context.Context, id string) (*model.Vnet, error)
	UpdateVnet(ctx context.Context, req *v1.UpdateVnetRequest) error
	CreateVnet(ctx context.Context, req *v1.CreateVnetRequest, userId string) error
	DeleteVnet(ctx context.Context, req *v1.DeleteVnetRequest) error
	EnableVnet(ctx context.Context, req *v1.EnableVnetRequest) error
	DisableVnet(ctx context.Context, req *v1.DisableVnetRequest) error
	CheckVnetTokenExists(ctx context.Context, token string, excludeVnetId string) (bool, error)
	GetOnlineTunnels(ctx context.Context, userId string) (int, error)
	GetOnlineDevicesCount(ctx context.Context, userId string) (int, error)
	GetRunningVnetCount(ctx context.Context, userId string) (int, error)
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
	vnetLock       sync.Mutex
	vnetRepository repository.VnetRepository
}

func (s *vnetService) GetVnetByUserId(ctx context.Context, id string) (*[]model.Vnet, error) {
	return s.vnetRepository.GetVnetByUserId(ctx, id)
}

func (s *vnetService) GetVnetByVnetId(ctx context.Context, id string) (*model.Vnet, error) {
	return s.vnetRepository.GetVnetByVnetId(ctx, id)
}

func (s *vnetService) UpdateVnet(ctx context.Context, req *v1.UpdateVnetRequest) error {
	vnet, err := s.vnetRepository.GetVnetByVnetId(ctx, req.VnetId)
	if err != nil {
		return err
	}
	vnet.Comment = req.Comment
	vnet.Enabled = req.Enabled
	vnet.Token = req.Token
	vnet.Password = req.Password
	vnet.IpRange = req.IpRange
	vnet.EnableDHCP = req.EnableDHCP
	vnet.ClientsLimit = req.ClientsLimit
	vnet.NeedUpdate = true
	if err := s.vnetRepository.UpdateVnet(ctx, vnet); err != nil {
		return err
	}
	return nil
}

func (s *vnetService) CreateVnet(ctx context.Context, req *v1.CreateVnetRequest, userId string) error {
	s.vnetLock.Lock()
	defer s.vnetLock.Unlock()
	vnet := &model.Vnet{
		VnetId:       req.VnetId,
		UserId:       userId,
		Comment:      req.Comment,
		Enabled:      req.Enabled,
		Token:        req.Token,
		Password:     req.Password,
		IpRange:      req.IpRange,
		EnableDHCP:   req.EnableDHCP,
		ClientsLimit: req.ClientsLimit,
		NeedUpdate:   true,
	}
	if err := s.vnetRepository.CreateVnet(ctx, vnet); err != nil {
		return err
	}
	return nil
}

func (s *vnetService) DeleteVnet(ctx context.Context, req *v1.DeleteVnetRequest) error {
	s.vnetLock.Lock()
	defer s.vnetLock.Unlock()
	vnet, err := s.vnetRepository.GetVnetByVnetId(ctx, req.VnetID)
	if err != nil {
		return err
	}
	if err := s.vnetRepository.DeleteVnet(ctx, vnet.VnetId); err != nil {
		return err
	}
	return nil
}

func (s *vnetService) EnableVnet(ctx context.Context, req *v1.EnableVnetRequest) error {
	vnet, err := s.vnetRepository.GetVnetByVnetId(ctx, req.VnetID)
	if err != nil {
		return err
	}
	vnet.Enabled = true
	vnet.NeedUpdate = true
	if err := s.vnetRepository.UpdateVnet(ctx, vnet); err != nil {
		return err
	}
	return nil
}

func (s *vnetService) DisableVnet(ctx context.Context, req *v1.DisableVnetRequest) error {
	vnet, err := s.vnetRepository.GetVnetByVnetId(ctx, req.VnetID)
	if err != nil {
		return err
	}
	vnet.Enabled = false
	vnet.NeedUpdate = true
	if err := s.vnetRepository.UpdateVnet(ctx, vnet); err != nil {
		return err
	}
	return nil
}

func (s *vnetService) CheckVnetTokenExists(ctx context.Context, token string, excludeVnetId string) (bool, error) {
	return s.vnetRepository.CheckVnetTokenExists(ctx, token, excludeVnetId)
}

func (s *vnetService) GetOnlineTunnels(ctx context.Context, userId string) (int, error) {
	return s.vnetRepository.GetOnlineTunnels(ctx, userId)
}

func (s *vnetService) GetOnlineDevicesCount(ctx context.Context, userId string) (int, error) {
	return s.vnetRepository.GetOnlineDevicesCount(ctx, userId)
}

func (s *vnetService) GetRunningVnetCount(ctx context.Context, userId string) (int, error) {
	return s.vnetRepository.GetRunningVnetCount(ctx, userId)
}
