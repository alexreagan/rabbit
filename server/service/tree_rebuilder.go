package service

import (
	"context"
	"github.com/alexreagan/rabbit/server/model/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
	"time"
)

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
	s.wg.Add(1)
	go func() {
		s.StartReBuilder()
		defer s.wg.Done()
	}()
}

func (s *TreeReBuilder) StartReBuilder() {
	log.Println("[TreeReBuilder] StartReBuilder...")

	// 时间定时器启动
	dur := viper.GetDuration("tree.rebuild.duration") * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Println("ctx done")
			return
		case <-ticker.C:
			node.HostGroup{}.ReBuildTree()
		}
	}
}

func InitTreeReBuilder() *TreeReBuilder {
	builder := &TreeReBuilder{}
	builder.ctx, builder.cancel = context.WithCancel(context.Background())
	return builder
}
