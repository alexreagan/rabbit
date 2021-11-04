package node

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/caas"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/alexreagan/rabbit/server/utils"
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
	Name  string `json:"name" form:"name"`
	Limit int    `json:"limit" form:"limit"`
	Page  int    `json:"page" form:"page"`
}

type APIGetCaasServiceListOutputs struct {
	List       []*CaasService `json:"list"`
	TotalCount int64          `json:"totalCount"`
}

type CaasService struct {
	Namespace          string `json:"namespace"`
	WorkspaceName      string `json:"workspaceName"`
	ClusterName        string `json:"clusterName"`
	PhysicalSystemName string `json:"physicalSystemName"`
	caas.Service
}

// @Summary 更新host group信息
// @Description
// @Produce json
// @Param APIGetCaasServiceListInputs query APIGetCaasServiceListInputs true "更新host group信息"
// @Success 200 {object} APIGetCaasServiceListOutputs
// @Failure 400 {object} APIGetCaasServiceListOutputs
// @Router /api/v1/caas/service/list [get]
func CaasServiceList(c *gin.Context) {
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
	db = db.Joins("left join `caas_namespace_service_rel` on `service_id` = `caas_service`.`id`")
	db = db.Joins("left join `caas_namespace` on `caas_namespace`.`id` = `caas_namespace_service_rel`.`namespace_id`")
	if inputs.Name != "" {
		db.Where("`caas_service`.`service_name` regexp ?", inputs.Name)
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
	ServiceId int64 `json:"service_id" form:"service_id"`
}

// @Summary 更新service下的pods信息
// @Description
// @Produce json
// @Param APIGetCaasServiceListInputs query APIGetCaasServiceListInputs true "更新service下的pods信息"
// @Success 200 {object} APIGetCaasServiceListOutputs
// @Failure 400 {object} APIGetCaasServiceListOutputs
// @Router /api/v1/caas/service/refresh_pods [get]
func CaasServiceRefreshPods(c *gin.Context) {
	var inputs APIGetCaasServiceRefreshPodsInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	db := g.Con().Portal.Debug()

	var ser caas.Service
	db.Model(caas.Service{}).Where("id = ?", inputs.ServiceId).Find(&ser)

	var rel caas.NamespaceServiceRel
	db.Model(caas.NamespaceServiceRel{}).Where("service_id = ?", inputs.ServiceId).Find(&rel)
	if rel.ServiceID == 0 {
		h.JSONR(c, h.BadStatus, "no service id")
		return
	}

	var namespace caas.NameSpace
	db.Model(caas.NameSpace{}).Where("id = ?", rel.NamespaceID).Find(&namespace)

	pods, err := service.GetPod(&namespace, &ser)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		log.Errorln(err)
		return
	}
	service.UpdatePods(&ser, pods)

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
func CaasNamespaceList(c *gin.Context) {
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
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
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
// @Router /api/v1/pod/list [get]
func CaasPodList(c *gin.Context) {
	var inputs APIGetPodListInputs
	inputs.Page = -1
	inputs.Limit = -1

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var pods []*caas.Pod
	var totalCount int64
	db := g.Con().Portal.Debug().Model(caas.Pod{})
	db = db.Select("distinct `caas_pod`.*")

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
func CaasWorkspaceList(c *gin.Context) {
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
// @Router /api/v1/caas/pod/:id [get]
func CaasPodGet(c *gin.Context) {
	id := c.Param("id")

	db := g.Con().Portal
	f := caas.Pod{}
	db.Model(f).Where("id = ?", id).First(&f)
	f.AdditionalAttrs()
	h.JSONR(c, f)
	return
}
