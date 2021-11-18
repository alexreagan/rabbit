package uic

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type Perm struct {
	ID       int64       `json:"id" gorm:"primary_key;column:id"`
	Name     string      `json:"name" gorm:"column:name;type:string;unique;size:80;index;not null;comment:权限名称"`
	CnName   string      `json:"cnName" gorm:"column:cn_name;type:string;size:80;comment:权限中文名称"`
	Remark   string      `json:"remark" gorm:"column:remark;type:string;size:255;comment:备注"`
	CreateAt gtime.GTime `json:"createAt" gorm:"column:create_at;default:null;comment:"`
}

func (r Perm) TableName() string {
	return "perm"
}
