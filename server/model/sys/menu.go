package sys

import (
	"github.com/alexreagan/rabbit/g"
)

type Menu struct {
	//gorm.Model
	ID int64 `json:"menuId" gorm:"primary_key;column:id"`
	//Type     int64  `json:"type" gorm:"column:type;comment:"`
	Name     string `json:"name" gorm:"column:name;type:string;size:256;comment:"`
	Url      string `json:"url" gorm:"column:url;type:string;size:512;comment:"`
	ParentID int64  `json:"parentID" gorm:"column:parent_id;comment:"`
	Icon     string `json:"icon" gorm:"column:icon;type:string;size:512;comment:"`
	OrderNum int64  `json:"orderNum" gorm:"column:order_num;comment:"`
	//Open     string `json:"open" gorm:"column:open;type:string;size:128;comment:"`
	Children []*Menu `json:"list" gorm:"-"`
	//Perms []*uic.Perm `json:"perms" gorm:"-"`
}

func (this Menu) TableName() string {
	return "menu"
}

func (this Menu) BuildTree() []*Menu {
	var menus []*Menu
	var rootMenu []*Menu
	var menuMap map[int64]*Menu
	menuMap = make(map[int64]*Menu)

	db := g.Con().Portal
	db.Model(Menu{}).Order("order_num").Find(&menus)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}
	for _, menu := range menus {
		if menu.ParentID == 0 {
			rootMenu = append(rootMenu, menu)
		} else if _, ok := menuMap[menu.ParentID]; ok {
			menuMap[menu.ParentID].Children = append(menuMap[menu.ParentID].Children, menu)
		}
	}
	return rootMenu
}
