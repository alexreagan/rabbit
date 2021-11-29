package alert

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	alertGroup := r.Group("/api/v1/alert")
	alertGroup.GET("/list", AlertList)
	alertGroup.GET("/physical_system_choices", AlertPhysicalSystemChoices)
}
