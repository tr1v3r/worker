package main

import (
	"context"
	"net/url"
	"time"

	"github.com/riverchu/pkg/log"
	ws "github.com/riverchu/pkg/websocket"

	"github.com/riverchu/worker/config"
)

const (
	serverScheme = "ws"
	serverAddr   = "localhost"
	serverPath   = "/ws"
)

var server = &url.URL{Scheme: serverScheme, Host: serverAddr + ":" + config.WebServerPort(), Path: serverPath}

func main() {
	work()
}

func work() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	c, _, err := ws.ConnectWebsocket(ctx, server.String(), nil)
	if err != nil {
		log.Error("connect server websocket fail: %s", err)
		return
	}
	if c == nil {
		log.Error("connect server websocket fail: got nil")
		return
	}
	defer ws.Close(c)

	_ = ws.Write(c, []byte("ping"))
	_ = ws.Write(c, []byte(`{"cmd": "pwd"}`))
	_ = ws.Write(c, []byte(`{"cmd": "whoami"}`))

	go func() {
		for ts := range time.Tick(time.Second) {
			err := ws.Write(c, []byte(ts.String()))
			if err != nil {
				log.Error("write msg fail: %s", err)
			}
		}
	}()

	for msg := range ws.Read(c) {
		log.Info("recv: %s", string(msg))
	}
}
