package caas

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/caas"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/alexreagan/rabbit/server/worker"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type APIGetEnvListInputs struct {
	Name  string `json:"name" form:"name"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

type APIGetCaasServiceListInputs struct {
	ServiceName string `json:"serviceName" form:"serviceName"`
	TagIDs     []int64 `json:"tagIDs[]" form:"tagIDs[]"`
	Limit       int    `json:"limit" form:"limit"`
	Page        int    `json:"page" form:"page"`
}

type APIGetCaasServiceListOutputs struct {
	List       []*CaasService `json:"list"`
	TotalCount int64          `json:"totalCount"`
}

type CaasService struct {
	Namespace          string     `json:"namespace"`
	WorkspaceName      string     `json:"workspaceName"`
	ClusterName        string     `json:"clusterName"`
	PhysicalSystemName string     `json:"physicalSystemName"`
	Tags               []*app.Tag `json:"tags"`
	caas.Service
}

// @Summary 获取caas service列表
// @Description
// @Produce json
// @Param APIGetCaasServiceListInputs query APIGetCaasServiceListInputs true "获取caas service列表"
// @Success 200 {object} APIGetCaasServiceListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/service/list [get]
func ServiceList(c *gin.Context) {
	var inputs APIGetCaasServiceListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}
	var services []*CaasService
	var totalCount int64
	tx := g.Con().Portal.Model(caas.Service{}).Debug()
	tx = tx.Select("`caas_service`.*, `caas_namespace`.`namespace`, `caas_namespace`.`workspace_name`, `caas_namespace`.`cluster_name`, `caas_namespace`.`physical_system_name`")
	tx = tx.Joins("left join `caas_namespace_service_rel` on `service` = `caas_service`.`id`")
	tx = tx.Joins("left join `caas_namespace` on `caas_namespace`.`id` = `caas_namespace_service_rel`.`namespace`")
	tx = tx.Joins("left join `caas_service_tag_rel` on `caas_service`.id=`caas_service_tag_rel`.`service`")
	if inputs.ServiceName != "" {
		tx.Where("`caas_service`.`service_name` regexp ?", inputs.ServiceName)
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
		tx = tx.Where("`caas_service_tag_rel`.`tag` in (?)", inputs.TagIDs)
		tx = tx.Group("`caas_service_tag_rel`.`service`")
		tx = tx.Having("group_concat(`caas_service_tag_rel`.`tag` order by `caas_service_tag_rel`.`tag`) = ?", strings.Join(tmp, ","))
	} else {
		tx = tx.Group("`caas_service`.`id`")
	}
	tx.Count(&totalCount)
	tx.Offset(offset).Limit(limit).Find(&services)

	for _, s := range services {
		s.Tags = s.RelatedTags()
	}

	resp := &APIGetCaasServiceListOutputs{
		List:       services,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetCaasServiceRefreshPodsInputs struct {
	ServiceID int64 `json:"service_id" form:"service_id"`
}

// @Summary 更新service下的pods信息
// @Description
// @Produce json
// @Param APIGetCaasServiceListInputs query APIGetCaasServiceListInputs true "更新service下的pods信息"
// @Success 200 {object} APIGetCaasServiceListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/service/refresh_pods [get]
func ServiceRefreshPods(c *gin.Context) {
	var inputs APIGetCaasServiceRefreshPodsInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Debug()

	var ser caas.Service
	tx.Model(caas.Service{}).Where("id = ?", inputs.ServiceID).Find(&ser)

	var rel caas.NamespaceServiceRel
	tx.Model(caas.NamespaceServiceRel{}).Where("service = ?", inputs.ServiceID).Find(&rel)
	if rel.Service == 0 {
		h.JSONR(c, h.BadStatus, "no service id")
		return
	}

	var namespace caas.NameSpace
	tx.Model(caas.NameSpace{}).Where("id = ?", rel.NameSpace).Find(&namespace)

	pods, err := worker.GetPod(&namespace, &ser)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		log.Errorln(err)
		return
	}
	worker.UpdatePods(&ser, pods)

	h.JSONR(c, h.BadStatus, pods)
	return
}

type APIGetCaasNamespaceListInputs struct {
	Namespace          string `json:"namespace" form:"namespace"`
	WorkspaceName      string `json:"workspaceName" form:"workspaceName"`
	ClusterName        string `json:"clusterName" form:"clusterName"`
	PhysicalSystemName string `json:"physicalSystemName" form:"physicalSystemName"`
	Limit              int    `json:"limit" form:"limit"`
	Page               int    `json:"page" form:"page"`
}

type APIGetCaasNamespaceListOutputs struct {
	List       []*caas.NameSpace `json:"list"`
	TotalCount int64             `json:"totalCount"`
}

// @Summary 获取caas项目空间信息
// @Description
// @Produce json
// @Param APIGetCaasNamespaceListInputs query APIGetCaasNamespaceListInputs true "获取caas项目空间信息"
// @Success 200 {object} APIGetCaasNamespaceListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/namespace/list [get]
func NamespaceList(c *gin.Context) {
	var inputs APIGetCaasNamespaceListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}
	var namespaces []*caas.NameSpace
	var totalCount int64
	tx := g.Con().Portal.Model(caas.NameSpace{}).Debug()
	if inputs.Namespace != "" {
		tx = tx.Where("`namespace` regexp ?", inputs.Namespace)
	}
	if inputs.WorkspaceName != "" {
		tx = tx.Where("`workspace_name` regexp ?", inputs.WorkspaceName)
	}
	if inputs.ClusterName != "" {
		tx = tx.Where("cluster_name regexp ?", inputs.ClusterName)
	}
	if inputs.PhysicalSystemName != "" {
		tx = tx.Where("physical_system_name regexp ?", inputs.PhysicalSystemName)
	}
	tx.Count(&totalCount)
	tx.Offset(offset).Limit(limit).Find(&namespaces)

	resp := &APIGetCaasNamespaceListOutputs{
		List:       namespaces,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetPodListInputs struct {
	Name        string `json:"name" form:"name"`
	Namespace   string `json:"namespace" form:"namespace"`
	ServiceName string `json:"serviceName" form:"serviceName"`
	Limit       int    `json:"limit" form:"limit"`
	Page        int    `json:"page" form:"page"`
	OrderBy     string `json:"orderBy" form:"orderBy"`
	Order       string `json:"order" form:"order"`
}

type APIGetPodListOutputs struct {
	List       []*caas.Pod `json:"list"`
	TotalCount int64       `json:"totalCount"`
}

// @Summary pod列表接口
// @Description
// @Produce json
// @Param APIGetPodListInputs query APIGetPodListInputs true "根据查询条件分页查询机器列表"
// @Success 200 {object} APIGetPodListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/pod/list [get]
func PodList(c *gin.Context) {
	var inputs APIGetPodListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var pods []*caas.Pod
	var totalCount int64
	tx := g.Con().Portal.Debug().Model(caas.Pod{})
	tx = tx.Select("distinct `caas_pod`.*")
	if inputs.Name != "" {
		tx = tx.Where("name regexp ?", inputs.Name)
	}
	if inputs.Namespace != "" {
		tx = tx.Where("namespace = ?", inputs.Namespace)
	}
	if inputs.ServiceName != "" {
		tx = tx.Where("service_name = ?", inputs.ServiceName)
	}
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx.Count(&totalCount)
	tx = tx.Offset(offset).Limit(limit)
	tx.Find(&pods)

	resp := &APIGetPodListOutputs{
		List:       pods,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetCaasWorkspaceListInputs struct {
	Name  string `json:"name" form:"name"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

type APIGetCaasWorkspaceListOutputs struct {
	List       []*caas.WorkSpace `json:"list"`
	TotalCount int64             `json:"totalCount"`
}

// @Summary 获取caas组织空间信息
// @Description
// @Produce json
// @Param APIGetCaasWorkspaceListInputs query APIGetCaasWorkspaceListInputs true "获取caas组织空间信息"
// @Success 200 {object} APIGetCaasWorkspaceListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/workspace/list [get]
func WorkspaceList(c *gin.Context) {
	var inputs APIGetCaasWorkspaceListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}
	var workspaces []*caas.WorkSpace
	var totalCount int64
	tx := g.Con().Portal.Model(caas.WorkSpace{}).Debug()
	if inputs.Name != "" {
		tx = tx.Where("`name` regexp ?", inputs.Name)
	}
	tx.Count(&totalCount)
	tx.Offset(offset).Limit(limit).Find(&workspaces)

	resp := &APIGetCaasWorkspaceListOutputs{
		List:       workspaces,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIGetPodGetInputs struct {
	ID int64 `json:"id"`
}

// @Summary 根据机器ID获取机器详细信息
// @Description
// @Produce json
// @Param id path int true "根据机器ID获取机器详细信息"
// @Success 200 {object} caas.Pod
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/pod/info [get]
func PodInfo(c *gin.Context) {
	id := c.Query("id")

	tx := g.Con().Portal
	f := caas.Pod{}
	tx.Model(f).Where("id = ?", id).First(&f)
	f.AdditionalAttrs()
	h.JSONR(c, f)
	return
}

// @Summary 获取service详细信息
// @Description
// @Produce json
// @Param id query int64 true "获取service详细信息"
// @Success 200 {object} caas.Service
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/service/info [get]
func ServiceInfo(c *gin.Context) {
	id := c.Query("id")

	var srv caas.Service
	tx := g.Con().Portal.Model(caas.Service{})
	tx = tx.Where("`caas_service`.`id` = ?", id)
	tx = tx.Find(&srv)

	var srvInfo CaasService
	tx = g.Con().Portal.Model(caas.Service{})
	tx = tx.Select("`caas_namespace`.`namespace`, `caas_namespace`.`workspace_name`, `caas_namespace`.`cluster_name`, `caas_namespace`.`physical_system_name`")
	tx = tx.Joins("left join `caas_namespace_service_rel` on `service` = `caas_service`.`id`")
	tx = tx.Joins("left join `caas_namespace` on `caas_namespace`.`id` = `caas_namespace_service_rel`.`namespace`")
	tx = tx.Where("`caas_service`.`id` = ?", id)
	tx = tx.Find(&srvInfo)

	resp := CaasService{
		Namespace:          srvInfo.Namespace,
		WorkspaceName:      srvInfo.WorkspaceName,
		ClusterName:        srvInfo.ClusterName,
		PhysicalSystemName: srvInfo.PhysicalSystemName,
		Service:            srv,
		Tags:               service.CaasService.GetServiceRelatedTags(&srv),
	}

	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIPutServiceUpdateInputs struct {
	ID     int64   `json:"id" form:"id"`
	TagIDs []int64 `json:"tagIDs" form:"tagIDs"`
	Owner  string  `json:"owner" form:"owner"`
}

// @Summary 更新service信息
// @Description
// @Produce json
// @Success 200 {object} APIPutServiceUpdateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/service/update [put]
func ServiceUpdate(c *gin.Context) {
	var inputs APIPutServiceUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Begin()

	ser := caas.Service{}
	if err := tx.Model(caas.Service{}).Where("id = ?", inputs.ID).Find(&ser).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	if err := tx.Model(caas.Service{}).Where("id = ?", inputs.ID).Updates(&caas.Service{
		Owner: inputs.Owner,
	}).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	if err := tx.Model(caas.ServiceTagRel{}).Where(&caas.ServiceTagRel{Service: inputs.ID}).Delete(&caas.ServiceTagRel{}).Error; err != nil {
		tx.Rollback()
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	for _, tagID := range inputs.TagIDs {
		if err := tx.Create(&caas.ServiceTagRel{Service: inputs.ID, Tag: tagID}).Error; err != nil {
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

type APIGetCaasAppListInputs struct {
	AppName string `json:"appName" form:"appName"`
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
}

type APIGetCaasAppListOutputs struct {
	List       []*App `json:"list"`
	TotalCount int64  `json:"totalCount"`
}

type App struct {
	caas.App
	NamespaceName string `json:"namespaceName"`
}

// @Summary 获取caas应用信息
// @Description
// @Produce json
// @Param APIGetCaasAppListInputs query APIGetCaasAppListInputs true "获取caas应用信息"
// @Success 200 {object} APIGetCaasAppListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/app/list [get]
func AppList(c *gin.Context) {
	var inputs APIGetCaasAppListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}
	var apps []*App
	var totalCount int64
	tx := g.Con().Portal.Model(caas.App{}).Debug()
	tx = tx.Select("`caas_app`.`id`, `caas_app`.`app_name`, `caas_app`.`description`, `caas_app`.`create_time`, `caas_app`.`update_time`, `caas_namespace`.`namespace` as namespace_name")
	tx = tx.Joins("left join `caas_namespace` on `caas_app`.`namespace_id` = `caas_namespace`.`id`")
	if inputs.AppName != "" {
		tx = tx.Where("`caas_app`.`app_name` regexp ?", inputs.AppName)
	}
	tx.Count(&totalCount)
	tx.Offset(offset).Limit(limit).Find(&apps)

	resp := &APIGetCaasAppListOutputs{
		List:       apps,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 获取app详细信息
// @Description
// @Produce json
// @Param id query int64 true "获取app详细信息"
// @Success 200 {object} caas.App
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/caas/app/info [get]
func AppInfo(c *gin.Context) {
	id := c.Query("id")

	var caasApp *caas.App
	tx := g.Con().Portal.Model(caas.App{})
	tx = tx.Where("id = ?", id)
	tx = tx.Find(&caasApp)

	h.JSONR(c, http.StatusOK, caasApp)
	return
}
