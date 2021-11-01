package uic

import (
	utils "github.com/Valiben/gin_unit_test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"rabbit/g"
	"rabbit/server/model/uic"
	"testing"
)

func init() {
	g.ParseConfig("../../../config/cfg.json")
	g.InitDBPool()

	router := gin.Default()
	Routes(router)
	utils.SetRouter(router)
}

func TestList(t *testing.T) {
	var data []uic.User
	err := utils.TestHandlerUnMarshalResp("GET", "/api/v1/user/list", "json", nil, &data)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	assert.NotNil(t, data)
}

func TestInfo(t *testing.T) {
	var data []uic.User
	err := utils.TestHandlerUnMarshalResp("GET", "/api/v1/user/info", "json", nil, &data)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	assert.NotNil(t, data)
}
