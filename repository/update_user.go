package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/WalletService/model"
)

func (r *Repository) UpdateUser(ctx context.Context, in model.User) error {
	var (
		query     string
		setFields []string
		params    []interface{}
		offset    int = 0
	)

	if in.PhoneNumber != "" {
		setFields = append(setFields, fmt.Sprintf(setUserPhoneNumberF, offset+1))
		params = append(
			params,
			in.PhoneNumber,
		)
		offset++
	}

	if in.FullName != "" {
		setFields = append(setFields, fmt.Sprintf(setUserFullNameF, offset+1))
		params = append(
			params,
			in.FullName,
		)
		offset++
	}

	if in.Balance > 0 {
		setFields = append(setFields, fmt.Sprintf(setUserBalanceF, offset+1))
		params = append(
			params,
			in.Balance,
		)
		offset++
	}

	setFields = append(setFields, fmt.Sprintf(setUserUpdatedTimeF, offset+1))
	params = append(
		params,
		time.Now(),
	)
	offset++

	query = fmt.Sprintf(queryUpdateUserF, strings.Join(setFields, ","))

	query += fmt.Sprintf(whereUserID, offset+1)
	params = append(
		params,
		in.ID,
	)
	offset++

	var (
		result sql.Result
		err    error
	)
	if result, err = r.exec.ExecContext(ctx, query, params...); err != nil {
		return err
	}

	// Check the affected rows count
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// No rows updated means user does not exist
	if affectedRows == 0 {
		return errors.New("user not found")
	}

	return nil
}
