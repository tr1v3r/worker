package scan

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/riverchu/pkg/log"
	ws "github.com/riverchu/pkg/websocket"
	wm "github.com/riverchu/workmanager"

	"github.com/riverchu/worker/base"
	"github.com/riverchu/worker/biz/service/scan/scanner/dir"
	"github.com/riverchu/worker/config"
)

var (
	DirStep      wm.WorkStep = "dirsearch"
	TransferStep wm.WorkStep = "transfer"
)

// StepRunner func(work Work, workTarget WorkTarget, nexts ...func(WorkTarget))

func init() {
	wm.RegisterWorker(dir.DirSearchWorker, dir.DirSearchWorkerBuilder)

	wm.RegisterStep(DirStep, stepRunner, TransferStep)
	wm.RegisterStep(TransferStep, wm.TransferRunner(func(target wm.WorkTarget) {
		conn := config.GetWSConn(target.Token())
		if conn == nil {
			log.Error("cannot find ws conn(%s)", target.Token())
			return
		}

		data, _ := json.Marshal(target)
		err := ws.Write(conn, data)
		if err != nil {
			log.Error("write to websocket fail: %s", err)
		}
		return
	}))

	wm.Serve()
}

func stepRunner(work wm.Work, target wm.WorkTarget, nexts ...func(wm.WorkTarget)) {
	results, err := work(target, map[wm.WorkerName]wm.WorkerConfig{
		dir.DirSearchWorker: &wm.DummyConfig{},
	})
	if err != nil {
		log.Error("dir search fail: %s", err)
		return
	}

	for _, result := range results {
		for _, next := range nexts {
			next(result)
		}
	}
}

func Recv(step wm.WorkStep, s *base.ScanMeta) error {
	task := wm.NewTask(context.Background())
	task.(*wm.Task).TaskToken = s.TaskToken
	wm.AddTask(task)
	return wm.Recv(step, &base.ScanTarget{ScanMeta: *s})
}

// Parse ...
func Parse(data []byte) (*base.ScanMeta, error) {
	var m base.ScanMeta
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal detail fail: %w", err)
	}
	return &m, nil
}
