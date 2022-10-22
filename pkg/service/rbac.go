package service

import (
	"strconv"

	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/repository"
	"github.com/qingwave/weave/pkg/utils/request"
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
	return rbac.rbacRepository.List()
}

func (rbac *rbacService) Create(role *model.Role) (*model.Role, error) {
	return rbac.rbacRepository.Create(role)
}

func (rbac *rbacService) Get(id string) (*model.Role, error) {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return rbac.rbacRepository.GetRoleByID(rid)
}

func (rbac *rbacService) Update(id string, role *model.Role) (*model.Role, error) {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	role.ID = uint(rid)
	return rbac.rbacRepository.Update(role)
}

func (rbac *rbacService) Delete(id string) error {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return rbac.rbacRepository.Delete(uint(rid))
}

func (rbac *rbacService) ListResources() ([]model.Resource, error) {
	return rbac.rbacRepository.ListResources()
}

func (rbac *rbacService) ListOperations() ([]model.Operation, error) {
	return []model.Operation{
		model.AllOperation,
		model.EditOperation,
		model.ViewOperation,
		request.CreateOperation,
		request.PatchOperation,
		request.UpdateOperation,
		request.GetOperation,
		request.ListOperation,
		request.DeleteOperation,
		"log",
		"exec",
		"proxy",
	}, nil
}
