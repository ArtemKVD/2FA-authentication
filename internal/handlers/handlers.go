package handlers

import (
	"log"
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

	valid, err := h.authService.DataVerify(username, password)
	if err != nil {
		log.Printf("Login error: %v", err)
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": "Invalid username or password",
		})
		return
	}

	if !valid {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": "Invalid username or password",
		})
		return
	}

	err = h.authService.SendAuthCode(username)
	if err != nil {
		log.Printf("Failed to send auth code: %v", err)
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": "Failed to send verification code",
		})
		return
	}

	c.HTML(http.StatusOK, "verify.html", gin.H{
		"Username": username,
	})
}

func (h *AuthHandler) HandleVerify(c *gin.Context) {
	username := c.PostForm("username")
	code := c.PostForm("code")

	valid, err := h.authService.VerifyCode(username, code)
	log.Printf(code)
	if err != nil {
		log.Printf("Ошибка при проверке кода: %v", err)
		c.HTML(http.StatusOK, "verify.html", gin.H{
			"Username": username,
			"Error":    "error check code",
		})
		return
	}

	if !valid {
		c.HTML(http.StatusOK, "verify.html", gin.H{
			"Username": username,
			"Error":    "wrong code",
		})
		return
	}

	c.HTML(http.StatusOK, "success.html", gin.H{
		"Username": username,
	})
}
