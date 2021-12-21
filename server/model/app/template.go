package app

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type Template struct {
	ID       int64       `json:"id" gorm:"primary_key;column:id"`
	Name     string      `json:"name" gorm:"column:name;unique;type:string;size:128;comment:模板名称"`
	Type     string      `json:"type" gorm:"column:type;type:enum('tree');default:'tree';comment:模板类型"`
	Remark   string      `json:"remark" gorm:"column:remark;type:text;comment:模板描述"`
	State    string      `json:"state" gorm:"column:state;type:enum('enable','disable');default:disable;comment:模板状态，同一时刻最多只有一个template处于启用状态"`
	Content  string      `json:"content" gorm:"column:content;type:text;comment:模板配置"`
	Creator  string      `json:"creator" gorm:"column:creator;type:string;size:64;comment:创建人"`
	CreateAt gtime.GTime `json:"createAt" gorm:"column:create_at;comment:创建时间"`
	UpdateAt gtime.GTime `json:"updateAt" gorm:"column:update_at;comment:更新时间"`
}

func (t Template) TableName() string {
	return "template"
}

type G6Point struct {
	X float64 `json:"x" form:"x" gorm:"x"`
	Y float64 `json:"y" form:"y" gorm:"y"`
}

type G6Edge struct {
	End        G6Point `json:"end" form:"end"`
	EndPoint   G6Point `json:"endPoint" form:"endPoint"`
	ID         string  `json:"id" form:"id"`
	Shape      string  `json:"shape" form:"shape"`
	Source     int64   `json:"source" form:"source"`
	SourceID   int64   `json:"sourceId" form:"sourceId"` // edge.sourceID 指向 node.id
	Start      G6Point `json:"start" form:"start"`
	StartPoint G6Point `json:"startPoint" form:"startPoint"`
	Target     int64   `json:"target" form:"target"`
	TargetID   int64   `json:"targetId" form:"targetId"`
	Type       string  `json:"type" form:"type"`
}

type G6Group struct {
}

type G6Node struct {
	ID int64 `json:"id" form:"id"`
	//TagID      int64       `json:"tagID" form:"tagID"`
	Name       string      `json:"name" form:"name"`
	Label      string      `json:"label" form:"label"`
	Size       []string    `json:"size" form:"size"`
	Type       string      `json:"type" form:"type"`
	X          float64     `json:"x" form:"x"`
	Y          float64     `json:"y" form:"y"`
	Shape      string      `json:"shape" form:"shape"`
	Color      string      `json:"color" form:"color"`
	Image      string      `json:"image" form:"image"`
	StateImage string      `json:"stateImage" form:"stateImage"`
	OffsetX    float64     `json:"offsetX" form:"offsetX"`
	OffsetY    float64     `json:"offsetY" form:"offsetY"`
	InPoints   [][]float64 `json:"inPoints" form:"inPoints"`
	OutPoints  [][]float64 `json:"outPoints" form:"outPoints"`
}

type G6Graph struct {
	Edges  []*G6Edge  `json:"edges" form:"edges"`
	Groups []*G6Group `json:"groups" form:"groups"`
	Nodes  []*G6Node  `json:"nodes" form:"nodes"`
}

//type TemplatePoint struct {
//	ID             int64   `json:"id" gorm:"primary_key;column:id"`
//	X              float64 `json:"x" form:"x" gorm:"x"`
//	Y              float64 `json:"y" form:"y" gorm:"y"`
//	TemplateEdgeID int64   `json:"template_edge_id" gorm:"template_edge_id"`
//}
//
//func (this TemplatePoint) TableName() string {
//	return "template_point"
//}
//
//type TemplateEdge struct {
//	ID         int64  `json:"id" gorm:"primary_key;column:id"`
//	Template   int64  `json:"template" gorm:"column:template"`
//	End        string `json:"end" gorm:"end;type:json;"`
//	EndPoint   string `json:"endPoint" gorm:"end_point;type:json;"`
//	G6Edge       string `json:"edge" gorm:"edge"`
//	Shape      string `json:"shape" gorm:"shape"`
//	Source     string `json:"source" gorm:"source"`
//	SourceID   string `json:"sourceId" gorm:"source_id"`
//	Start      string `json:"start" gorm:"start;type:json;"`
//	StartPoint string `json:"startPoint" gorm:"start_point;type:json;"`
//	Target     string `json:"target" gorm:"target"`
//	TargetID   string `json:"targetId" gorm:"target_id"`
//	Type       string `json:"type" gorm:"'type'"`
//}
//
//func (this TemplateEdge) TableName() string {
//	return "template_edge"
//}
//
//type TemplateNode struct {
//	Template   int64  `json:"template" gorm:"primary_key;column:template"`
//	ID         string `json:"id" gorm:"primary_key;column:id"`
//	Name       string `json:"name" gorm:"name;type:varchar(256);index;comment:"`
//	Label      string `json:"label" gorm:"label;type:varchar(256);comment:"`
//	Size       string `json:"size" gorm:"size;type:json"`
//	Type       string `json:"type" gorm:"'type';type:varchar(256);comment:"`
//	X          int    `json:"x" gorm:"x"`
//	Y          int    `json:"y" gorm:"y"`
//	Shape      string `json:"shape" gorm:"shape"`
//	Color      string `json:"color" gorm:"color"`
//	OffsetX    int    `json:"offsetX" gorm:"offset_x"`
//	OffsetY    int    `json:"offsetY" gorm:"offset_y"`
//	Image      string `json:"image" gorm:"image"`
//	StateImage string `json:"stateImage" gorm:"state_image"`
//	InPoints   string `json:"inPoints" gorm:"in_points;type:json;"`
//	OutPoints  string `json:"outPoints" gorm:"out_points;type:json;"`
//}
//
//func (this TemplateNode) TableName() string {
//	return "template_node"
//}
//
//type TemplateGroup struct {
//	Template int64 `json:"template" gorm:"column:template"`
//}
//
//func (this TemplateGroup) TableName() string {
//	return "template_group"
//}
