package main

import (
	"os"

	"go-hexagonal-scheduler/internal/core/cron"
	"go-hexagonal-scheduler/internal/infrastructure/configs/injector"
	"go-hexagonal-scheduler/internal/infrastructure/middleware/log"
	"go-hexagonal-scheduler/pkg/kit/enums"

	zerolog "github.com/rs/zerolog/log"
)

func main() {
	container := injector.BuildContainer()

	log.InitLogger(enums.App)

	err := container.Invoke(func(scheduler *cron.SchedulerTransaction) {
		scheduler.Start()
	})

	if err != nil {
		zerolog.Error().Msgf("Failed to start the application: %v", err)
		os.Exit(1)
	}

	select {} // Keep the main goroutine alive
}
