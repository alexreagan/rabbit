package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/uic"
	"sync"
)

type userService struct {
	mu sync.Mutex
}

func (s *userService) Get(jgygUserID string) (*uic.User, error) {
	var user uic.User
	tx := g.Con().Uic.Model(uic.User{})
	if err := tx.Where("jgyg_user_id = ?", jgygUserID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func newUserService() *userService {
	return &userService{}
}
