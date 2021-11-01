package portal

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	menuGroup := r.Group("/api/v1/menu")
	menuGroup.GET("/nav", MenuNav)
	menuGroup.GET("/list", MenuList)
	menuGroup.GET("/info/:id", MenuInfo)
	menuGroup.POST("/update", MenuUpdate)

	envGroup := r.Group("/api/v1/env")
	envGroup.GET("/list", EnvList)
}
