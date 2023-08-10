package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
