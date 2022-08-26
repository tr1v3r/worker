package dir

import (
	"context"
	"fmt"

	"github.com/riverchu/worker/base"
	wm "github.com/riverchu/workmanager"
)

const DirSearchWorker wm.WorkerName = "dirsearch"

// DirSearchWorkerBuilder builder
func DirSearchWorkerBuilder(ctx context.Context, args map[string]interface{}) wm.Worker {
	return &dirSearchWorker{
		finish: make(chan struct{}, 1),
	}
}

type dirSearchWorker struct {
	wm.DummyWorker

	finish chan struct{}
}

func (w dirSearchWorker) Work(target wm.WorkTarget) ([]wm.WorkTarget, error) {
	t, ok := target.(*base.ScanTarget)
	if !ok {
		return nil, fmt.Errorf("target is not ScanMeta")
	}

	return []wm.WorkTarget{&base.ScanTarget{ScanMeta: base.ScanMeta{
		Type:      "work_result",
		TaskToken: target.Token(),
		Target:    t.Target,
	}}}, nil
}

func (w dirSearchWorker) AfterWork() { w.finish <- struct{}{} }

func (w dirSearchWorker) Finished() <-chan struct{} { return w.finish }
