package handler

import (
	"encoding/json"

	"github.com/riverchu/pkg/log"
	"github.com/riverchu/worker/base"
	"github.com/riverchu/worker/biz/service/command"
)

func WSHandle(msg []byte) []byte {
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

	return HandleMeta(meta)
}

func HandleMeta(meta *base.Meta) []byte {
	switch meta.Step {
	case base.StepCommand:
		result, err := command.Parse(meta.Detail)
		if err != nil {
			return []byte(err.Error())
		}
		return command.Exec(result)
	case base.StepScan:
		return nil
	default:
		return nil
	}
}
