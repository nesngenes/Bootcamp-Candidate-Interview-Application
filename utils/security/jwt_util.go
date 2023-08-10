package security

import (
	"fmt"
	"interview_bootcamp/config"
	"interview_bootcamp/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user model.Users) (string, error) {
	AccessTokenLifeTime := time.Duration(1) * time.Minute
	tokenConfig := TokenConfig{
		ApplicationName:     "enigma-interview-bootcamp",
		JwtSigntureKey:      "test-kunci-masuk",
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: AccessTokenLifeTime,
	}

	now := time.Now().UTC()
	end := now.Add(tokenConfig.AccessTokenLifeTime)

	tokenMyClaims := TokenMyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: tokenConfig.ApplicationName,
		},
		Username: user.UserName,
	}

	tokenMyClaims.IssuedAt = jwt.NewNumericDate(now)
	tokenMyClaims.ExpiresAt = jwt.NewNumericDate(end)

	token := jwt.NewWithClaims(tokenConfig.JwtSigningMethod, tokenMyClaims)
	sString, err := token.SignedString(tokenConfig.JwtSigntureKey)
	if err != nil {
		return "gagal", err
	}

	return sString, nil
}

func VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	cfg, _ := config.NewConfig()

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || method != cfg.JwtSigningMethod {
			return nil, fmt.Errorf("token tidak valid")
		}
		return cfg.JwtSigntureKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok || !token.Valid || claims["iss"] != cfg.ApplicationName {
		return nil, fmt.Errorf("token tidak valid")
	}

	return claims, nil
}
