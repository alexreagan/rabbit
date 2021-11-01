package uic

type Role struct {
	//gorm.Model
	ID     int64  `json:"id" gorm:"primary_key;column:id"`
	Name   string `json:"name" gorm:"column:name;type:string;size:80;index;not null;comment:角色名称"`
	CnName string `json:"cnName" gorm:"column:cn_name;type:string;size:80;comment:中文名称"`
	Desc   string `json:"desc" gorm:"column:description;type:string;size:255;comment:角色描述"`
}

func (r Role) TableName() string {
	return "role"
}
