package handler

import (
	"encoding/json"
	"os/exec"

	"github.com/riverchu/pkg/log"
	"github.com/riverchu/worker/base"
)

func WSHandle(msg []byte) []byte {
	if len(msg) == len("ping") || string(msg) == "pong" {
		return []byte("pong")
	}
	if len(msg) > 0 && msg[0] != '{' {
		return msg
	}

	var cmd base.Command
	err := json.Unmarshal(msg, &cmd)
	if err != nil {
		log.Error("unmarshal command fail: %s", err)
		return []byte("unmarshal command fail")
	}

	c := exec.Command(cmd.Cmd, cmd.Args...)
	result, err := c.CombinedOutput()
	if err != nil {
		return []byte(err.Error())
	}
	return result
}
