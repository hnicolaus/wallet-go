package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/WalletService/generated"
	"github.com/WalletService/model"
	"github.com/WalletService/usecase"
	"github.com/WalletService/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

var (
	floatPtr = func(in float32) *float32 {
		return &in
	}

	intPtr = func(in int64) *int64 {
		return &in
	}

	stringPtr = func(in string) *string {
		return &in
	}

	transactionTypePtr = func(t generated.TransactionType) *generated.TransactionType {
		return &t
	}

	convertToUUID = func(in string) uuid.UUID {
		result, _ := uuid.Parse(in)
		return result
	}
)

func TestRegisterUser(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	int64Ptr := func(in int64) *int64 {
		return &in
	}

	errorConflictUserPhoneNumber := pq.Error{
		Code: "23505",
	}

	tests := []struct {
		name                               string
		mockUsecase                        func(controller *gomock.Controller) *usecase.MockUsecaseInterface
		requestBody                        generated.User
		fnConvertRegisterUserRequestToUser func(generated.User) (model.User, []string)
		wantResponse                       generated.RegisterUserResponse
		wantHttpStatusCode                 int
	}{
		{
			name: "success",
			requestBody: generated.User{
				FullName:    stringPtr("User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (model.User, []string) {
				user := model.User{
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().RegisterUser(gomock.Any(), model.User{
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}).Return(int64(123), nil)

				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  true,
					Messages: []string{successMsg},
				},
				User: generated.User{
					Id: int64Ptr(123),
				},
			},
			wantHttpStatusCode: http.StatusCreated,
		},
		{
			name: "fail-insert-user",
			requestBody: generated.User{
				FullName:    stringPtr("User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (model.User, []string) {
				user := model.User{
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().RegisterUser(gomock.Any(), model.User{
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}).Return(int64(0), errors.New("error-insert-user"))

				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"error-insert-user"},
				},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fail-insert-user-conflict-phone-number",
			requestBody: generated.User{
				FullName:    stringPtr("User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (model.User, []string) {
				user := model.User{
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().RegisterUser(gomock.Any(), model.User{
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}).Return(int64(0), &errorConflictUserPhoneNumber)

				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"phone number is already registered to an existing user"},
				},
			},
			wantHttpStatusCode: http.StatusConflict,
		},
		{
			name: "fail-insert-user-conflict-phone-number",
			requestBody: generated.User{
				FullName:    stringPtr("User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (model.User, []string) {
				user := model.User{
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().RegisterUser(gomock.Any(), model.User{
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}).Return(int64(0), &errorConflictUserPhoneNumber)

				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{duplicatePhoneNumberErrorMsg},
				},
			},
			wantHttpStatusCode: http.StatusConflict,
		},
		{
			name: "fail-invalid-input",
			requestBody: generated.User{
				FullName:    stringPtr("User"),
				PhoneNumber: stringPtr("+62812"),
				Password:    stringPtr("P455w"),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (model.User, []string) {
				return model.User{}, []string{"invalid full name", "invalid phone number", "invalid password"}
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)
				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"invalid full name", "invalid phone number", "invalid password"},
				},
			},
			wantHttpStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Usecase: test.mockUsecase(controller),
			}

			requestBodyJSON, _ := json.Marshal(test.requestBody)
			requestBody := []byte(requestBodyJSON)
			requestBodyBuffer := bytes.NewBuffer(requestBody)

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/v1/user", requestBodyBuffer)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)

			fnConvertRegisterUserRequestToUser = test.fnConvertRegisterUserRequestToUser

			gotHttpStatusCode, gotResponse := handler.registerUser(ctx)

			if gotHttpStatusCode != test.wantHttpStatusCode {
				t.Errorf("handler.RegisterUser() httpStatusCode = %v, wantHttpStatusCode %v", gotHttpStatusCode, test.wantHttpStatusCode)
			}

			if !reflect.DeepEqual(test.wantResponse, gotResponse) {
				t.Errorf("handler.RegisterUser() response = %v, wantResponse %v", gotResponse, test.wantResponse)
			}

		})
	}
}

func TestUserLogin(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	int64Ptr := func(in int64) *int64 {
		return &in
	}

	tests := []struct {
		name               string
		mockUsecase        func(controller *gomock.Controller) *usecase.MockUsecaseInterface
		requestBody        generated.User
		wantResponse       generated.UserLoginResponse
		wantCtxUserID      int64
		wantHttpStatusCode int
	}{
		{
			name: "success",
			requestBody: generated.User{
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("Password123!."),
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().UserLogin(gomock.Any(), "+628123456789", "Password123!.").Return(int64(123), nil)

				return mock
			},
			wantResponse: generated.UserLoginResponse{
				Header: generated.ResponseHeader{
					Success:  true,
					Messages: []string{successMsg},
				},
				User: generated.User{
					Id: int64Ptr(123),
				},
			},
			wantCtxUserID:      123,
			wantHttpStatusCode: http.StatusOK,
		},
		{
			name: "fail-usecase-error",
			requestBody: generated.User{
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("Password123.!"),
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().UserLogin(gomock.Any(), "+628123456789", "Password123.!").Return(int64(123), errors.New("invalid password"))

				return mock
			},
			wantResponse: generated.UserLoginResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"invalid password"},
				},
			},
			wantHttpStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Usecase: test.mockUsecase(controller),
			}

			requestBodyJSON, _ := json.Marshal(test.requestBody)
			requestBody := []byte(requestBodyJSON)
			requestBodyBuffer := bytes.NewBuffer(requestBody)

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/v1/user/login", requestBodyBuffer)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)

			gotHttpStatusCode, gotResponse := handler.userLogin(ctx)

			if gotHttpStatusCode != test.wantHttpStatusCode {
				t.Errorf("handler.UserLogin() httpStatusCode = %v, wantHttpStatusCode %v", gotHttpStatusCode, test.wantHttpStatusCode)
			}

			if !reflect.DeepEqual(test.wantResponse, gotResponse) {
				t.Errorf("handler.UserLogin() response = %v, wantResponse %v", gotResponse, test.wantResponse)
			}

			if test.wantResponse.Header.Success {
				gotCtxUserID, _ := ctx.Get(string(utils.JWTClaimUserID)).(int64)
				if gotCtxUserID != test.wantCtxUserID {
					t.Errorf("handler.UserLogin() gotCtxUserID = %v, wantCtxUserID %v", gotCtxUserID, test.wantCtxUserID)
				}

				gotCtxPermissions, _ := ctx.Get(string(utils.JWTClaimPermissions)).([]utils.JWTPermission)
				wantCtxPermissions := []utils.JWTPermission{utils.JWTPermissionGetUser, utils.JWTPermissionPerformTransaction}
				if !reflect.DeepEqual(gotCtxPermissions, wantCtxPermissions) {
					t.Errorf("handler.UserLogin() gotCtxPermissions = %v, wantCtxPermissions %v", gotCtxPermissions, wantCtxPermissions)
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name           string
		mockUsecase    func(controller *gomock.Controller) *usecase.MockUsecaseInterface
		ctxPermissions []utils.JWTPermission
		ctxUserID      int64

		wantResponse       generated.GetUserResponse
		wantHttpStatusCode int
	}{
		{
			name: "success",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			ctxUserID: 123,
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().GetUser(gomock.Any(), int64(123)).Return(model.User{
					ID:          123,
					FullName:    "User",
					PhoneNumber: "+628123456789",
					Balance:     888.90,
					Password:    "$2a$12$41bm0d9VyLDKALovox4S9.FoNezvO9tB8ck94/0fEyKcYIFmV8guq",
				}, nil)

				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  true,
					Messages: []string{successMsg},
				},
				User: generated.User{
					Id:          intPtr(123),
					FullName:    stringPtr("User"),
					PhoneNumber: stringPtr("+628123456789"),
					Balance:     floatPtr(888.90),
				},
			},
			wantHttpStatusCode: http.StatusOK,
		},
		{
			name:           "fail-not-authorized-no-permission",
			ctxPermissions: []utils.JWTPermission{},
			ctxUserID:      123,
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)
				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"not authorized: missing required permission"},
				},
			},
			wantHttpStatusCode: http.StatusForbidden,
		},
		{
			name:           "fail-not-authorized-wrong-permission",
			ctxPermissions: []utils.JWTPermission{utils.JWTPermissionPerformTransaction},
			ctxUserID:      123,
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)
				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"not authorized: missing required permission"},
				},
			},
			wantHttpStatusCode: http.StatusForbidden,
		},
		{
			name:           "fail-not-authorized-no-user-id",
			ctxPermissions: []utils.JWTPermission{utils.JWTPermissionGetUser},
			ctxUserID:      0,
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)
				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"missing user_id"},
				},
			},
			wantHttpStatusCode: http.StatusForbidden,
		},
		{
			name: "fail-get-user",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			ctxUserID: 123,
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().GetUser(gomock.Any(), int64(123)).Return(model.User{}, errors.New("error-get-user"))

				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"error-get-user"},
				},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Usecase: test.mockUsecase(controller),
			}

			e := echo.New()
			request := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)

			if test.ctxUserID != 0 {
				ctx.Set(string(utils.JWTClaimUserID), test.ctxUserID)
			}
			if len(test.ctxPermissions) > 0 {
				ctx.Set(string(utils.JWTClaimPermissions), test.ctxPermissions)
			}

			gotHttpStatusCode, gotResponse := handler.getUser(ctx)

			if gotHttpStatusCode != test.wantHttpStatusCode {
				t.Errorf("handler.GetUser() httpStatusCode = %v, wantHttpStatusCode %v", gotHttpStatusCode, test.wantHttpStatusCode)
			}

			if !reflect.DeepEqual(test.wantResponse, gotResponse) {
				t.Errorf("handler.GetUser() response = %v, wantResponse %v", gotResponse, test.wantResponse)
			}
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name                                           string
		mockUsecase                                    func(controller *gomock.Controller) *usecase.MockUsecaseInterface
		ctxPermissions                                 []utils.JWTPermission
		ctxUserID                                      int64
		requestBody                                    generated.Transaction
		fnConvertCreateTransactionRequestToTransaction func(int64, generated.Transaction) (model.Transaction, []string)

		wantResponse       generated.TransactionResponse
		wantHttpStatusCode int
	}{
		{
			name: "success",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionPerformTransaction,
			},
			ctxUserID: 123,
			requestBody: generated.Transaction{
				Amount:      floatPtr(100000),
				RecipientId: intPtr(2),
				Type:        transactionTypePtr(generated.TransferOut),
				Description: stringPtr("Traktir Makan"),
				Password:    stringPtr("Admin1234!"),
			},
			fnConvertCreateTransactionRequestToTransaction: func(int64, generated.Transaction) (model.Transaction, []string) {
				transaction := model.Transaction{
					Amount:      100000,
					RecipientID: 2,
					Type:        model.TransactionTypeTransferOut,
					Description: "Traktir Makan",
					Password:    "Admin1234!",
				}

				errorMsgs := []string{}

				return transaction, errorMsgs
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().CreateUserTransaction(gomock.Any(), model.Transaction{
					Amount:      100000,
					RecipientID: 2,
					Type:        model.TransactionTypeTransferOut,
					Description: "Traktir Makan",
					Password:    "Admin1234!",
				}).Return(convertToUUID("3d6e668f-ad02-40ff-8540-90c1528a7c88"), nil)

				return mock
			},
			wantResponse: generated.TransactionResponse{
				Header: generated.ResponseHeader{
					Success:  true,
					Messages: []string{successMsg},
				},
				Transaction: generated.Transaction{
					Id: stringPtr("3d6e668f-ad02-40ff-8540-90c1528a7c88"),
				},
			},
			wantHttpStatusCode: http.StatusOK,
		},
		{
			name: "fail-not-authorized-permission",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			ctxUserID: 123,
			requestBody: generated.Transaction{
				Amount:      floatPtr(100000),
				RecipientId: intPtr(2),
				Type:        transactionTypePtr(generated.TransferOut),
				Description: stringPtr("Traktir Makan"),
				Password:    stringPtr("Admin1234!"),
			},
			fnConvertCreateTransactionRequestToTransaction: func(int64, generated.Transaction) (model.Transaction, []string) {
				return model.Transaction{}, []string{}
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				return usecase.NewMockUsecaseInterface(controller)
			},
			wantResponse: generated.TransactionResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"not authorized: missing required permission"},
				},
				Transaction: generated.Transaction{},
			},
			wantHttpStatusCode: http.StatusForbidden,
		},
		{
			name: "fail-create-transaction",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionPerformTransaction,
			},
			ctxUserID: 123,
			requestBody: generated.Transaction{
				Amount:      floatPtr(100000),
				RecipientId: intPtr(2),
				Type:        transactionTypePtr(generated.TransferOut),
				Description: stringPtr("Traktir Makan"),
				Password:    stringPtr("Admin1234!"),
			},
			fnConvertCreateTransactionRequestToTransaction: func(int64, generated.Transaction) (model.Transaction, []string) {
				transaction := model.Transaction{
					Amount:      100000,
					RecipientID: 2,
					Type:        model.TransactionTypeTransferOut,
					Description: "Traktir Makan",
					Password:    "Admin1234!",
				}

				errorMsgs := []string{}

				return transaction, errorMsgs
			},
			mockUsecase: func(controller *gomock.Controller) *usecase.MockUsecaseInterface {
				mock := usecase.NewMockUsecaseInterface(controller)

				mock.EXPECT().CreateUserTransaction(gomock.Any(), model.Transaction{
					Amount:      100000,
					RecipientID: 2,
					Type:        model.TransactionTypeTransferOut,
					Description: "Traktir Makan",
					Password:    "Admin1234!",
				}).Return(uuid.Nil, errors.New("error-create-transaction"))

				return mock
			},
			wantResponse: generated.TransactionResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"error-create-transaction"},
				},
				Transaction: generated.Transaction{},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Usecase: test.mockUsecase(controller),
			}

			requestBodyJSON, _ := json.Marshal(test.requestBody)
			requestBody := []byte(requestBodyJSON)
			requestBodyBuffer := bytes.NewBuffer(requestBody)

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/v1/{userID}/transactions", requestBodyBuffer)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)
			if test.ctxUserID != 0 {
				ctx.Set(string(utils.JWTClaimUserID), test.ctxUserID)
			}
			if len(test.ctxPermissions) > 0 {
				ctx.Set(string(utils.JWTClaimPermissions), test.ctxPermissions)
			}

			fnConvertCreateTransactionRequestToTransaction = test.fnConvertCreateTransactionRequestToTransaction

			gotHttpStatusCode, gotResponse := handler.createUserTransaction(ctx, 123)

			if gotHttpStatusCode != test.wantHttpStatusCode {
				t.Errorf("handler.CreateUserTransaction() httpStatusCode = %v, wantHttpStatusCode %v", gotHttpStatusCode, test.wantHttpStatusCode)
			}

			if !reflect.DeepEqual(test.wantResponse, gotResponse) {
				t.Errorf("handler.CreateUserTransaction() response = %v, wantResponse %v", gotResponse, test.wantResponse)
			}
		})
	}
}
