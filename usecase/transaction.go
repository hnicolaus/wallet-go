package usecase

import (
	"context"
	"errors"

	"github.com/WalletService/model"
	"github.com/WalletService/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *Usecase) CreateUserTransaction(ctx context.Context, transaction model.Transaction) (newTransactionID uuid.UUID, err error) {
	if transaction.Amount <= 0 {
		return uuid.Nil, errors.New("invalid amount")
	}

	// Validate User exists
	user, err := uc.Repository.GetUser(ctx, transaction.UserID)
	if err != nil {
		return uuid.Nil, err
	}

	// Perform Transaction based on its Type
	switch transaction.Type {
	case model.TransactionTypeTransferOut:
		return uc.performTransferOut(ctx, user, transaction)
	case model.TransactionTypeTopUp:
		return uc.performTopUp(ctx, user, transaction)
	default:
		return uuid.Nil, errors.New("unknown transaction type")
	}
}

func (uc *Usecase) performTransferOut(ctx context.Context, user model.User, transaction model.Transaction) (newTransactionID uuid.UUID, err error) {
	// Validate requested password matches with User's password
	inputPassword, errorList := utils.ValidatePassword(&transaction.Password)
	if len(errorList) > 0 {
		return uuid.Nil, errors.New(errorList[0])
	} else if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)) != nil {
		return uuid.Nil, errors.New("invalid password")
	}

	// Validate User has balance to cover transaction
	if user.Balance < transaction.Amount {
		return uuid.Nil, errors.New("balance not enough")
	}

	// Validate requested Recipient's ID should not be 0 for TransferOut Transaction
	if transaction.RecipientID == 0 {
		return uuid.Nil, errors.New("must specify recipient for TransferOut")
	}

	// Validate should not perform TransferOut to oneself
	if transaction.RecipientID == user.ID {
		return uuid.Nil, errors.New("cannot perform TransferOut to self")
	}

	// Validate that Recipient exists
	recipient, err := uc.Repository.GetUser(ctx, transaction.RecipientID)
	if err != nil {
		return uuid.Nil, err
	}

	// Subtract Users's balance
	user.Balance = user.Balance - transaction.Amount

	// Increment Recipient's balance
	recipient.Balance = recipient.Balance + transaction.Amount

	// Perform the following in a single DB transaction to ensure atomicity of all TransferOut operations
	if err = utils.WithDbTx(context.Background(), uc.Repository, func(ctx context.Context) error {
		// 1. Lock User data to prevent race condition
		if err := uc.Repository.LockUser(ctx, user.ID); err != nil {
			return err
		}

		// 2. Subtract User's balance
		if err := uc.Repository.UpdateUser(ctx, model.UpdateUserRequest{
			UserID: user.ID,
			Balance: model.UpdateBalanceRequest{
				Amount: transaction.Amount,
				Type:   model.UpdateBalanceDecrement,
			},
		}); err != nil {
			return err
		}

		// 3. Lock Recipient data to prevent race condition
		if err := uc.Repository.LockUser(ctx, recipient.ID); err != nil {
			return err
		}

		// 4. Increment Recipient's balance
		if err := uc.Repository.UpdateUser(ctx, model.UpdateUserRequest{
			UserID: recipient.ID,
			Balance: model.UpdateBalanceRequest{
				Amount: transaction.Amount,
				Type:   model.UpdateBalanceIncrement,
			},
		}); err != nil {
			return err
		}

		// 5. Insert a Successful Transaction record
		transaction.Status = model.TransactionStatusSuccessful
		if newTransactionID, err = uc.Repository.InsertTransaction(ctx, transaction); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// Insert a Failed transaction record if failed
		transaction.Status = model.TransactionStatusFailed

		// Fire and forget because a Failed log for a TransferOut operation is equal to no log
		go uc.Repository.InsertTransaction(ctx, transaction)

		return uuid.Nil, err
	}

	return newTransactionID, err
}

func (uc *Usecase) performTopUp(ctx context.Context, user model.User, transaction model.Transaction) (newTransactionID uuid.UUID, err error) {
	// Increment User's balance
	user.Balance = user.Balance + transaction.Amount

	// Set RecipientID to User's ID for TopUp Transaction
	transaction.RecipientID = user.ID

	// Perform the following in a single DB transaction to ensure atomicity of all TopUp operations
	if err = utils.WithDbTx(context.Background(), uc.Repository, func(ctx context.Context) error {
		// 1. Lock User data to prevent race condition
		if err := uc.Repository.LockUser(ctx, user.ID); err != nil {
			return err
		}

		// 2. Increment User's balance
		if err := uc.Repository.UpdateUser(ctx, model.UpdateUserRequest{
			UserID: user.ID,
			Balance: model.UpdateBalanceRequest{
				Amount: transaction.Amount,
				Type:   model.UpdateBalanceIncrement,
			},
		}); err != nil {
			return err
		}

		// 3. Insert a Successful Transaction record
		transaction.Status = model.TransactionStatusSuccessful
		if newTransactionID, err = uc.Repository.InsertTransaction(ctx, transaction); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// Insert a Failed transaction record if failed
		transaction.Status = model.TransactionStatusFailed
		uc.Repository.InsertTransaction(ctx, transaction)

		return uuid.Nil, err
	}

	return newTransactionID, err
}
