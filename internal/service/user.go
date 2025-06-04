// 在此模块定义业务逻辑的实现
// 不关心请求的来源，接受结构化了的请求
// 按照业务规则进行操作并返回结果
// 在此模块的所有方法都不需要鉴权

package service

import (
	"context"
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/model"
	"hyacinth-backend/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
	"sync"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userId string, req *v1.ChangePasswordRequest) error
	PurchasePackage(ctx context.Context, userId string, req *v1.PurchasePackageRequest) error
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
	CheckUsernameExists(ctx context.Context, username string, excludeUserId string) (bool, error)
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
) UserService {
	return &userService{
		userRepo:      userRepo,
		Service:       service,
		registerMutex: sync.Mutex{},
	}
}

type userService struct {
	userRepo      repository.UserRepository
	registerMutex sync.Mutex // 确保注册操作的线程安全
	*Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	s.registerMutex.Lock()
	defer s.registerMutex.Unlock()
	// check username
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrEmailAlreadyUse
	}

	user, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrUsernameAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate user ID
	userId, err := s.sid.GenString()
	if err != nil {
		return err
	}
	user = &model.User{
		UserId:           userId,
		Username:         req.Username,
		Email:            req.Email,
		Password:         string(hashedPassword),
		UserGroup:        1,                              // 默认为普通用户
		RemainingTraffic: model.DefaultTrafficForNewUser, // 为新用户提供默认初始流量
	}
	// Transaction demo
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.UsernameOrEmail)
	if err != nil || user == nil {
		user, err = s.userRepo.GetByUsername(ctx, req.UsernameOrEmail)
		if err != nil || user == nil {
			return nil, v1.ErrUnauthorized
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}
	token, err := s.jwt.GenToken(user.UserId, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return nil, err
	}

	return &v1.LoginResponseData{
		AccessToken: token,
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	activeTunnels := 0

	availableTraffic := user.FormatRemainingTraffic()

	onlineDevices := 0

	var privilegeExpiryStr *string
	if user.PrivilegeExpiry != nil {
		expiry := user.PrivilegeExpiry.Format("2006-01-02 15:04:05")
		privilegeExpiryStr = &expiry
	}

	return &v1.GetProfileResponseData{
		UserId:           user.UserId,
		Username:         user.Username,
		Email:            user.Email,
		UserGroup:        user.UserGroup,
		UserGroupName:    user.GetUserGroupName(),
		PrivilegeExpiry:  privilegeExpiryStr,
		IsVip:            user.IsVip(),
		ActiveTunnels:    activeTunnels,
		AvailableTraffic: availableTraffic,
		OnlineDevices:    onlineDevices,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Username = req.Username

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, userId string, req *v1.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	// 验证当前密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return v1.ErrOriginalPasswordNotMatch
	}

	// 生成新密码的哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) PurchasePackage(ctx context.Context, userId string, req *v1.PurchasePackageRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	// 如果没有指定购买时长
	duration := req.Duration
	if duration <= 0 {
		return v1.ErrBadRequest
	}

	now := time.Now()

	// 定义不同套餐的流量额度（字节）
	var monthlyTraffic int64
	switch req.PackageType {
	case 2: // 青铜用户 - 50GB/月
		monthlyTraffic = model.BronzeMonthlyTraffic
	case 3: // 白银用户 - 200GB/月
		monthlyTraffic = model.SilverMonthlyTraffic
	case 4: // 黄金用户 - 无限流量（设置为1TB）
		monthlyTraffic = model.GoldMonthlyTraffic
	default:
		return v1.ErrBadRequest
	}

	// 计算总流量（购买时长 * 月流量）
	totalTraffic := monthlyTraffic

	if req.PackageType == user.UserGroup {
		// 相同用户组，特权时间增加指定月数，流量不变
		if user.PrivilegeExpiry == nil || user.PrivilegeExpiry.Before(now) {
			// 如果没有特权或已过期，从当前时间开始计算
			expiry := now.AddDate(0, duration, 0)
			user.PrivilegeExpiry = &expiry
			user.RemainingTraffic = totalTraffic
		} else {
			// 如果特权未过期，在现有基础上增加指定月数
			expiry := user.PrivilegeExpiry.AddDate(0, duration, 0)
			user.PrivilegeExpiry = &expiry
			// user.RemainingTraffic += totalTraffic
		}
	} else {
		// 购买不同的用户组（升级或降级），重置为指定时长，设置新的流量额度
		user.UserGroup = req.PackageType
		expiry := now.AddDate(0, duration, 0)
		user.PrivilegeExpiry = &expiry
		user.RemainingTraffic = totalTraffic
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userId)
}

func (s *userService) CheckUsernameExists(ctx context.Context, username string, excludeUserId string) (bool, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	// 如果找到的用户就是要排除的用户，则认为用户名不冲突
	if user.UserId == excludeUserId {
		return false, nil
	}
	return true, nil
}
