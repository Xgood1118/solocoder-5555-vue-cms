package handlers

import (
	"net/http"

	"cms/middleware"
	"cms/models"
	"cms/store"
	"cms/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if !utils.ValidatePassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters with letters and numbers"})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}

	existingUser, _ := store.GetUserByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	existingEmail, _ := store.GetUserByEmail(req.Email)
	if existingEmail != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleAuthor,
		Avatar:       "",
	}

	if err := store.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	CreateAuditLog(c, "register", "user", user.ID, "user registered")

	user.PasswordHash = ""
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var user *models.User
	var err error

	if req.Username != "" {
		user, err = store.GetUserByUsername(req.Username)
	} else if req.Email != "" {
		user, err = store.GetUserByEmail(req.Email)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email required"})
		return
	}

	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	tokenPair, err := utils.GenerateTokenPair(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	CreateAuditLog(c, "login", "user", user.ID, "user logged in")

	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"access_token": tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":   tokenPair.ExpiresIn,
	})
}

func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tokenPair, err := utils.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
	})
}

func Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	utils.RevokeRefreshToken(req.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	user, err := store.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"user_id": userID,
		"username": username,
		"role": role,
	})
}

func GetCSRFToken(c *gin.Context) {
	middleware.CSRFTokenHandler(c)
}

func GitHubLogin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "GitHub OAuth requires client ID and secret configuration",
		"message": "Set GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET environment variables to enable GitHub login",
	})
}

func GitHubCallback(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "GitHub OAuth callback not configured",
	})
}

func CreateAuditLog(c *gin.Context, action, resource, resourceID, detail string) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	userIDStr := ""
	usernameStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}
	if username != nil {
		usernameStr = username.(string)
	}

	log := &models.AuditLog{
		UserID:     userIDStr,
		Username:   usernameStr,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IP:         c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		Detail:     detail,
	}

	store.CreateAuditLog(log)
}
