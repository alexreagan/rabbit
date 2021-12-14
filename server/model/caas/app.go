package caas

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type App struct {
	ID          int64       `json:"id" gorm:"primary_key;column:id"`
	AppName     string      `json:"appName" gorm:"column:app_name;type:string;size:128;comment:名称"`
	NameSpaceId int64       `json:"namespaceID" gorm:"column:namespace_id;comment:项目空间ID"`
	Description string      `json:"description" gorm:"column:description;type:string;size:128;comment:描述"`
	CreateTime  gtime.GTime `json:"createTime" gorm:"column:create_time;default:null;comment:创建时间"`
	UpdateTime  gtime.GTime `json:"updateTime" gorm:"column:update_time;default:null;comment:数据更新时间"`
}

func (this App) TableName() string {
	return "caas_app"
}
