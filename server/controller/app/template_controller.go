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
	ID     int64          `json:"id" form:"id"`
	Edges  []*app.G6Edge  `json:"edges" form:"edges"`
	Groups []*app.G6Group `json:"groups" form:"groups"`
	Nodes  []*app.G6Node  `json:"nodes" form:"nodes"`
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
		Groups: inputs.Groups,
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

	//// template edge
	//var deletedTemplateEdges []*app.TemplateEdge
	//dt := db.Model(app.TemplateEdge{}).Debug()
	//dt = dt.Where("template = ?", inputs.ID)
	//if dt = dt.Delete(&deletedTemplateEdges); dt.Error != nil {
	//	dt.Rollback()
	//	h.JSONR(c, h.ExpecStatus, dt.Error)
	//	return
	//}
	//for _, edge := range inputs.Edges {
	//	edgeEnd, _ := json.Marshal(edge.End)
	//	edgeEndPoint, _ := json.Marshal(edge.EndPoint)
	//	edgeStart, _ := json.Marshal(edge.Start)
	//	edgeStartPoint, _ := json.Marshal(edge.StartPoint)
	//	if dt = db.Model(app.TemplateEdge{}).Create(&app.TemplateEdge{
	//		Template:   inputs.ID,
	//		End:        string(edgeEnd),
	//		EndPoint:   string(edgeEndPoint),
	//		G6Edge:       edge.ID,
	//		Shape:      edge.Shape,
	//		Source:     edge.Source,
	//		SourceID:   edge.SourceID,
	//		Start:      string(edgeStart),
	//		StartPoint: string(edgeStartPoint),
	//		Target:     edge.Target,
	//		TargetID:   edge.TargetID,
	//		Type:       edge.Type,
	//	}); dt.Error != nil {
	//		dt.Rollback()
	//		h.JSONR(c, h.ExpecStatus, dt.Error)
	//		return
	//	}
	//}
	//
	//// template node
	//var deletedTemplateNodes []*app.TemplateNode
	//dt = db.Model(app.TemplateNode{})
	//dt = dt.Where("template = ?", inputs.ID)
	//if dt = dt.Delete(&deletedTemplateNodes); dt.Error != nil {
	//	dt.Rollback()
	//	h.JSONR(c, h.ExpecStatus, dt.Error)
	//	return
	//}
	//for _, node := range inputs.Nodes {
	//	nodeSize, _ := json.Marshal(node.Size)
	//	nodeInPoints, _ := json.Marshal(node.InPoints)
	//	nodeOutPoints, _ := json.Marshal(node.OutPoints)
	//	if dt = db.Model(app.TemplateNode{}).Create(&app.TemplateNode{
	//		Template:   inputs.ID,
	//		ID:         node.ID,
	//		Name:       node.Name,
	//		Label:      node.Label,
	//		Size:       string(nodeSize),
	//		Type:       node.Type,
	//		X:          node.X,
	//		Y:          node.Y,
	//		Shape:      node.Shape,
	//		Color:      node.Color,
	//		Image:      node.Image,
	//		StateImage: node.StateImage,
	//		OffsetX:    node.OffsetX,
	//		OffsetY:    node.OffsetY,
	//		InPoints:   string(nodeInPoints),
	//		OutPoints:  string(nodeOutPoints),
	//	}); dt.Error != nil {
	//		dt.Rollback()
	//		h.JSONR(c, h.ExpecStatus, dt.Error)
	//		return
	//	}
	//}
	//
	//// template group
	//var deletedTemplateGroups []*app.TemplateGroup
	//dt = db.Model(app.TemplateGroup{})
	//dt = dt.Where("template = ?", inputs.ID)
	//if dt = dt.Delete(&deletedTemplateGroups); dt.Error != nil {
	//	dt.Rollback()
	//	h.JSONR(c, h.ExpecStatus, dt.Error)
	//	return
	//}
	//for _, grp := range inputs.Groups {
	//	if dt = db.Model(app.TemplateGroup{}).Create(grp); dt.Error != nil {
	//		dt.Rollback()
	//		h.JSONR(c, h.ExpecStatus, dt.Error)
	//		return
	//	}
	//}

	h.JSONR(c, h.OKStatus, template)
	return
}

type APIGetV3TreeInputs struct {
	TagIDs []int64 `json:"tagIDs[]" form:"tagIDs[]"`
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

	globalTemplateGraphNode := service.TemplateService.GlobalTagGraphNodeV3()
	if globalTemplateGraphNode == nil {
		template, _ := service.TemplateService.ValidTemplate()
		g6Graph, _ := service.TemplateService.UnSerialize(template.Content)
		globalTemplateGraphNode = service.TemplateService.BuildGraphV3(g6Graph)
	}

	// 根节点
	if len(inputs.TagIDs) == 0 {
		h.JSONR(c, http.StatusOK, globalTemplateGraphNode)
		return
	}

	// 其他节点
	// 找到inputs的末级节点
	graphNode := globalTemplateGraphNode
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

	globalTemplateGraphNode := service.TemplateService.GlobalTagGraphNodeV3()
	if globalTemplateGraphNode == nil {
		template, _ := service.TemplateService.ValidTemplate()
		g6Graph, _ := service.TemplateService.UnSerialize(template.Content)
		globalTemplateGraphNode = service.TemplateService.BuildGraphV3(g6Graph)
	}

	// 根节点
	if len(inputs.TagIDs) == 0 {
		var resp []*service.TagGraphNode
		for _, n := range globalTemplateGraphNode.Nexts() {
			resp = append(resp, n)
		}
		h.JSONR(c, http.StatusOK, resp)
		return
	}

	// 其他节点
	// 找到inputs的末级节点
	graphNode := globalTemplateGraphNode
	for _, id := range inputs.TagIDs {
		graphNode = graphNode.Next[id]
	}

	// 末级节点的子节点
	var resp []interface{}
	// 子节点
	for _, x := range graphNode.Nexts() {
		resp = append(resp, x)
	}

	// 未打到子标签的host
	for _, x := range graphNode.UnTaggedHosts {
		resp = append(resp, x)
	}

	// 未打到子标签的pod
	for _, x := range graphNode.UnTaggedPods {
		resp = append(resp, x)
	}

	h.JSONR(c, h.OKStatus, resp)
	return
}
