package alarm

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	alarmGroup := r.Group("/api/v1/alarm")
	alarmGroup.GET("/list", List)
	alarmGroup.GET("/physical_system_choices", PhysicalSystemChoices)
}
