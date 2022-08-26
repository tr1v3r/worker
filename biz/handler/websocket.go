package handler

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/riverchu/pkg/log"

	"github.com/riverchu/worker/base"
	"github.com/riverchu/worker/biz/service/command"
	"github.com/riverchu/worker/biz/service/scan"
	"github.com/riverchu/worker/config"
)

func WSHandle(conn *websocket.Conn, msg []byte) []byte {
	if len(msg) == len("ping") || string(msg) == "pong" {
		return []byte("pong")
	}
	if len(msg) > 0 && msg[0] != '{' {
		return msg
	}

	var meta = new(base.Meta)
	err := json.Unmarshal(msg, meta)
	if err != nil {
		log.Error("unmarshal command fail: %s", err)
		return []byte("unmarshal command fail")
	}

	return HandleMeta(conn, meta)
}

func HandleMeta(conn *websocket.Conn, meta *base.Meta) []byte {
	switch meta.Step {
	case base.StepCommand:
		cmd, err := command.Parse(meta.Detail)
		if err != nil {
			return []byte(err.Error())
		}
		return command.Exec(cmd)
	case base.StepScan:
		target, err := scan.Parse(meta.Detail)
		if err != nil {
			return []byte(err.Error())
		}
		config.SetWSConn(target.TaskToken, conn)
		return []byte(scan.Recv(target).Error())
	default:
		return nil
	}
}
