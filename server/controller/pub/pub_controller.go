package pub

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/pub"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetPubListInputs struct {
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
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
	db := g.Con().Portal.Debug().Model(pub.Pub{})
	db.Count(&totalCount)
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&pubs)

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

// @Summary 创建新发布单
// @Description
// @Produce json
// @Param APIPostPubUpdateInputs body APIPostPubUpdateInputs true "创建新发布单"
// @Success 200 {object} APIPostPubUpdateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/pub/create [post]
func Create(c *gin.Context) {
	var inputs APIPostPubUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	user, _ := h.GetUser(c)

	var deployUnit *app.Tag
	g.Con().Portal.Model(app.Tag{}).Where("id = ?", inputs.DeployUnitID).Find(&deployUnit)

	p := pub.Pub{
		DeployUnitID:          inputs.DeployUnitID,
		DeployUnitName:        deployUnit.Name,
		VersionDate:           inputs.VersionDate,
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
		CreateAt:              gtime.Now(),
	}
	tx := g.Con().Portal
	if dt := tx.Model(pub.Pub{}).Create(&p); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	h.JSONR(c, h.OKStatus, inputs)
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

	var deployUnit *app.Tag
	g.Con().Portal.Model(app.Tag{}).Where("id = ?", inputs.DeployUnitID).Find(&deployUnit)

	p := pub.Pub{
		ID:                    inputs.ID,
		DeployUnitID:          inputs.DeployUnitID,
		DeployUnitName:        deployUnit.Name,
		VersionDate:           inputs.VersionDate,
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
	}
	tx := g.Con().Portal.Model(pub.Pub{}).Debug()
	if dt := tx.Where("id = ?", inputs.ID).Updates(p); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
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
