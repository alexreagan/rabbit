package portal

import (
	"fmt"
	"github.com/alexreagan/rabbit/server/controller/node"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type APIGetVariableItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type APIGetVariableInputs struct {
	Name string `json:"name"`
}

type APIGetVariableOutputs struct {
	List       []*APIGetVariableItem `json:"list"`
	TotalCount int64                 `json:"totalCount"`
}

func EnvList(c *gin.Context) {
	var inputs node.APIGetEnvListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var variables []*APIGetVariableItem
	for key, val := range viper.GetStringMapString(fmt.Sprintf("%s", inputs.Name)) {
		variables = append(variables, &APIGetVariableItem{
			Label: val,
			Value: key,
		})
	}

	h.JSONR(c, http.StatusOK, variables)
	return
}
