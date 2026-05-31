package main

import (
	"os"
	"time"

	"github.com/mkch/gg"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/paint/brush"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

var ticker *time.Ticker

type MyWindow struct {
	window.Window
	bkBrush *brush.Brush
}

func (w *MyWindow) OnInit() {
	w.Window.OnInit()
	w.bkBrush = gg.Must(brush.New(&win32.LOGBRUSH{
		Style: win32.BS_HATCHED,
		Color: win32.RGB(255, 0, 0),
		Hatch: win32.HS_DIAGCROSS,
	}))
}

func (w *MyWindow) OnDestroy() {
	w.bkBrush.Release()
	w.Window.OnDestroy()
}

func (w *MyWindow) OnPaint() {
	dc := gg.Must(paint.NewPaintDC(w.HWND()))
	defer func() { gg.MustOK(dc.EndPaint()) }()
	win32.FillRect(dc.HDC(), dc.Rect(), w.bkBrush.HBRUSH())
}

func NewMyWindow(spec *window.Spec) (win *MyWindow) {
	win = &MyWindow{Window: window.Window{Spec: spec}}
	gw.Init(win)
	return
}

func main() {
	os.Exit(gw.Run(ui, Cleanup))
}

func ui(app *app.App) {

	win := NewMyWindow(&window.Spec{
		Text:  "Static demo",
		Style: win32.WS_OVERLAPPEDWINDOW,
		X:     metrics.Px(win32.CW_USEDEFAULT),
		Width: metrics.Dip(500), Height: metrics.Dip(300),
		OnDestroy: func() {
			app.Quit(0)
		},
	})

	ctrl := static.New(&static.Spec{
		Parent: win,
		Text:   "Hello, World!",
		Style:  win32.WS_VISIBLE | static.SS_CENTER | static.SS_CENTERIMAGE,
		X:      metrics.Dip(200), Y: metrics.Dip(30),
		Width: metrics.Dip(100), Height: metrics.Dip(60),
		//BackgroundColor: &color,
	})

	ticker = time.NewTicker(time.Millisecond * 500)
	go func() {
		for {
			str := (<-ticker.C).Local().Format("15:04:05")
			// Run SetText() in UI goroutine.
			app.Post(func() {
				if !ctrl.Valid() {
					return
				}
				ctrl.SetText(str)
			})
		}
	}()

	win.Show(win32.SW_SHOW)
}

func Cleanup(app *app.App) {
	ticker.Stop()
}
