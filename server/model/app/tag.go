package app

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
	"sort"
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
	Name         string      `json:"name" gorm:"column:name;type:string;size:128;comment:名称"`
	CnName       string      `json:"cnName" gorm:"column:cn_name;unique;type:string;size:128;comment:中文名称"`
	CategoryID   int64       `json:"categoryID" gorm:"column:category_id;comment:所属类别id"`
	Remark       string      `json:"remark" gorm:"column:remark;type:string;size:1024;comment:tag描述"`
	CreateAt     gtime.GTime `json:"createAt" gorm:"column:create_at;comment:创建时间"`
	CategoryName string      `json:"categoryName"`
	Label        string      `json:"label" gorm:"column:label;type:string;size:128;comment:画布上的展现文字"`
	Size         string      `json:"size" gorm:"column:size;type:string;size:64;comment:大小，譬如170*34"`
	Type         string      `json:"type" gorm:"column:type;type:string;size:64;default:node;comment:类型"`
	Color        string      `json:"color" gorm:"column:color;type:string;size:64;default:#1890ff;comment:颜色"`
	Shape        string      `json:"shape" gorm:"column:shape;type:string;size:256;default:customNode;comment:形状"`
	Image        string      `json:"image" gorm:"column:image;type:varchar(512);comment:背景图片url地址"`
	StateImage   string      `json:"stateImage" gorm:"column:state_image;type:varchar(512);comment:背景图片url地址"`
	X            int         `json:"x" gorm:"column:x;default:0;comment:"`
	Y            int         `json:"y" gorm:"column:y;default:0;comment:"`
	InPoints     string      `json:"inPoints" gorm:"column:in_points;comment:"`
	OutPoints    string      `json:"outPoints" gorm:"column:out_points;comment:"`
	IsDoingStart bool        `json:"isDoingStart" gorm:"column:is_doing_start;default:false;comment:是否是开始节点"`
	IsDoingEnd   bool        `json:"isDoingEnd" gorm:"column:is_doing_end;default:false;comment:是否是结束节点"`
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

func (t Tags) Sort() {
	sort.Sort(t)
}
