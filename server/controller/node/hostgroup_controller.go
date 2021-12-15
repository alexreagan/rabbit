package node

import (
	"encoding/json"
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/node"
	u "github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type APIGetHostGroupListInputs struct {
	Name  string `json:"name" form:"name"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

type APIGetHostGroupListOutputs struct {
	List       []*node.HostGroup `json:"list"`
	TotalCount int64             `json:"totalCount"`
}

// @Summary 获取host group列表
// @Description
// @Produce json
// @Param APIGetHostGroupListInputs query APIGetHostGroupListInputs true "获取host group列表"
// @Success 200 {object} APIGetHostGroupListOutputs
// @Failure 400 {object} APIGetHostGroupListOutputs
// @Router /api/v1/host_group/list [get]
func HostGroupList(c *gin.Context) {
	var inputs APIGetHostGroupListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}
	var hostGroups []*node.HostGroup
	var totalCount int64
	db := g.Con().Portal.Model(node.HostGroup{})
	if inputs.Name != "" {
		db.Where("path regexp ?", inputs.Name)
	}
	db.Count(&totalCount)
	db.Offset(offset).Limit(limit).Find(&hostGroups)

	resp := &APIGetHostGroupListOutputs{
		List:       hostGroups,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 获取所有host group信息
// @Description
// @Produce json
// @Success 200 {object} APIGetHostGroupListOutputs
// @Failure 400 {object} APIGetHostGroupListOutputs
// @Router /api/v1/host_group/all [get]
func HostGroupAll(c *gin.Context) {
	var hostGroups []*node.HostGroup
	var totalCount int64
	db := g.Con().Portal.Table(node.HostGroup{}.TableName()).Find(&hostGroups)
	db.Count(&totalCount)

	resp := &APIGetHostGroupListOutputs{
		List:       hostGroups,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIPostCreateHostGroup struct {
	Name          string `json:"name" form:"name" binding:"required"`
	Type          string `json:"type" form:"type" binding:"required"`
	CaasServiceId int64  `json:"caasServiceId" form:"caasServiceId"`
	ParentName    string `json:"parentName" form:"parentName"`
	ParentId      int64  `json:"parentId" form:"parentId"`
	Desc          string `json:"desc" form:"desc"`
}

// @Summary 创建host group信息
// @Description
// @Produce json
// @Param APIGetHostGroupListInputs query APIGetHostGroupListInputs true "更新host group信息"
// @Success 200 {object} APIGetHostGroupListOutputs
// @Failure 400 {object} APIGetHostGroupListOutputs
// @Router /api/v1/host_group/create [post]
func HostGroupCreate(c *gin.Context) {
	var inputs APIPostCreateHostGroup
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	user, _ := h.GetUser(c)

	// get group path
	pGroup := node.HostGroup{}
	pGroup.ID = inputs.ParentId
	pGroupPathArray := pGroup.GetPath()
	pGroupPathArray = append(pGroupPathArray, inputs.Name)
	hostGroupPathArrayBytes, _ := json.Marshal(pGroupPathArray)

	// create host group
	hostGroup := node.HostGroup{
		Name:          inputs.Name,
		Type:          inputs.Type,
		ParentName:    inputs.ParentName,
		ParentId:      inputs.ParentId,
		CaasServiceId: inputs.CaasServiceId,
		Path:          strings.Join(pGroupPathArray, node.GroupPathSeperator),
		PathArray:     string(hostGroupPathArrayBytes),
		Desc:          inputs.Desc,
		CreateUser:    user.UserName,
	}
	if dt := g.Con().Portal.Create(&hostGroup); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}
	hostGroup.ReBuildTree()

	h.JSONR(c, hostGroup)
	return
}

type APIBindHostToHostGroupInput struct {
	HostId  int64 `json:"hostId" form:"hostId" binding:"required"`
	GroupId int64 `json:"groupId" form:"groupId"  binding:"required"`
}

// @Summary 获取host group信息
// @Description
// @Produce json
// @Param id query int64 true "根据ID获取host group树状信息"
// @Success 200 {object} node.HostGroupRel
// @Failure 417 "internal error"
// @Router /api/v1/host_group/bind_host [post]
func BindHostToHostGroup(c *gin.Context) {
	var inputs APIBindHostToHostGroupInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Model(node.HostGroupRel{}).Begin()
	if dt := tx.Where("host_id = ?", inputs.HostId).Delete(&node.HostGroupRel{}); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		dt.Rollback()
		return
	}
	rel := &node.HostGroupRel{
		HostID:  inputs.HostId,
		GroupID: inputs.GroupId,
	}
	if dt := tx.Debug().Create(rel); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		dt.Rollback()
		return
	}
	tx.Commit()

	h.JSONR(c, h.OKStatus, rel)
	return
}

type APIGetHostGroupGetInputs struct {
	ID int64 `json:"id" form:"id"`
}

// @Summary 根据机器群组名称获取机器群组详细信息
// @Description
// @Produce json
// @Param id query int true "根据机器群组名称获取机器群组详细信息"
// @Success 200 {object} node.HostGroup
// @Failure 400 {object} node.HostGroup
// @Router /api/v1/host_group/get [get]
func HostGroupGet(c *gin.Context) {
	var inputs APIGetHostGroupGetInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	hostGroup := node.HostGroup{}
	if dt := g.Con().Portal.Where("id = ?", inputs.ID).Find(&hostGroup); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	h.JSONR(c, h.OKStatus, hostGroup)
	return
}

type APIPutHostGroupInputs struct {
	ID int64 `json:"id" form:"id" binding:"required"`
	APIPostCreateHostGroup
}

// @Summary 更新host group信息
// @Description
// @Produce json
// @Param APIPutHostGroupInputs query APIPutHostGroupInputs true "更新host group信息"
// @Success 200 {object} node.HostGroup
// @Failure 400 {object} node.HostGroup
// @Router /api/v1/host_group/update [put]
func HostGroupPut(c *gin.Context) {
	var inputs APIPutHostGroupInputs
	err := c.Bind(&inputs)
	switch {
	case err != nil:
		h.JSONR(c, h.BadStatus, err)
		return
	case u.HasDangerousCharacters(inputs.Name):
		h.JSONR(c, h.BadStatus, "group name is invalid")
		return
	}

	user, _ := h.GetUser(c)

	hostGroup := node.HostGroup{}
	hostGroup.ID = inputs.ID
	if dt := g.Con().Portal.Find(&hostGroup); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	// get group path
	pGroup := node.HostGroup{}
	pGroup.ID = inputs.ParentId
	pGroupPathArray := pGroup.GetPath()
	pGroupPathArray = append(pGroupPathArray, inputs.Name)
	hostGroupPathArrayBytes, _ := json.Marshal(pGroupPathArray)

	// update attr
	hostGroup.Name = inputs.Name
	hostGroup.Type = inputs.Type
	hostGroup.CaasServiceId = inputs.CaasServiceId
	hostGroup.ParentName = inputs.ParentName
	hostGroup.ParentId = inputs.ParentId
	hostGroup.Desc = inputs.Desc
	hostGroup.Path = strings.Join(pGroupPathArray, node.GroupPathSeperator)
	hostGroup.PathArray = string(hostGroupPathArrayBytes)
	hostGroup.CreateUser = user.UserName
	dt := g.Con().Portal.Model(&hostGroup).Updates(hostGroup)
	if dt.Error != nil {
		h.JSONR(c, h.BadStatus, dt.Error)
		return
	}

	// 调整叶子节点
	hostGroup.UpdateChildrenPath()

	h.JSONR(c, h.OKStatus, hostGroup)
	return
}

// @Summary 删除host group信息
// @Description
// @Produce json
// @Param id path int64 true "删除host group信息"
// @Success 200 {object} node.HostGroup
// @Failure 417 "internal error"
// @Router /api/v1/host_group/delete [post]
func HostGroupDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	hostGroup := node.HostGroup{}
	hostGroup.ID = id
	if dt := g.Con().Portal.Delete(&hostGroup); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	// get group path
	var rel []*node.HostGroupRel
	dt := g.Con().Portal.Model(node.HostGroupRel{}).Where("host_id = ?", id).Delete(&rel)
	if dt.Error != nil {
		h.JSONR(c, h.BadStatus, dt.Error)
		return
	}

	h.JSONR(c, h.OKStatus, hostGroup)
	return
}

// @Summary 根据group路径获取group下所有的host
// @Description
// @Produce json
// @Param path query string true "group路径"
// @Success 200 {object} node.HostGroup
// @Failure 417 "internal error"
// @Router /api/v1/host_group/related_hosts [get]
func HostGroupRelatedHosts(c *gin.Context) {
	groupPath := c.Query("path")
	if groupPath == "" {
		h.JSONR(c, h.BadStatus, "param path is required!")
		return
	}

	hostGroup := node.HostGroup{}
	hostGroup.Path = groupPath
	if dt := g.Con().Portal.Debug().Model(hostGroup).Where(&hostGroup).Find(&hostGroup); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	h.JSONR(c, h.OKStatus, hostGroup.RelatedHosts())
	return
}
