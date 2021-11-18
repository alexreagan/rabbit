package uic

type RolePermRel struct {
	Role int64 `json:"role" gorm:"column:role"`
	Perm int64 `json:"perm" gorm:"column:perm"`
}

func (this RolePermRel) TableName() string {
	return "role_perm_rel"
}
