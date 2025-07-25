package server

import (
	"2FA/internal/auth"
	"2FA/internal/handlers"

	"github.com/gin-gonic/gin"
)

func NewServer(authHandler *handlers.AuthHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")

	router.GET("/login", authHandler.ShowLoginPage)
	router.POST("/login", authHandler.HandleLogin)

	private := router.Group("/")
	private.Use(auth.JWTAuth(jwtSecret))
	{
		private.POST("/verify", authHandler.HandleVerify)
	}

	return router
}
