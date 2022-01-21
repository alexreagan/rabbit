package pub

import (
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Condition struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type NextUser struct {
	ID string `json:"ID"`
	Name string `json:"name"`
	IDPrcActionID string `json:"IDPrcActionID"`
	UsrIDLandNm string `json:"UsrIDLandNm"`
	CurUserInstID string `json:"curUserInstID"`
	CurUserInstNm string `json:"curUserInstNm"`
}

type ApiPostProcCreateInputs struct {
	TemplateID string `json:"templateID" form:"templateID"`
	TaskID string `json:"taskID" form:"taskID"`
	Remark string `json:"remark" form:"remark"`
	OpinDesc string `json:"opinDesc" form:"opinDesc"`
	NextUserFlag string `json:"nextUserFlag" form:"nextUserFlag"`
	ButtonName string `json:"buttonName" form:"buttonName"`
	UserID string `json:"userID" form:"userID"`
	UserName string `json:"userName" form:"userName"`
	UsrIDLandNm string `json:"usrIDLandNm" form:"usrIDLandNm"`
	CurUsrInstID string `json:"curUsrInstID" form:"curUsrInstID"`
	CurUsrInstNm string `json:"curUsrInstNm" form:"curUsrInstNm"`
	NextUserGrp  []*NextUser `json:"nextUserGrp" form:"nextUserGrp"`
	Conditions []*Condition `json:"conditions" form:"conditions"`
	PrjID  string `json:"prjID" form:"prjID"`
	PrjSn  string `json:"prjSn" form:"prjSn"`
	ToDoTmTpCd  string `json:"toDoTmTpCd" form:"toDoTmTpCd"`
	ToDoTmTtl  string `json:"toDoTmTtl" form:"toDoTmTtl"`
	BlngInstID  string `json:"blngInstID" form:"blngInstID"`
	DmnGrpID  string `json:"dmnGrpID" form:"dmnGrpID"`
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

	nextUserGrp := make([]*service.NextUser, 0)
	for _, nxg := range nextUserGrp {
		nextUserGrp = append(nextUserGrp, &service.NextUser{
			ID: nxg.ID,
			Name: nxg.Name,
			IDPrcActionID: nxg.IDPrcActionID,
			UsrIDLandNm: nxg.UsrIDLandNm,
			CurUserInstID: nxg.CurUserInstID,
			CurUserInstNm: nxg.CurUserInstNm,
		})
	}

	conditions := make([]*service.Condition, 0)
	for _, cond := range conditions {
		conditions = append(conditions, &service.Condition{
			Key: cond.Key,
			Value: cond.Value,
		})
	}

	resp, e := service.ProcService.ProcCreate(service.ProcCreateInputs{
		TemplateID:inputs.TemplateID,
		TaskID:inputs.TaskID,
		Remark:inputs.Remark,
		OpinDesc:inputs.OpinDesc,
		NextUserFlag:inputs.NextUserFlag,
		ButtonName:inputs.ButtonName,
		UserID:inputs.UserID,
		UserName:inputs.UserName,
		UsrIDLandNm:inputs.UsrIDLandNm,
		CurUsrInstID:inputs.CurUsrInstID,
		CurUsrInstNm:inputs.CurUsrInstNm,
		NextUserGrp:nextUserGrp,
		Conditions:conditions,
		PrjID:inputs.PrjID,
		PrjSn:inputs.PrjSn,
		ToDoTmTpCd:inputs.ToDoTmTpCd,
		ToDoTmTtl:inputs.ToDoTmTtl,
		BlngInstID:inputs.BlngInstID,
		DmnGrpID:inputs.DmnGrpID,
	})
	if e != nil {
		h.JSONR(c, http.StatusExpectationFailed, e)
		return
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}


type ApiPostProcExecuteInputs struct {
	ProcessInstID string `json:"processInstID" form:"processInstID"`
	TemplateID string `json:"templateID" form:"templateID"`
	TaskID string `json:"taskID" form:"taskID"`
	Remark string `json:"remark" form:"remark"`
	OpinCode string `json:"opinCode" form:"opinCode"`
	OpinDesc string `json:"opinDesc" form:"opinDesc"`
	NextUserFlag string `json:"nextUserFlag" form:"nextUserFlag"`
	ButtonName string `json:"buttonName" form:"buttonName"`
	UserID string `json:"userID" form:"userID"`
	UserName string `json:"userName" form:"userName"`
	UsrIDLandNm string `json:"usrIDLandNm" form:"usrIDLandNm"`
	CurUsrInstID string `json:"curUsrInstID" form:"curUsrInstID"`
	CurUsrInstNm string `json:"curUsrInstNm" form:"curUsrInstNm"`
	NextUserGrp  []*NextUser `json:"nextUserGrp" form:"nextUserGrp"`
	Conditions []*Condition `json:"conditions" form:"conditions"`
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

	nextUserGrp := make([]*service.NextUser, 0)
	for _, nxg := range nextUserGrp {
		nextUserGrp = append(nextUserGrp, &service.NextUser{
			ID: nxg.ID,
			Name: nxg.Name,
			IDPrcActionID: nxg.IDPrcActionID,
			UsrIDLandNm: nxg.UsrIDLandNm,
			CurUserInstID: nxg.CurUserInstID,
			CurUserInstNm: nxg.CurUserInstNm,
		})
	}

	conditions := make([]*service.Condition, 0)
	for _, cond := range conditions {
		conditions = append(conditions, &service.Condition{
			Key: cond.Key,
			Value: cond.Value,
		})
	}

	resp, e := service.ProcService.ProcExecute(service.ProcExecuteInputs{
		ProcessInstID:inputs.ProcessInstID,
		TemplateID:inputs.TemplateID,
		TaskID:inputs.TaskID,
		Remark:inputs.Remark,
		OpinCode:inputs.OpinCode,
		OpinDesc:inputs.OpinDesc,
		NextUserFlag:inputs.NextUserFlag,
		ButtonName:inputs.ButtonName,
		UserID:inputs.UserID,
		UserName:inputs.UserName,
		UsrIDLandNm:inputs.UsrIDLandNm,
		CurUsrInstID:inputs.CurUsrInstID,
		CurUsrInstNm:inputs.CurUsrInstNm,
		NextUserGrp:nextUserGrp,
		Conditions:conditions,
	})
	if e != nil {
		h.JSONR(c, http.StatusExpectationFailed, e)
		return
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type ApiGetProcNextNodeInfoInputs struct {
	ProcessInstID string `json:"processInstID" form:"processInstID"`
	TemplateID string `json:"templateID" form:"templateID"`
	TaskTD string `json:"taskID" form:"taskID"`
	Conditions []*Condition `json:"conditions" form:"conditions"`
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
			Key: cond.Key,
			Value: cond.Value,
		})
	}

	resp, e := service.ProcService.NextNodeInfo(service.NextNodeInfoInputs{
		ProcessInstID: inputs.ProcessInstID,
		TemplateID: inputs.TemplateID,
		TaskTD: inputs.TaskTD,
		Conditions: conditions,
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
	TaskID string `json:"taskID" form:"taskID"`
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
	TemplateID string `json:"templateID" form:"templateID"`
	ProcessInstID string `json:"processInstID" form:"processInstID"`
	TaskID string `json:"TaskID" form:"TaskID"`
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

	resp, e := service.ProcService.GetHistDetailList(inputs.ProcessInstID, inputs.TaskID)
	if e != nil {
		h.JSONR(c, http.StatusExpectationFailed, e)
		return
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}
