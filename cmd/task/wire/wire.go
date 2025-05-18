//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"hyacinth-backend/internal/repository"
	"hyacinth-backend/internal/server"
	"hyacinth-backend/internal/task"
	"hyacinth-backend/pkg/app"
	"hyacinth-backend/pkg/log"
	"hyacinth-backend/pkg/sid"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewUsageRepository,
	repository.NewVNetRepository,
)

var taskSet = wire.NewSet(
	task.NewTask,
	task.NewUserTask,
)
var serverSet = wire.NewSet(
	server.NewTaskServer,
)

// build App
func newApp(
	task *server.TaskServer,
) *app.App {
	return app.NewApp(
		app.WithServer(task),
		app.WithName("demo-task"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		taskSet,
		serverSet,
		newApp,
		sid.NewSid,
	))
}
