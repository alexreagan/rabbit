package node

import (
	"encoding/json"
	"rabbit/g"
	"rabbit/server/model/caas"
)

var tree []*HostGroup
var nodeMap map[int64]*HostGroup
var GroupPathSeperator = "/"

type HostGroup struct {
	ID            int64        `json:"id" gorm:"primary_key;column:id"`
	Type          string       `json:"type" gorm:"column:type;type:enum('vmGroup','containerGroup');default:'vmGroup';comment:"`
	Name          string       `json:"name" gorm:"column:name;type:string;size:128;comment:"`
	ParentName    string       `json:"parentName" gorm:"column:parent_name;type:string;size:128;comment:"`
	ParentId      int64        `json:"parentId" gorm:"column:parent_id;comment:"`
	Path          string       `json:"path" gorm:"column:path;type:string;size:512;comment:"`
	PathArray     string       `json:"pathArray" gorm:"column:path_array;type:json;comment:"`
	CaasServiceId int64        `json:"caasServiceId" gorm:"column:caas_service_id;comment:"`
	Desc          string       `json:"desc" gorm:"column:desc;type:string;size:256;comment:"`
	Icon          string       `json:"icon" gorm:"column:icon;type:string;size:512;comment:"`
	CreateUser    string       `json:"createUser" gorm:"column:create_user;type:string;size:32;comment:"`
	Children      []*HostGroup `json:"children" gorm:"-"`
	IsLeaf        bool         `json:"isLeaf" gorm:"-"`
	IsWarning     bool         `json:"isWarning" gorm:"-"`
}

func (this HostGroup) TableName() string {
	return "host_group"
}

func (this HostGroup) GetParentName() string {
	var hostGroup HostGroup
	db := g.Con().Portal.Table(HostGroup{}.TableName())
	db = db.Where("id = ?", this.ParentId)
	db.Find(&hostGroup)
	return hostGroup.Name
}

func (this HostGroup) BuildTree(id int64) ([]*HostGroup, map[int64]*HostGroup) {
	if tree != nil {
		return tree, nodeMap
	}
	var hostGroups []*HostGroup
	nodeMap = make(map[int64]*HostGroup)

	db := g.Con().Portal.Debug().Model(HostGroup{})
	if id != 0 {
		db = db.Where("id = ? or parent_id = ?", id, id)
	}
	db = db.Order("name")
	db.Find(&hostGroups)

	// 组建树状结构
	for _, grp := range hostGroups {
		// 群组默认为叶子节点，没达到报警条件
		grp.IsLeaf = true
		grp.IsWarning = false
		nodeMap[grp.ID] = grp
	}
	for _, grp := range hostGroups {
		if grp.ParentId == 0 {
			tree = append(tree, grp)
		} else if _, ok := nodeMap[grp.ParentId]; ok {
			nodeMap[grp.ParentId].Children = append(nodeMap[grp.ParentId].Children, grp)
			nodeMap[grp.ParentId].IsLeaf = len(nodeMap[grp.ParentId].Children) == 0

			// 设置报警状态
			if grp.MeetWarningCondition() == true {
				// 达到报警条件，当前节点设置为报警
				grp.IsWarning = true

				// 所有的父节点设置为报警
				t := grp
				for {
					if t.ParentId == 0 {
						break
					}
					nodeMap[t.ParentId].IsWarning = true
					t = nodeMap[t.ParentId]
				}
			}
		}
	}
	return tree, nodeMap
}

func (this HostGroup) ReBuildTree() ([]*HostGroup, map[int64]*HostGroup) {
	tree = nil
	nodeMap = nil
	return this.BuildTree(0)
}

func (this HostGroup) GetPath() []string {
	var pathArray []string
	id := this.ID
	for {
		hostGroup := &HostGroup{}
		db := g.Con().Portal.Table(HostGroup{}.TableName())
		db.Where("id=?", id).Find(&hostGroup)
		pathArray = append(pathArray, hostGroup.Name)

		if hostGroup.ParentId == 0 {
			break
		}
		id = hostGroup.ParentId
	}
	var reversePathArray []string
	for i := len(pathArray) - 1; i >= 0; i-- {
		reversePathArray = append(reversePathArray, pathArray[i])
	}
	return reversePathArray
}

func (this HostGroup) GetJsonPath() string {
	reversePath := this.GetPath()
	path, _ := json.Marshal(reversePath)
	return string(path)
}

func (this HostGroup) GetChildren() []*HostGroup {
	if nodeMap == nil {
		this.BuildTree(0)
	}
	return nodeMap[this.ID].Children
}

func (this HostGroup) RelatedHosts() []*Host {
	var hostGroupRels []*HostGroupRel
	g.Con().Portal.Table(HostGroupRel{}.TableName()).Where("`group_id` = ?", this.ID).Find(&hostGroupRels)
	var hostIds []int64
	for _, t := range hostGroupRels {
		hostIds = append(hostIds, t.HostID)
	}
	var hosts []*Host
	g.Con().Portal.Table(Host{}.TableName()).Where("id in (?)", hostIds).Find(&hosts)

	// 添加报警标识
	for _, host := range hosts {
		host.AdditionalAttrs()
	}
	return hosts
}

func (this HostGroup) MeetWarningCondition() bool {
	hosts := this.RelatedHosts()
	for _, host := range hosts {
		if host.MeetWarningCondition() == true {
			return true
		}
	}
	return false
}

func (this HostGroup) RelatedPods() []*caas.Pod {
	var pods []*caas.Pod
	tx := g.Con().Portal.Model(caas.Pod{}).Debug()
	tx = tx.Select("`caas_pod`.*")
	tx = tx.Joins("left join `caas_service_pod_rel` on `caas_pod`.`id` = `caas_service_pod_rel`.`pod_id`")
	tx = tx.Where("`caas_service_pod_rel`.`service_id` = ?", this.CaasServiceId)
	tx = tx.Find(&pods)

	// 添加报警标识
	for _, pod := range pods {
		pod.AdditionalAttrs()
	}
	return pods
}
