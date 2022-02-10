package pub

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/pub"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type Condition struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NextUser struct {
	ID            string `json:"ID"`
	Name          string `json:"name"`
	IDPrcActionID string `json:"IDPrcActionID"`
	UsrIDLandNm   string `json:"UsrIDLandNm"`
	CurUserInstID string `json:"curUserInstID"`
	CurUserInstNm string `json:"curUserInstNm"`
}

type ApiPostProcCreateInputs struct {
	TemplateID   string       `json:"templateID" form:"templateID"`
	TaskID       string       `json:"taskID" form:"taskID"`
	Remark       string       `json:"remark" form:"remark"`
	OpinDesc     string       `json:"opinDesc" form:"opinDesc"`
	NextUserFlag string       `json:"nextUserFlag" form:"nextUserFlag"`
	ButtonName   string       `json:"buttonName" form:"buttonName"`
	UserID       string       `json:"userID" form:"userID"`
	UserName     string       `json:"userName" form:"userName"`
	UsrIDLandNm  string       `json:"usrIDLandNm" form:"usrIDLandNm"`
	CurUsrInstID string       `json:"curUsrInstID" form:"curUsrInstID"`
	CurUsrInstNm string       `json:"curUsrInstNm" form:"curUsrInstNm"`
	NextUserGrp  []*NextUser  `json:"nextUserGrp" form:"nextUserGrp"`
	Conditions   []*Condition `json:"conditions" form:"conditions"`
	PrjID        string       `json:"prjID" form:"prjID"`
	PrjSn        string       `json:"prjSn" form:"prjSn"`
	ToDoTmTpCd   string       `json:"toDoTmTpCd" form:"toDoTmTpCd"`
	ToDoTmTtl    string       `json:"toDoTmTtl" form:"toDoTmTtl"`
	BlngInstID   string       `json:"blngInstID" form:"blngInstID"`
	DmnGrpID     string       `json:"dmnGrpID" form:"dmnGrpID"`
}

type ApiPostProcCreateOutputs struct {
}

// @Summary 流程发起
// @Description
// @Produce json
// @Param ApiPostProcCreateInputs body ApiPostProcCreateInputs true "发起信息"
// @Success 200 {object} ApiPostProcCreateOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/proc/create [post]
func ProcCreate(c *gin.Context) {
	var inputs ApiPostProcCreateInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	nextUserGrp := make([]service.ProcManagerApiNextUser, 0)
	for _, nxg := range inputs.NextUserGrp {
		nextUserGrp = append(nextUserGrp, service.ProcManagerApiNextUser{
			ID:              nxg.ID,
			NAME:            nxg.Name,
			IDPRC_ACTION_ID: nxg.IDPrcActionID,
			USR_ID_LAND_NM:  nxg.UsrIDLandNm,
			CUR_USR_INST_ID: nxg.CurUserInstID,
			CUR_USR_INST_NM: nxg.CurUserInstNm,
		})
	}

	conditions := make([]service.ProcManagerApiCondition, 0)
	for _, cond := range inputs.Conditions {
		conditions = append(conditions, service.ProcManagerApiCondition{
			KEY:   cond.Key,
			VALUE: cond.Value,
		})
	}

	if service.ProcManagerApiService.Addr == "" {
		service.ProcManagerApiService.Addr = viper.GetString("procManager.addr")
	}

	resp, e := service.ProcManagerApiService.ProcManagerApiProcCreate(service.ProcManagerApiCreateInputs{
		TEMPLATE_ID:     inputs.TemplateID,
		TASK_ID:         inputs.TaskID,
		REMARK:          inputs.Remark,
		USER_ID:         inputs.UserID,
		CUR_USR_INST_ID: inputs.CurUsrInstID,
		CUR_USR_INST_NM: inputs.UsrIDLandNm,
		PRJ_ID:          inputs.PrjID,
		PRJ_SN:          inputs.PrjSn,
		TO_DO_TM_TTL:    inputs.ToDoTmTtl,
		BUTTON_NAME:     inputs.ButtonName,
		USER_NAME:       inputs.UserName,
		USR_ID_LAND_NM:  inputs.UsrIDLandNm,
		NEXT_USER_GRP:   nextUserGrp,
		CONDITIONS:      conditions,
	})
	if e != nil {
		h.JSONR(c, http.StatusExpectationFailed, e)
		return
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type ApiPostProcExecuteInputs struct {
	ProcessInstID string       `json:"processInstID" form:"processInstID"`
	TemplateID    string       `json:"templateID" form:"templateID"`
	TaskID        string       `json:"taskID" form:"taskID"`
	Remark        string       `json:"remark" form:"remark"`
	OpinCode      string       `json:"opinCode" form:"opinCode"`
	OpinDesc      string       `json:"opinDesc" form:"opinDesc"`
	NextUserFlag  string       `json:"nextUserFlag" form:"nextUserFlag"`
	ButtonName    string       `json:"buttonName" form:"buttonName"`
	UserID        string       `json:"userID" form:"userID"`
	UserName      string       `json:"userName" form:"userName"`
	UsrIDLandNm   string       `json:"usrIDLandNm" form:"usrIDLandNm"`
	CurUsrInstID  string       `json:"curUsrInstID" form:"curUsrInstID"`
	CurUsrInstNm  string       `json:"curUsrInstNm" form:"curUsrInstNm"`
	PrjID         string       `json:"prjID" form:"prjID"`
	PrjSn         string       `json:"prjSn" form:"prjSn"`
	ToDoID        string       `json:"toDoID" form:"toDoID"`
	ToDoTmTtl     string       `json:"toDoTmTtl" form:"toDoTmTtl"`
	NextUserGrp   []*NextUser  `json:"nextUserGrp" form:"nextUserGrp"`
	Conditions    []*Condition `json:"conditions" form:"conditions"`
}

type ApiPostProcExecuteOutputs struct {
}

// @Summary 流程处理
// @Description
// @Produce json
// @Param ApiPostProcExecuteInputs body ApiPostProcExecuteInputs true "处理信息"
// @Success 200 {object} ApiPostProcCreateOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/proc/execute [post]
func ProcExecute(c *gin.Context) {
	var inputs ApiPostProcExecuteInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	nextUserGrp := make([]service.ProcManagerApiNextUser, 0)
	for _, nxg := range inputs.NextUserGrp {
		nextUserGrp = append(nextUserGrp, service.ProcManagerApiNextUser{
			ID:              nxg.ID,
			NAME:            nxg.Name,
			PRC_ACTION_ID:   nxg.IDPrcActionID,
			IDPRC_ACTION_ID: nxg.IDPrcActionID,
			USR_ID_LAND_NM:  nxg.UsrIDLandNm,
			CUR_USR_INST_ID: nxg.CurUserInstID,
			CUR_USR_INST_NM: nxg.CurUserInstNm,
		})
	}

	conditions := make([]service.ProcManagerApiCondition, 0)
	for _, cond := range inputs.Conditions {
		conditions = append(conditions, service.ProcManagerApiCondition{
			KEY:   cond.Key,
			VALUE: cond.Value,
		})
	}

	if service.ProcManagerApiService.Addr == "" {
		service.ProcManagerApiService.Addr = viper.GetString("procManager.addr")
	}

	resp, e := service.ProcManagerApiService.ProcManagerApiExecute(service.ProcManageApiProcExecuteInputs{
		PROCESS_INST_ID: inputs.ProcessInstID,
		TEMPLATE_ID:     inputs.TemplateID,
		TASK_ID:         inputs.TaskID,
		REMARK:          inputs.Remark,
		OPIN_DESC:       inputs.OpinDesc,
		PRJ_ID:          inputs.PrjID,
		PRJ_SN:          inputs.PrjSn,
		TO_DO_TM_TTL:    inputs.ToDoTmTtl,
		TODO_ID:         inputs.ToDoID,
		NEXT_USER_GRP:   nextUserGrp,
		BUTTON_NAME:     inputs.ButtonName,
		CONDITIONS:      conditions,
	})
	if e != nil {
		h.JSONR(c, http.StatusExpectationFailed, e)
		return
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type ApiGetProcNextNodeInfoInputs struct {
	ProcessInstID string       `json:"processInstID" form:"processInstID"`
	TemplateID    string       `json:"templateID" form:"templateID"`
	TaskTD        string       `json:"taskID" form:"taskID"`
	Conditions    []*Condition `json:"conditions" form:"conditions"`
}

type ApiGetProcNextNodeInfoOutputs struct {
}

// @Summary 查询下一环节信息
// @Description
// @Produce json
// @Param ApiGetProcNextNodeInfoInputs body ApiGetProcNextNodeInfoInputs true "查询下一环节信息"
// @Success 200 {object} ApiGetProcNextNodeInfoOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/proc/nextNodeInfo [get]
func NextNodeInfo(c *gin.Context) {
	var inputs ApiGetProcNextNodeInfoInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	conditions := make([]*service.Condition, 0)
	for _, cond := range conditions {
		conditions = append(conditions, &service.Condition{
			Key:   cond.Key,
			Value: cond.Value,
		})
	}

	resp, e := service.ProcService.NextNodeInfo(service.NextNodeInfoInputs{
		ProcessInstID: inputs.ProcessInstID,
		TemplateID:    inputs.TemplateID,
		TaskTD:        inputs.TaskTD,
		Conditions:    conditions,
	})
	if e != nil {
		h.JSONR(c, http.StatusExpectationFailed, e)
		return
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type ApiGetProcGetPersonByNodeInputs struct {
	TemplateID string `json:"templateID" form:"templateID"`
	TaskID     string `json:"taskID" form:"taskID"`
}

// @Summary 根据节点获取审批用户
// @Description
// @Produce json
// @Param ApiGetProcNextNodeInfoInputs body ApiGetProcNextNodeInfoInputs true "根据节点获取审批用户"
// @Success 200 {object} ApiGetProcNextNodeInfoOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/proc/getPersonByNode [get]
func GetPersonByNode(c *gin.Context) {
	var inputs ApiGetProcGetPersonByNodeInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	resp, e := service.ProcService.GetPersonByNode(inputs.TemplateID, inputs.TaskID)
	if e != nil {
		h.JSONR(c, http.StatusExpectationFailed, e)
		return
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type ApiGetProcGetHistDetailListInputs struct {
	TemplateID    string `json:"templateID" form:"templateID"`
	ProcessInstID string `json:"processInstID" form:"processInstID"`
	TaskID        string `json:"taskID" form:"taskID"`
	BelongInstID  string `json:"BelongInstID" form:"BelongInstID"`
	SelectMode    string `json:"SelectMode" form:"SelectMode"`
}

// @Summary 获取审批历史
// @Description
// @Produce json
// @Param ApiGetProcGetHistDetailListInputs body ApiGetProcGetHistDetailListInputs true "获取审批历史"
// @Success 200 {object} ApiGetProcNextNodeInfoOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/proc/getHistDetailList [get]
func GetHistDetailList(c *gin.Context) {
	var inputs ApiGetProcGetHistDetailListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	resp, e := service.ProcService.GetHistDetailList(inputs.ProcessInstID, inputs.TaskID, inputs.BelongInstID, inputs.SelectMode)
	if e != nil {
		h.JSONR(c, http.StatusExpectationFailed, e)
		return
	}
	h.JSONR(c, http.StatusOK, resp)
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
func ProcInfo(c *gin.Context) {
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
	} else {
		resp = &APIGetPubProcInfoOutputs{
			TemplateID:    "600100PubAudit",
			ProcessInstID: "",
			TaskID:        "",
		}
	}
	h.JSONR(c, h.OKStatus, resp)
	return
}
