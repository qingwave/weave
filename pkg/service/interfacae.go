package service

import "github.com/qingwave/weave/pkg/model"

type UserService interface {
	List() (model.Users, error)
	Create(*model.User) (*model.User, error)
	Get(string) (*model.User, error)
	CreateOAuthUser(user *model.User) (*model.User, error)
	Update(string, *model.User) (*model.User, error)
	Delete(string) error
	Validate(*model.User) error
	Auth(*model.AuthUser) (*model.User, error)
	Default(*model.User)
	GetGroups(string) ([]model.Group, error)
}

type GroupService interface {
	List() ([]model.Group, error)
	Create(*model.User, *model.Group) (*model.Group, error)
	Get(string) (*model.Group, error)
	Update(string, *model.Group) (*model.Group, error)
	Delete(string) error
	GetUsers(gid string) ([]model.UserRole, error)
	AddUser(user *model.UserRole, gid string) error
	DelUser(user *model.UserRole, gid string) error
}
