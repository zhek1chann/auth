package converter

import (
	"auth/internal/model"
	modelRepo "auth/internal/repository/user/model"
)

func ToUserFromRepo(note *modelRepo.User) *model.User {
	return &model.User{
		ID:        note.ID,
		Info:      ToUserInfoFromRepo(note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:        info.Name,
		PhoneNumber: info.PhoneNumber,
		Role:        info.Role,
	}
}
