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
	hostGroup.GET("/select", HostSelect)

	hgGroup := r.Group("/api/v1/host_group")
	hgGroup.GET("/list", HostGroupList)
	hgGroup.GET("/all", HostGroupAll)
	hgGroup.GET("/get", HostGroupGet)
	hgGroup.POST("/create", HostGroupCreate)
	hgGroup.PUT("/update", HostGroupPut)
	hgGroup.POST("/delete/:id", HostGroupDelete)
	hgGroup.POST("/bind_host", BindHostToHostGroup)
	hgGroup.GET("/related_hosts", HostGroupRelatedHosts)

	hostApplyRequestGroup := r.Group("/api/v1/host_apply_request")
	hostApplyRequestGroup.GET("/list", HostApplyRequestList)
	hostApplyRequestGroup.GET("/info", HostApplyRequestInfo)
	hostApplyRequestGroup.POST("/create", HostApplyRequestCreate)
	hostApplyRequestGroup.PUT("/assign", HostApplyRequestAssign)
}
