package caas

type ServiceTagRel struct {
	Service int64 `json:"service" gorm:"column:service;index;comment:"`
	Tag     int64 `json:"tag" gorm:"column:tag;index;comment:"`
}

func (this ServiceTagRel) TableName() string {
	return "caas_service_tag_rel"
}
