package model

import (
	"time"
)

type BaseUser struct {
	Name  string `json:"name" gorm:"type:varchar(256)"`
	Email string `json:"email"`
}

type User struct {
	BaseUser
	ID        int        `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time  `json:"create_time"`
	UpdatedAt time.Time  `json:"update_time"`
	DeletedAt *time.Time `json:"-"` // soft delete
}

type Users []User

type UserService interface {
	List() (Users, error)
	Create(*User) (*User, error)
	Get(string) (*User, error)
	Update(string, *User) (*User, error)
	Delete(string) error
	Validate(*User) error
	Default(*User)
}

type UserRepository interface {
	GetUserByID(id int) (*User, error)
	List() (Users, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)
	Delete(user *User) error
	Migrate() error
}
