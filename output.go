// Package worker ...
package worker

import (
	"context"
	"fmt"

	wm "github.com/tr1v3r/workmanager"
)

var (
	_ wm.Worker = new(Printer)
)

func PrinterBuilder(ctx context.Context, args map[string]any) wm.Worker { return new(Printer) }

// Printer print worker
type Printer struct{}

func (p *Printer) WithContext(context.Context) wm.Worker { return p }
func (p *Printer) Work(targets ...wm.WorkTarget) ([]wm.WorkTarget, error) {
	for _, target := range targets {
		fmt.Printf("got result:\n%s\n", target)
	}
	return nil, nil
}
