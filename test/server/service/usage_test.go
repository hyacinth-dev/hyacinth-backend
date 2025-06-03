package service_test

import (
	"context"
	"errors"
	"testing"

	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/service"
	mock_repository "hyacinth-backend/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupUsageService(t *testing.T) (service.UsageService, *mock_repository.MockUsageRepository) {
	ctrl := gomock.NewController(t)

	mockUsageRepo := mock_repository.NewMockUsageRepository(ctrl)
	srv := service.NewService(nil, logger, sf, j)
	usageService := service.NewUsageService(srv, mockUsageRepo)

	return usageService, mockUsageRepo
}

func TestUsageService_GetUsage_24h(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "vnet_123",
		Range:  "24h",
	}

	expectedUsages := &[]v1.UsageData{
		{Date: "01-15 08:00", Usage: 1024},
		{Date: "01-15 09:00", Usage: 2048},
		{Date: "01-15 10:00", Usage: 1536},
		{Date: "01-15 11:00", Usage: 3072},
		{Date: "01-15 12:00", Usage: 2560},
	}

	// Mock期望：获取24小时使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(expectedUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 5)
	assert.Equal(t, "01-15 08:00", result.Usages[0].Date)
	assert.Equal(t, 1024, result.Usages[0].Usage)
	assert.Equal(t, "01-15 12:00", result.Usages[4].Date)
	assert.Equal(t, 2560, result.Usages[4].Usage)
}

func TestUsageService_GetUsage_7d(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "",
		Range:  "7d",
	}

	expectedUsages := &[]v1.UsageData{
		{Date: "2024-01-15", Usage: 10240},
		{Date: "2024-01-16", Usage: 15360},
		{Date: "2024-01-17", Usage: 8192},
		{Date: "2024-01-18", Usage: 20480},
		{Date: "2024-01-19", Usage: 12288},
		{Date: "2024-01-20", Usage: 18432},
		{Date: "2024-01-21", Usage: 14336},
	}

	// Mock期望：获取7天使用量数据（所有虚拟网络）
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(expectedUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 7)
	assert.Equal(t, "2024-01-15", result.Usages[0].Date)
	assert.Equal(t, 10240, result.Usages[0].Usage)
	assert.Equal(t, "2024-01-21", result.Usages[6].Date)
	assert.Equal(t, 14336, result.Usages[6].Usage)
}

func TestUsageService_GetUsage_30d(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "vnet_456",
		Range:  "30d",
	}

	expectedUsages := &[]v1.UsageData{
		{Date: "2024-01-01", Usage: 102400},
		{Date: "2024-01-02", Usage: 153600},
		{Date: "2024-01-03", Usage: 81920},
		{Date: "2024-01-04", Usage: 204800},
		{Date: "2024-01-05", Usage: 122880},
	}

	// Mock期望：获取30天使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(expectedUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 5)
	assert.Equal(t, "2024-01-01", result.Usages[0].Date)
	assert.Equal(t, 102400, result.Usages[0].Usage)
}

func TestUsageService_GetUsage_Month(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "",
		Range:  "month",
	}

	expectedUsages := &[]v1.UsageData{
		{Date: "2023-02", Usage: 1048576},
		{Date: "2023-03", Usage: 2097152},
		{Date: "2023-04", Usage: 1572864},
		{Date: "2023-05", Usage: 3145728},
		{Date: "2023-06", Usage: 2621440},
		{Date: "2023-07", Usage: 4194304},
		{Date: "2023-08", Usage: 3670016},
		{Date: "2023-09", Usage: 5242880},
		{Date: "2023-10", Usage: 4718592},
		{Date: "2023-11", Usage: 6291456},
		{Date: "2023-12", Usage: 5767168},
		{Date: "2024-01", Usage: 7340032},
	}

	// Mock期望：获取按月统计的使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(expectedUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 12)
	assert.Equal(t, "2023-02", result.Usages[0].Date)
	assert.Equal(t, 1048576, result.Usages[0].Usage)
	assert.Equal(t, "2024-01", result.Usages[11].Date)
	assert.Equal(t, 7340032, result.Usages[11].Usage)
}

func TestUsageService_GetUsage_All(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "vnet_789",
		Range:  "all",
	}

	expectedUsages := &[]v1.UsageData{
		{Date: "2021", Usage: 10485760},
		{Date: "2022", Usage: 20971520},
		{Date: "2023", Usage: 31457280},
		{Date: "2024", Usage: 41943040},
	}

	// Mock期望：获取按年统计的所有使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(expectedUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 4)
	assert.Equal(t, "2021", result.Usages[0].Date)
	assert.Equal(t, 10485760, result.Usages[0].Usage)
	assert.Equal(t, "2024", result.Usages[3].Date)
	assert.Equal(t, 41943040, result.Usages[3].Usage)
}

func TestUsageService_GetUsage_EmptyResult(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "new_user_123",
		VnetId: "new_vnet_123",
		Range:  "30d",
	}

	emptyUsages := &[]v1.UsageData{}

	// Mock期望：新用户或新虚拟网络没有使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(emptyUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 0)
}

func TestUsageService_GetUsage_AllVnets(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "", // 空值表示所有虚拟网络
		Range:  "7d",
	}

	expectedUsages := &[]v1.UsageData{
		{Date: "2024-01-15", Usage: 30720}, // 来自多个虚拟网络的汇总数据
		{Date: "2024-01-16", Usage: 45056},
		{Date: "2024-01-17", Usage: 28672},
		{Date: "2024-01-18", Usage: 61440},
		{Date: "2024-01-19", Usage: 40960},
		{Date: "2024-01-20", Usage: 53248},
		{Date: "2024-01-21", Usage: 47104},
	}

	// Mock期望：获取用户所有虚拟网络的使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(expectedUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 7)
	assert.Equal(t, "2024-01-15", result.Usages[0].Date)
	assert.Equal(t, 30720, result.Usages[0].Usage)
}

func TestUsageService_GetUsage_RepositoryError(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "vnet_123",
		Range:  "7d",
	}

	// Mock期望：数据库查询时发生错误
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(nil, errors.New("database connection error"))

	result, err := usageService.GetUsage(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database connection error", err.Error())
}

func TestUsageService_GetUsage_InvalidRange(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "vnet_123",
		Range:  "invalid_range",
	}

	// Mock期望：无效的时间范围参数，repository可能返回错误
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(nil, errors.New("invalid time range"))

	result, err := usageService.GetUsage(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid time range", err.Error())
}

func TestUsageService_GetUsage_UserNotExists(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "nonexistent_user",
		VnetId: "vnet_123",
		Range:  "30d",
	}

	emptyUsages := &[]v1.UsageData{}

	// Mock期望：用户不存在，返回空的使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(emptyUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 0)
}

func TestUsageService_GetUsage_VnetNotExists(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "user_123",
		VnetId: "nonexistent_vnet",
		Range:  "7d",
	}

	emptyUsages := &[]v1.UsageData{}

	// Mock期望：虚拟网络不存在，返回空的使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(emptyUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 0)
}

func TestUsageService_GetUsage_LargeDataset(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "heavy_user_123",
		VnetId: "",
		Range:  "month",
	}

	// 生成大量数据来测试性能
	expectedUsages := &[]v1.UsageData{}
	for i := 1; i <= 100; i++ {
		*expectedUsages = append(*expectedUsages, v1.UsageData{
			Date:  "2024-01",
			Usage: i * 1024 * 1024, // 每个数据点1MB递增
		})
	}

	// Mock期望：获取大量使用量数据
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(expectedUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 100)
	assert.Equal(t, "2024-01", result.Usages[0].Date)
	assert.Equal(t, 1024*1024, result.Usages[0].Usage)
	assert.Equal(t, 100*1024*1024, result.Usages[99].Usage)
}

func TestUsageService_GetUsage_ZeroUsage(t *testing.T) {
	usageService, mockUsageRepo := setupUsageService(t)

	ctx := context.Background()
	req := &v1.GetUsageRequest{
		UserId: "inactive_user_123",
		VnetId: "vnet_123",
		Range:  "7d",
	}

	expectedUsages := &[]v1.UsageData{
		{Date: "2024-01-15", Usage: 0},
		{Date: "2024-01-16", Usage: 0},
		{Date: "2024-01-17", Usage: 0},
		{Date: "2024-01-18", Usage: 0},
		{Date: "2024-01-19", Usage: 0},
		{Date: "2024-01-20", Usage: 0},
		{Date: "2024-01-21", Usage: 0},
	}

	// Mock期望：获取零使用量的数据（非活跃用户）
	mockUsageRepo.EXPECT().GetUsage(ctx, req.UserId, req.VnetId, req.Range).Return(expectedUsages, nil)

	result, err := usageService.GetUsage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Usages, 7)
	for i, usage := range result.Usages {
		assert.Equal(t, 0, usage.Usage, "Usage at index %d should be 0", i)
		assert.Contains(t, usage.Date, "2024-01-")
	}
}
