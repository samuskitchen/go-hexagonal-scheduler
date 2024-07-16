package cron

import (
	"context"
	"time"

	"go-hexagonal-scheduler/internal/core/ports/in"
	"go-hexagonal-scheduler/pkg/kit"
	"go-hexagonal-scheduler/pkg/kit/enums"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type SchedulerTransaction struct {
	scheduler          *gocron.Scheduler
	transactionService in.TransactionService
	timeZones          map[string]string
}

func NewSchedulerTransaction(transactionService in.TransactionService) *SchedulerTransaction {
	return &SchedulerTransaction{
		transactionService: transactionService,
		scheduler:          gocron.NewScheduler(time.UTC),
		timeZones:          kit.TimeZones,
	}
}

func (sct *SchedulerTransaction) Start() {
	subLogger := log.With().
		Str("method", "cron.Start").
		Logger()
	subLogger.Info().Msg("INIT")

	for country, tz := range sct.timeZones {
		location, err := time.LoadLocation(tz)
		if err != nil {
			subLogger.Error().Msgf("Error loading location for %s: %v", country, err)
			continue
		}

		sct.scheduler.ChangeLocation(location)

		doOne, err := sct.scheduler.Name(enums.TaskNameOne).Every(enums.TaskRunEveryOne).Milliseconds().Do(sct.runJob, country, enums.TaskNameOne)
		if err != nil {
			subLogger.Error().Msgf("Error scheduler name %s: error: %v", doOne.GetName(), err)
			return
		}

		doTwo, err := sct.scheduler.Name(enums.TaskNameTwo).Every(enums.TaskRunEveryTwo).Milliseconds().Do(sct.runJob, country, enums.TaskNameTwo)
		if err != nil {
			subLogger.Error().Msgf("Error scheduler name %s: error: %v", doTwo.GetName(), err)
			return
		}
	}

	subLogger.Info().Msg("FIN_OK")
	sct.scheduler.StartAsync()
}

func (sct *SchedulerTransaction) runJob(country, taskName string) {
	subLogger := log.With().
		Str("method", "cron.runJob").
		Str("country", country).
		Str("taskName", taskName).
		Logger()
	subLogger.Info().Msg("INIT")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	go func(ctx context.Context) {
		err := sct.transactionService.FetchTransactionsWithProcessOk(ctx, country, taskName)
		if err != nil {
			subLogger.Error().Msgf("Error fetching transactions for %s: %v", country, err)
			return
		}

		defer cancel()
	}(ctx)

	subLogger.Info().Msg("FIN_OK")

}
