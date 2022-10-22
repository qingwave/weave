package repository

import (
	"context"

	"github.com/qingwave/weave/pkg/model"
	"gorm.io/gorm/clause"
)

type Repository interface {
	User() UserRepository
	Group() GroupRepository
	Post() PostRepository
	RBAC() RBACRepository
	Close() error
	Ping(ctx context.Context) error
	Init() error
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
	AddRole(role *model.Role, user *model.User) error
	DelRole(role *model.Role, user *model.User) error
	GetGroups(*model.User) ([]model.Group, error)
	Migrate() error
}

type GroupRepository interface {
	GetGroupByID(uint) (*model.Group, error)
	GetGroupByName(string) (*model.Group, error)
	List() ([]model.Group, error)
	Create(*model.User, *model.Group) (*model.Group, error)
	CreateGroups(groups []model.Group, conds ...clause.Expression) error
	Update(*model.Group) (*model.Group, error)
	Delete(uint) error
	GetUsers(*model.Group) (model.Users, error)
	AddUser(user *model.User, group *model.Group) error
	DelUser(user *model.User, group *model.Group) error
	AddRole(role *model.Role, group *model.Group) error
	DelRole(role *model.Role, group *model.Group) error
	RoleBinding(role *model.Role, group *model.Group) error
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

type RBACRepository interface {
	List() ([]model.Role, error)
	ListResources() ([]model.Resource, error)
	Create(role *model.Role) (*model.Role, error)
	CreateResource(resource *model.Resource) (*model.Resource, error)
	CreateResources(resources []model.Resource, conds ...clause.Expression) error
	GetRoleByID(id int) (*model.Role, error)
	GetResource(id int) (*model.Resource, error)
	GetRoleByName(name string) (*model.Role, error)
	Update(role *model.Role) (*model.Role, error)
	Delete(id uint) error
	DeleteResource(id uint) error
	Migrate() error
}
