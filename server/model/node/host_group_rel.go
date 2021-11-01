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
	"rabbit/g"
)

// +---------+------------------+------+-----+---------+-------+
// | Field   | Type             | Null | Key | Default | Extra |
// +---------+------------------+------+-----+---------+-------+
// | group_id  | int(11) unsigned | NO   | PRI | NULL    |       |
// | host_od   | int(11)          | NO   | PRI | NULL    |       |
// +---------+------------------+------+-----+---------+-------+

type HostGroupRel struct {
	GroupID int64 `gorm:"column:group_id"`
	HostID  int64 `gorm:"column:host_id"`
}

func (this HostGroupRel) TableName() string {
	return "host_group_rel"
}

func (this HostGroupRel) Existing() bool {
	var tGrpHost HostGroupRel
	db := g.Con()
	db.Portal.Table(this.TableName()).Where("group_id = ? AND host_id = ?", this.GroupID, this.HostID).Scan(&tGrpHost)
	if tGrpHost.GroupID != 0 {
		return true
	} else {
		return false
	}
}
