package node

import (
	"errors"
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/gtime"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetTagListInputs struct {
	Name       string `json:"name" form:"name"`
	Remark     string `json:"remark" form:"remark"`
	CategoryID string `json:"categoryID" form:"categoryID"`
	Limit      int    `json:"limit" form:"limit"`
	Page       int    `json:"page" form:"page"`
	OrderBy    string `json:"orderBy" form:"orderBy"`
	Order      string `json:"order" form:"order"`
}

type APIGetTagListOutputs struct {
	List       []*node.Tag `json:"list"`
	TotalCount int64       `json:"totalCount"`
}

func (input APIGetTagListInputs) checkInputsContain() error {
	return nil
}

// @Summary tag列表
// @Description
// @Produce json
// @Param APIGetTagListInputs query APIGetTagListInputs true "根据查询条件分页查询tag列表"
// @Success 200 {object} APIGetTagListOutputs
// @Failure 400 {object} APIGetTagListOutputs
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

	var tags []*node.Tag
	var totalCount int64
	db := g.Con().Portal.Debug().Model(node.Tag{})
	db = db.Select("`tag`.*, `tag_category`.name as category_name, `ptag`.`name` as parent_name")
	db = db.Joins("left join `tag_category` on `tag`.`category_id` = `tag_category`.`id`")
	db = db.Joins("left join `tag` as `ptag` on `ptag`.`id` = `tag`.`parent_id`")
	if inputs.Name != "" {
		db = db.Where("`tag`.`name` regexp ?", inputs.Name)
	}
	if inputs.CategoryID != "" {
		db = db.Where("`tag`.`category_id` = ?", inputs.CategoryID)
	}
	if inputs.Remark != "" {
		db = db.Where("`tag`.`remark` regexp ?", inputs.Remark)
	}

	db.Count(&totalCount)
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&tags)

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
// @Success 200 {object} node.TagCate
// @Failure 400 json error
// @Router /api/v1/tag/info [get]
func TagInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		h.JSONR(c, h.BadStatus, errors.New("parameter id is required"))
		return
	}

	var tag *node.Tag
	db := g.Con().Portal.Debug().Model(node.Tag{})
	db.Where("id = ?", id).Find(&tag)

	h.JSONR(c, http.StatusOK, tag)
	return
}

type APIPostTagCreateInputs struct {
	ID         int64  `json:"id" form:"id"`
	Name       string `json:"name" form:"name"`
	CnName     string `json:"cnName" form:"cnName"`
	CategoryID int64  `json:"categoryID" form:"categoryID"`
	//ParentID   int64  `json:"parentID" form:"parentID"`
	Remark string `json:"remark" form:"remark"`
}

// @Summary 创建新tag
// @Description
// @Produce json
// @Param APIPostTagCreateInputs formData APIPostTagCreateInputs true "创建新tag"
// @Success 200 {object} APIPostTagCreateInputs
// @Failure 400 {object} APIPostTagCreateInputs
// @Router /api/v1/tag/create [post]
func TagCreate(c *gin.Context) {
	var inputs APIPostTagCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tag := node.Tag{
		Name:       inputs.Name,
		CnName:     inputs.CnName,
		CategoryID: inputs.CategoryID,
		Remark:     inputs.Remark,
		CreateAt:   gtime.Now(),
	}
	tx := g.Con().Portal
	if dt := tx.Model(node.Tag{}).Create(&tag); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	h.JSONR(c, h.OKStatus, tag)
	return
}

// @Summary 更新tag
// @Description
// @Produce json
// @Param APIPostTagCreateInputs formData APIPostTagCreateInputs true "更新tag"
// @Success 200 {object} APIPostTagCreateInputs
// @Failure 400 {object} APIPostTagCreateInputs
// @Router /api/v1/tag/update [put]
func TagUpdate(c *gin.Context) {
	var inputs APIPostTagCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tag := node.Tag{
		ID:         inputs.ID,
		Name:       inputs.Name,
		CnName:     inputs.CnName,
		CategoryID: inputs.CategoryID,
		//ParentID:   inputs.ParentID,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	tx := g.Con().Portal
	if dt := tx.Model(node.Tag{}).Where("id = ?", inputs.ID).Updates(&tag); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
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
	List       []*node.TagCategory `json:"list"`
	TotalCount int64               `json:"totalCount"`
}

// @Summary tag category列表
// @Description
// @Produce json
// @Param APIGetTagCategoryListInputs query APIGetTagCategoryListInputs true "根据查询条件分页查询tag category列表"
// @Success 200 {object} APIGetTagCategoryListOutputs
// @Failure 400 {object} APIGetTagCategoryListOutputs
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

	var categorys []*node.TagCategory
	var totalCount int64
	db := g.Con().Portal.Debug().Model(node.TagCategory{})
	if inputs.Name != "" {
		db = db.Where("`name` regexp ?", inputs.Name)
	}

	db.Count(&totalCount)
	if inputs.OrderBy != "" {
		db = db.Order(utils.Camel2Case(inputs.OrderBy) + " " + inputs.Order)
	}
	db = db.Offset(offset).Limit(limit)
	db.Find(&categorys)

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
// @Success 200 {object} node.TagCategory
// @Failure 400 json error
// @Router /api/v1/tag_category/info [get]
func TagCategoryInfo(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		h.JSONR(c, h.BadStatus, errors.New("parameter id is required"))
		return
	}

	var tc *node.TagCategory
	db := g.Con().Portal.Debug().Model(node.TagCategory{})
	db.Where("id = ?", id).Find(&tc)

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
// @Param APIPostTagCategoryCreateInputs formData APIPostTagCategoryCreateInputs true "创建tag category"
// @Success 200 {object} node.TagCategory
// @Failure 400 {object} node.TagCategory
// @Router /api/v1/tag_category/create [post]
func TagCategoryCreate(c *gin.Context) {
	var inputs APIPostTagCategoryCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tagCategory := node.TagCategory{
		ID:       inputs.ID,
		Name:     inputs.Name,
		CnName:   inputs.CnName,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	tx := g.Con().Portal
	if dt := tx.Model(node.TagCategory{}).Create(&tagCategory); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	h.JSONR(c, h.OKStatus, tagCategory)
	return
}

// @Summary 更新tag category
// @Description
// @Produce json
// @Param APIPostTagCategoryCreateInputs formData APIPostTagCategoryCreateInputs true "更新tag category"
// @Success 200 {object} node.TagCategory
// @Failure 400 {object} node.TagCategory
// @Router /api/v1/tag_category/update [put]
func TagCategoryUpdate(c *gin.Context) {
	var inputs APIPostTagCategoryCreateInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	tagCategory := node.TagCategory{
		ID:       inputs.ID,
		Name:     inputs.Name,
		CnName:   inputs.CnName,
		Remark:   inputs.Remark,
		CreateAt: gtime.Now(),
	}
	tx := g.Con().Portal
	if dt := tx.Model(node.TagCategory{}).Where("id = ?", inputs.ID).Updates(&tagCategory); dt.Error != nil {
		h.JSONR(c, h.ExpecStatus, dt.Error)
	}

	h.JSONR(c, h.OKStatus, tagCategory)
	return
}

type APIGetTagCategoryTagsOutputs struct {
	List       []*node.Tag `json:"list"`
	TotalCount int64       `json:"totalCount"`
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
// @Failure 400 json error
// @Router /api/v1/tag_category/tags [get]
func TagCategoryTags(c *gin.Context) {
	var inputs APIGetTagCategoryTagsInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var tags []*node.Tag
	var totalCount int64
	db := g.Con().Portal.Model(node.Tag{})
	db = db.Select("`tag`.*, `tag_category`.name as category_name")
	db = db.Joins("left join `tag_category` on `tag`.`category_id` = `tag_category`.`id`")
	if inputs.CategoryID != "" {
		db = db.Where("`tag`.category_id = ?", inputs.CategoryID)
	}
	if inputs.CategoryName != "" {
		db = db.Where("`tag_category`.name = ?", inputs.CategoryName)
	}
	db.Count(&totalCount)
	db.Find(&tags)

	resp := &APIGetTagCategoryTagsOutputs{
		List:       tags,
		TotalCount: totalCount,
	}
	h.JSONR(c, h.OKStatus, resp)
	return
}
