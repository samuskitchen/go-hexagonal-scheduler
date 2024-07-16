package in

import (
	"context"
)

type TransactionService interface {
	FetchTransactionsWithProcessOk(ctx context.Context, country, taskName string) error
}
