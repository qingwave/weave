package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"qinng.io/weave/pkg/models"
	"qinng.io/weave/pkg/sevices"
)

type UserController struct {
	userService sevices.UserService
}

func (u *UserController) List(c echo.Context) error {
	users, err := u.userService.List()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func (u *UserController) Get(c echo.Context) error {
	user, err := u.userService.Get(c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserController) Create(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	if user.Email == "" {
		user.Email = fmt.Sprintf("%s@gmail.com", user.Name)
	}
	user, err := u.userService.Create(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserController) Update(c echo.Context) error {
	new := new(models.User)
	if err := c.Bind(new); err != nil {
		return err
	}

	user, err := u.userService.Update(c.Param("id"), new)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserController) Delete(c echo.Context) error {
	if err := u.userService.Delete(c.Param("id")); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
