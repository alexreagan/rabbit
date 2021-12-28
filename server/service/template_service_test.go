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

func TestBuildTree(t *testing.T) {
	template, _ := TemplateService.Get(1)
	graph := TemplateService.BuildTemplateGraph(template)
	for _, node := range graph.Next {
		log.Printf("%#v", node)
	}
}
