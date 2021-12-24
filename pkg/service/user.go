package service

import (
	"errors"
	"fmt"
	"strconv"

	"weave/pkg/model"
)

type userService struct {
	userRepository model.UserRepository
}

func NewUserService(userRepository model.UserRepository) model.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) List() (model.Users, error) {
	return u.userRepository.List()
}

func (u *userService) Create(user *model.User) (*model.User, error) {
	return u.userRepository.Create(user)
}

func (u *userService) Get(id string) (*model.User, error) {
	return u.getUserByID(id)
}

func (u *userService) Update(id string, new *model.User) (*model.User, error) {
	old, err := u.getUserByID(id)
	if err != nil {
		return nil, err
	}

	if new.ID != 0 && old.ID != new.ID {
		return nil, fmt.Errorf("update user %s not match", id)
	}

	return u.userRepository.Update(new)
}

func (u *userService) Delete(id string) error {
	user, err := u.getUser(id)
	if err != nil {
		return err
	}
	return u.userRepository.Delete(user)
}

func (u *userService) Validate(user *model.User) error {
	if user == nil {
		return errors.New("user is empty")
	}
	if user.Name == "" {
		return errors.New("user name is empty")
	}
	return nil
}

func (u *userService) Default(user *model.User) {
	if user == nil || user.Name == "" {
		return
	}
	if user.Email == "" {
		user.Email = fmt.Sprintf("%s@qinng.io", user.Name)
	}
}

func (u *userService) getUserByID(id string) (*model.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return u.userRepository.GetUserByID(uid)
}

func (u *userService) getUser(id string) (*model.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID: uid,
	}, nil
}
