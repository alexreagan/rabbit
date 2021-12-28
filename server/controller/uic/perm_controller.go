package uic

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/uic"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetPermListInputs struct {
	Name    string `json:"name" form:"name"`
	CnName  string `json:"cnName" form:"cnName"`
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}

type APIGetPermListOutputs struct {
	List       []*uic.Perm `json:"list"`
	TotalCount int64       `json:"totalCount"`
}

// @Summary 权限列表接口
// @Description
// @Produce json
// @Param request query APIGetPermListInputs true "根据查询条件分页查询权限列表"
// @Success 200 {object} APIGetPermListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/perm/list [get]
func PermList(c *gin.Context) {
	var inputs APIGetPermListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var perms []*uic.Perm
	var totalCount int64
	db := g.Con().Portal.Model(uic.Perm{})
	if inputs.Name != "" {
		db = db.Where("name regexp ?", inputs.Name)
	}
	if inputs.CnName != "" {
		db = db.Where("cn_name regexp ?", inputs.CnName)
	}
	db.Count(&totalCount)

	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&perms)

	resp := &APIGetPermListOutputs{
		List:       perms,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetPermCreateInputs struct {
	ID     int64  `json:"id" form:"id"`
	Name   string `json:"name" form:"name"`
	CnName string `json:"cnName" form:"cnName"`
	Remark string `json:"remark" form:"remark"`
}

type APIGetPermCreateOutputs struct {
	Perm *uic.Perm `json:"perm"`
}

// @Summary 新建权限接口
// @Description
// @Produce json
// @Param request query APIGetPermCreateInputs true "新建权限接口"
// @Success 200 {object} APIGetPermCreateOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/perm/create [post]
func PermCreate(c *gin.Context) {
	var inputs APIGetPermCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	db := g.Con().Portal.Model(uic.Perm{})
	perm := &uic.Perm{
		Name:     inputs.Name,
		CnName:   inputs.CnName,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	db = db.Create(perm)

	resp := &APIGetPermCreateOutputs{
		Perm: perm,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 更新权限接口
// @Description
// @Produce json
// @Param request query APIGetPermCreateInputs true "更新权限接口"
// @Success 200 {object} APIGetPermCreateOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/perm/update [put]
func PermUpdate(c *gin.Context) {
	var inputs APIGetPermCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	perm := &uic.Perm{
		ID:       inputs.ID,
		Name:     inputs.Name,
		CnName:   inputs.CnName,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	db := g.Con().Portal.Model(uic.Perm{})
	db = db.Where("id = ?", inputs.ID).Updates(perm)
	if db.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, db.Error)
		return
	}

	resp := &APIGetPermCreateOutputs{
		Perm: perm,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetPermInfoOutputs struct {
	Perm *uic.Perm `json:"perm"`
}

// @Summary 查看权限接口
// @Description
// @Produce json
// @Param request query string true "查看权限接口"
// @Success 200 {object} APIGetPermInfoOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/perm/info [get]
func PermInfo(c *gin.Context) {
	id := c.Query("id")

	var perm *uic.Perm
	db := g.Con().Portal.Model(uic.Perm{})
	db = db.Where("id = ?", id).Find(&perm)
	if db.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, db.Error)
		return
	}

	resp := &APIGetPermInfoOutputs{
		Perm: perm,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetPermMyselfOutputs struct {
	Perms []*uic.Perm `json:"perms"`
}

// @Summary 查看当前用户的权限
// @Description
// @Produce json
// @Success 200 {object} APIGetPermInfoOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/perm/myself [get]
func PermMyself(c *gin.Context) {
	user, _ := h.GetUser(c)

	var perms []*uic.Perm
	db := g.Con().Portal.Model(uic.Perm{}).Debug()
	db = db.Select("`perm`.*")
	db = db.Joins("left join `role_perm_rel` on `perm`.`id` = `role_perm_rel`.`perm`")
	db = db.Joins("left join `user_role_rel` on `role_perm_rel`.`role` = `user_role_rel`.`role`")
	db = db.Where("`user_role_rel`.`user` = ?", user.ID)
	db = db.Find(&perms)
	if db.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, db.Error)
		return
	}

	resp := &APIGetPermMyselfOutputs{
		Perms: perms,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}
