package caas

type ServicePortRel struct {
	Service int64 `json:"service" gorm:"column:service;index;comment:"`
	Port    int64 `json:"port" gorm:"column:port;index;comment:"`
}

func (this ServicePortRel) TableName() string {
	return "caas_service_port_rel"
}
