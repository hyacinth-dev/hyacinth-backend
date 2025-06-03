package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/model"
	"hyacinth-backend/internal/service"
	mock_repository "hyacinth-backend/test/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupUserService(t *testing.T) (service.UserService, *mock_repository.MockUserRepository, *mock_repository.MockTransaction) {
	ctrl := gomock.NewController(t)

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	return userService, mockUserRepo, mockTm
}

func TestUserService_Register(t *testing.T) {
	userService, mockUserRepo, mockTm := setupUserService(t)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	// Mock期望：检查邮箱是否存在（返回不存在）
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)

	// Mock期望：检查用户名是否存在（返回不存在）
	mockUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, nil)

	// Mock期望：事务执行成功，在事务中创建用户
	mockTm.EXPECT().Transaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		// 模拟事务内的Create调用
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
		return fn(ctx)
	})

	err := userService.Register(ctx, req)

	assert.NoError(t, err)
}

func TestUserService_Register_EmailExists(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	// Mock期望：邮箱已存在
	existingUser := &model.User{
		UserId: "existing_user",
		Email:  req.Email,
	}
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(existingUser, nil)

	err := userService.Register(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, v1.ErrEmailAlreadyUse, err)
}

func TestUserService_Register_UsernameExists(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	// Mock期望：邮箱不存在
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)

	// Mock期望：用户名已存在
	existingUser := &model.User{
		UserId:   "existing_user",
		Username: req.Username,
	}
	mockUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(existingUser, nil)

	err := userService.Register(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, v1.ErrUsernameAlreadyUse, err)
}

func TestUserService_Register_GetByEmailError(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	// Mock期望：查询邮箱时发生错误
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, errors.New("database error"))

	err := userService.Register(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, v1.ErrInternalServerError, err)
}

func TestUserService_Register_GetByUsernameError(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	// Mock期望：邮箱不存在
	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)

	// Mock期望：查询用户名时发生错误
	mockUserRepo.EXPECT().GetByUsername(ctx, req.Username).Return(nil, errors.New("database error"))

	err := userService.Register(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, v1.ErrInternalServerError, err)
}

func TestUserService_Login(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	req := &v1.LoginRequest{
		UsernameOrEmail: "test@example.com",
		Password:        "password123",
	}

	// 生成测试用的哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := &model.User{
		UserId:   "user_123",
		Username: "testuser",
		Email:    req.UsernameOrEmail,
		Password: string(hashedPassword),
	}

	// Mock期望：通过邮箱找到用户
	mockUserRepo.EXPECT().GetByEmail(ctx, req.UsernameOrEmail).Return(user, nil)

	result, err := userService.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	req := &v1.LoginRequest{
		UsernameOrEmail: "nonexistent@example.com",
		Password:        "password123",
	}

	// Mock期望：通过邮箱找不到用户
	mockUserRepo.EXPECT().GetByEmail(ctx, req.UsernameOrEmail).Return(nil, gorm.ErrRecordNotFound)

	// Mock期望：通过用户名也找不到用户
	mockUserRepo.EXPECT().GetByUsername(ctx, req.UsernameOrEmail).Return(nil, gorm.ErrRecordNotFound)

	result, err := userService.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, v1.ErrUnauthorized, err)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	req := &v1.LoginRequest{
		UsernameOrEmail: "test@example.com",
		Password:        "wrongpassword",
	}

	// 生成不同密码的哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := &model.User{
		UserId:   "user_123",
		Username: "testuser",
		Email:    req.UsernameOrEmail,
		Password: string(hashedPassword),
	}

	// Mock期望：找到用户但密码错误
	mockUserRepo.EXPECT().GetByEmail(ctx, req.UsernameOrEmail).Return(user, nil)

	result, err := userService.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserService_GetProfile(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"

	user := &model.User{
		UserId:           userId,
		Username:         "testuser",
		Email:            "test@example.com",
		UserGroup:        1,
		RemainingTraffic: 1024 * 1024 * 1024, // 1GB
	}

	// Mock期望：通过ID找到用户
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(user, nil)

	result, err := userService.GetProfile(ctx, userId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userId, result.UserId)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestUserService_GetProfile_UserNotFound(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "nonexistent_user"

	// Mock期望：用户不存在
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(nil, gorm.ErrRecordNotFound)

	result, err := userService.GetProfile(ctx, userId)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserService_UpdateProfile(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.UpdateProfileRequest{
		Username: "newusername",
		Email:    "newemail@example.com",
	}

	existingUser := &model.User{
		UserId:   userId,
		Username: "oldusername",
		Email:    "oldemail@example.com",
	}

	// Mock期望：先获取现有用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	// Mock期望：更新用户信息
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err := userService.UpdateProfile(ctx, userId, req)

	assert.NoError(t, err)
}

func TestUserService_UpdateProfile_UserNotFound(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "nonexistent_user"
	req := &v1.UpdateProfileRequest{
		Username: "newusername",
		Email:    "newemail@example.com",
	}

	// Mock期望：用户不存在
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(nil, gorm.ErrRecordNotFound)

	err := userService.UpdateProfile(ctx, userId, req)

	assert.Error(t, err)
}

func TestUserService_GetUserByID(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"

	expectedUser := &model.User{
		UserId:   userId,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Mock期望：通过ID找到用户
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(expectedUser, nil)

	result, err := userService.GetUserByID(ctx, userId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userId, result.UserId)
	assert.Equal(t, "testuser", result.Username)
}

func TestUserService_CheckUsernameExists(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	username := "testuser"
	excludeUserId := "user_123"

	existingUser := &model.User{
		UserId:   "other_user",
		Username: username,
	}

	// Mock期望：找到使用该用户名的其他用户
	mockUserRepo.EXPECT().GetByUsername(ctx, username).Return(existingUser, nil)

	exists, err := userService.CheckUsernameExists(ctx, username, excludeUserId)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserService_CheckUsernameExists_NotExists(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	username := "newusername"
	excludeUserId := "user_123"

	// Mock期望：没有找到使用该用户名的用户
	mockUserRepo.EXPECT().GetByUsername(ctx, username).Return(nil, nil)

	exists, err := userService.CheckUsernameExists(ctx, username, excludeUserId)

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestUserService_CheckUsernameExists_SameUser(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	username := "testuser"
	excludeUserId := "user_123"

	existingUser := &model.User{
		UserId:   excludeUserId, // 同一个用户
		Username: username,
	}

	// Mock期望：找到的是同一个用户
	mockUserRepo.EXPECT().GetByUsername(ctx, username).Return(existingUser, nil)

	exists, err := userService.CheckUsernameExists(ctx, username, excludeUserId)

	assert.NoError(t, err)
	assert.False(t, exists) // 同一个用户不算存在冲突
}

func TestUserService_ChangePassword(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	// 生成旧密码的哈希
	hashedOldPassword, err := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	assert.NoError(t, err)

	req := &v1.ChangePasswordRequest{
		CurrentPassword: oldPassword,
		NewPassword:     newPassword,
	}

	existingUser := &model.User{
		UserId:   userId,
		Password: string(hashedOldPassword),
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	// Mock期望：更新密码
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err = userService.ChangePassword(ctx, userId, req)

	assert.NoError(t, err)
}

func TestUserService_ChangePassword_WrongOldPassword(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	wrongOldPassword := "wrongpassword"
	newPassword := "newpassword"

	// 生成不同的旧密码哈希
	hashedOldPassword, err := bcrypt.GenerateFromPassword([]byte("correctoldpassword"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	req := &v1.ChangePasswordRequest{
		CurrentPassword: wrongOldPassword,
		NewPassword:     newPassword,
	}

	existingUser := &model.User{
		UserId:   userId,
		Password: string(hashedOldPassword),
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	err = userService.ChangePassword(ctx, userId, req)

	assert.Error(t, err)
}

func TestUserService_PurchasePackage_Success_Upgrade(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.PurchasePackageRequest{
		PackageType: 2, // 青铜套餐
		Duration:    3, // 3个月
	}

	existingUser := &model.User{
		UserId:    userId,
		UserGroup: 1, // 普通用户升级到青铜
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	// Mock期望：更新用户信息
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *model.User) error {
		// 验证用户组已更新
		assert.Equal(t, 2, user.UserGroup)
		// 验证特权到期时间已设置（应该是3个月后）
		assert.NotNil(t, user.PrivilegeExpiry)
		// 验证剩余流量已设置为青铜套餐流量
		assert.Equal(t, model.BronzeMonthlyTraffic, user.RemainingTraffic)
		return nil
	})

	err := userService.PurchasePackage(ctx, userId, req)

	assert.NoError(t, err)
}

func TestUserService_PurchasePackage_Success_SameLevel_ExpiredPrivilege(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.PurchasePackageRequest{
		PackageType: 3, // 白银套餐
		Duration:    1, // 1个月
	}

	// 特权已过期的时间
	expiredTime := time.Now().AddDate(0, 0, -1) // 昨天过期
	existingUser := &model.User{
		UserId:          userId,
		UserGroup:       3, // 白银用户续费
		PrivilegeExpiry: &expiredTime,
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	// Mock期望：更新用户信息
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *model.User) error {
		// 验证用户组未变
		assert.Equal(t, 3, user.UserGroup)
		// 验证特权到期时间重新设置为1个月后
		assert.NotNil(t, user.PrivilegeExpiry)
		assert.True(t, user.PrivilegeExpiry.After(time.Now()))
		// 验证剩余流量重新设置为白银套餐流量
		assert.Equal(t, model.SilverMonthlyTraffic, user.RemainingTraffic)
		return nil
	})

	err := userService.PurchasePackage(ctx, userId, req)

	assert.NoError(t, err)
}

func TestUserService_PurchasePackage_Success_SameLevel_ValidPrivilege(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.PurchasePackageRequest{
		PackageType: 4, // 黄金套餐
		Duration:    2, // 2个月
	}

	// 特权还有效的时间（还有10天到期）
	validTime := time.Now().AddDate(0, 0, 10)
	existingUser := &model.User{
		UserId:          userId,
		UserGroup:       4, // 黄金用户续费
		PrivilegeExpiry: &validTime,
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	// Mock期望：更新用户信息
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *model.User) error {
		// 验证用户组未变
		assert.Equal(t, 4, user.UserGroup)
		// 验证特权到期时间在原基础上增加2个月
		assert.NotNil(t, user.PrivilegeExpiry)
		expectedExpiry := validTime.AddDate(0, 2, 0)
		assert.True(t, user.PrivilegeExpiry.Equal(expectedExpiry) || user.PrivilegeExpiry.After(expectedExpiry.Add(-time.Minute)))
		return nil
	})

	err := userService.PurchasePackage(ctx, userId, req)

	assert.NoError(t, err)
}

func TestUserService_PurchasePackage_UserNotFound(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "nonexistent_user"
	req := &v1.PurchasePackageRequest{
		PackageType: 2,
		Duration:    1,
	}

	// Mock期望：用户不存在
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(nil, gorm.ErrRecordNotFound)

	err := userService.PurchasePackage(ctx, userId, req)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserService_PurchasePackage_InvalidPackageType(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.PurchasePackageRequest{
		PackageType: 5, // 无效的套餐类型
		Duration:    1,
	}

	existingUser := &model.User{
		UserId:    userId,
		UserGroup: 1,
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	err := userService.PurchasePackage(ctx, userId, req)

	assert.Error(t, err)
	assert.Equal(t, v1.ErrBadRequest, err)
}

func TestUserService_PurchasePackage_InvalidDuration_Zero(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.PurchasePackageRequest{
		PackageType: 2,
		Duration:    0, // 无效的购买时长
	}

	existingUser := &model.User{
		UserId:    userId,
		UserGroup: 1,
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	err := userService.PurchasePackage(ctx, userId, req)

	assert.Error(t, err)
	assert.Equal(t, v1.ErrBadRequest, err)
}

func TestUserService_PurchasePackage_InvalidDuration_Negative(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.PurchasePackageRequest{
		PackageType: 3,
		Duration:    -1, // 负数购买时长
	}

	existingUser := &model.User{
		UserId:    userId,
		UserGroup: 1,
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	err := userService.PurchasePackage(ctx, userId, req)

	assert.Error(t, err)
	assert.Equal(t, v1.ErrBadRequest, err)
}

func TestUserService_PurchasePackage_UpdateError(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	ctx := context.Background()
	userId := "user_123"
	req := &v1.PurchasePackageRequest{
		PackageType: 2,
		Duration:    1,
	}

	existingUser := &model.User{
		UserId:    userId,
		UserGroup: 1,
	}

	// Mock期望：获取用户信息
	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

	// Mock期望：更新用户信息失败
	updateError := errors.New("database update error")
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(updateError)

	err := userService.PurchasePackage(ctx, userId, req)

	assert.Error(t, err)
	assert.Equal(t, updateError, err)
}

func TestUserService_PurchasePackage_DifferentPackageTypes(t *testing.T) {
	userService, mockUserRepo, _ := setupUserService(t)

	testCases := []struct {
		name         string
		packageType  int
		expectedFlow int64
	}{
		{
			name:         "Bronze Package",
			packageType:  2,
			expectedFlow: model.BronzeMonthlyTraffic,
		},
		{
			name:         "Silver Package",
			packageType:  3,
			expectedFlow: model.SilverMonthlyTraffic,
		},
		{
			name:         "Gold Package",
			packageType:  4,
			expectedFlow: model.GoldMonthlyTraffic,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			userId := "user_123"
			req := &v1.PurchasePackageRequest{
				PackageType: tc.packageType,
				Duration:    1,
			}

			existingUser := &model.User{
				UserId:    userId,
				UserGroup: 1, // 普通用户升级
			}

			// Mock期望：获取用户信息
			mockUserRepo.EXPECT().GetByID(ctx, userId).Return(existingUser, nil)

			// Mock期望：更新用户信息
			mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *model.User) error {
				// 验证用户组已更新
				assert.Equal(t, tc.packageType, user.UserGroup)
				// 验证剩余流量已设置为对应套餐流量
				assert.Equal(t, tc.expectedFlow, user.RemainingTraffic)
				return nil
			})

			err := userService.PurchasePackage(ctx, userId, req)

			assert.NoError(t, err)
		})
	}
}
