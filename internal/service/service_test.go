package service

import (
	"context"
	"errors"
	"testing"

	"github.com/dannamer/JavaCode-test/internal/model"
	"github.com/dannamer/JavaCode-test/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWalletService_WalletTransaction_DepositSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoWallet(ctrl)
	walletUUID := uuid.New()
	initialBalance := decimal.NewFromInt32(100)
	transaction := model.Transaction{
		WalletID:      walletUUID,
		OperationType: model.Deposit,
		Amount:        decimal.NewFromInt32(50),
	}

	mockRepo.EXPECT().GetWallet(context.Background(), walletUUID).Return(model.Wallet{
		UUID:    walletUUID,
		Balance: initialBalance,
	}, nil)

	mockRepo.EXPECT().ProcessTransaction(context.Background(), gomock.Any(), transaction).Return(walletUUID, nil)

	walletService := NewWalletService(mockRepo)

	err := walletService.WalletTransaction(context.Background(), transaction)

	assert.NoError(t, err)
}

func TestWalletService_WalletTransaction_WithdrawSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoWallet(ctrl)
	walletUUID := uuid.New()
	initialBalance := decimal.NewFromInt32(100)
	transaction := model.Transaction{
		WalletID:      walletUUID,
		OperationType: model.Withdraw,
		Amount:        decimal.NewFromInt32(50),
	}

	mockRepo.EXPECT().GetWallet(context.Background(), walletUUID).Return(model.Wallet{
		UUID:    walletUUID,
		Balance: initialBalance,
	}, nil)

	mockRepo.EXPECT().ProcessTransaction(context.Background(), gomock.Any(), transaction).Return(walletUUID, nil)

	walletService := NewWalletService(mockRepo)

	err := walletService.WalletTransaction(context.Background(), transaction)

	assert.NoError(t, err)
}

func TestWalletService_WalletTransaction_WithdrawInsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoWallet(ctrl)
	walletUUID := uuid.New()
	initialBalance := decimal.NewFromInt32(30)
	transaction := model.Transaction{
		WalletID:      walletUUID,
		OperationType: model.Withdraw,
		Amount:        decimal.NewFromInt32(50),
	}

	mockRepo.EXPECT().GetWallet(context.Background(), walletUUID).Return(model.Wallet{
		UUID:    walletUUID,
		Balance: initialBalance,
	}, nil)

	walletService := NewWalletService(mockRepo)

	err := walletService.WalletTransaction(context.Background(), transaction)

	assert.EqualError(t, err, "insufficient funds")
}

func TestWalletService_WalletTransaction_GetWalletError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoWallet(ctrl)
	walletUUID := uuid.New()
	transaction := model.Transaction{
		WalletID:      walletUUID,
		OperationType: model.Deposit,
		Amount:        decimal.NewFromInt32(50),
	}

	mockRepo.EXPECT().GetWallet(context.Background(), walletUUID).Return(model.Wallet{}, errors.New("no rows in result set"))

	walletService := NewWalletService(mockRepo)

	err := walletService.WalletTransaction(context.Background(), transaction)

	assert.EqualError(t, err, "no rows in result set")
}

func TestWalletService_WalletTransaction_ProcessTransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoWallet(ctrl)
	walletUUID := uuid.New()
	initialBalance := decimal.NewFromInt32(100)
	transaction := model.Transaction{
		WalletID:      walletUUID,
		OperationType: model.Deposit,
		Amount:        decimal.NewFromInt32(50),
	}

	mockRepo.EXPECT().GetWallet(context.Background(), walletUUID).Return(model.Wallet{
		UUID:    walletUUID,
		Balance: initialBalance,
	}, nil)

	mockRepo.EXPECT().ProcessTransaction(context.Background(), gomock.Any(), transaction).Return(uuid.Nil, errors.New("process error"))

	walletService := NewWalletService(mockRepo)

	err := walletService.WalletTransaction(context.Background(), transaction)

	assert.EqualError(t, err, "process error")
}

func TestWalletService_GetWalletBalance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoWallet(ctrl)
	walletUUID := uuid.New()
	expectedWallet := model.Wallet{
		UUID:    walletUUID,
		Balance: decimal.NewFromInt32(100),
	}

	mockRepo.EXPECT().GetWallet(context.Background(), walletUUID).Return(expectedWallet, nil)

	walletService := NewWalletService(mockRepo)

	wallet, err := walletService.GetWalletBalance(context.Background(), walletUUID)

	assert.NoError(t, err)
	assert.Equal(t, expectedWallet, wallet)
}

func TestWalletService_GetWalletBalance_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoWallet(ctrl)
	walletUUID := uuid.New()

	mockRepo.EXPECT().GetWallet(context.Background(), walletUUID).Return(model.Wallet{}, errors.New("no rows in result set"))

	walletService := NewWalletService(mockRepo)

	_, err := walletService.GetWalletBalance(context.Background(), walletUUID)

	assert.EqualError(t, err, "no rows in result set")
}
