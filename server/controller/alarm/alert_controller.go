package alarm

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model"
	"github.com/alexreagan/rabbit/server/model/alarm"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetAlarmListInputs struct {
	IP             string `json:"ip" form:"ip"`
	PhysicalSystem string `json:"physicalSystem" form:"physicalSystem"`
	Resolved       string `json:"resolved" form:"resolved"`
	Limit          int    `json:"limit" form:"limit"`
	Page           int    `json:"page" form:"page"`
	OrderBy        string `json:"orderBy" form:"orderBy"`
	Order          string `json:"order" form:"order"`
}

type APIGetAlarmListOutputs struct {
	List       []*alarm.Alarm `json:"list"`
	TotalCount int64          `json:"totalCount"`
}

// @Summary 监控报警接口
// @Description
// @Produce json
// @Param APIGetAlarmListInputs query APIGetAlarmListInputs true "根据查询条件分页查询报警列表"
// @Success 200 {object} APIGetAlarmListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/alarm/list [get]
func List(c *gin.Context) {
	var inputs APIGetAlarmListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}

	var alarms []*alarm.Alarm
	var totalCount int64
	db := g.Con().Portal.Debug().Model(alarm.Alarm{})
	db = db.Select("distinct `alarm`.*")
	if inputs.IP != "" {
		db = db.Where("`alarm`.`prod_ip` regexp ?", inputs.IP)
	}
	if inputs.PhysicalSystem != "" {
		db = db.Where("`alarm`.`sub_sys_en_name` = ?", inputs.PhysicalSystem)
	}
	if inputs.Resolved != "" {
		if inputs.Resolved == "true" {
			db = db.Where("`alarm`.`resolved` = 1")
		} else {
			db = db.Where("`alarm`.`resolved` = 0")
		}
	}
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db.Count(&totalCount)
	db = db.Offset(offset).Limit(limit)
	db.Find(&alarms)

	resp := &APIGetAlarmListOutputs{
		List:       alarms,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 物理子系统类别
// @Description
// @Produce json
// @Success 200 {object} model.APIGetVariableOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/alarm/physical_system_choices [get]
func PhysicalSystemChoices(c *gin.Context) {
	var data []*model.APIGetVariableItem
	db := g.Con().Portal.Model(alarm.Alarm{}).Debug()
	db = db.Select("distinct `sub_sys_name` as `label`, `sub_sys_en_name` as `value`")
	db = db.Order("`sub_sys_en_name`")
	db = db.Find(&data)
	resp := model.APIGetVariableOutputs{
		List:       data,
		TotalCount: int64(len(data)),
	}
	h.JSONR(c, resp)
	return
}
