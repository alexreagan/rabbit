package portal

import (
	"fmt"
	"github.com/alexreagan/rabbit/server/controller/node"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func EnvList(c *gin.Context) {
	var inputs node.APIGetEnvListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var variables []*h.APIGetVariableItem
	for key, val := range viper.GetStringMapString(fmt.Sprintf("%s", inputs.Name)) {
		variables = append(variables, &h.APIGetVariableItem{
			Label: val,
			Value: key,
		})
	}

	h.JSONR(c, http.StatusOK, variables)
	return
}
