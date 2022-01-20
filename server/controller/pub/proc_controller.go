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
