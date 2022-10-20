package service

import (
	"strconv"

	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/repository"
)

type rbacService struct {
	rbacRepository repository.RBACRepository
}

func NewRBACService(rbacRepository repository.RBACRepository) RBACService {
	return &rbacService{
		rbacRepository: rbacRepository,
	}
}

func (rbac *rbacService) List() ([]model.Role, error) {
	return rbac.rbacRepository.ListRoles()
}

func (rbac *rbacService) Create(role *model.Role) (*model.Role, error) {
	return rbac.rbacRepository.CreateRole(role)
}

func (rbac *rbacService) Get(id string) (*model.Role, error) {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return rbac.rbacRepository.GetRole(rid)
}

func (rbac *rbacService) Update(id string, role *model.Role) (*model.Role, error) {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	role.ID = uint(rid)
	return rbac.rbacRepository.UpdateRole(role)
}

func (rbac *rbacService) Delete(id string) error {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return rbac.rbacRepository.DeleteRole(uint(rid))
}
