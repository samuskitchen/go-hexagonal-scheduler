package repository

import (
	"context"

	"go-hexagonal-scheduler/internal/core/domain"
	"go-hexagonal-scheduler/internal/core/ports/out"
	"go-hexagonal-scheduler/pkg/kit/enums"
	errorsCustom "go-hexagonal-scheduler/pkg/kit/errors"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionRepository struct {
	mongo *mongo.Client
}

func NewTransactionRepository(mongoClient *mongo.Client) out.TransactionRepository {
	return &transactionRepository{
		mongo: mongoClient,
	}
}

func (tr *transactionRepository) GetTransactionsWithProcessOk(ctx context.Context) ([]domain.TransactionResponse, error) {
	// Get the database name from viper configuration
	databaseName := enums.MongodbDatabase

	// Get the MongoDB collection
	collection := tr.mongo.Database(databaseName).Collection(enums.PreOrderTransaction)

	// Define the MongoDB aggregation pipeline
	pipeline := []bson.D{
		// Filter
		matchTransaction(),

		// select (pre-order-transaction)
		projectOne(),

		// Perform a left outer join with pre-order-response
		lookup(),

		// unwind
		unwind(),

		// limit registers
		limit(),

		// select result (pre-order-transaction)
		projectOne(),
	}

	// Execute the aggregation pipeline
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Error().Msgf("[collection.Aggregate] %s", err.Error())
		return []domain.TransactionResponse{}, errorsCustom.TransactionErrorGetting
	}

	// Decode the result into a slice of domain.TransactionResponse
	// var transactionResponse []domain.TransactionResponse
	var listOfTransactionResponse []domain.TransactionResponse
	if err = cursor.All(ctx, &listOfTransactionResponse); err != nil {
		log.Error().Msgf("Error cursor All %v:", err)
		return []domain.TransactionResponse{}, errorsCustom.TransactionErrorGetting
	}

	// Check if any customer properties were found
	if len(listOfTransactionResponse) == 0 {
		return []domain.TransactionResponse{}, errorsCustom.TransactionDoesNotExist
	}

	return listOfTransactionResponse, nil
}

// matchClientID Filter in transaction
func matchTransaction() bson.D {
	return bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{Key: "purchaseDate", Value: "2024-07-09"},
				{Key: "cancel", Value: false},
				{Key: "status", Value: "PROCESS_OK"},
			},
		},
	}
}

// projectOne fields
func projectOne() bson.D {
	return bson.D{
		{
			Key: "$project",
			Value: bson.D{
				{Key: "channel", Value: 1},
				{Key: "customer", Value: 1},
				{Key: "messageUniqueID", Value: 1},
				{Key: "docType", Value: 1},
				{Key: "voucherID", Value: 1},
				{Key: "route", Value: 1},
			},
		},
	}
}

// lookup fields
func lookup() bson.D {
	return bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: enums.PreOrderResponse},
				{Key: "localField", Value: "messageUniqueID"},
				{Key: "foreignField", Value: "messageUniqueID"},
				{Key: "as", Value: "response"},
			},
		},
	}
}

// unwind fields
func unwind() bson.D {
	return bson.D{
		{
			Key:   "$unwind",
			Value: "$response",
		},
	}
}

// limit fields
func limit() bson.D {
	return bson.D{
		{
			Key:   "$limit",
			Value: 150,
		},
	}
}

// projectOne fields
func projectTwo() bson.D {
	return bson.D{
		{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "channel", Value: 1},
				{Key: "customer", Value: 1},
				{Key: "messageUniqueID", Value: 1},
				{Key: "docType", Value: 1},
				{Key: "voucherID", Value: 1},
				{Key: "route", Value: 1},
				{Key: "country", Value: "$response.country"},
				{Key: "salesDocument", Value: "$response.salesDocument"},
			},
		},
	}
}
