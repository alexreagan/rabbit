package node

import (
	"encoding/json"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/alarm"
	"github.com/alexreagan/rabbit/server/model/caas"
	"strings"
)

var globalTree []*NodeGroup
var globalNodeMap map[int64]*NodeGroup
var GroupPathSeperator = "/"

type NodeGroup struct {
	ID                int64        `json:"id" gorm:"primary_key;column:id"`
	Type              string       `json:"type" gorm:"column:type;type:enum('vmGroup','containerGroup');default:'vmGroup';comment:"`
	Name              string       `json:"name" gorm:"column:name;type:string;size:128;comment:"`
	ParentName        string       `json:"parentName" gorm:"column:parent_name;type:string;size:128;comment:"`
	ParentID          int64        `json:"parentID" gorm:"column:parent_id;comment:"`
	Path              string       `json:"path" gorm:"column:path;type:string;size:512;comment:"`
	PathArray         string       `json:"pathArray" gorm:"column:path_array;type:json;comment:"`
	CaasServiceID     int64        `json:"caasServiceId" gorm:"column:caas_service_id;comment:"`
	Desc              string       `json:"desc" gorm:"column:desc;type:string;size:256;comment:"`
	CreateUser        string       `json:"createUser" gorm:"column:create_user;type:string;size:32;comment:"`
	Children          []*NodeGroup `json:"children" gorm:"-"`
	SubGroupCount     int          `json:"subGroupCount" gorm:"-"`
	RelatedNodeCount  int          `json:"relatedNodeCount" gorm:"-"`
	RelatedPodCount   int          `json:"relatedPodCount" gorm:"-"`
	ChildrenNodeCount int          `json:"childrenNodeCount" gorm:"-"`
	ChildrenPodCount  int          `json:"childrenPodCount" gorm:"-"`
	IsWarning         bool         `json:"isWarning" gorm:"-"`
}

func (this NodeGroup) TableName() string {
	return "node_group"
}

//type nodeGroupPro struct {
//	NodeGroup
//	Children          []*nodeGroupPro `json:"children"`
//	SubGroupCount     int             `json:"subGroupCount"`
//	RelatedNodesCount  int             `json:"relatedNodeCount"`
//	RelatedPodCount   int             `json:"relatedPodCount"`
//	ChildrenNodeCount int             `json:"childrenNodeCount"`
//	ChildrenPodCount  int             `json:"childrenPodCount"`
//	IsWarning         bool            `json:"isWarning"`
//}

func (this *NodeGroup) UpdateChildrenPath() {
	groupPathArray := this.GetPath()
	children := this.GetRTChildren()
	for _, child := range children {
		tGroupPath := groupPathArray
		tGroupPath = append(tGroupPath, child.Name)
		groupPathArrayBytes, _ := json.Marshal(tGroupPath)

		tx := g.Con().Portal.Model(NodeGroup{})
		tx = tx.Where("id = ?", child.ID).Updates(
			&NodeGroup{
				Path:      strings.Join(tGroupPath, GroupPathSeperator),
				PathArray: string(groupPathArrayBytes),
			})
		child.UpdateChildrenPath()
	}
	return
}

func (this NodeGroup) GetParentName() string {
	var nodeGroup NodeGroup
	tx := g.Con().Portal.Model(NodeGroup{})
	tx = tx.Where("id = ?", this.ParentID)
	tx.Find(&nodeGroup)
	return nodeGroup.Name
}

func (this NodeGroup) BuildTree(id int64) ([]*NodeGroup, map[int64]*NodeGroup) {
	if globalTree != nil {
		return globalTree, globalNodeMap
	}
	var nodeGroups []*NodeGroup
	globalNodeMap = make(map[int64]*NodeGroup)

	tx := g.Con().Portal.Model(NodeGroup{})
	if id != 0 {
		tx = tx.Where("id = ? or parent_id = ?", id, id)
	}
	tx = tx.Order("name")
	tx.Find(&nodeGroups)

	// 组建树状结构
	for _, grp := range nodeGroups {
		// 群组默认为叶子节点，没达到报警条件
		grp.IsWarning = false
		globalNodeMap[grp.ID] = grp
	}
	for _, grp := range nodeGroups {
		if grp.ParentID == 0 {
			globalTree = append(globalTree, grp)
		} else if _, ok := globalNodeMap[grp.ParentID]; ok {
			globalNodeMap[grp.ParentID].Children = append(globalNodeMap[grp.ParentID].Children, grp)

			// 设置报警状态
			if grp.MeetWarningCondition() == true {
				// 达到报警条件，当前节点设置为报警
				grp.IsWarning = true

				// 所有的父节点设置为报警
				t := grp
				for {
					if t.ParentID == 0 {
						break
					}
					globalNodeMap[t.ParentID].IsWarning = true
					t = globalNodeMap[t.ParentID]
				}
			}
		}
	}

	for _, grp := range globalNodeMap {
		grp.SubGroupCount = len(grp.Children)
		grp.RelatedNodeCount = len(NodeGroup{ID: grp.ID}.RelatedNodes())
		grp.RelatedPodCount = len(NodeGroup{ID: grp.ID, CaasServiceID: grp.CaasServiceID}.RelatedPods())

		// 叶子节点加权
		if len(grp.Children) == 0 {
			grp.ChildrenNodeCount = grp.RelatedNodeCount
			grp.ChildrenPodCount = grp.RelatedPodCount

			t := grp
			for {
				if t.ParentID == 0 {
					break
				}
				globalNodeMap[t.ParentID].ChildrenNodeCount += grp.ChildrenNodeCount
				globalNodeMap[t.ParentID].ChildrenPodCount += grp.ChildrenPodCount

				t = globalNodeMap[t.ParentID]
			}
		}
	}
	return globalTree, globalNodeMap
}

func (this NodeGroup) ReBuildTree() ([]*NodeGroup, map[int64]*NodeGroup) {
	globalTree = nil
	globalNodeMap = nil
	return this.BuildTree(0)
}

func (this NodeGroup) GetPath() []string {
	var pathArray []string
	id := this.ID
	for {
		nodeGroup := &NodeGroup{}
		tx := g.Con().Portal.Model(NodeGroup{})
		tx.Where("id = ?", id).Find(&nodeGroup)
		pathArray = append(pathArray, nodeGroup.Name)

		if nodeGroup.ParentID == 0 {
			break
		}
		id = nodeGroup.ParentID
	}
	var reversePathArray []string
	for i := len(pathArray) - 1; i >= 0; i-- {
		reversePathArray = append(reversePathArray, pathArray[i])
	}
	return reversePathArray
}

func (this NodeGroup) GetJsonPath() string {
	reversePath := this.GetPath()
	path, _ := json.Marshal(reversePath)
	return string(path)
}

func (this NodeGroup) GetRTChildren() []*NodeGroup {
	var nodeGroups []*NodeGroup
	tx := g.Con().Portal.Model(NodeGroup{})
	tx = tx.Where("parent_id = ?", this.ID).Find(&nodeGroups)
	return nodeGroups
}

func (this NodeGroup) GetChildren() []*NodeGroup {
	if globalNodeMap == nil {
		this.BuildTree(0)
	}
	return globalNodeMap[this.ID].Children
}

func (this NodeGroup) RelatedNodes() []*Node {
	// 报警信息
	alarms := alarm.Alarm{}.LatestRecords()
	alarmMap := make(map[string]*alarm.Alarm)
	for _, alm := range alarms {
		alarmMap[alm.ProdIP] = alm
	}

	// 当前组关联的机器
	var nodes []*Node
	tx := g.Con().Portal.Model(Node{})
	tx = tx.Select("`node`.*")
	tx = tx.Joins("left join `node_group_rel` on `node_group_rel`.`node_id` = `node`.`id`")
	tx = tx.Where("`node_group_rel`.`group_id` = ?", this.ID)
	tx = tx.Find(&nodes)

	// 报警信息
	for _, n := range nodes {
		n.Type = "node"

		alt, ok := alarmMap[n.IP]
		if ok {
			n.IsWarning = alt.Resolved == false
		} else {
			n.IsWarning = false
		}
	}
	return nodes
}

// 判断群组是否满足报警条件
func (this NodeGroup) MeetWarningCondition() bool {
	nodes := this.RelatedNodes()
	for _, n := range nodes {
		if n.IsWarning == true {
			return true
		}
	}
	return false
}

func (this NodeGroup) RelatedPods() []*caas.Pod {
	var pods []*caas.Pod
	tx := g.Con().Portal.Model(caas.Pod{}).Debug()
	tx = tx.Select("`caas_pod`.*")
	tx = tx.Joins("left join `caas_service_pod_rel` on `caas_pod`.`id` = `caas_service_pod_rel`.`pod`")
	tx = tx.Where("`caas_service_pod_rel`.`service` = ?", this.CaasServiceID)
	tx = tx.Find(&pods)

	// 添加报警标识
	for _, pod := range pods {
		pod.AdditionalAttrs()
	}
	return pods
}
