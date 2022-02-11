package sys

import "github.com/alexreagan/rabbit/server/model/gtime"

type Notice struct {
	ID          int64       `json:"id" gorm:"primary_key;column:id"`
	Title       string      `json:"title" gorm:"column:title;type:string;size:256;comment:标题"`
	Content     string      `json:"content" gorm:"column:content;type:text;comment:内容"`
	TimeBegin   gtime.GTime `json:"timeBegin" gorm:"column:time_begin;default:null;comment:开始时间"`
	TimeEnd     gtime.GTime `json:"timeEnd" gorm:"column:time_end;default:null;comment:结束时间"`
	Creator     string      `json:"creator" gorm:"column:creator;type:string;size:64;comment:创建人"`
	CreatorName string      `json:"creatorName" gorm:"column:creator_name;type:string;size:64;comment:创建人"`
	CreateAt    gtime.GTime `json:"createAt" gorm:"column:create_at;default:null;comment:创建时间"`
}

func (this Notice) TableName() string {
	return "notice"
}
