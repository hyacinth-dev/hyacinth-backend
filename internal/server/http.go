// 在此模块定义请求路由
// 路由分四类：
// 1. 无需权限的路由，如：注册、登录、以及主页数据的展示
// 2. 非严格权限的路由，如：获取用户信息、获取使用情况等
// 3. 严格权限的路由，如：更新用户信息等
// 4. 管理员权限的路由，如：管理用户等

package server

import (
	apiV1 "hyacinth-backend/api/v1"
	"hyacinth-backend/docs"
	"hyacinth-backend/internal/handler"
	"hyacinth-backend/internal/middleware"
	"hyacinth-backend/internal/repository"
	"hyacinth-backend/pkg/jwt"
	"hyacinth-backend/pkg/log"
	"hyacinth-backend/pkg/server/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	userRepo repository.UserRepository,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	v1 := s.Group("/v1")
	{
		// No route group has permission
		noAuthRouter := v1.Group("/")
		{
			noAuthRouter.POST("/register", userHandler.Register)
			noAuthRouter.POST("/login", userHandler.Login)
		}
		// Non-strict permission routing group
		noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		{
			noStrictAuthRouter.GET("/user", userHandler.GetProfile)
			noStrictAuthRouter.GET("/user/usage", userHandler.GetUsage)
			noStrictAuthRouter.GET("/user/vnet", userHandler.GetVNet)
			noStrictAuthRouter.POST("/user/vnet", userHandler.CreateVNet)
			noStrictAuthRouter.DELETE("/user/vnet/:id", userHandler.DeleteVNet)
		}
		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
			strictAuthRouter.PUT("/user/vnet/:VNETID", userHandler.UpdateVNet)
		}
		adminAuthRouter := v1.Group("/admin").Use(middleware.AdminAuth(jwt, logger, userRepo))
		{
			//adminAuthRouter.GET("/user", func(ctx *gin.Context) {})
			adminAuthRouter.GET("/usage/overview", adminHandler.GetTotalUsage)
			adminAuthRouter.GET("/usage/page/:page", adminHandler.GetUsagePage)
			adminAuthRouter.GET("/usage/:id", adminHandler.GetUsage)
			adminAuthRouter.GET("/vnet", adminHandler.AdminGetVNet)
			adminAuthRouter.POST("/vnet/:USERID", adminHandler.AdminCreateVNet)
			adminAuthRouter.PUT("/vnet/:USERID/:VNETID", adminHandler.AdminUpdateVNet)
			adminAuthRouter.DELETE("/vnet/:USERID/:VNETID", adminHandler.AdminDeleteVNet)
		}
	}

	return s
}
