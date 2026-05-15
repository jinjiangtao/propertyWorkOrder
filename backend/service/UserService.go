package service

import (
	"backend/entity"
	"backend/model"
	"errors"
)

func Register(username, password string) error {
	exists, err := model.CheckUsernameExists(username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}

	user := &entity.User{
		Username: username,
		Password: password,
		Role:     entity.RoleUser,
	}

	return model.CreateUser(user)
}

func Login(username, password string) (*entity.User, error) {
	user, err := model.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Password != password {
		return nil, errors.New("用户名或密码错误")
	}

	return user, nil
}

func AdminLogin(username, password string) (*entity.User, error) {
	user, err := model.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Role != entity.RoleAdmin {
		return nil, errors.New("无管理员权限")
	}

	if user.Password != password {
		return nil, errors.New("用户名或密码错误")
	}

	return user, nil
}