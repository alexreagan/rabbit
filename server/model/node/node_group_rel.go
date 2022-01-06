// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"github.com/alexreagan/rabbit/g"
)

// +---------+------------------+------+-----+---------+-------+
// | Field   | Type             | Null | Key | Default | Extra |
// +---------+------------------+------+-----+---------+-------+
// | group_id  | int(11) unsigned | NO   | PRI | NULL    |       |
// | node_id   | int(11)          | NO   | PRI | NULL    |       |
// +---------+------------------+------+-----+---------+-------+

type NodeGroupRel struct {
	GroupID int64 `gorm:"column:group_id"`
	NodeID  int64 `gorm:"column:node_id"`
}

func (this NodeGroupRel) TableName() string {
	return "node_group_rel"
}

func (this NodeGroupRel) Existing() bool {
	var tGrpNode NodeGroupRel
	tx := g.Con().Portal.Model(this)
	tx.Where("group_id = ? AND node_id = ?", this.GroupID, this.NodeID).Scan(&tGrpNode)
	if tGrpNode.GroupID != 0 {
		return true
	} else {
		return false
	}
}
