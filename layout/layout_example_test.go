package layout_test

import (
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/layout"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

func Example_layout() {
	gw.Run(func(app *app.App) {
		// Create the window and the button.
		win, _ := window.New(&window.Spec{
			Text:      "Resize to see layout in action",
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(400),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		})

		btn, _ := button.New(&button.Spec{
			Parent: win,
			Text:   "Always center in window.",
			Style:  win32.WS_VISIBLE,
			Width:  metrics.Dip(200),
			Height: metrics.Dip(50),
		})

		// Create a Center layout with the button as the item.
		center := &layout.Center{
			Child: &layout.Intrinsic{Hwnd: btn.HWND()},
		}
		// Build the layout element tree.
		tree, _ := layout.Build(center)
		// Perform the layout for the first time.
		layout.PerformWindow(tree, win.HWND())
		// Perform the layout whenever the window is resized.
		win.SetOnSizeListener(func(event events.SizeEvent) {
			size, _ := layout.EventSize(win.HWND(), event)
			layout.Perform(tree, size)
		})

		win.Show(win32.SW_SHOW)

	}, nil)
}
