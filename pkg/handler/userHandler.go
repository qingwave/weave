package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"qinng.io/weave/pkg/utils"
)

type User struct {
	ID        int        `json:"id" gorm:"AUTO_INCREMENT;primary_key"`
	Name      string     `json:"name" gorm:"type:varchar(256)"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"create_time"`
	UpdatedAt time.Time  `json:"update_time"`
	DeletedAt *time.Time `json:"-"` // 软删除
}

type AllUser struct {
	*User
	DeletedTime *time.Time `json:"delete_time"`
}

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) (*UserHandler, error) {
	if err := db.AutoMigrate(&User{}).Error; err != nil {
		return nil, err
	}
	return &UserHandler{
		db: db,
	}, nil
}

func (u *UserHandler) findUserByID(id string) (*User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, utils.NewHttpError(http.StatusBadRequest, "InvalidID", "请在URL中提供合法的ID")
	}
	user := new(User)
	if err := u.db.First(user, uid).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserHandler) List(c echo.Context) error {
	users := make([]*User, 0)
	if err := u.db.Order("name").Find(&users).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func (u *UserHandler) ListAll(c echo.Context) error {
	nt := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(nt)
	users := make([]*User, 0)
	if err := u.db.Unscoped().Where("deleted_at < ?", nt).Order("name").Find(&users).Error; err != nil {
		return err
	}
	allUsers := make([]*AllUser, 0)
	for _, u := range users {
		allUsers = append(allUsers, &AllUser{
			User:        u,
			DeletedTime: u.DeletedAt,
		})
	}
	return c.JSON(http.StatusOK, allUsers)
}

func (u *UserHandler) Get(c echo.Context) error {
	user, err := u.findUserByID(c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Create(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}

	if user.Email == "" {
		user.Email = fmt.Sprintf("%s@gmail.com", user.Name)
	}

	if err := u.db.Create(user).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Update(c echo.Context) error {
	user, err := u.findUserByID(c.Param("id"))
	if err != nil {
		return err
	}

	newUser := new(User)
	if err := c.Bind(newUser); err != nil {
		return err
	}

	if newUser.ID != 0 && user.ID != newUser.ID {
		return utils.NewHttpError(http.StatusBadRequest, "BadRequest", "请求ID不匹配")
	}

	if err := u.db.Model(&User{}).Updates(newUser).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newUser)
}

func (u *UserHandler) Delete(c echo.Context) error {
	user, err := u.findUserByID(c.Param("id"))
	if err != nil {
		return err
	}

	if err := u.db.Delete(user).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
