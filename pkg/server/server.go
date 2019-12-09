package server

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"qinng.io/weave/pkg/handler"
	"qinng.io/weave/pkg/utils"
)

type Server struct {
	echo    *echo.Echo
	db      *gorm.DB
	uh      *handler.UserHandler
	address string
}

func NewServer(address string) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.HTTPErrorHandler = utils.HttpErrorHandler

	db, err := gorm.Open("mysql", "root:123456@tcp(localhost:3306)/weave?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		e.Logger.Fatalf("init db failed: %v", err)
		panic(err)
	}
	utils.Logger.Info("start db")

	uh, err := handler.NewUserHandler(db)
	if err != nil {
		e.Logger.Fatalf("init  failed: %v", err)
		panic(err)
	}

	s := Server{
		echo:    e,
		db:      db,
		uh:      uh,
		address: address,
	}

	e.GET("/", s.index)
	e.GET("/users", uh.List)
	e.GET("/allusers", uh.ListAll)
	e.POST("/users", uh.Create)
	e.GET("/users/:id", uh.Get)
	e.POST("/users/:id", uh.Update)
	e.DELETE("/users/:id", uh.Delete)

	return &s
}

func (s *Server) Run() {
	s.echo.Logger.Fatal(s.echo.Start(s.address))
}

func (s *Server) index(c echo.Context) error {
	s.echo.Logger.Warnf("start %s", time.Now().String())
	return c.HTML(http.StatusOK, "<center><strong>Hello, Weave!</strong></center>")
}
