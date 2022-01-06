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

type APIGetNodeGroupListInputs struct {
	Name  string `json:"name" form:"name"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

type APIGetNodeGroupListOutputs struct {
	List       []*node.NodeGroup `json:"list"`
	TotalCount int64             `json:"totalCount"`
}

// @Summary 获取node group列表
// @Description
// @Produce json
// @Param APIGetNodeGroupListInputs query APIGetNodeGroupListInputs true "获取node group列表"
// @Success 200 {object} APIGetNodeGroupListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_group/list [get]
func NodeGroupList(c *gin.Context) {
	var inputs APIGetNodeGroupListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}
	var nodeGroups []*node.NodeGroup
	var totalCount int64
	tx := g.Con().Portal.Model(node.NodeGroup{})
	if inputs.Name != "" {
		tx.Where("path regexp ?", inputs.Name)
	}
	tx.Count(&totalCount)
	tx.Offset(offset).Limit(limit).Find(&nodeGroups)

	resp := &APIGetNodeGroupListOutputs{
		List:       nodeGroups,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 获取所有node group信息
// @Description
// @Produce json
// @Success 200 {object} APIGetNodeGroupListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_group/all [get]
func NodeGroupAll(c *gin.Context) {
	var nodeGroups []*node.NodeGroup
	var totalCount int64
	tx := g.Con().Portal.Table(node.NodeGroup{}.TableName()).Find(&nodeGroups)
	tx.Count(&totalCount)

	resp := &APIGetNodeGroupListOutputs{
		List:       nodeGroups,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIPostCreateNodeGroup struct {
	Name          string `json:"name" form:"name" binding:"required"`
	Type          string `json:"type" form:"type" binding:"required"`
	CaasServiceID int64  `json:"caasServiceId" form:"caasServiceId"`
	ParentName    string `json:"parentName" form:"parentName"`
	ParentID      int64  `json:"parentID" form:"parentID"`
	Desc          string `json:"desc" form:"desc"`
}

// @Summary 创建node group
// @Description
// @Produce json
// @Param APIGetNodeGroupListInputs query APIGetNodeGroupListInputs true "更新node group"
// @Success 200 {object} APIGetNodeGroupListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_group/create [post]
func NodeGroupCreate(c *gin.Context) {
	var inputs APIPostCreateNodeGroup
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	user, _ := h.GetUser(c)

	// get group path
	pGroup := node.NodeGroup{}
	pGroup.ID = inputs.ParentID
	pGroupPathArray := pGroup.GetPath()
	pGroupPathArray = append(pGroupPathArray, inputs.Name)
	nodeGroupPathArrayBytes, _ := json.Marshal(pGroupPathArray)

	// create node group
	nodeGroup := node.NodeGroup{
		Name:          inputs.Name,
		Type:          inputs.Type,
		ParentName:    inputs.ParentName,
		ParentID:      inputs.ParentID,
		CaasServiceID: inputs.CaasServiceID,
		Path:          strings.Join(pGroupPathArray, node.GroupPathSeperator),
		PathArray:     string(nodeGroupPathArrayBytes),
		Desc:          inputs.Desc,
		CreateUser:    user.UserName,
	}
	if tx := g.Con().Portal.Create(&nodeGroup); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
		return
	}
	nodeGroup.ReBuildTree()

	h.JSONR(c, nodeGroup)
	return
}

type APIBindNodeToNodeGroupInput struct {
	NodeID  int64 `json:"nodeID" form:"nodeID" binding:"required"`
	GroupID int64 `json:"groupID" form:"groupID"  binding:"required"`
}

// @Summary 绑定node到node group
// @Description
// @Produce json
// @Param id query int64 true "绑定node到node group"
// @Success 200 {object} node.NodeGroupRel
// @Failure 417 "internal error"
// @Router /api/v1/node_group/bind_node [post]
func BindNodeToNodeGroup(c *gin.Context) {
	var inputs APIBindNodeToNodeGroupInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Model(node.NodeGroupRel{}).Begin()
	if err := tx.Where("node_id = ?", inputs.NodeID).Delete(&node.NodeGroupRel{}).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}
	rel := &node.NodeGroupRel{
		NodeID:  inputs.NodeID,
		GroupID: inputs.GroupID,
	}
	if err := tx.Create(rel).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		tx.Rollback()
		return
	}
	tx.Commit()

	h.JSONR(c, h.OKStatus, rel)
	return
}

type APIGetnodeGroupGetInputs struct {
	ID int64 `json:"id" form:"id"`
}

// @Summary 根据机器群组名称获取机器群组详细信息
// @Description
// @Produce json
// @Param id query int true "根据机器群组名称获取机器群组详细信息"
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_group/get [get]
func NodeGroupGet(c *gin.Context) {
	var inputs APIGetnodeGroupGetInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	nodeGroup := node.NodeGroup{}
	if tx := g.Con().Portal.Where("id = ?", inputs.ID).Find(&nodeGroup); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
		return
	}

	h.JSONR(c, h.OKStatus, nodeGroup)
	return
}

type APIPutNodeGroupInputs struct {
	ID int64 `json:"id" form:"id" binding:"required"`
	APIPostCreateNodeGroup
}

// @Summary 更新node group
// @Description
// @Produce json
// @Param APIPutNodeGroupInputs query APIPutNodeGroupInputs true "更新node group"
// @Success 200 {object} node.NodeGroup
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_group/update [put]
func NodeGroupPut(c *gin.Context) {
	var inputs APIPutNodeGroupInputs
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

	nodeGroup := node.NodeGroup{}
	nodeGroup.ID = inputs.ID
	if tx := g.Con().Portal.Find(&nodeGroup); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
		return
	}

	// get group path
	pGroup := node.NodeGroup{}
	pGroup.ID = inputs.ParentID
	pGroupPathArray := pGroup.GetPath()
	pGroupPathArray = append(pGroupPathArray, inputs.Name)
	nodeGroupPathArrayBytes, _ := json.Marshal(pGroupPathArray)

	// update attr
	nodeGroup.Name = inputs.Name
	nodeGroup.Type = inputs.Type
	nodeGroup.CaasServiceID = inputs.CaasServiceID
	nodeGroup.ParentName = inputs.ParentName
	nodeGroup.ParentID = inputs.ParentID
	nodeGroup.Desc = inputs.Desc
	nodeGroup.Path = strings.Join(pGroupPathArray, node.GroupPathSeperator)
	nodeGroup.PathArray = string(nodeGroupPathArrayBytes)
	nodeGroup.CreateUser = user.UserName
	tx := g.Con().Portal.Model(&nodeGroup).Updates(&nodeGroup)
	if tx.Error != nil {
		h.JSONR(c, h.BadStatus, tx.Error)
		return
	}

	// 调整叶子节点
	nodeGroup.UpdateChildrenPath()

	h.JSONR(c, h.OKStatus, nodeGroup)
	return
}

// @Summary 删除node group
// @Description
// @Produce json
// @Param id path int64 true "删除node group"
// @Success 200 {object} node.NodeGroup
// @Failure 417 "internal error"
// @Router /api/v1/node_group/delete [post]
func NodeGroupDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	nodeGroup := node.NodeGroup{}
	nodeGroup.ID = id
	if tx := g.Con().Portal.Delete(&nodeGroup); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
		return
	}

	// get group path
	var rel []*node.NodeGroupRel
	tx := g.Con().Portal.Model(node.NodeGroupRel{}).Where("node_id = ?", id).Delete(&rel)
	if tx.Error != nil {
		h.JSONR(c, h.BadStatus, tx.Error)
		return
	}

	h.JSONR(c, h.OKStatus, nodeGroup)
	return
}

// @Summary 根据group路径获取group下所有的node
// @Description
// @Produce json
// @Param path query string true "group路径"
// @Success 200 {object} node.NodeGroup
// @Failure 417 "internal error"
// @Router /api/v1/node_group/related_nodes [get]
func NodeGroupRelatedNodes(c *gin.Context) {
	groupPath := c.Query("path")
	if groupPath == "" {
		h.JSONR(c, h.BadStatus, "param path is required!")
		return
	}

	nodeGroup := node.NodeGroup{}
	nodeGroup.Path = groupPath
	if tx := g.Con().Portal.Debug().Model(nodeGroup).Where(&nodeGroup).Find(&nodeGroup); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
		return
	}

	h.JSONR(c, h.OKStatus, nodeGroup.RelatedNodes())
	return
}
