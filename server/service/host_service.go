package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/node"
	"sort"
	"strconv"
	"strings"
)

type hostService struct {
}

func (s *hostService) GetAll() []*node.Host {
	var hosts []*node.Host
	db := g.Con().Portal.Model(node.Host{})
	db.Find(&hosts)
	return hosts
}

func (s *hostService) HostsHavingTagIDs(tagIDs []int64) []*node.Host {
	var tIDs []int
	for _, i := range tagIDs {
		tIDs = append(tIDs, int(i))
	}
	sort.Ints(tIDs)

	var tmp []string
	for _, i := range tIDs {
		tmp = append(tmp, strconv.Itoa(i))
	}

	var hosts node.Hosts
	db := g.Con().Portal.Model(node.Host{}).Debug()
	db = db.Joins("left join `host_tag_rel` on `host`.id=`host_tag_rel`.`host`")
	db = db.Where("`host_tag_rel`.`tag` in (?)", tagIDs)
	db = db.Group("`host_tag_rel`.`host`")
	db = db.Having("group_concat(`host_tag_rel`.`tag` order by `host_tag_rel`.`tag`) = ?", strings.Join(tmp, ","))
	db.Find(&hosts)
	hosts.Sort()
	return hosts
}

// 已关联了tag的host
func (s *hostService) HostsRelatedTags() []*node.Host {
	var hosts []*node.Host
	db := g.Con().Portal.Model(node.Host{}).Debug()
	db = db.Select("distinct `host`.*")
	db = db.Joins("left join `host_tag_rel` on `host`.id=`host_tag_rel`.`host`")
	db = db.Where("`host_tag_rel`.`tag` is not null")
	db.Find(&hosts)
	return hosts
}

func newHostService() *hostService {
	return &hostService{}
}
