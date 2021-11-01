// +build !windows

package server

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"rabbit/server/controller/node"
	"rabbit/server/controller/portal"
	"rabbit/server/controller/uic"
	"rabbit/server/utils"

	//_ "rabbit/docs"
	"github.com/fvbock/endless"
)

func Start() {

	// gin middleware
	r := gin.Default()
	r.Use(utils.CORS())
	r.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	// router
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// set routers
	uic.Routes(r)
	portal.Routes(r)
	node.Routes(r)

	// start server graceful
	endless.ListenAndServe(viper.GetString("serv.addr"), r)
}
