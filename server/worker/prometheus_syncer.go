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
type NodeBootTimeConfig PrometheusConfig

// 总内存
type NodeMemoryMemTotalBytesConfig PrometheusConfig

// 总核数
type NodeCpuSecondsTotalConfig PrometheusConfig

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
	SumTimeInfo                     NodeBootTimeConfig                    `json:"sumTimeInfo"`
	NodeMemoryMemTotalBytes         NodeMemoryMemTotalBytesConfig         `json:"nodeMemoryMemTotalBytes"`
	NodeCpuCount                    NodeCpuSecondsTotalConfig             `json:"nodeCpuCount"`
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
	var syncNodeUnameInfoResp SyncNodeUnameInfoResp
	e = json.Unmarshal([]byte(resp), &syncNodeUnameInfoResp)

	tx := g.Con().Portal
	for _, n := range syncNodeUnameInfoResp.Data.Result {
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
	return &syncNodeUnameInfoResp, e
}

// 运行时间
type SyncNodeBootTimeRespDataResultMetric struct {
	Instance string `json:"instance"`
}

type SyncNodeBootTimeRespDataResult struct {
	Metric SyncNodeBootTimeRespDataResultMetric `json:"metric"`
	Value  []interface{}                        `json:"value"`
}

type SyncNodeBootTimeRespData struct {
	ResultType string                           `json:"resultType"`
	Result     []SyncNodeBootTimeRespDataResult `json:"result"`
}

type SyncNodeBootTimeResp struct {
	Status string                   `json:"status"`
	Data   SyncNodeBootTimeRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeBootTime() (*SyncNodeBootTimeResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeBootTimeInfo...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "sum(time() - node_boot_time_seconds{})by(instance)")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeBootTimeResp SyncNodeBootTimeResp
	e = json.Unmarshal([]byte(resp), &syncNodeBootTimeResp)

	tx := g.Con().Portal
	for _, n := range syncNodeBootTimeResp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		bootTime, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:       instance,
				BootTime: bootTime,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:       instance,
				BootTime: bootTime,
			})
		}
	}
	return &syncNodeBootTimeResp, e
}

// 总内存
type SyncNodeMemoryMemTotalBytesRespDataResultMetric struct {
	Instance string `json:"instance"`
	//Job      string `json:"job"`
}

type SyncNodeMemoryMemTotalBytesRespDataResult struct {
	Metric SyncNodeMemoryMemTotalBytesRespDataResultMetric `json:"metric"`
	Value  []interface{}                                   `json:"value"`
}

type SyncNodeMemoryMemTotalBytesRespData struct {
	ResultType string                                      `json:"resultType"`
	Result     []SyncNodeMemoryMemTotalBytesRespDataResult `json:"result"`
}

type SyncNodeMemoryMemTotalBytesResp struct {
	Status string                              `json:"status"`
	Data   SyncNodeMemoryMemTotalBytesRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeMemoryMemTotalBytes() (*SyncNodeMemoryMemTotalBytesResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeMemoryMemTotalBytes...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "node_memory_MemTotal_bytes{} - 0")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeMemoryMemTotalBytesResp SyncNodeMemoryMemTotalBytesResp
	e = json.Unmarshal([]byte(resp), &syncNodeMemoryMemTotalBytesResp)

	tx := g.Con().Portal
	for _, n := range syncNodeMemoryMemTotalBytesResp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		memTotalBytes, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:            instance,
				MemTotalBytes: memTotalBytes,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:            instance,
				MemTotalBytes: memTotalBytes,
			})
		}
	}
	return &syncNodeMemoryMemTotalBytesResp, e
}

// 总核数
type SyncNodeCpuCountRespDataResultMetric struct {
	Instance string `json:"instance"`
}

type SyncNodeCpuCountRespDataResult struct {
	Metric SyncNodeCpuCountRespDataResultMetric `json:"metric"`
	Value  []interface{}                        `json:"value"`
}

type SyncNodeCpuCountRespData struct {
	ResultType string                           `json:"resultType"`
	Result     []SyncNodeCpuCountRespDataResult `json:"result"`
}

type SyncNodeCpuCountResp struct {
	Status string                   `json:"status"`
	Data   SyncNodeCpuCountRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeCpuCount() (*SyncNodeCpuCountResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeCpuCount...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "count(node_cpu_seconds_total{mode='system'}) by (instance)")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeCpuCountResp SyncNodeCpuCountResp
	e = json.Unmarshal([]byte(resp), &syncNodeCpuCountResp)

	tx := g.Con().Portal
	for _, n := range syncNodeCpuCountResp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		cpuCount, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:       instance,
				CpuCount: cpuCount,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:       instance,
				CpuCount: cpuCount,
			})
		}
	}
	return &syncNodeCpuCountResp, e
}

// 5分钟负载
type SyncNodeLoad5RespDataResultMetric struct {
	Instance string `json:"instance"`
}

type SyncNodeLoad5RespDataResult struct {
	Metric SyncNodeLoad5RespDataResultMetric `json:"metric"`
	Value  []interface{}                     `json:"value"`
}

type SyncNodeLoad5RespData struct {
	ResultType string                        `json:"resultType"`
	Result     []SyncNodeLoad5RespDataResult `json:"result"`
}

type SyncNodeLoad5Resp struct {
	Status string                `json:"status"`
	Data   SyncNodeLoad5RespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeLoad5() (*SyncNodeLoad5Resp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeLoad5...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "node_load5{}")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeLoad5Resp SyncNodeLoad5Resp
	e = json.Unmarshal([]byte(resp), &syncNodeLoad5Resp)

	tx := g.Con().Portal
	for _, n := range syncNodeLoad5Resp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		load5, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:    instance,
				Load5: load5,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:    instance,
				Load5: load5,
			})
		}
	}
	return &syncNodeLoad5Resp, e
}

// CPU使用率
type SyncNodeCPUAvailableRespDataResultMetric struct {
	Instance string `json:"instance"`
}

type SyncNodeCPUAvailableRespDataResult struct {
	Metric SyncNodeCPUAvailableRespDataResultMetric `json:"metric"`
	Value  []interface{}                            `json:"value"`
}

type SyncNodeCPUAvailableRespData struct {
	ResultType string                               `json:"resultType"`
	Result     []SyncNodeCPUAvailableRespDataResult `json:"result"`
}

type SyncNodeCPUAvailableResp struct {
	Status string                       `json:"status"`
	Data   SyncNodeCPUAvailableRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeCPUAvailable() (*SyncNodeCPUAvailableResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeCPUAvailable...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "(1 - avg(rate(node_cpu_seconds_total{mode='idle'}[5m])) by (instance)) * 100")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeCPUAvailableResp SyncNodeCPUAvailableResp
	e = json.Unmarshal([]byte(resp), &syncNodeCPUAvailableResp)

	tx := g.Con().Portal
	for _, n := range syncNodeCPUAvailableResp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		cpuAvailable, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:           instance,
				CPUAvailable: cpuAvailable,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:           instance,
				CPUAvailable: cpuAvailable,
			})
		}
	}
	return &syncNodeCPUAvailableResp, e
}

// 内存使用率
type SyncNodeMemAvailableRespDataResultMetric struct {
	Instance string `json:"instance"`
}

type SyncNodeMemAvailableRespDataResult struct {
	Metric SyncNodeMemAvailableRespDataResultMetric `json:"metric"`
	Value  []interface{}                            `json:"value"`
}

type SyncNodeMemAvailableRespData struct {
	ResultType string                               `json:"resultType"`
	Result     []SyncNodeMemAvailableRespDataResult `json:"result"`
}

type SyncNodeMemAvailableResp struct {
	Status string                       `json:"status"`
	Data   SyncNodeMemAvailableRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeMemAvailable() (*SyncNodeMemAvailableResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeMemAvailable...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "(1 - (node_memory_MemAvailable_bytes{} / (node_memory_MemTotal_bytes{}))) * 100")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeMemAvailableBytesResp SyncNodeMemAvailableResp
	e = json.Unmarshal([]byte(resp), &syncNodeMemAvailableBytesResp)

	tx := g.Con().Portal
	for _, n := range syncNodeMemAvailableBytesResp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		memAvailable, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:           instance,
				MemAvailable: memAvailable,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:           instance,
				MemAvailable: memAvailable,
			})
		}
	}
	return &syncNodeMemAvailableBytesResp, e
}

// 分区使用率
type SyncNodeFileSystemAvailableRespDataResultMetric struct {
	Instance string `json:"instance"`
}

type SyncNodeFileSystemAvailableRespDataResult struct {
	Metric SyncNodeFileSystemAvailableRespDataResultMetric `json:"metric"`
	Value  []interface{}                                   `json:"value"`
}

type SyncNodeFileSystemAvailableRespData struct {
	ResultType string                                      `json:"resultType"`
	Result     []SyncNodeFileSystemAvailableRespDataResult `json:"result"`
}

type SyncNodeFileSystemAvailableResp struct {
	Status string                              `json:"status"`
	Data   SyncNodeFileSystemAvailableRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeFileSystemAvailable() (*SyncNodeFileSystemAvailableResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeFileSystemAvailable...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "max((node_filesystem_size_bytes{fstype=~'ext.?|xfs'} - node_filesystem_free_bytes{fstype=~'ext.?|xfs'}) * 100 / (node_filesystem_avail_bytes{fstype=~'ext.?|xfs'} + (node_filesystem_size_bytes{fstype=~'ext.?|xfs'} - node_filesystem_free_bytes{fstype=~'ext.?|xfs'}))) by (instance)")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeFileSystemAvailableResp SyncNodeFileSystemAvailableResp
	e = json.Unmarshal([]byte(resp), &syncNodeFileSystemAvailableResp)

	tx := g.Con().Portal
	for _, n := range syncNodeFileSystemAvailableResp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		fileSystemAvailable, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:                  instance,
				FileSystemAvailable: fileSystemAvailable,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:                  instance,
				FileSystemAvailable: fileSystemAvailable,
			})
		}
	}
	return &syncNodeFileSystemAvailableResp, e
}

// 最大读取
type SyncNodeMaxDiskReadBytesRespDataResultMetric struct {
	Instance string `json:"instance"`
}

type SyncNodeMaxDiskReadBytesRespDataResult struct {
	Metric SyncNodeMaxDiskReadBytesRespDataResultMetric `json:"metric"`
	Value  []interface{}                                `json:"value"`
}

type SyncNodeMaxDiskReadBytesRespData struct {
	ResultType string                                   `json:"resultType"`
	Result     []SyncNodeMaxDiskReadBytesRespDataResult `json:"result"`
}

type SyncNodeMaxDiskReadBytesResp struct {
	Status string                           `json:"status"`
	Data   SyncNodeMaxDiskReadBytesRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeMaxDiskReadBytes() (*SyncNodeMaxDiskReadBytesResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeMaxDiskReadBytes...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "max(rate(node_disk_read_bytes_total{}[5m])) by (instance)")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeMaxDiskReadBytesResp SyncNodeMaxDiskReadBytesResp
	e = json.Unmarshal([]byte(resp), &syncNodeMaxDiskReadBytesResp)

	tx := g.Con().Portal
	for _, n := range syncNodeMaxDiskReadBytesResp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		maxDiskReadBytes, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:               instance,
				MaxDiskReadBytes: maxDiskReadBytes,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:               instance,
				MaxDiskReadBytes: maxDiskReadBytes,
			})
		}
	}
	return &syncNodeMaxDiskReadBytesResp, e
}

// 最大写入
type SyncNodeMaxDiskWrittenBytesRespDataResultMetric struct {
	Instance string `json:"instance"`
}

type SyncNodeMaxDiskWrittenBytesRespDataResult struct {
	Metric SyncNodeMaxDiskWrittenBytesRespDataResultMetric `json:"metric"`
	Value  []interface{}                                   `json:"value"`
}

type SyncNodeMaxDiskWrittenBytesRespData struct {
	ResultType string                                      `json:"resultType"`
	Result     []SyncNodeMaxDiskWrittenBytesRespDataResult `json:"result"`
}

type SyncNodeMaxDiskWrittenBytesResp struct {
	Status string                              `json:"status"`
	Data   SyncNodeMaxDiskWrittenBytesRespData `json:"data"`
}

func (s *PrometheusSyncer) SyncNodeMaxDiskWrittenBytes() (*SyncNodeMaxDiskWrittenBytesResp, error) {
	log.Infoln("[PrometheusSyncer] SyncNodeMaxDiskReadBytes...")
	uri := fmt.Sprintf("%s/api/v1/query", prometheusSyncerConfig.Addr)
	req := httplib.Get(uri)
	req.Param("query", "max(rate(node_disk_written_bytes_total{}[5m])) by (instance)")
	resp, e := req.String()

	if e != nil {
		log.Errorln(e.Error())
		return nil, e
	}
	var syncNodeMaxDiskWrittenBytesResp SyncNodeMaxDiskWrittenBytesResp
	e = json.Unmarshal([]byte(resp), &syncNodeMaxDiskWrittenBytesResp)

	tx := g.Con().Portal
	for _, n := range syncNodeMaxDiskWrittenBytesResp.Data.Result {
		instance := strings.Split(n.Metric.Instance, ":")[0]
		maxDiskWrittenBytes, ok := n.Value[1].(string)
		if !ok {
			break
		}
		var nod node.Node
		tx.Model(node.Node{}).Where(node.Node{IP: instance}).First(&nod)
		if nod.ID == 0 {
			tx.Model(node.Node{}).Create(&node.Node{
				IP:                  instance,
				MaxDiskWrittenBytes: maxDiskWrittenBytes,
			})
		} else {
			tx.Model(node.Node{}).Where(node.Node{IP: instance}).Updates(&node.Node{
				IP:                  instance,
				MaxDiskWrittenBytes: maxDiskWrittenBytes,
			})
		}
	}
	return &syncNodeMaxDiskWrittenBytesResp, e
}

func (s *PrometheusSyncer) StartSyncer() {
	var err error
	prometheusSyncerConfig, err = loadPrometheusSyncerConfigFromDB()
	if err != nil {
		log.Error(err)
		return
	}
	if prometheusSyncerConfig.Enable == false {
		return
	}

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
			s.SyncNodeBootTime()
			s.SyncNodeCpuCount()
			s.SyncNodeMemoryMemTotalBytes()
			s.SyncNodeLoad5()
			s.SyncNodeCPUAvailable()
			s.SyncNodeMemAvailable()
			s.SyncNodeFileSystemAvailable()
			s.SyncNodeMaxDiskReadBytes()
			s.SyncNodeMaxDiskWrittenBytes()
		}
	}
}

func InitPrometheusSyncer() *PrometheusSyncer {
	syncer := &PrometheusSyncer{}
	syncer.ctx, syncer.cancel = context.WithCancel(context.Background())
	return syncer
}
