//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"hyacinth-backend/internal/handler"
	"hyacinth-backend/internal/job"
	"hyacinth-backend/internal/repository"
	"hyacinth-backend/internal/server"
	"hyacinth-backend/internal/service"
	"hyacinth-backend/pkg/app"
	"hyacinth-backend/pkg/jwt"
	"hyacinth-backend/pkg/log"
	"hyacinth-backend/pkg/server/http"
	"hyacinth-backend/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,

	repository.NewUsageRepository,
	repository.NewVnetRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,

	service.NewUsageService,
	service.NewAdminService,
	service.NewVnetService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewAdminHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
)
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
