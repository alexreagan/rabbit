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


type NextNodeInfoInputs struct {
	ProcessInstID string `json:"PROCESS_INST_ID"`
	TemplateID string `json:"TEMPLATE_ID"`
	TaskTD string `json:"TASK_ID"`
	Conditions []*Condition `json:"CONDITIONS"`
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
