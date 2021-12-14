package node

import (
	"github.com/alexreagan/rabbit/server/controller/app"
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

	tGroup := r.Group("/api/v1/tree")
	tGroup.GET("", app.Tree)
	tGroup.GET("/rebuild", app.TreeRebuild)
	//tGroup.POST("/dragging", TreeDragging)

	tagGroup := r.Group("/api/v1/tag")
	tagGroup.GET("/list", app.TagList)
	tagGroup.GET("/info", app.TagInfo)
	tagGroup.POST("/create", app.TagCreate)
	tagGroup.PUT("/update", app.TagUpdate)

	tcGroup := r.Group("/api/v1/tag_category")
	tcGroup.GET("/list", app.TagCategoryList)
	tcGroup.GET("/info", app.TagCategoryInfo)
	tcGroup.POST("/create", app.TagCategoryCreate)
	tcGroup.PUT("/update", app.TagCategoryUpdate)
	tcGroup.PUT("/tags", app.TagCategoryTags)

	v2TreeGroup := r.Group("/api/v2/tree")
	v2TreeGroup.GET("", app.V2Tree)

	hostApplyRequestGroup := r.Group("/api/v1/host_apply_request")
	hostApplyRequestGroup.GET("/list", HostApplyRequestList)
	hostApplyRequestGroup.GET("/info", HostApplyRequestInfo)
	hostApplyRequestGroup.POST("/create", HostApplyRequestCreate)
	hostApplyRequestGroup.PUT("/assign", HostApplyRequestAssign)
}
