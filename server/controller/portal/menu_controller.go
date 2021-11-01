package portal

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit/g"
	h "rabbit/server/helper"
	"rabbit/server/model/portal"
	"rabbit/server/model/uic"
)

type APIGetMenuNavListInputs struct {
}

type APIGetMenuNavListOutputs struct {
	Menus       []*portal.Menu    `json:"menus"`
	Permissions []*uic.Permission `json:"permissions"`
}

func MenuNav(c *gin.Context) {
	var inputs APIGetMenuNavListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var permissions []*uic.Permission
	g.Con().Portal.Table(uic.Permission{}.TableName()).Find(&permissions)

	resp := &APIGetMenuNavListOutputs{
		Menus:       portal.Menu{}.BuildTree(),
		Permissions: permissions,
	}

	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetMenuListOutputs struct {
	Menus []*portal.Menu `json:"menus"`
}

func MenuList(c *gin.Context) {
	var inputs APIGetMenuNavListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var menus []*portal.Menu
	g.Con().Portal.Table(portal.Menu{}.TableName()).Find(&menus)

	resp := &APIGetMenuListOutputs{
		Menus: menus,
	}

	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIPostMenuUpdateInputs struct {
	ID       int64  `json:"menuId" form:"menuId" binding:"required"`
	Type     int64  `json:"type" form:"type"`
	Name     string `json:"name" form:"name"`
	Url      string `json:"url" form:"url"`
	ParentId int64  `json:"parentId" form:"parentId"`
	Icon     string `json:"icon" form:"icon"`
}

// @Summary 更新机器信息
// @Description
// @Produce json
// @Param menuId formData string true "根据ID更新菜单信息"
// @Param name formData string false "更新Name"
// @Param Url formData string false "更新Url"
// @Param ParentId formData string false "更新ParentId"
// @Param Icon formData string false "更新Icon"
// @Success 200 {object} APIPostMenuUpdateInputs
// @Failure 400 {object} APIPostMenuUpdateInputs
// @Router /api/v1/menu/update [post]
func MenuUpdate(c *gin.Context) {
	var inputs APIPostMenuUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	g.Con().Portal.Table(portal.Menu{}.TableName()).Where("id = ?", inputs.ID).Updates(inputs)
	h.JSONR(c, http.StatusOK, inputs)
	return
}

// @Summary 根据ID获取菜单信息
// @Description
// @Produce json
// @Param id path int true "根据ID获取菜单信息"
// @Success 200 {object} portal.Menu
// @Failure 400 {object} portal.Menu
// @Router /api/v1/menu/info/:id [get]
func MenuInfo(c *gin.Context) {
	id := c.Param("id")
	menu := portal.Menu{}
	g.Con().Portal.Table(menu.TableName()).Where("id = ?", id).First(&menu)
	h.JSONR(c, menu)
	return
}
