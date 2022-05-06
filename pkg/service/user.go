package service

import (
	"errors"
	"fmt"
	"strconv"

	"weave/pkg/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	MinPasswordLength = 6
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
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(password)
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
	new.ID = old.ID

	if len(new.Password) > 0 {
		password, err := bcrypt.GenerateFromPassword([]byte(new.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		new.Password = string(password)
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
	if len(user.Password) < MinPasswordLength {
		return fmt.Errorf("password length must great than %d", MinPasswordLength)
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

func (u *userService) Auth(auser *model.AuthUser) (*model.User, error) {
	if auser == nil || auser.Name == "" || auser.Password == "" {
		return nil, fmt.Errorf("name or password is empty")
	}

	user, err := u.userRepository.GetUserByName(auser.Name)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auser.Password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) CreateOAuthUser(user *model.User) (*model.User, error) {
	if len(user.AuthInfos) == 0 {
		return nil, fmt.Errorf("empty auth info")
	}
	authInfo := user.AuthInfos[0]
	old, err := u.userRepository.GetUserByAuthID(authInfo.AuthType, authInfo.AuthId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.userRepository.Create(user)
		}
		return nil, err
	}

	return old, nil
}

func (u *userService) getUserByID(id string) (*model.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return u.userRepository.GetUserByID(uint(uid))
}

func (u *userService) getUser(id string) (*model.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: uint(uid)}, nil
}
