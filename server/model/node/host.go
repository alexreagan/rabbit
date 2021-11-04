package node

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/spf13/viper"
	"time"
)

type Host struct {
	//gorm.Model
	ID                   int64        `json:"id" gorm:"primary_key;column:id"`
	Name                 string       `json:"name" gorm:"column:name;type:string;size:128;comment:"`
	IP                   string       `json:"ip" gorm:"column:ip;type:string;size:128;comment:IP"`
	PhysicalSystem       string       `json:"physicalSystem" gorm:"column:physical_system;type:string;size:128;comment:所属物理子系统"`
	PhysicalSystemArea   string       `json:"physicalSystemArea" gorm:"column:physical_system_area;type:string;size:128;comment:所属物理子系统"`
	PhysicalSystemCnName string       `json:"physicalSystemCnName" gorm:"column:physical_system_cn_name;type:string;size:128;"`
	PhysicalSystemEnName string       `json:"physicalSystemEnName" gorm:"column:physical_system_en_name;type:string;size:128;"`
	LogicSystem          string       `json:"logicSystem" gorm:"column:logic_system;type:string;size:128;"`
	LogicSystemCnName    string       `json:"logicSystemCnName" gorm:"column:logic_system_cn_name;type:string;size:128;"`
	Department           string       `json:"department" gorm:"column:department;type:string;size:128;comment:所属团队"`
	ApplyUser            string       `json:"applyUser" gorm:"column:apply_user;type:string;size:128;comment:"`
	AreaName             string       `json:"areaName" gorm:"column:area_name;type:string;size:128;comment:"`
	CpuNumber            int          `json:"cpuNumber" gorm:"column:cpu_number;type:int"`
	DeployDate           time.Time    `json:"deployDate" gorm:"column:deploy_date"`
	DevAreaCode          string       `json:"devAreaCode" gorm:"column:dev_area_code;type:string;size:128;comment:"`
	FunDesc              string       `json:"funDesc" gorm:"column:fun_desc;type:string;size:512;comment:"`
	InstanceId           string       `json:"instanceId" gorm:"column:instance_id;type:string;size:128;comment:"`
	ManagerA             string       `json:"managerA" gorm:"column:manager_a;type:string;size:128;comment:"`
	ManagerB             string       `json:"managerB" gorm:"column:manager_b;type:string;size:128;comment:"`
	MemorySize           int          `json:"memorySize" gorm:"column:memory_size"`
	OsVersion            string       `json:"osVersion" gorm:"column:os_version;type:string;size:128;comment:"`
	PartTypeCode         string       `json:"partTypeCode" gorm:"column:part_type_code;type:string;size:128;comment:"`
	ManIp                string       `json:"manIp" gorm:"column:man_ip;type:string;size:128;comment:"`
	CloudPoolCode        string       `json:"cloudPoolCode" gorm:"column:cloud_pool_code;type:string;size:128;comment:"`
	CloudPoolName        string       `json:"cloudPoolName" gorm:"column:cloud_pool_name;type:string;size:128;comment:"`
	CoreTotalNum         int          `json:"coreTotalNum" gorm:"column:core_total_num;comment:"`
	DatabaseVersion      string       `json:"databaseVersion" gorm:"column:database_version;type:string;size:128;comment:"`
	DevCenterName        string       `json:"devCenterName" gorm:"column:dev_center_name;type:string;size:128;comment:"`
	DevTypeCode          string       `json:"devTypeCode" gorm:"column:dev_type_code;type:string;size:128;comment:"`
	CpuUsage             float64      `json:"cpuUsage" gorm:"column:cpu_usage;comment:"`
	MemoryUsage          float64      `json:"memoryUsage" gorm:"column:memory_usage;comment:"`
	FsUsage              float64      `json:"fsUsage" gorm:"column:fs_usage;comment:"`
	Oracle               string       `json:"oracle" gorm:"column:oracle;type:string;size:128;comment:"`
	PartType             string       `json:"partType" gorm:"column:part_type;type:string;size:128;comment:"`
	ServSpaceCodeList    string       `json:"servSpaceCodeList" gorm:"column:serv_space_code_list;type:string;size:128;comment:"`
	ServSpaceNameList    string       `json:"servSpaceNameList" gorm:"column:serv_space_name_list;type:string;size:128;comment:"`
	SetupMode            string       `json:"setupMode" gorm:"column:setup_mode;type:string;size:32;comment:"`
	SrvStatus            string       `json:"srvStatus" gorm:"column:srv_status;type:string;size:32;comment:"`
	Status               string       `json:"status" gorm:"column:status;type:string;size:32;comment:"`
	VirtFcNum            string       `json:"virtFcNum" gorm:"column:virt_fc_num;type:string;size:128;comment:"`
	VirtNetNum           string       `json:"virtNetNum" gorm:"column:virt_net_num;type:string;size:128;comment:"`
	DevOwner             string       `json:"devOwner" gorm:"column:dev_owner;type:string;size:128;comment:"`
	Desc                 string       `json:"desc" gorm:"column:desc;type:string;size:256;comment:"`
	CreateTime           time.Time    `json:"createTime" gorm:"column:create_time;default:null;comment:"`
	UpdateTime           time.Time    `json:"updateTime" gorm:"column:update_time;default:null;comment:"`
	Groups               []*HostGroup `json:"groups" gorm:"-"`
	IsWarning            bool         `json:"isWarning" gorm:"-"`
	Type                 string       `json:"type" gorm:"-"`
}

func (this Host) TableName() string {
	return "host"
}

func (this Host) Existing() (int64, bool) {
	db := g.Con().Portal
	db.Table(this.TableName()).Where("name = ?", this.Name).Scan(&this)
	if this.ID != 0 {
		return this.ID, true
	} else {
		return 0, false
	}
}

func (this Host) RelatedGroups() []*HostGroup {
	var hostGroupRels []*HostGroupRel
	g.Con().Portal.Table(HostGroupRel{}.TableName()).Where("`host_id` = ?", this.ID).Find(&hostGroupRels)
	var groupIds []int64
	for _, t := range hostGroupRels {
		groupIds = append(groupIds, t.GroupID)
	}
	var hostGroups []*HostGroup
	g.Con().Portal.Table(HostGroup{}.TableName()).Where("id in (?)", groupIds).Find(&hostGroups)
	return hostGroups
}

func (this *Host) AdditionalAttrs() *Host {
	this.Type = "host"
	this.IsWarning = this.MeetWarningCondition()
	return this
}

func (this Host) MeetWarningCondition() bool {
	if this.CpuUsage >= viper.GetFloat64("alarm.threshold.cpu") {
		return true
	}
	if this.FsUsage >= viper.GetFloat64("alarm.threshold.fs") {
		return true
	}
	return false
}
