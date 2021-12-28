package sys

import (
	"github.com/alexreagan/rabbit/g"
	log "github.com/sirupsen/logrus"
	"testing"
)

func init() {
	g.ParseConfig("../../../config/cfg.json")
	g.InitDBPool()
}

func TestBuildTree(t *testing.T) {
	tree := Menu{}.BuildTree()
	log.Printf("[BuildTemplateGraph] BuildTemplateGraph: %+v", tree)
}
