package domain

type TransactionResponse struct {
	Channel         string `bson:"channel"`
	Customer        string `bson:"customer"`
	MessageUniqueID string `bson:"messageUniqueID"`
	DocType         string `bson:"docType"`
	VoucherID       string `bson:"voucherID"`
	Route           int64  `bson:"route"`
	Country         string `bson:"country"`
	SalesDocument   string `bson:"salesDocument"`
}
