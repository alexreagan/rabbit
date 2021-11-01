package caas

import "time"

type WorkSpace struct {
	ID          int64     `json:"id" gorm:"primary_key;column:id"`
	Name        string    `json:"name" gorm:"column:name;type:string;size:128;comment:"`
	OrgName     string    `json:"orgName" gorm:"column:org_name;type:string;size:128;comment:"`
	OrgProposer string    `json:"orgProposer" gorm:"column:org_proposer;type:string;size:128;comment:"`
	Description string    `json:"description" gorm:"column:description;type:string;size:128;comment:"`
	Deleted     int64     `json:"deleted" gorm:"column:deleted;comment:"`
	Zones       string    `json:"zones" gorm:"column:zones;type:string;size:128;comment:"`
	NsCount     int64     `json:"nsCount" gorm:"column:ns_count;comment:"`
	UpdateTime  time.Time `json:"updateTime" gorm:"column:update_time;default:null;comment:"`
}

func (this WorkSpace) TableName() string {
	return "caas_workspace"
}
