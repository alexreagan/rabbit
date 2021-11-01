package caas

type ServicePodRel struct {
	ServiceID int64 `json:"serviceId" gorm:"column:service_id;index;comment:"`
	PodID     int64 `json:"podId" gorm:"column:pod_id;index;comment:"`
}

func (this ServicePodRel) TableName() string {
	return "caas_service_pod_rel"
}
