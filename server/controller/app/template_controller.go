package app

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// @Summary 展现模板所有数据
// @Description
// @Produce json
// @Param APIGetTemplateListInputs query APIGetTemplateListInputs true "展现模板所有数据"
// @Success 200 {object} APIGetTemplateListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/template/all [get]
func TemplateAll(c *gin.Context) {
	var inputs APIGetTemplateListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var templates []*app.Template
	var totalCount int64
	tx := g.Con().Portal.Model(app.Template{})
	if inputs.Name != "" {
		tx = tx.Where("`template`.`name` = ?", inputs.Name)
	}
	if inputs.Remark != "" {
		tx = tx.Where("`template`.`remark` regexp ?", inputs.Remark)
	}
	tx = tx.Where("`template`.`state` = ?", "enable")

	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx.Find(&templates)

	resp := &APIGetTemplateListOutputs{
		List:       templates,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 展现模板列表接口
// @Description
// @Produce json
// @Param APIGetTemplateListInputs query APIGetTemplateListInputs true "展现模板列表接口"
// @Success 200 {object} APIGetTemplateListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
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
	tx := g.Con().Portal.Debug().Model(app.Template{})
	if inputs.Name != "" {
		tx = tx.Where("`template`.`name` = ?", inputs.Name)
	}
	if inputs.Remark != "" {
		tx = tx.Where("`template`.`remark` regexp ?", inputs.Remark)
	}

	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx = tx.Offset(offset).Limit(limit)
	tx.Find(&templates)

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
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/template/info [get]
func TemplateInfo(c *gin.Context) {
	id := c.Query("id")
	idx, _ := strconv.ParseInt(id, 10, 64)
	template, err := service.TemplateService.Get(idx)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}
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
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
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
		Creator:  user.JgygUserID,
		CreateAt: gtime.Now(),
		UpdateAt: gtime.Now(),
	}
	if tx = tx.Model(app.Template{}).Create(&template); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
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

	template, err := service.TemplateService.Get(inputs.ID)
	if err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	user, _ := h.GetUser(c)
	template = &app.Template{
		ID:       inputs.ID,
		Name:     inputs.Name,
		Remark:   inputs.Remark,
		State:    inputs.State,
		Creator:  user.JgygUserID,
		CreateAt: gtime.Now(),
		UpdateAt: gtime.Now(),
	}
	if err = service.TemplateService.Updates(template); err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, h.OKStatus, template)
	return
}

type APIPostTemplateDesignInputs struct {
	ID     int64           `json:"id" form:"id"`
	Edges  []*app.G6Edge   `json:"edges" form:"edges"`
	Nodes  []*app.G6Node   `json:"nodes" form:"nodes"`
	Combos []*app.G6Combos `json:"combos" form:"combos"`
}

// @Summary 模板设计
// @Description
// @Produce json
// @Param APIPostTemplateDesignInputs body APIPostTemplateDesignInputs true "模板设计"
// @Success 200 {object} app.Template
// @Failure 400 "bad parameters"
// @Failure 417 "internal error"
// @Router /api/v1/template/design [post]
func TemplateDesign(c *gin.Context) {
	var inputs APIPostTemplateDesignInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	// 模板
	user, _ := h.GetUser(c)
	g6Graph := app.G6Graph{
		Nodes:  inputs.Nodes,
		Edges:  inputs.Edges,
		Combos: inputs.Combos,
	}
	byt, _ := service.TemplateService.Serialize(g6Graph)

	template := &app.Template{
		ID:       inputs.ID,
		Content:  string(byt),
		Creator:  user.JgygUserID,
		UpdateAt: gtime.Now(),
	}
	if err := service.TemplateService.Updates(template); err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	// 更新graph
	service.TemplateService.BuildTemplateGraph(template)
	h.JSONR(c, h.OKStatus, template)
	return
}

type APIGetV3TreeInputs struct {
	TemplateID int64   `json:"templateID" form:"templateID"`
	TagIDs     []int64 `json:"tagIDs[]" form:"tagIDs[]"`
}

// @Summary V3版根据tags路径获取tag信息
// @Description
// @Produce json
// @Param APIGetV3TreeInputs query APIGetV3TreeInputs true "V3版根据tags路径获取tag信息"
// @Success 200 {object} []service.TagGraphNode
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v3/tree/node [get]
func V3TreeNode(c *gin.Context) {
	var inputs APIGetV3TreeInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	// get graph node
	var template *app.Template
	var graphNode *service.TagGraphNode
	if inputs.TemplateID == 0 {
		template, _ = service.TemplateService.ValidTemplate()
	} else {
		template, _ = service.TemplateService.Get(inputs.TemplateID)
	}

	templateGraphNodeMap := service.TemplateService.GlobalTemplateGraphMap()
	if templateGraphNodeMap == nil {
		graphNode = service.TemplateService.BuildTemplateGraph(template)
	} else {
		if _, ok := templateGraphNodeMap[inputs.TemplateID]; !ok {
			graphNode = service.TemplateService.BuildTemplateGraph(template)
		} else {
			graphNode = templateGraphNodeMap[inputs.TemplateID]
		}
	}

	// 根节点
	if len(inputs.TagIDs) == 0 {
		h.JSONR(c, http.StatusOK, graphNode)
		return
	}

	// 其他节点
	// 找到inputs的末级节点
	//graphNode := templateGraphNodeMap[inputs.TemplateID]
	for _, id := range inputs.TagIDs {
		graphNode = graphNode.Next[id]
	}
	h.JSONR(c, http.StatusOK, graphNode)
	return
}

// @Summary V3版根据tags获取tags下所有的机器
// @Description
// @Produce json
// @Param APIGetV3TreeInputs query APIGetV3TreeInputs true "根据tags获取tags下所有的机器"
// @Success 200 {object} []interface{}
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v3/tree/children [get]
func V3TreeChildren(c *gin.Context) {
	var inputs APIGetV3TreeInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	// get graph node
	var template *app.Template
	var graphNode *service.TagGraphNode
	if inputs.TemplateID == 0 {
		template, _ = service.TemplateService.ValidTemplate()
	} else {
		template, _ = service.TemplateService.Get(inputs.TemplateID)
	}

	templateGraphNodeMap := service.TemplateService.GlobalTemplateGraphMap()
	if templateGraphNodeMap == nil {
		graphNode = service.TemplateService.BuildTemplateGraph(template)
	} else {
		if _, ok := templateGraphNodeMap[inputs.TemplateID]; !ok {
			graphNode = service.TemplateService.BuildTemplateGraph(template)
		} else {
			graphNode = templateGraphNodeMap[inputs.TemplateID]
		}
	}

	// 根节点
	if len(inputs.TagIDs) == 0 {
		h.JSONR(c, http.StatusOK, service.BuildChildrenInformation(graphNode))
		return
	}

	// 其他节点
	// 找到inputs的末级节点
	for _, id := range inputs.TagIDs {
		graphNode = graphNode.Next[id]
	}

	h.JSONR(c, h.OKStatus, service.BuildChildrenInformation(graphNode))
	return
}
