package controller

import (
	"fmt"
	"net/http"

	"weave/pkg/authorization"
	"weave/pkg/common"
	"weave/pkg/model"

	"github.com/gin-gonic/gin"
)

type RbacController struct{}

func NewRbacController() Controller {
	return &RbacController{}
}

// @Summary List rbac policy
// @Description List rbac policy
// @Produce json
// @Tags rbac
// @Security JWT
// @Success 200 {object} common.Response{data=[]string}
// @Param ptype    query     string  false  "ptype: p/g/g2"
// @Router /api/v1/policies policies [get]
func (rbac *RbacController) List(c *gin.Context) {
	ptype := c.Query("ptype")
	policies := make([][]string, 0)

	switch ptype {
	case model.DefaultPolicyType:
		policies = authorization.Enforcer.GetPolicy()
	case model.UserGroupPolicyType:
		policies = authorization.Enforcer.GetGroupingPolicy()
	case model.ResourcePolicyType:
		policies = authorization.Enforcer.GetNamedGroupingPolicy(model.ResourcePolicyType)
	default:
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("unsupported policy type: [%v]", ptype))
	}

	common.ResponseSuccess(c, policies)
}

// @Summary Handle rbac policy
// @Description Handle rbac policy and storage
// @Accept json
// @Produce json
// @Tags rbac
// @Security JWT
// @Param rbac policy body model.Policy true "rbac policy info"
// @Success 200 {object} common.Response
// @Router /api/v1/policies policies [post]
func (rbac *RbacController) Handle(c *gin.Context) {
	params := &model.Policy{}
	if err := c.BindJSON(&params); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	var ok bool
	var err error
	switch params.Action {
	case model.AddPolicyAction:
		ok, err = addPolicy(params.Type, toInterfaceSlice(params.Policy))
	case model.UpdatePolicyAction:
		ok, err = updatePolicy(params.Type, params.OldPolicy, params.Policy)
	case model.RemovePolicyAction:
		ok, err = removePolicy(params.Type, toInterfaceSlice(params.Policy))
	default:
		err = fmt.Errorf("invaild action %s", params.Action)
	}

	if err != nil || !ok {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to %s policy: %v", params.Action, err))
		return
	}

	common.ResponseSuccess(c, nil)
}

func (rbac *RbacController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/policies", rbac.List)
	api.POST("/policies", rbac.Handle)
}

func addPolicy(ptype string, policy []interface{}) (bool, error) {
	switch ptype {
	case model.DefaultPolicyType:
		return authorization.Enforcer.AddPolicy(policy...)
	case model.UserGroupPolicyType:
		return authorization.Enforcer.AddGroupingPolicy(policy...)
	case model.ResourcePolicyType:
		return authorization.Enforcer.AddNamedGroupingPolicy(ptype, policy...)
	}
	return false, fmt.Errorf("unsupported policy type %s", ptype)
}

func removePolicy(ptype string, policy []interface{}) (bool, error) {
	switch ptype {
	case model.DefaultPolicyType:
		return authorization.Enforcer.RemovePolicy(policy...)
	case model.UserGroupPolicyType:
		return authorization.Enforcer.RemoveGroupingPolicy(policy...)
	case model.ResourcePolicyType:
		return authorization.Enforcer.RemoveNamedGroupingPolicy(ptype, policy...)
	}
	return false, fmt.Errorf("unsupported policy type %s", ptype)
}

func updatePolicy(ptype string, old, new []string) (bool, error) {
	switch ptype {
	case model.DefaultPolicyType:
		return authorization.Enforcer.UpdatePolicy(old, new)
	case model.UserGroupPolicyType:
		return authorization.Enforcer.UpdateGroupingPolicy(old, new)
	case model.ResourcePolicyType:
		return authorization.Enforcer.UpdateNamedGroupingPolicy(ptype, old, new)
	}
	return false, fmt.Errorf("unsupported policy type %s", ptype)
}

func toInterfaceSlice(params []string) []interface{} {
	res := make([]interface{}, len(params))
	for i, param := range params {
		res[i] = param
	}
	return res
}
