package service

import (
	"fmt"
	"strconv"
	"weave/pkg/authorization"
	"weave/pkg/model"
	"weave/pkg/repository"
)

type groupService struct {
	userRepository  repository.UserRepository
	groupRepository repository.GroupRepository
}

func NewGroupService(groupRepository repository.GroupRepository, userRepository repository.UserRepository) GroupService {
	return &groupService{
		groupRepository: groupRepository,
		userRepository:  userRepository,
	}
}

func (g *groupService) List() ([]model.Group, error) {
	// names, err := authorization.Enforcer.GetAllDomains()
	// if err != nil {
	// 	return nil, err
	// }
	// groups := make([]model.Group, 0)
	// for _, name := range names {
	// 	groups = append(groups, model.Group{Name: name})
	// }
	// return groups, nil

	return g.groupRepository.List()
}

func (g *groupService) Create(user *model.User, group *model.Group) (*model.Group, error) {
	group, err := g.groupRepository.Create(user, group)
	if err != nil {
		return nil, err
	}

	// add tenant policy
	authorization.Enforcer.AddPolicies([][]string{
		{authorization.AdminRole, group.Name, "tenant_sys_resource", "*", "get,update,delete"},
		{authorization.AdminRole, group.Name, "tenant_resource", "*", "*"},
		{authorization.EditRole, group.Name, "tenant_resource", "*", "*"},
		{authorization.ViewRole, group.Name, "tenant_resource", "*", "get,list"},
	})

	if _, err := authorization.Enforcer.AddGroupingPolicy(user.Name, authorization.AdminRole, group.Name); err != nil {
		return nil, err
	}

	return group, nil
}

func (g *groupService) Get(id string) (*model.Group, error) {
	gid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return g.groupRepository.GetGroupByID(uint(gid))
}

func (g *groupService) Update(id string, group *model.Group) (*model.Group, error) {
	gid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	group.ID = uint(gid)
	return g.groupRepository.Update(group)
}

func (g *groupService) Delete(id string) error {
	gid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	group, err := g.groupRepository.GetGroupByID(uint(gid))
	if err != nil {
		return err
	}
	for _, val := range authorization.DefaultGroups {
		if val == group.Name {
			return fmt.Errorf("system group %s cannot be deleted", val)
		}
	}
	if _, err := authorization.Enforcer.DeleteDomains(group.Name); err != nil {
		return err
	}
	return g.groupRepository.Delete(uint(gid))
}

func (g *groupService) GetUsers(id string) ([]model.UserRole, error) {
	group, err := g.Get(id)
	if err != nil {
		return nil, err
	}

	users := make([]model.UserRole, 0)
	for _, role := range authorization.DefaultRoles {
		for _, user := range authorization.Enforcer.GetUsersForRoleInDomain(role, group.Name) {
			users = append(users, model.UserRole{
				Name: user,
				Role: role,
			})
		}
	}

	return users, nil
}

func (g *groupService) AddUser(ur *model.UserRole, id string) error {
	var err error
	if ur.ID == 0 && ur.Name == "" || ur.Role == "" {
		return fmt.Errorf("invaild user info")
	}

	user := ur.GetUser()
	if ur.Name == "" {
		user, err = g.userRepository.GetUserByID(ur.ID)
	} else if ur.ID == 0 {
		user, err = g.userRepository.GetUserByName(ur.Name)
	}
	if err != nil {
		return err
	}

	group, err := g.Get(id)
	if err != nil {
		return err
	}

	g.groupRepository.AddUser(user, group)
	_, err = authorization.Enforcer.AddGroupingPolicy(user.Name, ur.Role, group.Name)

	return err
}

func (g *groupService) DelUser(ur *model.UserRole, id string) error {
	var err error
	if ur.ID == 0 && ur.Name == "" {
		return fmt.Errorf("invaild user info")
	}

	user := ur.GetUser()
	if ur.Name == "" {
		user, err = g.userRepository.GetUserByID(ur.ID)
	} else if ur.ID == 0 {
		user, err = g.userRepository.GetUserByName(ur.Name)
	}
	if err != nil {
		return err
	}

	group, err := g.Get(id)
	if err != nil {
		return err
	}

	ok, err := authorization.Enforcer.RemoveGroupingPolicy(user.Name, ur.Role, group.Name)
	if ok && len(authorization.Enforcer.GetRolesForUserInDomain(user.Name, group.Name)) == 0 {
		g.groupRepository.DelUser(user, group)
	}

	return err
}
