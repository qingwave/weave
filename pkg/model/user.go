package model

import (
	"encoding/json"
	"time"
)

const (
	UserAssociation         = "Users"
	UserAuthInfoAssociation = "AuthInfos"
	GroupAssociation        = "Groups"
)

type User struct {
	ID        uint       `json:"id" gorm:"autoIncrement;primaryKey"`
	Name      string     `json:"name" gorm:"size:100;not null;unique"`
	Password  string     `json:"-" gorm:"size:256;"`
	Email     string     `json:"email" gorm:"size:256;"`
	Avatar    string     `json:"avatar" gorm:"size:256;"`
	AuthInfos []AuthInfo `json:"authInfos" gorm:"foreignKey:UserId;references:ID"`
	Groups    []Group    `json:"groups" gorm:"many2many:user_groups;"`
	Roles     []Role     `json:"roles" gorm:"many2many:user_roles;"`

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
	Url          string    `json:"url" gorm:"size:256"`
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

type UserRole struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (u *UserRole) GetUser() *User {
	return &User{
		ID:   u.ID,
		Name: u.Name,
	}
}

type UserInfo struct {
	User
	InRoot bool   `json:"inRoot"`
	Role   string `json:"role"`
}

type Users []User
