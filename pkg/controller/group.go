package controller

import (
	"fmt"
	"net/http"

	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/service"
	"github.com/qingwave/weave/pkg/utils/trace"

	"github.com/gin-gonic/gin"
)

type GroupController struct {
	groupService service.GroupService
}

func NewGroupController(groupService service.GroupService) Controller {
	return &GroupController{
		groupService: groupService,
	}
}

// @Summary List group
// @Description List group
// @Produce json
// @Tags group
// @Security JWT
// @Success 200 {object} common.Response{data=[]model.Group}
// @Router /api/v1/groups [get]
func (g *GroupController) List(c *gin.Context) {
	common.TraceStep(c, "start list group")
	groups, err := g.groupService.List()
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.TraceStep(c, "list group done")
	common.ResponseSuccess(c, groups)
}

// @Summary Get group
// @Description Get group
// @Produce json
// @Tags group
// @Security JWT
// @Param id path int true "group id"
// @Success 200 {object} common.Response{data=model.Group}
// @Router /api/v1/groups/{id} [get]
func (g *GroupController) Get(c *gin.Context) {
	group, err := g.groupService.Get(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, group)
}

// @Summary Create group
// @Description Create group and storage
// @Accept json
// @Produce json
// @Tags group
// @Security JWT
// @Param group body model.CreatedGroup true "group info"
// @Success 200 {object} common.Response{data=model.Group}
// @Router /api/v1/groups [post]
func (g *GroupController) Create(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
		return
	}

	createdGroup := new(model.CreatedGroup)
	if err := c.BindJSON(createdGroup); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	group := createdGroup.GetGroup(user.ID)
	common.TraceStep(c, "start create group", trace.Field{"group", group.Name})
	defer common.TraceStep(c, "create group done", trace.Field{"group", group.Name})

	group, err := g.groupService.Create(user, group)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	common.ResponseSuccess(c, group)
}

// @Summary Update group
// @Description Update group and storage
// @Accept json
// @Produce json
// @Tags group
// @Security JWT
// @Param group body model.UpdatedUser true "group info"
// @Param id   path      int  true  "group id"
// @Success 200 {object} common.Response{data=model.Group}
// @Router /api/v1/groups/{id} [put]
func (g *GroupController) Update(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
		return
	}

	id := c.Param("id")

	new := new(model.UpdatedGroup)
	if err := c.BindJSON(new); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.TraceStep(c, "start update group", trace.Field{"group", new.Name})
	defer common.TraceStep(c, "update group done", trace.Field{"group", new.Name})

	group, err := g.groupService.Update(id, new.GetGroup(user.ID))
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	common.ResponseSuccess(c, group)
}

// @Summary Delete group
// @Description Delete group
// @Produce json
// @Tags group
// @Security JWT
// @Param id path int true "group id"
// @Success 200 {object} common.Response
// @Router /api/v1/groups/{id} [delete]
func (g *GroupController) Delete(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
		return
	}

	if err := g.groupService.Delete(c.Param("id")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

// @Summary Get users
// @Description Get users
// @Produce json
// @Tags group
// @Security JWT
// @Param id path int true "group id"
// @Success 200 {object} common.Response
// @Router /api/v1/groups/{id}/users [get]
func (g *GroupController) GetUsers(c *gin.Context) {
	users, err := g.groupService.GetUsers(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, users)
}

// @Summary Add user
// @Description Add user to group
// @Produce json
// @Tags group
// @Security JWT
// @Param id path int true "group id"
// @Param user body model.User true "user info"
// @Success 200 {object} common.Response
// @Router /api/v1/groups/{id}/users [post]
func (g *GroupController) AddUser(c *gin.Context) {
	user := new(model.User)
	if err := c.BindJSON(user); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	if err := g.groupService.AddUser(user, c.Param("id")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

// @Summary Delete user
// @Description Delete user from group
// @Produce json
// @Tags group
// @Security JWT
// @Param id path int true "group id"
// @Param user body model.User true "user info"
// @Param name    query     string  true  "user name"
// @Param name    query     string  true  "user role"
// @Success 200 {object} common.Response
// @Router /api/v1/groups/{id}/users [delete]
func (g *GroupController) DelUser(c *gin.Context) {
	user := new(model.User)
	user.Name = c.Query("name")

	if err := g.groupService.DelUser(user, c.Param("id")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

// @Summary Add role
// @Description Add role to group
// @Produce json
// @Tags group
// @Security JWT
// @Param id path int true "group id"
// @Param rid path int true "role id"
// @Success 200 {object} common.Response
// @Router /api/v1/groups/{id}/roles/{rid} [post]
func (g *GroupController) AddRole(c *gin.Context) {
	if err := g.groupService.AddRole(c.Param("id"), c.Param("rid")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

// @Summary Delete role
// @Description delete role from group
// @Produce json
// @Tags group
// @Security JWT
// @Param id path int true "group id"
// @Param rid path int true "role id"
// @Success 200 {object} common.Response
// @Router /api/v1/groups/{id}/roles/{rid} [delete]
func (g *GroupController) DelRole(c *gin.Context) {
	if err := g.groupService.DelRole(c.Param("id"), c.Param("rid")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

func (g *GroupController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/groups", g.List)
	api.POST("/groups", g.Create)
	api.GET("/groups/:id", g.Get)
	api.PUT("/groups/:id", g.Update)
	api.DELETE("/groups/:id", g.Delete)
	api.GET("/groups/:id/users", g.GetUsers)
	api.POST("/groups/:id/users", g.AddUser)
	api.DELETE("/groups/:id/users", g.DelUser)
	api.POST("/groups/:id/roles/:rid", g.AddRole)
	api.DELETE("/groups/:id/roles/:rid", g.DelRole)
}

func (g *GroupController) Name() string {
	return "Group"
}
