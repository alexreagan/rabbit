package app

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetTemplateListInputs struct {
	Name    string `json:"name" form:"name"`
	Remark  string `json:"remark" form:"remark"`
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}

type APIGetTemplateListOutputs struct {
	List       []*app.Template `json:"list"`
	TotalCount int64           `json:"totalCount"`
}

// @Summary 展现模板列表接口
// @Description
// @Produce json
// @Param APIGetTemplateListInputs query APIGetTemplateListInputs true "展现模板列表接口"
// @Success 200 {object} APIGetTemplateListOutputs
// @Failure 400 "error"
// @Router /api/v1/template/list [get]
func TemplateList(c *gin.Context) {
	var inputs APIGetTemplateListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var templates []*app.Template
	var totalCount int64
	db := g.Con().Portal.Debug().Model(app.Template{})
	if inputs.Name != "" {
		db = db.Where("`template`.`name` = ?", inputs.Name)
	}
	if inputs.Remark != "" {
		db = db.Where("`template`.`remark` regexp ?", inputs.Remark)
	}

	db.Count(&totalCount)
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&templates)

	resp := &APIGetTemplateListOutputs{
		List:       templates,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 根据机器ID获取模板详细信息
// @Description
// @Produce json
// @Param id query int true "根据机器ID获取模板详细信息"
// @Success 200 {object} app.Template
// @Failure 400 {object} app.Template
// @Router /api/v1/template/info [get]
func TemplateInfo(c *gin.Context) {
	id := c.Query("id")
	template := app.Template{}
	g.Con().Portal.Model(template).Where("id = ?", id).First(&template)
	h.JSONR(c, template)
	return
}

type APIPostTemplateUpdateInputs struct {
	ID     int64  `json:"id" form:"id"`
	Name   string `json:"name" form:"name"`
	Remark string `json:"remark" form:"remark"`
	State  string `json:"state" form:"state"`
}

// @Summary 创建新模板
// @Description
// @Produce json
// @Param APIPostTemplateUpdateInputs body APIPostTemplateUpdateInputs true "创建新模板"
// @Success 200 {object} APIPostTemplateUpdateInputs
// @Failure 400 {object} APIPostTemplateUpdateInputs
// @Router /api/v1/template/create [post]
func TemplateCreate(c *gin.Context) {
	var inputs APIPostTemplateUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	user, _ := h.GetUser(c)

	tx := g.Con().Portal
	template := app.Template{
		Name:     inputs.Name,
		Remark:   inputs.Remark,
		State:    inputs.State,
		Creator:  user.JgygUserId,
		CreateAt: gtime.Now(),
	}
	if dt := tx.Model(app.Template{}).Create(&template); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	h.JSONR(c, h.OKStatus, template)
	return
}

// @Summary 更新模板信息
// @Description
// @Produce json
// @Param APIPostTemplateUpdateInputs body APIPostTemplateUpdateInputs true "更新模板信息"
// @Success 200 {object} app.Template
// @Failure 400 "bad parameters"
// @Failure 417 "internal error"
// @Router /api/v1/template/update [put]
func TemplateUpdate(c *gin.Context) {
	var inputs APIPostTemplateUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal
	template := app.Template{}
	if dt := tx.Model(app.Template{}).Where("id = ?", inputs.ID).Find(&template); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	user, _ := h.GetUser(c)
	template.Name = inputs.Name
	template.Remark = inputs.Remark
	template.State = inputs.State
	template.Creator = user.JgygUserId
	template.UpdateAt = gtime.Now()
	if dt := tx.Model(app.Template{}).Where("id = ?", inputs.ID).Updates(template); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	h.JSONR(c, h.OKStatus, template)
	return
}
