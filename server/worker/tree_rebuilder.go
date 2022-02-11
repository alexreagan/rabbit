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

type TreeReBuilderConfig struct {
	Enabled  bool `json:"enabled"`
	Duration int  `json:"duration"`
}

func initTreeReBuilderConfigFromDB() (*TreeReBuilderConfig, error) {
	value, err := service.ParamService.Get("tree.rebuild")
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, errors.New("tree.rebuild is empty")
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
	log.Infoln("[TreeReBuilder] closing...")
	s.cancel()
	s.wg.Wait()
	log.Infoln("[TreeReBuilder] closed...")
}

func (s *TreeReBuilder) Start() {
	//if viper.GetBool("tree.rebuild.enabled") == false {
	//	return
	//}
	var err error
	treeRebuilderConfig, err = initTreeReBuilderConfigFromDB()
	if err != nil {
		log.Error(err)
		return
	}
	if treeRebuilderConfig.Enabled == false {
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
		treeRebuilderConfig, err = initTreeReBuilderConfigFromDB()
		if err != nil {
			log.Error(err)
			return
		}
		if treeRebuilderConfig.Enabled == false {
			return
		}
		s.StartReBuilder()
		defer s.wg.Done()
	}()
}

func (s *TreeReBuilder) StartReBuilder() {
	log.Infoln("[TreeReBuilder] StartReBuilder...")

	// 时间定时器启动
	//dur := viper.GetDuration("tree.rebuild.duration") * time.Second
	dur := time.Duration(treeRebuilderConfig.Duration) * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Infoln("[TreeReBuilder] ctx done")
			return
		case <-ticker.C:
			// ReBuildGraphV2
			service.TagService.ReBuildGraphV2()

			// build template graphs
			service.TemplateService.BuildGraphs()
		}
	}
}

func InitTreeReBuilder() *TreeReBuilder {
	builder := &TreeReBuilder{}
	builder.ctx, builder.cancel = context.WithCancel(context.Background())
	return builder
}
