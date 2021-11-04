package node

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	hostGroup := r.Group("/api/v1/host")
	hostGroup.GET("/get", HostGet)
	hostGroup.GET("/list", HostList)
	hostGroup.GET("/info/:id", HostInfo)
	hostGroup.POST("/create", HostCreate)
	hostGroup.PUT("/update", HostUpdate)
	hostGroup.PUT("/batch/update", HostBatchUpdate)

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

	chartGroup := r.Group("/api/v1/chart")
	chartGroup.GET("/bar", ChartBar)
	chartGroup.GET("/pie", ChartPie)
	chartGroup.GET("/vm/stat", ChartVMStat)
	chartGroup.GET("/container/stat", ChartContainerStat)

	caasGroup := r.Group("/api/v1/caas")
	caasGroup.GET("/workspace/list", CaasWorkspaceList)
	caasGroup.GET("/namespace/list", CaasNamespaceList)
	caasGroup.GET("/service/list", CaasServiceList)
	caasGroup.GET("/service/refresh_pods", CaasServiceRefreshPods)

	podGroup := r.Group("/api/v1/caas/pod")
	podGroup.GET("/list", CaasPodList)
	podGroup.GET("/:id", CaasPodGet)
}
