package node

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type APIGetHostApplyRequestListInputs struct {
	Applier string `json:"applier" form:"applier"`
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}

type APIGetHostApplyRequestListOutputs struct {
	List       []*node.HostApplyRequest `json:"list"`
	TotalCount int64                    `json:"totalCount"`
}

// @Summary 机器资源申请列表
// @Description
// @Produce json
// @Param APIGetHostApplyRequestListInputs query APIGetHostApplyRequestListInputs true "根据查询条件分页查询机器资源申请列表"
// @Success 200 {object} APIGetHostApplyRequestListOutputs
// @Failure 400 {object} APIGetHostApplyRequestListOutputs
// @Router /api/v1/host_apply_request/list [get]
func HostApplyRequestList(c *gin.Context) {
	var inputs APIGetHostApplyRequestListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var requests []*node.HostApplyRequest
	var totalCount int64
	db := g.Con().Portal.Debug().Model(node.HostApplyRequest{})
	db = db.Select("distinct `host_apply_request`.*")
	if inputs.Applier != "" {
		db = db.Where("`host_apply_request`.`applier` = ?", inputs.Applier)
	}

	db.Count(&totalCount)
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&requests)

	resp := &APIGetHostApplyRequestListOutputs{
		List:       requests,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 根据请求单ID获取详细信息
// @Description
// @Produce json
// @Param id query int true "根据请求单ID获取详细信息"
// @Success 200 {object} node.Host
// @Failure 400 {object} node.Host
// @Router /api/v1/host_apply_request/info [get]
func HostApplyRequestInfo(c *gin.Context) {
	id := c.Query("id")
	req := node.HostApplyRequest{}
	g.Con().Portal.Model(req).Where("id = ?", id).First(&req)
	h.JSONR(c, req)
	return
}

type APIPostHostApplyRequestCreateInputs struct {
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
// @Param APIPostHostApplyRequestCreateInputs formData APIPostHostApplyRequestCreateInputs true "创建机器申请单"
// @Success 200 json node.HostApplyRequest
// @Failure 400 json error
// @Router /api/v1/host_apply_request/create [post]
func HostApplyRequestCreate(c *gin.Context) {
	var inputs APIPostHostApplyRequestCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	user, _ := h.GetUser(c)

	cpu, _ := strconv.ParseInt(inputs.CPU, 10, 64)
	memory, _ := strconv.ParseInt(inputs.Memory, 10, 64)
	req := node.HostApplyRequest{
		Name:      inputs.Name,
		CPU:       cpu,
		Memory:    memory,
		Applier:   inputs.Applier,
		Remark:    inputs.Remark,
		Creator:   user.JgygUserId,
		CreateAt:  gtime.Now(),
		ReleaseAt: gtime.NewGTime(inputs.ReleaseAt),
		State:     node.HostApplyStateSubmitted,
	}

	tx := g.Con().Portal.Debug()
	if dt := tx.Model(req).Create(&req); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	h.JSONR(c, h.OKStatus, req)
	return

}
