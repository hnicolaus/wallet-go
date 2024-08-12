package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/WalletService/model"
	"github.com/WalletService/repository"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

var (
	convertToUUID = func(in string) uuid.UUID {
		result, _ := uuid.Parse(in)
		return result
	}
)

func Test_performTransferOut(t *testing.T) {
	tests := []struct {
		name              string
		inputUser         model.User
		inputTransaction  model.Transaction
		mockRepository    func(controller *gomock.Controller) *repository.MockRepositoryInterface
		wantTransactionID uuid.UUID
		wantErr           bool
	}{
		{
			name: "success",
			inputUser: model.User{
				ID:          1234,
				FullName:    "User",
				PhoneNumber: "+628123456789",
				Password:    "$2a$12$35ELZtgOq3iFR6awq.jsDuV5Dr.0XU5k7iUQuShfeLTWRHGFr//fq",
				Balance:     1000000,
			},
			inputTransaction: model.Transaction{
				UserID:      1234,
				Amount:      250000,
				RecipientID: 6789,
				Type:        model.TransactionTypeTransferOut,
				Description: "Traktir Makan",
				Password:    "Admin1234!",
			},
			mockRepository: func(ctrl *gomock.Controller) *repository.MockRepositoryInterface {
				m := repository.NewMockRepositoryInterface(ctrl)

				sqlDb := repository.NewMockSqlDbInterface(ctrl)
				sqlTx := repository.NewMockSqlTxInterface(ctrl)

				// sql_txn operations
				m.EXPECT().GetSqlDb().Return(sqlDb, nil).Times(1)
				sqlDb.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(sqlTx, nil).Times(1)
				m.EXPECT().SetExecutor(sqlTx).Times(1)
				sqlTx.EXPECT().Commit().Return(nil).Times(1)
				m.EXPECT().SetExecutor(sqlDb).Times(1)

				m.EXPECT().GetUser(gomock.Any(), int64(6789)).Return(model.User{
					ID: 6789,
				}, nil).Times(1)

				m.EXPECT().LockUser(gomock.Any(), int64(1234)).Return(nil).Times(1)

				m.EXPECT().UpdateUser(gomock.Any(), model.UpdateUserRequest{
					UserID: 1234,
					Balance: model.UpdateBalanceRequest{
						Amount: 250000,
						Type:   model.UpdateBalanceDecrement,
					},
				}).Return(nil).Times(1)

				m.EXPECT().LockUser(gomock.Any(), int64(6789)).Return(nil).Times(1)

				m.EXPECT().UpdateUser(gomock.Any(), model.UpdateUserRequest{
					UserID: 6789,
					Balance: model.UpdateBalanceRequest{
						Amount: 250000,
						Type:   model.UpdateBalanceIncrement,
					},
				}).Return(nil).Times(1)

				m.EXPECT().InsertTransaction(gomock.Any(), model.Transaction{
					UserID:      1234,
					Amount:      250000,
					RecipientID: 6789,
					Type:        model.TransactionTypeTransferOut,
					Status:      model.TransactionStatusSuccessful,
					Description: "Traktir Makan",
				}).Return(convertToUUID("3d6e668f-ad02-40ff-8540-90c1528a7c88"), nil).Times(1)

				return m
			},
			wantTransactionID: convertToUUID("3d6e668f-ad02-40ff-8540-90c1528a7c88"),
			wantErr:           false,
		},
		{
			name: "failed-repo-call-should-rollback",
			inputUser: model.User{
				ID:          1234,
				FullName:    "User",
				PhoneNumber: "+628123456789",
				Password:    "$2a$12$35ELZtgOq3iFR6awq.jsDuV5Dr.0XU5k7iUQuShfeLTWRHGFr//fq",
				Balance:     1000000,
			},
			inputTransaction: model.Transaction{
				UserID:      1234,
				Amount:      250000,
				RecipientID: 6789,
				Type:        model.TransactionTypeTransferOut,
				Description: "Traktir Makan",
				Password:    "Admin1234!",
			},
			mockRepository: func(ctrl *gomock.Controller) *repository.MockRepositoryInterface {
				m := repository.NewMockRepositoryInterface(ctrl)

				sqlDb := repository.NewMockSqlDbInterface(ctrl)
				sqlTx := repository.NewMockSqlTxInterface(ctrl)

				// sql_txn operations
				m.EXPECT().GetSqlDb().Return(sqlDb, nil).Times(1)
				sqlDb.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(sqlTx, nil).Times(1)
				m.EXPECT().SetExecutor(sqlTx).Times(1)
				sqlTx.EXPECT().Rollback().Return(nil).Times(1)
				m.EXPECT().SetExecutor(sqlDb).Times(1)

				m.EXPECT().GetUser(gomock.Any(), int64(6789)).Return(model.User{
					ID: 6789,
				}, nil).Times(1)

				m.EXPECT().LockUser(gomock.Any(), int64(1234)).Return(nil).Times(1)

				m.EXPECT().UpdateUser(gomock.Any(), model.UpdateUserRequest{
					UserID: 1234,
					Balance: model.UpdateBalanceRequest{
						Amount: 250000,
						Type:   model.UpdateBalanceDecrement,
					},
				}).Return(errors.New("error-update-user")).Times(1)

				m.EXPECT().InsertTransaction(gomock.Any(), model.Transaction{
					UserID:      1234,
					Amount:      250000,
					RecipientID: 6789,
					Type:        model.TransactionTypeTransferOut,
					Status:      model.TransactionStatusFailed,
					Description: "Traktir Makan",
				}).Return(convertToUUID("3d6e668f-ad02-40ff-8540-90c1528a7c88"), nil).Times(1)

				return m
			},
			wantTransactionID: uuid.Nil,
			wantErr:           true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)

			usecase := &Usecase{
				Repository: test.mockRepository(controller),
			}

			gotTransactionID, gotErr := usecase.performTransferOut(context.Background(), test.inputUser, test.inputTransaction)
			if gotTransactionID != test.wantTransactionID {
				t.Errorf("usecase.performTransferOut() gotTransactionID = %v, wantTransactionID %v", gotTransactionID, test.wantTransactionID)
				return
			}
			if (gotErr != nil) != test.wantErr {
				t.Errorf("usecase.performTransferOut() gotErr = %v, wantErr %v", gotErr, test.wantErr)
				return
			}
		})
	}
}

func Test_performTopUp(t *testing.T) {
	tests := []struct {
		name              string
		inputUser         model.User
		inputTransaction  model.Transaction
		mockRepository    func(controller *gomock.Controller) *repository.MockRepositoryInterface
		wantTransactionID uuid.UUID
		wantErr           bool
	}{
		{
			name: "success",
			inputUser: model.User{
				ID:          1234,
				FullName:    "User",
				PhoneNumber: "+628123456789",
				Password:    "$2a$12$35ELZtgOq3iFR6awq.jsDuV5Dr.0XU5k7iUQuShfeLTWRHGFr//fq",
				Balance:     1000000,
			},
			inputTransaction: model.Transaction{
				UserID:      1234,
				Amount:      250000,
				RecipientID: 6789,
				Type:        model.TransactionTypeTopUp,
				Description: "Top Up",
			},
			mockRepository: func(ctrl *gomock.Controller) *repository.MockRepositoryInterface {
				m := repository.NewMockRepositoryInterface(ctrl)

				sqlDb := repository.NewMockSqlDbInterface(ctrl)
				sqlTx := repository.NewMockSqlTxInterface(ctrl)

				// sql_txn operations
				m.EXPECT().GetSqlDb().Return(sqlDb, nil).Times(1)
				sqlDb.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(sqlTx, nil).Times(1)
				m.EXPECT().SetExecutor(sqlTx).Times(1)
				sqlTx.EXPECT().Commit().Return(nil).Times(1)
				m.EXPECT().SetExecutor(sqlDb).Times(1)

				m.EXPECT().LockUser(gomock.Any(), int64(1234)).Return(nil).Times(1)

				m.EXPECT().UpdateUser(gomock.Any(), model.UpdateUserRequest{
					UserID: 1234,
					Balance: model.UpdateBalanceRequest{
						Amount: 250000,
						Type:   model.UpdateBalanceIncrement,
					},
				}).Return(nil).Times(1)

				m.EXPECT().InsertTransaction(gomock.Any(), model.Transaction{
					UserID:      1234,
					Amount:      250000,
					RecipientID: 1234,
					Type:        model.TransactionTypeTopUp,
					Status:      model.TransactionStatusSuccessful,
					Description: "Top Up",
				}).Return(convertToUUID("3d6e668f-ad02-40ff-8540-90c1528a7c88"), nil).Times(1)

				return m
			},
			wantTransactionID: convertToUUID("3d6e668f-ad02-40ff-8540-90c1528a7c88"),
			wantErr:           false,
		},
		{
			name: "failed-repo-call-should-rollback",
			inputUser: model.User{
				ID:          1234,
				FullName:    "User",
				PhoneNumber: "+628123456789",
				Password:    "$2a$12$35ELZtgOq3iFR6awq.jsDuV5Dr.0XU5k7iUQuShfeLTWRHGFr//fq",
				Balance:     1000000,
			},
			inputTransaction: model.Transaction{
				UserID:      1234,
				Amount:      250000,
				RecipientID: 6789,
				Type:        model.TransactionTypeTopUp,
				Description: "Top Up",
			},
			mockRepository: func(ctrl *gomock.Controller) *repository.MockRepositoryInterface {
				m := repository.NewMockRepositoryInterface(ctrl)

				sqlDb := repository.NewMockSqlDbInterface(ctrl)
				sqlTx := repository.NewMockSqlTxInterface(ctrl)

				// sql_txn operations
				m.EXPECT().GetSqlDb().Return(sqlDb, nil).Times(1)
				sqlDb.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(sqlTx, nil).Times(1)
				m.EXPECT().SetExecutor(sqlTx).Times(1)
				sqlTx.EXPECT().Rollback().Return(nil).Times(1)
				m.EXPECT().SetExecutor(sqlDb).Times(1)

				m.EXPECT().LockUser(gomock.Any(), int64(1234)).Return(nil).Times(1)

				m.EXPECT().UpdateUser(gomock.Any(), model.UpdateUserRequest{
					UserID: 1234,
					Balance: model.UpdateBalanceRequest{
						Amount: 250000,
						Type:   model.UpdateBalanceIncrement,
					},
				}).Return(errors.New("error-update-user")).Times(1)

				m.EXPECT().InsertTransaction(gomock.Any(), model.Transaction{
					UserID:      1234,
					Amount:      250000,
					RecipientID: 1234,
					Type:        model.TransactionTypeTopUp,
					Status:      model.TransactionStatusFailed,
					Description: "Top Up",
				}).Return(convertToUUID("3d6e668f-ad02-40ff-8540-90c1528a7c88"), nil).Times(1)

				return m
			},
			wantTransactionID: uuid.Nil,
			wantErr:           true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)

			usecase := &Usecase{
				Repository: test.mockRepository(controller),
			}

			gotTransactionID, gotErr := usecase.performTopUp(context.Background(), test.inputUser, test.inputTransaction)
			if gotTransactionID != test.wantTransactionID {
				t.Errorf("usecase.performTransferOut() gotTransactionID = %v, wantTransactionID %v", gotTransactionID, test.wantTransactionID)
				return
			}
			if (gotErr != nil) != test.wantErr {
				t.Errorf("usecase.performTransferOut() gotErr = %v, wantErr %v", gotErr, test.wantErr)
				return
			}
		})
	}
}
