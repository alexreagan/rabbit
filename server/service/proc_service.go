package service

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/toolkits/net/httplib"
	"sync"
)

type procService struct {
	Addr string
	mu sync.Mutex
}

type Condition struct {
	Key string `json:"KEY"`
	Value string `json:"VALUE"`
}


type NextUser struct {
	ID string `json:"ID"`
	Name string `json:"name"`
	IDPrcActionID string `json:"IDPrcActionID"`
	UsrIDLandNm string `json:"UsrIDLandNm"`
	CurUserInstID string `json:"curUserInstID"`
	CurUserInstNm string `json:"curUserInstNm"`
}

type NextNodeInfoInputs struct {
	ProcessInstID string `json:"PROCESS_INST_ID"`
	TemplateID string `json:"TEMPLATE_ID"`
	TaskTD string `json:"TASK_ID"`
	Conditions []*Condition `json:"CONDITIONS"`
}

type ProcCreateInputs struct {
	TemplateID    string `json: "TEMPLATE_ID"`
	TaskID        string `json: "TASK_ID"`
	Remark        string `json: "REMARK"`
	OpinDesc      string `json: "OPIN_DESC"`
	NextUserFlag  string `json: "NEXT_USER_FLAG"`
	ButtonName    string `json: "BUTTON_NAME"`
	UserID        string `json: "USER_ID"`
	UserName      string `json: "USER_NAME"`
	UsrIDLandNm  string `json: "USR_ID_LAND_NM"`
	CurUsrInstID string `json: "CUR_USR_INST_ID"`
	CurUsrInstNm string `json: "CUR_USR_INST_NM"`
	NextUserGrp   []*NextUser `json: "NEXT_USER_GRP"`
	Conditions    []*Condition `json: "CONDITIONS"`
	PrjID         string `json: "PRJ_ID"`
	PrjSn         string `json: "PRJ_SN"`
	ToDoTmTpCd    string `json: "TO_DO_TM_TP_CD"`
	ToDoTmTtl     string `json: "TO_DO_TM_TTL"`
	BlngInstID    string `json: "BLNG_INST_ID"`
	DmnGrpID      string `json: "DMN_GRP_ID"`
}

type ProcExecuteInputs struct {
	ProcessInstID string `json: "PROCESS_INST_ID"`
	TemplateID    string `json: "TEMPLATE_ID"`
	TaskID        string `json: "TASK_ID"`
	Remark        string `json: "REMARK"`
	OpinCode      string `json: "OPIN_CODE"`
	OpinDesc      string `json: "OPIN_DESC"`
	NextUserFlag  string `json: "NEXT_USER_FLAG"`
	ButtonName    string `json: "BUTTON_NAME"`
	UserID        string `json: "USER_ID"`
	UserName      string `json: "USER_NAME"`
	UsrIDLandNm  string `json: "USR_ID_LAND_NM"`
	CurUsrInstID string `json: "CUR_USR_INST_ID"`
	CurUsrInstNm string `json: "CUR_USR_INST_NM"`
	NextUserGrp   []*NextUser `json: "NEXT_USER_GRP"`
	Conditions    []*Condition `json: "CONDITIONS"`
}

func (s *procService) ProcCreate(param ProcCreateInputs) (string, error) {
	uri := fmt.Sprintf("%s/procmanager/api/procCreate", s.Addr)
	req := httplib.Post(uri)
	jsonData, _ := json.Marshal(param)
	req.Param("jsonData", string(jsonData))
	return req.String()
}

func (s *procService) ProcExecute(param ProcExecuteInputs) (string, error) {
	uri := fmt.Sprintf("%s/procmanager/api/procExecute", s.Addr)
	req := httplib.Post(uri)
	jsonData, _ := json.Marshal(param)
	req.Param("jsonData", string(jsonData))
	return req.String()
}

func (s *procService) NextNodeInfo(param NextNodeInfoInputs) (string, error) {

	uri := fmt.Sprintf("%s/procmanager/api/procNextNodeInfo", s.Addr)
	req := httplib.Post(uri)
	byts, _ := json.Marshal(param)
	req.Param("jsonData", string(byts))
	return req.String()
}

func (s *procService) GetPersonByNode(templateID string, taskID string) (interface{}, interface{}) {
	uri := fmt.Sprintf("%s/procmanager/api/getPersonByNode", s.Addr)
	req := httplib.Post(uri)
	param := make(map[string]string, 0)
	param["TEMPLATE_ID"] = templateID
	param["TASK_ID"] = taskID
	byts, _ := json.Marshal(param)
	req.Param("jsonData", string(byts))
	return req.String()
}

func (s *procService) GetHistDetailList(processInstID string, taskID string) (interface{}, interface{}) {
	uri := fmt.Sprintf("%s/procmanager/api/procGetHistDetailList", s.Addr)
	req := httplib.Post(uri)
	param := make(map[string]string, 0)
	param["PROCESS_INST_ID"] = processInstID
	param["TASK_ID"] = taskID
	byts, _ := json.Marshal(param)
	req.Param("jsonData", string(byts))
	return req.String()
}

func newProcService() *procService {
	addr := viper.GetString("procManager.addr")
	return &procService{
		Addr: addr,
	}
}
