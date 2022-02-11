package service

import (
	"encoding/xml"
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/uic"
	log "github.com/sirupsen/logrus"
	"testing"
)

func init() {
	g.ParseConfig("../../config/cfg.json")
}

func TestWfeService(t *testing.T) {
	//wfe := newWfeService()
	//wfe.init()
	response := `
<TX>
	<TX_HEADER>
		<SYS_HDR_LEN></SYS_HDR_LEN>
		<SYS_PKG_VRSN>01</SYS_PKG_VRSN>
		<SYS_REQ_SEC_ID>108048</SYS_REQ_SEC_ID>
		<SYS_SND_SEC_ID>108048</SYS_SND_SEC_ID>
		<SYS_TX_TYPE>00000</SYS_TX_TYPE>
		<SYS_EVT_TRACE_ID>1147319008</SYS_EVT_TRACE_ID>
		<SYS_SND_SERIAL_NO>0000000000</SYS_SND_SERIAL_NO>
		<SYS_PKG_TYPE>1</SYS_PKG_TYPE>
		<SYS_MSG_LEN></SYS_MSG_LEN>
		<SYS_IS_ENCRYPTED>0</SYS_IS_ENCRYPTED>
		<SYS_ENCRYPT_TYPE>3</SYS_ENCRYPT_TYPE>
		<SYS_COMPRESS_TYPE>0</SYS_COMPRESS_TYPE>
		<SYS_EMB_MSG_LEN></SYS_EMB_MSG_LEN>
		<SYS_RECV_TIME>20220129104341504</SYS_RECV_TIME>
		<SYS_RESP_TIME>20220129104341900</SYS_RESP_TIME>
		<SYS_PKG_STS_TYPE>01</SYS_PKG_STS_TYPE>
		<SYS_RESP_CODE>0000000000</SYS_RESP_CODE>
		<SYS_RESP_DESC_LEN></SYS_RESP_DESC_LEN>
		<SYS_RESP_DESC>成功</SYS_RESP_DESC>
	</TX_HEADER>
	<TX_BODY>
		<COMMON>
			<FILE_LIST_PACK>
				<FILE_NUM></FILE_NUM>
				<FILE_MODE></FILE_MODE>
				<FILE_NODE></FILE_NODE>
				<FILE_NAME_PACK></FILE_NAME_PACK>
				<FILE_PATH_PACK></FILE_PATH_PACK>
			</FILE_LIST_PACK>
			<COMB>
				<ERR_MSG_NUM>0</ERR_MSG_NUM>
				<CMPT_TRCNO></CMPT_TRCNO>
				<TOTAL_PAGE>0</TOTAL_PAGE>
				<TOTAL_REC>0</TOTAL_REC>
				<CURR_TOTAL_PAGE>0</CURR_TOTAL_PAGE>
				<CURR_TOTAL_REC>0</CURR_TOTAL_REC>
				<STS_TRACE_ID></STS_TRACE_ID>
			</COMB>
		</COMMON>
		<ENTITY>
			<PROCESS_INST_ID>20220129104341000000046478</PROCESS_INST_ID>
			<EVT_TRACE_ID></EVT_TRACE_ID>
		</ENTITY>
	</TX_BODY>
	<TX_EMB></TX_EMB>
</TX>
`
	var wfeResponse WfeCreateResponse
	err := xml.Unmarshal([]byte(response), &wfeResponse)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Printf("%+v", wfeResponse)
}

func TestWfeProcExecute(t *testing.T) {
	var nextUserGrp []*NextUserGrp
	nextUserGrp = append(nextUserGrp, &NextUserGrp{
		Type:         "G",
		ID:           "nextUser.JgygUserID",
		Name:         "nextUser.CnName",
		PrcActionID:  "",
		UsrIDLandNm:  "nextUser.UserName",
		CurUsrInstID: "nextUserInst.InstID",
		CurUsrInstNm: "nextUserInst.Name",
	})

	var conditions []*CONDITION
	conditions = append(conditions, &CONDITION{
		Type:  "G",
		Key:   "cond.Key",
		Value: "cond.Value",
	})

	req := WfeExecuteRequest{
		TXHeader: &TXHeader{
			SysHdrLen:        100,
			SysPkgVrsn:       "v1",
			SysTtlLen:        400,
			SysReqSecID:      "123456",
			SysSndSecID:      "987654",
			SysTxCode:        "A0902S102",
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
		},
		TXBody: &TXBody{
			Common: &TXBodyCommon{
				Com1: &TXBodyCommonCom1{
					TxnInsID:          "370616138",
					TxnIttChnlID:      "001",
					TxnIttChnlCgyCode: "20180030",
					TxnDT:             "20140728",
					TxnTM:             "110545",
					TxnStffID:         "1111111",
					MultiTenanCyID:    "CN000",
					LngID:             "zh-cn",
				},
				Com4: &TXBodyCommonCom4{},
			},
			Entity: &TXBodyEntity{
				AppEntity: &TxBodyEntityAppEntity{
					ProcessInstID: "202201291144",
					TemplateID:    "600100PubAudit",
					TaskID:        "10102",
					Remark:        "1111111",
					OpinCode:      "1",
					OpinDesc:      "test",
					ButtonName:    "提交",
					NextUserFlag:  "0",
					ExecutorInfo: &ExecutorInfo{
						ID:           "1111111",
						Name:         "小强",
						UsrIDLandNm:  "xiaoqiang.zh",
						CurUsrInstID: "1111111111",
						CurUsrInstNm: "小强之家",
					},
					NextUserGrp: nextUserGrp,
					Conditions:  conditions,
					BsnComInfo:  nil,
					//BsnComInfo: service.BsnComInfo{
					//	PrjID:      inputs.PrjID,
					//	PrjSN:      inputs.PrjSN,
					//	TodoTmTpCd: inputs.TodoTmTpCd,
					//	TodoTmTtl:  inputs.TodoTmTtl,
					//	BlngInstID: inputs.BlngInstID,
					//	DmnGrpID:   inputs.DmnGrpID,
					//},
				},
			},
		},
	}
	resp, _ := WfeService.execute(&req)
	log.Printf("resp: %+v", resp)
}

func TestWfeTodoList(t *testing.T) {
	u := uic.User{
		JgygUserID: "23598915",
	}
	resp, _ := WfeService.Todos(&u, "", "", "",
		"", "", "", nil, "1", "10")

	log.Printf("resp: %+v", resp)
}
