package app

import (
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.Use(utils.AuthSessionMidd)

	tGroup := r.Group("/api/v1/tree")
	tGroup.GET("/children", TreeChildren)
	tGroup.GET("/rebuild", TreeRebuild)

	v2TreeGroup := r.Group("/api/v2")
	v2TreeGroup.GET("/tree/children", V2TreeChildren)

	v3TreeGroup := r.Group("/api/v3")
	v3TreeGroup.GET("/tree/children", V3TreeChildren)
	v3TreeGroup.GET("/tree/node", V3TreeNode)

	tagGroup := r.Group("/api/v1/tag")
	tagGroup.GET("/list", TagList)
	tagGroup.GET("/all", TagAll)
	tagGroup.GET("/info", TagInfo)
	tagGroup.POST("/create", TagCreate)
	tagGroup.PUT("/update", TagUpdate)

	tcGroup := r.Group("/api/v1/tag_category")
	tcGroup.GET("/all", TagCategoryAll)
	tcGroup.GET("/list", TagCategoryList)
	tcGroup.GET("/info", TagCategoryInfo)
	tcGroup.POST("/create", TagCategoryCreate)
	tcGroup.PUT("/update", TagCategoryUpdate)
	tcGroup.PUT("/tags", TagCategoryTags)

	templateGroup := r.Group("/api/v1/template")
	templateGroup.GET("/list", TemplateList)
	templateGroup.GET("/all", TemplateAll)
	templateGroup.GET("/info", TemplateInfo)
	templateGroup.POST("/create", TemplateCreate)
	templateGroup.PUT("/update", TemplateUpdate)
	templateGroup.POST("/design", TemplateDesign)
	templateGroup.GET("/tags", TemplateTags)
}
