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
	mu   sync.Mutex
}

type Condition struct {
	Key   string `json:"key" form:"key,omitempty"`
	Value string `json:"value" form:"value,omitempty"`
}

type NextNodeInfoInputs struct {
	ProcessInstID string       `json:"PROCESS_INST_ID"`
	TemplateID    string       `json:"TEMPLATE_ID"`
	TaskTD        string       `json:"TASK_ID"`
	Conditions    []*Condition `json:"CONDITIONS"`
}

func (s *procService) NextNodeInfo(param NextNodeInfoInputs) (string, error) {

	uri := fmt.Sprintf("%s/procmanager/api/procNextNodeInfo", viper.GetString("proc_manager.addr"))
	req := httplib.Post(uri)
	byts, _ := json.Marshal(param)
	req.Param("jsonData", string(byts))
	return req.String()
}

func (s *procService) GetPersonByNode(templateID string, taskID string) (interface{}, interface{}) {
	uri := fmt.Sprintf("%s/procmanager/api/getPersonByNode", viper.GetString("proc_manager.addr"))
	req := httplib.Post(uri)
	param := make(map[string]string, 0)
	param["TEMPLATE_ID"] = templateID
	param["TASK_ID"] = taskID
	byts, _ := json.Marshal(param)
	req.Param("jsonData", string(byts))
	return req.String()
}

func (s *procService) GetHistDetailList(processInstID string, taskID string, belongInstID string, selectMode string) (interface{}, interface{}) {
	uri := fmt.Sprintf("%s/procmanager/api/procGetHistDetailList", s.Addr)
	req := httplib.Post(uri)
	param := make(map[string]string, 0)
	param["PROCESS_INST_ID"] = processInstID
	param["TASK_ID"] = taskID
	param["BLNG_INST_ID"] = belongInstID
	param["SELECT_MODE"] = selectMode
	byts, _ := json.Marshal(param)
	req.Param("jsonData", string(byts))
	return req.String()
}

func newProcService() *procService {
	return &procService{}
}
