package sdk

import (
	"testing"
)

func TestGetNodesByGroup(t *testing.T) {
	c := &YangtzeClient{
		Addr: "http://localhost:8080",
	}
	nodes, err := c.GetNodesByGroup("ROOT/DMP")
	for _, n := range nodes {
		t.Logf("%+v", n)
	}
	t.Logf("%+v", err)
}
