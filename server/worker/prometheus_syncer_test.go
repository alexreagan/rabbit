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
}

func TestInitPrometheusSyncerConfigFromDB(t *testing.T) {
	cfg, _ := loadPrometheusSyncerConfigFromDB()
	log.Printf("[TestInitPrometheusSyncerConfigFromDB] config: %+v", cfg)
}

func TestUnmarshalSyncNodeUnameInfo(t *testing.T) {
	content := `
{
	"status": "success",
	"data": {
		"resultType": "vector",
		"result": [
			{
				"metric": {
					"__name__": "node_uname_info",
					"domainname": "(none)",
					"instance": "128.180.66.1:19100",
					"job": "prometheus-client",
					"machine": "x86_64",
					"nodename": "bd41appdpwb1001",
					"release": "2.6.32-573.e16.x86_64",
					"sysname": "Linux",
					"version": "#1 SMP Wed Jul 1 18:23:37 EDT 2015"
				},
				"value": [1641364000.775, "1"]
			}
		]
	}
}`
	var syncNodeUnameInfo SyncNodeUnameInfoResp
	e := json.Unmarshal([]byte(content), &syncNodeUnameInfo)
	if e != nil {
		log.Error(e)
	} else {
		log.Printf("=======%+v", syncNodeUnameInfo)
	}
	return
}
