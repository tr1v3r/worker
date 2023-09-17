package worker

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tr1v3r/pkg/log"
	wm "github.com/tr1v3r/workmanager"
)

var _ wm.Worker = new(Pinger)

var PingerBuilder wm.WorkerBuilder = func(ctx context.Context, args map[string]any) wm.Worker { return &Pinger{ctx: ctx, count: 5} }

// Printer print worker
type Pinger struct {
	ctx   context.Context
	count int
}

func (p *Pinger) WithContext(context.Context) wm.Worker { return p }
func (p *Pinger) Work(targets ...wm.WorkTarget) (results []wm.WorkTarget, err error) {
	for _, target := range targets {
		if t, ok := target.(*HostTarget); ok {
			delay, err := p.ping(t.HostName)
			if err != nil {
				return nil, fmt.Errorf("ping host %s fail: %w", t.HostName, err)
			}
			t.ICMPDelay = delay

			results = append(results, t)
		}
	}
	return
}

func (p *Pinger) ping(host string) (time.Duration, error) {
	// 解析主机名
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return 0, fmt.Errorf("resolve ip addr: %w", err)
	}

	// 创建连接
	conn, err := net.DialIP("ip4:icmp", nil, addr)
	if err != nil {
		return 0, fmt.Errorf("dail ip fail: %w", err)
	}
	defer conn.Close()

	cost, reach := time.Duration(0), 0

	// 发送ICMP请求
	for seq := 1; seq <= p.count; seq++ {
		if err := p.ctx.Err(); err != nil {
			return 0, fmt.Errorf("context error: %w", err)
		}

		msg := []byte{8, 0, 0, 0, 0, 0, byte(seq), 0, 0, 0, 0, 0, 0, 0, 0, 0}
		response := make([]byte, 1024)

		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 2)) // 设置读取超时时间

		s := time.Now()
		// 发送请求
		_, _ = conn.Write(msg)
		// 接收响应
		if _, err := conn.Read(response); err != nil {
			log.Warn("Ping %s: request timed out", host)
		} else {
			cost += time.Since(s)
			reach++
			log.Info("Ping %s: received response from %s, cost %s", host, addr.String(), time.Since(s))
		}

		time.Sleep(time.Second) // 休眠一秒钟再发送下一个请求
	}

	if reach == 0 {
		return 0, fmt.Errorf("host %s seems unreachable, tried %d failed", host, p.count)
	}
	return cost / time.Duration(reach), nil
}
