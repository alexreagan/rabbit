package app

import "github.com/alexreagan/rabbit/server/model/gtime"

type Template struct {
	ID       int64       `json:"id" gorm:"primary_key;column:id"`
	Name     string      `json:"name" gorm:"column:name;unique;type:string;size:128;comment:tag类别名称"`
	Remark   string      `json:"remark" gorm:"column:remark;type:string;size:1024;comment:tag类别描述"`
	State    string      `json:"state" gorm:"column:state;type:enum('enable','disable');default:disable;comment:tag类别描述"`
	Creator  string      `json:"creator" gorm:"column:creator;type:string;size:64;comment:创建人"`
	CreateAt gtime.GTime `json:"createAt" gorm:"column:create_at;comment:创建时间"`
	UpdateAt gtime.GTime `json:"updateAt" gorm:"column:update_at;comment:创建时间"`
}

func (t Template) TableName() string {
	return "template"
}
