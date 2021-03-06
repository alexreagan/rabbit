package app

import (
	"errors"
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetTagListInputs struct {
	Name         string `json:"name" form:"name"`
	Remark       string `json:"remark" form:"remark"`
	CategoryID   string `json:"categoryID" form:"categoryID"`
	CategoryName string `json:"categoryName" form:"categoryName"`
	Limit        int    `json:"limit" form:"limit"`
	Page         int    `json:"page" form:"page"`
	OrderBy      string `json:"orderBy" form:"orderBy"`
	Order        string `json:"order" form:"order"`
}

type APIGetTagListOutputs struct {
	List       []*app.Tag `json:"list"`
	TotalCount int64      `json:"totalCount"`
}

func (input APIGetTagListInputs) checkInputsContain() error {
	return nil
}

type APIGetTagAllInputs struct {
	Name         string `json:"name" form:"name"`
	Remark       string `json:"remark" form:"remark"`
	CategoryName string `json:"categoryName" form:"categoryName"`
	OrderBy      string `json:"orderBy" form:"orderBy"`
	Order        string `json:"order" form:"order"`
}

// @Summary 全部tag数据
// @Description
// @Produce json
// @Param APIGetTagAllInputs query APIGetTagAllInputs true "查询全部tag列表"
// @Success 200 {object} APIGetTagListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag/all [get]
func TagAll(c *gin.Context) {
	var inputs APIGetTagAllInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var tags []*app.Tag
	var totalCount int64
	tx := g.Con().Portal.Model(app.Tag{})
	if inputs.Name != "" {
		tx = tx.Where("`tag`.`name` regexp ?", inputs.Name)
	}
	if inputs.CategoryName != "" {
		tx = tx.Where("`tag`.`category_name` = ?", inputs.CategoryName)
	}
	if inputs.Remark != "" {
		tx = tx.Where("`tag`.`remark` regexp ?", inputs.Remark)
	}
	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	} else {
		tx = tx.Order("`tag`.`name`")
	}
	tx.Find(&tags)

	resp := &APIGetTagListOutputs{
		List:       tags,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary tag列表
// @Description
// @Produce json
// @Param APIGetTagListInputs query APIGetTagListInputs true "根据查询条件分页查询tag列表"
// @Success 200 {object} APIGetTagListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag/list [get]
func TagList(c *gin.Context) {
	var inputs APIGetTagListInputs

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

	var tags []*app.Tag
	var totalCount int64
	tx := g.Con().Portal.Model(app.Tag{})
	if inputs.Name != "" {
		tx = tx.Where("`tag`.`name` regexp ?", inputs.Name)
	}
	if inputs.CategoryID != "" {
		tx = tx.Where("`tag`.`category_id` = ?", inputs.CategoryID)
	}
	if inputs.CategoryName != "" {
		tx = tx.Where("`tag`.`category_name` = ?", inputs.CategoryName)
	}
	if inputs.Remark != "" {
		tx = tx.Where("`tag`.`remark` regexp ?", inputs.Remark)
	}

	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	} else {
		tx = tx.Order("`tag`.`name`")
	}
	tx = tx.Offset(offset).Limit(limit)
	tx.Find(&tags)

	resp := &APIGetTagListOutputs{
		List:       tags,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary tag 详情
// @Description
// @Produce json
// @Param id query string true "tag id"
// @Success 200 {object} app.G6Edge
// @Failure 400 "bad arguments"
// @Router /api/v1/tag/info [get]
func TagInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		h.JSONR(c, h.BadStatus, errors.New("parameter id is required"))
		return
	}

	var tag *app.Tag
	tx := g.Con().Portal.Debug().Model(app.Tag{})
	tx.Where("id = ?", id).Find(&tag)

	h.JSONR(c, http.StatusOK, tag)
	return
}

type APIPostTagCreateInputs struct {
	ID         int64  `json:"id" form:"id"`
	Name       string `json:"name" form:"name"`
	CnName     string `json:"cnName" form:"cnName"`
	CategoryID int64  `json:"categoryID" form:"categoryID"`
	Remark     string `json:"remark" form:"remark"`
}

// @Summary 创建新tag
// @Description
// @Produce json
// @Param APIPostTagCreateInputs body APIPostTagCreateInputs true "创建新tag"
// @Success 200 {object} APIPostTagCreateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag/create [post]
func TagCreate(c *gin.Context) {
	var inputs APIPostTagCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var tagCategory *app.TagCategory
	g.Con().Portal.Model(app.TagCategory{}).Where("id = ?", inputs.CategoryID).Find(&tagCategory)

	tag := app.Tag{
		Name:         inputs.Name,
		CnName:       inputs.CnName,
		CategoryID:   inputs.CategoryID,
		CategoryName: tagCategory.Name,
		Remark:       inputs.Remark,
		CreateAt:     gtime.Now(),
	}
	tx := g.Con().Portal
	if tx = tx.Model(app.Tag{}).Create(&tag); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
	}

	h.JSONR(c, h.OKStatus, tag)
	return
}

// @Summary 更新tag
// @Description
// @Produce json
// @Param APIPostTagCreateInputs body APIPostTagCreateInputs true "更新tag"
// @Success 200 {object} APIPostTagCreateInputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag/update [put]
func TagUpdate(c *gin.Context) {
	var inputs APIPostTagCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var tagCategory *app.TagCategory
	g.Con().Portal.Model(app.TagCategory{}).Where("id = ?", inputs.CategoryID).Find(&tagCategory)

	tag := app.Tag{
		ID:           inputs.ID,
		Name:         inputs.Name,
		CnName:       inputs.CnName,
		CategoryID:   inputs.CategoryID,
		CategoryName: tagCategory.Name,
		Remark:       inputs.Remark,
		CreateAt:     gtime.Now(),
	}
	tx := g.Con().Portal
	if tx = tx.Model(app.Tag{}).Where("id = ?", inputs.ID).Updates(&tag); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
	}

	h.JSONR(c, h.OKStatus, tag)
	return
}

type APIGetTagCategoryListInputs struct {
	Name    string `json:"name" form:"name"`
	Limit   int    `json:"limit" form:"limit"`
	Page    int    `json:"page" form:"page"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}

type APIGetTagCategoryListOutputs struct {
	List       []*app.TagCategory `json:"list"`
	TotalCount int64              `json:"totalCount"`
}

// @Summary tag category列表
// @Description
// @Produce json
// @Param APIGetTagCategoryListInputs query APIGetTagCategoryListInputs true "根据查询条件分页查询tag category列表"
// @Success 200 {object} APIGetTagCategoryListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag_category/all [get]
func TagCategoryAll(c *gin.Context) {
	var inputs APIGetTagCategoryListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var categorys []*app.TagCategory
	var totalCount int64
	tx := g.Con().Portal.Debug().Model(app.TagCategory{})
	if inputs.Name != "" {
		tx = tx.Where("`name` regexp ?", inputs.Name)
	}

	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx.Find(&categorys)

	resp := &APIGetTagCategoryListOutputs{
		List:       categorys,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary tag category列表
// @Description
// @Produce json
// @Param APIGetTagCategoryListInputs query APIGetTagCategoryListInputs true "根据查询条件分页查询tag category列表"
// @Success 200 {object} APIGetTagCategoryListOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag_category/list [get]
func TagCategoryList(c *gin.Context) {
	var inputs APIGetTagCategoryListInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}
	offset, limit, err := h.PageParser(inputs.Page, inputs.Limit)
	if err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var categorys []*app.TagCategory
	var totalCount int64
	tx := g.Con().Portal.Debug().Model(app.TagCategory{})
	if inputs.Name != "" {
		tx = tx.Where("`name` regexp ?", inputs.Name)
	}

	tx.Count(&totalCount)
	if inputs.OrderBy != "" {
		tx = tx.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	tx = tx.Offset(offset).Limit(limit)
	tx.Find(&categorys)

	resp := &APIGetTagCategoryListOutputs{
		List:       categorys,
		TotalCount: totalCount,
	}
	h.JSONR(c, http.StatusOK, resp)
	return
}

// @Summary tag category 详情
// @Description
// @Produce json
// @Param id query string true "tag category id"
// @Success 200 {object} app.TagCategory
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag_category/info [get]
func TagCategoryInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		h.JSONR(c, h.BadStatus, errors.New("parameter id is required"))
		return
	}

	var tc *app.TagCategory
	tx := g.Con().Portal.Debug().Model(app.TagCategory{})
	tx.Where("id = ?", id).Find(&tc)

	h.JSONR(c, http.StatusOK, tc)
	return
}

type APIPostTagCategoryCreateInputs struct {
	ID         int64  `json:"id" form:"id"`
	Name       string `json:"name" form:"name"`
	CnName     string `json:"cnName" form:"cnName"`
	CategoryID int64  `json:"categoryID" form:"categoryID"`
	Remark     string `json:"remark" form:"remark"`
}

// @Summary 创建tag category
// @Description
// @Produce json
// @Param APIPostTagCategoryCreateInputs body APIPostTagCategoryCreateInputs true "创建tag category"
// @Success 200 {object} app.TagCategory
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag_category/create [post]
func TagCategoryCreate(c *gin.Context) {
	var inputs APIPostTagCategoryCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tagCategory := app.TagCategory{
		ID:       inputs.ID,
		Name:     inputs.Name,
		CnName:   inputs.CnName,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	tx := g.Con().Portal
	if tx = tx.Model(app.TagCategory{}).Create(&tagCategory); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
	}

	h.JSONR(c, h.OKStatus, tagCategory)
	return
}

// @Summary 更新tag category
// @Description
// @Produce json
// @Param APIPostTagCategoryCreateInputs body APIPostTagCategoryCreateInputs true "更新tag category"
// @Success 200 {object} app.TagCategory
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag_category/update [put]
func TagCategoryUpdate(c *gin.Context) {
	var inputs APIPostTagCategoryCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tagCategory := app.TagCategory{
		ID:       inputs.ID,
		Name:     inputs.Name,
		CnName:   inputs.CnName,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	tx := g.Con().Portal
	if tx = tx.Model(app.TagCategory{}).Where("id = ?", inputs.ID).Updates(&tagCategory); tx.Error != nil {
		h.JSONR(c, h.ExpecStatus, tx.Error)
	}

	h.JSONR(c, h.OKStatus, tagCategory)
	return
}

type APIGetTagCategoryTagsOutputs struct {
	List       []*app.Tag `json:"list"`
	TotalCount int64      `json:"totalCount"`
}

type APIGetTagCategoryTagsInputs struct {
	CategoryID   string `json:"categoryID" form:"categoryID"`
	CategoryName string `json:"categoryName" form:"categoryName"`
}

// @Summary 查看某一tag category下的tag
// @Description
// @Produce json
// @Param id query APIGetTagCategoryTagsInputs true "查看tag category下的tag"
// @Success 200 {object} APIGetTagCategoryTagsOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tag_category/tags [get]
func TagCategoryTags(c *gin.Context) {
	var inputs APIGetTagCategoryTagsInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var tags []*app.Tag
	var totalCount int64
	tx := g.Con().Portal.Model(app.Tag{})
	tx = tx.Select("`tag`.*, `tag_category`.name as category_name")
	tx = tx.Joins("left join `tag_category` on `tag`.`category_id` = `tag_category`.`id`")
	if inputs.CategoryID != "" {
		tx = tx.Where("`tag`.category_id = ?", inputs.CategoryID)
	}
	if inputs.CategoryName != "" {
		tx = tx.Where("`tag_category`.name = ?", inputs.CategoryName)
	}
	tx.Count(&totalCount)
	tx.Find(&tags)

	resp := &APIGetTagCategoryTagsOutputs{
		List:       tags,
		TotalCount: totalCount,
	}
	h.JSONR(c, h.OKStatus, resp)
	return
}
