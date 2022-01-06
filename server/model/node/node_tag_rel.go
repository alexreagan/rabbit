package node

import (
	"github.com/alexreagan/rabbit/g"
)

type NodeTagRel struct {
	Node int64 `gorm:"primary_key;column:node"`
	Tag  int64 `gorm:"primary_key;column:tag"`
}

func (this NodeTagRel) TableName() string {
	return "node_tag_rel"
}

func (this NodeTagRel) Existing() bool {
	var rel NodeTagRel
	tx := g.Con().Portal.Model(this)
	tx.Where("tag = ? AND node = ?", this.Tag, this.Node).Scan(&rel)
	if rel.Tag != 0 {
		return true
	} else {
		return false
	}
}
