package service_test

import (
	"context"
	"errors"
	"testing"

	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/model"
	"hyacinth-backend/internal/service"
	mock_repository "hyacinth-backend/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupVnetService(t *testing.T) (service.VnetService, *mock_repository.MockVnetRepository, *mock_repository.MockTransaction) {
	ctrl := gomock.NewController(t)

	mockVnetRepo := mock_repository.NewMockVnetRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	vnetService := service.NewVnetService(srv, mockVnetRepo)

	return vnetService, mockVnetRepo, mockTm
}

func TestVnetService_GetVnetByUserId(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"

	expectedVnets := &[]model.Vnet{
		{
			VnetId:        "vnet_1",
			UserId:        userId,
			Comment:       "测试虚拟网络1",
			Enabled:       true,
			Token:         "token_1",
			Password:      "password_1",
			IpRange:       "192.168.1.0/24",
			EnableDHCP:    true,
			ClientsLimit:  10,
			ClientsOnline: 5,
			NeedUpdate:    false,
		},
		{
			VnetId:        "vnet_2",
			UserId:        userId,
			Comment:       "测试虚拟网络2",
			Enabled:       false,
			Token:         "token_2",
			Password:      "password_2",
			IpRange:       "192.168.2.0/24",
			EnableDHCP:    false,
			ClientsLimit:  5,
			ClientsOnline: 0,
			NeedUpdate:    true,
		},
	}

	// Mock期望：通过用户ID获取虚拟网络列表
	mockVnetRepo.EXPECT().GetVnetByUserId(ctx, userId).Return(expectedVnets, nil)

	result, err := vnetService.GetVnetByUserId(ctx, userId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)
	assert.Equal(t, "vnet_1", (*result)[0].VnetId)
	assert.Equal(t, "vnet_2", (*result)[1].VnetId)
}

func TestVnetService_GetVnetByUserId_Empty(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"

	emptyVnets := &[]model.Vnet{}

	// Mock期望：用户没有虚拟网络
	mockVnetRepo.EXPECT().GetVnetByUserId(ctx, userId).Return(emptyVnets, nil)

	result, err := vnetService.GetVnetByUserId(ctx, userId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 0)
}

func TestVnetService_GetVnetByVnetId(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	vnetId := "vnet_123"

	expectedVnet := &model.Vnet{
		VnetId:        vnetId,
		UserId:        "user_123",
		Comment:       "测试虚拟网络",
		Enabled:       true,
		Token:         "test_token",
		Password:      "test_password",
		IpRange:       "192.168.1.0/24",
		EnableDHCP:    true,
		ClientsLimit:  10,
		ClientsOnline: 3,
		NeedUpdate:    false,
	}

	// Mock期望：通过VnetId获取虚拟网络
	mockVnetRepo.EXPECT().GetVnetByVnetId(ctx, vnetId).Return(expectedVnet, nil)

	result, err := vnetService.GetVnetByVnetId(ctx, vnetId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, vnetId, result.VnetId)
	assert.Equal(t, "user_123", result.UserId)
	assert.Equal(t, "测试虚拟网络", result.Comment)
}

func TestVnetService_GetVnetByVnetId_NotFound(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	vnetId := "nonexistent_vnet"

	// Mock期望：虚拟网络不存在
	mockVnetRepo.EXPECT().GetVnetByVnetId(ctx, vnetId).Return(nil, gorm.ErrRecordNotFound)

	result, err := vnetService.GetVnetByVnetId(ctx, vnetId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestVnetService_CreateVnet(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.CreateVnetRequest{
		VnetProfile: v1.VnetProfile{
			VnetId:       "new_vnet_123",
			Comment:      "新虚拟网络",
			Enabled:      true,
			Token:        "new_token",
			Password:     "new_password",
			IpRange:      "192.168.3.0/24",
			EnableDHCP:   true,
			ClientsLimit: 15,
		},
	}

	// Mock期望：创建虚拟网络成功
	mockVnetRepo.EXPECT().CreateVnet(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, vnet *model.Vnet) error {
		assert.Equal(t, req.VnetId, vnet.VnetId)
		assert.Equal(t, userId, vnet.UserId)
		assert.Equal(t, req.Comment, vnet.Comment)
		assert.Equal(t, req.Enabled, vnet.Enabled)
		assert.Equal(t, req.Token, vnet.Token)
		assert.Equal(t, req.Password, vnet.Password)
		assert.Equal(t, req.IpRange, vnet.IpRange)
		assert.Equal(t, req.EnableDHCP, vnet.EnableDHCP)
		assert.Equal(t, req.ClientsLimit, vnet.ClientsLimit)
		assert.True(t, vnet.NeedUpdate)
		return nil
	})

	err := vnetService.CreateVnet(ctx, req, userId)

	assert.NoError(t, err)
}

func TestVnetService_CreateVnet_Error(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.CreateVnetRequest{
		VnetProfile: v1.VnetProfile{
			VnetId:       "duplicate_vnet",
			Comment:      "重复的虚拟网络",
			Enabled:      true,
			Token:        "duplicate_token",
			Password:     "password",
			IpRange:      "192.168.4.0/24",
			EnableDHCP:   false,
			ClientsLimit: 8,
		},
	}

	// Mock期望：创建虚拟网络失败（例如重复的VnetId）
	mockVnetRepo.EXPECT().CreateVnet(ctx, gomock.Any()).Return(errors.New("duplicate key error"))

	err := vnetService.CreateVnet(ctx, req, userId)

	assert.Error(t, err)
	assert.Equal(t, "duplicate key error", err.Error())
}

func TestVnetService_UpdateVnet(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	req := &v1.UpdateVnetRequest{
		VnetProfile: v1.VnetProfile{
			VnetId:       "vnet_123",
			Comment:      "更新的虚拟网络",
			Enabled:      false,
			Token:        "updated_token",
			Password:     "updated_password",
			IpRange:      "192.168.5.0/24",
			EnableDHCP:   false,
			ClientsLimit: 20,
		},
	}

	existingVnet := &model.Vnet{
		VnetId:        req.VnetId,
		UserId:        "user_123",
		Comment:       "旧的虚拟网络",
		Enabled:       true,
		Token:         "old_token",
		Password:      "old_password",
		IpRange:       "192.168.1.0/24",
		EnableDHCP:    true,
		ClientsLimit:  10,
		ClientsOnline: 5,
		NeedUpdate:    false,
	}

	// Mock期望：先获取现有虚拟网络
	mockVnetRepo.EXPECT().GetVnetByVnetId(ctx, req.VnetId).Return(existingVnet, nil)

	// Mock期望：更新虚拟网络
	mockVnetRepo.EXPECT().UpdateVnet(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, vnet *model.Vnet) error {
		assert.Equal(t, req.Comment, vnet.Comment)
		assert.Equal(t, req.Enabled, vnet.Enabled)
		assert.Equal(t, req.Token, vnet.Token)
		assert.Equal(t, req.Password, vnet.Password)
		assert.Equal(t, req.IpRange, vnet.IpRange)
		assert.Equal(t, req.EnableDHCP, vnet.EnableDHCP)
		assert.Equal(t, req.ClientsLimit, vnet.ClientsLimit)
		assert.True(t, vnet.NeedUpdate)
		return nil
	})

	err := vnetService.UpdateVnet(ctx, req)

	assert.NoError(t, err)
}

func TestVnetService_UpdateVnet_VnetNotFound(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	req := &v1.UpdateVnetRequest{
		VnetProfile: v1.VnetProfile{
			VnetId:       "nonexistent_vnet",
			Comment:      "不存在的虚拟网络",
			Enabled:      true,
			Token:        "token",
			Password:     "password",
			IpRange:      "192.168.6.0/24",
			EnableDHCP:   true,
			ClientsLimit: 5,
		},
	}

	// Mock期望：虚拟网络不存在
	mockVnetRepo.EXPECT().GetVnetByVnetId(ctx, req.VnetId).Return(nil, gorm.ErrRecordNotFound)

	err := vnetService.UpdateVnet(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestVnetService_DeleteVnet(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	req := &v1.DeleteVnetRequest{
		VnetID: "vnet_to_delete",
	}

	existingVnet := &model.Vnet{
		VnetId:   req.VnetID,
		UserId:   "user_123",
		Comment:  "要删除的虚拟网络",
		Enabled:  true,
		Token:    "token",
		Password: "password",
		IpRange:  "192.168.7.0/24",
	}

	// Mock期望：先获取要删除的虚拟网络
	mockVnetRepo.EXPECT().GetVnetByVnetId(ctx, req.VnetID).Return(existingVnet, nil)

	// Mock期望：删除虚拟网络
	mockVnetRepo.EXPECT().DeleteVnet(ctx, req.VnetID).Return(nil)

	err := vnetService.DeleteVnet(ctx, req)

	assert.NoError(t, err)
}

func TestVnetService_DeleteVnet_VnetNotFound(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	req := &v1.DeleteVnetRequest{
		VnetID: "nonexistent_vnet",
	}

	// Mock期望：虚拟网络不存在
	mockVnetRepo.EXPECT().GetVnetByVnetId(ctx, req.VnetID).Return(nil, gorm.ErrRecordNotFound)

	err := vnetService.DeleteVnet(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestVnetService_EnableVnet(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	req := &v1.EnableVnetRequest{
		VnetID: "vnet_to_enable",
	}

	existingVnet := &model.Vnet{
		VnetId:     req.VnetID,
		UserId:     "user_123",
		Comment:    "要启用的虚拟网络",
		Enabled:    false, // 当前是禁用状态
		Token:      "token",
		Password:   "password",
		IpRange:    "192.168.8.0/24",
		NeedUpdate: false,
	}

	// Mock期望：先获取虚拟网络
	mockVnetRepo.EXPECT().GetVnetByVnetId(ctx, req.VnetID).Return(existingVnet, nil)

	// Mock期望：更新虚拟网络状态
	mockVnetRepo.EXPECT().UpdateVnet(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, vnet *model.Vnet) error {
		assert.True(t, vnet.Enabled)
		assert.True(t, vnet.NeedUpdate)
		return nil
	})

	err := vnetService.EnableVnet(ctx, req)

	assert.NoError(t, err)
}

func TestVnetService_DisableVnet(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	req := &v1.DisableVnetRequest{
		VnetID: "vnet_to_disable",
	}

	existingVnet := &model.Vnet{
		VnetId:     req.VnetID,
		UserId:     "user_123",
		Comment:    "要禁用的虚拟网络",
		Enabled:    true, // 当前是启用状态
		Token:      "token",
		Password:   "password",
		IpRange:    "192.168.9.0/24",
		NeedUpdate: false,
	}

	// Mock期望：先获取虚拟网络
	mockVnetRepo.EXPECT().GetVnetByVnetId(ctx, req.VnetID).Return(existingVnet, nil)

	// Mock期望：更新虚拟网络状态
	mockVnetRepo.EXPECT().UpdateVnet(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, vnet *model.Vnet) error {
		assert.False(t, vnet.Enabled)
		assert.True(t, vnet.NeedUpdate)
		return nil
	})

	err := vnetService.DisableVnet(ctx, req)

	assert.NoError(t, err)
}

func TestVnetService_CheckVnetTokenExists(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	token := "existing_token"
	excludeVnetId := "vnet_123"

	// Mock期望：Token已存在
	mockVnetRepo.EXPECT().CheckVnetTokenExists(ctx, token, excludeVnetId).Return(true, nil)

	exists, err := vnetService.CheckVnetTokenExists(ctx, token, excludeVnetId)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestVnetService_CheckVnetTokenExists_NotExists(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	token := "unique_token"
	excludeVnetId := "vnet_123"

	// Mock期望：Token不存在
	mockVnetRepo.EXPECT().CheckVnetTokenExists(ctx, token, excludeVnetId).Return(false, nil)

	exists, err := vnetService.CheckVnetTokenExists(ctx, token, excludeVnetId)

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestVnetService_GetOnlineTunnels(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"
	expectedCount := 3

	// Mock期望：获取在线隧道数量
	mockVnetRepo.EXPECT().GetOnlineTunnels(ctx, userId).Return(expectedCount, nil)

	count, err := vnetService.GetOnlineTunnels(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestVnetService_GetOnlineDevicesCount(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"
	expectedCount := 8

	// Mock期望：获取在线设备数量
	mockVnetRepo.EXPECT().GetOnlineDevicesCount(ctx, userId).Return(expectedCount, nil)

	count, err := vnetService.GetOnlineDevicesCount(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestVnetService_GetRunningVnetCount(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"
	expectedCount := 2

	// Mock期望：获取运行中的虚拟网络数量
	mockVnetRepo.EXPECT().GetRunningVnetCount(ctx, userId).Return(expectedCount, nil)

	count, err := vnetService.GetRunningVnetCount(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestVnetService_GetOnlineTunnels_Error(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"

	// Mock期望：获取在线隧道数量时发生错误
	mockVnetRepo.EXPECT().GetOnlineTunnels(ctx, userId).Return(0, errors.New("database error"))

	count, err := vnetService.GetOnlineTunnels(ctx, userId)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Equal(t, "database error", err.Error())
}

func TestVnetService_GetOnlineDevicesCount_Error(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"

	// Mock期望：获取在线设备数量时发生错误
	mockVnetRepo.EXPECT().GetOnlineDevicesCount(ctx, userId).Return(0, errors.New("connection error"))

	count, err := vnetService.GetOnlineDevicesCount(ctx, userId)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Equal(t, "connection error", err.Error())
}

func TestVnetService_GetRunningVnetCount_Error(t *testing.T) {
	vnetService, mockVnetRepo, _ := setupVnetService(t)

	ctx := context.Background()
	userId := "user_123"

	// Mock期望：获取运行中的虚拟网络数量时发生错误
	mockVnetRepo.EXPECT().GetRunningVnetCount(ctx, userId).Return(0, errors.New("query error"))

	count, err := vnetService.GetRunningVnetCount(ctx, userId)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Equal(t, "query error", err.Error())
}
