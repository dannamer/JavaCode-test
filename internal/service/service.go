package service

import (
	"context"
	"errors"
	"sync"

	"github.com/dannamer/JavaCode-test/internal/model"
	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mock/service_mock.go -package=mock
type RepoWallet interface {
	GetWallet(ctx context.Context, UUID uuid.UUID) (model.Wallet, error)
	ProcessTransaction(ctx context.Context, wallet model.Wallet, transaction model.Transaction) (uuid.UUID, error)
}

type WalletService struct {
	RepoWallet
	Mu sync.Mutex
}

func NewWalletService(repoWallet RepoWallet) WalletService {
	return WalletService{RepoWallet: repoWallet}
}

func (s *WalletService) WalletTransaction(ctx context.Context, transaction model.Transaction) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	wallet, err := s.GetWallet(ctx, transaction.WalletID)
	if err != nil {
		return err
	}

	if transaction.OperationType == model.Withdraw && wallet.Balance.LessThan(transaction.Amount) {
		return errors.New("insufficient funds")
	}

	switch transaction.OperationType {
	case model.Deposit:
		wallet.Balance = wallet.Balance.Add(transaction.Amount)
	case model.Withdraw:
		wallet.Balance = wallet.Balance.Sub(transaction.Amount)
	}

	if _, err = s.ProcessTransaction(ctx, wallet, transaction); err != nil {
		return err
	}
	return nil
}

func (s *WalletService) GetWalletBalance(ctx context.Context, UUID uuid.UUID) (model.Wallet, error) {
	wallet, err := s.GetWallet(ctx, UUID)
	if err != nil {
		return model.Wallet{}, err
	}
	return wallet, nil
}
