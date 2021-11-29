package sys

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/sys"
	"github.com/alexreagan/rabbit/server/model/uic"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetMenuNavListInputs struct {
}

type APIGetMenuNavListOutputs struct {
	Menus       []*sys.Menu `json:"menus"`
	Permissions []*uic.Perm `json:"permissions"`
}

// @Summary menu列表
// @Description
// @Produce json
// @Success 200 {object} APIGetMenuNavListOutputs
// @Failure 400 json error
// @Router /api/v1/menu/nav [get]
func MenuNav(c *gin.Context) {
	var inputs APIGetMenuNavListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	resp := &APIGetMenuNavListOutputs{
		Menus: sys.Menu{}.BuildTree(),
	}

	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetMenuListOutputs struct {
	Menus []*sys.Menu `json:"menus"`
}

func MenuList(c *gin.Context) {
	var inputs APIGetMenuNavListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var menus []*sys.Menu
	g.Con().Portal.Table(sys.Menu{}.TableName()).Find(&menus)

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

	g.Con().Portal.Table(sys.Menu{}.TableName()).Where("id = ?", inputs.ID).Updates(inputs)
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
	menu := sys.Menu{}
	g.Con().Portal.Table(menu.TableName()).Where("id = ?", id).First(&menu)
	h.JSONR(c, menu)
	return
}
