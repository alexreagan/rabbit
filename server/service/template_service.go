package service

import (
	"encoding/json"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/app"
)

type templateService struct{}

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

var globalTagGraphNodeV3 *TagGraphNode

func (s *templateService) GlobalTagGraphNodeV3() *TagGraphNode {
	return globalTagGraphNodeV3
}

// 创建树
func (t *templateService) BuildGraphV3(g6Graph *app.G6Graph) *TagGraphNode {
	headMap := make(map[int64]*TagGraphNode)
	nodeMap := make(map[int64]*TagGraphNode)

	// tag路由图
	globalTagGraphNodeV3 = newTagGraphNode(&app.Tag{})

	// 初始化
	for _, n := range g6Graph.Nodes {
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
		// 根据指向关系重建树
		nodeMap[edge.SourceID].Next[edge.TargetID] = nodeMap[edge.TargetID]

		// 将尾节点删除，剩余的只有头节点的数据就是开始节点
		delete(headMap, edge.TargetID)
	}
	// globalTagGraphNode初始节点赋值
	for k, _ := range headMap {
		globalTagGraphNodeV3.Next[k] = nodeMap[k]
	}

	// buildTaggedInformationV3
	buildTaggedInformationV3(globalTagGraphNodeV3.Path, globalTagGraphNodeV3)

	// 补充额外信息
	buildUnTaggedInformation(globalTagGraphNodeV3)

	// 补充children信息
	buildChildrenInformation(globalTagGraphNodeV3)

	return globalTagGraphNodeV3
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
