package caas

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

// 项目空间
type NameSpace struct {
	ID                 int64       `json:"id" gorm:"primary_key;column:id"`
	Namespace          string      `json:"namespace" gorm:"column:namespace;type:string;size:128;comment:"`
	WorkSpaceID        int64       `json:"workspaceId" gorm:"column:workspace_id;comment:工作空间ID"`
	WorkSpaceName      string      `json:"workspaceName" gorm:"column:workspace_name;type:string;size:128;comment:"`
	ClusterID          int64       `json:"clusterId" gorm:"column:cluster_id;comment:"`
	ClusterName        string      `json:"clusterName" gorm:"column:cluster_name;type:string;size:128;comment:"`
	PhysicalSystemID   int64       `json:"phSubSystemId" gorm:"column:physical_system;comment:"`
	PhysicalSystemName string      `json:"phSubSystemName" gorm:"column:physical_system_name;type:string;size:128;comment:"`
	MetaData           string      `json:"metaData" gorm:"column:meta_data;type:string;size:128;comment:"`
	Cpu                int64       `json:"cpu" gorm:"column:cpu;comment:"`
	Gpu                int64       `json:"gpu" gorm:"column:gpu;comment:"`
	Memory             int64       `json:"memory" gorm:"column:memory;comment:"`
	SharedVolume       int64       `json:"sharedVolume" gorm:"column:shared_volume;comment:"`
	LocalVolume        int64       `json:"localVolume" gorm:"column:local_volume;comment:"`
	Zones              string      `json:"zones" gorm:"column:zones;type:string;size:128;comment:"`
	CreateTime         gtime.GTime `json:"createTime" gorm:"column:create_time;default:null;comment:"`
	UpdateTime         gtime.GTime `json:"updateTime" gorm:"column:update_time;default:null;comment:"`
}

func (this NameSpace) TableName() string {
	return "caas_namespace"
}
