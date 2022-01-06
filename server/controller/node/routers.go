package node

import (
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.Use(utils.AuthSessionMidd)

	ng := r.Group("/api/v1/node")
	ng.GET("/get", NodeGet)
	ng.GET("/list", NodeList)
	ng.GET("/detail", NodeDetail)
	ng.GET("/info", NodeInfo)
	ng.POST("/create", NodeCreate)
	ng.PUT("/update", NodeUpdate)
	ng.PUT("/batch/update", NodeBatchUpdate)
	ng.GET("/physical_system_choices", NodePhysicalSystemChoices)
	ng.GET("/area_choices", NodeAreaChoices)
	ng.GET("/select", NodeSelect)

	ngGroup := r.Group("/api/v1/node_group")
	ngGroup.GET("/list", NodeGroupList)
	ngGroup.GET("/all", NodeGroupAll)
	ngGroup.GET("/get", NodeGroupGet)
	ngGroup.POST("/create", NodeGroupCreate)
	ngGroup.PUT("/update", NodeGroupPut)
	ngGroup.POST("/delete/:id", NodeGroupDelete)
	ngGroup.POST("/bind_node", BindNodeToNodeGroup)
	ngGroup.GET("/related_nodes", NodeGroupRelatedNodes)

	nodeApplyRequestGroup := r.Group("/api/v1/node_apply_request")
	nodeApplyRequestGroup.GET("/list", NodeApplyRequestList)
	nodeApplyRequestGroup.GET("/info", NodeApplyRequestInfo)
	nodeApplyRequestGroup.POST("/create", NodeApplyRequestCreate)
	nodeApplyRequestGroup.PUT("/assign", NodeApplyRequestAssign)
}
