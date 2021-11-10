package node

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type Alert struct {
	ID            int64       `json:"id" gorm:"primary_key;column:id"`
	AlertLevel    string      `json:"alertLevel" gorm:"column:alert_level;type:string;size:128;comment:"`
	AlertName     string      `json:"alertName" gorm:"column:alert_name;type:string;size:128;comment:"`
	CloudPoolName string      `json:"cloudPoolName" gorm:"column:cloud_pool_name;type:string;size:128;comment:"`
	FiringTime    gtime.GTime `json:"firingTime" gorm:"column:firing_time;default:null;comment:"`
	ProdIP        string      `json:"prodIp" gorm:"column:prod_ip;type:string;size:128;comment:"`
	Resolved      bool        `json:"resolved" `
	ResolvedTime  gtime.GTime `json:"resolvedTime" gorm:"column:resolved_time;default:null;comment:"`
	StrategyID    int64       `json:"strategyId" gorm:"column:strategy_id;comment:"`
	StrategyName  string      `json:"strategyName" gorm:"column:strategy_name;type:string;size:128;comment:"`
	StrategyType  string      `json:"strategyType" gorm:"column:strategy_type;type:string;size:128;comment:"`
	SubSysEnName  string      `json:"subSysEnName" gorm:"column:sub_sys_en_name;type:string;size:128;comment:"`
	SubSysName    string      `json:"subSysName" gorm:"column:sub_sys_name;type:string;size:128;comment:"`
	U1            string      `json:"u1" gorm:"column:u1;type:string;size:128;comment:"`
	U2            string      `json:"u2" gorm:"column:u2;type:string;size:128;comment:"`
	UpdateTime    gtime.GTime `json:"updateTime" gorm:"column:update_time;default:null;comment:"`
}

func (this Alert) TableName() string {
	return "alert"
}

func (this Alert) LatestRecords() []*Alert {
	var alerts []*Alert
	db := g.Con().Portal.Debug()
	db = db.Model(Alert{})
	db = db.Select("`alert`.*")
	db = db.Joins("right join (select max(id) as id from alert group by prod_ip) as tbl on alert.id = tbl.id")
	db = db.Find(&alerts)
	return alerts
}
