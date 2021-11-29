package node

import (
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.Use(utils.AuthSessionMidd)

	hostGroup := r.Group("/api/v1/host")
	hostGroup.GET("/get", HostGet)
	hostGroup.GET("/list", HostList)
	hostGroup.GET("/detail", HostDetail)
	hostGroup.GET("/info", HostInfo)
	hostGroup.POST("/create", HostCreate)
	hostGroup.PUT("/update", HostUpdate)
	hostGroup.PUT("/batch/update", HostBatchUpdate)
	hostGroup.GET("/physical_system_choices", HostPhysicalSystemChoices)
	hostGroup.GET("/area_choices", HostAreaChoices)

	hgGroup := r.Group("/api/v1/host_group")
	hgGroup.GET("/list", HostGroupList)
	hgGroup.GET("/all", HostGroupAll)
	hgGroup.GET("/get", HostGroupGet)
	hgGroup.POST("/create", HostGroupCreate)
	hgGroup.PUT("/update", HostGroupPut)
	hgGroup.POST("/delete/:id", HostGroupDelete)
	hgGroup.POST("/bind_host", BindHostToHostGroup)
	hgGroup.GET("/related_hosts", HostGroupRelatedHosts)

	tGroup := r.Group("/api/v1/tree")
	tGroup.GET("", Tree)
	tGroup.GET("/rebuild", TreeRebuild)
	tGroup.POST("/dragging", TreeDragging)

	tagGroup := r.Group("/api/v1/tag")
	tagGroup.GET("/list", TagList)
	tagGroup.GET("/info", TagInfo)
	tagGroup.POST("/create", TagCreate)
	tagGroup.PUT("/update", TagUpdate)

	tcGroup := r.Group("/api/v1/tag_category")
	tcGroup.GET("/list", TagCategoryList)
	tcGroup.GET("/info", TagCategoryInfo)
	tcGroup.POST("/create", TagCategoryCreate)
	tcGroup.PUT("/update", TagCategoryUpdate)
	tcGroup.PUT("/tags", TagCategoryTags)

	chartGroup := r.Group("/api/v1/chart")
	chartGroup.GET("/bar", ChartBar)
	chartGroup.GET("/pie", ChartPie)
	chartGroup.GET("/vm/stat", ChartVMStat)
	chartGroup.GET("/container/stat", ChartContainerStat)

	caasGroup := r.Group("/api/v1/caas")
	caasGroup.GET("/workspace/list", CaasWorkspaceList)
	caasGroup.GET("/namespace/list", CaasNamespaceList)
	caasGroup.GET("/service/list", CaasServiceList)
	caasGroup.GET("/service/info", CaasServiceInfo)
	caasGroup.GET("/service/refresh_pods", CaasServiceRefreshPods)

	podGroup := r.Group("/api/v1/caas/pod")
	podGroup.GET("/list", CaasPodList)
	podGroup.GET("/:id", CaasPodGet)

	v2TreeGroup := r.Group("/api/v2/tree")
	v2TreeGroup.GET("", V2Tree)

	hostApplyRequestGroup := r.Group("/api/v1/host_apply_request")
	hostApplyRequestGroup.GET("/list", HostApplyRequestList)
	hostApplyRequestGroup.GET("/info", HostApplyRequestInfo)
	hostApplyRequestGroup.POST("/create", HostApplyRequestCreate)
}
