package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

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
	Rules     Rules  `json:"rules" gorm:"type:json"`
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
	Resource  string    `json:"resource"`
	Operation Operation `json:"operation"`
}

type Rules []Rule

func (r *Rules) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	result := Rules{}
	err := json.Unmarshal(bytes, &result)
	*r = result
	return err
}

func (r Rules) Value() (driver.Value, error) {
	b, err := json.Marshal(r)
	return string(b), err
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
