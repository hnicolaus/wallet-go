// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
	"database/sql"

	"github.com/WalletService/model"
	"github.com/google/uuid"
)

type RepositoryInterface interface {
	InsertUser(ctx context.Context, user model.User) (userID int64, err error)
	GetUsers(ctx context.Context, request model.UserFilter) (users []model.User, err error)
	GetUser(ctx context.Context, userID int64) (user model.User, err error)
	InsertTransaction(ctx context.Context, transaction model.Transaction) (transactionID uuid.UUID, err error)
	UpdateUser(ctx context.Context, user model.User) error
	DbTxnRepoInterface // to enable using db txn
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type DbTxnRepoInterface interface {
	GetSqlDb() (*sql.DB, error)
	SetExecutor(executor Executor)
}
