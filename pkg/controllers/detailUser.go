package controllers

import (
	"net/http"

	"github.com/labstack/echo"

	"qinng.io/weave/pkg/sevices"
)

type DetailUserController struct {
	detailUserService sevices.DetailUserService
}

func (u *DetailUserController) List(c echo.Context) error {
	users, err := u.detailUserService.List()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func (u *DetailUserController) Get(c echo.Context) error {
	user, err := u.detailUserService.Get(c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (u *DetailUserController) Delete(c echo.Context) error {
	if err := u.detailUserService.Delete(c.Param("id")); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
