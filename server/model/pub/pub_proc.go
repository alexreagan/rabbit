package pub

import "github.com/alexreagan/rabbit/server/model/gtime"

type PubProc struct {
	ID            int64       `json:"id" gorm:"primary_key;column:id"`
	PubID         int64       `json:"pubID" gorm:"column:pub_id;comment:发布单ID"`
	TemplateID    string      `json:"templateID" gorm:"column:template_id;type:string;size:128;comment:模板ID"`
	ProcessInstID string      `json:"processInstID" gorm:"column:process_inst_id;type:string;size:128;comment:流程ID"`
	TaskID        string      `json:"taskID" gorm:"column:task_id;type:string;size:128;comment:流程步骤ID"`
	Creator       string      `json:"creator" gorm:"column:creator;type:string;size:128;comment:创建人"`
	CreateAt      gtime.GTime `json:"createAt" gorm:"column:create_at;comment:创建时间"`
}

func (this PubProc) TableName() string {
	return "pub_proc"
}
