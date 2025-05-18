package service

import (
	"context"
	"fmt"
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/model"
	"hyacinth-backend/internal/repository"
	"log"
	"strconv"
	"time"
)

type AdminService interface {
	GetTotalUsage(ctx context.Context) (*v1.GetTotalUsageResponseData, error)
	GetUsagePage(ctx context.Context, page string, pageSize int) (*v1.GetUsagePageResponseData, error)
	GetUsage(ctx context.Context, userId string, req *v1.GetUsageRequest) (*v1.GetUsageResponseData, error)
	AdminGetVNet(ctx context.Context, req *v1.AdminGetVNetRequest) (*v1.AdminGetVNetResponseData, error)
	AdminUpdateVNet(ctx context.Context, vnetID string, req *v1.AdminUpdateVNetRequest) error
	AdminCreateVNet(ctx context.Context, userId string, req *v1.AdminCreateVNetRequest) error
	AdminDeleteVNet(ctx context.Context, vnetID string) error
}

func NewAdminService(
	service *Service,
	userRepo repository.UserRepository,
	usageRepo repository.UsageRepository,
	venetRepo repository.VNetRepository,
) AdminService {
	return &adminService{
		userRepo:  userRepo,
		usageRepo: usageRepo,
		vnetRepo:  venetRepo,
		Service:   service,
	}
}

type adminService struct {
	userRepo  repository.UserRepository
	usageRepo repository.UsageRepository
	vnetRepo  repository.VNetRepository
	*Service
}

func (s *adminService) GetTotalUsage(ctx context.Context) (*v1.GetTotalUsageResponseData, error) {
	total, err := s.usageRepo.GetTotalUsage(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetTotalUsageResponseData{Total: total}, nil
}

func (s *adminService) GetUsagePage(ctx context.Context, page string, pageSize int) (*v1.GetUsagePageResponseData, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, fmt.Errorf("invalid page value: %v", err)
	}
	// 获取总页数
	totalCount, err := s.userRepo.CountAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	offset := (pageInt - 1) * pageSize

	// 获取该页内所有用户
	users, err := s.userRepo.GetUsersWithPagination(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	var results []v1.UsagePageItem
	for _, user := range users {
		// 统计该用户的usage总量
		totalUsage, err := s.usageRepo.GetByUser(ctx, user.UserId)
		if err != nil {
			return nil, err
		}
		// 统计该用户vnet数量
		vnets, err := s.vnetRepo.GetAllByUser(ctx, user.UserId)
		if err != nil {
			return nil, err
		}

		results = append(results, v1.UsagePageItem{
			UserID:      user.UserId,
			UserName:    user.Username,
			NumNetworks: len(vnets),
			Usage:       totalUsage,
		})
	}

	return &v1.GetUsagePageResponseData{Items: results, PageCount: (totalCount + pageSize - 1) / pageSize}, nil
}

func (s *adminService) GetUsage(ctx context.Context, userId string, req *v1.GetUsageRequest) (*v1.GetUsageResponseData, error) {

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

func (s *adminService) AdminGetVNet(ctx context.Context, req *v1.AdminGetVNetRequest) (*v1.AdminGetVNetResponseData, error) {
	var (
		vnets []model.VNet
		err   error
	)
	userId := req.UserID
	if userId == "0" {
		vnets, err = s.vnetRepo.GetAll(ctx)
	} else {
		vnets, err = s.vnetRepo.GetAllByUser(ctx, userId)
	}
	if err != nil {
		return nil, err
	}

	resp := &v1.AdminGetVNetResponseData{}

	userMap := make(map[string]string)

	// 查所有用户
	if userId == "0" {
		userIds := make([]string, 0, len(vnets))
		seen := make(map[string]struct{})
		for _, v := range vnets {
			if _, ok := seen[v.UserId]; !ok {
				seen[v.UserId] = struct{}{}
				userIds = append(userIds, v.UserId)
			}
		}

		// 2. 批量查询所有用户
		users, err := s.userRepo.GetByIDs(ctx, userIds)
		if err != nil {
			log.Printf("AdminGetVNet: 批量获取用户失败: %v\n", err)
		} else {
			for _, user := range users {
				userMap[user.UserId] = user.Username
			}
		}
	} else {
		// 单用户查询
		user, err := s.userRepo.GetByID(ctx, userId)
		if err != nil || user == nil {
			userMap[userId] = "N/A"
			log.Printf("AdminGetVNet：未找到用户userid:%s，错误: %v\n", userId, err)
		} else {
			userMap[userId] = user.Username
		}
	}

	for _, v := range vnets {
		username := userMap[v.UserId]
		if username == "" {
			username = "N/A"
		}

		resp.Networks = append(resp.Networks, v1.AdminVNetData{
			ID:           v.ID,
			Name:         v.Name,
			Enabled:      v.Enabled,
			Token:        v.Token,
			Password:     v.Password,
			IpRange:      v.IpRange,
			EnableDHCP:   v.EnableDHCP,
			ClientsLimit: v.ClientsLimit,
			Clients:      v.Clients,
			UserID:       v.UserId,
			UserName:     username,
		})
	}
	return resp, nil
}

func (s *adminService) AdminUpdateVNet(ctx context.Context, vnetID string, req *v1.AdminUpdateVNetRequest) error {
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

func (s *adminService) AdminCreateVNet(ctx context.Context, userId string, req *v1.AdminCreateVNetRequest) error {
	vnet := &model.VNet{
		Name:         req.Name,
		Enabled:      req.Enabled,
		Token:        req.Token,
		Password:     req.Password,
		IpRange:      req.IpRange,
		EnableDHCP:   req.EnableDHCP,
		ClientsLimit: req.ClientsLimit,
		UserId:       userId,
	}

	if err := s.vnetRepo.Create(ctx, vnet); err != nil {
		return err
	}

	return nil
}

func (s *adminService) AdminDeleteVNet(ctx context.Context, vnetID string) error {
	if err := s.vnetRepo.Delete(ctx, vnetID); err != nil {
		return err
	}

	return nil
}
