package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tupyy/osdviewer/internal/service"
	"github.com/tupyy/osdviewer/internal/tui"
)

var (
	CommitID string
)

func main() {
	token := os.Getenv("OCM_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "ocm token not found. Please provide token as env var OCM_TOKEN.\n")
		os.Exit(1)
	}

	fmt.Printf("Build from commit %q\n", CommitID)

	fm := service.NewDefaultFleetManagerCache(service.NewFleetManager(token))
	app := tview.NewApplication()
	tui := tui.New(app, fm)

	// ESC exits
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			tui.Stop()
			app.Stop()
		}
		tui.HandleEventKey(event)
		return event
	})

	app.SetRoot(tui.Layout(), true)
	tui.Start()
	app.Run()
}
