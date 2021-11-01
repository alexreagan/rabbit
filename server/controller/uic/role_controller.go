package uic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit/g"
	h "rabbit/server/helper"
	"rabbit/server/model/uic"
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
