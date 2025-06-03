// 在此模块定义请求路由
// 路由分三类：
// 1. 无需权限的路由，如：注册、登录、以及主页数据的展示
// 2. 非严格权限的路由，如：获取用户信息、获取使用情况等
// 3. 严格权限的路由，如：更新用户信息等

package server

import (
	apiV1 "hyacinth-backend/api/v1"
	"hyacinth-backend/docs"
	"hyacinth-backend/internal/handler"
	"hyacinth-backend/internal/middleware"
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
			noStrictAuthRouter.GET("/user/group", userHandler.GetUserGroup)
			noStrictAuthRouter.GET("/usage", userHandler.GetUsage)
		}

		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
			strictAuthRouter.PUT("/user/password", userHandler.ChangePassword)
			strictAuthRouter.POST("/user/purchase", userHandler.PurchasePackage)

			// User VNet operations
			strictAuthRouter.GET("/vnet", userHandler.GetVNetList)
			strictAuthRouter.GET("/vnet/limit", userHandler.GetVNetLimitInfo)
			strictAuthRouter.POST("/vnet", userHandler.CreateVNet)
			strictAuthRouter.PUT("/vnet/:vnetId", userHandler.UpdateVNet)
			strictAuthRouter.DELETE("/vnet/:vnetId", userHandler.DeleteVNet)
		}
	}

	return s
}
