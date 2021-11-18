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
	u.POST("/create", UserCreate)

	// auth api
	userGroup := r.Group("/api/v1/user")
	userGroup.Use(utils.AuthSessionMidd)
	userGroup.GET("/list", UserList)
	userGroup.GET("/info", UserInfo)
	userGroup.GET("/myself", UserMyself)
	userGroup.PUT("/update", UserUpdate)

	// role
	roleGroup := r.Group("/api/v1/role")
	roleGroup.Use(utils.AuthSessionMidd)
	roleGroup.GET("/list", RoleList)
	roleGroup.GET("/info", RoleInfo)
	roleGroup.GET("/select", RoleSelect)
	roleGroup.POST("/create", RoleCreate)
	roleGroup.PUT("/update", RoleUpdate)

	// perm
	permGroup := r.Group("/api/v1/perm")
	permGroup.Use(utils.AuthSessionMidd)
	permGroup.GET("/list", PermList)
	permGroup.GET("/info", PermInfo)
	permGroup.POST("/create", PermCreate)
	permGroup.PUT("/update", PermUpdate)
	permGroup.GET("/myself", PermMyself)

	// department
	departGroup := r.Group("/api/v1/department")
	departGroup.Use(utils.AuthSessionMidd)
	departGroup.GET("/list", DepartmentLists)
}
