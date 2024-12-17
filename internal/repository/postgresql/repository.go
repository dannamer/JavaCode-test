package postgresql

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/dannamer/JavaCode-test/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type WalletRepo struct {
	PgxPool
}

func NewWalletRepo(postgresql PgxPool) WalletRepo {
	return WalletRepo{PgxPool: postgresql}
}

func (r *WalletRepo) GetWallet(ctx context.Context, UUID uuid.UUID) (model.Wallet, error) {
	sql, args, err := Builder().Select("uuid", "balance", "created_at").
		From("wallets").
		Where(squirrel.Eq{"uuid": UUID}).ToSql()
	if err != nil {
		logrus.Errorf("Failed to build query for GetWallet: %v", err)
		return model.Wallet{}, err
	}

	var wallet model.Wallet
	err = r.PgxPool.QueryRow(ctx, sql, args...).Scan(&wallet.UUID, &wallet.Balance, &wallet.CreatedAt)
	if err != nil {
		logrus.Errorf("Error executing query for GetWallet with UUID %s: %v", UUID, err)
		return model.Wallet{}, err
	}

	return wallet, nil
}

func (r *WalletRepo) ProcessTransaction(ctx context.Context, wallet model.Wallet, transaction model.Transaction) (uuid.UUID, error) {
	tx, err := r.PgxPool.Begin(ctx)
	if err != nil {
		logrus.Errorf("Failed to begin transaction: %v", err)
		return uuid.Nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			logrus.Errorf("Transaction rolled back due to error: %v", err)
		}
	}()

	err = r.UpdatedWallet(ctx, wallet, tx)
	if err != nil {
		logrus.Errorf("Failed to update wallet with UUID %s: %v", wallet.UUID, err)
		return uuid.Nil, err
	}

	transactionUUID, err := r.SaveTransaction(ctx, transaction, tx)
	if err != nil {
		logrus.Errorf("Failed to save transaction for WalletID %s: %v", transaction.WalletID, err)
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		logrus.Errorf("Failed to commit transaction: %v", err)
		return uuid.Nil, err
	}

	return transactionUUID, nil
}

func (r *WalletRepo) UpdatedWallet(ctx context.Context, wallet model.Wallet, tx pgx.Tx) error {
	sql, args, err := Builder().Update("wallets").
		Set("balance", wallet.Balance).
		Where(squirrel.Eq{"uuid": wallet.UUID}).
		ToSql()
	if err != nil {
		logrus.Errorf("Failed to build query for UpdatedWallet: %v", err)
		return err
	}

	res, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		logrus.Errorf("Error executing update query for wallet with UUID %s: %v", wallet.UUID, err)
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		logrus.Errorf("No rows updated for wallet with UUID %s", wallet.UUID)
	}
	return nil
}

func (r *WalletRepo) SaveTransaction(ctx context.Context, transaction model.Transaction, tx pgx.Tx) (uuid.UUID, error) {
	sql, args, err := Builder().Insert("transactions").
		Columns("wallet_uuid", "transaction_type", "amount").
		Values(transaction.WalletID, transaction.OperationType, transaction.Amount).
		Suffix("RETURNING uuid").ToSql()
	if err != nil {
		logrus.Errorf("Failed to build insert query for SaveTransaction: %v", err)
		return uuid.Nil, err
	}

	var transactionUUID uuid.UUID
	err = tx.QueryRow(ctx, sql, args...).Scan(&transactionUUID)
	if err != nil {
		logrus.Errorf("Error saving transaction for wallet %s: %v", transaction.WalletID, err)
		return uuid.Nil, err
	}

	return transactionUUID, nil
}
