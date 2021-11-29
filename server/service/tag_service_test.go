package service

import (
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

func TestBuildGraph(t *testing.T) {
	graph := TagService.BuildGraph()
	for _, node := range graph.Next {
		log.Printf("%#v", node)
	}
}
