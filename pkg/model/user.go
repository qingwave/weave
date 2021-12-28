package model

import (
	"time"
)

type User struct {
	ID        int        `json:"id" gorm:"autoIncrement;primaryKey"`
	Name      string     `json:"name" gorm:"type:varchar(50);unique;not null"`
	Password  string     `json:"-" gorm:"type:varchar(256);not null"`
	Email     string     `json:"email" gorm:"type:varchar(256);unique;not null"`
	CreatedAt time.Time  `json:"create_time"`
	UpdatedAt time.Time  `json:"update_time"`
	DeletedAt *time.Time `json:"-"` // soft delete
}

type CreatedUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *CreatedUser) GetUser() *User {
	return &User{
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
	}
}

type UpdatedUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *UpdatedUser) GetUser() *User {
	return &User{
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
	}
}

type AuthUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Users []User

type UserService interface {
	List() (Users, error)
	Create(*User) (*User, error)
	Get(string) (*User, error)
	Update(string, *User) (*User, error)
	Delete(string) error
	Validate(*User) error
	Auth(*AuthUser) (*User, error)
	Default(*User)
}

type UserRepository interface {
	GetUserByID(int) (*User, error)
	GetUserByName(string) (*User, error)
	List() (Users, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)
	Delete(user *User) error
	Migrate() error
}
