package repository

import (
	"github.com/qingwave/weave/pkg/database"
	"github.com/qingwave/weave/pkg/model"

	"gorm.io/gorm"
)

var (
	groupUpdateFields = []string{"describe"}
)

type groupRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func newGroupRepository(db *gorm.DB, rdb *database.RedisDB) GroupRepository {
	return &groupRepository{
		db:  db,
		rdb: rdb,
	}
}

func (g *groupRepository) List() ([]model.Group, error) {
	groups := make([]model.Group, 0)
	if err := g.db.Order("name").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *groupRepository) Create(user *model.User, group *model.Group) (*model.Group, error) {
	group.CreatorId = user.ID
	group.Users = []model.User{*user}
	err := g.db.Create(group).Error
	return group, err
}

func (g *groupRepository) GetUsers(group *model.Group) (model.Users, error) {
	users := make(model.Users, 0)
	err := g.db.Model(group).Association(model.UserAssociation).Find(&users)
	return users, err
}

func (g *groupRepository) AddUser(user *model.User, group *model.Group) error {
	return g.db.Model(group).Association(model.UserAssociation).Append(user)
}

func (g *groupRepository) DelUser(user *model.User, group *model.Group) error {
	return g.db.Model(group).Association(model.UserAssociation).Delete(user)
}

func (g *groupRepository) GetGroupByID(id uint) (*model.Group, error) {
	group := new(model.Group)
	if err := g.db.Preload(model.UserAssociation).First(group, id).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (g *groupRepository) GetGroupByName(name string) (*model.Group, error) {
	group := new(model.Group)
	if err := g.db.Preload(model.UserAssociation).Where("name = ?", name).First(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

func (g *groupRepository) Update(group *model.Group) (*model.Group, error) {
	err := g.db.Model(group).Select(groupUpdateFields).Updates(group).Error
	return group, err
}

func (g *groupRepository) Delete(id uint) error {
	return g.db.Delete(&model.Group{}, id).Error
}

func (g *groupRepository) Migrate() error {
	return g.db.AutoMigrate(&model.Group{})
}
