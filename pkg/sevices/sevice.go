package sevices

import (
	"qinng.io/weave/pkg/database"
	"qinng.io/weave/pkg/models"
)

type UserService interface {
	List() (models.Users, error)
	Create(*models.User) (*models.User, error)
	Get(string) (*models.User, error)
	Update(string, *models.User) (*models.User, error)
	Delete(string) error
}

type DetailUserService interface {
	List() (models.DetailUsers, error)
	Get(string) (*models.DetailUser, error)
	Delete(string) error
}

type WeaveSevice interface {
	Users() UserService
	DetailUser() DetailUserService
}

type weaveBackend struct {
	user       *userBackend
	detailUser *detailUserBackend
}

func (wb *weaveBackend) Users() UserService {
	return wb.user
}

func (wb *weaveBackend) DetailUser() DetailUserService {
	return wb.detailUser
}

func NewWeaveSevice() (WeaveSevice, error) {
	if err := database.Mysql.AutoMigrate(&models.User{}).Error; err != nil {
		return nil, err
	}
	return &weaveBackend{
		user:       &userBackend{},
		detailUser: &detailUserBackend{},
	}, nil
}
