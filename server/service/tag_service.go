package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/caas"
	"github.com/alexreagan/rabbit/server/model/node"
	"sort"
)

type TagGraphNode struct {
	app.Tag
	// 当前tag下关联的所有机器
	Path              []int64      `json:"path"`
	RelatedHosts      []*node.Host `json:"relatedHosts"`
	RelatedHostsCount int          `json:"relatedHostsCount"`
	// 当前tag下未关联到子tag的机器
	UnTaggedHosts      []*node.Host `json:"unTaggedHosts"`
	UnTaggedHostsCount int          `json:"UnTaggedHostsCount"`
	// 当前tag下关联的所有Pod
	RelatedPods      []*caas.Pod `json:"relatedPods"`
	RelatedPodsCount int         `json:"relatedPodsCount"`
	// 当前tag下未关联到子tag的Pod
	UnTaggedPods      []*caas.Pod `json:"unTaggedPods"`
	UnTaggedPodsCount int         `json:"UnTaggedPodsCount"`
	Next              map[int64]*TagGraphNode
	Children          []interface{} `json:"children"`
}

func (t *TagGraphNode) Nexts() []*TagGraphNode {
	var nodes TagGraphNodes
	for _, x := range t.Next {
		nodes = append(nodes, x)
	}
	nodes.Sort()
	return nodes
}

type TagGraphNodes []*TagGraphNode

func (t TagGraphNodes) Len() int { return len(t) }

func (t TagGraphNodes) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

func (t TagGraphNodes) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TagGraphNodes) Sort() {
	sort.Sort(t)
}

func newTagGraphNode(n *app.Tag) *TagGraphNode {
	return &TagGraphNode{
		Tag:  *n,
		Next: make(map[int64]*TagGraphNode),
	}
}

var globalTagGraphNodeV2 *TagGraphNode

type tagService struct {
}

func (s *tagService) GlobalTagGraphNodeV2() *TagGraphNode {
	return globalTagGraphNodeV2
}

func (s *tagService) Get(id int64) *app.Tag {
	var tag app.Tag
	db := g.Con().Portal.Model(tag)
	db = db.Select("`tag`.*, `tag_category`.name as category_name")
	db = db.Joins("left join `tag_category` on `tag_category`.id = `tag`.`category_id`")
	db = db.Where("`tag`.id = ?", id)
	db.Find(&tag)
	return &tag
}

func (s *tagService) GetAll() []*app.Tag {
	var tags []*app.Tag
	db := g.Con().Portal.Model(app.Tag{})
	db = db.Select("`tag`.*, `tag_category`.name as category_name")
	db = db.Joins("left join `tag_category` on `tag_category`.id = `tag`.`category_id`")
	db.Find(&tags)
	return tags
}

// 将tags按照category分桶
// 返回slice
//		key: tag所属的层级(0,1,2....)
//		value: [][]*node.tag
func BucketTags(categoryNames []string, tags []*app.Tag) [][]*app.Tag {
	bucketTags := make([][]*app.Tag, len(categoryNames))
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

func (s *tagService) ReBuildGraph() *TagGraphNode {
	globalTagGraphNodeV2 = nil
	return s.BuildGraphV2()
}

func buildTaggedInformationV2(bucketTags [][]*app.Tag, idx int, nodePath []int64, node *TagGraphNode) {
	if idx >= len(bucketTags) {
		return
	}

	for _, tag := range bucketTags[idx] {
		if _, ok := node.Next[tag.ID]; !ok {
			node.Next[tag.ID] = newTagGraphNode(tag)
			node.Next[tag.ID].Path = append(nodePath, tag.ID)
			node.Next[tag.ID].RelatedHosts = HostService.HostsHavingTagIDs(node.Next[tag.ID].Path)
			node.Next[tag.ID].RelatedHostsCount = len(node.Next[tag.ID].RelatedHosts)
			node.Next[tag.ID].RelatedPods = CaasService.PodsHavingTagIDs(node.Next[tag.ID].Path)
			node.Next[tag.ID].RelatedPodsCount = len(node.Next[tag.ID].RelatedPods)
		}
		buildTaggedInformationV2(bucketTags, idx+1, node.Next[tag.ID].Path, node.Next[tag.ID])
	}
}

// 计算未打到子标签的机器信息
func buildUnTaggedInformation(n *TagGraphNode) {
	// 如果不存在next节点，说明不存在unTaggedHosts和unTaggedPods
	if n.Next == nil {
		return
	}

	var unTaggedHosts node.Hosts
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
	unTaggedHosts.Sort()
	n.UnTaggedHosts = unTaggedHosts
	n.UnTaggedHostsCount = len(n.UnTaggedHosts)

	// 子节点上所有的Pod
	var unTaggedPods caas.Pods
	podMap := make(map[int64]*caas.Pod)
	for _, nd := range n.Next {
		for _, pod := range nd.RelatedPods {
			if _, ok := podMap[pod.ID]; !ok {
				podMap[pod.ID] = pod
			}
		}
	}

	for _, pod := range n.RelatedPods {
		if _, ok := podMap[pod.ID]; !ok {
			unTaggedPods = append(unTaggedPods, pod)
		}
	}
	unTaggedPods.Sort()
	n.UnTaggedPods = unTaggedPods
	n.UnTaggedPodsCount = len(n.UnTaggedPods)

	// 递归计算下一个
	for _, nd := range n.Next {
		buildUnTaggedInformation(nd)
	}
}

func buildChildrenInformation(n *TagGraphNode) {
	children := make([]interface{}, 0, 0)
	for _, nd := range n.Next {
		children = append(children, app.Tag{
			ID:           nd.ID,
			Type:         "Children",
			Name:         nd.Name,
			CnName:       nd.CnName,
			CategoryID:   nd.CategoryID,
			CategoryName: nd.CategoryName,
			Remark:       nd.Remark,
		})
	}
	for _, host := range n.UnTaggedHosts {
		host.Type = "UnTaggedHost"
		children = append(children, host)
	}
	for _, pod := range n.UnTaggedPods {
		pod.Type = "UnTaggedPod"
		children = append(children, pod)
	}
	n.Children = children

	for _, nd := range n.Next {
		buildChildrenInformation(nd)
	}
}

func (s *tagService) RelatedHosts(tag *app.Tag) []*node.Host {
	var hosts []*node.Host
	db := g.Con().Portal.Model(node.Host{}).Debug()
	db = db.Select("distinct `host`.*")
	db = db.Joins("left join `host_tag_rel` on `host_tag_rel`.host = `host`.id")
	db = db.Where("`host_tag_rel`.tag = ?", tag.ID)
	db.Find(&hosts)
	return hosts
}

// 根据host上打的tag信息创建树
func (s *tagService) BuildGraphV2() *TagGraphNode {
	if globalTagGraphNodeV2 != nil {
		return globalTagGraphNodeV2
	}

	// 获取树结构
	categoryNames, err := ParamService.GetTreeOrder()
	if err != nil {
		return nil
	}

	// tag路由图
	globalTagGraphNodeV2 = newTagGraphNode(&app.Tag{})

	// 组织host tag路由图
	for _, host := range HostService.GetAll() {
		tags := host.RelatedTags()
		bucketTags := BucketTags(categoryNames, tags)

		buildTaggedInformationV2(bucketTags, 0, globalTagGraphNodeV2.Path, globalTagGraphNodeV2)
	}

	// 组织pod tag路由图
	for _, service := range CaasService.GetAllService() {
		tags := CaasService.GetServiceRelatedTags(service)
		bucketTags := BucketTags(categoryNames, tags)

		buildTaggedInformationV2(bucketTags, 0, globalTagGraphNodeV2.Path, globalTagGraphNodeV2)
	}

	// 补充未关联的节点信息
	buildUnTaggedInformation(globalTagGraphNodeV2)

	// 补充children信息
	buildChildrenInformation(globalTagGraphNodeV2)

	return globalTagGraphNodeV2
}

func newTagService() *tagService {
	return &tagService{}
}
