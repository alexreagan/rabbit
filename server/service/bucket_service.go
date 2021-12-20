package service

import (
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/node"
)

type bucketService struct {
}

func (s *bucketService) Sort(hosts []*node.Host, tags []*app.Tag) (map[int64]*app.Tag, []*node.Host) {
	// host tag map
	hostTagMap := make(map[int64]*app.Tag)
	for _, host := range hosts {
		for _, tag := range host.RelatedTags() {
			if _, ok := hostTagMap[tag.ID]; !ok {
				hostTagMap[tag.ID] = tag
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

	// host tag与tag桶的交集
	intersectTags := make([]*app.Tag, 0)
	for _, hostTag := range hostTagMap {
		for _, tag := range tagMap {
			if hostTag.ID == tag.ID {
				intersectTags = append(intersectTags, tag)
			}
		}
	}

	// tag 交集
	intersectTagMap := make(map[int64]*app.Tag)
	unTaggedHost := make([]*node.Host, 0)
	for _, tag := range intersectTags {
		intersectTagMap[tag.ID] = tag
	}
	for _, host := range hosts {
		tagged := false
		for _, tag := range host.RelatedTags() {
			if _, ok := intersectTagMap[tag.ID]; ok {
				//intersectTagMap[tag.Tag].RelatedHosts = append(intersectTagMap[tag.Tag].RelatedHosts, host)
				tagged = true
			}
		}
		if tagged == false {
			unTaggedHost = append(unTaggedHost, host)
		}
	}
	return intersectTagMap, unTaggedHost
}

func newBucketService() *bucketService {
	return &bucketService{}
}
