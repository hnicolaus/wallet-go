package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/WalletService/model"
)

func (r *Repository) GetUser(ctx context.Context, userID int64) (user model.User, err error) {
	users, err := r.GetUsers(ctx, model.UserFilter{UserID: userID})
	if err != nil {
		return user, err
	}

	if len(users) == 0 {
		return user, errors.New("user not found")
	}

	return users[0], nil
}

func (r *Repository) GetUsers(ctx context.Context, request model.UserFilter) (users []model.User, err error) {
	query, params, err := buildQueryGetUsers(request)
	if err != nil {
		return []model.User{}, err
	}

	rows, err := r.exec.QueryContext(ctx, query, params...)
	if err != nil {
		return []model.User{}, err
	}

	defer rows.Close()
	for rows.Next() {
		user := model.User{}

		if err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.PhoneNumber,
			&user.Balance,
			&user.Password,
			&user.CreatedTime,
			&user.UpdatedTime,
		); err != nil {
			return []model.User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func buildQueryGetUsers(in model.UserFilter) (string, []interface{}, error) {
	var (
		query  string = querySelectUsers
		params []interface{}
		offset int = 0
	)

	if in.PhoneNumber != "" {
		query += fmt.Sprintf(whereUserPhoneNumber, offset+1)
		params = append(
			params,
			in.PhoneNumber,
		)
		offset++
	}

	if in.UserID != 0 {
		query += fmt.Sprintf(whereUserID, offset+1)
		params = append(
			params,
			strconv.Itoa(int(in.UserID)),
		)
		offset++
	}

	return query, params, nil
}
