package worker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alexreagan/rabbit/server/service"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var treeRebuilderConfig *TreeReBuilderConfig

type TreeReBuilderRebuildConfig struct {
	Enable   bool `json:"enable"`
	Duration int  `json:"duration"`
}

type TreeReBuilderConfig struct {
	Rebuild TreeReBuilderRebuildConfig `json:"rebuild"`
}

func initTreeReBuilderConfigFromDB() (*TreeReBuilderConfig, error) {
	value, err := service.ParamService.Get("tree")
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, errors.New("tree is empty")
	}
	var config TreeReBuilderConfig
	err = json.Unmarshal([]byte(value), &config)
	return &config, err
}

type TreeReBuilder struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *TreeReBuilder) Ctx() context.Context {
	return s.ctx
}

func (s *TreeReBuilder) Close() {
	log.Infof("closing...")
	s.cancel()
	s.wg.Wait()
}

func (s *TreeReBuilder) Start() {
	//if viper.GetBool("tree.rebuild.enable") == false {
	//	return
	//}
	var err error
	treeRebuilderConfig, err = initTreeReBuilderConfigFromDB()
	if err != nil {
		log.Error(err)
		return
	}
	if treeRebuilderConfig.Rebuild.Enable == false {
		return
	}

	s.wg.Add(1)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(err)
			}
		}()
		s.StartReBuilder()
		defer s.wg.Done()
	}()
}

func (s *TreeReBuilder) StartReBuilder() {
	log.Println("[TreeReBuilder] StartReBuilder...")

	// 时间定时器启动
	//dur := viper.GetDuration("tree.rebuild.duration") * time.Second
	dur := time.Duration(treeRebuilderConfig.Rebuild.Duration) * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Println("ctx done")
			return
		case <-ticker.C:
			service.TagService.ReBuildGraph()
		}
	}
}

func InitTreeReBuilder() *TreeReBuilder {
	builder := &TreeReBuilder{}
	builder.ctx, builder.cancel = context.WithCancel(context.Background())
	return builder
}
