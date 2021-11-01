package caas

type ServicePortRel struct {
	ServiceID int64 `json:"serviceId" gorm:"column:service_id;index;comment:"`
	PortID    int64 `json:"portId" gorm:"column:port_id;index;comment:"`
}

func (this ServicePortRel) TableName() string {
	return "caas_service_port_rel"
}
