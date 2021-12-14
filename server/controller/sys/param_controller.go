package sys

import (
	"errors"
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/sys"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetParamListInputs struct {
	Key     string `json:"key" form:"key"`
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}

type APIGetParamListOutputs struct {
	List       []*sys.Param `json:"list"`
	TotalCount int64        `json:"totalCount"`
}

// @Summary 参数列表接口
// @Description
// @Produce json
// @Param APIGetParamListInputs query APIGetParamListInputs true "根据查询条件分页查询参数列表"
// @Success 200 {object} APIGetParamListOutputs
// @Failure 400 {object} APIGetParamListOutputs
// @Router /api/v1/param/list [get]
func ParamList(c *gin.Context) {
	var inputs APIGetParamListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var params []*sys.Param
	var totalCount int64
	db := g.Con().Portal.Debug().Model(sys.Param{})
	db = db.Where("deleted = ?", false)

	db.Count(&totalCount)
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&params)

	resp := &APIGetParamListOutputs{
		List:       params,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 参数详情
// @Description
// @Produce json
// @Param id query string true "param id"
// @Success 200 {object} sys.Param
// @Failure 400 json error
// @Router /api/v1/param/info [get]
func ParamInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		h.JSONR(c, h.BadStatus, errors.New("parameter id is required"))
		return
	}

	var param *sys.Param
	db := g.Con().Portal.Debug().Model(sys.Param{})
	db.Where("id = ?", id).Find(&param)

	h.JSONR(c, http.StatusOK, param)
	return
}

type APIPostParamCreateInputs struct {
	ID     int64  `json:"id" form:"id"`
	Key    string `json:"key" form:"key"`
	Value  string `json:"value" form:"value"`
	Remark string `json:"remark" form:"remark"`
}

// @Summary 创建新参数
// @Description
// @Produce json
// @Param APIPostParamCreateInputs body APIPostParamCreateInputs true "创建新参数"
// @Success 200 json sys.Param
// @Failure 400 json errors
// @Router /api/v1/param/create [post]
func ParamCreate(c *gin.Context) {
	var inputs APIPostParamCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	param := sys.Param{
		Key:      inputs.Key,
		Value:    inputs.Value,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	tx := g.Con().Portal
	if dt := tx.Model(sys.Param{}).Create(&param); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	h.JSONR(c, h.OKStatus, param)
	return
}

// @Summary 更新参数
// @Description
// @Produce json
// @Param APIPostParamCreateInputs body APIPostParamCreateInputs true "更新参数"
// @Success 200 json sys.Param
// @Failure 400 json errors
// @Router /api/v1/param/update [put]
func ParamUpdate(c *gin.Context) {
	var inputs APIPostParamCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	param := sys.Param{
		ID:       inputs.ID,
		Key:      inputs.Key,
		Value:    inputs.Value,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	tx := g.Con().Portal
	if dt := tx.Model(sys.Param{}).Where("id = ?", inputs.ID).Updates(&param); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	h.JSONR(c, h.OKStatus, param)
	return
}
