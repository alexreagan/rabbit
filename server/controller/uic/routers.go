package uic

import (
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	// session
	u := r.Group("/api/v1/user")
	u.POST("/login", Login)
	u.POST("/logout", Logout)
	u.POST("/create", CreateUser)

	// auth api
	userGroup := r.Group("/api/v1/user")
	userGroup.Use(utils.AuthSessionMidd)
	userGroup.GET("/list", List)
	userGroup.GET("/info", Info)

	// role
	roleGroup := r.Group("/api/v1/role")
	roleGroup.Use(utils.AuthSessionMidd)
	roleGroup.GET("/select", RoleSelect)

	// department
	departGroup := r.Group("/api/v1/department")
	departGroup.Use(utils.AuthSessionMidd)
	departGroup.GET("/list", DepartmentLists)
}
