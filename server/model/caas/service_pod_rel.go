package caas

type ServicePodRel struct {
	Service int64 `json:"service" gorm:"column:service;index;comment:"`
	Pod     int64 `json:"pod" gorm:"column:pod;index;comment:"`
}

func (this ServicePodRel) TableName() string {
	return "caas_service_pod_rel"
}
