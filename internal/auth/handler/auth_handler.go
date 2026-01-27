package handler

import (
	"errors"
	response "github.com/hardikm9850/GoChat/internal/http/validation"
	"github.com/hardikm9850/authkit/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/auth"
	"github.com/hardikm9850/GoChat/internal/auth/service"
)

type AuthHandler struct {
	service service.AuthService
}

func New(auth service.AuthService) *AuthHandler {
	return &AuthHandler{service: auth}
}

type registerRequest struct {
	Phone       string `json:"phone" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
}

type loginRequest struct {
	Phone       string `json:"phone" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

// refreshRequest represents refresh token request
// swagger:model RefreshRequest
type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// refreshResponse represents token response
// swagger:model RefreshResponse
type refreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @BasePath /api/v1

// @Summary Register a new user
// @Description Registers a new user with phone, password, and name
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body registerRequest true "Register Request"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 409 {object} map[string]string "User already exists"
// @Failure 500 {object} map[string]string "Registration failed"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	log.Println("Http 👉 /auth/register hit")
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.Register(req.CountryCode, req.Phone, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Received error %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusOK, authResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		UserID:       tokens.UserID,
	})
}

// @Summary Login user
// @Description Login user with phone and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Login Request"
// @Success 200 {object} authResponse "Access token returned"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Login failed"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError(err))
		return
	}

	tokens, err := h.service.Login(req.Phone, req.Password, req.CountryCode)

	if err != nil {
		log.Printf("error is %s\n", err)
		if errors.Is(err, auth.ErrUserDoesNotExists) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "account does not exist"})
			return
		}
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		log.Printf("error %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, authResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Refresh godoc
//
// @Summary      Refresh access token
// @Description  Rotates refresh token and issues a new access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body refreshRequest true "Refresh token payload"
// @Success      200 {object} refreshResponse
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      401 {object} map[string]string "Invalid / expired / reused refresh token"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /auth/v1/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tokens, err := h.service.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		mapRefreshError(c, err)
		return
	}

	c.JSON(http.StatusOK, refreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Logout godoc
//
// @Summary Logout user
// @Description Invalidates the refresh token for the logged-in user
// @Tags Auth
// @Security BearerAuth
// @Success 200 {object} map[string]string "logged out successfully"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "failed to logout"
// @Router /auth/v1/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.Logout(userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

}

func mapRefreshError(c *gin.Context, err error) {
	switch err {
	case auth.ErrInvalidToken,
		auth.ErrExpiredToken,
		auth.ErrTokenReuseDetected:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
