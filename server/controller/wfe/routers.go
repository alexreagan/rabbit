package wfe

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	procGroup := r.Group("/api/v1/wfe")
	procGroup.POST("/create", Create)
	procGroup.POST("/execute", Execute)
	procGroup.POST("/histDetails", HistDetails)
	procGroup.POST("/todos", Todos)
	procGroup.POST("/todo2doing", Todo2Doing)
	procGroup.POST("/nextNodeInfo", NextNodeInfo)
}
