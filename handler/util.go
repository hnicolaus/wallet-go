package handler

import (
	"errors"
	"strings"
	"unicode"

	"github.com/WalletService/generated"
	"github.com/WalletService/model"
	"github.com/WalletService/utils"
	"github.com/labstack/echo/v4"
)

var (
	// Define function wrappers so we can inject dummy function in UT
	fnValidatePhoneNumber func(*string) (string, []string) = validatePhoneNumber
	fnValidatePassword    func(*string) (string, []string) = validatePassword
	fnValidateFullName    func(*string) (string, []string) = validateFullName
)

const (
	successMsg                   = "request successful"
	duplicatePhoneNumberErrorMsg = "phone number is already registered to an existing user"
)

func authorize(ctx echo.Context, requiredPermission utils.JWTPermission) (userID int64, err error) {
	permissions, _ := ctx.Get(string(utils.JWTClaimPermissions)).([]utils.JWTPermission)

	hasRole := false
	for _, permission := range permissions {
		if permission == requiredPermission {
			hasRole = true
		}
	}

	if !hasRole {
		return userID, errors.New("not authorized: missing required permission")
	}

	userID, ok := ctx.Get(string(utils.JWTClaimUserID)).(int64)
	if !ok {
		return userID, errors.New("missing user_id")
	}

	return userID, nil
}

func validatePhoneNumber(input *string) (validPhoneNumber string, errorList []string) {
	phoneNumber := ""

	if input != nil {
		phoneNumber = strings.TrimSpace(*input)
	}

	// Verify "+62" prefix
	if !strings.HasPrefix(phoneNumber, "+62") {
		errorList = append(errorList, "phone_number should start with +62")
	}

	// Check the length of the phone number
	if len(phoneNumber) < 10 || len(phoneNumber) > 13 {
		errorList = append(errorList, "phone_number should be 10 to 13 digits")
	}

	// Check if all remaining characters are digits
	for i := 3; i < len(phoneNumber); i++ {
		c := phoneNumber[i]
		if !unicode.IsDigit(rune(c)) {
			errorList = append(errorList, "phone_number should only contain numbers")
			break
		}
	}

	if len(errorList) == 0 {
		validPhoneNumber = phoneNumber
	}

	return validPhoneNumber, errorList
}

func validateFullName(input *string) (validFullName string, errorList []string) {
	fullName := ""

	if input != nil {
		fullName = strings.TrimSpace(*input)
	}

	if len(fullName) < 3 || len(fullName) > 60 {
		errorList = append(errorList, "full_name should be 3 to 60 characters")
	}

	if len(errorList) == 0 {
		validFullName = fullName
	}

	return validFullName, errorList
}

func validatePassword(input *string) (validPassword string, errorList []string) {
	password := ""

	if input != nil {
		password = *input
	}

	if len(password) < 6 || len(password) > 64 {
		errorList = append(errorList, "password should be 6 to 64 characters")
	}

	containsCapital, containsNumber, containsSpecialAlphaNumeric := false, false, false
	for _, c := range password {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			containsSpecialAlphaNumeric = true
		}
		if unicode.IsUpper(c) {
			containsCapital = true
		}
		if unicode.IsNumber(c) {
			containsNumber = true
		}

		if containsSpecialAlphaNumeric && containsCapital && containsNumber {
			break
		}
	}

	if !containsCapital {
		errorList = append(errorList, "password should contain a capital letter")
	}
	if !containsNumber {
		errorList = append(errorList, "password should contain a number")
	}
	if !containsSpecialAlphaNumeric {
		errorList = append(errorList, "password should contain a special alphanumeric character")
	}

	if len(errorList) == 0 {
		validPassword = password
	}

	return validPassword, errorList
}

func convertRegisterUserRequestToUser(request generated.User) (user model.User, errorMsgs []string) {
	validPhoneNumber, phoneNumberErrorMsgs := fnValidatePhoneNumber(request.PhoneNumber)
	validFullName, fullNameErrorMsgs := fnValidateFullName(request.FullName)
	validPassword, passwordErrorMsgs := fnValidatePassword(request.Password)

	errorList := append(phoneNumberErrorMsgs, fullNameErrorMsgs...)
	errorList = append(errorList, passwordErrorMsgs...)

	if len(errorList) > 0 {
		return model.User{}, errorList
	}

	return model.User{
		FullName:    validFullName,
		PhoneNumber: validPhoneNumber,
		Password:    validPassword,
	}, nil
}

func convertCreateTransactionRequestToTransaction(userID int64, request generated.Transaction) (transaction model.Transaction, errorMsgs []string) {
	if request.Amount == nil || *request.Amount <= 0 {
		errorMsgs = append(errorMsgs, "amount should be > 0")
		return model.Transaction{}, errorMsgs
	}

	if request.Type == nil {
		errorMsgs = append(errorMsgs, "type is required")
		return model.Transaction{}, errorMsgs
	}

	transaction = model.Transaction{
		UserID: userID,
		Amount: *request.Amount,
		Type:   model.TransactionType(*request.Type),
	}

	if request.RecipientId != nil {
		transaction.RecipientID = *request.RecipientId
	}

	if request.Description != nil {
		transaction.Description = *request.Description
	}

	if request.Password != nil {
		transaction.Password = *request.Password
	}

	return transaction, nil
}
