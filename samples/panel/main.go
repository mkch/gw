package main

import (
	"time"

	"github.com/mkch/gg"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/paint/brush"
	"github.com/mkch/gw/panel"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

func main() {
	win := gg.Must(window.New(&window.Spec{
		Text:  "Panel demo",
		Style: win32.WS_OVERLAPPEDWINDOW,
		X:     metrics.Px(win32.CW_USEDEFAULT),
		Width: metrics.Dip(500), Height: metrics.Dip(300),
		OnDestroy: func() { app.Quit(0) },
	}))
	winBkBrush := gg.Must(brush.New(&win32.LOGBRUSH{
		Style: win32.BS_HATCHED,
		Color: win32.RGB(255, 0, 0),
		Hatch: win32.HS_DIAGCROSS,
	}))
	defer winBkBrush.Release()
	win.AddPaintCallback(func(paintData *paint.PaintData, prev func(*paint.PaintData)) {
		rect := gg.Must(win.GetClientRect())
		win32.FillRect(paintData.DC, rect, winBkBrush.HBRUSH())
	})

	panelCtrl := gg.Must(panel.New(win.HWND(), &panel.Spec{
		X: metrics.Dip(10), Y: metrics.Dip(10),
		Width: metrics.Dip(120), Height: metrics.Dip(80),
	}))

	panelCtrl.SetBackgroundColor(win32.RGB(220, 220, 220))

	staticCtrl := gg.Must(static.New(panelCtrl.HWND(), &static.Spec{
		Text:  "Hello, World!",
		Style: win32.WS_VISIBLE | static.SS_CENTER | static.SS_CENTERIMAGE,
		X:     metrics.Dip(10), Y: metrics.Dip(10),
		Width: metrics.Dip(100), Height: metrics.Dip(60),
	}))
	staticCtrl.SetBackgroundColor(win32.RGB(200, 200, 255))

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	go func() {
		for {
			str := (<-ticker.C).Local().Format("15:04:05")
			// Run SetText() in UI goroutine.
			app.Post(func() { staticCtrl.SetText(str) })
		}
	}()

	win.Show(win32.SW_SHOW)
	app.Run()
}
