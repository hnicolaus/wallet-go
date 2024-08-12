package repository

import (
	"context"
)

func (r *Repository) LockUser(ctx context.Context, userID int64) (err error) {
	var (
		params []interface{}
	)

	params = append(
		params,
		userID,
	)

	rows, err := r.exec.QueryContext(ctx, queryLockUser, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}
