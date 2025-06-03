package repository

import (
	"context"
	"regexp"
	"testing"

	"hyacinth-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupUsageRepository(t *testing.T) (repository.UsageRepository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	repo := repository.NewRepository(logger, db)
	usageRepo := repository.NewUsageRepository(repo)

	return usageRepo, mock
}

func TestUsageRepository_GetUsage_24h(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	vnetId := "vnet_123456"
	timeRange := "24h"

	rows := sqlmock.NewRows([]string{"date", "usage"}).
		AddRow("06-03 10:00", 1024).
		AddRow("06-03 11:00", 2048).
		AddRow("06-03 12:00", 4096)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE_FORMAT(created_at, '%m-%d %H:00') as `date`, SUM(`usage`) as `usage` FROM `usages` WHERE deleted_at IS NULL AND user_id = ? AND vnet_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 1 DAY) AND `usages`.`deleted_at` IS NULL GROUP BY DATE_FORMAT(created_at, '%m-%d %H:00') ORDER BY `date` ASC")).
		WithArgs(userId, vnetId).
		WillReturnRows(rows)

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 24)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsageRepository_GetUsage_7d(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	vnetId := ""
	timeRange := "7d"

	rows := sqlmock.NewRows([]string{"date", "usage"}).
		AddRow("2025-06-01", 10240).
		AddRow("2025-06-02", 20480).
		AddRow("2025-06-03", 40960)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE_FORMAT(created_at, '%Y-%m-%d') as `date`, SUM(`usage`) as `usage` FROM `usages` WHERE deleted_at IS NULL AND user_id = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) AND `usages`.`deleted_at` IS NULL GROUP BY DATE_FORMAT(created_at, '%Y-%m-%d') ORDER BY created_at ASC")).
		WithArgs(userId).
		WillReturnRows(rows)

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 7)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsageRepository_GetUsage_30d(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	vnetId := "vnet_123456"
	timeRange := "30d"

	rows := sqlmock.NewRows([]string{"date", "usage"}).
		AddRow("2025-05-04", 102400).
		AddRow("2025-05-05", 204800).
		AddRow("2025-05-06", 409600)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE_FORMAT(created_at, '%Y-%m-%d') as `date`, SUM(`usage`) as `usage` FROM `usages` WHERE deleted_at IS NULL AND user_id = ? AND vnet_id = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY) AND `usages`.`deleted_at` IS NULL GROUP BY DATE_FORMAT(created_at, '%Y-%m-%d') ORDER BY created_at ASC")).
		WithArgs(userId, vnetId).
		WillReturnRows(rows)

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 30)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsageRepository_GetUsage_Month(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	vnetId := ""
	timeRange := "month"

	rows := sqlmock.NewRows([]string{"date", "usage"}).
		AddRow("2024-01", 1024000).
		AddRow("2024-02", 2048000).
		AddRow("2024-03", 4096000)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE_FORMAT(created_at, '%Y-%m') as `date`, SUM(`usage`) as `usage` FROM `usages` WHERE deleted_at IS NULL AND user_id = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH) AND `usages`.`deleted_at` IS NULL GROUP BY DATE_FORMAT(created_at, '%Y-%m') ORDER BY created_at ASC")).
		WithArgs(userId).
		WillReturnRows(rows)

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 13)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsageRepository_GetUsage_WithUserIdOnly(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	vnetId := ""
	timeRange := "24h"

	rows := sqlmock.NewRows([]string{"date", "usage"}).
		AddRow("06-03 14:00", 8192).
		AddRow("06-03 15:00", 16384)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE_FORMAT(created_at, '%m-%d %H:00') as `date`, SUM(`usage`) as `usage` FROM `usages` WHERE deleted_at IS NULL AND user_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 1 DAY) AND `usages`.`deleted_at` IS NULL GROUP BY DATE_FORMAT(created_at, '%m-%d %H:00') ORDER BY `date` ASC")).
		WithArgs(userId).
		WillReturnRows(rows)

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 24) // Should fill missing hours with 0

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsageRepository_GetUsage_WithVnetIdOnly(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := ""
	vnetId := "vnet_123456"
	timeRange := "7d"

	rows := sqlmock.NewRows([]string{"date", "usage"}).
		AddRow("2024-06-01", 32768).
		AddRow("2024-06-03", 65536)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE_FORMAT(created_at, '%Y-%m-%d') as `date`, SUM(`usage`) as `usage` FROM `usages` WHERE deleted_at IS NULL AND vnet_id = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) AND `usages`.`deleted_at` IS NULL GROUP BY DATE_FORMAT(created_at, '%Y-%m-%d') ORDER BY created_at ASC")).
		WithArgs(vnetId).
		WillReturnRows(rows)

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 7) // Should fill missing days with 0

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsageRepository_GetUsage_InvalidTimeRange(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	vnetId := "vnet_123456"
	timeRange := "invalid"

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.NoError(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsageRepository_GetUsage_EmptyResult(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	vnetId := "vnet_123456"
	timeRange := "24h"

	rows := sqlmock.NewRows([]string{"date", "usage"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE_FORMAT(created_at, '%m-%d %H:00') as `date`, SUM(`usage`) as `usage` FROM `usages` WHERE deleted_at IS NULL AND user_id = ? AND vnet_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 1 DAY) AND `usages`.`deleted_at` IS NULL GROUP BY DATE_FORMAT(created_at, '%m-%d %H:00') ORDER BY `date` ASC")).
		WithArgs(userId, vnetId).
		WillReturnRows(rows)

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 24) // Should still return 24 hours with 0 usage

	// Check that all entries have 0 usage
	for _, usage := range *result {
		assert.Equal(t, 0, usage.Usage)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUsageRepository_GetUsage_DatabaseError(t *testing.T) {
	usageRepo, mock := setupUsageRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	vnetId := "vnet_123456"
	timeRange := "24h"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE_FORMAT(created_at, '%m-%d %H:00') as `date`, SUM(`usage`) as `usage` FROM `usages` WHERE deleted_at IS NULL AND user_id = ? AND vnet_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 1 DAY) AND `usages`.`deleted_at` IS NULL GROUP BY DATE_FORMAT(created_at, '%m-%d %H:00') ORDER BY `date` ASC")).
		WithArgs(userId, vnetId).
		WillReturnError(assert.AnError)

	result, err := usageRepo.GetUsage(ctx, userId, vnetId, timeRange)
	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}
