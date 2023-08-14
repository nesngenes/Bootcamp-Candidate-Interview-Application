package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	ApplicationName     string
	JwtSigntureKey      string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Username string `json: "username"`
}

type authHandler struct {
	Authorizationhandler string
}

func AuthMiddlerware() gin.HandlerFunc {
	return func(a *gin.Context) {
		if a.Request.URL.Path == "/api/v1/users" { //bila autentifikasi berhasil
			a.Next()
			fmt.Println("masuk halaman pengguna")
		} else { //bila autentifikasi gagal
			var unAuth authHandler
			if err := a.ShouldBindHeader(&unAuth); err != nil {
				a.JSON(http.StatusUnauthorized, gin.H{"pesan": "unauthorization"})
				a.Abort()
				return
			}

			if unAuth.Authorizationhandler != "token" {
				a.JSON(http.StatusUnauthorized, gin.H{"pesan": "unauthorization"})
				a.Abort()
				return
			}
		}
	}
}
