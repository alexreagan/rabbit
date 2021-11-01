package uic

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"rabbit/g"
	h "rabbit/server/helper"
	"rabbit/server/model/uic"
	"rabbit/server/utils"
	"time"
)

type APILoginInput struct {
	UserName string `json:"username"  form:"username" binding:"required"`
	Password string `json:"password"  form:"password" binding:"required"`
}

type APIAdminLoginInput struct {
	UserName string `json:"username"  form:"username" binding:"required"`
}

func Login(c *gin.Context) {
	inputs := APILoginInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, "name or password is blank")
		return
	}

	db := g.Con().Uic
	user := uic.User{}
	db.Table(user.TableName()).Where(uic.User{AdUserName: inputs.UserName}).Find(&user)
	if user.ID == 0 {
		h.JSONR(c, h.BadStatus, "no such user")
		return
	}

	switch viper.GetBool("ldap.enabled") {
	case true:
		err := utils.Authenticate(inputs.UserName, inputs.Password)
		if err != nil {
			h.JSONR(c, h.BadStatus, err.Error())
			return
		}
	default:
		if user.Password != utils.HashIt(inputs.Password) {
			h.JSONR(c, h.BadStatus, "password error")
			return
		}
	}

	var session uic.Session
	s := db.Table(session.TableName()).Where("uid = ?", user.ID).Scan(&session)
	if s.Error != nil && s.Error.Error() != "record not found" {
		h.JSONR(c, h.BadStatus, s.Error)
		return
	} else if session.ID == 0 {
		session.Sig = utils.GenerateUUID()
		session.Expired = int(time.Now().Unix()) + 3600*24*30
		session.Uid = int64(user.ID)
		db.Create(&session)
	}
	log.Debugf("login session: %v", session)
	resp := struct {
		Sig   string `json:"sig,omitempty"`
		Name  string `json:"name,omitempty"`
		Admin bool   `json:"admin"`
	}{session.Sig, user.UserName, user.IsAdmin()}
	h.JSONR(c, resp)
	return
}

func AdminLogin(c *gin.Context) {
	inputs := APIAdminLoginInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, "name is blank")
		return
	}
	name := inputs.UserName

	user := uic.User{
		UserName: name,
	}
	_, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}

	db := g.Con().Uic
	db.Where(&user).Find(&user)
	switch {
	case user.ID == 0:
		h.JSONR(c, h.BadStatus, "no such user")
		return
	case user.IsSuperuser == false:
		h.JSONR(c, h.BadStatus, "API_USER not admin, no permissions can do this")
		return
	}
	var session uic.Session
	s := db.Table(session.TableName()).Where("uid = ?", user.ID).Scan(&session)
	if s.Error != nil && s.Error.Error() != "record not found" {
		h.JSONR(c, h.BadStatus, s.Error)
		return
	} else if session.ID == 0 {
		session.Sig = utils.GenerateUUID()
		session.Expired = int(time.Now().Unix()) + 3600*24*30
		session.Uid = int64(user.ID)
		db.Create(&session)
	}
	log.Debugf("session: %v", session)
	resp := struct {
		Sig   string `json:"sig,omitempty"`
		Name  string `json:"name,omitempty"`
		Admin bool   `json:"admin"`
	}{session.Sig, user.UserName, user.IsAdmin()}
	h.JSONR(c, resp)
	return
}

func Logout(c *gin.Context) {
	wsession, err := h.GetSession(c)
	if err != nil {
		h.JSONR(c, h.BadStatus, err.Error())
		return
	}
	var session uic.Session
	var user uic.User
	db := g.Con().Uic
	db.Table(user.TableName()).Where(uic.User{UserName: wsession.Name}).Scan(&user)
	db.Table(session.TableName()).Where("sig = ? AND uid = ?", wsession.Sig, user.ID).Scan(&session)

	if session.ID == 0 {
		h.JSONR(c, h.BadStatus, "not found this kind of session in database.")
		return
	} else {
		r := db.Table(session.TableName()).Delete(&session)
		if r.Error != nil {
			h.JSONR(c, h.BadStatus, r.Error)
		}
		h.JSONR(c, h.OKStatus, "logout successful")
	}
	return
}

func AuthSession(c *gin.Context) {
	auth, err := h.SessionChecking(c)
	if err != nil || auth != true {
		h.JSONR(c, http.StatusUnauthorized, err)
		return
	}
	h.JSONR(c, "session is valid!")
	return
}

func CreateRoot(c *gin.Context) {
	password := c.DefaultQuery("password", "")
	if password == "" {
		h.JSONR(c, h.BadStatus, "password is empty, please check it")
		return
	}
	password = utils.HashIt(password)
	user := uic.User{
		UserName: "root",
		Password: password,
	}
	db := g.Con().Uic
	dt := db.Table(user.TableName()).Save(&user)
	if dt.Error != nil {
		h.JSONR(c, h.BadStatus, dt.Error)
		return
	}
	h.JSONR(c, "root created!")
	return
}
