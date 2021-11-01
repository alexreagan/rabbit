package uic

type RolePermission struct {
	ID         int64 `json:"id" gorm:"primary_key;column:id"`
	Role       int64 `gorm:"column:role"`
	Permission int64 `gorm:"column:permission"`
}

func (this RolePermission) TableName() string {
	return "role_perm_rel"
}
