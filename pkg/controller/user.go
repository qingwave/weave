package controller

import (
	"net/http"
	"strconv"

	"github.com/qingwave/weave/pkg/authorization"
	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/service"
	"github.com/qingwave/weave/pkg/utils/trace"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) Controller {
	return &UserController{
		userService: userService,
	}
}

// @Summary List user
// @Description List user and storage
// @Produce json
// @Tags user
// @Security JWT
// @Success 200 {object} common.Response{data=model.Users}
// @Router /api/v1/users [get]
func (u *UserController) List(c *gin.Context) {
	common.TraceStep(c, "start list user")
	users, err := u.userService.List()
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.TraceStep(c, "list user done")
	common.ResponseSuccess(c, users)
}

// @Summary Get user
// @Description Get user and storage
// @Produce json
// @Tags user
// @Security JWT
// @Param id path int true "user id"
// @Success 200 {object} common.Response{data=model.User}
// @Router /api/v1/users/{id} [get]
func (u *UserController) Get(c *gin.Context) {
	user, err := u.userService.Get(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, user)
}

// @Summary Create user
// @Description Create user and storage
// @Accept json
// @Produce json
// @Tags user
// @Security JWT
// @Param user body model.CreatedUser true "user info"
// @Success 200 {object} common.Response{data=model.User}
// @Router /api/v1/users [post]
func (u *UserController) Create(c *gin.Context) {
	createdUser := new(model.CreatedUser)
	if err := c.BindJSON(createdUser); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	user := createdUser.GetUser()
	if err := u.userService.Validate(user); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	u.userService.Default(user)
	common.TraceStep(c, "start create user", trace.Field{"user", user.Name})
	defer common.TraceStep(c, "create user done", trace.Field{"user", user.Name})
	user, err := u.userService.Create(user)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}

	common.ResponseSuccess(c, user)
}

// @Summary Update user
// @Description Update user and storage
// @Accept json
// @Produce json
// @Tags user
// @Security JWT
// @Param user body model.UpdatedUser true "user info"
// @Param id   path      int  true  "user id"
// @Success 200 {object} common.Response{data=model.User}
// @Router /api/v1/users/{id} [put]
func (u *UserController) Update(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil || (strconv.Itoa(int(user.ID)) != c.Param("id") && !authorization.IsClusterAdmin(user)) {
		common.ResponseFailed(c, http.StatusForbidden, nil)
		return
	}

	new := new(model.UpdatedUser)
	if err := c.BindJSON(new); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	logrus.Infof("get update user: %#v, user: %#v", new, new.GetUser())

	common.TraceStep(c, "start update user", trace.Field{"user", new.Name})
	defer common.TraceStep(c, "update user done", trace.Field{"user", new.Name})

	user, err := u.userService.Update(c.Param("id"), new.GetUser())
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	common.ResponseSuccess(c, user)
}

// @Summary Delete user
// @Description Delete user and storage
// @Produce json
// @Tags user
// @Security JWT
// @Param id path int true "user id"
// @Success 200 {object} common.Response
// @Router /api/v1/users/{id} [delete]
func (u *UserController) Delete(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil || (strconv.Itoa(int(user.ID)) != c.Param("id") && !authorization.IsClusterAdmin(user)) {
		common.ResponseFailed(c, http.StatusForbidden, nil)
		return
	}

	if err := u.userService.Delete(c.Param("id")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

// @Summary Get groups
// @Description Get groups
// @Produce json
// @Tags group
// @Security JWT
// @Param id path int true "user id"
// @Success 200 {object} common.Response
// @Router /api/v1/users/{id}/groups [get]
func (u *UserController) GetGroups(c *gin.Context) {
	groups, err := u.userService.GetGroups(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, groups)
}

// @Summary Add role
// @Description Add role to user
// @Produce json
// @Tags user
// @Security JWT
// @Param id path int true "user id"
// @Param rid path int true "role id"
// @Success 200 {object} common.Response
// @Router /api/v1/users/{id}/roles/{rid} [post]
func (u *UserController) AddRole(c *gin.Context) {
	if err := u.userService.AddRole(c.Param("id"), c.Param("rid")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

// @Summary Delete role
// @Description delete role from user
// @Produce json
// @Tags user
// @Security JWT
// @Param id path int true "user id"
// @Param rid path int true "role id"
// @Success 200 {object} common.Response
// @Router /api/v1/users/{id}/roles/{rid} [delete]
func (u *UserController) DelRole(c *gin.Context) {
	if err := u.userService.DelRole(c.Param("id"), c.Param("rid")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

func (u *UserController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/users", u.List)
	api.POST("/users", u.Create)
	api.GET("/users/:id", u.Get)
	api.PUT("/users/:id", u.Update)
	api.DELETE("/users/:id", u.Delete)
	api.GET("/users/:id/groups", u.GetGroups)
	api.POST("/users/:id/roles/:rid", u.AddRole)
	api.DELETE("/users/:id/roles/:rid", u.DelRole)
}

func (u *UserController) Name() string {
	return "User"
}
