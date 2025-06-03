package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"hyacinth-backend/internal/model"
	"hyacinth-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupVnetRepository(t *testing.T) (repository.VnetRepository, sqlmock.Sqlmock) {
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
	vnetRepo := repository.NewVnetRepository(repo)

	return vnetRepo, mock
}

func TestVnetRepository_CreateVnet(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	now := time.Now()
	vnet := &model.Vnet{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: now, Valid: false},
		},
		VnetId:        "vnet_123456",
		UserId:        "user_123456",
		Comment:       "Test VNet",
		Enabled:       true,
		Token:         "test_token",
		Password:      "test_password",
		IpRange:       "10.0.0.0/24",
		EnableDHCP:    true,
		ClientsLimit:  10,
		ClientsOnline: 0,
		NeedUpdate:    false,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `vnets` (`created_at`,`updated_at`,`deleted_at`,`vnet_id`,`user_id`,`comment`,`enabled`,`token`,`password`,`ip_range`,`enable_dhcp`,`clients_limit`,`clients_online`,`need_update`,`id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(vnet.CreatedAt, vnet.UpdatedAt, vnet.DeletedAt, vnet.VnetId, vnet.UserId, vnet.Comment, vnet.Enabled, vnet.Token, vnet.Password, vnet.IpRange, vnet.EnableDHCP, vnet.ClientsLimit, vnet.ClientsOnline, vnet.NeedUpdate, vnet.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := vnetRepo.CreateVnet(ctx, vnet)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_UpdateVnet(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	now := time.Now()
	vnet := &model.Vnet{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		VnetId:        "vnet_123456",
		UserId:        "user_123456",
		Comment:       "Updated Test VNet",
		Enabled:       false,
		Token:         "updated_token",
		Password:      "updated_password",
		IpRange:       "10.0.1.0/24",
		EnableDHCP:    false,
		ClientsLimit:  20,
		ClientsOnline: 5,
		NeedUpdate:    true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `vnets` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`vnet_id`=?,`user_id`=?,`comment`=?,`enabled`=?,`token`=?,`password`=?,`ip_range`=?,`enable_dhcp`=?,`clients_limit`=?,`clients_online`=?,`need_update`=? WHERE `vnets`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), vnet.DeletedAt, vnet.VnetId, vnet.UserId, vnet.Comment, vnet.Enabled, vnet.Token, vnet.Password, vnet.IpRange, vnet.EnableDHCP, vnet.ClientsLimit, vnet.ClientsOnline, vnet.NeedUpdate, vnet.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := vnetRepo.UpdateVnet(ctx, vnet)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_GetVnetByVnetId(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	vnetId := "vnet_123456"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "vnet_id", "user_id", "comment", "enabled", "token", "password", "ip_range", "enable_dhcp", "clients_limit", "clients_online", "need_update"}).
		AddRow(1, now, now, nil, "vnet_123456", "user_123456", "Test VNet", true, "test_token", "test_password", "10.0.0.0/24", true, 10, 3, false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `vnets` WHERE vnet_id = ? AND `vnets`.`deleted_at` IS NULL ORDER BY `vnets`.`id` LIMIT ?")).
		WithArgs(vnetId, 1).
		WillReturnRows(rows)

	vnet, err := vnetRepo.GetVnetByVnetId(ctx, vnetId)
	assert.NoError(t, err)
	assert.NotNil(t, vnet)
	assert.Equal(t, "vnet_123456", vnet.VnetId)
	assert.Equal(t, "user_123456", vnet.UserId)
	assert.Equal(t, "Test VNet", vnet.Comment)
	assert.True(t, vnet.Enabled)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_GetVnetByUserId(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "vnet_id", "user_id", "comment", "enabled", "token", "password", "ip_range", "enable_dhcp", "clients_limit", "clients_online", "need_update"}).
		AddRow(1, now, now, nil, "vnet_123456", "user_123456", "Test VNet 1", true, "test_token1", "test_password1", "10.0.0.0/24", true, 10, 3, false).
		AddRow(2, now, now, nil, "vnet_789012", "user_123456", "Test VNet 2", false, "test_token2", "test_password2", "10.0.1.0/24", false, 20, 0, true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `vnets` WHERE user_id = ? AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(userId).
		WillReturnRows(rows)

	vnets, err := vnetRepo.GetVnetByUserId(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, vnets)
	assert.Len(t, *vnets, 2)
	assert.Equal(t, "vnet_123456", (*vnets)[0].VnetId)
	assert.Equal(t, "vnet_789012", (*vnets)[1].VnetId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_DeleteVnet(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	vnetId := "vnet_123456"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `vnets` SET `deleted_at`=? WHERE vnet_id = ? AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), vnetId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := vnetRepo.DeleteVnet(ctx, vnetId)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_CheckVnetTokenExists(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	token := "test_token"
	excludeVnetId := "vnet_exclude"

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `vnets` WHERE (token = ? AND deleted_at IS NULL) AND vnet_id != ? AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(token, excludeVnetId).
		WillReturnRows(rows)

	exists, err := vnetRepo.CheckVnetTokenExists(ctx, token, excludeVnetId)
	assert.NoError(t, err)
	assert.True(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_CheckVnetTokenExists_NotFound(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	token := "nonexistent_token"
	excludeVnetId := ""

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `vnets` WHERE (token = ? AND deleted_at IS NULL) AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(token).
		WillReturnRows(rows)

	exists, err := vnetRepo.CheckVnetTokenExists(ctx, token, excludeVnetId)
	assert.NoError(t, err)
	assert.False(t, exists)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_GetOnlineTunnels(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	userId := "user_123456"

	rows := sqlmock.NewRows([]string{"count"}).AddRow(3)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `vnets` WHERE (enabled = ? AND deleted_at IS NULL) AND user_id = ? AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(true, userId).
		WillReturnRows(rows)

	count, err := vnetRepo.GetOnlineTunnels(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_GetOnlineTunnels_AllUsers(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	userId := "0" // Special case for all users

	rows := sqlmock.NewRows([]string{"count"}).AddRow(10)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `vnets` WHERE (enabled = ? AND deleted_at IS NULL) AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(true).
		WillReturnRows(rows)

	count, err := vnetRepo.GetOnlineTunnels(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, 10, count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_GetOnlineDevicesCount(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	userId := "user_123456"

	rows := sqlmock.NewRows([]string{"COALESCE(SUM(clients_online), 0)"}).AddRow(15)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(clients_online), 0) FROM `vnets` WHERE (enabled = ? AND deleted_at IS NULL) AND user_id = ? AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(true, userId).
		WillReturnRows(rows)

	count, err := vnetRepo.GetOnlineDevicesCount(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, 15, count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_GetOnlineDevicesCount_AllUsers(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	userId := "" // Empty string for all users

	rows := sqlmock.NewRows([]string{"COALESCE(SUM(clients_online), 0)"}).AddRow(50)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(clients_online), 0) FROM `vnets` WHERE (enabled = ? AND deleted_at IS NULL) AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(true).
		WillReturnRows(rows)

	count, err := vnetRepo.GetOnlineDevicesCount(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, 50, count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVnetRepository_GetRunningVnetCount(t *testing.T) {
	vnetRepo, mock := setupVnetRepository(t)

	ctx := context.Background()
	userId := "user_123456"

	rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `vnets` WHERE (user_id = ? AND enabled = ? AND deleted_at IS NULL) AND `vnets`.`deleted_at` IS NULL")).
		WithArgs(userId, true).
		WillReturnRows(rows)

	count, err := vnetRepo.GetRunningVnetCount(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	assert.NoError(t, mock.ExpectationsWereMet())
}
