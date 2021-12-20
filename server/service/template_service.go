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
	if dt := tx.Model(app.Template{}).Where("state = 'enable'").Find(&template); dt.Error != nil {
		return &template, dt.Error
	}
	return &template, nil
}

// 更新
func (t *templateService) Updates(template *app.Template) error {
	db := g.Con().Portal.Debug()
	if db = db.Model(app.Template{}).Where("id = ?", template.ID).Updates(template); db.Error != nil {
		return db.Error
	}
	return nil
}

var globalTagTreeNode *TagGraphNode

func (s *templateService) GlobalTagTreeNode() *TagGraphNode {
	return globalTagTreeNode
}

// 创建树
func (t *templateService) BuildTree(g6Graph *app.G6Graph) *TagGraphNode {
	headMap := make(map[int64]*TagGraphNode)
	nodeMap := make(map[int64]*TagGraphNode)

	// tag路由图
	globalTagGraphNode = newTagRouterGraphNode(&app.Tag{})

	// 初始化
	for _, n := range g6Graph.Nodes {
		nSize, _ := json.Marshal(n.Size)
		nInPoints, _ := json.Marshal(n.InPoints)
		nOutPoints, _ := json.Marshal(n.OutPoints)
		tag := &app.Tag{
			ID:           n.ID,
			Name:         n.Name,
			CnName:       n.Label,
			Label:        n.Label,
			Size:         string(nSize),
			Type:         n.Type,
			Color:        n.Color,
			Shape:        n.Shape,
			Image:        n.Image,
			StateImage:   n.StateImage,
			X:            n.X,
			Y:            n.Y,
			InPoints:     string(nInPoints),
			OutPoints:    string(nOutPoints),
			IsDoingStart: false,
			IsDoingEnd:   false,
		}
		nodeMap[n.ID] = newTagRouterGraphNode(tag)
		headMap[n.ID] = newTagRouterGraphNode(tag)
	}
	// 组织树状结构
	for _, edge := range g6Graph.Edges {
		// 根据指向关系重建树
		nodeMap[edge.SourceID].Next[edge.TargetID] = nodeMap[edge.TargetID]

		// 将尾节点删除，剩余的只有头节点的数据就是开始节点
		delete(headMap, edge.TargetID)
	}
	for k, _ := range headMap {
		globalTagGraphNode.Next[k] = nodeMap[k]
	}

	// 补充额外信息
	calUnTaggedInformation(globalTagGraphNode)

	return globalTagGraphNode
}

func newTemplateService() *templateService {
	return &templateService{}
}
