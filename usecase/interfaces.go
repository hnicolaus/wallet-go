// This file contains the interfaces for the usecase layer.
// The usecase layer is responsible for orchestrating calls to external dependencies, i.e. DB, message queue, cache
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package usecase

import (
	"context"

	"github.com/WalletService/model"
	"github.com/google/uuid"
)

type UsecaseInterface interface {
	RegisterUser(ctx context.Context, user model.User) (userID int64, err error)
	GetUser(ctx context.Context, userID int64) (user model.User, err error)
	GetUsers(ctx context.Context, request model.UserFilter) (users []model.User, err error)
	UserLogin(ctx context.Context, phoneNumber, password string) (userID int64, err error)
	CreateUserTransaction(ctx context.Context, transaction model.Transaction) (newTransactionID uuid.UUID, err error)
}
