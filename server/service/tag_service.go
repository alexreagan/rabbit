package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/node"
)

type TagRouterGraphNode struct {
	node.Tag
	// 当前tag下关联的所有机器
	Path []int64 `json:"path" gorm:"-"`
	RelatedHosts      []*node.Host `json:"relatedHosts" gorm:"-"`
	RelatedHostsCount int     `json:"relatedHostsCount" gorm:"-"`
	// 当前tag下未关联到子tag的机器
	UnTaggedHosts      []*node.Host `json:"unTaggedHosts" gorm:"-"`
	UnTaggedHostsCount int     `json:"UnTaggedHostsCount" gorm:"-"`
	Next map[int64]*TagRouterGraphNode
}

func newTagRouterGraphNode(n *node.Tag) *TagRouterGraphNode {
	return &TagRouterGraphNode{
		Tag:  *n,
		Next: make(map[int64]*TagRouterGraphNode),
	}
}

var globalTagRouterGraphNode *TagRouterGraphNode

type tagService struct {
}

func (s *tagService) GlobalTagRouterGraph() *TagRouterGraphNode {
	return globalTagRouterGraphNode
}

func (s *tagService) Get(id int64) *node.Tag {
	var tag node.Tag
	db := g.Con().Portal.Model(tag)
	db = db.Select("`tag`.*, `tag_category`.name as category_name")
	db = db.Joins("left join `tag_category` on `tag_category`.id = `tag`.`category_id`")
	db = db.Where("`tag`.id = ?", id)
	db.Find(&tag)
	return &tag
}

func (s *tagService) GetAll() []*node.Tag {
	var tags []*node.Tag
	db := g.Con().Portal.Model(node.Tag{})
	db = db.Select("`tag`.*, `tag_category`.name as category_name")
	db = db.Joins("left join `tag_category` on `tag_category`.id = `tag`.`category_id`")
	db.Find(&tags)
	return tags
}

// 将tags按照category分桶
// 返回slice
//		key: tag所属的层级(0,1,2....)
//		value: [][]*node.tag
func BucketTags(categoryNames []string, tags []*node.Tag) [][]*node.Tag {
	bucketTags := make([][]*node.Tag, len(categoryNames))
	for _, tag := range tags {
		cIdx := -1
		for idx, categoryName := range categoryNames {
			if categoryName == tag.CategoryName {
				cIdx = idx
				break
			}
		}
		if cIdx != -1 {
			bucketTags[cIdx] = append(bucketTags[cIdx], tag)
		}
	}
	return bucketTags
}

func (s *tagService) ReBuildGraph() *TagRouterGraphNode {
	globalTagRouterGraphNode = nil
	return s.BuildGraph()
}

func build(bucketTags [][]*node.Tag, idx int, nodePath []int64, node *TagRouterGraphNode) {
	if idx >= len(bucketTags) {
		return
	}

	for _, tag := range bucketTags[idx] {
		if _, ok := node.Next[tag.ID]; !ok {
			node.Next[tag.ID] = newTagRouterGraphNode(tag)
			node.Next[tag.ID].Path = append(nodePath, tag.ID)
			node.Next[tag.ID].RelatedHosts = HostService.HostsHavingTagIDs(node.Next[tag.ID].Path)
			node.Next[tag.ID].RelatedHostsCount = len(node.Next[tag.ID].RelatedHosts)
		}
		build(bucketTags, idx+1, node.Next[tag.ID].Path, node.Next[tag.ID])
	}
}

func calUnTaggedHosts(n *TagRouterGraphNode) {
	var unTaggedHosts []*node.Host
	hostMap := make(map[int64]*node.Host)
	for _, nd := range n.Next {
		for _, host := range nd.RelatedHosts {
			if _, ok := hostMap[host.ID]; !ok {
				hostMap[host.ID] = host
			}
		}
	}

	for _, host := range n.RelatedHosts {
		if _, ok := hostMap[host.ID]; !ok {
			unTaggedHosts = append(unTaggedHosts, host)
		}
	}
	n.UnTaggedHosts = unTaggedHosts
	n.UnTaggedHostsCount = len(n.UnTaggedHosts)
}

func (s *tagService) RelatedHosts(tag *node.Tag) []*node.Host {
	var hosts []*node.Host
	db := g.Con().Portal.Model(node.Host{}).Debug()
	db = db.Select("distinct `host`.*")
	db = db.Joins("left join `host_tag_rel` on `host_tag_rel`.host = `host`.id")
	db = db.Where("`host_tag_rel`.tag = ?", tag.ID)
	db.Find(&hosts)
	return hosts
}

// 根据host上打的tag信息创建树
func (s *tagService) BuildGraph() *TagRouterGraphNode {
	if globalTagRouterGraphNode != nil {
		return globalTagRouterGraphNode
	}

	// 获取树结构
	categoryNames, err := ParamService.GetTreeOrder()
	if err != nil {
		return nil
	}

	// tag路由图
	globalTagRouterGraphNode = newTagRouterGraphNode(&node.Tag{})

	// 组织tag路由图
	for _, host := range HostService.GetAll() {
		tags := host.RelatedTags()
		bucketTags := BucketTags(categoryNames, tags)

		build(bucketTags, 0, globalTagRouterGraphNode.Path, globalTagRouterGraphNode)
	}

	// 补充未关联的节点信息
	calUnTaggedHosts(globalTagRouterGraphNode)

	return globalTagRouterGraphNode
}

func newTagService() *tagService {
	return &tagService{}
}
