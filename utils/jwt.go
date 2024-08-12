package utils

import (
	"github.com/WalletService/generated"
	"github.com/dgrijalva/jwt-go"
)

type JWTPermission string

const (
	JWTPermissionGetUser            JWTPermission = "get_profile"
	JWTPermissionPerformTransaction JWTPermission = "perform_transaction"
)

type JWTClaimKey string

const (
	JWTClaimUserID      JWTClaimKey = "user_id"
	JWTClaimPermissions JWTClaimKey = "permissions"
)

// CustomClaims represents the claims you want to include in your JWT.
type CustomClaims struct {
	UserID      int64           `json:"user_id"`
	Permissions []JWTPermission `json:"permissions"`
	ExpiresAt   int64           `json:"exp"`
	jwt.StandardClaims
}

type JWTResponse struct {
	Header generated.ResponseHeader `json:"header"`
}
