package dir

import (
	"context"

	wm "github.com/riverchu/workmanager"
)

const DirSearchWorker wm.WorkerName = "dirsearch"

// DirSearchWorkerBuilder builder
func DirSearchWorkerBuilder(ctx context.Context, args map[string]interface{}) wm.Worker {
	return nil
}
