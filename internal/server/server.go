package server

import (
	"2FA/internal/handlers"
	"2FA/internal/services"

	"github.com/gin-gonic/gin"
)

func NewServer(authHandler *handlers.AuthHandler, authService *services.AuthService, jwtSecret string, jwtRefresh string) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")

	router.GET("/login", authHandler.ShowLoginPage)
	router.POST("/login", authHandler.HandleLogin)
	router.POST("/verify", authHandler.HandleVerify)

	private := router.Group("/")
	private.Use(handlers.JWTAuth(authService))
	{
		private.GET("/success", authHandler.HandleSuccess)
	}

	return router
}
