package worker

import (
	"encoding/json"
	"github.com/alexreagan/rabbit/g"
	log "github.com/sirupsen/logrus"
	"testing"
)

func init() {
	g.ParseConfig("../../config/cfg.json")
	g.InitDBPool()

	//router := gin.Default()
	//utils.SetRouter(router)
}

func TestUnmarshal(t *testing.T) {
	ws := &CaasWorkSpaceResult{}
	body := `
{
	"msg": "Success",
	"code": 0,
	"data": [{
		"id": 161,
		"createTime": "2021-06-23T14:11:05+08:00",
		"updateTime": null,
		"finishTime": "2021-06-23T14:11:05+08:00"
	}]
}
`
	json.Unmarshal([]byte(body), ws)
	log.Printf("[CaasSyncer] GetWorkSpaceObj: %+v", ws)
}
