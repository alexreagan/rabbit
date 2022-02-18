package wfe

import (
	"errors"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/gin-gonic/gin"
)

type APIPostCreateInputs struct {
	TemplateID  string                 `json:"templateID" form:"templateID"`
	TaskID      string                 `json:"taskID" form:"taskID"`
	Remark      string                 `json:"remark" form:"remark"`
	ButtonName  string                 `json:"buttonName" form:"buttonName"`
	PrjID       string                 `json:"prjID" form:"prjID"`
	PrjSN       string                 `json:"prjSN" form:"prjSN"`
	TodoTmTpCd  string                 `json:"todoTmTpCd" form:"todoTmTpCd,omitempty"`
	TodoTmTtl   string                 `json:"todoTmTtl" form:"todoTmTtl"`
	BlngInstID  string                 `json:"blngInstID" form:"blngInstID,omitempty"`
	DmnGrpID    string                 `json:"dmnGrpID" form:"dmnGrpID,omitempty"`
	NextUserGrp []string               `json:"nextUserGrp" form:"nextUserGrp,omitempty"`
	NextNode    service.TXNextNodeInfo `json:"nextNode" form:"nextNode,omitempty"`
}

// @Summary 创建流程A0902S102
// @Description
// @Produce json
// @Param APIPostCreateInputs body APIPostCreateInputs true "创建流程"
// @Success 200 {object} service.WfeCreateResponse
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/wfe/create [post]
func Create(c *gin.Context) {
	var inputs APIPostCreateInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	u, _ := h.GetUser(c)
	resp, err := service.WfeService.Create(&u, inputs.TemplateID, inputs.TaskID, inputs.Remark, inputs.ButtonName, inputs.NextUserGrp,
		inputs.NextNode, inputs.PrjID, inputs.PrjSN, inputs.TodoTmTpCd, inputs.TodoTmTtl, inputs.BlngInstID, inputs.DmnGrpID)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, resp)
	return
}

type APIPostExecuteInputs struct {
	ProcessInstID string                 `json:"processInstID" form:"processInstID"`
	TemplateID    string                 `json:"templateID" form:"templateID"`
	TaskID        string                 `json:"taskID" form:"taskID"`
	Remark        string                 `json:"remark" form:"remark"`
	OpinCode      string                 `json:"opinCode" form:"opinCode,omitempty"`
	OpinDesc      string                 `json:"opinDesc" form:"opinDesc,omitempty"`
	ExeFstTask    string                 `json:"exeFstTask" form:"exeFstTask,omitempty"`
	ButtonName    string                 `json:"buttonName" form:"buttonName"`
	PrjID         string                 `json:"prjID" form:"prjID"`
	PrjSN         int64                  `json:"prjSN" form:"prjSN"`
	NextUserGrp   []string               `json:"nextUserGrp" form:"nextUserGrp,omitempty"`
	NextNode      service.TXNextNodeInfo `json:"nextNode" form:"nextNode,omitempty"`
	Conditions    []service.Condition    `json:"conditions" form:"conditions,omitempty"`
	//TodoTmTpCd    string               `json:"todoTmTpCd" form:"todoTmTpCd,omitempty"`
	//TodoTmTtl     string               `json:"todoTmTtl" form:"todoTmTtl"`
	//BlngInstID    string               `json:"blngInstID" form:"blngInstID,omitempty"`
	//DmnGrpID      string               `json:"dmnGrpID" form:"dmnGrpID,omitempty"`
}

// @Summary 处理流程A0902S102
// @Description
// @Produce json
// @Param APIPostExecuteInputs body APIPostExecuteInputs true "处理流程"
// @Success 200 {object} service.WfeExecuteResponse
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/wfe/execute [post]
func Execute(c *gin.Context) {
	var inputs APIPostExecuteInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	u, _ := h.GetUser(c)
	resp, err := service.WfeService.Execute(&u, inputs.TemplateID, inputs.ProcessInstID, inputs.TaskID, inputs.Remark,
		inputs.OpinCode, inputs.OpinDesc, inputs.ButtonName, inputs.ExeFstTask, inputs.NextUserGrp, inputs.NextNode, inputs.Conditions)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, resp)
	return
}

type SortType struct {
	NameP8    string `json:"nameP8" form:"nameP8"`
	Direction string `json:"direction" form:"direction"`
}

type APIPostTodoListInputs struct {
	//ProcessInstIDList string   `json:"processInstIDList" form:"processInstIDList"`
	//PrjBelongTypeList string   `json:"prjBelongTypeList" form:"prjBelongTypeList"`
	//AvyOwrNm          string   `json:"avyOwrNm" form:"avyOwrNm"`
	//TodoType          string   `json:"TodoType" form:"TodoType"`
	SortFields  []*service.SortField `json:"sortFields" form:"sortFields"`
	TimeStart   string               `json:"timeStart" form:"timeStart"`
	TimeEnd     string               `json:"timeEnd" form:"timeEnd"`
	PrjID       string               `json:"prjID" form:"prjID"`
	PrjTypeList string               `json:"prjTypeList" form:"prjTypeList"`
	PrjNm       string               `json:"prjNm" form:"prjNm"`
	WfExtrNm    string               `json:"WfExtrNm" form:"WfExtrNm"`
	Page        string               `json:"page" from:"page"`
	Limit       string               `json:"limit" form:"limit"`
}

// @Summary 待办A0902S119
// @Description
// @Produce json
// @Param APIPostTodoListInputs body APIPostTodoListInputs true "待办A0902S119"
// @Success 200 {object} service.WfeTodosResponse
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/wfe/todos [post]
func Todos(c *gin.Context) {
	var inputs APIPostTodoListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	u, _ := h.GetUser(c)
	resp, err := service.WfeService.Todos(&u, inputs.TimeStart, inputs.TimeEnd, inputs.PrjID,
		inputs.PrjTypeList, inputs.PrjNm, inputs.WfExtrNm, inputs.SortFields, inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, resp)
	return
}

type APIPostHistDetailListInputs struct {
	ProcessInstID string `json:"processInstID" form:"processInstID"`
	//SelectMode    string `json:"selectMode" form:"selectMode"`
}

// @Summary 流程处理历史记录A0902S124
// @Description
// @Produce json
// @Param APIPostHistDetailListInputs body APIPostHistDetailListInputs true "处理流程"
// @Success 200 {object} service.WfeExecuteResponse
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/wfe/histDetails [post]
func HistDetails(c *gin.Context) {
	var inputs APIPostHistDetailListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	u, _ := h.GetUser(c)
	resp, err := service.WfeService.HistDetails(&u, inputs.ProcessInstID)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, resp)
	return
}

type APIPostNextNodeInfoInputs struct {
	TemplateID    string `json:"templateID" form:"templateID"`
	ProcessInstID string `json:"processInstID" form:"processInstID"`
	TaskID        string `json:"taskID" form:"taskID"`
}

func (input APIPostNextNodeInfoInputs) checkInputsContain() error {
	if input.TemplateID == "" {
		return errors.New("templateID should not empty")
	}
	return nil
}

// @Summary 下个节点信息A0902S112
// @Description
// @Produce json
// @Param APIPostNextNodeInfoInputs body APIPostNextNodeInfoInputs true "下个节点信息A0902S112"
// @Success 200 {object} service.WfeTodosResponse
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/wfe/nextNodeInfo [post]
func NextNodeInfo(c *gin.Context) {
	var inputs APIPostNextNodeInfoInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	u, _ := h.GetUser(c)
	resp, err := service.WfeService.NextNodeInfo(&u, inputs.TemplateID, inputs.ProcessInstID, inputs.TaskID)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, resp)
	return
}

type APIPostTodo2DoingInputs struct {
	TodoID string `json:"todoID" form:"todoID" binding:"required"`
}

// @Summary 待办转在办A0902S132
// @Description
// @Produce json
// @Param APIPostTodo2DoingInputs body APIPostTodo2DoingInputs true "待办转在办A0902S132"
// @Success 200 {object} service.WfeTodo2DoingResponse
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/wfe/todo2doing [post]
func Todo2Doing(c *gin.Context) {
	var inputs APIPostTodo2DoingInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	u, _ := h.GetUser(c)
	resp, err := service.WfeService.Todo2Doing(&u, inputs.TodoID)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, resp)
	return
}

type APIPostHasDoneListInputs struct {
	SortFields  []*service.SortField `json:"sortFields" form:"sortFields"`
	TimeStart   string               `json:"timeStart" form:"timeStart"`
	TimeEnd     string               `json:"timeEnd" form:"timeEnd"`
	PrjID       string               `json:"prjID" form:"prjID"`
	PrjTypeList string               `json:"prjTypeList" form:"prjTypeList"`
	PrjNm       string               `json:"prjNm" form:"prjNm"`
	WfExtrNm    string               `json:"WfExtrNm" form:"WfExtrNm"`
	Page        string               `json:"page" from:"page"`
	Limit       string               `json:"limit" form:"limit"`
}

// @Summary 已办A0902S120
// @Description
// @Produce json
// @Param APIPostHasDoneListInputs body APIPostHasDoneListInputs true "已办A0902S120"
// @Success 200 {object} service.WfeHasDoneResponse
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/wfe/hasDone [post]
func HasDone(c *gin.Context) {
	var inputs APIPostTodoListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	u, _ := h.GetUser(c)
	resp, err := service.WfeService.HasDone(&u, inputs.TimeStart, inputs.TimeEnd, inputs.PrjID,
		inputs.PrjTypeList, inputs.PrjNm, inputs.WfExtrNm, inputs.SortFields, inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, resp)
	return
}
