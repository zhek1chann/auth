package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserInfo struct {
	Name           string
	PhoneNumber          string
	HassedPassword string
	Role           int
}
