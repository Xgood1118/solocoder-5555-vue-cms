package utils

import (
	"errors"
	"time"

	"cms/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

var refreshTokenStore = make(map[string]string)

func GenerateTokenPair(userID, username, role string) (*TokenPair, error) {
	accessToken, err := generateAccessToken(userID, username, role)
	if err != nil {
		return nil, err
	}

	refreshTokenID := uuid.New().String()
	refreshToken, err := generateRefreshToken(userID, username, role, refreshTokenID)
	if err != nil {
		return nil, err
	}

	refreshTokenStore[refreshTokenID] = userID

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    2 * 60 * 60,
	}, nil
}

func generateAccessToken(userID, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewTimeFloat(float64(time.Now().Add(2 * time.Hour).Unix())),
			IssuedAt:  jwt.NewTimeFloat(float64(time.Now().Unix())),
			Issuer:    "vue-cms",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func generateRefreshToken(userID, username, role string, tokenID string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewTimeFloat(float64(time.Now().Add(7 * 24 * time.Hour).Unix())),
			IssuedAt:  jwt.NewTimeFloat(float64(time.Now().Unix())),
			Issuer:    "vue-cms",
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	tokenID := claims.ID
	if _, exists := refreshTokenStore[tokenID]; !exists {
		return nil, errors.New("refresh token revoked")
	}

	delete(refreshTokenStore, tokenID)

	return GenerateTokenPair(claims.UserID, claims.Username, claims.Role)
}

func RevokeRefreshToken(refreshToken string) error {
	claims, err := ParseToken(refreshToken)
	if err != nil {
		return err
	}
	delete(refreshTokenStore, claims.ID)
	return nil
}
