package portal

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/uic"
)

type Menu struct {
	//gorm.Model
	ID int64 `json:"menuId" gorm:"primary_key;column:id"`
	//Type     int64  `json:"type" gorm:"column:type;comment:"`
	Name     string `json:"name" gorm:"column:name;type:string;size:256;comment:"`
	Url      string `json:"url" gorm:"column:url;type:string;size:512;comment:"`
	ParentId int64  `json:"parentId" gorm:"column:parent_id;comment:"`
	Icon     string `json:"icon" gorm:"column:icon;type:string;size:512;comment:"`
	OrderNum int64  `json:"orderNum" gorm:"column:order_num;comment:"`
	//Open     string `json:"open" gorm:"column:open;type:string;size:128;comment:"`
	List  []*Menu           `json:"list" gorm:"-"`
	Perms []*uic.Permission `json:"perms" gorm:"-"`
}

func (this Menu) TableName() string {
	return "menu"
}

func (this Menu) BuildTree() []*Menu {
	var menuList []*Menu
	var rootMenu []*Menu
	var menuMap map[int64]*Menu
	menuMap = make(map[int64]*Menu)

	db := g.Con().Portal
	db.Table(Menu{}.TableName()).Find(&menuList)
	for _, menu := range menuList {
		menuMap[menu.ID] = menu
		if menu.ParentId == 0 {
			rootMenu = append(rootMenu, menu)
		} else if _, ok := menuMap[menu.ParentId]; ok {
			menuMap[menu.ParentId].List = append(menuMap[menu.ParentId].List, menu)
		}
	}
	return rootMenu
}
