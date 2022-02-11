package sys

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/sys"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetNoticeListInputs struct {
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}

type APIGetNoticeListOutputs struct {
	List       []*sys.Notice `json:"list"`
	TotalCount int64         `json:"totalCount"`
}

// @Summary 系统公告
// @Description
// @Produce json
// @Param APIGetNoticeListInputs body APIGetNoticeListOutputs true "根据查询条件查询系统公告"
// @Success 200 {object} APIGetNoticeListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/notice/list [get]
func NoticeList(c *gin.Context) {
	var inputs APIGetNoticeListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var notices []*sys.Notice
	var totalCount int64
	tx := g.Con().Portal.Model(sys.Notice{})

	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx = tx.Offset(offset).Limit(limit)
	tx.Find(&notices)

	resp := &APIGetNoticeListOutputs{
		List:       notices,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIPostNoticeCreateInputs struct {
	ID      int64         `json:"id" form:"id"`
	Title   string        `json:"title" form:"title"`
	Content string        `json:"content" form:"content"`
	Time    []gtime.GTime `json:"time" form:"time"`
}

type APIPostNoticeCreateOutputs struct {
	Notice *sys.Notice `json:"notice"`
}

// @Summary 系统公告创建
// @Description
// @Produce json
// @Param APIPostNoticeCreateInputs body sys.Notice true "系统公告创建"
// @Success 200 {object} sys.Notice
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/notice/create [post]
func NoticeCreate(c *gin.Context) {
	var inputs APIPostNoticeCreateInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	u, _ := h.GetUser(c)

	notice := &sys.Notice{
		Title:       inputs.Title,
		Content:     inputs.Content,
		TimeBegin:   inputs.Time[0],
		TimeEnd:     inputs.Time[1],
		Creator:     u.JgygUserID,
		CreatorName: u.CnName,
		CreateAt:    gtime.Now(),
	}
	tx := g.Con().Portal.Model(notice)
	if err := tx.Create(notice).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, http.StatusOK, notice)
	return
}

// @Summary 系统公告更新
// @Description
// @Produce json
// @Param APIPostNoticeCreateInputs body sys.Notice true "系统公告更新"
// @Success 200 {object} sys.Notice
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/notice/update [put]
func NoticeUpdate(c *gin.Context) {
	var inputs APIPostNoticeCreateInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	u, _ := h.GetUser(c)

	notice := &sys.Notice{
		Title:       inputs.Title,
		Content:     inputs.Content,
		TimeBegin:   inputs.Time[0],
		TimeEnd:     inputs.Time[1],
		Creator:     u.JgygUserID,
		CreatorName: u.CnName,
		CreateAt:    gtime.Now(),
	}
	tx := g.Con().Portal.Model(notice)
	if err := tx.Where("id = ?", inputs.ID).Updates(notice).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, http.StatusOK, notice)
	return
}

type APIGetNoticeInfoInputs struct {
	ID int64 `json:"id" form:"id"`
}

type APIGetNoticeInfoOutputs struct {
	Notice *sys.Notice `json:"notice"`
}

// @Summary 系统公告信息
// @Description
// @Produce json
// @Param APIGetNoticeInfoInputs body sys.Notice true "系统公告信息"
// @Success 200 {object} sys.Notice
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/notice/info [get]
func NoticeInfo(c *gin.Context) {
	var inputs APIGetNoticeInfoInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	var notice sys.Notice
	tx := g.Con().Portal.Model(&notice)
	if err := tx.Where("id = ?", inputs.ID).First(&notice).Error; err != nil {
		h.JSONR(c, h.ExpecStatus, err)
		return
	}

	h.JSONR(c, http.StatusOK, notice)
	return
}
