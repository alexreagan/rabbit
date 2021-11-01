package uic

type Permission struct {
	//gorm.Model
	ID     int64  `json:"id" gorm:"primary_key;column:id"`
	Name   string `json:"name" gorm:"column:name;type:string;size:80;index;not null;comment:权限名称"`
	CnName string `json:"cnName" gorm:"column:cn_name;type:string;size:80;comment:权限中文名称"`
	Desc   string `json:"desc" gorm:"column:description;type:string;size:255;comment:权限描述"`
}

func (r Permission) TableName() string {
	return "permission"
}
