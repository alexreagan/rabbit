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

var caasCleanConfig *CaasCleanConfig

type CaasCleanConfig struct {
	Enable   bool `json:"enable"`
	Duration int  `json:"duration"`
}

func loadCaasCleanConfigFromDB() (*CaasCleanConfig, error) {
	value, err := service.ParamService.Get("caas.clean")
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, errors.New("caas.clean is empty")
	}
	var config CaasCleanConfig
	err = json.Unmarshal([]byte(value), &config)
	return &config, nil
}

type CaasCleaner struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *CaasCleaner) Ctx() context.Context {
	return s.ctx
}

func (s *CaasCleaner) Close() {
	log.Infof("closing...")
	s.cancel()
	s.wg.Wait()
}

func (s *CaasCleaner) Start() {
	s.wg.Add(1)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(err)
			}
		}()
		s.StartClean()
		defer s.wg.Done()
	}()
}

func (s *CaasCleaner) Clean() {
	log.Debugf("[CaasCleaner] Clean...")
	latestTime := service.CaasService.GetNameSpaceLatestTime()
	oneDayBeforeLatestTime := latestTime.AddDate(0, 0, -1)
	service.CaasService.DeleteNameSpaceBeforeTime(oneDayBeforeLatestTime)
	service.CaasService.DeleteServiceBeforeTime(oneDayBeforeLatestTime)
	service.CaasService.DeletePodBeforeTime(oneDayBeforeLatestTime)
	service.CaasService.DeletePortBeforeTime(oneDayBeforeLatestTime)
}

func (s *CaasCleaner) StartClean() {
	log.Debugf("[CaasCleaner] StartClean...")

	// load config
	var err error
	caasCleanConfig, err = loadCaasCleanConfigFromDB()
	if err != nil {
		log.Error(err)
		return
	}
	if caasCleanConfig.Enable == false {
		return
	}
	// 启动
	s.Clean()

	// 清理定时器启动
	//cleanDur := viper.GetDuration("caas.clean.duration") * time.Second
	cleanDur := time.Duration(caasCleanConfig.Duration) * time.Second
	cleanTicker := time.NewTicker(cleanDur)
	for {
		select {
		case <-s.ctx.Done():
			log.Println("ctx done")
			return
		case <-cleanTicker.C:
			// load config
			caasCleanConfig, err = loadCaasCleanConfigFromDB()
			if err != nil {
				log.Error(err)
				return
			}
			if caasCleanConfig.Enable == false {
				return
			}
			s.Clean()
		}
	}
}

func InitCaasCleaner() *CaasCleaner {
	syncer := &CaasCleaner{}
	syncer.ctx, syncer.cancel = context.WithCancel(context.Background())
	return syncer
}
