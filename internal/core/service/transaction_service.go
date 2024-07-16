package service

import (
	"context"

	"go-hexagonal-scheduler/internal/core/ports/in"
	"go-hexagonal-scheduler/internal/core/ports/out"

	"github.com/rs/zerolog/log"
)

type transactionService struct {
	transactionRepository out.TransactionRepository
}

func NewTransactionService(transactionRepository out.TransactionRepository) in.TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
	}
}

func (ts *transactionService) FetchTransactionsWithProcessOk(ctx context.Context, country, taskName string) error {
	subLogger := log.With().
		Str("method", "service.FetchTransactionsWithProcessOk").
		Str("country", country).
		Str("taskName", taskName).
		Logger()
	subLogger.Info().Msg("INIT")

	transactions, err := ts.transactionRepository.GetTransactionsWithProcessOk(ctx)
	if err != nil {
		subLogger.Printf("Error fetching transactions for %s: %v", country, err)
		return err
	}

	subLogger.Printf("Fetched %d transactions for %s", len(transactions), country)
	subLogger.Info().Msg("FIN_OK")
	return err
}
