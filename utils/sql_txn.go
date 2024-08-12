package utils

import (
	"context"
	"fmt"

	"github.com/WalletService/repository"
)

// WithDbTx executes the provided function in a single DB txn.
// fn is expected to contain all the operations that need to be performed within a single DB txn.
func WithDbTx(ctx context.Context, dbTxnRepo repository.DbTxnRepoInterface, fn func(ctx context.Context) error) error {
	// Revert back to the original *sql.DB Executor after the txn commit/rollback
	db, err := dbTxnRepo.GetSqlDb()
	if err != nil {
		return err
	}
	defer dbTxnRepo.SetExecutor(db)

	// Start a new DB tx
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Set repository to use the DB tx as the executor
	dbTxnRepo.SetExecutor(tx)

	// Run all operations in a single tx
	err = fn(ctx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %v", err, rbErr)
		}
		return err
	}

	// Commit the tx
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
