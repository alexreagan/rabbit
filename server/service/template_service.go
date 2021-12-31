package service

import (
	"encoding/json"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/app"
	"sync"
)

type templateService struct {
	mu sync.Mutex
}

// 反序列化content为对象
func (t *templateService) UnSerialize(c string) (*app.G6Graph, error) {
	var g6Graph app.G6Graph
	err := json.Unmarshal([]byte(c), &g6Graph)
	return &g6Graph, err
}

// 序列化对象
func (t *templateService) Serialize(g6Graph app.G6Graph) ([]byte, error) {
	return json.Marshal(g6Graph)
}

// 查找
func (t *templateService) Get(id int64) (*app.Template, error) {
	tx := g.Con().Portal
	template := app.Template{}
	if dt := tx.Model(app.Template{}).Where("id = ?", id).Find(&template); dt.Error != nil {
		return &template, dt.Error
	}
	return &template, nil
}

func (t *templateService) All() ([]*app.Template, error) {
	tx := g.Con().Portal
	var templates []*app.Template
	if dt := tx.Model(app.Template{}).Find(&templates); dt.Error != nil {
		return templates, dt.Error
	}
	return templates, nil
}

// 处于有效状态的template，有且只有一个
func (t *templateService) ValidTemplate() (*app.Template, error) {
	tx := g.Con().Portal
	template := app.Template{}
	if dt := tx.Model(app.Template{}).Where("state = 'enable'").First(&template); dt.Error != nil {
		return &template, dt.Error
	}
	return &template, nil
}

// 更新
func (t *templateService) Updates(template *app.Template) error {
	db := g.Con().Portal
	if db = db.Model(app.Template{}).Where("id = ?", template.ID).Updates(template); db.Error != nil {
		return db.Error
	}
	return nil
}

var globalTemplateGraphMap map[int64]*TagGraphNode

func (s *templateService) GlobalTemplateGraphMap() map[int64]*TagGraphNode {
	return globalTemplateGraphMap
}

// 组建树
func (t *templateService) BuildGraphs() map[int64]*TagGraphNode {
	if globalTemplateGraphMap == nil {
		globalTemplateGraphMap = make(map[int64]*TagGraphNode)
	}

	templates, _ := t.All()
	for _, template := range templates {
		globalTemplateGraphMap[template.ID] = t.BuildTemplateGraph(template)
	}
	return globalTemplateGraphMap
}

// 组建树结构
func (t *templateService) BuildTemplateGraph(template *app.Template) *TagGraphNode {
	t.mu.Lock()
	defer t.mu.Unlock()

	// UnSerialize template content
	g6Graph, _ := TemplateService.UnSerialize(template.Content)

	// build graph
	headMap := make(map[int64]*TagGraphNode)
	nodeMap := make(map[int64]*TagGraphNode)

	// tag路由图
	tagGraphNode := newTagGraphNode(&app.Tag{})

	// 初始化
	for _, n := range g6Graph.Nodes {
		//nid, _ := strconv.ParseInt(n.ID, 10, 64)
		tag := &app.Tag{
			ID:     n.ID,
			Name:   n.Name,
			CnName: n.Label,
		}
		nodeMap[n.ID] = newTagGraphNode(tag)
		headMap[n.ID] = newTagGraphNode(tag)
	}
	// 组织树状结构
	for _, edge := range g6Graph.Edges {
		//edgeSourceID, _ := strconv.ParseInt(edge.SourceID, 10, 64)
		//edgeTargetID, _ := strconv.ParseInt(edge.TargetID, 10, 64)
		// 根据指向关系重建树
		nodeMap[edge.SourceID].Next[edge.TargetID] = nodeMap[edge.TargetID]

		// 将尾节点删除，剩余的只有头节点的数据就是开始节点
		delete(headMap, edge.TargetID)
	}
	// globalTagGraphNode初始节点赋值
	for k, _ := range headMap {
		tagGraphNode.Next[k] = nodeMap[k]
	}

	// buildTaggedInformationV3
	buildTaggedInformationV3(tagGraphNode.Path, tagGraphNode)

	// 补充额外信息
	buildUnTaggedInformation(tagGraphNode)

	// 补充children信息
	buildChildrenInformation(tagGraphNode)

	// 保存到全局变量
	if globalTemplateGraphMap == nil {
		globalTemplateGraphMap = make(map[int64]*TagGraphNode)
	}
	globalTemplateGraphMap[template.ID] = tagGraphNode

	return tagGraphNode
}

func buildTaggedInformationV3(nodePath []int64, node *TagGraphNode) {
	for _, tag := range node.Next {
		node.Next[tag.ID].Path = append(nodePath, tag.ID)
		node.Next[tag.ID].RelatedHosts = HostService.HostsHavingTagIDs(node.Next[tag.ID].Path)
		node.Next[tag.ID].RelatedHostsCount = len(node.Next[tag.ID].RelatedHosts)
		node.Next[tag.ID].RelatedPods = CaasService.PodsHavingTagIDs(node.Next[tag.ID].Path)
		node.Next[tag.ID].RelatedPodsCount = len(node.Next[tag.ID].RelatedPods)
		buildTaggedInformationV3(node.Next[tag.ID].Path, tag)
	}
}

func newTemplateService() *templateService {
	return &templateService{}
}
