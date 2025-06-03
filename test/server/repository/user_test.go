package repository

import (
	"context"
	"hyacinth-backend/pkg/log"
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

var logger *log.Logger

func setupRepository(t *testing.T) (repository.UserRepository, sqlmock.Sqlmock) {
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
	userRepo := repository.NewUserRepository(repo)

	return userRepo, mock
}

func TestUserRepository_Create(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	now := time.Now()
	user := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: gorm.DeletedAt{Time: now, Valid: false},
		},
		UserId:           "user_123456",
		Username:         "testuser",
		Password:         "hashedpassword",
		Email:            "test@example.com",
		UserGroup:        1,
		PrivilegeExpiry:  nil,
		RemainingTraffic: model.DefaultTrafficForNewUser,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`username`,`password`,`email`,`user_group`,`remaining_traffic`,`id`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.UserId, user.Username, user.Password, user.Email, user.UserGroup, user.RemainingTraffic, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	now := time.Now()
	user := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserId:           "user_123456",
		Username:         "updateduser",
		Password:         "hashedpassword",
		Email:            "updated@example.com",
		UserGroup:        2,
		RemainingTraffic: model.DefaultTrafficForNewUser,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`user_id`=?,`username`=?,`password`=?,`email`=?,`user_group`=?,`privilege_expiry`=?,`remaining_traffic`=? WHERE `users`.`deleted_at` IS NULL AND `id` = ?")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), user.DeletedAt, user.UserId, user.Username, user.Password, user.Email, user.UserGroup, user.PrivilegeExpiry, user.RemainingTraffic, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Update(ctx, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	userId := "user_123456"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "username", "password", "email", "user_group", "remaining_traffic"}).
		AddRow(1, now, now, nil, "user_123456", "testuser", "hashedpassword", "test@example.com", 1, model.DefaultTrafficForNewUser)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE user_id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(userId, 1).
		WillReturnRows(rows)

	user, err := userRepo.GetByID(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user_123456", user.UserId)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	userId := "nonexistent"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE user_id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(userId, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := userRepo.GetByID(ctx, userId)
	assert.Error(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	email := "test@example.com"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "username", "password", "email", "user_group", "remaining_traffic"}).
		AddRow(1, now, now, nil, "user_123456", "testuser", "hashedpassword", "test@example.com", 1, model.DefaultTrafficForNewUser)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(email, 1).
		WillReturnRows(rows)

	user, err := userRepo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "testuser", user.Username)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	email := "nonexistent@example.com"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(email, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := userRepo.GetByEmail(ctx, email)
	assert.NoError(t, err) // GetByEmail returns nil error when not found
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByUsername(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	username := "testuser"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "username", "password", "email", "user_group", "remaining_traffic"}).
		AddRow(1, now, now, nil, "user_123456", "testuser", "hashedpassword", "test@example.com", 1, model.DefaultTrafficForNewUser)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ? AND deleted_at IS NULL) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(username, 1).
		WillReturnRows(rows)

	user, err := userRepo.GetByUsername(ctx, username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByUsername_NotFound(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	username := "nonexistent"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ? AND deleted_at IS NULL) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(username, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := userRepo.GetByUsername(ctx, username)
	assert.NoError(t, err) // GetByUsername returns nil error when not found
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}
