package pub

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type Pub struct {
	ID                    int64       `json:"id" gorm:"primary_key;column:id"`
	DeployUnitID          int64       `json:"deployUnitID" gorm:"column:deploy_unit_id;comment:"`
	DeployUnitName        string      `json:"deployUnitName" gorm:"column:deploy_unit_name;type:string;size:128;comment:"`
	VersionDate           gtime.GTime `json:"versionDate" gorm:"column:versionDate;comment:"`
	Content               string      `json:"content" gorm:"column:content;type:text;comment:"`
	Requirement           string      `json:"requirement" gorm:"column:requirement;type:text;comment:"`
	AppDesign             string      `json:"appDesign" gorm:"column:app_design;type:enum('notstart','running','finished','trim');comment:"`
	AppAssemblyTestDesign string      `json:"appAssemblyTestDesign" gorm:"column:app_assembly_test_design;type:enum('notstart','running','finished','trim');comment:"`
	AppAssemblyTestCase   string      `json:"appAssemblyTestCase" gorm:"column:app_assembly_test_case;type:enum('notstart','running','finished','trim');comment:"`
	AppAssemblyTestReport string      `json:"appAssemblyTestReport" gorm:"column:app_assembly_test_report;type:enum('notstart','running','finished','trim');comment:"`
	UserTestCase          string      `json:"userTestCase" gorm:"column:user_test_case;type:enum('notstart','running','finished','trim');comment:"`
	UserTestReport        string      `json:"userTestReport" gorm:"column:user_test_report;type:enum('notstart','running','finished','trim');comment:"`
	CodeReview            string      `json:"codeReview" gorm:"column:code_review;type:enum('notstart','running','finished','trim');comment:"`
	PubControlTable       string      `json:"pubControlTable" gorm:"column:pub_control_table;type:enum('notstart','running','finished','trim');comment:"`
	PubShellReview        string      `json:"pubShellReview" gorm:"column:pub_shell_review;type:enum('notstart','running','finished','trim');comment:"`
	TrialOperationDesign  string      `json:"trialOperationDesign" gorm:"column:trial_operation_design;type:enum('notstart','running','finished','trim');comment:"`
	TrialOperationCase    string      `json:"trialOperationCase" gorm:"column:trial_operation_case;type:enum('notstart','running','finished','trim');comment:"`
	Creator               string      `json:"creator" gorm:"column:creator;type:string;size:128;comment:"`
	CreateAt              gtime.GTime `json:"createAt" gorm:"column:create_at;comment:"`
}

func (this Pub) TableName() string {
	return "pub"
}
