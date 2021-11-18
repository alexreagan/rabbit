package uic

import "github.com/alexreagan/rabbit/server/model/gtime"

type Role struct {
	ID        int64       `json:"id" gorm:"primary_key;column:id"`
	Name      string      `json:"name" gorm:"column:name;type:string;size:80;index;not null;comment:角色名称"`
	CnName    string      `json:"cnName" gorm:"column:cn_name;type:string;size:80;comment:中文名称"`
	Remark    string      `json:"remark" gorm:"column:remark;type:string;size:255;comment:描述"`
	CreatedAt gtime.GTime `json:"createdAt" gorm:"column:created_at;default:null;comment:"`
}

func (r Role) TableName() string {
	return "role"
}
