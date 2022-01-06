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
	tx := g.Con().Portal.Model(sys.Param{})
	tx = tx.Select("value")
	tx = tx.Where("`key` = ?", key)
	tx = tx.Find(&value)
	return value, tx.Error
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
	tx := g.Con().Portal.Model(sys.Param{})
	tx.Select("value")
	tx.Where("`key` = ?", "tree.order")
	if tx = tx.Find(&value); tx.Error != nil {
		return []string{}, errors.New("there's no such record")
	}

	var categoryNames []string
	err := json.Unmarshal([]byte(value), &categoryNames)
	return categoryNames, err
}

func newParamService() *paramService {
	return &paramService{}
}
