package repository

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/WalletService/model"
)

const (
	saltCost = 12
)

func (r *Repository) InsertUser(ctx context.Context, user model.User) (userID int64, err error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), saltCost)
	if err != nil {
		return userID, err
	}

	user.Password = string(hashedPasswordBytes)

	query, params := buildQueryInsertUsers([]model.User{user})

	err = r.exec.QueryRowContext(ctx, query, params...).Scan(&userID)

	return
}

func buildQueryInsertUsers(in []model.User) (string, []interface{}) {
	var (
		query  string = queryInsertUsers
		params []interface{}
		offset int = 0
	)

	for _, row := range in {
		query += fmt.Sprintf(
			valuesInsertUsersF,
			offset+1, offset+2, offset+3, offset+4, offset+5,
		)

		params = append(
			params,
			row.FullName,
			row.PhoneNumber,
			row.Password,
			row.Balance,
		)

		// created_time is DB internal timestamp for new row creation
		params = append(params, time.Now())

		offset = offset + 5
	}

	// trim the last comma
	query = fmt.Sprintln(query[0:len(query)-1], returnLastInsertedUserID)

	return query, params
}
