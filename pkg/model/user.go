package model

import (
	"time"
)

type User struct {
	ID        int        `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	Name      string     `json:"name" gorm:"type:varchar(50)"`
	Password  string     `json:"-" gorm:"type:varchar(256)"`
	Email     string     `json:"email" gorm:"type:varchar(256)"`
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
