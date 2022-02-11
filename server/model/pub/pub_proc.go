package pub

import "github.com/alexreagan/rabbit/server/model/gtime"

type PubProc struct {
	ID int64 `json:"id" gorm:"primary_key;column:id"`
	//PubID         int64  `json:"pubID" gorm:"column:pub_id;comment:发布单ID"`
	TemplateID    string      `json:"templateID" gorm:"column:template_id;type:string;size:128;comment:模板ID"`
	ProcessInstID string      `json:"processInstID" gorm:"column:process_inst_id;type:string;size:128;comment:流程ID"`
	TodoID        string      `json:"todoID" gorm:"column:todo_id;type:string;size:128;comment:待办ID"`
	TaskID        string      `json:"taskID" gorm:"column:task_id;type:string;size:128;comment:流程步骤ID"`
	Auditor       string      `json:"auditor" gorm:"column:auditor;type:string;size:128;comment:处理人"`
	AuditorName   string      `json:"auditorName" gorm:"column:auditor_name;type:string;size:128;comment:处理人姓名"`
	OpinCode      string      `json:"opinCode" gorm:"column:opin_code;type:string;size:128;comment:处理结论代码"`
	OpinDesc      string      `json:"opinDesc" gorm:"column:opin_desc;type:string;size:128;comment:处理结论描述"`
	Remark        string      `json:"remark" gorm:"column:remark;type:text;comment:处理意见"`
	NextNode      string      `json:"nextNode" gorm:"column:next_node;type:string;size:256;comment:下一步节点"`
	NextUserGrp   string      `json:"nextUserGrp" gorm:"column:next_user_grp;type:string;size:256;comment:下一步处理人"`
	ExeFstTask    string      `json:"exeFstTask" gorm:"column:exe_fst_task;type:string;size:512;comment:"`
	ButtonName    string      `json:"buttonName" gorm:"column:button_name;type:string;size:16;comment:操作"`
	PrjID         string      `json:"prjID" gorm:"column:prj_id;type:string;size:64;comment:发布单名称"`
	PrjSN         string      `json:"prjSN" gorm:"column:prj_sn;type:string;size:64;comment:发布单ID"`
	Conditions    string      `json:"conditions" gorm:"column:conditions;type:string;size:512;comment:条件"`
	CreateAt      gtime.GTime `json:"createAt" gorm:"column:create_at;comment:创建时间"`
}

func (this PubProc) TableName() string {
	return "pub_proc"
}
