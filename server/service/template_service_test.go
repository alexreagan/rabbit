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
	g6Graph, _ := TemplateService.UnSerialize(template.Content)
	graph := TemplateService.BuildGraphV3(g6Graph)
	for _, node := range graph.Next {
		log.Printf("%#v", node)
	}
}
