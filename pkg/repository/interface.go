package repository

import "github.com/qingwave/weave/pkg/model"

type Repository interface {
	User() UserRepository
	Group() GroupRepository
	Migrant
}

type Migrant interface {
	Migrate() error
}

type UserRepository interface {
	GetUserByID(uint) (*model.User, error)
	GetUserByAuthID(authType, authID string) (*model.User, error)
	GetUserByName(string) (*model.User, error)
	List() (model.Users, error)
	Create(*model.User) (*model.User, error)
	Update(*model.User) (*model.User, error)
	Delete(*model.User) error
	AddAuthInfo(authInfo *model.AuthInfo) error
	DelAuthInfo(authInfo *model.AuthInfo) error
	GetGroups(*model.User) ([]model.Group, error)
	Migrate() error
}

type GroupRepository interface {
	GetGroupByID(uint) (*model.Group, error)
	GetGroupByName(string) (*model.Group, error)
	List() ([]model.Group, error)
	Create(*model.User, *model.Group) (*model.Group, error)
	Update(*model.Group) (*model.Group, error)
	Delete(uint) error
	GetUsers(*model.Group) (model.Users, error)
	AddUser(user *model.User, group *model.Group) error
	DelUser(user *model.User, group *model.Group) error
	Migrate() error
}
