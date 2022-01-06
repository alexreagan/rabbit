package alarm

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type Alarm struct {
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

func (this Alarm) TableName() string {
	return "alarm"
}

func (this Alarm) LatestRecords() []*Alarm {
	var alarms []*Alarm
	tx := g.Con().Portal.Debug()
	tx = tx.Model(Alarm{})
	tx = tx.Select("`alarm`.*")
	tx = tx.Joins("right join (select max(id) as id from alarm group by prod_ip) as tbl on alarm.id = tbl.id")
	tx = tx.Find(&alarms)
	return alarms
}
