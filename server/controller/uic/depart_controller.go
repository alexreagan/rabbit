package uic

import (
	"errors"
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/uic"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIGetDepartListsInputs struct {
	Name     string `json:"name" form:"name"`
	Priority int    `json:"priority" form:"priority"`
	Status   string `json:"status" form:"status"`
	//id
	//number of reacord's limit on each page
	Limit int `json:"limit" form:"limit"`
	//pagging
	Page int `json:"page" form:"page"`
}

func (input APIGetDepartListsInputs) checkInputsContain() error {
	return nil
}

func (input APIGetDepartListsInputs) collectDBFilters(database *gorm.DB, tableName string, columns []string) *gorm.DB {
	filterDB := database.Table(tableName)
	// nil columns mean select all columns
	if columns != nil && len(columns) != 0 {
		filterDB = filterDB.Select(columns)
	}
	if input.Name != "" {
		filterDB = filterDB.Where("name like ?", "%"+input.Name+"%")
	}
	return filterDB
}

func DepartmentLists(c *gin.Context) {
	var inputs APIGetDepartListsInputs
	//set default
	inputs.Page = -1
	inputs.Limit = -1
	inputs.Priority = -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	if err := inputs.checkInputsContain(); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	//for get correct table name
	f := uic.Inst{}
	tx := inputs.collectDBFilters(g.Con().Uic, f.TableName(), nil)
	var data []uic.Inst
	//if no specific, will give return first 2000 records
	if inputs.Page == -1 && inputs.Limit == -1 {
		inputs.Limit = 2000
		tx = tx.Order("id DESC").Limit(inputs.Limit)
	} else if inputs.Limit == -1 {
		// set page but not set limit
		h.JSONR(c, h.BadStatus, errors.New("You set page but skip limit params, please check your input"))
		return
	} else {
		// set limit but not set page
		if inputs.Page == -1 {
			// limit invalid
			if inputs.Limit <= 0 {
				h.JSONR(c, h.BadStatus, errors.New("limit or page can not set to 0 or less than 0"))
				return
			}
			// set default page
			inputs.Page = 1
		} else {
			// set page and limit
			// page or limit invalid
			if inputs.Page <= 0 || inputs.Limit <= 0 {
				h.JSONR(c, h.BadStatus, errors.New("limit or page can not set to 0 or less than 0"))
				return
			}
		}
		//set the max limit of each page
		if inputs.Limit >= 50 {
			inputs.Limit = 50
		}
		step := (inputs.Page - 1) * inputs.Limit
		tx = tx.Order("id DESC").Offset(step).Limit(inputs.Limit)
	}
	tx.Find(&data)
	h.JSONR(c, data)
}
