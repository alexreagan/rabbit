package uic

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/uic"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetRoleSearchInputs struct {
	Name string `json:"name"`
}

type APIGetRoleSearchOutputs struct {
	List []*uic.Role `json:"list"`
}

func RoleSelect(c *gin.Context) {
	var inputs APIGetRoleSearchInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	db := g.Con().Uic
	var roles []*uic.Role
	db.Table(uic.Role{}.TableName()).Where("name like ?", "%"+inputs.Name+"%").Find(&roles)

	resp := &APIGetRoleSearchOutputs{
		List: roles,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}
