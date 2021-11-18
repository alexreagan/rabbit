package uic

type UserRoleRel struct {
	User int64 `json:"user" gorm:"primary_key;column:user"`
	Role int64 `json:"role" gorm:"primary_key;column:role"`
}

func (this UserRoleRel) TableName() string {
	return "user_role_rel"
}
