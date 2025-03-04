package auth

import (
	"auth/internal/model"
	descAuth "auth/pkg/auth_v1"
	"context"
)

type AuthImplementation struct {
	descAuth.UnimplementedAuthV1Server
	authService IAuthService
}

func NewAuthImplementation(AuthInterface IAuthService) *AuthImplementation {
	return &AuthImplementation{
		authService: AuthInterface,
	}
}

type IAuthService interface {
	Register(ctx context.Context, user *model.AuthUser) (int64, error)
	Login(ctx context.Context, phoneNumber string, password string) (accessToken string, refreshToken string, err error)
}
