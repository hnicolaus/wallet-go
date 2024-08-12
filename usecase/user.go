package usecase

import (
	"context"
	"errors"

	"github.com/WalletService/model"
	"golang.org/x/crypto/bcrypt"
)

func (uc *Usecase) RegisterUser(ctx context.Context, user model.User) (userID int64, err error) {
	return uc.Repository.InsertUser(ctx, user)
}

func (uc *Usecase) GetUser(ctx context.Context, userID int64) (user model.User, err error) {
	return uc.Repository.GetUser(ctx, userID)
}

func (uc *Usecase) GetUsers(ctx context.Context, request model.UserFilter) (users []model.User, err error) {
	return uc.Repository.GetUsers(ctx, request)
}

func (uc *Usecase) UserLogin(ctx context.Context, phoneNumber, password string) (userID int64, err error) {
	// Get User data
	users, err := uc.Repository.GetUsers(ctx, model.UserFilter{PhoneNumber: phoneNumber})
	if err != nil {
		return 0, err
	}
	if len(users) == 0 {
		return 0, errors.New("user does not exist")
	}

	// Phone number is unique, so expecting only at most 1 user to be retrieved
	user := users[0]

	// Validate input password (plain) matches user's password (hashed and salted)
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return 0, errors.New("invalid password")
	}

	return user.ID, nil
}
