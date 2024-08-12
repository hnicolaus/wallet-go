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

func (r *Repository) UpdateUser(ctx context.Context, in model.UpdateUserRequest) error {
	var (
		query     string
		setFields []string
		params    []interface{}
		offset    int = 0
	)

	if in.Balance != (model.UpdateBalanceRequest{}) {
		switch in.Balance.Type {
		case model.UpdateBalanceIncrement:
			setFields = append(setFields, fmt.Sprintf(incrementUserBalanceF, offset+1))
		case model.UpdateBalanceDecrement:
			setFields = append(setFields, fmt.Sprintf(decrementUserBalanceF, offset+1))

		default:
			return errors.New("unknow balance update type")
		}
		params = append(
			params,
			in.Balance.Amount,
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
		in.UserID,
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
