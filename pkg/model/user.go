package model

import (
	"encoding/json"
	"time"
)

const (
	UserAuthInfoAssociation = "AuthInfos"
)

type User struct {
	ID        uint       `json:"id" gorm:"autoIncrement;primaryKey"`
	Name      string     `json:"name" gorm:"size:100;not null"`
	Password  string     `json:"-" gorm:"size:256;"`
	Email     string     `json:"email" gorm:"size:256;"`
	Avatar    string     `json:"avatar" gorm:"size:256;"`
	AuthInfos []AuthInfo `json:"authInfos" gorm:"foreignKey:UserId;references:ID"`

	BaseModel
}

func (*User) TableName() string {
	return "users"
}

func (u *User) CacheKey() string {
	return u.TableName() + ":id"
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

type AuthInfo struct {
	ID           uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	UserId       uint      `json:"userId" gorm:"size:256"`
	AuthType     string    `json:"authType" gorm:"size:256"`
	AuthId       string    `json:"authId" gorm:"size:256"`
	AccessToken  string    `json:"-" gorm:"size:256"`
	RefreshToken string    `json:"-" gorm:"size:256"`
	Expiry       time.Time `json:"-"`

	BaseModel
}

func (*AuthInfo) TableName() string {
	return "auth_infos"
}

type CreatedUser struct {
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	Avatar    string     `json:"avatar"`
	AuthInfos []AuthInfo `json:"authInfos"`
}

func (u *CreatedUser) GetUser() *User {
	return &User{
		Name:      u.Name,
		Password:  u.Password,
		Email:     u.Email,
		Avatar:    u.Avatar,
		AuthInfos: u.AuthInfos,
	}
}

type UpdatedUser struct {
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	AuthInfos []AuthInfo `json:"authInfos"`
}

func (u *UpdatedUser) GetUser() *User {
	return &User{
		Name:      u.Name,
		Password:  u.Password,
		Email:     u.Email,
		AuthInfos: u.AuthInfos,
	}
}

type AuthUser struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	SetCookie bool   `json:"setCookie"`
	AuthType  string `json:"authType"`
	AuthCode  string `json:"authCode"`
}

type Users []User

type UserService interface {
	List() (Users, error)
	Create(*User) (*User, error)
	Get(string) (*User, error)
	CreateOAuthUser(user *User) (*User, error)
	Update(string, *User) (*User, error)
	Delete(string) error
	Validate(*User) error
	Auth(*AuthUser) (*User, error)
	Default(*User)
}

type UserRepository interface {
	GetUserByID(uint) (*User, error)
	GetUserByAuthID(authType, authID string) (*User, error)
	GetUserByName(string) (*User, error)
	List() (Users, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)
	Delete(user *User) error
	AddAuthInfo(authInfo *AuthInfo) error
	DelAuthInfo(authInfo *AuthInfo) error
	Migrate() error
}
