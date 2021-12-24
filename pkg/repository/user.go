package repository

import (
	"weave/pkg/model"

	"gorm.io/gorm"
)

var (
	userCreateField = []string{"name", "email"}
	userUpdateField = []string{"name", "email"}
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		db: db,
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
	return user, nil
}

func (u *userRepository) Update(user *model.User) (*model.User, error) {
	if err := u.db.Model(&model.User{}).Select(userUpdateField).Updates(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) Delete(user *model.User) error {
	return u.db.Delete(user).Error
}

func (u *userRepository) GetUserByID(id int) (*model.User, error) {
	user := new(model.User)
	if err := u.db.First(user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}
