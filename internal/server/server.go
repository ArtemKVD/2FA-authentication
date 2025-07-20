package server

import (
	"2FA/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServer(authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("web/templates/*")

	api := router.Group("/api")
	{
		api.POST("/auth/init", authHandler.InitAuth)
		api.POST("/auth/verify", authHandler.VerifyAuth)
	}

	router.GET("/login", authHandler.ShowLoginPage)
	router.POST("/login", authHandler.HandleLogin)
	router.POST("/verify", authHandler.HandleVerify)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	return router
}
