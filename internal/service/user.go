// 在此模块定义业务逻辑的实现
// 不关心请求的来源，接受结构化了的请求
// 按照业务规则进行操作并返回结果
// 在此模块的所有方法都不需要鉴权

package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/model"
	"hyacinth-backend/internal/repository"
	"time"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponseData, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
	GetUsage(ctx context.Context, userId string, req *v1.GetUsageRequest) (*v1.GetUsageResponseData, error)
	GetVNet(ctx context.Context, userId string) (*v1.GetVNetResponseData, error)
	UpdateVNet(ctx context.Context, vnetID string, req *v1.UpdateVNetRequest) error
	CreateVNet(ctx context.Context, userId string, req *v1.CreateVNetRequest) error
	DeleteVNet(ctx context.Context, vnetID string) error
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
	usageRepo repository.UsageRepository,
	venetRepo repository.VNetRepository,
) UserService {
	return &userService{
		userRepo:  userRepo,
		usageRepo: usageRepo,
		vnetRepo:  venetRepo,
		Service:   service,
	}
}

type userService struct {
	userRepo  repository.UserRepository
	usageRepo repository.UsageRepository
	vnetRepo  repository.VNetRepository
	*Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if user != nil {
		return v1.ErrEmailAlreadyUse
	}

	// check username
	user, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if user != nil {
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
		UserId:   userId,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		IsAdmin:  false, // 默认非管理员
		IsVip:    false, // 默认非VIP
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
	user, err := s.userRepo.GetByEmail(ctx, req.EmailOrUsername)
	if err == nil && user == nil {
		user, err = s.userRepo.GetByUsername(ctx, req.EmailOrUsername)
	}
	if err != nil {
		return nil, v1.ErrInternalServerError
	} else if user == nil {
		return nil, v1.ErrUnauthorized
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
		IsAdmin:     user.IsAdmin,
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:   user.UserId,
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
		IsVip:    user.IsVip,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Username = req.Username
	user.IsVip = req.IsVip
	user.IsAdmin = req.IsAdmin

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUsage(ctx context.Context, userId string, req *v1.GetUsageRequest) (*v1.GetUsageResponseData, error) {

	var startDate, endDate string
	switch req.Range {
	case "month":
		// 获取最近12个月的数据
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, -12, 0).Format("2006-01-02")
	case "30days":
		// 获取最近30天的数据
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	case "7days":
		// 获取最近7天的数据
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	default:
		return nil, v1.ErrBadRequest
	}

	usages, err := s.usageRepo.GetUsageByUserAndRange(ctx, userId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var responseData v1.GetUsageResponseData
	if req.Range == "month" {
		// 按月聚合流量数据
		usageByMonth := make(map[string]int)
		for _, usage := range usages {
			month := usage.Date[:7] // 取年月部分
			usageByMonth[month] += usage.UsageCount
		}

		for month, usage := range usageByMonth {
			responseData.Usages = append(responseData.Usages, v1.UsageData{
				DateOrMonth: month,
				Usage:       usage,
			})
		}
	} else {
		// 按天返回数据
		for _, usage := range usages {
			responseData.Usages = append(responseData.Usages, v1.UsageData{
				DateOrMonth: usage.Date[5:],
				Usage:       usage.UsageCount,
			})
		}
	}

	return &responseData, nil
}

func (s *userService) GetVNet(ctx context.Context, userId string) (*v1.GetVNetResponseData, error) {
	vnets, err := s.vnetRepo.GetAllByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	resp := &v1.GetVNetResponseData{}
	for _, v := range vnets {
		resp.Networks = append(resp.Networks, v1.VNetData{
			ID:           v.ID,
			Name:         v.Name,
			Enabled:      v.Enabled,
			Token:        v.Token,
			Password:     v.Password,
			IpRange:      v.IpRange,
			EnableDHCP:   v.EnableDHCP,
			ClientsLimit: v.ClientsLimit,
			Clients:      v.Clients,
		})
	}
	return resp, nil
}

func (s *userService) UpdateVNet(ctx context.Context, vnetID string, req *v1.UpdateVNetRequest) error {
	vnet, err := s.vnetRepo.GetByID(ctx, vnetID)
	if err != nil {
		return err
	}

	vnet.Name = req.Name
	vnet.Enabled = req.Enabled
	vnet.Token = req.Token
	vnet.Password = req.Password
	vnet.IpRange = req.IpRange
	vnet.EnableDHCP = req.EnableDHCP
	vnet.ClientsLimit = req.ClientsLimit

	if err = s.vnetRepo.Update(ctx, vnet); err != nil {
		return err
	}

	return nil
}

func (s *userService) CreateVNet(ctx context.Context, userId string, req *v1.CreateVNetRequest) error {
	vnet := &model.VNet{
		ID:           req.Name,
		Name:         req.Name,
		Enabled:      req.Enabled,
		Token:        req.Token,
		Password:     req.Password,
		IpRange:      req.IpRange,
		EnableDHCP:   req.EnableDHCP,
		ClientsLimit: req.ClientsLimit,
		UserId:       userId,
	}

	// 检测id是否重名
	existingVnet, err := s.vnetRepo.GetByID(ctx, vnet.ID)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if existingVnet != nil {
		return v1.ErrVNetAlreadyExists
	}
	// 恢复软删除的vnet，避免create重复id报错
	deletedVnet, err := s.vnetRepo.GetByIDWithDeleted(ctx, vnet.ID)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if deletedVnet != nil {
		if err = s.vnetRepo.Update(ctx, vnet); err != nil {
			return err
		}
	} else {
		if err = s.vnetRepo.Create(ctx, vnet); err != nil {
			return err
		}
	}

	return nil
}

func (s *userService) DeleteVNet(ctx context.Context, vnetID string) error {
	if err := s.vnetRepo.Delete(ctx, vnetID); err != nil {
		return err
	}

	return nil
}
