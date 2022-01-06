package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/service"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/net/httplib"
	"strings"
	"sync"
	"time"
)

type PrometheusConfig struct {
	Expr string `json:"expr"`
}

// 主机名
type NodeUnameInfoConfig PrometheusConfig

// 运行时间
type SumTimeInfoConfig PrometheusConfig

// 总内存
type NodeMemoryMemTotalBytesConfig PrometheusConfig

// 总核数
type NodeCpuCountConfig PrometheusConfig

// 5分钟负载
type NodeLoad5Config PrometheusConfig

// CPU使用率
type NodeCpuUsageConfig PrometheusConfig

// 内存使用率
type NodeMemoryAvailableBytesConfig PrometheusConfig

// 分区使用率
type NodeFileSystemSizeBytesConfig PrometheusConfig

// 最大读取
type MaxRateNodeDiskReadBytesConfig PrometheusConfig

// 最大写入
type MaxRateNodeDiskWrittenBytesConfig PrometheusConfig

// 连接数
type NodeNetStatTcpCurrEstabConfig PrometheusConfig

// Time Wait
type NodeSockstatTCPTWConfig PrometheusConfig

// 下载带宽
type MaxRateNodeNetWorkReceiveBytesConfig PrometheusConfig

// 上传带宽
type MaxRateNodeNetWorkTransmitBytesConfig PrometheusConfig

var prometheusSyncerConfig *PrometheusSyncerConfig

type PrometheusSyncerConfig struct {
	Enable                          bool                                  `json:"enable"`
	Addr                            string                                `json:"addr"`
	Duration                        int64                                 `json:"duration"`
	NodeUnameInfo                   NodeUnameInfoConfig                   `json:"nodeUnameInfo"`
	SumTimeInfo                     SumTimeInfoConfig                     `json:"sumTimeInfo"`
	NodeMemoryMemTotalBytes         NodeMemoryMemTotalBytesConfig         `json:"nodeMemoryMemTotalBytes"`
	NodeCpuCount                    NodeCpuCountConfig                    `json:"nodeCpuCount"`
	NodeLoad5                       NodeLoad5Config                       `json:"nodeLoad5"`
	NodeCpuUsage                    NodeCpuUsageConfig                    `json:"nodeCpuUsage"`
	NodeMemoryAvailableBytes        NodeMemoryAvailableBytesConfig        `json:"nodeMemoryAvailableBytes"`
	NodeFileSystemSizeBytes         NodeFileSystemSizeBytesConfig         `json:"nodeFileSystemSizeBytes"`
	MaxRateNodeDiskReadBytes        MaxRateNodeDiskReadBytesConfig        `json:"maxRateNodeDiskReadBytes"`
	MaxRateNodeDiskWrittenBytes     MaxRateNodeDiskWrittenBytesConfig     `json:"maxRateNodeDiskWrittenBytes"`
	NodeNetStatTcpCurrEstab         NodeNetStatTcpCurrEstabConfig         `json:"nodeNetStatTcpCurrEstab"`
	NodeSockstatTCPTW               NodeSockstatTCPTWConfig               `json:"nodeSockstatTCPTW"`
	MaxRateNodeNetWorkReceiveBytes  MaxRateNodeNetWorkReceiveBytesConfig  `json:"MaxRateNodeNetWorkReceiveBytes"`
	MaxRateNodeNetWorkTransmitBytes MaxRateNodeNetWorkTransmitBytesConfig `json:"MaxRateNodeNetWorkTransmitBytes"`
}

func loadPrometheusSyncerConfigFromDB() (*PrometheusSyncerConfig, error) {
	value, err := service.ParamService.Get("prometheus.syncer")
	if err != nil {
		return nil, err
	}
	var config PrometheusSyncerConfig
	err = json.Unmarshal([]byte(value), &config)
	return &config, err
}

type PrometheusSyncer struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *PrometheusSyncer) Ctx() context.Context {
	return s.ctx
}

func (s *PrometheusSyncer) Close() {
	log.Infoln("[PrometheusSyncer] closing...")
	s.cancel()
	s.wg.Wait()
	log.Infoln("[PrometheusSyncer] closed...")
}

func (s *PrometheusSyncer) Start() {
	var err error
	prometheusSyncerConfig, err = loadPrometheusSyncerConfigFromDB()
	if err != nil {
		log.Error(err)
		return
	}
	if prometheusSyncerConfig.Enable == false {
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
		s.StartSyncer()
		defer s.wg.Done()
	}()
}

type SyncNodeUnameInfoRespDataResultMetric struct {
	Name       string `json:"__name__"`
	DomainName string `json:"domainname"`
	Instance   string `json:"instance"`
	Job        string `json:"job"`
	Machine    string `json:"machine"`
	NodeName   string `json:"nodename"`
	Release    string `json:"release"`
	SysName    string `json:"sysname"`
	Version    string `json:"version"`
}

type SyncNodeUnameInfoRespDataResult struct {
	Metric SyncNodeUnameInfoRespDataResultMetric `json:"metric"`
	Value  []interface{}                         `json:"value"`
}

type SyncNodeUnameInfoRespData struct {
	Result     []SyncNodeUnameInfoRespDataResult `json:"result"`
	ResultType string                            `json:"resultType"`
}

type SyncNodeUnameInfoResp struct {
	Status string                    `json:"status"`
	Data   SyncNodeUnameInfoRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeUnameInfo() (*SyncNodeUnameInfoResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeUnameInfo...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "node_uname_info{}")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeUnameInfo SyncNodeUnameInfoResp
	e = json.Unmarshal([]byte(resp), &syncNodeUnameInfo)

	tx := g.Con().Portal
	for _, n := range syncNodeUnameInfo.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:      instance,
				Name:    n.Metric.NodeName,
				Machine: n.Metric.Machine,
				Release: n.Metric.Release,
				SysName: n.Metric.SysName,
				Version: n.Metric.Version,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:      instance,
				Name:    n.Metric.NodeName,
				Machine: n.Metric.Machine,
				Release: n.Metric.Release,
				SysName: n.Metric.SysName,
				Version: n.Metric.Version,
			})
		}
	}
	return &syncNodeUnameInfo, e
}

func (s *PrometheusSyncer) StartSyncer() {
	var err error

	// 时间定时器启动
	dur := time.Duration(prometheusSyncerConfig.Duration) * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Infoln("[PrometheusSyncer] ctx done")
			return
		case <-ticker.C:
			// load config
			prometheusSyncerConfig, err = loadPrometheusSyncerConfigFromDB()
			if err != nil {
				log.Error(err)
				return
			}
			if prometheusSyncerConfig.Enable == false {
				return
			}

			// start sync
			s.SyncNodeUnameInfo()
		}
	}
}

func InitPrometheusSyncer() *PrometheusSyncer {
	syncer := &PrometheusSyncer{}
	syncer.ctx, syncer.cancel = context.WithCancel(context.Background())
	return syncer
}
