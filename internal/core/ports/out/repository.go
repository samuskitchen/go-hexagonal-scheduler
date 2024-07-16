package out

import (
	"context"
	"go-hexagonal-scheduler/internal/core/domain"
)

type TransactionRepository interface {
	GetTransactionsWithProcessOk(ctx context.Context) ([]domain.TransactionResponse, error)
}
