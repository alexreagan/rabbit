package pub

import "github.com/alexreagan/rabbit/server/model/gtime"

type Buzz struct {
	ID          int64       `json:"id" gorm:"primary_key;column:id"`
	Name        string      `json:"name" gorm:"column:name;type:string;size:128;comment:"`
	WfeTemplate string      `json:"wfeTemplate" gorm:"column:wfe_template;type:string;size:512;comment:使用的流程模板"`
	Creator     string      `json:"creator" gorm:"column:creator;type:string;size:64;comment:创建人"`
	CreateAt    gtime.GTime `json:"createAt" gorm:"column:createAt;comment:创建时间"`
}

func (this Buzz) TableName() string {
	return "buzz"
}
