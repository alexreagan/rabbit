package procmanager

import (
	"fmt"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/gin-gonic/gin"
)

const (
	ULI_WEB_SESSION_ID = "uliweb_session_id"
)

type APIPostProcManagerLoginInputs struct {
	TYPE     string
	NEXT     string
	USERNAME string
	PASSWORD string
}

// @Summary 模拟ProcManager Login
// @Description
// @Produce json
// @Param APIPostProcManagerLoginInputs body APIPostProcManagerLoginInputs true "登录表单"
// @Success 200 {object} APIPostRespSess
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /login [POST]
func ProcManagerLogin(c *gin.Context) {
	var inputs APIPostProcManagerLoginInputs
	if err := c.Bind(&inputs); err != nil {
		c.JSON(h.BadStatus, err)
		return
	}
	uliWebSessionId, err := c.Cookie(ULI_WEB_SESSION_ID)
	if err != nil || uliWebSessionId == "" {
		c.SetCookie(ULI_WEB_SESSION_ID, "session:s1aad2d", 3600, "/", "localhost", false, true)
	}
	c.SetCookie(ULI_WEB_SESSION_ID, "session:s1aad2d", 3600, "/", "localhost", false, true)
	uliWebSessionId, err = c.Cookie(ULI_WEB_SESSION_ID)
	fmt.Printf("api_controller, ProcManagerLogin, uliWebSessionId: %v\n", uliWebSessionId)
	c.Redirect(302, "http://localhost:8080/")
	// user, _ := h.GetUser(c)
	// c.JSON(h.OKStatus, inputs)
	return
}

type APIPostProcManagerNextUser struct {
	CUR_USR_INST_ID string
	CUR_USR_INST_NM string
	USR_ID_LAND_NM  string
	PRC_ACTION_ID   string
	ID              string
	NAME            string
}
type APIPostProcManagerCondition struct {
	KEY   string
	VALUE string
}

type APIPostProcManagerCreateInputs struct {
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
	NEXT_USER_GRP   []APIPostProcManagerNextUser
	CONDITIONS      []APIPostProcManagerCondition
}

type APIPostProcManagerCreateOutputs struct {
	SYS_RESP_DESC     string
	SYS_TX_TYPE       string
	SYS_PKG_TYPE      string
	SYS_PKG_VRSN      string
	SYS_RESP_DESC_LEN string
	PROCESS_INST_ID   string
	SYS_TX_STATUS     string
	SYS_RESP_CODE     string
}

func ProcCreate(c *gin.Context) {
	var inputs APIPostProcManagerCreateInputs
	if err := c.Bind(&inputs); err != nil {
		c.JSON(h.BadStatus, err)
		return
	}
	// user, _ := h.GetUser(c)

	var resp *APIPostProcManagerCreateOutputs
	c.JSON(h.OKStatus, resp)
	return
}

func ProcExecute(c *gin.Context) {
	var inputs APIPostProcManagerLoginInputs
	if err := c.Bind(&inputs); err != nil {
		c.JSON(h.BadStatus, err)
		return
	}
	// user, _ := h.GetUser(c)
	c.JSON(h.OKStatus, inputs)
	return
}

func ProcInstTodoInfo(c *gin.Context) {
	var inputs APIPostProcManagerLoginInputs
	if err := c.Bind(&inputs); err != nil {
		c.JSON(h.BadStatus, err)
		return
	}
	// user, _ := h.GetUser(c)
	c.JSON(h.OKStatus, inputs)
	return
}

type APIPostProcManagerNodeInfo struct {
	NODE_NAME             string
	NODE_ID               string
	NEXT_USER_FLAG        string
	DEFAULT_USERS_DISEDIT string
}

type APIPostProcManagerNodeOutputs struct {
	SYS_RESP_CODE  string
	NEXT_NODE_INFO []APIPostProcManagerNodeInfo
	SYS_RESP_DESC  string
}

func ProcNextNodeInfo(c *gin.Context) {
	var inputs APIPostProcManagerLoginInputs
	if err := c.Bind(&inputs); err != nil {
		c.JSON(h.BadStatus, err)
		return
	}
	nextNodeInfos := make([]APIPostProcManagerNodeInfo, 0)
	nextNodeInfos = append(nextNodeInfos, APIPostProcManagerNodeInfo{
		NODE_NAME:             "初审版本发布单",
		NODE_ID:               "10102",
		NEXT_USER_FLAG:        "0",
		DEFAULT_USERS_DISEDIT: "0",
	})
	// user, _ := h.GetUser(c)
	resp := APIPostProcManagerNodeOutputs{
		SYS_RESP_CODE:  "000000000000",
		NEXT_NODE_INFO: nextNodeInfos,
		SYS_RESP_DESC:  "成功",
	}
	c.JSON(h.OKStatus, resp)
	return
}
