package gw_test

import (
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint/brush"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

// MainWindow is the main application window.
type MainWindow struct {
	window.Window
	bkColor win32.COLORREF
	bkBrush *brush.Brush
}

func (w *MainWindow) OnInit() (err error) {
	err = w.Window.OnInit()
	if err != nil {
		return err
	}
	// Create a solid brush with the specified background color.
	w.bkBrush, err = brush.New(&win32.LOGBRUSH{
		Style: win32.BS_SOLID,
		Color: w.bkColor,
	})
	return
}

func (w *MainWindow) OnPaint(evt *events.PaintEvent) {
	dc, _ := evt.Begin()
	// Draw the background using the background brush.
	win32.FillRect(dc.HDC(), dc.Rect(), w.bkBrush.HBRUSH())
}

func (w *MainWindow) OnDestroy() {
	// Call the default implementation.
	w.Window.OnDestroy()
	// Release the background brush.
	w.bkBrush.Release()
	// Quit the message loop when this window is destroyed.
	w.App().Quit(1)
}

// NewMainWindow creates and initializes a MainWindow instance.
func NewMainWindow(title string, bkColor win32.COLORREF) (*MainWindow, error) {
	return gw.Init(&MainWindow{
		Window: window.Window{
			Spec: &window.Spec{
				Style:  win32.WS_OVERLAPPEDWINDOW,
				Text:   title,
				X:      gw.CW_USEDEFAULT,
				Width:  metrics.Dip(500),
				Height: metrics.Dip(300),
			},
		},
		bkColor: bkColor,
	})
}

func Example_wrapper() {
	gw.Run(func(app *app.App) {
		mainWindow, _ := NewMainWindow("Main Window", win32.RGB(128, 128, 0))
		mainWindow.Show(win32.SW_SHOWNORMAL)
	}, nil)
}
