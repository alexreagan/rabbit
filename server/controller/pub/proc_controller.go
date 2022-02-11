package pub

//
//import (
//	h "github.com/alexreagan/rabbit/server/helper"
//	"github.com/alexreagan/rabbit/server/service"
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//type Condition struct {
//	Key   string `json:"key"`
//	Value string `json:"value"`
//}
//
//type ApiGetProcNextNodeInfoInputs struct {
//	ProcessInstID string       `json:"processInstID" form:"processInstID"`
//	TemplateID    string       `json:"templateID" form:"templateID"`
//	TaskTD        string       `json:"taskID" form:"taskID"`
//	Conditions    []*Condition `json:"conditions" form:"conditions"`
//}
//
//type ApiGetProcNextNodeInfoOutputs struct {
//}
//
//// @Summary 查询下一环节信息
//// @Description
//// @Produce json
//// @Param ApiGetProcNextNodeInfoInputs body ApiGetProcNextNodeInfoInputs true "查询下一环节信息"
//// @Success 200 {object} ApiGetProcNextNodeInfoOutputs
//// @Failure 400 "bad arguments"
//// @Failure 417 "internal error"
//// @Router /api/v1/proc/nextNodeInfo [get]
//func NextNodeInfo(c *gin.Context) {
//	var inputs ApiGetProcNextNodeInfoInputs
//
//	if err := c.Bind(&inputs); err != nil {
//		h.JSONR(c, h.BadStatus, err)
//		return
//	}
//
//	conditions := make([]*service.Condition, 0)
//	for _, cond := range conditions {
//		conditions = append(conditions, &service.Condition{
//			Key:   cond.Key,
//			Value: cond.Value,
//		})
//	}
//
//	resp, e := service.ProcService.NextNodeInfo(service.NextNodeInfoInputs{
//		ProcessInstID: inputs.ProcessInstID,
//		TemplateID:    inputs.TemplateID,
//		TaskTD:        inputs.TaskTD,
//		Conditions:    conditions,
//	})
//	if e != nil {
//		h.JSONR(c, http.StatusExpectationFailed, e)
//		return
//	}
//	h.JSONR(c, http.StatusOK, resp)
//	return
//}
//
//type ApiGetProcGetPersonByNodeInputs struct {
//	TemplateID string `json:"templateID" form:"templateID"`
//	TaskID     string `json:"taskID" form:"taskID"`
//}
//
//// @Summary 根据节点获取审批用户
//// @Description
//// @Produce json
//// @Param ApiGetProcNextNodeInfoInputs body ApiGetProcNextNodeInfoInputs true "根据节点获取审批用户"
//// @Success 200 {object} ApiGetProcNextNodeInfoOutputs
//// @Failure 400 "bad arguments"
//// @Failure 417 "internal error"
//// @Router /api/v1/proc/getPersonByNode [get]
//func GetPersonByNode(c *gin.Context) {
//	var inputs ApiGetProcGetPersonByNodeInputs
//
//	if err := c.Bind(&inputs); err != nil {
//		h.JSONR(c, h.BadStatus, err)
//		return
//	}
//
//	resp, e := service.ProcService.GetPersonByNode(inputs.TemplateID, inputs.TaskID)
//	if e != nil {
//		h.JSONR(c, http.StatusExpectationFailed, e)
//		return
//	}
//	h.JSONR(c, http.StatusOK, resp)
//	return
//}
//
//type ApiGetProcGetHistDetailListInputs struct {
//	TemplateID    string `json:"templateID" form:"templateID"`
//	ProcessInstID string `json:"processInstID" form:"processInstID"`
//	TaskID        string `json:"taskID" form:"taskID"`
//	BelongInstID  string `json:"BelongInstID" form:"BelongInstID"`
//	SelectMode    string `json:"SelectMode" form:"SelectMode"`
//}
//
//// @Summary 获取审批历史
//// @Description
//// @Produce json
//// @Param ApiGetProcGetHistDetailListInputs body ApiGetProcGetHistDetailListInputs true "获取审批历史"
//// @Success 200 {object} ApiGetProcNextNodeInfoOutputs
//// @Failure 400 "bad arguments"
//// @Failure 417 "internal error"
//// @Router /api/v1/proc/getHistDetailList [get]
//func GetHistDetailList(c *gin.Context) {
//	var inputs ApiGetProcGetHistDetailListInputs
//
//	if err := c.Bind(&inputs); err != nil {
//		h.JSONR(c, h.BadStatus, err)
//		return
//	}
//
//	resp, e := service.ProcService.GetHistDetailList(inputs.ProcessInstID, inputs.TaskID, inputs.BelongInstID, inputs.SelectMode)
//	if e != nil {
//		h.JSONR(c, http.StatusExpectationFailed, e)
//		return
//	}
//	h.JSONR(c, http.StatusOK, resp)
//	return
//}
