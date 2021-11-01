package portal

type MenuPermission struct {
	Menu       uint `gorm:"column:menu"`
	Permission uint `gorm:"column:permission"`
}

func (this MenuPermission) TableName() string {
	return "menu_permission"
}
