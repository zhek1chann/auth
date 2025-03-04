package converter

import (
	"auth/internal/model"
	descUser "auth/pkg/auth_v1"
)

func ToCreateUserInfoFromDesc(user *descUser.RegisterUserInfo) *model.AuthUser {
	return &model.AuthUser{
		Info:     ToUserInfoFromDesc(user.Info),
		Password: user.Password,
	}
}

//func ToUserFromService(user *model.User) *descUser.User {
//	var updatedAt *timestamppb.Timestamp
//	if user.UpdatedAt.Valid {
//		updatedAt = timestamppb.New(user.UpdatedAt.Time)
//	}
//
//	return &descUser.User{
//		Id:        user.ID,
//		Info:      ToUserInfoFromService(user.Info),
//		CreatedAt: timestamppb.New(user.CreatedAt),
//		UpdatedAt: updatedAt,
//	}
//}

func ToUserInfoFromService(info *model.UserInfo) *descUser.UserInfo {
	return &descUser.UserInfo{
		Name:        info.Name,
		PhoneNumber: info.PhoneNumber,
	}
}

func ToUserInfoFromDesc(info *descUser.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:        info.Name,
		PhoneNumber: info.PhoneNumber,
		Role:        int(info.Role),
	}
}
