package node

import (
	"github.com/alexreagan/rabbit/g"
)

type HostTagRel struct {
	Host int64 `gorm:"primary_key;column:host"`
	Tag  int64 `gorm:"primary_key;column:tag"`
}

func (this HostTagRel) TableName() string {
	return "host_tag_rel"
}

func (this HostTagRel) Existing() bool {
	var rel HostTagRel
	db := g.Con()
	db.Portal.Table(this.TableName()).Where("tag = ? AND host = ?", this.Tag, this.Host).Scan(&rel)
	if rel.Tag != 0 {
		return true
	} else {
		return false
	}
}
