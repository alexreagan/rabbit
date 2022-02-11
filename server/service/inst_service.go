package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/uic"
	"sync"
)

type instService struct {
	mu sync.Mutex
}

func (s *instService) GetUserInst(u *uic.User) (*uic.Inst, error) {
	var inst uic.Inst
	tx := g.Con().Uic.Model(uic.Inst{})
	if err := tx.Where("inst_id = ?", u.InstID).First(&inst).Error; err != nil {
		return nil, err
	}
	return &inst, nil
}

func newInstService() *instService {
	return &instService{}
}
