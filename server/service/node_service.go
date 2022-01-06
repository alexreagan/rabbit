package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/node"
	"sort"
	"strconv"
	"strings"
)

type nodeService struct {
}

func (s *nodeService) GetAll() []*node.Node {
	var nodes []*node.Node
	tx := g.Con().Portal.Model(node.Node{})
	tx.Find(&nodes)
	return nodes
}

func (s *nodeService) NodesHavingTagIDs(tagIDs []int64) []*node.Node {
	var tIDs []int
	for _, i := range tagIDs {
		tIDs = append(tIDs, int(i))
	}
	sort.Ints(tIDs)

	var tmp []string
	for _, i := range tIDs {
		tmp = append(tmp, strconv.Itoa(i))
	}

	var nodes node.Nodes
	tx := g.Con().Portal.Model(node.Node{})
	tx = tx.Joins("left join `node_tag_rel` on `node`.id=`node_tag_rel`.`node`")
	tx = tx.Where("`node_tag_rel`.`tag` in (?)", tagIDs)
	tx = tx.Group("`node_tag_rel`.`node`")
	tx = tx.Having("group_concat(`node_tag_rel`.`tag` order by `node_tag_rel`.`tag`) = ?", strings.Join(tmp, ","))
	tx.Find(&nodes)
	nodes.Sort()
	return nodes
}

// 已关联了tag的node
func (s *nodeService) NodesRelatedTags() []*node.Node {
	var nodes []*node.Node
	tx := g.Con().Portal.Model(node.Node{})
	tx = tx.Select("distinct `node`.*")
	tx = tx.Joins("left join `node_tag_rel` on `node`.id=`node_tag_rel`.`node`")
	tx = tx.Where("`node_tag_rel`.`tag` is not null")
	tx.Find(&nodes)
	return nodes
}

func newNodeService() *nodeService {
	return &nodeService{}
}
