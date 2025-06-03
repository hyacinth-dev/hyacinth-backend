package handler

import (
	v1 "hyacinth-backend/api/v1"
	"hyacinth-backend/internal/handler"
	"hyacinth-backend/internal/middleware"
	"hyacinth-backend/internal/model"
	mock_service "hyacinth-backend/test/mocks/service"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

// createTestRouter 为每个测试创建独立的路由器实例
func createTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	testRouter := gin.New()
	testRouter.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
	)
	return testRouter
}

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RegisterRequest{
		Username: "testuser",
		Email:    "testuser@gmail.com",
		Password: "123456",
	}

	testRouter := createTestRouter()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 设置期望的方法调用
	mockUserService.EXPECT().Register(gomock.Any(), &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.POST("/register", userHandler.Register)

	obj := newHttpExcept(t, testRouter).POST("/register").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.LoginRequest{
		UsernameOrEmail: "testuser@gmail.com",
		Password:        "123456",
	}

	// 创建正确的响应数据结构
	loginResponse := &v1.LoginResponseData{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiJ0ZXN0dXNlciIsImV4cCI6MTczODIyMDUxNCwibmJmIjoxNzMwNDQ0NTE0LCJpYXQiOjE3MzA0NDQ1MTR9.test",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 设置期望的方法调用
	mockUserService.EXPECT().Login(gomock.Any(), &params).Return(loginResponse, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.POST("/login", userHandler.Login)

	obj := newHttpExcept(t, testRouter).POST("/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	obj.Value("data").Object().Value("accessToken").IsEqual(loginResponse.AccessToken)
}

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testuser"
	email := "testuser@gmail.com"
	activeTunnels := 2
	onlineDevices := 3

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 设置期望的方法调用
	profileData := &v1.GetProfileResponseData{
		UserId:        userId,
		Username:      username,
		Email:         email,
		UserGroup:     1,
		ActiveTunnels: 0, // 初始值，会在handler中被更新
		OnlineDevices: 0, // 初始值，会在handler中被更新
	}

	mockUserService.EXPECT().GetProfile(gomock.Any(), userId).Return(profileData, nil)
	mockVnetService.EXPECT().GetOnlineTunnels(gomock.Any(), userId).Return(activeTunnels, nil)
	mockVnetService.EXPECT().GetOnlineDevicesCount(gomock.Any(), userId).Return(onlineDevices, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.NoStrictAuth(jwt, logger))
	testRouter.GET("/user", userHandler.GetProfile)

	obj := newHttpExcept(t, testRouter).GET("/user").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("userId").IsEqual(userId)
	objData.Value("username").IsEqual(username)
	objData.Value("email").IsEqual(email)
	objData.Value("activeTunnels").IsEqual(activeTunnels)
	objData.Value("onlineDevices").IsEqual(onlineDevices)
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.UpdateProfileRequest{
		Username: "newusername",
		Email:    "newemail@gmail.com",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟当前用户信息
	currentUser := &model.User{
		UserId:   userId,
		Username: "oldusername",
		Email:    "oldemail@gmail.com",
	}

	// 设置期望的方法调用
	mockUserService.EXPECT().GetUserByID(gomock.Any(), userId).Return(currentUser, nil)
	mockUserService.EXPECT().CheckUsernameExists(gomock.Any(), params.Username, userId).Return(false, nil)
	mockUserService.EXPECT().UpdateProfile(gomock.Any(), userId, &params).Return(nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.PUT("/user", userHandler.UpdateProfile)

	obj := newHttpExcept(t, testRouter).PUT("/user").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_PurchasePackage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.PurchasePackageRequest{
		PackageType: 3, // 白银套餐
		Duration:    1, // 1个月
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟当前用户信息（普通用户升级到白银）
	currentUser := &model.User{
		UserId:    userId,
		Username:  "testuser",
		Email:     "test@gmail.com",
		UserGroup: 1, // 普通用户
	}

	// 设置期望的方法调用
	mockUserService.EXPECT().GetUserByID(gomock.Any(), userId).Return(currentUser, nil)
	mockUserService.EXPECT().PurchasePackage(gomock.Any(), userId, &params).Return(nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/user/purchase", userHandler.PurchasePackage)

	obj := newHttpExcept(t, testRouter).POST("/user/purchase").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_PurchasePackage_SameLevel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.PurchasePackageRequest{
		PackageType: 2, // 青铜套餐
		Duration:    3, // 3个月
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟当前用户信息（青铜用户续费）
	currentUser := &model.User{
		UserId:    userId,
		Username:  "testuser",
		Email:     "test@gmail.com",
		UserGroup: 2, // 青铜用户
	}

	// 设置期望的方法调用
	mockUserService.EXPECT().GetUserByID(gomock.Any(), userId).Return(currentUser, nil)
	mockUserService.EXPECT().PurchasePackage(gomock.Any(), userId, &params).Return(nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/user/purchase", userHandler.PurchasePackage)

	obj := newHttpExcept(t, testRouter).POST("/user/purchase").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_PurchasePackage_DowngradeWithClientsLimitExceeded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.PurchasePackageRequest{
		PackageType: 2, // 降级到青铜套餐
		Duration:    1, // 1个月
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟当前用户信息（白银用户降级到青铜）
	currentUser := &model.User{
		UserId:    userId,
		Username:  "testuser",
		Email:     "test@gmail.com",
		UserGroup: 3, // 白银用户
	}

	// 模拟虚拟网络数据，其中有一个网络的连接数超过青铜套餐限制
	vnets := &[]model.Vnet{
		{
			VnetId:       "vnet1",
			UserId:       userId,
			ClientsLimit: 8, // 超过青铜套餐的5人限制
			Enabled:      true,
		},
		{
			VnetId:       "vnet2",
			UserId:       userId,
			ClientsLimit: 3, // 未超过限制
			Enabled:      true,
		},
	}

	// 设置期望的方法调用
	mockUserService.EXPECT().GetUserByID(gomock.Any(), userId).Return(currentUser, nil)
	mockVnetService.EXPECT().GetVnetByUserId(gomock.Any(), userId).Return(vnets, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/user/purchase", userHandler.PurchasePackage)

	obj := newHttpExcept(t, testRouter).POST("/user/purchase").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object()
	// 验证返回的错误码不为0
	obj.Value("code").Number().Gt(0)
}

func TestUserHandler_PurchasePackage_InvalidPackageType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.PurchasePackageRequest{
		PackageType: 5, // 无效的套餐类型
		Duration:    1,
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/user/purchase", userHandler.PurchasePackage)

	obj := newHttpExcept(t, testRouter).POST("/user/purchase").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object()
	// 验证返回的错误码不为0
	obj.Value("code").Number().Gt(0)
}

func TestUserHandler_PurchasePackage_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.PurchasePackageRequest{
		PackageType: 2,
		Duration:    1,
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 设置期望的方法调用，模拟用户未找到
	mockUserService.EXPECT().GetUserByID(gomock.Any(), userId).Return(nil, v1.ErrNotFound)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/user/purchase", userHandler.PurchasePackage)

	obj := newHttpExcept(t, testRouter).POST("/user/purchase").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().
		Object()
	// 验证返回的错误码不为0
	obj.Value("code").Number().Gt(0)
}

func TestUserHandler_GetUsage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟使用量数据
	usageData := &v1.GetUsageResponseData{
		Usages: []v1.UsageData{
			{Date: "2024-01", Usage: 1024},
			{Date: "2024-02", Usage: 2048},
		},
	}

	req := v1.GetUsageRequest{
		UserId: userId,
		Range:  "12months",
	}

	// 设置期望的方法调用
	mockUsageService.EXPECT().GetUsage(gomock.Any(), &req).Return(usageData, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.NoStrictAuth(jwt, logger))
	testRouter.GET("/usage", userHandler.GetUsage)

	obj := newHttpExcept(t, testRouter).GET("/usage").
		WithQuery("range", "12months").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("usages").Array().Length().IsEqual(2)
}

func TestUserHandler_GetVNetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟虚拟网络数据
	vnets := &[]model.Vnet{
		{
			VnetId:        "vnet1",
			UserId:        userId,
			Comment:       "测试网络1",
			Enabled:       true,
			Token:         "token1",
			Password:      "password1",
			IpRange:       "192.168.1.0/24",
			EnableDHCP:    true,
			ClientsLimit:  5,
			ClientsOnline: 2,
		},
		{
			VnetId:        "vnet2",
			UserId:        userId,
			Comment:       "测试网络2",
			Enabled:       false,
			Token:         "token2",
			Password:      "password2",
			IpRange:       "192.168.2.0/24",
			EnableDHCP:    false,
			ClientsLimit:  10,
			ClientsOnline: 0,
		},
	}

	// 设置期望的方法调用
	mockVnetService.EXPECT().GetVnetByUserId(gomock.Any(), userId).Return(vnets, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.GET("/vnet", userHandler.GetVNetList)

	obj := newHttpExcept(t, testRouter).GET("/vnet").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("vnets").Array().Length().IsEqual(2)
}

func TestUserHandler_CreateVNet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.CreateVnetRequest{
		VnetProfile: v1.VnetProfile{
			Comment:      "新建测试网络",
			Enabled:      true,
			Token:        "newtoken",
			Password:     "newpassword",
			IpRange:      "192.168.100.0/24",
			EnableDHCP:   true,
			ClientsLimit: 5,
		},
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟用户信息
	futureTime := time.Now().Add(30 * 24 * time.Hour) // 30天后过期
	currentUser := &model.User{
		UserId:          userId,
		Username:        "testuser",
		Email:           "test@gmail.com",
		UserGroup:       2,           // 青铜用户
		PrivilegeExpiry: &futureTime, // 设置未过期的特权
	}

	// 设置期望的方法调用
	mockVnetService.EXPECT().CheckVnetTokenExists(gomock.Any(), params.Token, "").Return(false, nil)
	mockUserService.EXPECT().GetUserByID(gomock.Any(), userId).Return(currentUser, nil)
	mockVnetService.EXPECT().GetRunningVnetCount(gomock.Any(), userId).Return(1, nil)
	mockVnetService.EXPECT().CreateVnet(gomock.Any(), gomock.Any(), userId).Return(nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/vnet", userHandler.CreateVNet)

	obj := newHttpExcept(t, testRouter).POST("/vnet").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_CreateVNet_TokenExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.CreateVnetRequest{
		VnetProfile: v1.VnetProfile{
			Comment:      "新建测试网络",
			Enabled:      true,
			Token:        "existingtoken",
			Password:     "newpassword",
			IpRange:      "192.168.100.0/24",
			EnableDHCP:   true,
			ClientsLimit: 5,
		},
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 设置期望的方法调用 - token已存在
	mockVnetService.EXPECT().CheckVnetTokenExists(gomock.Any(), params.Token, "").Return(true, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.POST("/vnet", userHandler.CreateVNet)

	obj := newHttpExcept(t, testRouter).POST("/vnet").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object()
	obj.Value("code").Number().Gt(0)
}

func TestUserHandler_UpdateVNet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vnetId := "vnet1"
	params := v1.UpdateVnetRequest{
		VnetProfile: v1.VnetProfile{
			Comment:      "更新的测试网络",
			Enabled:      true,
			Token:        "updatedtoken",
			Password:     "updatedpassword",
			IpRange:      "192.168.101.0/24",
			EnableDHCP:   false,
			ClientsLimit: 8,
		},
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟用户信息
	futureTime := time.Now().Add(30 * 24 * time.Hour) // 30天后过期
	currentUser := &model.User{
		UserId:          userId,
		Username:        "testuser",
		Email:           "test@gmail.com",
		UserGroup:       3,           // 白银用户
		PrivilegeExpiry: &futureTime, // 设置未过期的特权
	}

	// 模拟现有的虚拟网络
	existingVnet := &model.Vnet{
		VnetId:   vnetId,
		UserId:   userId,
		Enabled:  false,
		Comment:  "原始网络",
		Token:    "oldtoken",
		Password: "oldpassword",
	}

	// 设置期望的方法调用
	mockVnetService.EXPECT().GetVnetByVnetId(gomock.Any(), vnetId).Return(existingVnet, nil)
	mockUserService.EXPECT().GetUserByID(gomock.Any(), userId).Return(currentUser, nil)
	mockVnetService.EXPECT().GetRunningVnetCount(gomock.Any(), userId).Return(1, nil)
	mockVnetService.EXPECT().CheckVnetTokenExists(gomock.Any(), params.Token, vnetId).Return(false, nil)
	mockVnetService.EXPECT().UpdateVnet(gomock.Any(), gomock.Any()).Return(nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.PUT("/vnet/:vnetId", userHandler.UpdateVNet)

	obj := newHttpExcept(t, testRouter).PUT("/vnet/"+vnetId).
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_UpdateVNet_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vnetId := "nonexistent"
	params := v1.UpdateVnetRequest{
		VnetProfile: v1.VnetProfile{
			Comment:      "更新的测试网络",
			Enabled:      true,
			Token:        "updatedtoken",
			Password:     "updatedpassword",
			IpRange:      "192.168.101.0/24",
			EnableDHCP:   false,
			ClientsLimit: 8,
		},
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 设置期望的方法调用 - vnet不存在
	mockVnetService.EXPECT().GetVnetByVnetId(gomock.Any(), vnetId).Return(nil, v1.ErrNotFound)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.PUT("/vnet/:vnetId", userHandler.UpdateVNet)

	obj := newHttpExcept(t, testRouter).PUT("/vnet/"+vnetId).
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusNotFound).
		JSON().
		Object()
	obj.Value("code").Number().Gt(0)
}

func TestUserHandler_DeleteVNet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vnetId := "vnet1"

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟现有的虚拟网络
	existingVnet := &model.Vnet{
		VnetId:   vnetId,
		UserId:   userId,
		Enabled:  true,
		Comment:  "待删除的网络",
		Token:    "deletetoken",
		Password: "deletepassword",
	}

	// 设置期望的方法调用
	mockVnetService.EXPECT().GetVnetByVnetId(gomock.Any(), vnetId).Return(existingVnet, nil)
	mockVnetService.EXPECT().DeleteVnet(gomock.Any(), gomock.Any()).Return(nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.DELETE("/vnet/:vnetId", userHandler.DeleteVNet)

	obj := newHttpExcept(t, testRouter).DELETE("/vnet/"+vnetId).
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_DeleteVNet_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vnetId := "vnet1"

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟其他用户的虚拟网络
	existingVnet := &model.Vnet{
		VnetId:   vnetId,
		UserId:   "otheruser", // 不同的用户
		Enabled:  true,
		Comment:  "其他用户的网络",
		Token:    "othertoken",
		Password: "otherpassword",
	}

	// 设置期望的方法调用
	mockVnetService.EXPECT().GetVnetByVnetId(gomock.Any(), vnetId).Return(existingVnet, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.DELETE("/vnet/:vnetId", userHandler.DeleteVNet)

	obj := newHttpExcept(t, testRouter).DELETE("/vnet/"+vnetId).
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusForbidden).
		JSON().
		Object()
	obj.Value("code").Number().Gt(0)
}

func TestUserHandler_GetUserGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟用户基本信息
	profileData := &v1.GetProfileResponseData{
		UserId:    userId,
		Username:  "testuser",
		Email:     "test@gmail.com",
		UserGroup: 3, // 白银用户
	}

	// 设置期望的方法调用
	mockUserService.EXPECT().GetProfile(gomock.Any(), userId).Return(profileData, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.NoStrictAuth(jwt, logger))
	testRouter.GET("/user/group", userHandler.GetUserGroup)

	obj := newHttpExcept(t, testRouter).GET("/user/group").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("userGroup").IsEqual(3)
}

func TestUserHandler_GetVNetLimitInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 模拟用户信息
	futureTime := time.Now().Add(30 * 24 * time.Hour) // 30天后过期
	currentUser := &model.User{
		UserId:          userId,
		Username:        "testuser",
		Email:           "test@gmail.com",
		UserGroup:       2,           // 青铜用户
		PrivilegeExpiry: &futureTime, // 设置未过期的特权
	}

	currentRunningCount := 2

	// 设置期望的方法调用
	mockUserService.EXPECT().GetUserByID(gomock.Any(), userId).Return(currentUser, nil)
	mockVnetService.EXPECT().GetRunningVnetCount(gomock.Any(), userId).Return(currentRunningCount, nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.GET("/vnet/limit", userHandler.GetVNetLimitInfo)

	obj := newHttpExcept(t, testRouter).GET("/vnet/limit").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("currentCount").IsEqual(currentRunningCount)
	objData.Value("userGroup").IsEqual(2)
}

func TestUserHandler_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.ChangePasswordRequest{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 设置期望的方法调用
	mockUserService.EXPECT().ChangePassword(gomock.Any(), userId, &params).Return(nil)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.PUT("/user/password", userHandler.ChangePassword)

	obj := newHttpExcept(t, testRouter).PUT("/user/password").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_ChangePassword_WrongCurrentPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.ChangePasswordRequest{
		CurrentPassword: "wrongpassword",
		NewPassword:     "newpassword",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUsageService := mock_service.NewMockUsageService(ctrl)
	mockVnetService := mock_service.NewMockVnetService(ctrl)

	// 设置期望的方法调用 - 返回密码错误
	mockUserService.EXPECT().ChangePassword(gomock.Any(), userId, &params).Return(v1.ErrOriginalPasswordNotMatch)

	testRouter := createTestRouter()

	userHandler := handler.NewUserHandler(hdl, mockUserService, mockUsageService, mockVnetService)
	testRouter.Use(middleware.StrictAuth(jwt, logger))
	testRouter.PUT("/user/password", userHandler.ChangePassword)

	obj := newHttpExcept(t, testRouter).PUT("/user/password").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object()
	obj.Value("code").Number().Gt(0)
}
