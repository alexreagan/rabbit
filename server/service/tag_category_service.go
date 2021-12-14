package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/app"
)

type tagCategoryService struct {
}

func (s *tagCategoryService) GetByName(name string) *app.TagCategory {
	var category app.TagCategory
	db := g.Con().Portal.Model(category)
	db.Where("name = ?", name)
	db.Find(&category)
	return &category
}

func (s *tagCategoryService) GetTagsByCategory(t *app.TagCategory) []*app.Tag {
	var tags []*app.Tag
	var totalCount int64
	db := g.Con().Portal.Model(app.Tag{})
	db = db.Select("`tag`.*, `tag_category`.name as category_name")
	db = db.Joins("left join `tag_category` on `tag`.`category_id` = `tag_category`.`id`")
	db = db.Where("`tag_category`.id = ?", t.ID)
	db = db.Order("`tag`.name")
	db.Count(&totalCount)
	db.Find(&tags)
	return tags
}

func (s *tagCategoryService) GetTagsByCategoryID(categoryID int64) []*app.Tag {
	var tags []*app.Tag
	db := g.Con().Portal.Model(app.Tag{})
	db = db.Select("`tag`.*, `tag_category`.name as category_name")
	db = db.Joins("left join `tag_category` on `tag_category`.id = `tag`.`category_id`")
	db = db.Where("`tag`.category_id = ?", categoryID)
	db = db.Order("`tag`.name")
	db.Find(&tags)
	return tags
}

func newTagCategoryService() *tagCategoryService {
	return &tagCategoryService{}
}
