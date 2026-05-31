package gw_test

import (
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

// MainWindow is a Window that quits the message loop when destroyed.
type MainWindow struct {
	window.Window
}

func (w *MainWindow) OnDestroy() {
	// Call the default implementation.
	w.Window.OnDestroy()
	// Quit the message loop when this window is destroyed.
	w.App().Quit(1)
}

// NewMainWindow creates and initializes a MainWindow instance.
func NewMainWindow() *MainWindow {
	w := &MainWindow{window.Window{Spec: &window.Spec{
		X:      metrics.Dip(100),
		Y:      metrics.Dip(100),
		Width:  metrics.Dip(500),
		Height: metrics.Dip(300),
	}}}
	gw.Init(w)
	return w
}

func Example_wrapper() {
	gw.Run(func(app *app.App) {
		mainWindow := NewMainWindow()
		mainWindow.SetText("Main Window")
		mainWindow.Show(win32.SW_SHOWNORMAL)
	}, nil)
}
