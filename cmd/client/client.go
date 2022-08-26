package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/riverchu/pkg/log"
	ws "github.com/riverchu/pkg/websocket"

	"github.com/riverchu/worker/base"
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

var c *websocket.Conn

func work() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err error
	c, _, err = ws.ConnectWebsocket(ctx, server.String(), nil)
	if err != nil {
		log.Error("connect server websocket fail: %s", err)
		return
	}
	if c == nil {
		log.Error("connect server websocket fail: got nil")
		return
	}
	defer ws.Close(c)

	go func() {
		for msg := range ws.Read(c) {
			log.Info("recv: %s", string(msg))
		}
	}()

	var command string
	var args string
	for {
		fmt.Scanln(&command, &args)
		if command == "" {
			continue
		}

		log.Info("got command: %q with args: %q", command, args)

		if command == "exit" {
			return
		}
		if args != "" {
			executeCmd(&base.CommandMeta{Cmd: command, Args: []string{args}})
		} else {
			executeCmd(&base.CommandMeta{Cmd: command})
		}
	}
}

func executeCmd(cmd *base.CommandMeta) error {
	cmdByte, _ := json.Marshal(cmd)
	info, _ := json.Marshal(&base.Meta{Step: base.StepCommand, Detail: cmdByte})
	return ws.Write(c, info)
}
