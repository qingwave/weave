package controller

import (
	"net/http"

	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/service"

	"github.com/gin-gonic/gin"
)

type RBACController struct {
	rbacService service.RBACService
}

func NewRbacController(rbacService service.RBACService) Controller {
	return &RBACController{rbacService: rbacService}
}

// @Summary List rbac role
// @Description List rbac role
// @Produce json
// @Tags rbac
// @Security JWT
// @Success 200 {object} common.Response{data=[]model.Role}
// @Router /api/v1/roles [get]
func (rbac *RBACController) List(c *gin.Context) {
	roles, err := rbac.rbacService.List()
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	common.ResponseSuccess(c, roles)
}

// @Summary Create rbac role
// @Description Create rbac role
// @Accept json
// @Produce json
// @Tags rbac
// @Security JWT
// @Param role body model.Role true "rbac role info"
// @Success 200 {object} common.Response
// @Router /api/v1/roles [post]
func (rbac *RBACController) Create(c *gin.Context) {
	role := &model.Role{}
	if err := c.BindJSON(role); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	role, err := rbac.rbacService.Create(role)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	common.ResponseSuccess(c, role)
}

// @Summary Get role
// @Description Get role
// @Produce json
// @Tags role
// @Security JWT
// @Param id path int true "role id"
// @Success 200 {object} common.Response{data=model.Role}
// @Router /api/v1/roles/{id} [get]
func (rbac *RBACController) Get(c *gin.Context) {
	role, err := rbac.rbacService.Get(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, role)
}

// @Summary Update rbac role
// @Description Update rbac role
// @Accept json
// @Produce json
// @Tags rbac
// @Security JWT
// @Param role body model.Role true "rbac role info"
// @Success 200 {object} common.Response
// @Param id path int true "role id"
// @Router /api/v1/roles/:id [put]
func (rbac *RBACController) Update(c *gin.Context) {
	role := &model.Role{}
	if err := c.BindJSON(role); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")
	role, err := rbac.rbacService.Update(id, role)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	common.ResponseSuccess(c, role)
}

// @Summary Delete role
// @Description Delete role
// @Produce json
// @Tags role
// @Security JWT
// @Param id path int true "role id"
// @Success 200 {object} common.Response
// @Router /api/v1/roles/{id} [delete]
func (rbac *RBACController) Delete(c *gin.Context) {
	if err := rbac.rbacService.Delete(c.Param("id")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

func (rbac *RBACController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/roles", rbac.List)
	api.POST("/roles", rbac.Create)
	api.GET("/roles/:id", rbac.Get)
	api.PUT("/roles/:id", rbac.Update)
	api.DELETE("/roles/:id", rbac.Delete)
}

func (rbac *RBACController) Name() string {
	return "RBAC"
}
