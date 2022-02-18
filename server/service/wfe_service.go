package service

import (
	"encoding/xml"
	"fmt"
	"github.com/alexreagan/rabbit/server/model/uic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/toolkits/net/httplib"
	"os"
	"path/filepath"
	"time"
)

type wfeService struct {
	serviceConfig map[string]string
	endPointConf  map[string]string
	dataConf      map[string]string
	segmentConf   map[string]string
}

func (s *wfeService) dataTransformWalkFunc(path string, f os.FileInfo, err error) error {
	log.Println(path)
	if f == nil {
		return err
	}
	log.Println(f.Name())
	return nil
}

func (s *wfeService) outBoundWalkFunc(path string, f os.FileInfo, err error) error {
	log.Println(path)
	if f == nil {
		return err
	}
	log.Println(f.Name())
	return nil
}

func (s *wfeService) init() {
	xmlDirDataTransform := viper.GetString("wfe.config.data_transform")
	filepath.Walk(xmlDirDataTransform, s.dataTransformWalkFunc)

	xmlDirOutBound := viper.GetString("wfe.config.outbound")
	filepath.Walk(xmlDirOutBound, s.outBoundWalkFunc)
}

type ExecutorInfo struct {
	ID           string `json:"ID" xml:"ID"`
	Name         string `json:"NAME" xml:"NAME"`
	AvyOwr       string `json:"AVY_OWR" xml:"AVY_OWR,omitempty"`
	UsrIDLandNm  string `json:"USR_ID_LAND_NM" xml:"USR_ID_LAND_NM,omitempty"`
	CurUsrInstID string `json:"CUR_USR_INST_ID" xml:"CUR_USR_INST_ID,omitempty"`
	CurUsrInstNm string `json:"CUR_USR_INST_NM" xml:"CUR_USR_INST_NM,omitempty"`
	BlngInstID   string `json:"BLNG_INST_ID" xml:"BLNG_INST_ID,omitempty"`
	BlngInstNm   string `json:"BLNG_INST_NM" xml:"BLNG_INST_NM,omitempty"`
}

type NextUserGrp struct {
	Type         string `json:"type" xml:"type,attr"`
	ID           string `json:"ID" xml:"ID"`
	Name         string `json:"NAME" xml:"NAME"`
	PrcActionID  string `json:"PRC_ACTION_ID" xml:"PRC_ACTION_ID"`
	UsrIDLandNm  string `json:"USR_ID_LAND_NM" xml:"USR_ID_LAND_NM"`
	CurUsrInstID string `json:"CUR_USR_INST_ID" xml:"CUR_USR_INST_ID"`
	CurUsrInstNm string `json:"CUR_USR_INST_NM" xml:"CUR_USR_INST_NM"`
}

type CONDITION struct {
	Type  string `json:"type" xml:"type,attr"`
	Key   string `json:"KEY" xml:"KEY"`
	Value string `json:"VALUE" xml:"VALUE"`
}

type BsnComInfo struct {
	PrjID      string `json:"PRJ_ID" xml:"PRJ_ID"`
	PrjSN      string `json:"PRJ_SN" xml:"PRJ_SN"`
	TodoTmTpCd string `json:"TO_DO_TM_TP_CD" xml:"TO_DO_TM_TP_CD,omitempty"`
	TodoTmTtl  string `json:"TO_DO_TM_TTL" xml:"TO_DO_TM_TTL"`
	BlngInstID string `json:"BLNG_INST_ID" xml:"BLNG_INST_ID,omitempty"`
	DmnGrpID   string `json:"DMN_GRP_ID" xml:"DMN_GRP_ID,omitempty"`
}

type TxBodyEntityAppEntity struct {
	ProcessInstID string         `json:"PROCESS_INST_ID" xml:"PROCESS_INST_ID"`
	TemplateID    string         `json:"TEMPLATE_ID" xml:"TEMPLATE_ID"`
	TaskID        string         `json:"TASK_ID" xml:"TASK_ID"`
	Remark        string         `json:"REMARK" xml:"REMARK"`
	OpinCode      string         `json:"OPIN_CODE" xml:"OPIN_CODE,omitempty"`
	OpinDesc      string         `json:"OPIN_DESC" xml:"OPIN_DESC,omitempty"`
	NextUserFlag  string         `json:"NEXT_USER_FLAG" xml:"NEXT_USER_FLAG,omitempty"`
	ExeFstTsk     string         `json:"EXE_FST_TSK" xml:"EXE_FST_TSK,omitempty"`
	ButtonName    string         `json:"BUTTON_NAME" xml:"BUTTON_NAME"`
	ExecutorInfo  *ExecutorInfo  `json:"EXECUTER_INFO" xml:"EXECUTER_INFO"`
	NextUserGrp   []*NextUserGrp `json:"NEXT_USER_GRP" xml:"NEXT_USER_GRP,omitempty"`
	Conditions    []*CONDITION   `json:"CONDITION" xml:"CONDITION,omitempty"`
	BsnComInfo    *BsnComInfo    `json:"BSN_COM_INFO" xml:"BSN_COM_INFO,omitempty"`
}

type TXHeader struct {
	SysHdrLen        int64  `json:"SYS_HDR_LEN" xml:"SYS_HDR_LEN"`
	SysPkgVrsn       string `json:"SYS_PKG_VRSN" xml:"SYS_PKG_VRSN"`
	SysTtlLen        int64  `json:"SYS_TTL_LEN" xml:"SYS_TTL_LEN"`
	SysReqSecID      string `json:"SYS_REQ_SEC_ID" xml:"SYS_REQ_SEC_ID"`
	SysSndSecID      string `json:"SYS_SND_SEC_ID" xml:"SYS_SND_SEC_ID"`
	SysTxCode        string `json:"SYS_TX_CODE" xml:"SYS_TX_CODE,omitempty"`
	SysTxVrsn        string `json:"SYS_TX_VRSN" xml:"SYS_TX_VRSN,omitempty"`
	SysTxType        string `json:"SYS_TX_TYPE" xml:"SYS_TX_TYPE"`
	SysReserved      string `json:"SYS_RESERVED" xml:"SYS_RESERVED,omitempty"`
	SysEvtTraceID    string `json:"SYS_EVT_TRACE_ID" xml:"SYS_EVT_TRACE_ID"`
	SysSndSerialNo   string `json:"SYS_SND_SERIAL_NO" xml:"SYS_SND_SERIAL_NO"`
	SysPkgType       string `json:"SYS_PKG_TYPE" xml:"SYS_PKG_TYPE"`
	SysMsgLen        int64  `json:"SYS_MSG_LEN" xml:"SYS_MSG_LEN"`
	SysIsEncrypted   string `json:"SYS_IS_ENCRYPTED" xml:"SYS_IS_ENCRYPTED"`
	SysEncryptType   string `json:"SYS_ENCRYPT_TYPE" xml:"SYS_ENCRYPT_TYPE"`
	SysCompressType  string `json:"SYS_COMPRESS_TYPE" xml:"SYS_COMPRESS_TYPE"`
	SysEmbMsgLen     int64  `json:"SYS_EMB_MSG_LEN" xml:"SYS_EMB_MSG_LEN"`
	SysReqTime       string `json:"SYS_REQ_TIME" xml:"SYS_REQ_TIME"`
	SysTimeLeft      string `json:"SYS_TIME_LEFT" xml:"SYS_TIME_LEFT,omitempty"`
	SysRecvTime      string `json:"SYS_RECV_TIME" xml:"SYS_RECV_TIME,omitempty"`
	SysRespTime      string `json:"SYS_RESP_TIME" xml:"SYS_RESP_TIME,omitempty"`
	SysPkgStsType    string `json:"SYS_PKG_STS_TYPE" xml:"SYS_PKG_STS_TYPE"`
	SysTxStatus      string `json:"SYS_TX_STATUS" xml:"SYS_TX_STATUS,omitempty"`
	SysRespCode      string `json:"SYS_RESP_CODE" xml:"SYS_RESP_CODE,omitempty"`
	SysRespDescLen   string `json:"SYS_RESP_DESC_LEN" xml:"SYS_RESP_DESC_LEN,omitempty"`
	SysRespDesc      string `json:"SYS_RESP_DESC" xml:"SYS_RESP_DESC,omitempty"`
	SysSecContextLen int64  `json:"SYS_SEC_CONTEXT_LEN" xml:"SYS_SEC_CONTEXT_LEN,omitempty"`
	SysSecContext    string `json:"SYS_SEC_CONTEXT" xml:"SYS_SEC_CONTEXT,omitempty"`
}

type TXBodyCommonComB struct {
	ErrMsgNum     string `json:"ERR_MSG_NUM" xml:"ERR_MSG_NUM"`
	CmptTrcNo     string `json:"CMPT_TRCNO" xml:"CMPT_TRCNO"`
	TotalPage     int64  `json:"TOTAL_PAGE" xml:"TOTAL_PAGE"`
	TotalRec      int64  `json:"TOTAL_REC" xml:"TOTAL_REC"`
	CurrTotalPage int64  `json:"CURR_TOTAL_PAGE" xml:"CURR_TOTAL_PAGE"`
	CurrTotalRec  int64  `json:"CURR_TOTAL_REC" xml:"CURR_TOTAL_REC"`
	StsTraceID    string `json:"STS_TRACE_ID" xml:"STS_TRACE_ID"`
}

type TXBodyCommonCom1 struct {
	TxnInsID          string `json:"TXN_INSID" xml:"TXN_INSID"`
	TxnIttChnlID      string `json:"TXN_ITT_CHNL_ID" xml:"TXN_ITT_CHNL_ID"`
	TxnIttChnlCgyCode string `json:"TXN_ITT_CHNL_CGY_CODE" xml:"TXN_ITT_CHNL_CGY_CODE"`
	TxnDT             string `json:"TxnDT" xml:"TXN_DT"`
	TxnTM             string `json:"TxnTM" xml:"TXN_TM"`
	TxnStffID         string `json:"TXN_STFF_ID" xml:"TXN_STFF_ID,omitempty"`
	MultiTenanCyID    string `json:"MULTI_TENANCY_ID" xml:"MULTI_TENANCY_ID"`
	LngID             string `json:"LNG_ID" xml:"LNG_ID"`
}

type TXBodyCommonCom4 struct {
	RecInPage string `json:"REC_IN_PAGE" xml:"REC_IN_PAGE,omitempty"`
	PageJump  string `json:"PAGE_JUMP" xml:"PAGE_JUMP,omitempty"`
}

type TXBodyCommon struct {
	ComB *TXBodyCommonComB `json:"COMB" xml:"COMB,omitempty"`
	Com1 *TXBodyCommonCom1 `json:"COM1" xml:"COM1,omitempty"`
	Com4 *TXBodyCommonCom4 `json:"COM4" xml:"COM4,omitempty"`
}

type TXBodyEntity struct {
	AppEntity *TxBodyEntityAppEntity `json:"APP_ENTITY" xml:"APP_ENTITY,omitempty"`
}

type TXBody struct {
	Common *TXBodyCommon `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type TXEmb struct {
}

type WfeCreateRequest struct {
	XMLName  xml.Name  `xml:"TX"`
	TXHeader *TXHeader `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXBody   `json:"TX_BODY" xml:"TX_BODY"`
}

type TXResponseBodyEntity struct {
	ProcessInstID string `json:"PROCESS_INST_ID" xml:"PROCESS_INST_ID,omitempty"`
	EvtTraceID    string `json:"EVT_TRACE_ID" xml:"EVT_TRACE_ID,omitempty"`
}

type TXResponseBody struct {
	Common *TXBodyCommon         `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXResponseBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeCreateResponse struct {
	XMLName  xml.Name        `json:"TX" xml:"TX"`
	TXHeader *TXHeader       `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXResponseBody `json:"TX_BODY" xml:"TX_BODY"`
	TXEmb    *TXEmb          `json:"TX_EMB" xml:"TX_EMB"`
}

func wfeRequest(inputs interface{}) ([]byte, error) {
	byts, _ := xml.MarshalIndent(inputs, " ", "    ")
	url := fmt.Sprintf("%s/itdm-prod-all/http", viper.GetString("wfe.server.addr"))
	request := httplib.Post(url)
	request.SetTimeout(viper.GetDuration("wfe.server.conn_timeout")*time.Millisecond, viper.GetDuration("wfe.server.rw_timeout")*time.Millisecond)
	request.Body(byts)
	response, err := request.Bytes()
	log.Printf("wfe request: %s, response, %s", string(byts), string(response))
	return response, err
}

func (s *wfeService) Create(u *uic.User, templateID string, taskID string, remark string, buttonName string, nextUserIDs []string, nextNode TXNextNodeInfo,
	prjID string, prjSN string, todoTmTpCd string, todoTmTtl string, blngInstID string, dmnGrpID string) (*WfeCreateResponse, error) {
	inst, _ := InstService.GetUserInst(u)

	var nextUserGrp []*NextUserGrp
	for _, userID := range nextUserIDs {
		nextUser, _ := UserService.Get(userID)
		if nextUser != nil {
			nextUserInst, _ := InstService.GetUserInst(nextUser)
			nextUserGrp = append(nextUserGrp, &NextUserGrp{
				Type:         "G",
				ID:           nextUser.JgygUserID,
				Name:         nextUser.CnName,
				PrcActionID:  nextNode.NodeID,
				UsrIDLandNm:  nextUser.UserName,
				CurUsrInstID: nextUserInst.InstID,
				CurUsrInstNm: nextUserInst.Name,
			})
		}
	}

	request := &WfeCreateRequest{
		TXHeader: WfeService.GenTXHeader("A0902S101"),
		TXBody: &TXBody{
			Common: s.GenTXBodyCommon(u),
			Entity: &TXBodyEntity{
				AppEntity: &TxBodyEntityAppEntity{
					TemplateID:   templateID,
					TaskID:       taskID,
					Remark:       remark,
					ButtonName:   buttonName,
					NextUserFlag: "0",
					ExecutorInfo: &ExecutorInfo{
						ID:           u.JgygUserID,
						Name:         u.CnName,
						UsrIDLandNm:  u.UserName,
						CurUsrInstID: inst.InstID,
						CurUsrInstNm: inst.Name,
					},
					NextUserGrp: nextUserGrp,
					BsnComInfo: &BsnComInfo{
						PrjID:      prjID,
						PrjSN:      prjSN,
						TodoTmTpCd: todoTmTpCd,
						TodoTmTtl:  todoTmTtl,
						BlngInstID: blngInstID,
						DmnGrpID:   dmnGrpID,
					},
				},
			},
		},
	}

	return s.create(request)
}

// 创建
func (s *wfeService) create(inputs *WfeCreateRequest) (*WfeCreateResponse, error) {
	response, err := wfeRequest(inputs)
	if err != nil {
		return nil, err
	}

	var wfeResponse WfeCreateResponse
	err = xml.Unmarshal(response, &wfeResponse)
	if err != nil {
		return nil, err
	}
	return &wfeResponse, err
}

type WfeExecuteRequest struct {
	XMLName  xml.Name  `xml:"TX"`
	TXHeader *TXHeader `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXBody   `json:"TX_BODY" xml:"TX_BODY"`
}

type WfeExecuteResponse struct {
	XMLName  xml.Name        `json:"TX" xml:"TX"`
	TXHeader *TXHeader       `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXResponseBody `json:"TX_BODY" xml:"TX_BODY"`
	TXEmb    *TXEmb          `json:"TX_EMB" xml:"TX_EMB"`
}

func (s *wfeService) Execute(u *uic.User, templateID string, processInstID string, taskID string, remark string,
	opinCode string, opinDesc string, buttonName string, exeFstTask string, nextUserIDs []string, nextNode TXNextNodeInfo,
	conditions []Condition) (*WfeExecuteResponse, error) {
	inst, err := InstService.GetUserInst(u)
	if err != nil {
		return nil, err
	}

	var nextUserGrp []*NextUserGrp
	for _, userID := range nextUserIDs {
		if userID == "" {
			continue
		}
		nextUser, err2 := UserService.Get(userID)
		if err2 != nil {
			continue
		}
		if nextUser != nil {
			nextUserInst, err3 := InstService.GetUserInst(nextUser)
			if err3 != nil {
				continue
			}
			nextUserGrp = append(nextUserGrp, &NextUserGrp{
				Type:         "G",
				ID:           nextUser.JgygUserID,
				Name:         nextUser.CnName,
				PrcActionID:  nextNode.NodeID,
				UsrIDLandNm:  nextUser.UserName,
				CurUsrInstID: nextUserInst.InstID,
				CurUsrInstNm: nextUserInst.Name,
			})
		}
	}

	var conds []*CONDITION
	for _, cond := range conditions {
		conds = append(conds, &CONDITION{
			Type:  "G",
			Key:   cond.Key,
			Value: cond.Value,
		})
	}

	appEntity := &TxBodyEntityAppEntity{
		ProcessInstID: processInstID,
		TemplateID:    templateID,
		TaskID:        taskID,
		Remark:        remark,
		OpinCode:      opinCode,
		OpinDesc:      opinDesc,
		ButtonName:    buttonName,
		NextUserFlag:  nextNode.NextUserFlag,
		ExeFstTsk:     exeFstTask,
		ExecutorInfo: &ExecutorInfo{
			ID:           u.JgygUserID,
			Name:         u.CnName,
			UsrIDLandNm:  u.UserName,
			CurUsrInstID: inst.InstID,
			CurUsrInstNm: inst.Name,
		},
		NextUserGrp: nextUserGrp,
		Conditions:  conds,
	}

	request := &WfeExecuteRequest{
		TXHeader: WfeService.GenTXHeader("A0902S102"),
		TXBody: &TXBody{
			Common: WfeService.GenTXBodyCommon(u),
			Entity: &TXBodyEntity{
				AppEntity: appEntity,
			},
		},
	}

	return s.execute(request)
}

// 执行
func (s *wfeService) execute(inputs *WfeExecuteRequest) (*WfeExecuteResponse, error) {
	response, err := wfeRequest(inputs)
	if err != nil {
		return nil, err
	}

	var wfeResponse WfeExecuteResponse
	err = xml.Unmarshal(response, &wfeResponse)
	if err != nil {
		return nil, err
	}
	return &wfeResponse, err
}

func (s *wfeService) GenTXHeader(outBoundID string) *TXHeader {
	return &TXHeader{
		SysHdrLen:        100,
		SysPkgVrsn:       "v1",
		SysTtlLen:        400,
		SysReqSecID:      "123456",
		SysSndSecID:      "987654",
		SysTxCode:        outBoundID,
		SysTxVrsn:        "v0",
		SysTxType:        "00000",
		SysReserved:      "AA",
		SysEvtTraceID:    "1147319008",
		SysSndSerialNo:   "0000000000",
		SysPkgType:       "A",
		SysMsgLen:        100,
		SysIsEncrypted:   "N",
		SysEncryptType:   "R",
		SysCompressType:  "T",
		SysEmbMsgLen:     200,
		SysReqTime:       "20130318092311000",
		SysTimeLeft:      "000030000",
		SysPkgStsType:    "00",
		SysSecContextLen: 8,
		SysSecContext:    "fengshijie",
	}
}

func (s *wfeService) GenTXBodyCommon(u *uic.User) *TXBodyCommon {
	bodyCommonCom1 := &TXBodyCommonCom1{
		TxnInsID:          "370616138",
		TxnIttChnlID:      "001",
		TxnIttChnlCgyCode: "20180030",
		TxnStffID:         u.JgygUserID,
		TxnDT:             "20140728",
		TxnTM:             "110545",
		MultiTenanCyID:    "CN000",
		LngID:             "zh-cn",
	}
	bodyCommonCom4 := &TXBodyCommonCom4{}

	body := &TXBodyCommon{
		Com1: bodyCommonCom1,
		Com4: bodyCommonCom4,
	}
	return body
}

func (s *wfeService) GenTXBodyCommonPro(u *uic.User, page string, limit string) *TXBodyCommon {
	bodyCommonCom1 := &TXBodyCommonCom1{
		TxnInsID:          "370616138",
		TxnIttChnlID:      "001",
		TxnIttChnlCgyCode: "20180030",
		TxnStffID:         u.JgygUserID,
		TxnDT:             "20140728",
		TxnTM:             "110545",
		MultiTenanCyID:    "CN000",
		LngID:             "zh-cn",
	}
	bodyCommonCom4 := &TXBodyCommonCom4{
		RecInPage: limit,
		PageJump:  page,
	}

	body := &TXBodyCommon{
		Com1: bodyCommonCom1,
		Com4: bodyCommonCom4,
	}
	return body
}

type TxHistDetailsBodyEntityAppEntity struct {
	ProcessInstID string `json:"PROCESS_INST_ID" xml:"PROCESS_INST_ID"`
	SelectMode    string `json:"SELECT_MODE" xml:"SELECT_MODE"`
}

type TXHistDetailsBodyEntity struct {
	AppEntity *TxHistDetailsBodyEntityAppEntity `json:"APP_ENTITY" xml:"APP_ENTITY,omitempty"`
}

type TXHistDetailsBody struct {
	Common *TXBodyCommon            `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXHistDetailsBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeHistDetailsRequest struct {
	XMLName  xml.Name           `xml:"TX"`
	TXHeader *TXHeader          `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXHistDetailsBody `json:"TX_BODY" xml:"TX_BODY"`
}

type UserInfo struct {
	Type        string `json:"type" xml:"type,attr"`
	ID          string `json:"ID" xml:"ID"`
	Name        string `json:"NAME" xml:"NAME"`
	BlngInstID  string `json:"BLNG_INST_ID" xml:"BLNG_INST_ID"`
	BlngInstNm  string `json:"BLNG_INST_NM" xml:"BLNG_INST_NM"`
	UsrIDLandNm string `json:"USR_ID_LAND_NM" xml:"USR_ID_LAND_NM"`
}

type ProcessInfo struct {
	Type               string        `json:"type" xml:"type,attr"`
	ProcessInstID      string        `json:"PROCESS_INST_ID" xml:"PROCESS_INST_ID"`
	PcsAvyID           string        `json:"PCS_AVY_ID" xml:"PCS_AVY_ID"`
	PcsAvyNM           string        `json:"PCS_AVY_NM" xml:"PCS_AVY_NM"`
	ProcessInstType    string        `json:"PROCESS_INST_TYPE" xml:"PROCESS_INST_TYPE"`
	StartTime          string        `json:"START_TIME" xml:"START_TIME"`
	EndTime            string        `json:"END_TIME" xml:"END_TIME"`
	ApproveResult      string        `json:"APPROVE_RESULT" xml:"APPROVE_RESULT"`
	ApproveDesc        string        `json:"APPROVE_DESC" xml:"APPROVE_DESC"`
	ApproveCode        string        `json:"APPROVE_CODE" xml:"APPROVE_CODE"`
	PcsAvyStatus       string        `json:"PCS_AVY_STATUS" xml:"PCS_AVY_STATUS"`
	ButtonName         string        `json:"BUTTON_NAME" xml:"BUTTON_NAME"`
	FlodType           string        `json:"FLOD_TYPE" xml:"FLOD_TYPE"`
	EndFlag            string        `json:"END_FLAG" xml:"END_FLAG"`
	SubPcsID           string        `json:"SUB_PCS_ID" xml:"SUB_PCS_ID"`
	HistID             string        `json:"HIST_ID" xml:"HIST_ID"`
	OrigiHistID        string        `json:"ORIGI_HIST_ID" xml:"ORIGI_HIST_ID"`
	IttendrsTodoSvcEcd string        `json:"ITTENDRS_TODO_SVC_ECD" xml:"ITTENDRS_TODO_SVC_ECD"`
	ExecuterInfo       *ExecutorInfo `json:"EXECUTER_INFO" xml:"EXECUTER_INFO"`
	UserInfo           []*UserInfo   `json:"USER_INFO" xml:"USER_INFO"`
}

type TXHistDetailsResponseBodyEntity struct {
	ProcessInfo []*ProcessInfo `json:"PROCESS_INFO" xml:"PROCESS_INFO"`
}

type TXHistDetailsResponseBody struct {
	Common *TXBodyCommon                    `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXHistDetailsResponseBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeHistDetailsResponse struct {
	XMLName  xml.Name                   `json:"TX" xml:"TX"`
	TXHeader *TXHeader                  `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXHistDetailsResponseBody `json:"TX_BODY" xml:"TX_BODY"`
	TXEmb    *TXEmb                     `json:"TX_EMB" xml:"TX_EMB"`
}

func (s *wfeService) HistDetails(u *uic.User, processInstID string) (*WfeHistDetailsResponse, error) {
	request := &WfeHistDetailsRequest{
		TXHeader: WfeService.GenTXHeader("A0902S124"),
		TXBody: &TXHistDetailsBody{
			Common: WfeService.GenTXBodyCommon(u),
			Entity: &TXHistDetailsBodyEntity{
				AppEntity: &TxHistDetailsBodyEntityAppEntity{
					ProcessInstID: processInstID,
					SelectMode:    "3",
				},
			},
		},
	}
	return s.histDetails(request)
}

// 历史记录
func (s *wfeService) histDetails(inputs *WfeHistDetailsRequest) (*WfeHistDetailsResponse, error) {
	response, err := wfeRequest(inputs)
	if err != nil {
		return nil, err
	}

	var wfeResponse WfeHistDetailsResponse
	err = xml.Unmarshal(response, &wfeResponse)
	if err != nil {
		return nil, err
	}
	return &wfeResponse, err
}

type SortField struct {
	Type     string `json:"type" xml:"type,attr"`
	FieldNm  string `json:"FIELD_NM" xml:"FIELD_NM,omitempty"`
	FieldAsc string `json:"FIELD_ASC" xml:"FIELD_ASC,omitempty"`
}

type TXTodosBodyEntityComSortEntity struct {
	SortFields []*SortField `json:"SORT_FIELDS" xml:"SORT_FIELDS,omitempty"`
}

type TXTodosBodyEntityAppEntity struct {
	//ProcessInstIDList string `json:"PROCESS_INST_ID_LIST" xml:"PROCESS_INST_ID_LIST"`
	TimeStart   string `json:"TIME_START" xml:"TIME_START,omitempty"`
	TimeEnd     string `json:"TIME_END" xml:"TIME_END,omitempty"`
	PrjID       string `json:"PRJ_ID" xml:"PRJ_ID,omitempty"`
	PrjTypeList string `json:"PRJ_TYPE_LIST" xml:"PRJ_TYPE_LIST,omitempty"`
	//PrjBelongTypeList string `json:"PRJ_BELONG_TYPE_LIST" xml:"PRJ_BELONG_TYPE_LIST,omitempty"`
	PrjNm    string `json:"PRJ_NM" xml:"PRJ_NM,omitempty"`
	WfExtrNm string `json:"WF_EXTR_NM" xml:"WF_EXTR_NM,omitempty"`
	//AvyOwrNm string `json:"AVY_OWR_NM" xml:"AVY_OWR_NM,omitempty"`
	TodoType string `json:"TODO_TYPE" xml:"TODO_TYPE"`
}

type TXTodosBodyEntity struct {
	ComSortEntity *TXTodosBodyEntityComSortEntity `json:"COM_SORT_ENTITY" xml:"COM_SORT_ENTITY,omitempty"`
	AppEntity     *TXTodosBodyEntityAppEntity     `json:"APP_ENTITY" xml:"APP_ENTITY,omitempty"`
}

type TXTodosBody struct {
	Common *TXBodyCommon      `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXTodosBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeTodosRequest struct {
	XMLName  xml.Name     `xml:"TX"`
	TXHeader *TXHeader    `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXTodosBody `json:"TX_BODY" xml:"TX_BODY"`
}

type TXTodosResponseBodyEntityTodoInfo struct {
	Type          string `json:"type" xml:"type,attr,omitempty"`
	ProcessInstID string `json:"PROCESS_INST_ID" xml:"PROCESS_INST_ID,omitempty"`
	TaskID        string `json:"TASK_ID" xml:"TASK_ID,omitempty"`
	PcsAvyNm      string `json:"PCS_AVY_NM" xml:"PCS_AVY_NM,omitempty"`
	WfExtrID      string `json:"WF_EXTR_ID" xml:"WF_EXTR_ID,omitempty"`
	WfExtrNm      string `json:"WF_EXTR_NM" xml:"WF_EXTR_NM,omitempty"`
	TodoStartTm   string `json:"TO_START_TM" xml:"TO_START_TM,omitempty"`
	TemplateID    string `json:"TEMPLATE_ID" xml:"TEMPLATE_ID,omitempty"`
	PcsAvyStatus  string `json:"PCS_AVY_STATUS" xml:"PCS_AVY_STATUS,omitempty"`
	PcsStatus     string `json:"PCS_STATUS" xml:"PCS_STATUS,omitempty"`
	BlngInstID    string `json:"BLNG_INST_ID" xml:"BLNG_INST_ID,omitempty"`
	SourceType    string `json:"SOURCE_TYPE" xml:"SOURCE_TYPE,omitempty"`
	FormEdit      string `json:"FORM_EDIT" xml:"FORM_EDIT,omitempty"`
	TodoSN        string `json:"TODO_SN" xml:"TODO_SN,omitempty"`
	PrjID         string `json:"PRJ_ID" xml:"PRJ_ID,omitempty"`
	PrjSN         string `json:"PRJ_SN" xml:"PRJ_SN,omitempty"`
	PrjNM         string `json:"PRJ_NM" xml:"PRJ_NM,omitempty"`
	PrjType       string `json:"PRJ_TYPE" xml:"PRJ_TYPE,omitempty"`
	DmnGrpID      string `json:"DMN_GRP_ID" xml:"DMN_GRP_ID,omitempty"`
	AgentType     string `json:"AGENT_TYPE" xml:"AGENT_TYPE,omitempty"`
	PrjBelongType string `json:"PRJ_BELONG_TYPE" xml:"PRJ_BELONG_TYPE,omitempty"`
	SkipPcsAvyID  string `json:"SKIP_PCS_AVY_ID" xml:"SKIP_PCS_AVY_ID,omitempty"`
	URL           string `json:"URL" xml:"URL,omitempty"`
}

type TXTodosResponseBodyEntity struct {
	TodoInfo []*TXTodosResponseBodyEntityTodoInfo `json:"TODO_INFO" xml:"TODO_INFO,omitempty"`
}

type TXTodosResponseBody struct {
	Common *TXBodyCommon              `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXTodosResponseBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeTodosResponse struct {
	XMLName  xml.Name             `json:"TX" xml:"TX"`
	TXHeader *TXHeader            `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXTodosResponseBody `json:"TX_BODY" xml:"TX_BODY"`
	TXEmb    *TXEmb               `json:"TX_EMB" xml:"TX_EMB"`
}

func (s *wfeService) Todos(u *uic.User, timeStart string, timeEnd string, prjID string, prjTypeList string,
	prjNm string, wfExtrNm string, sortFields []*SortField, page string, limit string) (*WfeTodosResponse, error) {
	if sortFields == nil {
		sortFields = append(sortFields, &SortField{
			Type:     "G",
			FieldNm:  "pcsAvyEfdt",
			FieldAsc: "desc",
		})
	}
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	request := &WfeTodosRequest{
		TXHeader: WfeService.GenTXHeader("A0902S119"),
		TXBody: &TXTodosBody{
			Common: WfeService.GenTXBodyCommonPro(u, page, limit),
			Entity: &TXTodosBodyEntity{
				ComSortEntity: &TXTodosBodyEntityComSortEntity{
					SortFields: sortFields,
				},
				AppEntity: &TXTodosBodyEntityAppEntity{
					TimeStart:   timeStart,
					TimeEnd:     timeEnd,
					PrjID:       prjID,
					PrjTypeList: prjTypeList,
					PrjNm:       prjNm,
					WfExtrNm:    wfExtrNm,
					TodoType:    "2",
				},
			},
		},
	}
	return s.todos(request)
}

func (s *wfeService) todos(inputs *WfeTodosRequest) (*WfeTodosResponse, error) {
	response, err := wfeRequest(inputs)
	if err != nil {
		return nil, err
	}

	var wfeResponse WfeTodosResponse
	err = xml.Unmarshal(response, &wfeResponse)
	if err != nil {
		return nil, err
	}
	return &wfeResponse, err
}

type WfeNextNodeInfoRequest struct {
	XMLName  xml.Name  `xml:"TX"`
	TXHeader *TXHeader `json:"TXHeader" xml:"TX_HEADER"`
	TXBody   *TXBody   `json:"TXBody" xml:"TX_BODY"`
}

type TXNextNodeInfo struct {
	Type         string `json:"type" xml:"type,attr,omitempty"`
	NodeID       string `json:"NODE_ID" xml:"NODE_ID,omitempty"`
	NodeName     string `json:"NODE_NAME" xml:"NODE_NAME,omitempty"`
	NextUserFlag string `json:"NEXT_USER_FLAG" xml:"NEXT_USER_FLAG,omitempty"`
	Multiple     string `json:"multiple" xml:"multiple,omitempty"`
}

type TXNextNodeInfoResponseBodyEntity struct {
	NextNodeInfo []*TXNextNodeInfo `json:"NEXT_NODE_INFO" xml:"NEXT_NODE_INFO,omitempty"`
}

type TXNextNodeInfoBody struct {
	Common *TXBodyCommon                     `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXNextNodeInfoResponseBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeNextNodeInfoResponse struct {
	XMLName  xml.Name            `json:"TX" xml:"TX"`
	TXHeader *TXHeader           `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXNextNodeInfoBody `json:"TX_BODY" xml:"TX_BODY"`
	TXEmb    *TXEmb              `json:"TX_EMB" xml:"TX_EMB"`
}

// 下个节点信息
func (s *wfeService) NextNodeInfo(u *uic.User, templateID string, processInstID string, taskID string) (*WfeNextNodeInfoResponse, error) {
	request := &WfeNextNodeInfoRequest{
		TXHeader: WfeService.GenTXHeader("A0902S112"),
		TXBody: &TXBody{
			Common: WfeService.GenTXBodyCommon(u),
			Entity: &TXBodyEntity{
				AppEntity: &TxBodyEntityAppEntity{
					ProcessInstID: processInstID,
					TemplateID:    templateID,
					TaskID:        taskID,
				},
			},
		},
	}
	return s.nextNodeInfo(request)
}

func (s *wfeService) nextNodeInfo(inputs *WfeNextNodeInfoRequest) (*WfeNextNodeInfoResponse, error) {
	response, err := wfeRequest(inputs)
	if err != nil {
		return nil, err
	}

	var wfeResponse WfeNextNodeInfoResponse
	err = xml.Unmarshal(response, &wfeResponse)
	if err != nil {
		return nil, err
	}
	return &wfeResponse, err
}

type TxTodo2DoingBodyEntityAppEntity struct {
	TodoID string `json:"TODO_ID" xml:"TODO_ID"`
}

type TXTodo2DoingBodyEntity struct {
	AppEntity *TxTodo2DoingBodyEntityAppEntity `json:"APP_ENTITY" xml:"APP_ENTITY,omitempty"`
}

type TXTodo2DoingBody struct {
	Common *TXBodyCommon           `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXTodo2DoingBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeTodo2DoingRequest struct {
	XMLName  xml.Name          `xml:"TX"`
	TXHeader *TXHeader         `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXTodo2DoingBody `json:"TX_BODY" xml:"TX_BODY"`
}

type TXTodo2DoingResponseBodyEntity struct {
	ResultDesc string `json:"RESULT_DESC" xml:"RESULT_DESC,omitempty"`
}

type TXTodo2DoingResponseBody struct {
	Common *TXBodyCommon                   `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXTodo2DoingResponseBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeTodo2DoingResponse struct {
	XMLName  xml.Name                  `json:"TX" xml:"TX"`
	TXHeader *TXHeader                 `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXTodo2DoingResponseBody `json:"TX_BODY" xml:"TX_BODY"`
	TXEmb    *TXEmb                    `json:"TX_EMB" xml:"TX_EMB"`
}

func (s *wfeService) Todo2Doing(u *uic.User, todoID string) (*WfeTodo2DoingResponse, error) {
	request := &WfeTodo2DoingRequest{
		TXHeader: WfeService.GenTXHeader("A0902S132"),
		TXBody: &TXTodo2DoingBody{
			Common: WfeService.GenTXBodyCommon(u),
			Entity: &TXTodo2DoingBodyEntity{
				AppEntity: &TxTodo2DoingBodyEntityAppEntity{
					TodoID: todoID,
				},
			},
		},
	}
	return s.todo2Doing(request)
}

func (s *wfeService) todo2Doing(inputs *WfeTodo2DoingRequest) (*WfeTodo2DoingResponse, error) {
	response, err := wfeRequest(inputs)
	if err != nil {
		return nil, err
	}

	var wfeResponse WfeTodo2DoingResponse
	err = xml.Unmarshal(response, &wfeResponse)
	if err != nil {
		return nil, err
	}
	return &wfeResponse, err
}

type TXHasDoneBodyEntityComSortEntity struct {
	SortFields []*SortField `json:"SORT_FIELDS" xml:"SORT_FIELDS,omitempty"`
}

type TXHasDoneBodyEntityAppEntity struct {
	TimeStart   string `json:"TIME_START" xml:"TIME_START,omitempty"`
	TimeEnd     string `json:"TIME_END" xml:"TIME_END,omitempty"`
	PrjID       string `json:"PRJ_ID" xml:"PRJ_ID,omitempty"`
	PrjTypeList string `json:"PRJ_TYPE_LIST" xml:"PRJ_TYPE_LIST,omitempty"`
	PrjNm       string `json:"PRJ_NM" xml:"PRJ_NM,omitempty"`
	WfExtrNm    string `json:"WF_EXTR_NM" xml:"WF_EXTR_NM,omitempty"`
	HistType    string `json:"HIST_TYPE" xml:"HIST_TYPE"`
}

type TXHasDoneBodyEntity struct {
	ComSortEntity *TXHasDoneBodyEntityComSortEntity `json:"COM_SORT_ENTITY" xml:"COM_SORT_ENTITY,omitempty"`
	AppEntity     *TXHasDoneBodyEntityAppEntity     `json:"APP_ENTITY" xml:"APP_ENTITY,omitempty"`
}

type TXHasDoneBody struct {
	Common *TXBodyCommon        `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXHasDoneBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeHasDoneRequest struct {
	XMLName  xml.Name       `xml:"TX"`
	TXHeader *TXHeader      `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXHasDoneBody `json:"TX_BODY" xml:"TX_BODY"`
}

type TXHasDoneResponseBodyEntityHasDoneInfo struct {
	Type          string `json:"type" xml:"type,attr,omitempty"`
	ProcessInstID string `json:"PROCESS_INST_ID" xml:"PROCESS_INST_ID,omitempty"`
	TaskID        string `json:"TASK_ID" xml:"TASK_ID,omitempty"`
	PcsAvyNm      string `json:"PCS_AVY_NM" xml:"PCS_AVY_NM,omitempty"`
	WfExtrID      string `json:"WF_EXTR_ID" xml:"WF_EXTR_ID,omitempty"`
	WfExtrNm      string `json:"WF_EXTR_NM" xml:"WF_EXTR_NM,omitempty"`
	TodoStartTm   string `json:"TO_START_TM" xml:"TO_START_TM,omitempty"`
	TemplateID    string `json:"TEMPLATE_ID" xml:"TEMPLATE_ID,omitempty"`
	PcsAvyStatus  string `json:"PCS_AVY_STATUS" xml:"PCS_AVY_STATUS,omitempty"`
	PcsStatus     string `json:"PCS_STATUS" xml:"PCS_STATUS,omitempty"`
	BlngInstID    string `json:"BLNG_INST_ID" xml:"BLNG_INST_ID,omitempty"`
	SourceType    string `json:"SOURCE_TYPE" xml:"SOURCE_TYPE,omitempty"`
	FormEdit      string `json:"FORM_EDIT" xml:"FORM_EDIT,omitempty"`
	TodoSN        string `json:"TODO_SN" xml:"TODO_SN,omitempty"`
	PrjID         string `json:"PRJ_ID" xml:"PRJ_ID,omitempty"`
	PrjSN         string `json:"PRJ_SN" xml:"PRJ_SN,omitempty"`
	PrjNM         string `json:"PRJ_NM" xml:"PRJ_NM,omitempty"`
	PrjType       string `json:"PRJ_TYPE" xml:"PRJ_TYPE,omitempty"`
	DmnGrpID      string `json:"DMN_GRP_ID" xml:"DMN_GRP_ID,omitempty"`
	AgentType     string `json:"AGENT_TYPE" xml:"AGENT_TYPE,omitempty"`
	PrjBelongType string `json:"PRJ_BELONG_TYPE" xml:"PRJ_BELONG_TYPE,omitempty"`
	SkipPcsAvyID  string `json:"SKIP_PCS_AVY_ID" xml:"SKIP_PCS_AVY_ID,omitempty"`
	URL           string `json:"URL" xml:"URL,omitempty"`
}

type TXHasDoneResponseBodyEntity struct {
	HasDoneInfo []*TXHasDoneResponseBodyEntityHasDoneInfo `json:"HASDONE_INFO" xml:"HASDONE_INFO,omitempty"`
}

type TXHasDoneResponseBody struct {
	Common *TXBodyCommon                `json:"COMMON" xml:"COMMON,omitempty"`
	Entity *TXHasDoneResponseBodyEntity `json:"ENTITY" xml:"ENTITY,omitempty"`
}

type WfeHasDoneResponse struct {
	XMLName  xml.Name             `json:"TX" xml:"TX"`
	TXHeader *TXHeader            `json:"TX_HEADER" xml:"TX_HEADER"`
	TXBody   *TXTodosResponseBody `json:"TX_BODY" xml:"TX_BODY"`
	TXEmb    *TXEmb               `json:"TX_EMB" xml:"TX_EMB"`
}

func (s *wfeService) HasDone(u *uic.User, timeStart string, timeEnd string, prjID string, prjTypeList string,
	prjNm string, wfExtrNm string, sortFields []*SortField, page string, limit string) (*WfeHasDoneResponse, error) {
	if sortFields == nil {
		sortFields = append(sortFields, &SortField{
			Type:     "G",
			FieldNm:  "pcsAvyEfdt",
			FieldAsc: "desc",
		})
	}
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	request := &WfeHasDoneRequest{
		TXHeader: WfeService.GenTXHeader("A0902S120"),
		TXBody: &TXHasDoneBody{
			Common: WfeService.GenTXBodyCommonPro(u, page, limit),
			Entity: &TXHasDoneBodyEntity{
				ComSortEntity: &TXHasDoneBodyEntityComSortEntity{
					SortFields: sortFields,
				},
				AppEntity: &TXHasDoneBodyEntityAppEntity{
					TimeStart:   timeStart,
					TimeEnd:     timeEnd,
					PrjID:       prjID,
					PrjTypeList: prjTypeList,
					PrjNm:       prjNm,
					WfExtrNm:    wfExtrNm,
					HistType:    "2",
				},
			},
		},
	}
	return s.hasDone(request)
}

func (s *wfeService) hasDone(inputs *WfeHasDoneRequest) (*WfeHasDoneResponse, error) {
	response, err := wfeRequest(inputs)
	if err != nil {
		return nil, err
	}

	var wfeResponse WfeHasDoneResponse
	err = xml.Unmarshal(response, &wfeResponse)
	if err != nil {
		return nil, err
	}
	return &wfeResponse, err
}

func newWfeService() *wfeService {
	return &wfeService{}
}
