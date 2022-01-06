package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/app"
)

type tagCategoryService struct {
}

func (s *tagCategoryService) GetByName(name string) *app.TagCategory {
	var category app.TagCategory
	tx := g.Con().Portal.Model(category)
	tx.Where("name = ?", name)
	tx.Find(&category)
	return &category
}

func (s *tagCategoryService) GetTagsByCategory(t *app.TagCategory) []*app.Tag {
	var tags []*app.Tag
	var totalCount int64
	tx := g.Con().Portal.Model(app.Tag{})
	tx = tx.Select("`tag`.*, `tag_category`.name as category_name")
	tx = tx.Joins("left join `tag_category` on `tag`.`category_id` = `tag_category`.`id`")
	tx = tx.Where("`tag_category`.id = ?", t.ID)
	tx = tx.Order("`tag`.name")
	tx.Count(&totalCount)
	tx.Find(&tags)
	return tags
}

func (s *tagCategoryService) GetTagsByCategoryID(categoryID int64) []*app.Tag {
	var tags []*app.Tag
	tx := g.Con().Portal.Model(app.Tag{})
	tx = tx.Select("`tag`.*, `tag_category`.name as category_name")
	tx = tx.Joins("left join `tag_category` on `tag_category`.id = `tag`.`category_id`")
	tx = tx.Where("`tag`.category_id = ?", categoryID)
	tx = tx.Order("`tag`.name")
	tx.Find(&tags)
	return tags
}

func newTagCategoryService() *tagCategoryService {
	return &tagCategoryService{}
}
