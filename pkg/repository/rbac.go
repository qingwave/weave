package repository

import (
	"github.com/qingwave/weave/pkg/database"
	"github.com/qingwave/weave/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type rbacRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func newRBACRepository(db *gorm.DB, rdb *database.RedisDB) RBACRepository {
	return &rbacRepository{
		db:  db,
		rdb: rdb,
	}
}

func (rbac *rbacRepository) List() ([]model.Role, error) {
	roles := make([]model.Role, 0)
	if err := rbac.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (rbac *rbacRepository) ListResources() ([]model.Resource, error) {
	resources := make([]model.Resource, 0)
	if err := rbac.db.Order("name").Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func (rbac *rbacRepository) Create(role *model.Role) (*model.Role, error) {
	err := rbac.db.Create(role).Error
	return role, err
}

func (rbac *rbacRepository) CreateResource(resource *model.Resource) (*model.Resource, error) {
	err := rbac.db.Create(resource).Error
	return resource, err
}

func (rbac *rbacRepository) CreateResources(resources []model.Resource, conds ...clause.Expression) error {
	err := rbac.db.Clauses(conds...).Create(resources).Error
	return err
}

func (rbac *rbacRepository) GetRoleByID(id int) (*model.Role, error) {
	role := &model.Role{}
	err := rbac.db.First(role, id).Error
	return role, err
}

func (rbac *rbacRepository) GetResource(id int) (*model.Resource, error) {
	res := &model.Resource{}
	err := rbac.db.First(res, id).Error
	return res, err
}

func (rbac *rbacRepository) GetRoleByName(name string) (*model.Role, error) {
	role := new(model.Role)
	if err := rbac.db.Preload(model.UserAssociation).Where("name = ?", name).First(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (rbac *rbacRepository) Update(role *model.Role) (*model.Role, error) {
	err := rbac.db.Updates(role).Error
	return role, err
}

func (rbac *rbacRepository) Delete(id uint) error {
	return rbac.db.Delete(&model.Role{}, id).Error
}

func (rbac *rbacRepository) DeleteResource(id uint) error {
	return rbac.db.Delete(&model.Resource{}, id).Error
}

func (rbac *rbacRepository) Migrate() error {
	return rbac.db.AutoMigrate(&model.Role{}, &model.Resource{})
}
