package sys

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	menuGroup := r.Group("/api/v1/menu")
	menuGroup.GET("/nav", MenuNav)
	menuGroup.GET("/list", MenuList)
	menuGroup.GET("/info/:id", MenuInfo)
	menuGroup.POST("/update", MenuUpdate)

	paramGroup := r.Group("/api/v1/param")
	paramGroup.GET("/list", ParamList)
	paramGroup.GET("/info", ParamInfo)
	paramGroup.POST("/create", ParamCreate)
	paramGroup.PUT("/update", ParamUpdate)
}
