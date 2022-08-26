package repository

import "github.com/qingwave/weave/pkg/model"

type Repository interface {
	User() UserRepository
	Group() GroupRepository
	Post() PostRepository
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

type PostRepository interface {
	GetPostByID(uint) (*model.Post, error)
	GetPostByName(string) (*model.Post, error)
	List() ([]model.Post, error)
	Create(*model.User, *model.Post) (*model.Post, error)
	Update(*model.Post) (*model.Post, error)
	Delete(uint) error
	GetTags(*model.Post) ([]model.Tag, error)
	GetCategories(*model.Post) ([]model.Category, error)
	IncView(id uint) error
	AddLike(pid, uid uint) error
	DelLike(pid, uid uint) error
	GetLike(pid, uid uint) (bool, error)
	GetLikeByUser(uid uint) ([]model.Like, error)
	AddComment(comment *model.Comment) (*model.Comment, error)
	DelComment(id string) error
	ListComment(pid string) ([]model.Comment, error)
	Migrate() error
}
