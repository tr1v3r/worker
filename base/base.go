package base

import (
	"encoding/json"

	wm "github.com/riverchu/workmanager"
)

// Meta ...
type Meta struct {
	Step   wm.WorkStep     `json:"step"`
	Detail json.RawMessage `json:"detail"`
}

const (
	// StepCommand ...
	StepCommand wm.WorkStep = "command"
	// StepScan ...
	StepScan wm.WorkStep = "scan"
)

// CommandMeta ...
type CommandMeta struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

// ScanMeta ...
type ScanMeta struct {
	Type      string `json:"type"`
	TaskToken string `json:"task_token"`
}
