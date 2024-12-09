package local_storage

import (
	"fmt"
	"github.com/iamvkosarev/sso/back/internal/model"
)

type LocalStorage struct {
	Users map[string]*model.User
}

func NewStorage() *LocalStorage {
	return &LocalStorage{}
}

func (l *LocalStorage) AddUser(user model.User) error {
	if l.Users == nil {
		l.Users = make(map[string]*model.User)
	}
	if _, ok := l.Users[user.Email]; ok {
		return fmt.Errorf("User with email %s already exists", user.Email)
	}
	l.Users[user.Email] = &user
	return nil
}

func (l *LocalStorage) GetUser(email string) (*model.User, error) {
	if l.Users == nil {
		return nil, fmt.Errorf("No users found")
	}
	user, ok := l.Users[email]
	if !ok {
		return nil, fmt.Errorf("User with email %s not found", email)
	}
	return user, nil
}
