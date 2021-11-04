package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// sync host information from kunyuan

type KunYuanSyncer struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *KunYuanSyncer) Ctx() context.Context {
	return s.ctx
}

func (s *KunYuanSyncer) Close() {
	log.Infof("closing...")
	s.cancel()
	s.wg.Wait()
}

func (s *KunYuanSyncer) Start() {
	s.wg.Add(1)
	go func() {
		s.StartBaseSyncer()
		defer s.wg.Done()
	}()

	s.wg.Add(1)
	go func() {
		s.StartMonitorSyncer()
		defer s.wg.Done()
	}()
}

func (s *KunYuanSyncer) StartBaseSyncer() {
	log.Println("StartBaseSyncer")
	s.SyncBase()
	dur := viper.GetDuration("kunyuan_syncer.base.duration") * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Println("ctx done")
			return
		case <-ticker.C:
			s.SyncBase()
		}
	}
}

func (s *KunYuanSyncer) StartMonitorSyncer() {
	log.Println("StartMonitorSyncer")
	s.SyncMonitor()
	dur := viper.GetDuration("kunyuan_syncer.monitor.duration") * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Println("ctx done")
			return
		case <-ticker.C:
			s.SyncMonitor()
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
	InstanceId        string `json:"instanceId,omitempty"`
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
	log.Println("Login...")
	lr := &KunYuanLoginResult{}
	loginUrl := viper.GetString("kunyuan_syncer.login.url")
	payload := make(url.Values)
	payload.Add("name", viper.GetString("kunyuan_syncer.login.user"))
	payload.Add("username", viper.GetString("kunyuan_syncer.login.user"))
	payload.Add("password", viper.GetString("kunyuan_syncer.login.password"))
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
	log.Printf("Login Result: %s", body)
	return lr, err
}

func (s *KunYuanSyncer) GetKunyuanBaseResult(abbr string, page int, lr *KunYuanLoginResult) (*SyncerBaseResult, error) {
	syncerResult := &SyncerBaseResult{}
	syncUrl := viper.GetString("kunyuan_syncer.base.url")
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
	log.Printf("[KunyuanSync] SyncBaseResult: %s", buf)
	err = json.Unmarshal(buf, &syncerResult)
	return syncerResult, err
}

func (s *KunYuanSyncer) SyncBase() {
	log.Println("SyncBase...")
	lr, err := s.Login()
	if err != nil || lr.AccessToken == "" || lr.ExpiresIn <= 0 {
		log.Errorln(err)
		return
	}

	abbrs := viper.GetStringSlice("kunyuan_syncer.base.physical_systems")
	for _, abbr := range abbrs {
		page := 1
		for {
			syncerResult, err := s.GetKunyuanBaseResult(abbr, page, lr)
			if err != nil {
				log.Errorln(err)
				break
			}
			if syncerResult.Data.Current > syncerResult.Data.Pages {
				break
			}
			page = syncerResult.Data.Current + 1

			db := g.Con().Portal
			for _, record := range syncerResult.Data.Records {
				vt, _ := time.ParseInLocation("2006-01-02T15:04:05", record.DeployDate, time.Local)
				h := &node.Host{
					Name:                 record.ServPartName,
					ApplyUser:            record.ApplyUser,
					AreaName:             record.AreaName,
					CpuNumber:            record.CpuNumber,
					MemorySize:           record.MemorySize,
					DeployDate:           vt,
					DevAreaCode:          record.DevAreaCode,
					FunDesc:              record.FunDesc,
					InstanceId:           record.InstanceId,
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
				}

				var hh node.Host
				if record.ProdIp != "" {
					db.Table(hh.TableName()).Where(node.Host{IP: record.ProdIp}).First(&hh)
					if hh.ID == 0 {
						db.Table(hh.TableName()).Create(h)
					} else {
						db.Table(hh.TableName()).Where(node.Host{IP: record.ProdIp}).Updates(h)
					}
				} else if record.ServPartName != "" {
					db.Table(hh.TableName()).Where(node.Host{Name: record.ServPartName}).First(&hh)
					if hh.ID == 0 {
						db.Table(hh.TableName()).Create(h)
					} else {
						db.Table(hh.TableName()).Where(node.Host{Name: record.ServPartName}).Updates(h)
					}
				}
			}

			// sleep 1s
			time.Sleep(time.Second)
		}
	}
}

func (s *KunYuanSyncer) GetKunyuanMonitorResult(abbr string, page int, lr *KunYuanLoginResult) (*SyncerMonitorResult, error) {
	syncerResult := &SyncerMonitorResult{}
	syncUrl := viper.GetString("kunyuan_syncer.monitor.url")
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
	log.Printf("[KunyuanSync] SyncMonitorResult: %s", buf)
	err = json.Unmarshal(buf, &syncerResult)
	return syncerResult, err
}

func (s *KunYuanSyncer) SyncMonitor() {
	log.Println("SyncMonitor...")
	lr, err := s.Login()
	if err != nil || lr.AccessToken == "" || lr.ExpiresIn <= 0 {
		log.Errorln(err)
		return
	}

	abbrs := viper.GetStringSlice("kunyuan_syncer.monitor.physical_systems")
	//cloudPools := viper.GetStringSlice("syncer.monitor.cloud_pools")
	for _, abbr := range abbrs {
		page := 1
		for {
			syncerResult, err := s.GetKunyuanMonitorResult(abbr, page, lr)
			if err != nil {
				log.Errorln(err)
				break
			}
			if syncerResult.Data.Current > syncerResult.Data.Pages {
				break
			}
			page = syncerResult.Data.Current + 1

			db := g.Con().Portal
			for _, record := range syncerResult.Data.Servers {
				coreTotalNum, _ := strconv.Atoi(record.CoreTotalNum)
				h := &node.Host{
					IP:   record.ProdIp,
					Name: record.ServPartName,
					//ApplyUser:            record.ApplyUser,
					//AreaName:             record.AreaName,
					//CpuNumber:            record.CpuNumber,
					//MemorySize:           record.MemorySize,
					//DevAreaCode:          record.DevAreaCode,
					//DeployDate:           vt,
					//FunDesc:             record.FunDesc,
					//ProdIp:               record.ProdIp,
					//OsVersion:            record.OsVersion,
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
				}

				var hh node.Host
				if record.ProdIp != "" {
					db.Table(hh.TableName()).Where(node.Host{IP: record.ProdIp}).First(&hh)
					if hh.ID == 0 {
						db.Table(hh.TableName()).Create(h)
					} else {
						db.Table(hh.TableName()).Where(node.Host{IP: record.ProdIp}).Updates(h)
					}
				} else if record.ServPartName != "" {
					db.Table(hh.TableName()).Where(node.Host{Name: record.ServPartName}).First(&hh)
					if hh.ID == 0 {
						db.Table(hh.TableName()).Create(h)
					} else {
						db.Table(hh.TableName()).Where(node.Host{Name: record.ServPartName}).Updates(h)
					}
				}
			}

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
