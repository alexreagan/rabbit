package node

import (
	"encoding/json"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/alert"
	"github.com/alexreagan/rabbit/server/model/caas"
	"strings"
)

var globalTree []*HostGroup
var globalNodeMap map[int64]*HostGroup
var GroupPathSeperator = "/"

type HostGroup struct {
	ID                int64        `json:"id" gorm:"primary_key;column:id"`
	Type              string       `json:"type" gorm:"column:type;type:enum('vmGroup','containerGroup');default:'vmGroup';comment:"`
	Name              string       `json:"name" gorm:"column:name;type:string;size:128;comment:"`
	ParentName        string       `json:"parentName" gorm:"column:parent_name;type:string;size:128;comment:"`
	ParentId          int64        `json:"parentId" gorm:"column:parent_id;comment:"`
	Path              string       `json:"path" gorm:"column:path;type:string;size:512;comment:"`
	PathArray         string       `json:"pathArray" gorm:"column:path_array;type:json;comment:"`
	CaasServiceId     int64        `json:"caasServiceId" gorm:"column:caas_service_id;comment:"`
	Desc              string       `json:"desc" gorm:"column:desc;type:string;size:256;comment:"`
	CreateUser        string       `json:"createUser" gorm:"column:create_user;type:string;size:32;comment:"`
	Children          []*HostGroup `json:"children" gorm:"-"`
	SubGroupCount     int          `json:"subGroupCount" gorm:"-"`
	RelatedHostCount  int          `json:"relatedHostCount" gorm:"-"`
	RelatedPodCount   int          `json:"relatedPodCount" gorm:"-"`
	ChildrenHostCount int          `json:"childrenHostCount" gorm:"-"`
	ChildrenPodCount  int          `json:"childrenPodCount" gorm:"-"`
	IsWarning         bool         `json:"isWarning" gorm:"-"`
}

func (this HostGroup) TableName() string {
	return "host_group"
}

//type HostGroupPro struct {
//	HostGroup
//	Children          []*HostGroupPro `json:"children"`
//	SubGroupCount     int             `json:"subGroupCount"`
//	RelatedHostsCount  int             `json:"relatedHostCount"`
//	RelatedPodCount   int             `json:"relatedPodCount"`
//	ChildrenHostCount int             `json:"childrenHostCount"`
//	ChildrenPodCount  int             `json:"childrenPodCount"`
//	IsWarning         bool            `json:"isWarning"`
//}

func (this *HostGroup) UpdateChildrenPath() {
	groupPathArray := this.GetPath()
	children := this.GetRTChildren()
	for _, child := range children {
		tGroupPath := groupPathArray
		tGroupPath = append(tGroupPath, child.Name)
		groupPathArrayBytes, _ := json.Marshal(tGroupPath)

		db := g.Con().Portal.Model(HostGroup{})
		db = db.Where("id = ?", child.ID).Updates(
			HostGroup{
				Path:      strings.Join(tGroupPath, GroupPathSeperator),
				PathArray: string(groupPathArrayBytes),
			})
		child.UpdateChildrenPath()
	}
	return
}

func (this HostGroup) GetParentName() string {
	var hostGroup HostGroup
	db := g.Con().Portal.Model(HostGroup{})
	db = db.Where("id = ?", this.ParentId)
	db.Find(&hostGroup)
	return hostGroup.Name
}

func (this HostGroup) BuildTree(id int64) ([]*HostGroup, map[int64]*HostGroup) {
	if globalTree != nil {
		return globalTree, globalNodeMap
	}
	var hostGroups []*HostGroup
	globalNodeMap = make(map[int64]*HostGroup)

	db := g.Con().Portal.Debug().Model(HostGroup{})
	if id != 0 {
		db = db.Where("id = ? or parent_id = ?", id, id)
	}
	db = db.Order("name")
	db.Find(&hostGroups)

	// 组建树状结构
	for _, grp := range hostGroups {
		// 群组默认为叶子节点，没达到报警条件
		grp.IsWarning = false
		globalNodeMap[grp.ID] = grp
	}
	for _, grp := range hostGroups {
		if grp.ParentId == 0 {
			globalTree = append(globalTree, grp)
		} else if _, ok := globalNodeMap[grp.ParentId]; ok {
			globalNodeMap[grp.ParentId].Children = append(globalNodeMap[grp.ParentId].Children, grp)

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
					globalNodeMap[t.ParentId].IsWarning = true
					t = globalNodeMap[t.ParentId]
				}
			}
		}
	}

	for _, grp := range globalNodeMap {
		grp.SubGroupCount = len(grp.Children)
		grp.RelatedHostCount = len(HostGroup{ID: grp.ID}.RelatedHosts())
		grp.RelatedPodCount = len(HostGroup{ID: grp.ID, CaasServiceId: grp.CaasServiceId}.RelatedPods())

		// 叶子节点加权
		if len(grp.Children) == 0 {
			grp.ChildrenHostCount = grp.RelatedHostCount
			grp.ChildrenPodCount = grp.RelatedPodCount

			t := grp
			for {
				if t.ParentId == 0 {
					break
				}
				globalNodeMap[t.ParentId].ChildrenHostCount += grp.ChildrenHostCount
				globalNodeMap[t.ParentId].ChildrenPodCount += grp.ChildrenPodCount

				t = globalNodeMap[t.ParentId]
			}
		}
	}
	return globalTree, globalNodeMap
}

func (this HostGroup) ReBuildTree() ([]*HostGroup, map[int64]*HostGroup) {
	globalTree = nil
	globalNodeMap = nil
	return this.BuildTree(0)
}

func (this HostGroup) GetPath() []string {
	var pathArray []string
	id := this.ID
	for {
		hostGroup := &HostGroup{}
		db := g.Con().Portal.Model(HostGroup{})
		db.Where("id = ?", id).Find(&hostGroup)
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

func (this HostGroup) GetRTChildren() []*HostGroup {
	var hostGroups []*HostGroup
	dt := g.Con().Portal.Model(HostGroup{})
	dt = dt.Where("parent_id = ?", this.ID).Find(&hostGroups)
	return hostGroups
}

func (this HostGroup) GetChildren() []*HostGroup {
	if globalNodeMap == nil {
		this.BuildTree(0)
	}
	return globalNodeMap[this.ID].Children
}

func (this HostGroup) RelatedHosts() []*Host {
	//var hostGroupRels []*HostGroupRel
	//g.Con().Portal.Model(HostGroupRel{}).Where("`group_id` = ?", this.ID).Find(&hostGroupRels)
	//var hostIds []int64
	//for _, t := range hostGroupRels {
	//	hostIds = append(hostIds, t.HostID)
	//}
	//var hosts []*Host
	//g.Con().Portal.Model(Host{}).Where("id in (?)", hostIds).Find(&hosts)

	//// 添加报警标识
	//for _, host := range hosts {
	//	host.AdditionalAttrs()
	//}

	// 报警信息
	alerts := alert.Alert{}.LatestRecords()
	alertMap := make(map[string]*alert.Alert)
	for _, alert := range alerts {
		alertMap[alert.ProdIP] = alert
	}

	// 当前组关联的机器
	var hosts []*Host
	db := g.Con().Portal.Debug()
	db = db.Model(Host{})
	db = db.Select("`host`.*")
	db = db.Joins("left join `host_group_rel` on `host_group_rel`.`host_id` = `host`.`id`")
	db = db.Where("`host_group_rel`.`group_id` = ?", this.ID)
	db = db.Find(&hosts)

	// 报警信息
	for _, h := range hosts {
		h.Type = "host"

		alt, ok := alertMap[h.IP]
		if ok {
			h.IsWarning = alt.Resolved == false
		} else {
			h.IsWarning = false
		}
	}
	return hosts
}

// 判断群组是否满足报警条件
func (this HostGroup) MeetWarningCondition() bool {
	hosts := this.RelatedHosts()
	for _, host := range hosts {
		if host.IsWarning == true {
			return true
		}
	}
	return false
}

func (this HostGroup) RelatedPods() []*caas.Pod {
	var pods []*caas.Pod
	tx := g.Con().Portal.Model(caas.Pod{}).Debug()
	tx = tx.Select("`caas_pod`.*")
	tx = tx.Joins("left join `caas_service_pod_rel` on `caas_pod`.`id` = `caas_service_pod_rel`.`pod`")
	tx = tx.Where("`caas_service_pod_rel`.`service` = ?", this.CaasServiceId)
	tx = tx.Find(&pods)

	// 添加报警标识
	for _, pod := range pods {
		pod.AdditionalAttrs()
	}
	return pods
}
