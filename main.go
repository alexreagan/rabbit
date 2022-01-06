package main

import (
	"flag"
	"fmt"
	_ "github.com/alexreagan/rabbit/docs"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server"
	"github.com/alexreagan/rabbit/server/model/alarm"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/caas"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/model/pub"
	"github.com/alexreagan/rabbit/server/model/sys"
	"github.com/alexreagan/rabbit/server/model/uic"
	"github.com/alexreagan/rabbit/server/worker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

// @title rabbit
// @description 一个简单的运维系统
func main() {

	g.BinaryName = BinaryName
	g.Version = Version
	g.GitCommit = GitCommit

	cfg := flag.String("c", "./config/cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Printf("Rabbit %s version %s, build %s\n", BinaryName, Version, GitCommit)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitLog()
	g.InitSentry()
	g.InitDBPool()
	if viper.GetBool("auto_migrate.enable") == true {
		//migrate database
		//g.Con().Uic.AutoMigrate(&uic.User{})
		//g.Con().Uic.AutoMigrate(&uic.Session{})
		g.Con().Portal.AutoMigrate(&uic.UserWhiteList{})
		g.Con().Portal.AutoMigrate(&uic.Role{})
		g.Con().Portal.AutoMigrate(&uic.Perm{})
		g.Con().Portal.AutoMigrate(&uic.RolePermRel{})
		g.Con().Portal.AutoMigrate(&uic.UserRoleRel{})
		//g.Con().Portal.AutoMigrate(&uic.Depart{})
		g.Con().Portal.AutoMigrate(&sys.Param{})
		g.Con().Portal.AutoMigrate(&sys.Menu{})
		g.Con().Portal.AutoMigrate(&sys.MenuPermission{})

		// node
		g.Con().Portal.AutoMigrate(&node.Node{})
		//g.Con().Portal.AutoMigrate(&node.NodeGroup{})
		//g.Con().Portal.AutoMigrate(&node.NodeGroupRel{})
		g.Con().Portal.AutoMigrate(&node.NodeTagRel{})
		// node apply request
		g.Con().Portal.AutoMigrate(&node.NodeApplyRequest{})

		// app
		g.Con().Portal.AutoMigrate(&app.Template{})
		g.Con().Portal.AutoMigrate(&app.Tag{})
		g.Con().Portal.AutoMigrate(&app.TagCategory{})

		// caas
		g.Con().Portal.AutoMigrate(&caas.WorkSpace{})
		g.Con().Portal.AutoMigrate(&caas.NameSpace{})
		g.Con().Portal.AutoMigrate(&caas.App{})
		g.Con().Portal.AutoMigrate(&caas.Service{})
		g.Con().Portal.AutoMigrate(&caas.NamespaceServiceRel{})
		g.Con().Portal.AutoMigrate(&caas.Port{})
		g.Con().Portal.AutoMigrate(&caas.ServicePortRel{})
		g.Con().Portal.AutoMigrate(&caas.Pod{})
		g.Con().Portal.AutoMigrate(&caas.ServicePodRel{})
		g.Con().Portal.AutoMigrate(&caas.ServiceTagRel{})
		g.Con().Portal.AutoMigrate(&alarm.Alarm{})

		// pub
		g.Con().Portal.AutoMigrate(&pub.Pub{})
	}

	// start gin server
	go server.Start()

	// sync nodes from kunyuan
	kunyuanSyncer := worker.InitKunYuanSyncer()
	kunyuanSyncer.Start()

	// sync nodes from caas
	caasSyncer := worker.InitCaasSyncer()
	caasSyncer.Start()

	// clean
	caasCleaner := worker.InitCaasCleaner()
	caasCleaner.Start()

	// tree ReBuilder
	treeReBuilder := worker.InitTreeReBuilder()
	treeReBuilder.Start()

	// prometheus syncer
	prometheusSyncer := worker.InitPrometheusSyncer()
	prometheusSyncer.Start()

	// process signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	select {
	case n := <-quit:
		log.Infof("receive signal %v, closing", n)
	case <-kunyuanSyncer.Ctx().Done():
		log.Infoln("kunyuanSyncer ctx done, closing")
	case <-caasSyncer.Ctx().Done():
		log.Infoln("caasSyncer ctx done, closing")
	case <-caasCleaner.Ctx().Done():
		log.Infoln("caasCleaner ctx done, closing")
	case <-treeReBuilder.Ctx().Done():
		log.Infoln("treeReBuilder ctx done, closing")
	case <-prometheusSyncer.Ctx().Done():
		log.Infoln("treeReBuilder ctx done, closing")
	}
	kunyuanSyncer.Close()
	caasSyncer.Close()
	caasCleaner.Close()
	treeReBuilder.Close()
	prometheusSyncer.Close()
}
