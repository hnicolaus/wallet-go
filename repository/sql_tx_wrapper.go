// This file contains the repository implementation layer.
package repository

import (
	"context"
	"database/sql"
)

type SqlTx struct {
	tx *sql.Tx
}

func (r *SqlTx) Commit() error {
	return r.tx.Commit()
}

func (r *SqlTx) Rollback() error {
	return r.tx.Rollback()
}

func (r *SqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.tx.ExecContext(ctx, query, args...)
}

func (r *SqlTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.tx.QueryContext(ctx, query, args...)
}

func (r *SqlTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.tx.QueryRowContext(ctx, query, args...)
}
