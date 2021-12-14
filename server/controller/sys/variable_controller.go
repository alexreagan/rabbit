package sys

import (
	"fmt"
	"github.com/alexreagan/rabbit/server/controller/caas"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func EnvList(c *gin.Context) {
	var inputs caas.APIGetEnvListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var variables []*model.APIGetVariableItem
	for key, val := range viper.GetStringMapString(fmt.Sprintf("%s", inputs.Name)) {
		variables = append(variables, &model.APIGetVariableItem{
			Label: val,
			Value: key,
		})
	}

	h.JSONR(c, http.StatusOK, variables)
	return
}
