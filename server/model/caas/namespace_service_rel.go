package caas

import "github.com/alexreagan/rabbit/g"

type NamespaceServiceRel struct {
	NamespaceID int64 `json:"namespaceId" gorm:"column:namespace_id;index;comment:"`
	ServiceID   int64 `json:"serviceId" gorm:"column:service_id;index;comment:"`
}

func (this NamespaceServiceRel) TableName() string {
	return "caas_namespace_service_rel"
}

func (this NamespaceServiceRel) Existing() bool {
	var nsr NamespaceServiceRel
	db := g.Con().Portal
	db.Model(this).Where("namespace_id = ? and service_id = ?", this.NamespaceID, this.ServiceID).Scan(&nsr)
	if nsr.NamespaceID != 0 {
		return true
	} else {
		return false
	}
}
