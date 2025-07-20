package handlers

import (
	"net/http"

	"2FA/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) InitAuth(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "auth initiated"})
}

func (h *AuthHandler) VerifyAuth(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "verified"})
}

func (h *AuthHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (h *AuthHandler) HandleLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(password) < 8 {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": "Incorrect password",
		})
		return
	}
	//check password later
	c.HTML(http.StatusOK, "verify.html", gin.H{
		"Username": username,
	})
}

func (h *AuthHandler) HandleVerify(c *gin.Context) {
	username := c.PostForm("username")
	code := c.PostForm("code")
	//send code later
	//check code later
	if code == "1" {
		c.String(http.StatusOK, "succes join: %s", username)
		return
	}

	c.HTML(http.StatusOK, "verify.html", gin.H{
		"Username": username,
		"Error":    "Неверный код подтверждения",
	})
}
