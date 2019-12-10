package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"qinng.io/weave/pkg/sevices"
)

type WeaveController struct {
	DetailUser *DetailUserController
	User       *UserController
}

func NewWeaveController() (*WeaveController, error) {
	svc, err := sevices.NewWeaveSevice()
	if err != nil {
		return nil, err
	}
	detailUserController := &DetailUserController{
		detailUserService: svc.DetailUser(),
	}
	userController := &UserController{
		userService: svc.Users(),
	}

	return &WeaveController{
		DetailUser: detailUserController,
		User:       userController,
	}, nil
}

func (w *WeaveController) Index(c echo.Context) error {
	return c.HTML(http.StatusOK, "<html><head><title>Weave Server</title></head><body><center><h1>Hello, Weave!</h1></center><hr><center>Weave/1.0</center></body><html>")
}
