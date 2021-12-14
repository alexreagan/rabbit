package caas

import (
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.Use(utils.AuthSessionMidd)

	caasGroup := r.Group("/api/v1/caas")
	caasGroup.GET("/workspace/list", WorkspaceList)
	caasGroup.GET("/namespace/list", NamespaceList)
	caasGroup.GET("/service/list", ServiceList)
	caasGroup.GET("/service/info", ServiceInfo)
	caasGroup.PUT("/service/update", ServiceUpdate)
	caasGroup.GET("/service/refresh_pods", ServiceRefreshPods)
	caasGroup.GET("/app/list", AppList)
	caasGroup.GET("/app/info", AppInfo)

	podGroup := r.Group("/api/v1/caas/pod")
	podGroup.GET("/list", PodList)
	podGroup.GET("/info", PodInfo)
}
