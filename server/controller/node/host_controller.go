package node

import (
	"fmt"
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type OrderBy struct {
	Prop  string `json:"prop" form:"prop"`
	Order string `json:"order" form:"order"`
}

type APIGetHostListInputs struct {
	IP                    string  `json:"ip" form:"ip"`
	PhysicalSystem        string  `json:"physicalSystem" form:"physicalSystem"`
	CpuNumber             int64   `json:"cpuNumber" form:"cpuNumber"`
	AreaName              string  `json:"areaName" form:"areaName"`
	CpuUsageUpperLimit    float64 `json:"cpuUsageUpperLimit" form:"cpuUsageUpperLimit"`
	CpuUsageLowerLimit    float64 `json:"cpuUsageLowerLimit" form:"cpuUsageLowerLimit"`
	FsUsageUpperLimit     float64 `json:"fsUsageUpperLimit" form:"fsUsageUpperLimit"`
	FsUsageLowerLimit     float64 `json:"fsUsageLowerLimit" form:"fsUsageLowerLimit"`
	MemoryUsageUpperLimit float64 `json:"memoryUsageUpperLimit" form:"memoryUsageUpperLimit"`
	MemoryUsageLowerLimit float64 `json:"memoryUsageLowerLimit" form:"memoryUsageLowerLimit"`
	Group                 string  `json:"group" form:"group"`
	BoundGroup            string  `json:"boundGroup" form:"boundGroup"`
	Status                string  `json:"status" form:"status"`
	Limit                 int     `json:"limit" form:"limit"`
	Page                  int     `json:"page" form:"page"`
	OrderBy               string  `json:"orderBy" form:"orderBy"`
	Order                 string  `json:"order" form:"order"`
}

type APIGetHostListOutputs struct {
	List       []*node.Host `json:"list"`
	TotalCount int64        `json:"totalCount"`
}

func (input APIGetHostListInputs) checkInputsContain() error {
	return nil
}

// @Summary 机器列表接口
// @Description
// @Produce json
// @Param APIGetHostListInputs query APIGetHostListInputs true "根据查询条件分页查询机器列表"
// @Success 200 {object} APIGetHostListOutputs
// @Failure 400 {object} APIGetHostListOutputs
// @Router /api/v1/host/list [get]
func HostList(c *gin.Context) {
	var inputs APIGetHostListInputs
	inputs.Page = -1
	inputs.Limit = -1

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var hosts []*node.Host
	var totalCount int64
	db := g.Con().Portal.Debug().Model(node.Host{})
	db = db.Select("distinct `host`.*")
	if inputs.Group != "" || inputs.BoundGroup != "" {
		db = db.Joins("left join `host_group_rel` on `host`.id=`host_group_rel`.`host_id`")
	}
	if inputs.Group != "" {
		db = db.Joins("left join `host_group` on `host_group_rel`.`group_id`=`host_group`.`id`")
		db = db.Where("`host_group`.`name` regexp ?", inputs.Group)
	}
	if inputs.IP != "" {
		db = db.Where("`host`.`ip` regexp ?", inputs.IP)
	}
	if inputs.PhysicalSystem != "" {
		db = db.Where("`host`.`physical_system` = ?", inputs.PhysicalSystem)
	}
	if inputs.CpuUsageLowerLimit != 0 {
		if inputs.CpuUsageLowerLimit > 1 {
			inputs.CpuUsageLowerLimit = inputs.CpuUsageLowerLimit / 100
		}
		db = db.Where("`host`.`cpu_usage` >= ?", inputs.CpuUsageLowerLimit)
	}
	if inputs.CpuUsageUpperLimit != 0 {
		if inputs.CpuUsageUpperLimit > 1 {
			inputs.CpuUsageUpperLimit = inputs.CpuUsageUpperLimit / 100
		}
		db = db.Where("`host`.`cpu_usage` < ?", inputs.CpuUsageUpperLimit)
	}
	if inputs.FsUsageLowerLimit != 0 {
		if inputs.FsUsageLowerLimit > 1 {
			inputs.FsUsageLowerLimit = inputs.FsUsageLowerLimit / 100
		}
		db = db.Where("`host`.`fs_usage` >= ?", inputs.FsUsageLowerLimit)
	}
	if inputs.FsUsageUpperLimit != 0 {
		if inputs.FsUsageUpperLimit > 1 {
			inputs.FsUsageUpperLimit = inputs.FsUsageUpperLimit / 100
		}
		db = db.Where("`host`.`fs_usage` < ?", inputs.FsUsageUpperLimit)
	}
	if inputs.MemoryUsageLowerLimit != 0 {
		if inputs.MemoryUsageLowerLimit > 1 {
			inputs.MemoryUsageLowerLimit = inputs.MemoryUsageLowerLimit / 100
		}
		db = db.Where("`host`.`memory_usage` >= ?", inputs.FsUsageLowerLimit)
	}
	if inputs.MemoryUsageUpperLimit != 0 {
		if inputs.MemoryUsageUpperLimit > 1 {
			inputs.MemoryUsageUpperLimit = inputs.MemoryUsageUpperLimit / 100
		}
		db = db.Where("`host`.`memory_usage` < ?", inputs.MemoryUsageUpperLimit)
	}
	if inputs.CpuNumber != 0 {
		db = db.Where("`host`.`cpu_number` = ?", inputs.CpuNumber)
	}
	if inputs.AreaName != "" {
		db = db.Where("`host`.`area_name` = ?", inputs.AreaName)
	}
	if inputs.BoundGroup != "" {
		if inputs.BoundGroup == "bound" {
			db = db.Where("`host_group_rel`.`group_id` is not null")

		} else if inputs.BoundGroup == "unbound" {
			db = db.Where("`host_group_rel`.`group_id` is null")
		}
	}

	db.Count(&totalCount)
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&hosts)

	for _, host := range hosts {
		host.CpuUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", host.CpuUsage*100), 64)
		host.FsUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", host.FsUsage*100), 64)
		host.MemoryUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", host.MemoryUsage*100), 64)
		host.Groups = host.RelatedGroups()
	}

	resp := &APIGetHostListOutputs{
		List:       hosts,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetHostGetInputs struct {
	Ip string `json:"ip" form:"ip"`
}

func HostGet(c *gin.Context) {
	var inputs APIGetHostGetInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	db := g.Con().Portal
	f := node.Host{}
	db.Table(f.TableName()).Where("ip = ?", inputs.Ip).First(&f)
	h.JSONR(c, f)
	return
}

type APIPostHostUpdateInputs struct {
	ID             int64   `json:"id" form:"id"`
	IP             string  `json:"ip" form:"ip" binding:"required"`
	Name           string  `json:"name" form:"name"`
	PhysicalSystem string  `json:"physicalSystem" form:"physicalSystem"`
	GroupIds       []int64 `json:"groupIds" form:"groupIds"`
	Path           string  `json:"path" form:"path"`
	DevOwner       string  `json:"devOwner" form:"devOwner"`
}

// @Summary 创建新机器
// @Description
// @Produce json
// @Param IP formData string true "创建新机器IP"
// @Param Hostname formData string false "创建HostName"
// @Param GroupID formData string false "创建HostGroup"
// @Param Tenant formData string false "创建Tenant"
// @Param Env formData string false "创建Env"
// @Param Project formData string false "创建Project"
// @Param Module formData string false "创建Module"
// @Param DevOwner formData string false "创建DevOwner"
// @Success 200 {object} APIPostHostUpdateInputs
// @Failure 400 {object} APIPostHostUpdateInputs
// @Router /api/v1/host/create [post]
func HostCreate(c *gin.Context) {
	var inputs APIPostHostUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	log.Printf("%+v", inputs)

	tx := g.Con().Portal.Begin()

	host := node.Host{
		IP:             inputs.IP,
		Name:           inputs.Name,
		PhysicalSystem: inputs.PhysicalSystem,
		DevOwner:       inputs.DevOwner,
	}
	if dt := tx.Model(node.Host{}).Create(&host); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	var hostGroups []node.HostGroup
	if dt := tx.Model(node.HostGroup{}).Where("id in (?)", inputs.GroupIds).Find(&hostGroups); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	dt := tx.Debug().Model(node.HostGroupRel{})
	for _, grp := range hostGroups {
		if dt = dt.Create(&node.HostGroupRel{HostID: host.ID, GroupID: grp.ID}); dt.Error != nil {
			h.JSONR(c, h.ExpecStatus, dt.Error)
			dt.Rollback()
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
// @Param Hostname formData string false "更新HostName"
// @Param GroupID formData string false "更新HostGroup"
// @Param DevOwner formData string false "更新DevOwner"
// @Success 200 {object} APIPostHostUpdateInputs
// @Failure 400 {object} APIPostHostUpdateInputs
// @Router /api/v1/host/update [put]
func HostUpdate(c *gin.Context) {
	var inputs APIPostHostUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Begin()

	host := node.Host{}
	if dt := tx.Model(node.Host{}).Where("ip = ?", inputs.IP).Find(&host); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	var hostGroups []node.HostGroup
	if dt := tx.Model(node.HostGroup{}).Where("id in (?)", inputs.GroupIds).Find(&hostGroups); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	if dt := tx.Model(node.Host{}).Where("ip = ?", inputs.IP).Updates(node.Host{
		DevOwner: inputs.DevOwner,
	}); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	dt := tx.Debug().Model(node.HostGroupRel{})
	if dt = dt.Where(&node.HostGroupRel{HostID: host.ID}).Delete(&node.HostGroupRel{}); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		dt.Rollback()
		return
	}
	for _, grp := range hostGroups {
		if dt = dt.Create(&node.HostGroupRel{HostID: host.ID, GroupID: grp.ID}); dt.Error != nil {
			h.JSONR(c, h.ExpecStatus, dt.Error)
			dt.Rollback()
			return
		}
	}
	tx.Commit()

	h.JSONR(c, h.OKStatus, inputs)
	return
}

type APIPostHostBatchUpdateInputs struct {
	IDs      []int64 `json:"ids" form:"ids"`
	GroupIds []int64 `json:"groupIds" form:"groupIds"`
	//Path     string  `json:"path" form:"path"`
	DevOwner string `json:"devOwner" form:"devOwner"`
}

// @Summary 更新机器信息
// @Description
// @Produce json
// @Param GroupID formData string false "更新HostGroup"
// @Param DevOwner formData string false "更新DevOwner"
// @Success 200 {object} APIPostHostBatchUpdateInputs
// @Failure 400 {object} APIPostHostBatchUpdateInputs
// @Router /api/v1/host/batch/update [put]
func HostBatchUpdate(c *gin.Context) {
	var inputs APIPostHostBatchUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	log.Printf("%+v", inputs)

	tx := g.Con().Portal.Begin()
	for _, id := range inputs.IDs {
		log.Println(id)
		if dt := tx.Model(node.Host{}).Where("id = ?", id).Updates(node.Host{
			DevOwner: inputs.DevOwner,
		}); dt.Error != nil {
			h.JSONR(c, h.ExpecStatus, dt.Error)
		}

		dt := tx.Model(node.HostGroupRel{})
		if dt = dt.Where(&node.HostGroupRel{HostID: id}).Delete(&node.HostGroupRel{}); dt.Error != nil {
			h.JSONR(c, h.ExpecStatus, dt.Error)
			dt.Rollback()
			return
		}
		for _, grpId := range inputs.GroupIds {
			if dt = dt.Create(&node.HostGroupRel{HostID: id, GroupID: grpId}); dt.Error != nil {
				h.JSONR(c, h.ExpecStatus, dt.Error)
				dt.Rollback()
				return
			}
		}
		tx.Commit()
	}

	h.JSONR(c, h.OKStatus, inputs)
	return
}

// @Summary 根据机器ID获取机器详细信息
// @Description
// @Produce json
// @Param id path int true "根据机器ID获取机器详细信息"
// @Success 200 {object} node.Host
// @Failure 400 {object} node.Host
// @Router /api/v1/host/info/:id [get]
func HostInfo(c *gin.Context) {
	id := c.Param("id")
	host := node.Host{}
	g.Con().Portal.Model(host).Where("id = ?", id).First(&host)
	host.Groups = host.RelatedGroups()
	h.JSONR(c, host)
	return
}

// @Summary 根据机器ID或IP获取机器详细信息
// @Description
// @Produce json
// @Param id query string true "根据机器ID获取机器详细信息"
// @Param ip query string true "根据机器IP获取机器详细信息"
// @Success 200 {object} node.Host
// @Failure 400 {object} node.Host
// @Router /api/v1/host/detail [get]
func HostDetail(c *gin.Context) {
	id := c.Query("id")
	ip := c.Query("ip")
	host := node.Host{}
	if id != "" {
		g.Con().Portal.Model(host).Where("id = ?", id).First(&host)
	} else if ip != "" {
		g.Con().Portal.Model(host).Where("ip = ?", ip).First(&host)
	}
	host.Groups = host.RelatedGroups()
	h.JSONR(c, host)
	return
}

// @Summary 物理子系统类别
// @Description
// @Produce json
// @Success 200 {object} APIGetVariableOutputs
// @Failure 400 {object} APIGetVariableOutputs
// @Router /api/v1/host/physical_system_choices [get]
func HostPhysicalSystemChoices(c *gin.Context) {
	var data []*h.APIGetVariableItem
	db := g.Con().Portal.Model(node.Host{}).Debug()
	db = db.Select("distinct `physical_system` as `label`, `physical_system` as `value`")
	db = db.Order("`physical_system`")
	db = db.Find(&data)
	resp := h.APIGetVariableOutputs{
		List:       data,
		TotalCount: int64(len(data)),
	}
	h.JSONR(c, resp)
	return
}
