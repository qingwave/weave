package sevices

import (
	"fmt"
	"strconv"

	"qinng.io/weave/pkg/database"
	"qinng.io/weave/pkg/models"
	"qinng.io/weave/pkg/utils"
)

type userBackend struct{}

func (u *userBackend) List() (models.Users, error) {
	users := make(models.Users, 0)
	if err := database.Mysql.Order("name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userBackend) Create(user *models.User) (*models.User, error) {
	if user.Email == "" {
		user.Email = fmt.Sprintf("%s@gmail.com", user.Name)
	}

	if err := database.Mysql.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userBackend) Get(id string) (*models.User, error) {
	return getUserByID(id)
}

func (u *userBackend) Update(id string, new *models.User) (*models.User, error) {
	old, err := getUserByID(id)
	if err != nil {
		return nil, err
	}

	if new.ID != 0 && old.ID != new.ID {
		return nil, utils.NotMatchError
	}

	if err := database.Mysql.Model(&models.User{}).Updates(new).Error; err != nil {
		return nil, err
	}

	return new, nil
}

func (u *userBackend) Delete(id string) error {
	user, err := getUserByID(id)
	if err != nil {
		return err
	}
	if err := database.Mysql.Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func getUserByID(id string) (*models.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	user := new(models.User)
	if err := database.Mysql.First(user, uid).Error; err != nil {
		return nil, err
	}
	return user, nil
}
