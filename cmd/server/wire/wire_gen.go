// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	jwtJWT := jwt.NewJwt(viperViper)
	handlerHandler := handler.NewHandler(logger)
	db := repository.NewDB(viperViper, logger)
	repositoryRepository := repository.NewRepository(logger, db)
	transaction := repository.NewTransaction(repositoryRepository)
	sidSid := sid.NewSid()
	serviceService := service.NewService(transaction, logger, sidSid, jwtJWT)
	userRepository := repository.NewUserRepository(repositoryRepository)
	userService := service.NewUserService(serviceService, userRepository)
	userHandler := handler.NewUserHandler(handlerHandler, userService)
	httpServer := server.NewHTTPServer(logger, viperViper, jwtJWT, userHandler)
	jobJob := job.NewJob(transaction, logger, sidSid)
	userJob := job.NewUserJob(jobJob, userRepository)
	jobServer := server.NewJobServer(logger, userJob)
	appApp := newApp(httpServer, jobServer)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository)

var serviceSet = wire.NewSet(service.NewService, service.NewUserService)

var handlerSet = wire.NewSet(handler.NewHandler, handler.NewUserHandler)

var jobSet = wire.NewSet(job.NewJob, job.NewUserJob)

var serverSet = wire.NewSet(server.NewHTTPServer, server.NewJobServer)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,

) *app.App {
	return app.NewApp(app.WithServer(httpServer, jobServer), app.WithName("demo-server"))
}
