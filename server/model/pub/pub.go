package pub

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

// 发布单
type Pub struct {
	ID                    int64       `json:"id" gorm:"primary_key;column:id"`
	DeployUnitID          string       `json:"deployUnitID" gorm:"column:deploy_unit_id;comment:"`
	DeployUnitName        string      `json:"deployUnitName" gorm:"column:deploy_unit_name;type:string;size:128;comment:"`
	VersionDate           gtime.GTime `json:"versionDate" gorm:"column:versionDate;comment:"`
	Git                   string      `json:"git" gorm:"column:git;type:string;size:512;comment:git地址"`
	CommitID              string      `json:"commitID" gorm:"column:commit_id;type:string;size:128;comment:commit id"`
	PackageAddress        string      `json:"packageAddress" gorm:"column:package_address;type:string;size:512;comment:版本包地址"`
	PubContent            string      `json:"pubContent" gorm:"column:pub_content;type:text;comment:"`
	PubStep               string      `json:"pubStep" gorm:"column:pub_step;type:text;comment:"`
	RollbackStep          string      `json:"rollbackStep" gorm:"column:rollback_step;type:text;comment:"`
	Requirement           string      `json:"requirement" gorm:"column:requirement;type:text;comment:"`
	AppDesign             string      `json:"appDesign" gorm:"column:app_design;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	AppAssemblyTestDesign string      `json:"appAssemblyTestDesign" gorm:"column:app_assembly_test_design;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	AppAssemblyTestCase   string      `json:"appAssemblyTestCase" gorm:"column:app_assembly_test_case;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	AppAssemblyTestReport string      `json:"appAssemblyTestReport" gorm:"column:app_assembly_test_report;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	UserTestCase          string      `json:"userTestCase" gorm:"column:user_test_case;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	UserTestReport        string      `json:"userTestReport" gorm:"column:user_test_report;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	CodeReview            string      `json:"codeReview" gorm:"column:code_review;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	PubControlTable       string      `json:"pubControlTable" gorm:"column:pub_control_table;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	PubShellReview        string      `json:"pubShellReview" gorm:"column:pub_shell_review;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	TrialOperationDesign  string      `json:"trialOperationDesign" gorm:"column:trial_operation_design;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	TrialOperationCase    string      `json:"trialOperationCase" gorm:"column:trial_operation_case;type:enum('unstart','running','finished','trim');default:'unstart';comment:"`
	Creator               string      `json:"creator" gorm:"column:creator;type:string;size:128;comment:创建人"`
	CreatorName           string      `json:"creatorName" gorm:"column:creator_name;type:string;size:128;comment:创建人姓名"`
	CreateAt              gtime.GTime `json:"createAt" gorm:"column:create_at;comment:创建时间"`
	Status                string      `json:"status" gorm:"column:status;type:enum('submitted','success','failure','rolledback');default:'submitted';comment:"`
	Implementer           string      `json:"implementer" gorm:"column:implementer;type:string;size:128;comment:实施人"`
	ImplementerName       string      `json:"implementerName" gorm:"column:implementer_name;type:string;size:128;comment:实施人姓名"`
	ImplementAt           gtime.GTime `json:"implementAt" gorm:"column:implement_at;comment:实施时间"`
}

func (this Pub) TableName() string {
	return "pub"
}
