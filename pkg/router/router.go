package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"qinng.io/weave/pkg/controllers"
	"qinng.io/weave/pkg/utils"
)

func Routers() (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = utils.HttpErrorHandler

	ctl, err := controllers.NewWeaveController()
	if err != nil {
		return nil, err
	}

	e.GET("/", ctl.Index)
	e.GET("/users", ctl.User.List)
	e.POST("/users", ctl.User.Create)
	e.GET("/users/:id", ctl.User.Get)
	e.PUT("/users/:id", ctl.User.Update)
	e.DELETE("/users/:id", ctl.User.Delete)

	e.GET("/dusers", ctl.DetailUser.List)
	e.GET("/dusers/:id", ctl.DetailUser.Get)
	e.DELETE("/dusers/:id", ctl.DetailUser.Delete)

	return e, nil
}
