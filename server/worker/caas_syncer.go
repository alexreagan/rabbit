package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/caas"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type CaasSyncer struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *CaasSyncer) Ctx() context.Context {
	return s.ctx
}

func (s *CaasSyncer) Close() {
	log.Infof("closing...")
	s.cancel()
	s.wg.Wait()
}

func (s *CaasSyncer) Start() {
	if viper.GetBool("caas_syncer.enable") == false {
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

func (s *CaasSyncer) Clean() {
	log.Debugf("[CaasSyncer] Clean...")
	latestTime := service.CaasService.GetNameSpaceLatestTime()
	oneDayBeforeLatestTime := latestTime.AddDate(0, 0, -1)
	service.CaasService.DeleteNameSpaceBeforeTime(oneDayBeforeLatestTime)
	service.CaasService.DeleteServiceBeforeTime(oneDayBeforeLatestTime)
	service.CaasService.DeletePodBeforeTime(oneDayBeforeLatestTime)
	service.CaasService.DeletePortBeforeTime(oneDayBeforeLatestTime)
}

func (s *CaasSyncer) StartClean() {
	log.Debugf("[CaasSyncer] StartClean...")
	// 启动
	s.Clean()

	// 清理定时器启动
	cleanDur := viper.GetDuration("caas_syncer.clean.duration") * time.Second
	cleanTicker := time.NewTicker(cleanDur)
	for {
		select {
		case <-s.ctx.Done():
			log.Println("ctx done")
			return
		case <-cleanTicker.C:
			s.Clean()
		}
	}
}

func (s *CaasSyncer) StartSyncer() {
	log.Println("[CaasSyncer] StartReBuilder...")
	// 启动
	s.Sync()

	// 时间定时器启动
	dur := viper.GetDuration("caas_syncer.sync.duration") * time.Second
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-s.ctx.Done():
			log.Println("ctx done")
			return
		case <-ticker.C:
			s.Sync()
		}
	}
}

func RemoveRepeated(arr []int64) []int64 {
	var result []int64
	m := make(map[int64]bool) //map的值不重要
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

func (s *CaasSyncer) Sync() {
	// workspace
	wsResult, err := GetWorkSpace()
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, ws := range wsResult.Data {
		log.Debugf("%+v", ws)
		// 更新数据库
		UpdateWorkspace(&ws)

		// namespace
		nsResult, err := GetNameSpace(ws.ID)
		if err != nil {
			log.Errorln(err)
			return
		}
		for _, ns := range nsResult.Data {
			// 更新数据库
			UpdateNamespace(&ns)
			// app
			UpdateApps(&ns)

			// service
			services, err := GetService(ns)
			if err != nil {
				log.Errorln(err)
				return
			}

			for _, ser := range services.Data {
				// 更新数据库
				UpdateService(&ser)

				// 更新namespace与service的关系
				UpdateNamespaceServiceRel(&ns, &ser)

				// port，以及service和pod的关系
				UpdateServicePorts(&ser)

				// 获取namespace,service下的pods
				pods, err := GetPod(&ns, &ser)
				if err != nil {
					log.Errorln(err)
					return
				}
				UpdatePods(&ser, pods)
			}
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second)
	}
}

type CaasLoginResultDataCluster struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type CaasLoginResultData struct {
	Cluster  []CaasLoginResultDataCluster `json:"cluster"`
	LdapUer  int64                        `json:"ldapUser"`
	Token    string                       `json:"token"`
	UserId   int64                        `json:"userId"`
	UserType int64                        `json:"userType"`
	Username string                       `json:"username"`
}

type CaasLoginResult struct {
	Code int64               `json:"code"`
	Data CaasLoginResultData `json:"data"`
	Msg  string              `json:"msg"`
}

// Login 登录
func Login() (*CaasLoginResult, error) {
	log.Debugf("[CaasSyncer] Login...")
	lr := &CaasLoginResult{}
	loginUrl := viper.GetString("caas_syncer.login.url")
	payload := make(map[string]string)
	payload["userName"] = viper.GetString("caas_syncer.login.user")
	payload["password"] = viper.GetString("caas_syncer.login.password")
	buf, err := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, loginUrl, bytes.NewReader(buf))
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Connection", "Keep-Alive")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return lr, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("[CaasSyncer] GetLoginResult: %s", body)
	err = json.Unmarshal(body, lr)
	log.Debugf("[CaasSyncer] GetLoginObj: %+v", lr)
	return lr, err
}

type CaasWorkSpaceResult struct {
	Code int64            `json:"code"`
	Data []caas.WorkSpace `json:"data"`
	Msg  string           `json:"msg"`
}

type App struct {
	ID          int64     `json:"Id"`
	AppName     string    `json:"AppName"`
	NameSpaceId int64     `json:"NamespaceID"`
	Description string    `json:"Description"`
	CreateTime  time.Time `json:"CreateTime"`
}

type CaasAppResult struct {
	Code int64  `json:"code"`
	Data []App  `json:"data"`
	Msg  string `json:"msg"`
}

// UpdateApps 更新应用列表
func UpdateApps(ns *caas.NameSpace) {
	appResult, err := GetApp(ns)
	if err != nil {
		log.Error(err)
		return
	}
	for _, app := range appResult.Data {
		UpdateApp(&app)
	}
}

func UpdateApp(app *App) {
	napp := caas.App{
		ID:          app.ID,
		AppName:     app.AppName,
		NameSpaceId: app.NameSpaceId,
		Description: app.Description,
		CreateTime:  gtime.NewGTime(app.CreateTime),
		UpdateTime:  gtime.Now(),
	}

	db := g.Con().Portal.Debug()
	tapp := caas.App{}
	db.Model(napp).Debug().Where(caas.App{ID: napp.ID}).First(&tapp)
	if tapp.ID == 0 {
		db.Model(napp).Create(&napp)
	} else {
		db.Model(napp).Updates(&napp)
	}
}

// GetApp 获取单页应用列表
func GetApp(ns *caas.NameSpace) (*CaasAppResult, error) {
	namespaceID := ns.ID
	clusterID := ns.ClusterId
	lr, err := Login()
	if err != nil || lr.Data.Token == "" {
		log.Errorln(err)
		return nil, err
	}
	token := lr.Data.Token

	log.Debugf("[CaasSyncer] GetApp...")
	app := &CaasAppResult{}
	appUrl := viper.GetString("caas_syncer.app.url")
	appUrl = fmt.Sprintf(appUrl, namespaceID)
	req, _ := http.NewRequest(http.MethodGet, appUrl, nil)
	query := req.URL.Query()
	//query.Add("current", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Authorization", token)
	req.Header.Add("clusterId", strconv.FormatInt(clusterID, 10))
	log.Debugf("request url: %s, params: %+v, headers: %+v", req.URL, query, req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return app, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("[CaasSyncer] GetAppResult: %s", body)
	err = json.Unmarshal(body, app)
	log.Debugf("[CaasSyncer] GetAppObj: %+v", app)

	return app, err
}

// GetWorkSpace 获取组织空间列表
func GetWorkSpace() (*CaasWorkSpaceResult, error) {
	lr, err := Login()
	if err != nil || lr.Data.Token == "" {
		log.Errorln(err)
		return nil, err
	}
	token := lr.Data.Token

	log.Debugf("[CaasSyncer] GetWorkSpace...")
	ws := &CaasWorkSpaceResult{}
	workspaceUrl := viper.GetString("caas_syncer.workspace.url")
	req, _ := http.NewRequest(http.MethodGet, workspaceUrl, nil)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return ws, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("[CaasSyncer] GetWorkSpaceResult: %s", body)
	err = json.Unmarshal(body, ws)
	log.Debugf("[CaasSyncer] GetWorkSpaceObj: %+v", ws)

	return ws, err
}

type CaasNameSpaceResult struct {
	Code int64            `json:"code"`
	Data []caas.NameSpace `json:"data"`
	Msg  string           `json:"msg"`
}

// GetNameSpace 获取组织空间列表下项目空间列表
func GetNameSpace(workspaceId int64) (*CaasNameSpaceResult, error) {
	lr, err := Login()
	if err != nil || lr.Data.Token == "" {
		log.Errorln(err)
		return nil, err
	}
	token := lr.Data.Token

	log.Debugf("[CaasSyncer] GetNameSpace...")
	ns := &CaasNameSpaceResult{}
	namespaceUrl := viper.GetString("caas_syncer.namespace.url")
	req, _ := http.NewRequest(http.MethodGet, namespaceUrl, nil)
	query := req.URL.Query()
	query.Add("workspace_id", strconv.FormatInt(workspaceId, 10))
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return ns, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("[CaasSyncer] GetNameSpaceResult: %s", body)
	err = json.Unmarshal(body, &ns)
	log.Debugf("[CaasSyncer] GetNameSpaceObj: %+v", ns)

	return ns, err
}

type CaasServiceResult struct {
	Code int64          `json:"code"`
	Data []caas.Service `json:"data"`
	Msg  string         `json:"msg"`
}

// GetService 获取项目空间下的服务列表
func GetService(ns caas.NameSpace) (*CaasServiceResult, error) {
	namespaceID := ns.ID
	clusterId := ns.ClusterId
	lr, err := Login()
	if err != nil || lr.Data.Token == "" {
		log.Errorln(err)
		return nil, err
	}
	token := lr.Data.Token

	log.Debugf("[CaasSyncer] GetService...")
	sr := &CaasServiceResult{}
	serviceUrl := viper.GetString("caas_syncer.service.url")
	serviceUrl = fmt.Sprintf(serviceUrl, namespaceID)
	req, _ := http.NewRequest(http.MethodGet, serviceUrl, nil)
	query := req.URL.Query()
	query.Add("type", "deployment")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Authorization", token)
	req.Header.Add("clusterId", strconv.FormatInt(clusterId, 10))
	log.Debugf("request url: %s, params: %+v, headers: %+v", req.URL, query, req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return sr, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("[CaasSyncer] GetServiceResult: %s", body)
	err = json.Unmarshal(body, &sr)
	log.Debugf("[CaasSyncer] GetServiceObj: %+v", sr)

	return sr, err
}

type CaasPodResult struct {
	Code int64      `json:"code"`
	Data []caas.Pod `json:"data"`
	Msg  string     `json:"msg"`
}

func UpdateWorkspace(ws *caas.WorkSpace) {
	ws.UpdateTime = gtime.NewGTime(time.Now())

	db := g.Con().Portal.Debug()
	t := caas.WorkSpace{}
	db.Model(ws).Where(caas.WorkSpace{ID: ws.ID}).First(&t)
	if t.ID == 0 {
		db.Model(ws).Create(&ws)
	} else {
		db.Model(ws).Where(caas.WorkSpace{ID: ws.ID}).Updates(&ws)
	}
}

func UpdateService(ser *caas.Service) {
	ser.UpdateTime = gtime.Now()

	db := g.Con().Portal.Debug()
	tser := caas.Service{}
	db.Model(ser).Debug().Where(caas.Service{Type: ser.Type, ServiceName: ser.ServiceName, AppId: ser.AppId}).First(&tser)
	if tser.ID == 0 {
		db.Model(ser).Create(&ser)
	} else {
		ser.ID = tser.ID
		db.Model(ser).Updates(&ser)
	}
}

func UpdateNamespace(ns *caas.NameSpace) {
	ns.UpdateTime = gtime.NewGTime(time.Now())

	db := g.Con().Portal.Debug()
	n := caas.NameSpace{}
	db.Model(ns).Where(caas.NameSpace{ID: ns.ID}).First(&n)
	if n.ID == 0 {
		db.Model(ns).Create(&ns)
	} else {
		db.Model(ns).Where(caas.NameSpace{ID: ns.ID}).Updates(&ns)
	}
}

func UpdateServicePorts(ser *caas.Service) {
	db := g.Con().Portal.Debug()

	var ports []int64
	for _, p := range ser.Ports {
		p.UpdateTime = gtime.NewGTime(time.Now())

		port := caas.Port{}
		db.Model(port).Where(caas.Port{Host: p.Host}).First(&port)
		if port.ID == 0 {
			db.Model(port).Create(&p)
		} else {
			p.ID = port.ID
			db.Model(port).Updates(&p)
		}
		ports = append(ports, p.ID)
	}
	// service port rel
	db.Model(caas.ServicePortRel{}).Debug().Where(&caas.ServicePortRel{Service: ser.ID}).Delete(&caas.ServicePortRel{})
	for _, p := range RemoveRepeated(ports) {
		db.Model(caas.ServicePortRel{}).Debug().Create(&caas.ServicePortRel{Service: ser.ID, Port: p})
	}
}

func UpdateNamespaceServiceRel(ns *caas.NameSpace, ser *caas.Service) {
	db := g.Con().Portal.Debug()
	rel := caas.NamespaceServiceRel{
		NameSpace: ns.ID,
		Service:   ser.ID,
	}
	if !rel.Existing() {
		db.Model(rel).Create(&rel)
	}
}

func UpdatePods(ser *caas.Service, pods *CaasPodResult) {
	db := g.Con().Portal.Debug()

	// 更新数据库
	var podIds []int64
	for _, p := range pods.Data {
		p.UpdateTime = gtime.NewGTime(time.Now())

		pod := caas.Pod{}
		db.Model(pod).Where(caas.Pod{Name: p.Name}).First(&pod)
		if pod.ID == 0 {
			db.Model(pod).Create(&p)
		} else {
			p.ID = pod.ID
			db.Model(pod).Updates(&p)
		}

		podIds = append(podIds, p.ID)
	}

	// service pod rel
	var rels []caas.ServicePodRel
	db.Model(caas.ServicePodRel{}).Where(&caas.ServicePodRel{Service: ser.ID}).Delete(&rels)
	for _, p := range RemoveRepeated(podIds) {
		db.Model(caas.ServicePodRel{}).Create(&caas.ServicePodRel{Service: ser.ID, Pod: p})
	}
}

// GetPod 获取服务下的实例
func GetPod(namespace *caas.NameSpace, service *caas.Service) (*CaasPodResult, error) {
	lr, err := Login()
	if err != nil || lr.Data.Token == "" {
		log.Errorln(err)
		return nil, err
	}
	token := lr.Data.Token

	log.Debugf("[CaasSyncer] GetPod...")
	inst := &CaasPodResult{}
	podUrl := viper.GetString("caas_syncer.pod.url")
	podUrl = fmt.Sprintf(podUrl, namespace.ID)
	req, _ := http.NewRequest(http.MethodGet, podUrl, nil)
	query := req.URL.Query()
	query.Add("appName", service.ServiceName)
	query.Add("type", service.Type)
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Authorization", token)
	req.Header.Set("clusterId", strconv.FormatInt(namespace.ClusterId, 10))
	log.Debugf("request url: %s, params: %+v, headers: %+v", req.URL, query, req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return inst, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("[CaasSyncer] GetPodResult: %s", body)
	err = json.Unmarshal(body, inst)
	log.Debugf("[CaasSyncer] GetPodObj: %+v", inst)

	return inst, err
}

func InitCaasSyncer() *CaasSyncer {
	syncer := &CaasSyncer{}
	syncer.ctx, syncer.cancel = context.WithCancel(context.Background())
	return syncer
}
