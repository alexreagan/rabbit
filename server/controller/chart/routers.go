package chart

import (
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.Use(utils.AuthSessionMidd)

	chartGroup := r.Group("/api/v1/chart")
	chartGroup.GET("/bar", ChartBar)
	chartGroup.GET("/pie", ChartPie)
	chartGroup.GET("/vm/stat", ChartVMStat)
	chartGroup.GET("/container/stat", ChartContainerStat)
}
