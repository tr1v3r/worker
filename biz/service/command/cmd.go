package command

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/riverchu/worker/base"
)

// Parse parse data to command meta info struct
func Parse(data []byte) (*base.CommandMeta, error) {
	var cmd base.CommandMeta
	if err := json.Unmarshal(data, &cmd); err != nil {
		return nil, fmt.Errorf("unmarshal detail fail: %w", err)
	}
	return &cmd, nil
}

// Exec execute command
func Exec(cmd *base.CommandMeta) []byte {
	c := exec.Command(cmd.Cmd, cmd.Args...)
	result, err := c.CombinedOutput()
	if err != nil {
		return []byte(err.Error())
	}
	return result
}
