package repository

import (
	"fmt"
	"strconv"

	"github.com/qingwave/weave/pkg/database"
	"github.com/qingwave/weave/pkg/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	userCreateField = []string{"name", "email", "password", "avatar", model.UserAuthInfoAssociation}
)

type userRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func newUserRepository(db *gorm.DB, rdb *database.RedisDB) UserRepository {
	return &userRepository{
		db:  db,
		rdb: rdb,
	}
}

func (u *userRepository) List() (model.Users, error) {
	users := make(model.Users, 0)
	if err := u.db.Preload(model.UserAuthInfoAssociation).Preload(model.GroupAssociation).Preload("Roles").Order("name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepository) Create(user *model.User) (*model.User, error) {
	if err := u.db.Select(userCreateField).Create(user).Error; err != nil {
		return nil, err
	}

	u.setCacheUser(user)

	return user, nil
}

func (u *userRepository) Update(user *model.User) (*model.User, error) {
	if err := u.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}

	u.rdb.HDel(user.CacheKey(), strconv.Itoa(int(user.ID)))

	return user, nil
}

func (u *userRepository) Delete(user *model.User) error {
	err := u.db.Select(model.UserAuthInfoAssociation).Delete(user).Error
	if err != nil {
		return err
	}
	u.rdb.HDel(user.CacheKey(), strconv.Itoa(int(user.ID)))
	return nil
}

func (u *userRepository) GetUserByID(id uint) (*model.User, error) {
	// TODO HSet not support expire, avoid roles and groups inconsistent
	// if user := u.getCacheUser(id); user != nil {
	// 	return user, nil
	// }

	user := new(model.User)
	if err := u.db.Omit("Password").Preload(model.UserAuthInfoAssociation).Preload("Groups").Preload("Groups.Roles").Preload("Roles").First(user, id).Error; err != nil {
		return nil, err
	}

	if err := u.setCacheUser(user); err != nil {
		logrus.Errorf("failed to set user: %v", err)
	}

	return user, nil
}

func (u *userRepository) GetUserByAuthID(authType, authID string) (*model.User, error) {
	authInfo := new(model.AuthInfo)
	if err := u.db.Where("auth_type = ? and auth_id = ?", authType, authID).First(authInfo).Error; err != nil {
		return nil, err
	}

	return u.GetUserByID(authInfo.UserId)
}

func (u *userRepository) GetUserByName(name string) (*model.User, error) {
	user := new(model.User)
	if err := u.db.Preload(model.UserAuthInfoAssociation).Preload("Groups").Preload("Groups.Roles").Preload("Roles").Where("name = ?", name).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) AddAuthInfo(authInfo *model.AuthInfo) error {
	if authInfo == nil {
		return nil
	}
	if authInfo.UserId == 0 {
		return fmt.Errorf("empty user id")
	}
	return u.db.Create(authInfo).Error
}

func (u *userRepository) DelAuthInfo(authInfo *model.AuthInfo) error {
	if authInfo == nil {
		return nil
	}
	return u.db.Delete(authInfo).Error
}

func (u *userRepository) AddRole(role *model.Role, user *model.User) error {
	return u.db.Model(user).Association("Roles").Append(role)
}

func (u *userRepository) DelRole(role *model.Role, user *model.User) error {
	return u.db.Model(user).Association("Roles").Delete(role)
}

func (u *userRepository) GetGroups(user *model.User) ([]model.Group, error) {
	groups := make([]model.Group, 0)
	err := u.db.Model(user).Association(model.GroupAssociation).Find(&groups)
	return groups, err
}

func (u *userRepository) Migrate() error {
	return u.db.AutoMigrate(&model.User{}, &model.AuthInfo{})
}

func (u *userRepository) setCacheUser(user *model.User) error {
	if user == nil {
		return nil
	}

	return u.rdb.HSet(user.CacheKey(), strconv.Itoa(int(user.ID)), user)
}

func (u *userRepository) getCacheUser(id uint) *model.User {
	user := new(model.User)
	key := user.CacheKey()
	field := strconv.Itoa(int(id))
	if err := u.rdb.HGet(key, field, user); err != nil {
		if err != database.RedisDisableError {
			logrus.Warnf("failed to hget field %s from key %s, %v", field, key, err)
		}
		return nil
	}

	return user
}
