package converter

import (
	"auth/internal/model"
	descAuth "auth/pkg/auth_v1"
)

func ToAuthUserFromAuthDesc(user *descAuth.RegisterUserInfo) *model.AuthUser {
	return &model.AuthUser{
		Info:     ToUserInfoFromAuthDesc(user.Info),
		Password: user.Password,
	}
}

func ToUserInfoFromAuthDesc(info *descAuth.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:        info.Name,
		PhoneNumber: info.PhoneNumber,
		Role:        int(info.Role),
	}
}
