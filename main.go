package main

import (
	"flag"
	"fmt"
	_ "github.com/alexreagan/rabbit/docs"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server"
	"github.com/alexreagan/rabbit/server/model/alert"
	"github.com/alexreagan/rabbit/server/model/caas"
	"github.com/alexreagan/rabbit/server/model/node"
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
// @description 一个简单的机器管理系统
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

		// host tag rel
		g.Con().Portal.AutoMigrate(&node.Host{})
		g.Con().Portal.AutoMigrate(&node.HostGroup{})
		g.Con().Portal.AutoMigrate(&node.HostGroupRel{})
		g.Con().Portal.AutoMigrate(&node.Tag{})
		g.Con().Portal.AutoMigrate(&node.TagCategory{})
		g.Con().Portal.AutoMigrate(&node.HostTagRel{})

		// host apply request
		g.Con().Portal.AutoMigrate(&node.HostApplyRequest{})

		// caas
		g.Con().Portal.AutoMigrate(&caas.WorkSpace{})
		g.Con().Portal.AutoMigrate(&caas.NameSpace{})
		g.Con().Portal.AutoMigrate(&caas.Service{})
		g.Con().Portal.AutoMigrate(&caas.NamespaceServiceRel{})
		g.Con().Portal.AutoMigrate(&caas.Port{})
		g.Con().Portal.AutoMigrate(&caas.ServicePortRel{})
		g.Con().Portal.AutoMigrate(&caas.Pod{})
		g.Con().Portal.AutoMigrate(&caas.ServicePodRel{})
		g.Con().Portal.AutoMigrate(&alert.Alert{})
	}

	// start gin server
	go server.Start()

	// sync hosts from kunyuan
	kunyuanSyncer := worker.InitKunYuanSyncer()
	kunyuanSyncer.Start()

	// sync hosts from caas
	caasSyncer := worker.InitCaasSyncer()
	caasSyncer.Start()

	// tree ReBuilder
	treeReBuilder := worker.InitTreeReBuilder()
	treeReBuilder.Start()

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
		log.Printf("receive signal %v, closing", n)
	case <-kunyuanSyncer.Ctx().Done():
		log.Println("kunyuanSyncer ctx done, closing")
	case <-caasSyncer.Ctx().Done():
		log.Println("caasSyncer ctx done, closing")
	case <-treeReBuilder.Ctx().Done():
		log.Println("treeReBuilder ctx done, closing")
	}
	kunyuanSyncer.Close()
	caasSyncer.Close()
	treeReBuilder.Close()
}
