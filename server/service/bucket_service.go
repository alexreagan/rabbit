package service

import (
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/node"
)

type bucketService struct {
}

func (s *bucketService) Sort(nodes []*node.Node, tags []*app.Tag) (map[int64]*app.Tag, []*node.Node) {
	// node tag map
	nodeTagMap := make(map[int64]*app.Tag)
	for _, n := range nodes {
		for _, tag := range n.RelatedTags() {
			if _, ok := nodeTagMap[tag.ID]; !ok {
				nodeTagMap[tag.ID] = tag
			}
		}
	}

	// tag map
	tagMap := make(map[int64]*app.Tag)
	for _, tag := range tags {
		if _, ok := tagMap[tag.ID]; !ok {
			tagMap[tag.ID] = tag
		}
	}

	// n tag与tag桶的交集
	intersectTags := make([]*app.Tag, 0)
	for _, nodeTag := range nodeTagMap {
		for _, tag := range tagMap {
			if nodeTag.ID == tag.ID {
				intersectTags = append(intersectTags, tag)
			}
		}
	}

	// tag 交集
	intersectTagMap := make(map[int64]*app.Tag)
	unTaggedNode := make([]*node.Node, 0)
	for _, tag := range intersectTags {
		intersectTagMap[tag.ID] = tag
	}
	for _, n := range nodes {
		tagged := false
		for _, tag := range n.RelatedTags() {
			if _, ok := intersectTagMap[tag.ID]; ok {
				//intersectTagMap[tag.Tag].RelatedNodes = append(intersectTagMap[tag.Tag].RelatedNodes, n)
				tagged = true
			}
		}
		if tagged == false {
			unTaggedNode = append(unTaggedNode, n)
		}
	}
	return intersectTagMap, unTaggedNode
}

func newBucketService() *bucketService {
	return &bucketService{}
}
