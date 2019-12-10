package models

import (
	"time"
)

type User struct {
	ID        int        `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	Name      string     `json:"name" gorm:"type:varchar(256)"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"create_time"`
	UpdatedAt time.Time  `json:"update_time"`
	DeletedAt *time.Time `json:"-"` // 软删除
}

type Users []User
