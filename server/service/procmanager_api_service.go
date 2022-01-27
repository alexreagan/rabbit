package service

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
)

const (
	ProcManagerApiMethodGet               = "GET"
	ProcManagerApiMethodPost              = "POST"
	ProcManagerApiMethodPut               = "PUT"
	ProcManagerApiMethodDelete            = "DELETE"
	ProcManagerApiContentType             = "application/x-www-form-urlencoded"
	ProcManagerApiUserAgent               = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	ProcManagerApiConnection              = "keep-alive"
	ProcManagerApiUpgradeInSecureRequests = "1"
	ProcManagerApiLoginType               = "1"
	ProcManagerApiSessionIdName           = "uliweb_session_id"
)

type procManagerApiService struct {
	Addr string
	mu   sync.Mutex
}

type ProcManagerApiLoginInputs struct {
	TYPE     string
	NEXT     string
	USERNAME string
	PASSWORD string
}

type ProcManagerApiSession struct {
	UliWebSessionID string
	BodyContentLen  int
}

type ProcManagerApiNextUser struct {
	ID              string
	NAME            string
	PRC_ACTION_ID   string
	IDPRC_ACTION_ID string
	USR_ID_LAND_NM  string
	CUR_USR_INST_ID string
	CUR_USR_INST_NM string
}

type ProcManagerApiCondition struct {
	KEY   string
	VALUE string
}

type ProcManagerApiCreateInputs struct {
	TEMPLATE_ID     string
	TASK_ID         string
	REMARK          string
	USER_ID         string
	CUR_USR_INST_ID string
	CUR_USR_INST_NM string
	PRJ_ID          string
	PRJ_SN          string
	TO_DO_TM_TTL    string
	BUTTON_NAME     string
	USER_NAME       string
	USR_ID_LAND_NM  string
	NEXT_USER_GRP   []ProcManagerApiNextUser
	CONDITIONS      []ProcManagerApiCondition
}

type ProcManagerApiCreateProcess struct {
	PROCESS_INST_ID  string
	SYS_EVT_TRACE_ID string
	SYS_RESP_DESC    string
	SYS_RECV_TIME    string
	SYS_RESP_CODE    string
	SYS_RESP_TIME    string
}

type ProcManagerApiTodoInputs struct {
	PROCESS_INST_ID string
}

type ProcManagerToDoItem struct {
	AVY_OWR_ID      string
	PRJ_ID          string
	TODO_SN         string
	PROCESS_INST_ID string
	PRJ_NM          string
	PCS_AVY_NM      string
	TASK_ID         string
	TODO_START_TM   string
	TEMPLATE_ID     string
	AVY_OWR_NM      string
}

type ProcManagerApiTodoRecords struct {
	SYS_RESP_DESC    string
	SYS_RECV_TIME    string
	SYS_RESP_CODE    string
	SYS_EVT_TRACE_ID string
	SYS_RESP_TIME    string
	TODO_INFO        []ProcManagerToDoItem
}

type ProcManageApiProcExecuteInputs struct {
	PROCESS_INST_ID string
	TEMPLATE_ID     string
	TASK_ID         string
	REMARK          string
	OPIN_DESC       string
	PRJ_ID          string
	PRJ_SN          string
	TO_DO_TM_TTL    string
	TODO_ID         string
	NEXT_USER_GRP   []ProcManagerApiNextUser
	BUTTON_NAME     string
	CONDITIONS      []ProcManagerApiCondition
}

type ProcManagerApiProcExecuteRecord struct {
	SYS_RESP_DESC    string
	SYS_TX_TYPE      string
	SYS_RECV_TIME    string
	SYS_TX_STATUS    string
	SYS_EVT_TRACE_ID string
	SYS_RESP_CODE    string
	RESULT_DESC      string
	SYS_RESP_TIME    string
}

// 流程处理: 根据编号和待办标识
func (s *procManagerApiService) procManagerApiExecute(param ProcManageApiProcExecuteInputs) (*ProcManagerApiProcExecuteRecord, error) {
	client := &http.Client{}
	targetUrl := fmt.Sprintf("%s/procmanager/api/procExecute", s.Addr)

	var ItemExecuteVal = ProcManageApiProcExecuteInputs{
		PROCESS_INST_ID: string(param.PROCESS_INST_ID),
		TEMPLATE_ID:     string(param.TEMPLATE_ID),
		TASK_ID:         string(param.TASK_ID),
		REMARK:          string(param.REMARK),
		OPIN_DESC:       string(param.OPIN_DESC),
		PRJ_ID:          string(param.PRJ_ID),
		PRJ_SN:          string(param.PRJ_SN),
		TO_DO_TM_TTL:    string(param.TO_DO_TM_TTL),
		TODO_ID:         string(param.TODO_ID),
		NEXT_USER_GRP:   param.NEXT_USER_GRP,
		BUTTON_NAME:     string(param.BUTTON_NAME),
		CONDITIONS:      param.CONDITIONS,
	}
	postBytes, _ := json.Marshal(ItemExecuteVal)
	req, _ := http.NewRequest(ProcManagerApiMethodPost, targetUrl, strings.NewReader("jsonData="+string(postBytes)))

	var uliWebSessionId string
	var pLogin ProcManagerApiLoginInputs
	uliWebSessionMap, err := s.procManagerApiLogin(pLogin)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	uliWebSessionId = string(uliWebSessionMap.UliWebSessionID)

	req.Header.Set("Content-Type", ProcManagerApiContentType)
	req.Header.Add("Cookie", ProcManagerApiSessionIdName+"="+uliWebSessionId)

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var executeRecord ProcManagerApiProcExecuteRecord
	err = json.Unmarshal(body, &executeRecord)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &executeRecord, err
}

// 获取流程编号对应的待办数据
func (s *procManagerApiService) procManagerApiTodoInfo(param ProcManagerApiTodoInputs) (*ProcManagerApiTodoRecords, error) {
	client := &http.Client{}
	targetUrl := fmt.Sprintf("%s/procmanager/api/procInstTodoInfo", s.Addr)

	var ItemToDoVal = ProcManagerApiTodoInputs{
		PROCESS_INST_ID: string(param.PROCESS_INST_ID),
	}
	postBytes, _ := json.Marshal(ItemToDoVal)
	req, _ := http.NewRequest(ProcManagerApiMethodPost, targetUrl, strings.NewReader("jsonData="+string(postBytes)))

	var uliWebSessionId string
	var pLogin ProcManagerApiLoginInputs
	uliWebSessionMap, err := s.procManagerApiLogin(pLogin)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	uliWebSessionId = string(uliWebSessionMap.UliWebSessionID)

	req.Header.Set("Content-Type", ProcManagerApiContentType)
	req.Header.Add("Cookie", ProcManagerApiSessionIdName+"="+uliWebSessionId)

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var toDoRecords ProcManagerApiTodoRecords
	err = json.Unmarshal(body, &toDoRecords)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &toDoRecords, err
}

// 流程发起
func (s *procManagerApiService) procManagerApiProcCreate(param ProcManagerApiCreateInputs) (*ProcManagerApiCreateProcess, error) {
	client := &http.Client{}
	targetUrl := fmt.Sprintf("%s/procmanager/api/procCreate", s.Addr)

	nextUserGrp := make([]ProcManagerApiNextUser, 0)
	for _, nextUserVal := range param.NEXT_USER_GRP {
		nextUserGrp = append(nextUserGrp, nextUserVal)
	}

	conditions := make([]ProcManagerApiCondition, 0)
	for _, condition := range param.CONDITIONS {
		conditions = append(conditions, condition)
	}

	var ItemCreateVal = ProcManagerApiCreateInputs{
		TEMPLATE_ID:     string(param.TEMPLATE_ID), // 600100PubAudit
		TASK_ID:         string(param.TASK_ID),     // 10101
		REMARK:          string(param.REMARK),
		USER_ID:         string(param.USER_ID),         // 23598915
		CUR_USR_INST_ID: string(param.CUR_USR_INST_ID), // 010200103
		CUR_USR_INST_NM: string(param.CUR_USR_INST_NM),
		PRJ_ID:          string(param.PRJ_ID), // 1001
		PRJ_SN:          string(param.PRJ_SN), // 1001
		TO_DO_TM_TTL:    string(param.TO_DO_TM_TTL),
		BUTTON_NAME:     string(param.BUTTON_NAME),
		USER_NAME:       string(param.USER_NAME),
		USR_ID_LAND_NM:  string(param.USR_ID_LAND_NM),
		NEXT_USER_GRP:   nextUserGrp,
		CONDITIONS:      conditions,
	}

	postBytes, _ := json.Marshal(ItemCreateVal)

	req, _ := http.NewRequest(ProcManagerApiMethodPost, targetUrl, strings.NewReader("jsonData="+string(postBytes)))

	var uliWebSessionId string
	var pLogin ProcManagerApiLoginInputs
	uliWebSessionMap, err := s.procManagerApiLogin(pLogin)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	uliWebSessionId = string(uliWebSessionMap.UliWebSessionID)

	req.Header.Set("Content-Type", ProcManagerApiContentType)
	req.Header.Add("Cookie", ProcManagerApiSessionIdName+"="+uliWebSessionId)

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var createProcessMap ProcManagerApiCreateProcess
	err = json.Unmarshal(body, &createProcessMap)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &createProcessMap, err
}

// 登录：获得Cookie
func (s *procManagerApiService) procManagerApiLogin(param ProcManagerApiLoginInputs) (*ProcManagerApiSession, error) {
	// http://{IP}:{port}/login
	u1, err := url.Parse(s.Addr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// u1.Scheme, u1.Opaque, u1.User, u1.Host, u1.Path, u1.RawQuery, u1.Frament
	host := u1.Host
	targetUrl := fmt.Sprintf("%s/login", s.Addr)
	nextUrl := fmt.Sprintf("%s/", s.Addr)

	var userName, passWord string
	if len(param.USERNAME) > 1 {
		userName = param.USERNAME
	} else {
		userName = viper.GetString("procManager.userName")
	}

	if len(param.PASSWORD) > 1 {
		passWord = param.PASSWORD
	} else {
		passWord = viper.GetString("procManager.passWord")
	}

	// 先登录{s.Addr}/procmanager/api项目, 才能使用登录成功跳转时携带的cookie继续调用流程引擎其它接口
	var uliWebSessionId string
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	client := &http.Client{Jar: jar}

	var LoginMapVal = ProcManagerApiLoginInputs{
		TYPE:     ProcManagerApiLoginType,
		NEXT:     string(nextUrl),
		USERNAME: string(userName),
		PASSWORD: string(passWord),
	}

	postString := "type=" + LoginMapVal.TYPE + "&next=" + LoginMapVal.NEXT + "&username=" + LoginMapVal.USERNAME + "&password=" + LoginMapVal.PASSWORD
	req, _ := http.NewRequest(ProcManagerApiMethodPost, targetUrl, strings.NewReader(string(postString)))

	req.Header.Set("Content-Type", ProcManagerApiContentType)
	req.Header.Add("User-Agent", ProcManagerApiUserAgent)
	req.Header.Add("Connection", ProcManagerApiConnection)
	req.Header.Add("Upgrade-Insecure-Requests", ProcManagerApiUpgradeInSecureRequests)
	req.Header.Add("Host", host)
	req.Header.Add("Referer", nextUrl)

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	var respCookies []*http.Cookie
	respCookies = resp.Cookies()

	for _, respV := range respCookies {
		if respV.Name == ProcManagerApiSessionIdName {
			uliWebSessionId = respV.Value
			break
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	bodyContentLen := len(string(body))
	return &ProcManagerApiSession{
		UliWebSessionID: uliWebSessionId,
		BodyContentLen:  bodyContentLen,
	}, nil
}

func newProcManagerApiService() *procManagerApiService {
	addr := viper.GetString("procManager.addr")
	return &procManagerApiService{
		Addr: addr,
	}
}
