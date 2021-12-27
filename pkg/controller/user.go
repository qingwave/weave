package controller

import (
	"weave/pkg/model"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService model.UserService
}

func NewUserController(userService model.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Summary List user
// @Description List user and storage
// @Produce json
// @Tags user
// @Success 200 {object} Response{data=model.Users}
// @Router /api/v1/users [get]
func (u *UserController) List(c *gin.Context) {
	users, err := u.userService.List()
	if err != nil {
		ResponseFailed(c, failed, err)
		return
	}
	ResponseSuccess(c, users)
}

// @Summary Get user
// @Description Get user and storage
// @Produce json
// @Tags user
// @Param id path int true "user id"
// @Success 200 {object} Response{data=model.User}
// @Router /api/v1/users/{id} [get]
func (u *UserController) Get(c *gin.Context) {
	user, err := u.userService.Get(c.Param("id"))
	if err != nil {
		ResponseFailed(c, failed, err)
		return
	}
	ResponseSuccess(c, user)
}

// @Summary Create user
// @Description Create user and storage
// @Accept json
// @Produce json
// @Tags user
// @Param user body model.CreatedUser true "user info"
// @Success 200 {object} Response{data=model.User}
// @Router /api/v1/users [post]
func (u *UserController) Create(c *gin.Context) {
	createdUser := new(model.CreatedUser)
	if err := c.BindJSON(createdUser); err != nil {
		ResponseFailed(c, failed, err)
		return
	}

	user := createdUser.GetUser()
	if err := u.userService.Validate(user); err != nil {
		ResponseFailed(c, failed, err)
		return
	}

	u.userService.Default(user)
	user, err := u.userService.Create(user)
	if err != nil {
		ResponseFailed(c, failed, err)
	}

	ResponseSuccess(c, user)
}

// @Summary Update user
// @Description Update user and storage
// @Accept json
// @Produce json
// @Tags user
// @Param user body model.UpdatedUser true "user info"
// @Param id   path      int  true  "user id"
// @Success 200 {object} Response{data=model.User}
// @Router /api/v1/users/{id} [put]
func (u *UserController) Update(c *gin.Context) {
	new := new(model.UpdatedUser)
	if err := c.Bind(new); err != nil {
		ResponseFailed(c, failed, err)
		return
	}

	user, err := u.userService.Update(c.Param("id"), new.GetUser())
	if err != nil {
		ResponseFailed(c, failed, err)
		return
	}

	ResponseSuccess(c, user)
}

// @Summary Delete user
// @Description Delete user and storage
// @Produce json
// @Tags user
// @Param id path int true "user id"
// @Success 200 {object} Response
// @Router /api/v1/users/{id} [delete]
func (u *UserController) Delete(c *gin.Context) {
	if err := u.userService.Delete(c.Param("id")); err != nil {
		ResponseFailed(c, failed, err)
		return
	}

	ResponseSuccess(c, nil)
}
