package g

import "time"

var (
	BinaryName string
	Version    string
	GitCommit  string
)

func VersionMsg() string {
	return Version + "@" + GitCommit
}

const (
	COLLECT_INTERVAL = time.Second
	URL_CHECK_HEALTH = "url.check.health"
	NET_PORT_LISTEN  = "net.port.listen"
	DU_BS            = "du.bs"
	PROC_NUM         = "proc.num"
)
