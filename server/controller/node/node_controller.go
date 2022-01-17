package node

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type OrderBy struct {
	Prop  string `json:"prop" form:"prop"`
	Order string `json:"order" form:"order"`
}

type APIGetNodeListInputs struct {
	IP                     string  `json:"ip" form:"ip"`
	PhysicalSystem         string  `json:"physicalSystem" form:"physicalSystem"`
	CpuCount               int64   `json:"cpuCount" form:"cpuCount"`
	AreaName               string  `json:"areaName" form:"areaName"`
	CpuAvailableUpperLimit float64 `json:"cpuAvailableUpperLimit" form:"cpuAvailableUpperLimit"`
	CpuAvailableLowerLimit float64 `json:"cpuAvailableLowerLimit" form:"cpuAvailableLowerLimit"`
	FsAvailableUpperLimit  float64 `json:"fsAvailableUpperLimit" form:"fsAvailableUpperLimit"`
	FsAvailableLowerLimit  float64 `json:"fsAvailableLowerLimit" form:"fsAvailableLowerLimit"`
	MemAvailableUpperLimit float64 `json:"memAvailableUpperLimit" form:"memAvailableUpperLimit"`
	MemAvailableLowerLimit float64 `json:"memAvailableLowerLimit" form:"memAvailableLowerLimit"`
	//G6Combos                 string  `json:"group" form:"group"`
	//BoundGroup            string  `json:"boundGroup" form:"boundGroup"`
	TagIDs     []int64 `json:"tagIDs[]" form:"tagIDs[]"`
	RelatedTag string  `json:"relatedTag" form:"relatedTag"`
	Status     string  `json:"status" form:"status"`
	Limit      int     `json:"limit" form:"limit"`
	Page       int     `json:"page" form:"page"`
	OrderBy    string  `json:"orderBy" form:"orderBy"`
	Order      string  `json:"order" form:"order"`
}

type APIGetNodeListOutputs struct {
	List       []*node.Node `json:"list"`
	TotalCount int64        `json:"totalCount"`
}

func (input APIGetNodeListInputs) checkInputsContain() error {
	return nil
}

// @Summary 机器列表接口
// @Description
// @Produce json
// @Param APIGetNodeListInputs query APIGetNodeListInputs true "根据查询条件分页查询机器列表"
// @Success 200 {object} APIGetNodeListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/list [get]
func NodeList(c *gin.Context) {
	var inputs APIGetNodeListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var nodes []*node.Node
	var totalCount int64
	tx := g.Con().Portal.Debug().Model(node.Node{})
	tx = tx.Distinct("`node`.*")
	tx = tx.Joins("left join `node_tag_rel` on `node`.id=`node_tag_rel`.`node`")
	if inputs.IP != "" {
		tx = tx.Where("`node`.`ip` regexp ?", inputs.IP)
	}
	if inputs.PhysicalSystem != "" {
		tx = tx.Where("`node`.`physical_system` = ?", inputs.PhysicalSystem)
	}
	if inputs.CpuAvailableLowerLimit != 0 {
		tx = tx.Where("`node`.`cpu_available` >= ?", inputs.CpuAvailableLowerLimit)
	}
	if inputs.CpuAvailableUpperLimit != 0 {
		tx = tx.Where("`node`.`cpu_available` < ?", inputs.CpuAvailableUpperLimit)
	}
	if inputs.FsAvailableLowerLimit != 0 {
		tx = tx.Where("`node`.`file_system_available` >= ?", inputs.FsAvailableLowerLimit)
	}
	if inputs.FsAvailableUpperLimit != 0 {
		tx = tx.Where("`node`.`file_system_available` < ?", inputs.FsAvailableUpperLimit)
	}
	if inputs.MemAvailableLowerLimit != 0 {
		tx = tx.Where("`node`.`mem_available` >= ?", inputs.MemAvailableLowerLimit)
	}
	if inputs.MemAvailableUpperLimit != 0 {
		tx = tx.Where("`node`.`mem_available` < ?", inputs.MemAvailableUpperLimit)
	}
	if inputs.CpuCount != 0 {
		tx = tx.Where("`node`.`cpu_count` = ?", inputs.CpuCount)
	}
	if inputs.AreaName != "" {
		tx = tx.Where("`node`.`area_name` = ?", inputs.AreaName)
	}
	if len(inputs.TagIDs) > 0 {
		var tIDs []int
		for _, i := range inputs.TagIDs {
			tIDs = append(tIDs, int(i))
		}
		sort.Ints(tIDs)

		var tmp []string
		for _, i := range tIDs {
			tmp = append(tmp, strconv.Itoa(i))
		}
		tx = tx.Where("`node_tag_rel`.`tag` in (?)", inputs.TagIDs)
		tx = tx.Group("`node_tag_rel`.`node`")
		tx = tx.Having("group_concat(`node_tag_rel`.`tag` order by `node_tag_rel`.`tag`) = ?", strings.Join(tmp, ","))
	} else {
		tx = tx.Group("`node`.`ip`")
	}
	if inputs.RelatedTag != "" {
		if inputs.RelatedTag == "related" {
			tx = tx.Where("`node_tag_rel`.`tag` is not null")

		} else if inputs.RelatedTag == "unrelated" {
			tx = tx.Where("`node_tag_rel`.`tag` is null")
		}
	}

	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx = tx.Offset(offset).Limit(limit)
	tx.Find(&nodes)

	for _, n := range nodes {
		n.Tags = n.RelatedTags()
	}

	resp := &APIGetNodeListOutputs{
		List:       nodes,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetNodeSelectInputs struct {
	Query string `json:"query" form:"query"`
}

type APIGetNodeSelectOutputs struct {
	List []*node.Node `json:"list"`
}

// @Summary 机器查找
// @Description
// @Produce json
// @Param APIGetNodeSelectInputs query APIGetNodeSelectInputs true "机器查找"
// @Success 200 {object} APIGetNodeListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/select [get]
func NodeSelect(c *gin.Context) {
	var inputs APIGetNodeSelectInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	var nodes []*node.Node
	tx := g.Con().Portal.Debug().Model(node.Node{})
	tx = tx.Distinct("`node`.*")
	if inputs.Query != "" {
		tx = tx.Where("`node`.`ip` regexp ?", inputs.Query)
		tx = tx.Or("`node`.`physical_system` = ?", inputs.Query)
	}
	tx = tx.Limit(10)
	tx.Find(&nodes)

	resp := &APIGetNodeSelectOutputs{
		List: nodes,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetNodeGetInputs struct {
	Ip string `json:"ip" form:"ip"`
}

func NodeGet(c *gin.Context) {
	var inputs APIGetNodeGetInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal
	f := node.Node{}
	tx.Table(f.TableName()).Where("ip = ?", inputs.Ip).First(&f)
	h.JSONR(c, f)
	return
}

type APIPostNodeUpdateInputs struct {
	ID             int64   `json:"id" form:"id"`
	IP             string  `json:"ip" form:"ip" binding:"required"`
	Name           string  `json:"name" form:"name"`
	PhysicalSystem string  `json:"physicalSystem" form:"physicalSystem"`
	TagIDs         []int64 `json:"tagIDs" form:"tagIDs"`
	State          string  `json:"state" form:"state"`
	DevOwner       string  `json:"devOwner" form:"devOwner"`
}

// @Summary 创建新机器
// @Description
// @Produce json
// @Param IP formData string true "创建新机器IP"
// @Param Nodename formData string false "创建NodeName"
// @Param GroupID formData string false "创建NodeGroup"
// @Param Tenant formData string false "创建Tenant"
// @Param Env formData string false "创建Env"
// @Param Project formData string false "创建Project"
// @Param Module formData string false "创建Module"
// @Param DevOwner formData string false "创建DevOwner"
// @Success 200 {object} APIPostNodeUpdateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/create [post]
func NodeCreate(c *gin.Context) {
	var inputs APIPostNodeUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Begin()

	n := node.Node{
		IP:             inputs.IP,
		Name:           inputs.Name,
		PhysicalSystem: inputs.PhysicalSystem,
		DevOwner:       inputs.DevOwner,
		State:          inputs.State,
	}
	if err := tx.Model(node.Node{}).Create(&n).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	var tags []app.Tag
	if err := tx.Model(app.Tag{}).Where("id in (?)", inputs.TagIDs).Find(&tags).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}
	for _, tag := range tags {
		if err := tx.Model(node.NodeTagRel{}).Create(&node.NodeTagRel{Node: n.ID, Tag: tag.ID}).Error; err != nil {
			tx.Rollback()
			h.JSONR(c, h.ExpecStatus, err)
			return
		}
	}
	tx.Commit()

	h.JSONR(c, h.OKStatus, inputs)
	return
}

// @Summary 更新机器信息
// @Description
// @Produce json
// @Param IP formData string true "根据IP更新机器信息"
// @Param Nodename formData string false "更新NodeName"
// @Param GroupID formData string false "更新NodeGroup"
// @Param DevOwner formData string false "更新DevOwner"
// @Success 200 {object} APIPostNodeUpdateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/update [put]
func NodeUpdate(c *gin.Context) {
	var inputs APIPostNodeUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Begin()

	n := node.Node{}
	if err := tx.Model(node.Node{}).Where("ip = ?", inputs.IP).Find(&n).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	if err := tx.Model(node.Node{}).Where("ip = ?", inputs.IP).Updates(&node.Node{
		DevOwner: inputs.DevOwner,
		State:    inputs.State,
	}).Error; tx.Error != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	if err := tx.Model(node.NodeTagRel{}).Where(&node.NodeTagRel{Node: n.ID}).Delete(&node.NodeTagRel{}).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	for _, tagID := range inputs.TagIDs {
		if err := tx.Create(&node.NodeTagRel{Node: n.ID, Tag: tagID}).Error; err != nil {
			tx.Rollback()
			h.JSONR(c, h.ExpecStatus, err)
			return
		}
	}
	tx.Commit()

	// 重建tag图
	service.TagService.ReBuildGraphV2()

	h.JSONR(c, h.OKStatus, inputs)
	return
}

type APIPostNodeBatchUpdateInputs struct {
	IDs      []int64 `json:"ids" form:"ids"`
	TagIDs   []int64 `json:"tagIDs" form:"tagIDs"`
	DevOwner string  `json:"devOwner" form:"devOwner"`
}

// @Summary 更新机器标签和负责人
// @Description
// @Produce json
// @Param APIPostNodeBatchUpdateInputs body APIPostNodeBatchUpdateInputs false "更新机器标签和负责人信息"
// @Success 200 {object} APIPostNodeBatchUpdateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/batch/update [put]
func NodeBatchUpdate(c *gin.Context) {
	var inputs APIPostNodeBatchUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Begin()
	for _, id := range inputs.IDs {
		if err := tx.Model(node.Node{}).Where("id = ?", id).Updates(&node.Node{
			DevOwner: inputs.DevOwner,
		}).Error; err != nil {
			tx.Rollback()
			h.JSONR(c, h.ExpecStatus, err)
			return
		}

		if err := tx.Model(node.NodeTagRel{}).Where(&node.NodeTagRel{Node: id}).Delete(&node.NodeTagRel{}).Error; err != nil {
			tx.Rollback()
			h.JSONR(c, h.ExpecStatus, err)
			return
		}

		for _, tagID := range inputs.TagIDs {
			if err := tx.Create(&node.NodeTagRel{Node: id, Tag: tagID}).Error; err != nil {
				tx.Rollback()
				h.JSONR(c, h.ExpecStatus, err)
				return
			}
		}
	}
	tx.Commit()

	h.JSONR(c, h.OKStatus, inputs)
	return
}

// @Summary 根据机器ID获取机器详细信息
// @Description
// @Produce json
// @Param id query int true "根据机器ID获取机器详细信息"
// @Success 200 {object} node.Node
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/info [get]
func NodeInfo(c *gin.Context) {
	id := c.Query("id")
	n := node.Node{}
	g.Con().Portal.Model(n).Where("id = ?", id).First(&n)
	n.Tags = n.RelatedTags()
	h.JSONR(c, n)
	return
}

// @Summary 根据机器ID或IP获取机器详细信息
// @Description
// @Produce json
// @Param id query string true "根据机器ID获取机器详细信息"
// @Param ip query string true "根据机器IP获取机器详细信息"
// @Success 200 {object} node.Node
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/detail [get]
func NodeDetail(c *gin.Context) {
	id := c.Query("id")
	ip := c.Query("ip")
	n := node.Node{}
	if id != "" {
		g.Con().Portal.Model(n).Where("id = ?", id).First(&n)
	} else if ip != "" {
		g.Con().Portal.Model(n).Where("ip = ?", ip).First(&n)
	}
	n.Groups = n.RelatedGroups()
	h.JSONR(c, n)
	return
}

// @Summary 物理子系统类别
// @Description
// @Produce json
// @Success 200 {object} model.APIGetVariableOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/physical_system_choices [get]
func NodePhysicalSystemChoices(c *gin.Context) {
	var data []*model.APIGetVariableItem
	tx := g.Con().Portal.Model(node.Node{}).Debug()
	tx = tx.Select("distinct `physical_system` as `label`, `physical_system` as `value`")
	tx = tx.Order("`physical_system`")
	tx = tx.Find(&data)
	resp := model.APIGetVariableOutputs{
		List:       data,
		TotalCount: int64(len(data)),
	}
	h.JSONR(c, resp)
	return
}

// @Summary 区域类别
// @Description
// @Produce json
// @Success 200 {object} model.APIGetVariableOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/node/area_choices [get]
func NodeAreaChoices(c *gin.Context) {
	var data []*model.APIGetVariableItem
	tx := g.Con().Portal.Model(node.Node{}).Debug()
	tx = tx.Select("distinct `area_name` as `label`, `area_name` as `value`")
	tx = tx.Where("area_name != ''")
	tx = tx.Find(&data)
	resp := model.APIGetVariableOutputs{
		List:       data,
		TotalCount: int64(len(data)),
	}
	h.JSONR(c, resp)
	return
}
