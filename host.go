package worker

import (
	"fmt"
	"net"
	"time"

	wm "github.com/tr1v3r/workmanager"
)

var _ wm.WorkTarget = new(HostTarget)

type HostTarget struct {
	HostMetaInfo
}

func (*HostTarget) Token() string { return "" }
func (*HostTarget) Key() string   { return "" }

var _ fmt.Stringer = new(HostMetaInfo)

type HostMetaInfo struct {
	HostName  string
	IP        net.IP
	ICMPDelay time.Duration
}

func (i *HostMetaInfo) String() string {
	return fmt.Sprintf("hostname: %s,\nip: %s,\nICMP delay: %s", i.HostName, i.IP, i.ICMPDelay)
}
