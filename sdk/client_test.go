package sdk

import (
	"testing"
)

func TestGetHostsByGroup(t *testing.T) {
	c := &YangtzeClient{}
	c.init("http://localhost:8080")
	hosts, err := c.GetHostsByGroup("ROOT/DMP")
	for _, host := range hosts {
		t.Logf("%+v", host)
	}
	t.Logf("%+v", err)
}
