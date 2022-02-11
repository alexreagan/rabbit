// +build !windows

package server

import (
	"github.com/alexreagan/rabbit/server/controller/alarm"
	"github.com/alexreagan/rabbit/server/controller/app"
	"github.com/alexreagan/rabbit/server/controller/caas"
	"github.com/alexreagan/rabbit/server/controller/chart"
	"github.com/alexreagan/rabbit/server/controller/node"
	"github.com/alexreagan/rabbit/server/controller/pub"
	"github.com/alexreagan/rabbit/server/controller/sys"
	"github.com/alexreagan/rabbit/server/controller/uic"
	"github.com/alexreagan/rabbit/server/controller/wfe"
	"github.com/alexreagan/rabbit/server/utils"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"

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
	sys.Routes(r)
	app.Routes(r)
	node.Routes(r)
	alarm.Routes(r)
	chart.Routes(r)
	caas.Routes(r)
	pub.Routes(r)
	wfe.Routes(r)

	// start server graceful
	endless.ListenAndServe(viper.GetString("serv.addr"), r)
}
