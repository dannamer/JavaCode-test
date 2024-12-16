package model

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	StatusBadRequest           = "Bad Request"
	StatusInvalidRequestBody   = "Invalid request body"
	StatusInvalidRequestData   = "Invalid request data. Please check the input parameters."
	StatusInsufficientFunds    = "insufficient funds"
	StatusWalletNotFound       = "Wallet with UUID %s not found"
	StatusInternalServerError  = "Internal Server Error"
	StatusTransactionSuccess   = "Transaction successful"
	StatusWalletBalanceSuccess = "Wallet balance successfully received"
	StatusInvalidUUIDFormat    = "Invalid wallet UUID format."
)