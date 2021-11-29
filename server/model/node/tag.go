package node

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type TagCategory struct {
	ID       int64       `json:"id" gorm:"primary_key;column:id"`
	Name     string      `json:"name" gorm:"column:name;unique;type:string;size:128;comment:tag类别名称"`
	CnName   string      `json:"cnName" gorm:"column:cn_name;unique;type:string;size:128;comment:tag中文名称"`
	Remark   string      `json:"remark" gorm:"column:remark;type:string;size:1024;comment:tag类别描述"`
	CreateAt gtime.GTime `json:"createAt" gorm:"column:create_at;comment:创建时间"`
}

func (t TagCategory) TableName() string {
	return "tag_category"
}

type Tag struct {
	ID           int64       `json:"id" gorm:"primary_key;column:id"`
	Name         string      `json:"name" gorm:"column:name;unique;type:string;size:128;comment:tag名称"`
	CnName       string      `json:"cnName" gorm:"column:cn_name;unique;type:string;size:128;comment:tag中文名称"`
	CategoryID   int64       `json:"categoryID" gorm:"column:category_id;comment:所属类别id"`
	Remark       string      `json:"remark" gorm:"column:remark;type:string;size:1024;comment:tag描述"`
	CreateAt     gtime.GTime `json:"createAt" gorm:"column:create_at;comment:创建时间"`
	CategoryName string      `json:"categoryName"`
}

func (t Tag) TableName() string {
	return "tag"
}

type Tags []*Tag

func (t Tags) Len() int { return len(t) }

func (t Tags) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
