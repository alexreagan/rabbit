package service

import (
	"encoding/json"
	"errors"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/sys"
	"strconv"
)

type paramService struct{}

func (s *paramService) get(key string) (string, error) {
	var value string
	db := g.Con().Portal.Debug().Model(sys.Param{})
	db = db.Select("value")
	db = db.Where("`key` = ?", key)
	db = db.Find(&value)
	return value, db.Error
}

func (s *paramService) Get(key string) (string, error) {
	return s.get(key)
}

func (s *paramService) GetInt(key string) (int, error) {
	value, err := s.get(key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

func (s *paramService) GetTreeOrder() ([]string, error) {
	var value string
	db := g.Con().Portal.Debug().Model(sys.Param{})
	db.Select("value")
	db.Where("`key` = ?", "tree.order")
	if db = db.Find(&value); db.Error != nil {
		return []string{}, errors.New("there's no such record")
	}

	var categoryNames []string
	err := json.Unmarshal([]byte(value), &categoryNames)
	return categoryNames, err
}

func newParamService() *paramService {
	return &paramService{}
}
