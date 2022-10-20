package model

import (
	"github.com/qingwave/weave/pkg/utils/request"
	"github.com/qingwave/weave/pkg/utils/set"
)

const (
	All = "*"
)

type Scope string

const (
	ClusterScope   Scope = "cluster"
	NamespaceScope Scope = "namespace"
)

type Role struct {
	ID        uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name      string `json:"name" gorm:"size:100;not null;unique"`
	Scope     Scope  `json:"scope" gorm:"size:100"`
	Namespace string `json:"namespace"  gorm:"size:100"`
	Rules     []Rule `json:"rules" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

const (
	AllOperation  Operation = "*"
	EditOperation Operation = "edit"
	ViewOperation Operation = "view"
)

type Operation string

var (
	EditOperationSet = set.NewString(request.CreateOperation, request.DeleteOperation, request.UpdateOperation, request.PatchOperation, request.GetOperation, request.ListOperation)
	ViewOperationSet = set.NewString(request.GetOperation, request.ListOperation)
)

func (op Operation) Contain(verb string) bool {
	switch op {
	case AllOperation:
		return true
	case EditOperation:
		return EditOperationSet.Has(verb)
	case ViewOperation:
		return ViewOperationSet.Has(verb)
	default:
		return string(op) == verb
	}
}

type Rule struct {
	ID        uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	RoleID    uint      `json:"roleId" gorm:"index:idx_role_rule,unique"`
	Resource  string    `json:"resource" gorm:"size:100;index:idx_role_rule,unique"`
	Operation Operation `json:"operation" gorm:"size:100;index:idx_role_rule,unique"`
}

const (
	ResourceKind = "resource"
	MenuKind     = "menu"
)

const (
	ContainerResource = "containers"
	PostResource      = "posts"
	UserResource      = "users"
	GroupResource     = "groups"
	RoleResource      = "roles"
	AuthResource      = "auth"
	NamespaceResource = "namespaces"
)

type Resource struct {
	ID    uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name  string `json:"name" gorm:"size:256;not null;unique"`
	Scope Scope  `json:"scope"`
	Kind  string `json:"kind"`
}
