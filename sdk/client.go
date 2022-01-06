package sdk

import (
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/toolkits/net/httplib"
	"time"
)

type YangtzeClient struct {
	Addr string `json:"addr"`
}

const ConnTimeout = 1 * time.Second
const ReadWriteTimeout = 5 * time.Second

// 根据group路径获取路径下所有的node
func (c *YangtzeClient) GetNodesByGroup(path string) ([]*node.Node, error) {
	uri := "/api/v1/node_group/related_nodes"
	req := httplib.Get(c.Addr + uri)
	req = req.SetTimeout(ConnTimeout, ReadWriteTimeout)
	req.Param("path", path)

	var resp []*node.Node
	err := req.ToJson(&resp)
	return resp, err
}

func NewYangtzeClient(addr string) *YangtzeClient {
	return &YangtzeClient{
		Addr: addr,
	}
}
