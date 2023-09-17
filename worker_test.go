package worker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/tr1v3r/pkg/log"
	"github.com/tr1v3r/worker"
	wm "github.com/tr1v3r/workmanager"
)

func TestPingerOutputer(t *testing.T) {
	var (
		check  wm.WorkStep = "check"
		output wm.WorkStep = "output"

		checker wm.WorkerName = "checker"
		printer wm.WorkerName = "printer"
	)

	var (
		checkStepRunner = func(ctx context.Context, work wm.Work, target wm.WorkTarget, nexts ...func(wm.WorkTarget)) {
			if err := ctx.Err(); err != nil {
				return
			}

			results, err := work(target, map[wm.WorkerName]wm.WorkerConfig{checker: new(wm.DummyConfig)})
			if err != nil {
				log.Error("work fail: %s", err)
				return
			}

			for _, result := range results {
				for _, next := range nexts {
					next(result)
				}
			}
		}
		outputStepRunner = func(ctx context.Context, work wm.Work, target wm.WorkTarget, _ ...func(wm.WorkTarget)) {
			if err := ctx.Err(); err != nil {
				return
			}

			_, err := work(target, map[wm.WorkerName]wm.WorkerConfig{printer: new(wm.DummyConfig)})
			if err != nil {
				log.Error("work fail: %s", err)
				return
			}
		}
	)

	mgr := wm.NewWorkerManager(context.Background())

	mgr.RegisterWorker(checker, worker.PingerBuilder)
	mgr.RegisterWorker(printer, worker.PrinterBuilder)

	mgr.RegisterStep(check, checkStepRunner, output)
	mgr.RegisterStep(output, outputStepRunner)

	mgr.Serve(check, output)

	target := new(worker.HostTarget)
	target.HostName = "qq.com"

	if err := mgr.Recv(check, target); err != nil {
		fmt.Printf("send target fail: %s", err)
		return
	}

	time.Sleep(10 * time.Second)
	log.Flush()
}
