package procmanager

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	// apiV1 := r.Group("/api/v1")
	// apiV1.Use(utils.AuthSessionMidd)

	mainGroup := r.Group("/")
	mainGroup.POST("login", ProcManagerLogin)

	procManagerApiGroup := r.Group("/procmanager/api")
	procManagerApiGroup.POST("/procCreate", ProcCreate)
	procManagerApiGroup.POST("/procExecute", ProcExecute)
	procManagerApiGroup.POST("/procInstTodoInfo", ProcInstTodoInfo)
	procManagerApiGroup.POST("/procNextNodeInfo", ProcNextNodeInfo)
}
