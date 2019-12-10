package sevices

import (
	"strconv"

	"qinng.io/weave/pkg/database"
	"qinng.io/weave/pkg/models"
)

type detailUserBackend struct{}

func (u *detailUserBackend) List() (models.DetailUsers, error) {
	users := make(models.Users, 0)
	if err := database.Mysql.Unscoped().Find(&users).Error; err != nil {
		return nil, err
	}

	detailUsers := make(models.DetailUsers, 0)
	for _, u := range users {
		detailUsers = append(detailUsers, models.DetailUser{
			User:        u,
			DeletedTime: u.DeletedAt,
		})
	}

	return detailUsers, nil
}

func (u *detailUserBackend) Get(id string) (*models.DetailUser, error) {
	return getDetailUserByID(id)
}

func (u *detailUserBackend) Delete(id string) error {
	user, err := getDetailUserByID(id)
	if err != nil {
		return err
	}
	if err := database.Mysql.Unscoped().Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func getDetailUserByID(id string) (*models.DetailUser, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	user := new(models.User)
	if err := database.Mysql.Unscoped().First(user, uid).Error; err != nil {
		return nil, err
	}
	return &models.DetailUser{
		User:        *user,
		DeletedTime: user.DeletedAt,
	}, nil
}
