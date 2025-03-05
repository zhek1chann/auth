package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JSONWebToken struct {
	jwtKey []byte
}

func NewJSONWebToken(jwtSecret string) *JSONWebToken {
	return &JSONWebToken{
		jwtKey: []byte(jwtSecret),
	}
}

type Claims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

func (sa *JSONWebToken) GetJWTKey() []byte {
	return sa.jwtKey
}

func (sa *JSONWebToken) GenerateJSONWebTokens(id int64, username string, role int) (string, string, error) {
	accessToken, err := sa.generateShortLivedJSONWebToken(id, username, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := sa.generateLongLivedJSONWebToken(id, username, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (sa *JSONWebToken) generateShortLivedJSONWebToken(id int64, username string, role int) (string, error) {
	expiration := time.Now().Add(5 * time.Minute)
	return sa.generateJSONWebToken(id, username, role, expiration)
}

func (sa *JSONWebToken) generateLongLivedJSONWebToken(id int64, username string, role int) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)
	return sa.generateJSONWebToken(id, username, role, expiration)
}

func (sa *JSONWebToken) generateJSONWebToken(id int64, username string, role int, expirationTime time.Time) (string, error) {
	claims := &Claims{
		Id:       id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sa.jwtKey)
}

func (sa *JSONWebToken) RefreshAccessToken(refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return sa.jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	newAccessToken, err := sa.generateShortLivedJSONWebToken(claims.Id, claims.Username, claims.Role)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}
