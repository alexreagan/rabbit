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

type APIGetRoleListInputs struct {
	Name    string `json:"name" form:"name"`
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}

type APIGetRoleListOutputs struct {
	List       []*uic.Role `json:"list"`
	TotalCount int64       `json:"totalCount"`
}

// @Summary 角色列表接口
// @Description
// @Produce json
// @Param request query APIGetRoleListInputs true "根据查询条件分页查询角色列表"
// @Success 200 {object} APIGetRoleListOutputs
// @Failure 400 {object} APIGetRoleListOutputs
// @Router /api/v1/role/list [get]
func RoleList(c *gin.Context) {
	var inputs APIGetRoleListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var roles []*uic.Role
	var totalCount int64
	db := g.Con().Portal.Model(uic.Role{})
	if inputs.Name != "" {
		db = db.Where("name regexp ?", inputs.Name)
	}
	db.Count(&totalCount)

	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&roles)

	resp := &APIGetRoleListOutputs{
		List:       roles,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetRoleSearchInputs struct {
	Name string `json:"name"`
}

type APIGetRoleSearchOutputs struct {
	List []*uic.Role `json:"list"`
}

func RoleSelect(c *gin.Context) {
	var inputs APIGetRoleSearchInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	db := g.Con().Portal
	var roles []*uic.Role
	db.Model(uic.Role{}).Where("name regexp ?", inputs.Name).Find(&roles)

	resp := &APIGetRoleSearchOutputs{
		List: roles,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetRoleCreateInputs struct {
	ID       int64   `json:"id" form:"id"`
	Name     string  `json:"name" form:"name"`
	CnName   string  `json:"cnName" form:"cnName"`
	Remark   string  `json:"remark" form:"remark"`
	PermList []int64 `json:"permList" form:"permList"`
}

type APIGetRoleCreateOutputs struct {
	Role *uic.Role `json:"role"`
}

// @Summary 新建权限接口
// @Description
// @Produce json
// @Param request query APIGetRoleCreateInputs true "新建权限接口"
// @Success 200 {object} APIGetRoleCreateOutputs
// @Failure 400 {object} APIGetRoleCreateOutputs
// @Router /api/v1/role/create [post]
func RoleCreate(c *gin.Context) {
	var inputs APIGetRoleCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	db := g.Con().Portal.Model(uic.Role{}).Begin()
	role := &uic.Role{
		Name:      inputs.Name,
		CnName:    inputs.CnName,
		Remark:    inputs.Remark,
		CreatedAt: gtime.Now(),
	}
	db = db.Create(role)

	var rels []*uic.RolePermRel
	dt := g.Con().Portal.Model(uic.RolePermRel{})
	dt.Where("id = ?", role.ID).Delete(&rels)

	for _, v := range inputs.PermList {
		dt.Create(&uic.RolePermRel{
			Role: role.ID,
			Perm: v,
		})
	}
	db.Commit()

	resp := &APIGetRoleCreateOutputs{
		Role: role,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 更新角色接口
// @Description
// @Produce json
// @Param request query APIGetRoleCreateInputs true "更新角色接口"
// @Success 200 {object} APIGetRoleCreateOutputs
// @Failure 400 {object} APIGetRoleCreateOutputs
// @Router /api/v1/role/update [put]
func RoleUpdate(c *gin.Context) {
	var inputs APIGetRoleCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	role := &uic.Role{
		ID:     inputs.ID,
		Name:   inputs.Name,
		CnName: inputs.CnName,
		Remark: inputs.Remark,
	}
	db := g.Con().Portal.Model(uic.Role{})
	db = db.Where("id = ?", inputs.ID).Updates(role)
	if db.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, db.Error)
		return
	}

	var rels []*uic.RolePermRel
	dt := g.Con().Portal.Model(uic.RolePermRel{})
	dt = dt.Where("role = ?", role.ID).Delete(&rels)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	}

	for _, v := range inputs.PermList {
		dt.Create(&uic.RolePermRel{
			Role: role.ID,
			Perm: v,
		})
	}

	resp := &APIGetRoleCreateOutputs{
		Role: role,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetRoleInfoOutputs struct {
	Role  *uic.Role   `json:"role"`
	Perms []*uic.Perm `json:"perms"`
}

// @Summary 查看权限接口
// @Description
// @Produce json
// @Param request query string true "查看权限接口"
// @Success 200 {object} APIGetRoleInfoOutputs
// @Failure 400 {object} APIGetRoleInfoOutputs
// @Router /api/v1/role/info [get]
func RoleInfo(c *gin.Context) {
	id := c.Query("id")

	var role *uic.Role
	db := g.Con().Portal.Model(uic.Role{})
	db = db.Where("id = ?", id).Find(&role)
	if db.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, db.Error)
		return
	}

	var perms []*uic.Perm
	dt := g.Con().Portal.Model(uic.Perm{}).Debug()
	dt = dt.Select("`perm`.*")
	dt = dt.Joins("left join `role_perm_rel` on `role_perm_rel`.`perm` = `perm`.`id`")
	dt = dt.Where("`role_perm_rel`.`role` = ?", id)
	dt = dt.Find(&perms)

	resp := &APIGetRoleInfoOutputs{
		Role:  role,
		Perms: perms,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}
