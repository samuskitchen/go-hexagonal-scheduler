package errors

var (
	TransactionDoesNotExist = ContactPhoneError{Msg: "transaction data not found"}
	TransactionErrorGetting = ContactPhoneError{Msg: "error getting transactions"}
)

type ContactPhoneError struct {
	Msg string
}

func (err ContactPhoneError) Error() string {
	return err.Msg
}
