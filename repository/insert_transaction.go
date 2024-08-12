package repository

import (
	"context"
	"time"

	"github.com/WalletService/model"
	"github.com/google/uuid"
)

func (r *Repository) InsertTransaction(ctx context.Context, transaction model.Transaction) (transactionID uuid.UUID, err error) {
	var (
		params []interface{}
	)

	params = append(
		params,
		uuid.New(),
		transaction.UserID,
		transaction.Amount,
		transaction.Type,
		transaction.RecipientID,
		transaction.Status,
		transaction.Description,
		time.Now(),
	)

	err = r.exec.QueryRowContext(ctx, queryInsertTransaction, params...).Scan(&transactionID)

	return
}
