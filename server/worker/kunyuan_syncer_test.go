package worker

import (
	"github.com/alexreagan/rabbit/g"
	log "github.com/sirupsen/logrus"
	"testing"
)

func init() {
	g.ParseConfig("../../config/cfg.json")
	g.InitDBPool()
}

func TestInitKunYuanSyncerConfigFromDB(t *testing.T) {
	cfg, _ := initKunYuanSyncerConfigFromDB()
	log.Printf("[TestInitKunYuanSyncerConfigFromDB] config: %+v", cfg)
}
