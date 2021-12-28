package repository

import (
	"strconv"

	"weave/pkg/database"
	"weave/pkg/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	userCreateField = []string{"name", "email", "password"}
)

type userRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func NewUserRepository(db *gorm.DB, rdb *database.RedisDB) model.UserRepository {
	return &userRepository{
		db:  db,
		rdb: rdb,
	}
}

func (u *userRepository) List() (model.Users, error) {
	users := make(model.Users, 0)
	if err := u.db.Order("name").Find(&users).Error; err != nil {
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

	u.rdb.HDel(user.CacheKey(), strconv.Itoa(user.ID))

	return user, nil
}

func (u *userRepository) Delete(user *model.User) error {
	err := u.db.Delete(user).Error
	if err != nil {
		return err
	}
	u.rdb.HDel(user.CacheKey(), strconv.Itoa(user.ID))
	return nil
}

func (u *userRepository) GetUserByID(id int) (*model.User, error) {
	if user := u.getCacheUser(id); user != nil {
		return user, nil
	}

	user := new(model.User)
	if err := u.db.First(user, id).Error; err != nil {
		return nil, err
	}

	if err := u.setCacheUser(user); err != nil {
		logrus.Errorf("failed to set user", err)
	}

	return user, nil
}

func (u *userRepository) GetUserByName(name string) (*model.User, error) {
	user := new(model.User)
	if err := u.db.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *userRepository) setCacheUser(user *model.User) error {
	if user == nil {
		return nil
	}

	return u.rdb.HSet(user.CacheKey(), strconv.Itoa(user.ID), user)
}

func (u *userRepository) getCacheUser(id int) *model.User {
	user := new(model.User)
	key := user.CacheKey()
	field := strconv.Itoa(id)
	if err := u.rdb.HGet(key, field, user); err != nil {
		if err != database.RedisDisableError {
			logrus.Warnf("failed to hget field %s from key %s, %v", field, key, err)
		}
		return nil
	}

	return user
}
