package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/alarm"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/service"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"
const tTimeFormat = "2006-01-02T15:04:05"

// sync host information from kunyuan

var kunyuanConfig *KunYuanSyncerConfig

type KunYuanSyncerLoginConfig struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type KunYuanSyncerBaseConfig struct {
	URL             string   `json:"url"`
	PhysicalSystems []string `json:"physical_systems"`
	Duration        int      `json:"duration"`
}

type KunYuanSyncerMonitorConfig struct {
	URL             string   `json:"url"`
	PhysicalSystems []string `json:"physical_systems"`
	CloudPools      []string `json:"cloud_pools"`
	Duration        int      `json:"duration"`
}

type KunYuanSyncerAlarmConfig struct {
	URL          string `json:"url"`
	IntervalDays int    `json:"interval_days"`
	Duration     int    `json:"duration"`
}

type KunYuanSyncerConfig struct {
	Enable  bool                       `json:"enable"`
	Login   KunYuanSyncerLoginConfig   `json:"login"`
	Base    KunYuanSyncerBaseConfig    `json:"base"`
	Monitor KunYuanSyncerMonitorConfig `json:"monitor"`
	Alarm   KunYuanSyncerAlarmConfig   `json:"alarm"`
}

func initKunYuanSyncerConfigFromDB() (*KunYuanSyncerConfig, error) {
	value, err := service.ParamService.Get("kunyuan.syncer")
	if err != nil {
		return nil, err
	}
	var config KunYuanSyncerConfig
	err = json.Unmarshal([]byte(value), &config)
	return &config, err
}

type KunYuanSyncer struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *KunYuanSyncer) Ctx() context.Context {
	return s.ctx
}

func (s *KunYuanSyncer) Close() {
	log.Infoln("[KunYuanSyncer] closing...")
	s.cancel()
	s.wg.Wait()
	log.Infoln("[KunYuanSyncer] closed...")
}

func (s *KunYuanSyncer) Start() {
	//if viper.GetBool("kunyuan.syncer.enable") == false {
	//	return
	//}
	var err error
	kunyuanConfig, err = initKunYuanSyncerConfigFromDB()
	if err != nil {
		log.Error(err)
		return
	}
	if kunyuanConfig.Enable == false {
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
		s.StartBaseSyncer()
		defer s.wg.Done()
	}()

	s.wg.Add(1)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(err)
			}
		}()
		s.StartMonitorSyncer()
		defer s.wg.Done()
	}()
	s.wg.Add(1)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(err)
			}
		}()
		s.StartAlarmSyncer()
		defer s.wg.Done()
	}()
}

func (s *KunYuanSyncer) StartBaseSyncer() {
	log.Infoln("StartBaseSyncer")
	s.SyncBase()
	//dur := viper.GetDuration("kunyuan.syncer.base.duration") * time.Second
	dur := time.Duration(kunyuanConfig.Base.Duration) * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Infoln("[KunYuanSyncer] [BaseSyncer] ctx done")
			return
		case <-ticker.C:
			s.SyncBase()
		}
	}
}

func (s *KunYuanSyncer) StartMonitorSyncer() {
	log.Infoln("StartMonitorSyncer")
	s.SyncMonitor()
	//dur := viper.GetDuration("kunyuan.syncer.monitor.duration") * time.Second
	dur := time.Duration(kunyuanConfig.Monitor.Duration) * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Infoln("[KunYuanSyncer] [MonitorSyncer] ctx done")
			return
		case <-ticker.C:
			s.SyncMonitor()
		}
	}
}

func (s *KunYuanSyncer) StartAlarmSyncer() {
	log.Infoln("StartAlarmSyncer")
	s.SyncAlarm()
	//dur := viper.GetDuration("kunyuan.syncer.alarm.duration") * time.Second
	dur := time.Duration(kunyuanConfig.Alarm.Duration) * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Infoln("[KunYuanSyncer] [AlarmSyncer] ctx done")
			return
		case <-ticker.C:
			s.SyncAlarm()
		}
	}
}

type SyncerBaseRecord struct {
	Hostname          string `json:"hostname,omitempty"`
	ServPartName      string `json:"servPartName,omitempty"`
	SubSystemEnName   string `json:"subSysEnName,omitempty"`
	SubSystemAreaName string `json:"subAreaName,omitempty"`
	SubSystemCnName   string `json:"subSysName,omitempty"`
	LogicSystemCnName string `json:"sysName,omitempty"`
	Department        string `json:"department,omitempty"`
	ApplyUser         string `json:"applyUser,omitempty"`
	AreaName          string `json:"areaName,omitempty"`
	CpuNumber         int    `json:"cpuNum,omitempty"`
	DeployDate        string `json:"deployDate,omitempty"`
	DevAreaCode       string `json:"devAreaCode,omitempty"`
	FunDesc           string `json:"funDesc,omitempty"`
	InstanceID        string `json:"instanceId,omitempty"`
	ManagerA          string `json:"managerA,omitempty"`
	MemorySize        int    `json:"memSizeMB,omitempty"`
	OsVersion         string `json:"osVersion,omitempty"`
	PartTypeCode      string `json:"partTypeCode,omitempty"`
	ProdIp            string `json:"prodIp,omitempty"`
	ManIp             string `json:"manIp,omitempty"`
}

type SyncerBaseData struct {
	Current int                `json:"current"`
	Pages   int                `json:"pages"`
	Records []SyncerBaseRecord `json:"records"`
}

type SyncerBaseResult struct {
	Code int            `json:"code"`
	Data SyncerBaseData `json:"data"`
	Msg  string         `json:"msg"`
}

type SyncerMonitorServers struct {
	ApplyUser         string  `json:"applyUser,omitempty"`
	AreaName          string  `json:"areaName,omitempty"`
	CloudPoolCode     string  `json:"cloudPoolCode,omitempty"`
	CloudPoolName     string  `json:"cloudPoolName,omitempty"`
	CoreTotalNum      string  `json:"coreTotalNum,omitempty"`
	CpuNumber         string  `json:"cpuNum,omitempty"`
	CpuUsage          float64 `json:"cpuUsage,omitempty"`
	DatabaseVersion   string  `json:"databaseVersion,omitempty"`
	DeployDate        string  `json:"deployDate,omitempty"`
	DevAreaCode       string  `json:"devAreaCode,omitempty"`
	DevCenterName     string  `json:"devCenterName,omitempty"`
	DevTypeCode       string  `json:"devTypeCode,omitempty"`
	FsUsage           float64 `json:"fsUsage,omitempty"`
	FunDesc           string  `json:"funDesc,omitempty"`
	MemorySize        string  `json:"memSizeMB,omitempty"`
	MemoryUsage       float64 `json:"memUsage,omitempty"`
	ManagerA          string  `json:"mgrA,omitempty"`
	ManagerB          string  `json:"mgrB,omitempty"`
	Oracle            string  `json:"oracle,omitempty"`
	OsVersion         string  `json:"osVersion,omitempty"`
	PartType          string  `json:"partType,omitempty"`
	PartTypeCode      string  `json:"partTypeCode,omitempty"`
	ServSpaceCodeList string  `json:"servSpaceCodeList,omitempty"`
	ServSpaceNameList string  `json:"servSpaceNameList,omitempty"`
	SetupMode         string  `json:"setupMode,omitempty"`
	SrvStatus         string  `json:"srvStatus,omitempty"`
	Status            string  `json:"status,omitempty"`
	ServPartName      string  `json:"servPartName,omitempty"`
	SubSysCode        string  `json:"subSysCode,omitempty"`
	SubSystemName     string  `json:"subSystemName,omitempty"`
	SubSystemEnName   string  `json:"subSysEnName,omitempty"`
	SubSystemAreaName string  `json:"supAreaName,omitempty"`
	LogicSystemEnName string  `json:"sysEnname,omitempty"`
	LogicSystemCnName string  `json:"sysName,omitempty"`
	ProdIp            string  `json:"prodIp,omitempty"`
	VirtFcNum         string  `json:"virtFcNum,omitempty"`
	VirtNetNum        string  `json:"virtNetNum,omitempty"`
}

type SyncerMonitorData struct {
	Current     int                    `json:"current"`
	Pages       int                    `json:"pages"`
	Size        int                    `json:"size"`
	SearchCount int                    `json:"searchCount"`
	Servers     []SyncerMonitorServers `json:"servers"`
}

type SyncerMonitorResult struct {
	Code int               `json:"code"`
	Data SyncerMonitorData `json:"data"`
	Msg  string            `json:"msg"`
}

type KunYuanLoginResult struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	License      string `json:"license"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (s *KunYuanSyncer) Login() (*KunYuanLoginResult, error) {
	log.Infoln("Login...")
	lr := &KunYuanLoginResult{}
	//loginUrl := viper.GetString("kunyuan.syncer.login.url")
	loginUrl := kunyuanConfig.Login.URL
	payload := make(url.Values)
	//payload.Add("name", viper.GetString("kunyuan.syncer.login.user"))
	//payload.Add("username", viper.GetString("kunyuan.syncer.login.user"))
	//payload.Add("password", viper.GetString("kunyuan.syncer.login.password"))
	payload.Add("name", kunyuanConfig.Login.User)
	payload.Add("username", kunyuanConfig.Login.User)
	payload.Add("password", kunyuanConfig.Login.Password)
	payload.Add("grant_type", "password")
	payload.Add("scope", "server")
	req, _ := http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(payload.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Authorization", "Basic dGVzdDp0ZXN0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return lr, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, lr)
	log.Debugf("Login Result: %s", body)
	return lr, err
}

func (s *KunYuanSyncer) GetKunyuanBaseResult(abbr string, page int, lr *KunYuanLoginResult) (*SyncerBaseResult, error) {
	syncerResult := &SyncerBaseResult{}
	//syncUrl := viper.GetString("kunyuan.syncer.base.url")
	syncUrl := kunyuanConfig.Base.URL
	req, _ := http.NewRequest(http.MethodGet, syncUrl, nil)
	query := req.URL.Query()
	query.Add("areaName", "")
	query.Add("cloudPoolCode", "PRIVATE_CLOUD")
	query.Add("subSysEnNames", abbr)
	query.Add("current", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", lr.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return syncerResult, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return syncerResult, err
	}
	log.Debugf("[KunyuanSync] SyncBaseResult: %s", buf)
	err = json.Unmarshal(buf, &syncerResult)
	return syncerResult, err
}

func (s *KunYuanSyncer) SyncBase() {
	log.Infoln("SyncBase...")
	lr, err := s.Login()
	if err != nil || lr.AccessToken == "" || lr.ExpiresIn <= 0 {
		log.Errorln(err)
		return
	}

	//abbrs := viper.GetStringSlice("kunyuan.syncer.base.physical_systems")
	abbrs := kunyuanConfig.Base.PhysicalSystems
	for _, abbr := range abbrs {
		page := 1
		for {
			syncerResult, err := s.GetKunyuanBaseResult(abbr, page, lr)
			if err != nil {
				log.Errorln(err)
				break
			}

			// 保存数据
			tx := g.Con().Portal
			for _, record := range syncerResult.Data.Records {
				vt, _ := time.ParseInLocation(tTimeFormat, record.DeployDate, time.Local)
				h := &node.Node{
					Name:                 record.ServPartName,
					ApplyUser:            record.ApplyUser,
					AreaName:             record.AreaName,
					CpuNumber:            record.CpuNumber,
					MemorySize:           record.MemorySize,
					DeployDate:           vt,
					DevAreaCode:          record.DevAreaCode,
					FunDesc:              record.FunDesc,
					InstanceID:           record.InstanceID,
					IP:                   record.ProdIp,
					ManIp:                record.ManIp,
					ManagerA:             record.ManagerA,
					OsVersion:            record.OsVersion,
					PartTypeCode:         record.PartTypeCode,
					PhysicalSystem:       record.SubSystemEnName,
					PhysicalSystemCnName: record.SubSystemCnName,
					PhysicalSystemArea:   record.SubSystemAreaName,
					LogicSystemCnName:    record.LogicSystemCnName,
					UpdateTime:           time.Now(),
					State:                node.NodeStatusServicing,
				}

				var hh node.Node
				if record.ProdIp != "" {
					tx.Model(hh).Where(node.Node{IP: record.ProdIp}).First(&hh)
					if hh.ID == 0 {
						tx.Model(hh).Create(h)
					} else {
						tx.Model(hh).Where(node.Node{IP: record.ProdIp}).Updates(h)
					}
				} else if record.ServPartName != "" {
					tx.Model(hh).Where(node.Node{Name: record.ServPartName}).First(&hh)
					if hh.ID == 0 {
						tx.Model(hh).Create(h)
					} else {
						tx.Model(hh).Where(node.Node{Name: record.ServPartName}).Updates(h)
					}
				}
			}

			// 跳出循环
			if syncerResult.Data.Current >= syncerResult.Data.Pages {
				break
			}
			page = syncerResult.Data.Current + 1

			// sleep 1s
			time.Sleep(time.Second)
		}
	}
}

func (s *KunYuanSyncer) GetKunyuanMonitorResult(abbr string, page int, lr *KunYuanLoginResult) (*SyncerMonitorResult, error) {
	syncerResult := &SyncerMonitorResult{}
	//syncUrl := viper.GetString("kunyuan.syncer.monitor.url")
	syncUrl := kunyuanConfig.Monitor.URL
	req, _ := http.NewRequest(http.MethodGet, syncUrl, nil)
	query := req.URL.Query()
	query.Add("isMine", "1")
	//query.Add("cloudPoolName", cloudPool)
	query.Add("subSysEnName", abbr)
	query.Add("areaName", "")
	query.Add("prodIp", "")
	query.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", lr.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return syncerResult, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return syncerResult, err
	}
	log.Debugf("[KunyuanSync] SyncMonitorResult: %s", buf)
	err = json.Unmarshal(buf, &syncerResult)
	return syncerResult, err
}

type SyncerAlertDataAlert struct {
	AlertLevel    string `json:"alertLevel"`
	AlertName     string `json:"alertName"`
	CloudPoolName string `json:"cloudPoolName"`
	FiringTime    string `json:"firing_time"`
	ID            int64  `json:"id"`
	ProdIP        string `json:"prodIp"`
	Resolved      bool   `json:"resolved"`
	ResolvedTime  string `json:"resolved_time"`
	StrategyID    int64  `json:"stragetyId"`
	StrategyName  string `json:"strategyName"`
	StrategyType  string `json:"strategyType"`
	SubSysEnName  string `json:"subSysEnName"`
	SubSysName    string `json:"subSysName"`
	U1            string `json:"u1"`
	U2            string `json:"u2"`
}

type SyncerAlertData struct {
	Current     int                    `json:"current"`
	Pages       int                    `json:"pages"`
	Size        int                    `json:"size"`
	SearchCount int                    `json:"searchCount"`
	Alerts      []SyncerAlertDataAlert `json:"alerts"`
}

type SyncerAlertResult struct {
	Code int             `json:"code"`
	Data SyncerAlertData `json:"data"`
	Msg  string          `json:"msg"`
}

func (s *KunYuanSyncer) GetKunyuanAlertResult(page int, lr *KunYuanLoginResult) (*SyncerAlertResult, error) {
	syncerResult := &SyncerAlertResult{}
	//syncUrl := viper.GetString("kunyuan.syncer.alarm.url")
	syncUrl := kunyuanConfig.Alarm.URL
	req, _ := http.NewRequest(http.MethodGet, syncUrl, nil)
	query := req.URL.Query()
	query.Add("isMine", "1")
	query.Add("cloudPoolName", "")
	query.Add("subSysEnName", "")
	query.Add("strategyTypeName", "")
	query.Add("strategyName", "")
	query.Add("alertName", "")
	query.Add("resolved", "")
	query.Add("prodIp", "")
	query.Add("page", strconv.Itoa(page))
	now := time.Now()
	startTime := now.AddDate(0, 0, -30)
	//if viper.GetInt("kunyuan.syncer.alarm.interval_days") != 0 {
	//	startTime = now.AddDate(0, 0, viper.GetInt("kunyuan.syncer.alarm.interval_days")*(-1))
	//}
	if kunyuanConfig.Alarm.IntervalDays != 0 {
		startTime = now.AddDate(0, 0, kunyuanConfig.Alarm.IntervalDays*-1)
	}
	formatTimeStart := startTime.Format(timeFormat)
	query.Add("firingTimeStart", formatTimeStart)
	formatTimeEnd := now.Format(timeFormat)
	query.Add("firingTimeEnd", formatTimeEnd)
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", lr.AccessToken))
	log.Debugf("request url: %s, params: %+v, headers: %+v", req.URL, query, req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return syncerResult, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return syncerResult, err
	}
	log.Debugf("[KunyuanSync] SyncAlarmResult: %s", buf)
	err = json.Unmarshal(buf, &syncerResult)
	log.Debugf("[KunyuanSync] SyncAlarmObj: %+v", syncerResult)
	return syncerResult, err
}

func (s *KunYuanSyncer) SyncAlarm() {
	log.Infoln("SyncAlarm...")
	lr, err := s.Login()
	if err != nil || lr.AccessToken == "" || lr.ExpiresIn <= 0 {
		log.Errorln(err)
		return
	}

	page := 1
	for {
		syncerResult, err := s.GetKunyuanAlertResult(page, lr)
		if err != nil {
			log.Errorln(err)
			break
		}

		tx := g.Con().Portal
		for _, record := range syncerResult.Data.Alerts {
			firingTime, _ := time.ParseInLocation(timeFormat, record.FiringTime, time.Local)
			resolvedTime, _ := time.ParseInLocation(timeFormat, record.ResolvedTime, time.Local)
			alm := &alarm.Alarm{
				ID:            record.ID,
				AlertLevel:    record.AlertLevel,
				AlertName:     record.AlertName,
				CloudPoolName: record.CloudPoolName,
				FiringTime:    gtime.NewGTime(firingTime),
				ProdIP:        record.ProdIP,
				Resolved:      record.Resolved,
				ResolvedTime:  gtime.NewGTime(resolvedTime),
				StrategyID:    record.StrategyID,
				StrategyName:  record.StrategyName,
				StrategyType:  record.StrategyType,
				SubSysEnName:  record.SubSysEnName,
				SubSysName:    record.SubSysName,
				U1:            record.U1,
				U2:            record.U2,
				UpdateTime:    gtime.NewGTime(time.Now()),
			}

			var oalm alarm.Alarm
			if record.ID != 0 {
				tx.Model(oalm).Where(alarm.Alarm{ID: record.ID}).First(&oalm)
				if oalm.ID == 0 {
					tx.Model(oalm).Create(alm)
				} else {
					tx.Model(oalm).Where(alarm.Alarm{ID: record.ID}).Updates(alm)
				}
			}
		}

		// 跳出
		if syncerResult.Data.Current >= syncerResult.Data.Pages {
			break
		}
		page = syncerResult.Data.Current + 1

		// sleep 1s
		time.Sleep(time.Second)
	}
}

func (s *KunYuanSyncer) SyncMonitor() {
	log.Infoln("SyncMonitor...")
	lr, err := s.Login()
	if err != nil || lr.AccessToken == "" || lr.ExpiresIn <= 0 {
		log.Errorln(err)
		return
	}

	//abbrs := viper.GetStringSlice("kunyuan.syncer.monitor.physical_systems")
	abbrs := kunyuanConfig.Monitor.PhysicalSystems
	for _, abbr := range abbrs {
		page := 1
		for {
			syncerResult, err := s.GetKunyuanMonitorResult(abbr, page, lr)
			if err != nil {
				log.Errorln(err)
				break
			}

			tx := g.Con().Portal
			for _, record := range syncerResult.Data.Servers {
				coreTotalNum, _ := strconv.Atoi(record.CoreTotalNum)
				h := &node.Node{
					IP:                   record.ProdIp,
					Name:                 record.ServPartName,
					CloudPoolCode:        record.CloudPoolCode,
					CloudPoolName:        record.CloudPoolName,
					CoreTotalNum:         coreTotalNum,
					CpuUsage:             record.CpuUsage,
					DatabaseVersion:      record.DatabaseVersion,
					DevCenterName:        record.DevCenterName,
					DevTypeCode:          record.DevTypeCode,
					FsUsage:              record.FsUsage,
					MemoryUsage:          record.MemoryUsage,
					ManagerA:             record.ManagerA,
					ManagerB:             record.ManagerB,
					Oracle:               record.Oracle,
					PartType:             record.PartType,
					PartTypeCode:         record.PartTypeCode,
					ServSpaceCodeList:    record.ServSpaceCodeList,
					ServSpaceNameList:    record.ServSpaceNameList,
					SetupMode:            record.SetupMode,
					SrvStatus:            record.SrvStatus,
					Status:               record.Status,
					PhysicalSystem:       record.SubSysCode,
					PhysicalSystemEnName: record.SubSystemEnName,
					PhysicalSystemCnName: record.SubSystemName,
					PhysicalSystemArea:   record.SubSystemAreaName,
					LogicSystem:          record.LogicSystemEnName,
					LogicSystemCnName:    record.LogicSystemCnName,
					VirtFcNum:            record.VirtFcNum,
					VirtNetNum:           record.VirtNetNum,
					UpdateTime:           time.Now(),
					State:                node.NodeStatusServicing,
					//ApplyUser:            record.ApplyUser,
					//AreaName:             record.AreaName,
					//CpuNumber:            record.CpuNumber,
					//MemorySize:           record.MemorySize,
					//DevAreaCode:          record.DevAreaCode,
					//DeployDate:           vt,
					//FunDesc:             record.FunDesc,
					//ProdIp:               record.ProdIp,
					//OsVersion:            record.OsVersion,
				}

				var hh node.Node
				if record.ProdIp != "" {
					tx.Model(hh).Where(node.Node{IP: record.ProdIp}).First(&hh)
					if hh.ID == 0 {
						tx.Model(hh).Create(h)
					} else {
						tx.Model(hh).Where(node.Node{IP: record.ProdIp}).Updates(h)
					}
				} else if record.ServPartName != "" {
					tx.Model(hh).Where(node.Node{Name: record.ServPartName}).First(&hh)
					if hh.ID == 0 {
						tx.Model(hh).Create(h)
					} else {
						tx.Model(hh).Where(node.Node{Name: record.ServPartName}).Updates(h)
					}
				}
			}

			// 跳出
			if syncerResult.Data.Current >= syncerResult.Data.Pages {
				break
			}
			page = syncerResult.Data.Current + 1

			// sleep 1s
			time.Sleep(time.Second)
		}
	}
}

func InitKunYuanSyncer() *KunYuanSyncer {
	syncer := &KunYuanSyncer{}
	syncer.ctx, syncer.cancel = context.WithCancel(context.Background())
	return syncer
}
