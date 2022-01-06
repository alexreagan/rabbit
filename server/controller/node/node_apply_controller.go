package node

import (
	"encoding/json"
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/model/uic"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type APIGetNodeApplyRequestListInputs struct {
	Applier string `json:"applier" form:"applier"`
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}

type APIGetNodeApplyRequestListData struct {
	node.NodeApplyRequest
	Nodes []*node.Node `json:"nodes"`
	Tags  []*app.Tag   `json:"tags"`
}

type APIGetNodeApplyRequestListOutputs struct {
	List       []*APIGetNodeApplyRequestListData `json:"list"`
	TotalCount int64                             `json:"totalCount"`
}

// @Summary 机器资源申请列表
// @Description
// @Produce json
// @Param APIGetNodeApplyRequestListInputs query APIGetNodeApplyRequestListInputs true "根据查询条件分页查询机器资源申请列表"
// @Success 200 {object} APIGetNodeApplyRequestListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_apply_request/list [get]
func NodeApplyRequestList(c *gin.Context) {
	var inputs APIGetNodeApplyRequestListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var requests []*node.NodeApplyRequest
	var totalCount int64
	tx := g.Con().Portal.Model(node.NodeApplyRequest{})
	tx = tx.Select("distinct `node_apply_request`.*")
	if inputs.Applier != "" {
		tx = tx.Where("`node_apply_request`.`applier` = ?", inputs.Applier)
	}

	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx = tx.Offset(offset).Limit(limit)
	tx.Find(&requests)

	var data []*APIGetNodeApplyRequestListData
	for _, req := range requests {
		var tags []*app.Tag
		if req.TagIDs != "" {
			var tagIDs []int64
			json.Unmarshal([]byte(req.TagIDs), &tagIDs)

			tx = g.Con().Portal.Model(app.Tag{})
			tx.Where("id in (?)", tagIDs).Find(&tags)
		}

		var nodes []*node.Node
		if req.NodeIDs != "" {
			var nodeIDs []int64
			json.Unmarshal([]byte(req.NodeIDs), &nodeIDs)
			tx = g.Con().Portal.Debug().Model(node.Node{})
			tx.Where("id in (?)", nodeIDs).Find(&nodes)
		}

		data = append(data, &APIGetNodeApplyRequestListData{
			NodeApplyRequest: *req,
			Tags:             tags,
			Nodes:            nodes,
		})
	}

	resp := &APIGetNodeApplyRequestListOutputs{
		List:       data,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 根据请求单ID获取详细信息
// @Description
// @Produce json
// @Param id query int true "根据请求单ID获取详细信息"
// @Success 200 {object} node.Node
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_apply_request/info [get]
func NodeApplyRequestInfo(c *gin.Context) {
	id := c.Query("id")
	req := node.NodeApplyRequest{}
	g.Con().Portal.Model(req).Where("id = ?", id).First(&req)
	h.JSONR(c, req)
	return
}

type APIPostNodeApplyRequestCreateInputs struct {
	ID        int64     `json:"id" form:"id"`
	Name      string    `json:"name" form:"name"`
	CPU       string    `json:"cpu" form:"cpu"`
	Memory    string    `json:"memory" form:"memory"`
	Remark    string    `json:"remark" form:"remark"`
	Applier   string    `json:"applier" form:"applier"`
	ReleaseAt time.Time `json:"releaseAt" form:"releaseAt"`
}

// @Summary 创建机器申请单
// @Description
// @Produce json
// @Param APIPostNodeApplyRequestCreateInputs body APIPostNodeApplyRequestCreateInputs true "创建机器申请单"
// @Success 200 json node.NodeApplyRequest
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_apply_request/create [post]
func NodeApplyRequestCreate(c *gin.Context) {
	var inputs APIPostNodeApplyRequestCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	user, _ := h.GetUser(c)

	var applierName string
	if inputs.Applier != "" {
		var applier uic.User
		g.Con().Uic.Model(uic.User{}).Where("jgyg_user_id = ?", inputs.Applier).First(&applier)
		applierName = applier.CnName
	}

	cpu, _ := strconv.ParseInt(inputs.CPU, 10, 64)
	memory, _ := strconv.ParseInt(inputs.Memory, 10, 64)
	req := node.NodeApplyRequest{
		Name:        inputs.Name,
		CPU:         cpu,
		Memory:      memory,
		Applier:     inputs.Applier,
		ApplierName: applierName,
		Remark:      inputs.Remark,
		Creator:     user.JgygUserID,
		CreatorName: user.CnName,
		CreateAt:    gtime.Now(),
		ReleaseAt:   gtime.NewGTime(inputs.ReleaseAt),
		State:       node.NodeApplyStateSubmitted,
	}

	tx := g.Con().Portal
	if tx = tx.Model(req).Create(&req); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
		return
	}

	h.JSONR(c, h.OKStatus, req)
	return
}

type APIPostNodeApplyRequestAssignInputs struct {
	ID      int64   `json:"id" form:"id"`
	TagIDs  []int64 `json:"tagIDs" form:"tagIDs"`
	NodeIDs []int64 `json:"nodeIDs" form:"nodeIDs"`
	State   string  `json:"state" form:"state"`
}

// @Summary 创建机器申请单处理
// @Description
// @Produce json
// @Param APIPostNodeApplyRequestAssignInputs body APIPostNodeApplyRequestAssignInputs true "创建机器申请单处理"
// @Success 200 json node.NodeApplyRequest
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node_apply_request/assign [put]
func NodeApplyRequestAssign(c *gin.Context) {
	var inputs APIPostNodeApplyRequestAssignInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	user, _ := h.GetUser(c)

	nodeIDArrayBytes, _ := json.Marshal(inputs.NodeIDs)
	tagIDsArrayBytes, _ := json.Marshal(inputs.TagIDs)
	req := &node.NodeApplyRequest{
		Assigner:     user.JgygUserID,
		AssignerName: user.CnName,
		AssignAt:     gtime.Now(),
		State:        node.NodeApplyStateSuccess,
		NodeIDs:      string(nodeIDArrayBytes),
		TagIDs:       string(tagIDsArrayBytes),
	}

	tx := g.Con().Portal.Begin()
	if err := tx.Model(req).Where("id = ?", inputs.ID).Updates(req).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	for _, nodeID := range inputs.NodeIDs {
		for _, tagID := range inputs.TagIDs {
			rel := &node.NodeTagRel{
				Node: nodeID,
				Tag:  tagID,
			}
			if !rel.Existing() {
				if err := tx.Model(rel).Create(rel).Error; err != nil {
					tx.Rollback()
					h.JSONR(c, h.ExpecStatus, err)
					return
				}
			}
		}
	}
	tx.Commit()

	h.JSONR(c, h.OKStatus, req)
	return
}
