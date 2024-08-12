package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/WalletService/generated"
	"github.com/WalletService/model"
	"github.com/WalletService/utils"
	"github.com/labstack/echo/v4"
)

func Test_validatePhoneNumber(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name  string
		input *string

		wantValidPhoneNumber string
		wantErrorList        []string
		wantError            bool
	}{
		{
			name:                 "success",
			input:                stringPtr("+628123456789"),
			wantValidPhoneNumber: "+628123456789",
			wantErrorList:        nil,
		},
		{
			name:                 "nil",
			input:                nil,
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should start with +62",
				"phone_number should be 10 to 13 digits",
			},
		},
		{
			name:                 "empty",
			input:                stringPtr(""),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should start with +62",
				"phone_number should be 10 to 13 digits",
			},
		},
		{
			name:                 "fail-rule-1-length-min",
			input:                stringPtr("+62812345  "),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should be 10 to 13 digits",
			},
		},
		{
			name:                 "fail-rule-1-length-max",
			input:                stringPtr("  +6281234567890"),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should be 10 to 13 digits",
			},
		},
		{
			name:                 "fail-rule-1-non-numbers",
			input:                stringPtr("+628123456a"),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should only contain numbers",
			},
		},
		{
			name:                 "fail-rule-2",
			input:                stringPtr("08123456789"),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should start with +62",
			},
		},
		{
			name:                 "fail-all-rules",
			input:                stringPtr("0345abc"),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should start with +62",
				"phone_number should be 10 to 13 digits",
				"phone_number should only contain numbers",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotValidPhoneNumber, gotErrorList := validatePhoneNumber(test.input)
			if !reflect.DeepEqual(gotValidPhoneNumber, test.wantValidPhoneNumber) {
				t.Errorf("util.validatePhoneNumber() gotValidPhoneNumber = %v, wantValidPhoneNumber %v", gotValidPhoneNumber, test.wantValidPhoneNumber)
			}

			if !reflect.DeepEqual(gotErrorList, test.wantErrorList) {
				t.Errorf("util.validatePhoneNumber() gotErrorList = %v, wantErrorList %v", gotErrorList, test.wantErrorList)
			}
		})
	}
}

func Test_validateFullName(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name  string
		input *string

		wantValidFullName string
		wantErrorList     []string
	}{
		{
			name:              "success",
			input:             stringPtr("abcd"),
			wantValidFullName: "abcd",
			wantErrorList:     nil,
		},
		{
			name:              "nil",
			input:             nil,
			wantValidFullName: "",
			wantErrorList: []string{
				"full_name should be 3 to 60 characters",
			},
		},
		{
			name:              "fail-length-min",
			input:             stringPtr(" ab"),
			wantValidFullName: "",
			wantErrorList: []string{
				"full_name should be 3 to 60 characters",
			},
		},
		{
			name:              "fail-length-max",
			input:             stringPtr("    S2EeAKi6fze0JVsVbo6OR9uxmzdy89Kiy59z4Wzi2jTdomVUSUIh8G1GmHpJF "),
			wantValidFullName: "",
			wantErrorList: []string{
				"full_name should be 3 to 60 characters",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotValidFullName, gotErrorList := validateFullName(test.input)
			if !reflect.DeepEqual(gotValidFullName, test.wantValidFullName) {
				t.Errorf("util.validateFullName() gotValidFullName = %v, wantValidFullName %v", gotValidFullName, test.wantValidFullName)
			}

			if !reflect.DeepEqual(gotErrorList, test.wantErrorList) {
				t.Errorf("util.validateFullName() gotErrorList = %v, wantErrorList %v", gotErrorList, test.wantErrorList)
			}
		})
	}
}

func Test_validatePassword(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name  string
		input *string

		wantValidPassword string
		wantErrorList     []string
	}{
		{
			name:              "success",
			input:             stringPtr("A1.123"),
			wantValidPassword: "A1.123",
			wantErrorList:     nil,
		},
		{
			name:              "fail-all-rules",
			input:             stringPtr(""),
			wantValidPassword: "",
			wantErrorList: []string{
				"password should be 6 to 64 characters",
				"password should contain a capital letter",
				"password should contain a number",
				"password should contain a special alphanumeric character",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotValidPassword, gotErrorList := validatePassword(test.input)
			if !reflect.DeepEqual(gotValidPassword, test.wantValidPassword) {
				t.Errorf("util.validatePassword() gotValidPassword = %v, wantValidPassword %v", gotValidPassword, test.wantValidPassword)
			}

			if !reflect.DeepEqual(gotErrorList, test.wantErrorList) {
				t.Errorf("util.validatePassword() gotErrorList = %v, wantErrorList %v", gotErrorList, test.wantErrorList)
			}
		})
	}
}

func Test_convertRegisterUserRequestToUser(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name                  string
		input                 generated.User
		fnValidatePhoneNumber func(*string) (string, []string)
		fnValidatePassword    func(*string) (string, []string)
		fnValidateFullName    func(*string) (string, []string)

		wantUser      model.User
		wantErrorMsgs []string
	}{
		{
			name: "success",
			input: generated.User{
				FullName:    stringPtr("User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnValidatePhoneNumber: func(*string) (string, []string) {
				return "+628123456789", []string{}
			},
			fnValidateFullName: func(*string) (string, []string) {
				return "User", []string{}
			},
			fnValidatePassword: func(*string) (string, []string) {
				return "P455w0rd!.", []string{}
			},
			wantUser: model.User{
				FullName:    "User",
				PhoneNumber: "+628123456789",
				Password:    "P455w0rd!.",
			},
			wantErrorMsgs: nil,
		},
		{
			name: "fail-all-validations",
			input: generated.User{
				FullName:    stringPtr("User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnValidatePhoneNumber: func(*string) (string, []string) {
				return "", []string{"invalid-phone-number-1", "invalid-phone-number-2", "invalid-phone-number-3"}
			},
			fnValidateFullName: func(*string) (string, []string) {
				return "", []string{"invalid-full-name-1", "invalid-full-name-2"}
			},
			fnValidatePassword: func(*string) (string, []string) {
				return "", []string{"invalid-password"}
			},
			wantUser: model.User{},
			wantErrorMsgs: []string{
				"invalid-phone-number-1", "invalid-phone-number-2", "invalid-phone-number-3",
				"invalid-full-name-1", "invalid-full-name-2",
				"invalid-password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fnValidateFullName = test.fnValidateFullName
			fnValidatePassword = test.fnValidatePassword
			fnValidatePhoneNumber = test.fnValidatePhoneNumber

			gotUser, gotErrorMsgs := convertRegisterUserRequestToUser(test.input)
			if !reflect.DeepEqual(gotUser, test.wantUser) {
				t.Errorf("util.convertRegisterUserRequestToUser() gotUser = %v, wantUser %v", gotUser, test.wantUser)
			}

			if !reflect.DeepEqual(gotErrorMsgs, test.wantErrorMsgs) {
				t.Errorf("util.convertRegisterUserRequestToUser() gotErrorMsgs = %v, wantErrorMsgs %v", gotErrorMsgs, test.wantErrorMsgs)
			}
		})
	}
}

func Test_convertCreateTransactionRequestToTransaction(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name        string
		inputUserID int64
		input       generated.Transaction

		wantTransaction model.Transaction
		wantErrorMsgs   []string
	}{
		{
			name:        "success",
			inputUserID: 123,
			input: generated.Transaction{
				Amount:      floatPtr(100000),
				RecipientId: intPtr(2),
				Type:        transactionTypePtr(generated.TransferOut),
				Description: stringPtr("Traktir Makan"),
				Password:    stringPtr("Admin1234!"),
			},
			wantTransaction: model.Transaction{
				UserID:      123,
				Amount:      100000,
				RecipientID: 2,
				Type:        model.TransactionTypeTransferOut,
				Description: "Traktir Makan",
				Password:    "Admin1234!",
			},
			wantErrorMsgs: nil,
		},
		{
			name:        "success-nil-inputs",
			inputUserID: 123,
			input: generated.Transaction{
				Amount:      floatPtr(100000),
				Type:        transactionTypePtr(generated.TopUp),
				RecipientId: nil,
				Description: nil,
				Password:    nil,
			},
			wantTransaction: model.Transaction{
				UserID:      123,
				Amount:      100000,
				Type:        model.TransactionTypeTopUp,
				RecipientID: 0,
				Description: "",
				Password:    "",
			},
			wantErrorMsgs: nil,
		},
		{
			name:        "invalid-amount-zero",
			inputUserID: 123,
			input: generated.Transaction{
				Amount:      floatPtr(0),
				RecipientId: intPtr(2),
				Type:        transactionTypePtr(generated.TransferOut),
				Description: stringPtr("Traktir Makan"),
				Password:    stringPtr("Admin1234!"),
			},
			wantTransaction: model.Transaction{},
			wantErrorMsgs:   []string{"amount should be > 0"},
		},
		{
			name:        "invalid-amount-nil",
			inputUserID: 123,
			input: generated.Transaction{
				RecipientId: intPtr(2),
				Type:        transactionTypePtr(generated.TransferOut),
				Description: stringPtr("Traktir Makan"),
				Password:    stringPtr("Admin1234!"),
			},
			wantTransaction: model.Transaction{},
			wantErrorMsgs:   []string{"amount should be > 0"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotTransaction, gotErrorMsgs := convertCreateTransactionRequestToTransaction(test.inputUserID, test.input)
			if !reflect.DeepEqual(gotTransaction, test.wantTransaction) {
				t.Errorf("util.convertCreateTransactionRequestToTransaction() gotTransaction = %v, wanwantTransactiontUser %v", gotTransaction, test.wantTransaction)
			}

			if !reflect.DeepEqual(gotErrorMsgs, test.wantErrorMsgs) {
				t.Errorf("util.convertCreateTransactionRequestToTransaction() gotErrorMsgs = %v, wantErrorMsgs %v", gotErrorMsgs, test.wantErrorMsgs)
			}
		})
	}
}

func Test_authorize(t *testing.T) {
	tests := []struct {
		name               string
		ctxUserID          int64
		ctxPermissions     []utils.JWTPermission
		requiredPermission utils.JWTPermission

		wantUserID int64
		wantErr    error
	}{
		{
			name:      "success",
			ctxUserID: 123,
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
				utils.JWTPermissionPerformTransaction,
			},
			requiredPermission: utils.JWTPermissionGetUser,
			wantUserID:         123,
			wantErr:            nil,
		},
		{
			name:      "fail-not-authorized-permission",
			ctxUserID: 123,
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionPerformTransaction,
			},
			requiredPermission: utils.JWTPermissionGetUser,
			wantUserID:         0,
			wantErr:            errors.New("not authorized: missing required permission"),
		},
		{
			name:      "fail-not-authorized-user-id",
			ctxUserID: 0,
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			requiredPermission: utils.JWTPermissionGetUser,
			wantUserID:         0,
			wantErr:            errors.New("missing user_id"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/v1/unit-test", nil)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)
			if test.ctxUserID != 0 {
				ctx.Set(string(utils.JWTClaimUserID), test.ctxUserID)
			}
			if len(test.ctxPermissions) > 0 {
				ctx.Set(string(utils.JWTClaimPermissions), test.ctxPermissions)
			}

			gotUserID, gotErr := authorize(ctx, test.requiredPermission)
			if gotUserID != test.wantUserID {
				t.Errorf("util.authorize() gotUserID = %v, wantUserID %v", gotUserID, test.wantUserID)
			}

			if !reflect.DeepEqual(gotErr, test.wantErr) {
				t.Errorf("util.authorize() gotErr = %v, wantErr %v", gotErr, test.wantErr)
			}
		})
	}
}
