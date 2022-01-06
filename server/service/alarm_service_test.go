package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/alarm"
	"github.com/alexreagan/rabbit/server/model/node"
	"testing"
)

func TestTransaction(t *testing.T) {
	tx := g.Con().Portal
	//alarms := make([]*alarm.Alarm, 0)
	//for _, alm := range alarms {
	//	if err := tx.Model(alarm.Alarm{}).Create(&alm).Error; err != nil {
	//		log.Error(err)
	//	}
	//}
	var alarms []*alarm.Alarm
	tx = tx.Model(alarm.Alarm{}).Find(&alarms)

	var nodes []*node.Node
	tx = tx.Model(node.Node{}).Find(&nodes)
}
