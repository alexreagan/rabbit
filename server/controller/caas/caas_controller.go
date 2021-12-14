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
)

type APIGetEnvListInputs struct {
	Name  string `json:"name" form:"name"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

type APIGetCaasServiceListInputs struct {
	ServiceName string `json:"serviceName" form:"serviceName"`
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

// @Summary 更新host group信息
// @Description
// @Produce json
// @Param APIGetCaasServiceListInputs query APIGetCaasServiceListInputs true "更新host group信息"
// @Success 200 {object} APIGetCaasServiceListOutputs
// @Failure 400 {object} APIGetCaasServiceListOutputs
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
	db := g.Con().Portal.Model(caas.Service{}).Debug()
	db = db.Select("`caas_service`.*, `caas_namespace`.`namespace`, `caas_namespace`.`workspace_name`, `caas_namespace`.`cluster_name`, `caas_namespace`.`physical_system_name`")
	db = db.Joins("left join `caas_namespace_service_rel` on `service` = `caas_service`.`id`")
	db = db.Joins("left join `caas_namespace` on `caas_namespace`.`id` = `caas_namespace_service_rel`.`namespace`")
	if inputs.ServiceName != "" {
		db.Where("`caas_service`.`service_name` regexp ?", inputs.ServiceName)
	}
	db.Count(&totalCount)
	db.Offset(offset).Limit(limit).Find(&services)

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
// @Failure 400 {object} APIGetCaasServiceListOutputs
// @Router /api/v1/caas/service/refresh_pods [get]
func ServiceRefreshPods(c *gin.Context) {
	var inputs APIGetCaasServiceRefreshPodsInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	db := g.Con().Portal.Debug()

	var ser caas.Service
	db.Model(caas.Service{}).Where("id = ?", inputs.ServiceID).Find(&ser)

	var rel caas.NamespaceServiceRel
	db.Model(caas.NamespaceServiceRel{}).Where("service = ?", inputs.ServiceID).Find(&rel)
	if rel.Service == 0 {
		h.JSONR(c, h.BadStatus, "no service id")
		return
	}

	var namespace caas.NameSpace
	db.Model(caas.NameSpace{}).Where("id = ?", rel.NameSpace).Find(&namespace)

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
// @Failure 400 {object} APIGetCaasNamespaceListOutputs
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
	db := g.Con().Portal.Model(caas.NameSpace{}).Debug()
	if inputs.Namespace != "" {
		db = db.Where("`namespace` regexp ?", inputs.Namespace)
	}
	if inputs.WorkspaceName != "" {
		db = db.Where("`workspace_name` regexp ?", inputs.WorkspaceName)
	}
	if inputs.ClusterName != "" {
		db = db.Where("cluster_name regexp ?", inputs.ClusterName)
	}
	if inputs.PhysicalSystemName != "" {
		db = db.Where("physical_system_name regexp ?", inputs.PhysicalSystemName)
	}
	db.Count(&totalCount)
	db.Offset(offset).Limit(limit).Find(&namespaces)

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
// @Success 200 {object} APIGetHostListOutputs
// @Failure 400 {object} APIGetHostListOutputs
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
	db := g.Con().Portal.Debug().Model(caas.Pod{})
	db = db.Select("distinct `caas_pod`.*")
	if inputs.Name != "" {
		db = db.Where("name regexp ?", inputs.Name)
	}
	if inputs.Namespace != "" {
		db = db.Where("namespace = ?", inputs.Namespace)
	}
	if inputs.ServiceName != "" {
		db = db.Where("service_name = ?", inputs.ServiceName)
	}
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db.Count(&totalCount)
	db = db.Offset(offset).Limit(limit)
	db.Find(&pods)

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
// @Failure 400 {object} APIGetCaasWorkspaceListOutputs
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
	db := g.Con().Portal.Model(caas.WorkSpace{}).Debug()
	if inputs.Name != "" {
		db = db.Where("`name` regexp ?", inputs.Name)
	}
	db.Count(&totalCount)
	db.Offset(offset).Limit(limit).Find(&workspaces)

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
// @Failure 400 {object} caas.Pod
// @Router /api/v1/caas/pod/info [get]
func PodInfo(c *gin.Context) {
	id := c.Query("id")

	db := g.Con().Portal
	f := caas.Pod{}
	db.Model(f).Where("id = ?", id).First(&f)
	f.AdditionalAttrs()
	h.JSONR(c, f)
	return
}

// @Summary 获取service详细信息
// @Description
// @Produce json
// @Param id query int64 true "获取service详细信息"
// @Success 200 {object} caas.Service
// @Failure 400 {object} caas.Service
// @Router /api/v1/caas/service/info [get]
func ServiceInfo(c *gin.Context) {
	id := c.Query("id")

	var srv caas.Service
	db := g.Con().Portal.Model(caas.Service{})
	db = db.Where("`caas_service`.`id` = ?", id)
	db = db.Find(&srv)

	var srvInfo CaasService
	db = g.Con().Portal.Model(caas.Service{})
	db = db.Select("`caas_namespace`.`namespace`, `caas_namespace`.`workspace_name`, `caas_namespace`.`cluster_name`, `caas_namespace`.`physical_system_name`")
	db = db.Joins("left join `caas_namespace_service_rel` on `service` = `caas_service`.`id`")
	db = db.Joins("left join `caas_namespace` on `caas_namespace`.`id` = `caas_namespace_service_rel`.`namespace`")
	db = db.Where("`caas_service`.`id` = ?", id)
	db = db.Find(&srvInfo)

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
// @Failure 400 {object} APIPutServiceUpdateInputs
// @Router /api/v1/caas/service/update [put]
func ServiceUpdate(c *gin.Context) {
	var inputs APIPutServiceUpdateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tx := g.Con().Portal.Begin()

	ser := caas.Service{}
	if dt := tx.Model(caas.Service{}).Where("id = ?", inputs.ID).Find(&ser); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	if dt := tx.Model(caas.Service{}).Where("id = ?", inputs.ID).Updates(caas.Service{
		Owner: inputs.Owner,
	}); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		return
	}

	dt := tx.Debug().Model(caas.ServiceTagRel{})
	if dt = dt.Where(&caas.ServiceTagRel{Service: inputs.ID}).Delete(&caas.ServiceTagRel{}); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
		dt.Rollback()
		return
	}

	for _, tagID := range inputs.TagIDs {
		if dt = dt.Create(&caas.ServiceTagRel{Service: inputs.ID, Tag: tagID}); dt.Error != nil {
			h.JSONR(c, h.ExpecStatus, dt.Error)
			dt.Rollback()
			return
		}
	}
	tx.Commit()

	// 重建tag图
	service.TagService.ReBuildGraph()

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
// @Param APIGetCaasAppListInputs body APIGetCaasAppListInputs true "获取caas应用信息"
// @Success 200 {object} APIGetCaasAppListOutputs
// @Failure 400 {object} APIGetCaasAppListOutputs
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
	db := g.Con().Portal.Model(caas.App{}).Debug()
	db = db.Select("`caas_app`.`id`, `caas_app`.`app_name`, `caas_app`.`description`, `caas_app`.`create_time`, `caas_app`.`update_time`, `caas_namespace`.`namespace` as namespace_name")
	db = db.Joins("left join `caas_namespace` on `caas_app`.`namespace_id` = `caas_namespace`.`id`")
	if inputs.AppName != "" {
		db = db.Where("`caas_app`.`app_name` regexp ?", inputs.AppName)
	}
	db.Count(&totalCount)
	db.Offset(offset).Limit(limit).Find(&apps)

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
// @Failure 400 {object} caas.App
// @Router /api/v1/caas/app/info [get]
func AppInfo(c *gin.Context) {
	id := c.Query("id")

	var app *caas.App
	db := g.Con().Portal.Model(caas.App{})
	db = db.Where("id = ?", id)
	db = db.Find(&app)

	h.JSONR(c, http.StatusOK, app)
	return
}
