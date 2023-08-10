package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	ApplicationName     string
	JwtSigntureKey      string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type TokenMyClaims struct {
	jwt.RegisteredClaims
	Username string `json: "username"`
}
