package auth

import (
	"context"
	"log"

	"auth/internal/converter"
	descAuth "auth/pkg/auth_v1"
	"auth/pkg/validator"
)

func (i *AuthImplementation) Register(ctx context.Context, req *descAuth.RegisterRequest) (*descAuth.RegisterResponse, error) {
	form := struct {
		*descAuth.RegisterUserInfo
		validator.Validator
	}{
		RegisterUserInfo: req.GetUser(),
	}
	// TODO: validator for phone number, and equal for password
	form.CheckField(validator.NotBlank(form.Info.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Info.PhoneNumber), "phone-number", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirm-password", "This field cannot be blank")

	if !form.Valid() {
		return nil, form.TransformToGrpcError()
	}

	id, err := i.authService.Register(ctx, converter.ToAuthUserFromAuthDesc(form.RegisterUserInfo))
	if err != nil {
		return nil, err
	}

	log.Printf("registered user with id: %d", id)

	return &descAuth.RegisterResponse{
		Id: id,
	}, nil
}

func (i *AuthImplementation) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
	form := struct {
		validator.Validator
	}{}

	// TODO: validator for phone number, and equal for password
	form.CheckField(validator.NotBlank(req.GetPhoneNumber()), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(req.GetPassword()), "password", "This field cannot be blank")

	if !form.Valid() {
		return nil, form.TransformToGrpcError()
	}

	accessToken, refreshToken, err := i.authService.Login(ctx, req.GetPhoneNumber(), req.GetPassword())

	if err != nil {
		return nil, err
	}
	return &descAuth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

// func (s *AuthImplementation) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
// 	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
// 	if err != nil {
// 		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
// 	}

// 	// Можем слазать в базу или в кэш за доп данными пользователя

// 	refreshToken, err := utils.GenerateToken(model.UserInfo{
// 		Username: claims.Username,
// 		// Это пример, в реальности роль должна браться из базы или кэша
// 		Role: "admin",
// 	},
// 		[]byte(refreshTokenSecretKey),
// 		refreshTokenExpiration,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &descAuth.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
// }

// func (s *AuthImplementation) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
// 	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
// 	if err != nil {
// 		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
// 	}

// 	// Можем слазать в базу или в кэш за доп данными пользователя

// 	accessToken, err := utils.GenerateToken(model.UserInfo{
// 		Username: claims.Username,
// 		// Это пример, в реальности роль должна браться из базы или кэша
// 		Role: "admin",
// 	},
// 		[]byte(accessTokenSecretKey),
// 		accessTokenExpiration,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &descAuth.GetAccessTokenResponse{AccessToken: accessToken}, nil
// }
