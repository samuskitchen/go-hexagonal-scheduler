package injector

import (
	"fmt"

	"go-hexagonal-scheduler/internal/adapters/outgoing/repository"
	"go-hexagonal-scheduler/internal/core/cron"
	"go-hexagonal-scheduler/internal/core/service"
	"go-hexagonal-scheduler/internal/infrastructure/configs/storage"

	"go.uber.org/dig"
)

var Container *dig.Container

func BuildContainer() *dig.Container {
	Container = dig.New()

	// DB
	checkError(Container.Provide(storage.ConnInstance))

	// Cron
	checkError(Container.Provide(cron.NewSchedulerTransaction))

	// Service
	checkError(Container.Provide(service.NewTransactionService))

	// Repository
	checkError(Container.Provide(repository.NewTransactionRepository))

	return Container
}

func checkError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Error injecting %v", err))
	}
}
