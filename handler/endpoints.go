package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/WalletService/generated"
	"github.com/WalletService/model"
	"github.com/WalletService/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	// Define function wrappers so we can inject dummy function in UT
	fnConvertRegisterUserRequestToUser             func(generated.User) (model.User, []string)                      = convertRegisterUserRequestToUser
	fnConvertCreateTransactionRequestToTransaction func(int64, generated.Transaction) (model.Transaction, []string) = convertCreateTransactionRequestToTransaction
)

// RegisterUser creates a new User with a unique phoneNumber and valid password format.
func (s *Server) RegisterUser(ctx echo.Context) error {
	return ctx.JSON(s.registerUser(ctx))
}
func (s *Server) registerUser(ctx echo.Context) (int, generated.RegisterUserResponse) {
	var (
		context = context.Background()

		response = generated.RegisterUserResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	request := generated.User{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusBadRequest, response
	}

	user, errorList := fnConvertRegisterUserRequestToUser(request)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return http.StatusBadRequest, response
	}

	userID, err := s.Usecase.RegisterUser(context, user)
	if err != nil {
		if utils.IsUniqueConstraintViolation(err) {
			response.Header.Messages = []string{duplicatePhoneNumberErrorMsg}
			return http.StatusConflict, response
		}

		response.Header.Messages = []string{err.Error()}
		return http.StatusInternalServerError, response
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.User.Id = &userID
	return http.StatusCreated, response
}

// UserLogin will return JWT to be used by User in subsequent requests requiring elevated permissions, i.e. performing Transaction.
// NOTE: Check Authenticated cmd/main.go that adds JWT token to response header after successful UserLogin attempt.
func (s *Server) UserLogin(ctx echo.Context) error {
	return ctx.JSON(s.userLogin(ctx))
}
func (s *Server) userLogin(ctx echo.Context) (int, generated.UserLoginResponse) {
	var (
		context = context.Background()

		response = generated.UserLoginResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	request := generated.User{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusBadRequest, response
	}

	// Get user's phone number from request body
	validPhoneNumber, errorList := validatePhoneNumber(request.PhoneNumber)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return http.StatusBadRequest, response
	}

	// Validate password format is valid
	validPassword, errorList := validatePassword(request.Password)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return http.StatusBadRequest, response
	}

	// Logs in as user
	userID, err := s.Usecase.UserLogin(context, validPhoneNumber, validPassword)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusBadRequest, response
	}

	// Set data to Echo context so we can rely on AuthenticatedMiddleware to generate and return JWT in the Authorization header
	ctx.Set(string(utils.JWTClaimUserID), userID)
	ctx.Set(string(utils.JWTClaimPermissions), []utils.JWTPermission{utils.JWTPermissionGetUser, utils.JWTPermissionPerformTransaction})

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.User.Id = &userID
	return http.StatusOK, response
}

// GetUser retrieves User data, i.e. wallet balance.
// NOTE: Check AuthenticationMiddleware cmd/main.go that authenticates the JWT token
func (s *Server) GetUser(ctx echo.Context) error {
	return ctx.JSON(s.getUser(ctx))
}
func (s *Server) getUser(ctx echo.Context) (int, generated.GetUserResponse) {
	var (
		context = context.Background()

		response = generated.GetUserResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	// Authorize and get userID of the requester
	userID, err := authorize(ctx, utils.JWTPermissionGetUser)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusForbidden, response
	}

	// Get data for the userID
	user, err := s.Usecase.GetUser(context, userID)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusInternalServerError, response
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.User = generated.User{
		Id:          &user.ID,
		FullName:    &user.FullName,
		PhoneNumber: &user.PhoneNumber,
		Balance:     &user.Balance,
	}

	return http.StatusOK, response
}

// CreateUserTransaction allows authenticated and authorized User to perform Transaction, i.e. TransferOut, TopUp, etc.
// NOTE: Check AuthenticationMiddleware cmd/main.go that authenticates the JWT token
func (s *Server) CreateUserTransaction(ctx echo.Context, pathUserID int) error {
	return ctx.JSON(s.createUserTransaction(ctx, int64(pathUserID)))
}
func (s *Server) createUserTransaction(ctx echo.Context, pathUserID int64) (int, generated.TransactionResponse) {
	var (
		context = context.Background()

		response = generated.TransactionResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	// Authenticate and get userID of the requester
	userID, err := authorize(ctx, utils.JWTPermissionPerformTransaction)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusForbidden, response
	} else if pathUserID != int64(userID) {
		response.Header.Messages = []string{"JWT userID mismatched with request userID"}
		return http.StatusForbidden, response
	}

	request := generated.Transaction{}

	err = json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusBadRequest, response
	}

	transaction, errorList := fnConvertCreateTransactionRequestToTransaction(userID, request)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return http.StatusBadRequest, response
	}

	newTransactionID, err := s.Usecase.CreateUserTransaction(context, transaction)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusInternalServerError, response
	}

	if newTransactionID != uuid.Nil {
		newTransactionIDString := newTransactionID.String()
		response.Transaction = generated.Transaction{
			Id: &newTransactionIDString,
		}
	}
	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	return http.StatusOK, response
}
