// This file contains the repository implementation layer.
package repository

import (
	"context"
	"database/sql"
)

type SqlDb struct {
	db *sql.DB
}

func (r *SqlDb) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

func (r *SqlDb) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, query, args...)
}

func (r *SqlDb) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.db.QueryRowContext(ctx, query, args...)
}

func (r *SqlDb) BeginTx(ctx context.Context, opts *sql.TxOptions) (SqlTxInterface, error) {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &SqlTx{tx: tx}, nil
}
