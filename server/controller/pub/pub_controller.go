package pub

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/pub"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//const pubProcTemplateID = "600100PubAudit"
const pubProcTemplateID = "040500ChgTskBizReq"
const pubProcStartTaskID = "10101"

type APIGetPubListInputs struct {
	DeployUnitID int    `json:"deployUnitID" form:"deployUnitID"`
	Creator      string `json:"creator" form:"creator"`
	Limit        int    `json:"limit" form:"limit"`
	Page         int    `json:"page" form:"page"`
	OrderBy      string `json:"orderBy" form:"orderBy"`
	Order        string `json:"order" form:"order"`
}

type APIGetPubListOutputs struct {
	List       []*pub.Pub `json:"list"`
	TotalCount int64      `json:"totalCount"`
}

// @Summary 发布列表接口
// @Description
// @Produce json
// @Param APIGetPubListInputs query APIGetPubListInputs true "根据查询条件分页查询发布列表"
// @Success 200 {object} APIGetPubListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/pub/list [get]
func List(c *gin.Context) {
	var inputs APIGetPubListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var pubs []*pub.Pub
	var totalCount int64
	tx := g.Con().Portal.Model(pub.Pub{})
	if inputs.DeployUnitID != 0 {
		tx = tx.Where("deploy_unit_id = ?", inputs.DeployUnitID)
	}
	if inputs.Creator != "" {
		tx = tx.Where("creator_name regexp ?", inputs.Creator)
	}
	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx = tx.Offset(offset).Limit(limit)
	tx.Find(&pubs)

	resp := &APIGetPubListOutputs{
		List:       pubs,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIPostPubUpdateInputs struct {
	ID                    int64       `json:"id" form:"id"`
	DeployUnitID          int64       `json:"deployUnitID" form:"deployUnitID"`
	DeployUnitName        string      `json:"deployUnitName" form:"deployUnitName"`
	VersionDate           gtime.GTime `json:"versionDate" form:"versionDate"`
	PubContent            string      `json:"pubContent" form:"pubContent"`
	Git                   string      `json:"git" form:"git"`
	CommitID              string      `json:"commitID" form:"commitID"`
	PackageAddress        string      `json:"packageAddress" form:"packageAddress"`
	PubStep               string      `json:"pubStep" form:"pubStep"`
	RollbackStep          string      `json:"rollbackStep" form:"rollbackStep"`
	Requirement           string      `json:"requirement" form:"requirement"`
	AppDesign             string      `json:"appDesign" form:"appDesign"`
	AppAssemblyTestDesign string      `json:"appAssemblyTestDesign" form:"appAssemblyTestDesign"`
	AppAssemblyTestCase   string      `json:"appAssemblyTestCase" form:"appAssemblyTestCase"`
	AppAssemblyTestReport string      `json:"appAssemblyTestReport" form:"appAssemblyTestReport"`
	UserTestCase          string      `json:"userTestCase" form:"userTestCase"`
	UserTestReport        string      `json:"userTestReport" form:"userTestReport"`
	CodeReview            string      `json:"codeReview" form:"codeReview"`
	PubControlTable       string      `json:"pubControlTable" form:"pubControlTable"`
	PubShellReview        string      `json:"pubShellReview" form:"pubShellReview"`
	TrialOperationDesign  string      `json:"trialOperationDesign" form:"trialOperationDesign"`
	TrialOperationCase    string      `json:"trialOperationCase" form:"trialOperationCase"`
}

type APIPostPubCreateInputs struct {
	DeployUnitID   int64                  `json:"deployUnitID" form:"deployUnitID"`
	DeployUnitName string                 `json:"deployUnitName" form:"deployUnitName"`
	VersionDate    time.Time              `json:"versionDate" form:"versionDate"`
	PubContent     string                 `json:"pubContent" form:"pubContent"`
	Git            string                 `json:"git" form:"git"`
	CommitID       string                 `json:"commitID" form:"commitID"`
	PackageAddress string                 `json:"packageAddress" form:"packageAddress"`
	PubStep        string                 `json:"pubStep" form:"pubStep"`
	RollbackStep   string                 `json:"rollbackStep" form:"rollbackStep"`
	TemplateID     string                 `json:"templateID" form:"templateID,omitempty"`
	ProcessInstID  string                 `json:"processInstID" form:"processInstID,omitempty"`
	TaskID         string                 `json:"taskID" form:"taskID,omitempty"`
	Remark         string                 `json:"remark" form:"remark,omitempty"`
	ButtonName     string                 `json:"buttonName" form:"buttonName,omitempty"`
	NextNode       service.TXNextNodeInfo `json:"nextNode" form:"nextNode,omitempty"`
	//NextUserFlag   string      `json:"nextUserFlag" form:"nextUserFlag,omitempty"`
	NextUserGrp []string `json:"nextUserGrp[]" form:"nextUserGrp[],omitempty"`
}

// @Summary 创建新发布单并启动发布流程
// @Description
// @Produce json
// @Param APIPostPubUpdateInputs body APIPostPubUpdateInputs true "创建新发布单并启动发布流程"
// @Success 200 {object} APIPostPubUpdateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/pub/create [post]
func Create(c *gin.Context) {
	var inputs APIPostPubCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	u, _ := h.GetUser(c)

	var deployUnit *app.Tag
	g.Con().Portal.Model(app.Tag{}).Where("id = ?", inputs.DeployUnitID).Find(&deployUnit)

	tx := g.Con().Portal.Begin()
	// 发布单信息
	p := pub.Pub{
		DeployUnitID:   inputs.DeployUnitID,
		DeployUnitName: deployUnit.Name,
		VersionDate:    gtime.NewGTime(inputs.VersionDate),
		Git:            inputs.Git,
		CommitID:       inputs.CommitID,
		PackageAddress: inputs.PackageAddress,
		PubContent:     inputs.PubContent,
		PubStep:        inputs.PubStep,
		RollbackStep:   inputs.RollbackStep,
		Creator:        u.JgygUserID,
		CreatorName:    u.CnName,
		CreateAt:       gtime.Now(),
	}
	if err := tx.Model(pub.Pub{}).Create(&p).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	// 发布单流程信息
	nextUserGrp, _ := json.Marshal(inputs.NextUserGrp)
	nextNode, _ := json.Marshal(inputs.NextNode)
	prjSN := strconv.FormatInt(p.ID, 10)
	prjID := fmt.Sprintf("%s-版本发布单", prjSN)
	pubProc := pub.PubProc{
		TemplateID:    inputs.TemplateID,
		ProcessInstID: inputs.ProcessInstID,
		TaskID:        inputs.TaskID,
		Auditor:       u.JgygUserID,
		AuditorName:   u.CnName,
		Remark:        inputs.Remark,
		NextUserGrp:   string(nextUserGrp),
		ButtonName:    inputs.ButtonName,
		NextNode:      string(nextNode),
		PrjID:         prjID,
		PrjSN:         prjSN,
		CreateAt:      gtime.Now(),
	}
	if err := tx.Model(pub.PubProc{}).Create(&pubProc).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	// 流程创建
	resp, err := service.WfeService.Create(&u, pubProcTemplateID, pubProcStartTaskID, inputs.Remark, inputs.ButtonName,
		inputs.NextUserGrp, inputs.NextNode, prjID, prjSN, "", prjID, "", "")
	if err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	tx.Commit()
	h.JSONR(c, h.OKStatus, resp)
	return
}

type Condition struct {
	Key   string `json:"key" form:"key,omitempty"`
	Value string `json:"value" form:"value,omitempty"`
}

type APIPostPubExecuteInputs struct {
	TemplateID    string                 `json:"templateID" form:"templateID"`
	ProcessInstID string                 `json:"processInstID" form:"processInstID"`
	TodoID        string                 `json:"todoID" form:"todoID"`
	TaskID        string                 `json:"taskID" form:"taskID"`
	Remark        string                 `json:"remark" form:"remark,omitempty"`
	OpinCode      string                 `json:"opinCode" form:"opinCode,omitempty"`
	OpinDesc      string                 `json:"opinDesc" form:"opinDesc,omitempty"`
	ExeFstTask    string                 `json:"exeFstTask" form:"exeFstTask,omitempty"`
	ButtonName    string                 `json:"buttonName" form:"buttonName,omitempty"`
	NextNode      service.TXNextNodeInfo `json:"nextNode" form:"nextNode,omitempty"`
	//NextUserFlag  string              `json:"nextUserFlag" form:"nextUserFlag,omitempty"`
	PrjID       string              `json:"prjID" form:"prjID,omitempty"`
	PrjSN       string              `json:"prjSN" form:"prjSN,omitempty"`
	NextUserGrp []string            `json:"nextUserGrp[]" form:"nextUserGrp[],omitempty"`
	Conditions  []service.Condition `json:"conditions[]" form:"conditions[],omitempty"`
}

// @Summary 处理发布单信息
// @Description
// @Produce json
// @Param IP formData string true "处理发布单信息"
// @Success 200 {object} service.WfeExecuteResponse
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/pub/execute [post]
func Execute(c *gin.Context) {
	var inputs APIPostPubExecuteInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	u, _ := h.GetUser(c)

	nextUserGrp, _ := json.Marshal(inputs.NextUserGrp)
	conditions, _ := json.Marshal(inputs.Conditions)
	// 发布单流程信息
	tx := g.Con().Portal.Begin()
	pubProc := pub.PubProc{
		TemplateID:    inputs.TemplateID,
		ProcessInstID: inputs.ProcessInstID,
		TaskID:        inputs.TaskID,
		Auditor:       u.JgygUserID,
		AuditorName:   u.CnName,
		OpinCode:      inputs.OpinCode,
		OpinDesc:      inputs.OpinDesc,
		Remark:        inputs.Remark,
		NextUserGrp:   string(nextUserGrp),
		ExeFstTask:    inputs.ExeFstTask,
		ButtonName:    inputs.ButtonName,
		PrjID:         inputs.PrjID,
		PrjSN:         inputs.PrjSN,
		Conditions:    string(conditions),
		CreateAt:      gtime.Now(),
	}
	if err := tx.Model(pub.PubProc{}).Create(&pubProc).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	// 流程待办转在办
	resp, err := service.WfeService.Todo2Doing(&u, inputs.TodoID)
	if err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}
	if resp.TXBody.Entity.ResultDesc == "0" {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, errors.New("待办转在办失败"))
		return
	}
	// 流程处理
	resp2, err := service.WfeService.Execute(&u, inputs.TemplateID, inputs.ProcessInstID, inputs.TaskID, inputs.Remark, inputs.OpinCode,
		inputs.OpinDesc, inputs.ButtonName, inputs.ExeFstTask, inputs.NextUserGrp, inputs.NextNode, inputs.Conditions)
	if err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	tx.Commit()
	h.JSONR(c, h.OKStatus, resp2)
	return
}

// @Summary 更新发布单信息
// @Description
// @Produce json
// @Param IP formData string true "更新发布单信息"
// @Success 200 {object} APIPostPubUpdateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/pub/update [put]
func Update(c *gin.Context) {
	var inputs APIPostPubUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	user, _ := h.GetUser(c)

	var deployUnit *app.Tag
	g.Con().Portal.Model(app.Tag{}).Where("id = ?", inputs.DeployUnitID).Find(&deployUnit)

	p := pub.Pub{
		ID:                    inputs.ID,
		DeployUnitID:          inputs.DeployUnitID,
		DeployUnitName:        deployUnit.Name,
		VersionDate:           inputs.VersionDate,
		Git:                   inputs.Git,
		CommitID:              inputs.CommitID,
		PackageAddress:        inputs.PackageAddress,
		PubContent:            inputs.PubContent,
		PubStep:               inputs.PubStep,
		RollbackStep:          inputs.RollbackStep,
		Requirement:           inputs.Requirement,
		AppDesign:             inputs.AppDesign,
		AppAssemblyTestDesign: inputs.AppAssemblyTestDesign,
		AppAssemblyTestCase:   inputs.AppAssemblyTestCase,
		AppAssemblyTestReport: inputs.AppAssemblyTestReport,
		UserTestCase:          inputs.UserTestCase,
		UserTestReport:        inputs.UserTestReport,
		CodeReview:            inputs.CodeReview,
		PubControlTable:       inputs.PubControlTable,
		PubShellReview:        inputs.PubShellReview,
		TrialOperationDesign:  inputs.TrialOperationDesign,
		TrialOperationCase:    inputs.TrialOperationCase,
		Creator:               user.JgygUserID,
		CreatorName:           user.CnName,
		CreateAt:              gtime.Now(),
	}
	tx := g.Con().Portal.Model(pub.Pub{})
	if err := tx.Where("id = ?", inputs.ID).Updates(&p).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, inputs)
	return
}

// @Summary 根据ID获取发布单详细信息
// @Description
// @Produce json
// @Param id query int true "根据ID获取发布单详细信息"
// @Success 200 {object} pub.Pub
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/pub/info [get]
func Info(c *gin.Context) {
	id := c.Query("id")
	p := pub.Pub{}
	g.Con().Portal.Model(p).Where("id = ?", id).First(&p)
	h.JSONR(c, p)
	return
}

type APIPostPubAssignInputs struct {
	ID     int64  `json:"id" form:"id"`
	Status string `json:"status" form:"status"`
}

// @Summary 更新发布单实施状态信息
// @Description
// @Produce json
// @Param APIPostPubAssignInputs body APIPostPubAssignInputs true "更新发布单实施状态信息"
// @Success 200 {object} APIPostPubUpdateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/pub/assign [put]
func Assign(c *gin.Context) {
	var inputs APIPostPubAssignInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	user, _ := h.GetUser(c)
	p := pub.Pub{
		ID:              inputs.ID,
		Status:          inputs.Status,
		Implementer:     user.JgygUserID,
		ImplementerName: user.UserName,
		ImplementAt:     gtime.Now(),
	}
	tx := g.Con().Portal.Model(pub.Pub{})
	if err := tx.Where("id = ?", inputs.ID).Updates(&p).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, inputs)
	return
}

type APIGetPubProcInfoInputs struct {
	ID string `json:"ID" form:"ID"`
}

type APIGetPubProcInfoOutputs struct {
	TemplateID    string `json:"templateID" form:"templateID"`
	ProcessInstID string `json:"processInstID" form:"processInstID"`
	TaskID        string `json:"taskID" form:"taskID"`
}

// @Summary 发布单流程信息
// @Description
// @Produce json
// @Param APIGetPubProcInfoInputs body APIGetPubProcInfoOutputs true "发布单流程信息"
// @Success 200 {object} APIGetPubProcInfoOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/pub/proc/info [get]
func PubProcInfo(c *gin.Context) {
	var inputs APIGetPubProcInfoInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	var p pub.PubProc
	tx := g.Con().Portal.Model(pub.PubProc{})
	if err := tx.Where("pub_id = ?", inputs.ID).Order("create_at desc").First(&p).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	var resp *APIGetPubProcInfoOutputs
	if p.ID != 0 {
		resp = &APIGetPubProcInfoOutputs{
			TemplateID:    p.TemplateID,
			ProcessInstID: p.ProcessInstID,
			TaskID:        p.TaskID,
		}
	}
	h.JSONR(c, h.OKStatus, resp)
	return
}
