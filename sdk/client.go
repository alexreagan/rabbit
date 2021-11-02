package sdk

import (
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/toolkits/net/httplib"
	"time"
)

type YangtzeClient struct {
	Addr string `json:"addr"`
}

// 根据group路径获取路径下所有的host
func (c *YangtzeClient) GetHostsByGroup(path string) ([]*node.Host, error) {
	uri := "/api/v1/host_group/related_hosts"
	req := httplib.Get(c.Addr + uri)
	req = req.SetTimeout(1*time.Second, 5*time.Second)
	req.Param("path", path)

	var resp []*node.Host
	err := req.ToJson(&resp)
	return resp, err
}
