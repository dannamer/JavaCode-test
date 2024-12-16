package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type Transaction struct {
	WalletID      uuid.UUID       `json:"walletId"`
	OperationType OperationType   `json:"operationType"`
	Amount        decimal.Decimal `json:"amount"`
}

func (t *Transaction) ValidateWalletID() bool {
	return t.WalletID != uuid.Nil
}

func (t *Transaction) ValidateOperationType() bool {
	return t.OperationType == Deposit || t.OperationType == Withdraw
}

func (t *Transaction) ValidateAmount() bool {
	return t.Amount.GreaterThan(decimal.Zero)
}

func (t *Transaction) Validate() bool {
	return t.ValidateWalletID() && t.ValidateOperationType() && t.ValidateAmount()
}
