package alert

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/alert"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetAlertListInputs struct {
	IP             string `json:"ip" form:"ip"`
	PhysicalSystem string `json:"physicalSystem" form:"physicalSystem"`
	Resolved       string `json:"resolved" form:"resolved"`
	Limit          int    `json:"limit" form:"limit"`
	Page           int    `json:"page" form:"page"`
	OrderBy        string `json:"orderBy" form:"orderBy"`
	Order          string `json:"order" form:"order"`
}

type APIGetAlertListOutputs struct {
	List       []*alert.Alert `json:"list"`
	TotalCount int64          `json:"totalCount"`
}

// @Summary 监控报警接口
// @Description
// @Produce json
// @Param APIGetAlertListInputs query APIGetAlertListInputs true "根据查询条件分页查询报警列表"
// @Success 200 {object} APIGetAlertListOutputs
// @Failure 400 {object} APIGetAlertListOutputs
// @Router /api/v1/alert/list [get]
func AlertList(c *gin.Context) {
	var inputs APIGetAlertListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}

	var alerts []*alert.Alert
	var totalCount int64
	db := g.Con().Portal.Debug().Model(alert.Alert{})
	db = db.Select("distinct `alert`.*")
	if inputs.IP != "" {
		db = db.Where("`alert`.`prod_ip` regexp ?", inputs.IP)
	}
	if inputs.PhysicalSystem != "" {
		db = db.Where("`alert`.`sub_sys_en_name` = ?", inputs.PhysicalSystem)
	}
	if inputs.Resolved != "" {
		if inputs.Resolved == "true" {
			db = db.Where("`alert`.`resolved` = 1")
		} else {
			db = db.Where("`alert`.`resolved` = 0")
		}
	}
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db.Count(&totalCount)
	db = db.Offset(offset).Limit(limit)
	db.Find(&alerts)

	resp := &APIGetAlertListOutputs{
		List:       alerts,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 物理子系统类别
// @Description
// @Produce json
// @Success 200 {object} APIGetVariableOutputs
// @Failure 400 {object} APIGetVariableOutputs
// @Router /api/v1/alert/physical_system_choices [get]
func AlertPhysicalSystemChoices(c *gin.Context) {
	var data []*h.APIGetVariableItem
	db := g.Con().Portal.Model(alert.Alert{}).Debug()
	db = db.Select("distinct `sub_sys_name` as `label`, `sub_sys_en_name` as `value`")
	db = db.Order("`sub_sys_en_name`")
	db = db.Find(&data)
	resp := h.APIGetVariableOutputs{
		List:       data,
		TotalCount: int64(len(data)),
	}
	h.JSONR(c, resp)
	return
}
