package pub

import (
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.Use(utils.AuthSessionMidd)

	pubGroup := r.Group("/api/v1/pub")
	pubGroup.GET("/list", List)
	pubGroup.GET("/info", Info)
	pubGroup.POST("/create", Create)
	pubGroup.PUT("/update", Update)
	pubGroup.PUT("/assign", Assign)
	pubGroup.GET("/proc/info", ProcInfo)

	procGroup := r.Group("/api/v1/proc")
	procGroup.POST("/create", ProcCreate)
	procGroup.POST("/execute", ProcExecute)
	procGroup.POST("/nextNodeInfo", NextNodeInfo)
	procGroup.GET("/getPersonByNode", GetPersonByNode)
	procGroup.GET("/getHistDetailList", GetHistDetailList)
}
