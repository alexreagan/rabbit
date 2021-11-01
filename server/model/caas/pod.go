package caas

import "time"

type Pod struct {
	ID          int64     `json:"id" gorm:"primary_key;column:id"`
	HostIP      string    `json:"hostIp" gorm:"column:host_ip;type:string;size:128;comment:"`
	HostName    string    `json:"hostName" gorm:"column:host_name;type:string;size:128;comment:"`
	PodIP       string    `json:"podIp" gorm:"column:pod_ip;type:string;size:128;comment:"`
	Name        string    `json:"name" gorm:"column:name;type:string;size:128;comment:"`
	NameSpace   string    `json:"namespace" gorm:"column:namespace;type:string;size:128;comment:"`
	ServiceName string    `json:"serviceName" gorm:"column:service_name;type:string;size:128;comment:"`
	Status      string    `json:"status" gorm:"column:status;type:string;size:128;comment:"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;default:null;comment:"`
	UpdateTime  time.Time `json:"updateTime" gorm:"column:update_time;default:null;comment:"`
	Type        string    `json:"type" gorm:"-"`
	IsWarning   bool      `json:"isWarning" gorm:"-"`
}

func (this Pod) TableName() string {
	return "caas_pod"
}

func (this Pod) MeetWarningCondition() bool {
	return this.Status != "Running"
}

func (this *Pod) AdditionalAttrs() *Pod {
	this.Type = "pod"
	this.IsWarning = this.MeetWarningCondition()
	return this
}
