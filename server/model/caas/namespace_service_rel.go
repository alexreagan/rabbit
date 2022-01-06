package caas

import "github.com/alexreagan/rabbit/g"

type NamespaceServiceRel struct {
	NameSpace int64 `json:"namespace" gorm:"column:namespace;index;comment:"`
	Service   int64 `json:"service" gorm:"column:service;index;comment:"`
}

func (this NamespaceServiceRel) TableName() string {
	return "caas_namespace_service_rel"
}

func (this NamespaceServiceRel) Existing() bool {
	var nsr NamespaceServiceRel
	tx := g.Con().Portal.Model(this)
	tx.Where("namespace = ? and service = ?", this.NameSpace, this.Service).Scan(&nsr)
	if nsr.NameSpace != 0 {
		return true
	} else {
		return false
	}
}
