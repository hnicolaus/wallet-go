package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/WalletService/generated"
	"github.com/WalletService/handler"
	"github.com/WalletService/repository"
	"github.com/WalletService/usecase"
	"github.com/WalletService/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	jwtExpiryDuration = time.Minute * 5 // Token expires in 5 minutes by default
)

func init() {
	// Ensure time is set to Asia/Jakarta time.
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	time.Local = location
}

func main() {
	e := echo.New()
	e.Pre(AuthenticationMiddleware) // Register pre-handler middleware
	e.Use(AuthenticatedMiddleware)  // Register post-handler middleware

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	var repo repository.RepositoryInterface = repository.NewRepository(os.Getenv("DATABASE_URL"))
	var usecase usecase.UsecaseInterface = usecase.NewUsecase(repo)
	return handler.NewServer(usecase)
}

// AuthenticationMiddleware validates incoming JWT (RS256 algorithm) using public key.
// See command in Dockerfile: openssl rsa -in /tmp/rsa -pubout -out /tmp/rsa.pub
func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		whitelistedEndpoints := []string{
			"GET - /v1/user",
			`POST - /v1/user/\d+/transactions`,
		}

		if isEndpointWhitelisted(ctx, whitelistedEndpoints) {
			claims, err := func() (*utils.CustomClaims, error) {
				authHeader := ctx.Request().Header.Get("Authorization")
				if authHeader == "" {
					return nil, errors.New("missing authorization header")
				}

				if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
					return nil, errors.New("token format is invalid")
				}

				token := authHeader[7:]

				publicKeyData, err := os.ReadFile("../rsa.pub")
				if err != nil {
					return nil, err
				}

				publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
				if err != nil {
					return nil, err
				}

				// Parse & validate token using publicKey
				tok, err := jwt.ParseWithClaims(token, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					// Use publicKey for token verification
					return publicKey, nil
				})
				if err != nil {
					return nil, err
				}

				// Check if token is valid (has not expired)
				claims, ok := tok.Claims.(*utils.CustomClaims)
				if ok && tok.Valid {
					if time.Now().Unix() >= claims.ExpiresAt {
						return nil, errors.New("JWT has expired")
					}
				}

				return claims, nil
			}()
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, utils.JWTResponse{
					Header: generated.ResponseHeader{
						Messages: []string{err.Error()},
						Success:  false,
					},
				})
			}

			// Set custom claims to context so handler can use the values, i.e. authorization
			ctx.Set(string(utils.JWTClaimUserID), claims.UserID)
			ctx.Set(string(utils.JWTClaimPermissions), claims.Permissions)
		}

		return next(ctx)
	}
}

// AuthenticatedMiddleware generates JWT (RS256 algorithm) using private key
// The JWT will contain userID and its claims
// The JWT will be included in the `Authentication` response header only if handler is returning status OK
// See command in Dockerfile: openssl genrsa -out /tmp/rsa 4096
func AuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		whitelistedEndpoints := []string{
			"POST - /v1/user/login",
		}

		if isEndpointWhitelisted(ctx, whitelistedEndpoints) {
			// use `Before` so this middleware can write token to response header right before handler writes to response body
			ctx.Response().Before(func() {
				// Only add JWT in `Authorization` response header if handler returns OK
				if ctx.Response().Status == http.StatusOK {
					privateKeyData, err := os.ReadFile("../rsa")
					if err != nil {
						return
					}

					privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
					if err != nil {
						return
					}

					// Authorize JWT
					permissions, ok := ctx.Get(string(utils.JWTClaimPermissions)).([]utils.JWTPermission)
					if !ok {
						return
					}

					userID, ok := ctx.Get(string(utils.JWTClaimUserID)).(int64)
					if !ok {
						return
					}

					claims := utils.CustomClaims{
						UserID:      userID,
						Permissions: permissions,
						ExpiresAt:   time.Now().Add(jwtExpiryDuration).Unix(),
					}

					token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
					jwtToken, err := token.SignedString(privateKey) // sign the token with the RSA private key
					if err != nil {
						return
					}

					bearerToken := fmt.Sprintf("Bearer %s", jwtToken)

					ctx.Response().Header().Set(echo.HeaderAuthorization, bearerToken)
				}
			})
		}

		return next(ctx)
	}
}

func isEndpointWhitelisted(ctx echo.Context, whitelist []string) bool {
	endpoint := fmt.Sprintf("%s - %s", ctx.Request().Method, ctx.Request().URL.Path)
	for _, pattern := range whitelist {
		matched, _ := regexp.MatchString(pattern, endpoint)
		if matched {
			return true
		}
	}
	return false
}
