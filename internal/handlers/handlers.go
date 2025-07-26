package handlers

import (
	"log"
	"net/http"
	"time"

	"2FA/internal/auth"
	"2FA/internal/models"
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
	user, err := h.authService.GetByUsername(username)
	if err != nil {
		log.Printf("error get user by username")
	}

	tokenPair, err := h.authService.GenerateTokenPair(user.ID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "verify.html", gin.H{
			"Error": "Failed to generate tokens",
		})
		return
	}

	setTokenCookies(c, tokenPair, h.authService.JwtExpiration, h.authService.JwtRefreshExpiration)

	c.HTML(http.StatusOK, "success.html", gin.H{
		"Username": username,
	})
}

func (h *AuthHandler) HandleSuccess(c *gin.Context) {
	c.HTML(http.StatusOK, "success.html", gin.H{
		"Username": "welcome",
	})
}

func setTokenCookies(c *gin.Context, pair *models.TokenPair, accessExpiry, refreshExpiry time.Duration) {
	c.SetCookie("access_token", pair.AccessToken, int(accessExpiry.Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", pair.RefreshToken, int(refreshExpiry.Seconds()), "/", "", false, true)
}

func JWTAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err == nil && accessToken != "" {
			claims, err := auth.ParseToken(accessToken, authService.JwtSecret)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Next()
				return
			}
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		newPair, err := authService.RefreshTokens(refreshToken)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		setTokenCookies(c, newPair, authService.JwtExpiration, authService.JwtRefreshExpiration)
		c.Set("user_id", newPair.UserID)
		c.Next()
	}
}
