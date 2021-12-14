package uic

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/uic"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type APIGetUserListInputs struct {
	UserName string `json:"username" form:"username"`
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

type APIGetUserListOutputs struct {
	List       []*uic.User `json:"list"`
	TotalCount int64       `json:"totalCount"`
}

func (input APIGetUserListInputs) checkInputsContain() error {
	return nil
}

// @Summary 用户列表接口
// @Description
// @Produce json
// @Param APIGetUserListInputs query APIGetUserListInputs true "根据查询条件分页查询用户列表"
// @Success 200 {object} APIGetUserListOutputs
// @Failure 400 {object} APIGetUserListOutputs
// @Router /api/v1/user/list [get]
func UserList(c *gin.Context) {
	var inputs APIGetUserListInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	//for get correct table name
	db := g.Con().Uic.Table(uic.User{}.TableName())
	db = db.Where("username like ? ", inputs.UserName+"%")
	db = db.Or("ad_username like ?", inputs.UserName+"%")
	db = db.Or("nickname like ?", inputs.UserName+"%")

	var totalCount int64
	db.Count(&totalCount)
	var users []*uic.User
	db = db.Order("id DESC").Offset(offset).Limit(limit)
	db.Find(&users)

	resp := &APIGetUserListOutputs{
		List:       users,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIUserCreateInputs struct {
	UserName   string      `json:"username" form:"username" binding:"required"`
	CnName     string      `json:"cnName" form:"cn_name" binding:"required"`
	Password   string      `json:"password" form:"password" binding:"required"`
	JgygUserId string      `json:"jgygUserId" form:"jgyg_user_id" binding:"required"`
	Birthday   gtime.GTime `json:"birthday" form:"birthday"`
}

func UserCreate(c *gin.Context) {
	var inputs APIUserCreateInputs
	err := c.Bind(&inputs)

	switch {
	case err != nil:
		h.JSONR(c, http.StatusBadRequest, err)
		return
	case utils.HasDangerousCharacters(inputs.CnName):
		h.JSONR(c, http.StatusBadRequest, "name pattern is invalid")
		return
	}
	db := g.Con().Uic
	var user uic.User
	db.Table(user.TableName()).Where(&uic.User{UserName: inputs.UserName}).Scan(&user)
	if user.ID != 0 {
		h.JSONR(c, http.StatusBadRequest, "name is already existing")
		return
	}
	password := utils.HashIt(inputs.Password)
	user = uic.User{
		UserName: inputs.UserName,
		Password: password,
		CnName:   inputs.CnName,
		//Birthday:   inputs.Birthday,
		JgygUserId: inputs.JgygUserId,
		AdUserName: inputs.JgygUserId,
	}

	dt := db.Table(user.TableName()).Create(&user)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, dt.Error)
		return
	}

	var session uic.Session
	response := map[string]string{}
	s := db.Table(session.TableName()).Where("uid = ?", user.ID).Scan(&session)
	if s.Error != nil && s.Error.Error() != "record not found" {
		h.JSONR(c, http.StatusBadRequest, s.Error)
		return
	} else if session.ID == 0 {
		session.Sig = utils.GenerateUUID()
		session.Expired = int(time.Now().Unix()) + 3600*24*30
		session.Uid = int64(user.ID)
		db.Create(&session)
	}
	log.Debugf("%v", session)
	response["sig"] = session.Sig
	response["name"] = user.UserName
	h.JSONR(c, http.StatusOK, response)
	return
}

type APIGetUserInfoOutputs struct {
	User  *uic.User   `json:"user"`
	Roles []*uic.Role `json:"roles"`
}

// @Summary 查看用户信息
// @Description
// @Produce json
// @Param request query string true "查看用户信息"
// @Success 200 {object} APIGetUserInfoOutputs
// @Failure 400 {object} APIGetUserInfoOutputs
// @Router /api/v1/user/info [get]
func UserInfo(c *gin.Context) {
	_, err := h.GetUser(c)
	if err != nil && err.Error() != "record not found" {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}

	id := c.Query("id")
	var user *uic.User
	db := g.Con().Uic
	db = db.Where("id = ?", id).Find(&user)

	var roles []*uic.Role
	dt := g.Con().Portal.Model(uic.Role{})
	dt = dt.Select("`role`.*")
	dt = dt.Joins("left join `user_role_rel` on `user_role_rel`.`role` = `role`.id")
	dt = dt.Where("`user_role_rel`.`user` = ?", id)
	dt = dt.Find(&roles)

	resp := &APIGetUserInfoOutputs{
		User:  user,
		Roles: roles,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary 查看当前用户信息
// @Description
// @Produce json
// @Success 200 {object} uic.User
// @Failure 400 {object} uic.User
// @Router /api/v1/user/myself [get]
func UserMyself(c *gin.Context) {
	user, err := h.GetUser(c)
	if err != nil && err.Error() != "record not found" {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}

	h.JSONR(c, http.StatusOK, user)
	return
}

type APIPutUserUpdateInputs struct {
	ID       int64   `json:"id" form:"id"`
	RoleList []int64 `json:"roleList" form:"roleList"`
}

// @Summary 更新用户信息
// @Description
// @Produce json
// @Param APIPutUserUpdateInputs query APIPutUserUpdateInputs true "根据查询条件分页查询用户列表"
// @Success 200 {object} uic.User
// @Failure 400 {object} uic.User
// @Router /api/v1/user/update [put]
func UserUpdate(c *gin.Context) {
	var inputs APIPutUserUpdateInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	var user *uic.User
	db := g.Con().Uic
	db = db.Where("id = ?", inputs.ID).Find(&user)
	if db.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, db.Error)
		return
	}

	var rels []*uic.UserRoleRel
	dt := g.Con().Portal.Model(uic.UserRoleRel{})
	dt.Where("user = ?", inputs.ID).Delete(&rels)

	for _, v := range inputs.RoleList {
		dt.Create(&uic.UserRoleRel{
			User: inputs.ID,
			Role: v,
		})
	}
	h.JSONR(c, http.StatusOK, user)
	return
}
