package sys

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type Param struct {
	ID       int64       `json:"menuId" gorm:"primary_key;column:id"`
	Key      string      `json:"key" gorm:"column:key;type:string;size:256;comment:参数key"`
	Value    string      `json:"value" gorm:"column:value;type:string;size:1024;comment:参数value"`
	Remark   string      `json:"remark" gorm:"column:remark;type:string;size:1024;comment:参数说明"`
	CreateAt gtime.GTime `json:"createAt" gorm:"column:create_at;default:null;comment:创建时间"`
	Deleted  bool        `json:"deleted" gorm:"column:deleted;comment:是否已删除"`
}

func (this Param) TableName() string {
	return "param"
}
