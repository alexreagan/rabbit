package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/app"
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
	graph := TagService.BuildGraphV2()
	var resp []*TagGraphNode
	for _, n := range graph.Next {
		log.Printf("%#v", n)
		resp = append(resp, &TagGraphNode{
			Tag: app.Tag{
				ID:           n.ID,
				Name:         n.Name,
				CnName:       n.CnName,
				CategoryID:   n.CategoryID,
				Remark:       n.Remark,
				CategoryName: n.CategoryName,
			},
			Path:               n.Path,
			RelatedNodes:       n.RelatedNodes,
			RelatedNodesCount:  n.RelatedNodesCount,
			UnTaggedNodes:      n.UnTaggedNodes,
			UnTaggedNodesCount: n.UnTaggedNodesCount,
			RelatedPods:        n.RelatedPods,
			RelatedPodsCount:   n.RelatedPodsCount,
			UnTaggedPods:       n.UnTaggedPods,
			UnTaggedPodsCount:  n.UnTaggedPodsCount,
		})
	}
	log.Printf("%+v", resp)
}
