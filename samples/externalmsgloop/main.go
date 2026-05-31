package main

import (
	"os"
	"runtime"
	"time"

	"github.com/mkch/gw/app"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -ico main.ico -manifest manifest.xml
//go:generate rsrc -arch 386 -ico main.ico -manifest manifest.xml

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	const TickerDuration = time.Millisecond * 100
	ticker := time.NewTicker(TickerDuration)

	app := app.NewBase()

	mainWindow := window.New(&window.Spec{
		Text:      "Hello, World!",
		Style:     win32.WS_OVERLAPPEDWINDOW,
		X:         metrics.Px(win32.CW_USEDEFAULT),
		Width:     metrics.Dip(500),
		Height:    metrics.Dip(300),
		OnDestroy: func() { app.Quit(0) },
	})

	timeStatic := static.New(&static.Spec{
		Parent: mainWindow,
		Text:   "Time",
		Style:  win32.WS_VISIBLE | static.SS_CENTER | static.SS_CENTERIMAGE,
		X:      metrics.Dip(200), Y: metrics.Dip(30),
		Width: metrics.Dip(100), Height: metrics.Dip(60),
	})

	mainWindow.Show(win32.SW_SHOW)

	go func() {
		for t := range ticker.C {
			str := t.Local().Format("15:04:05")
			// Run SetText() in UI goroutine.
			app.Post(func() {
				if !timeStatic.Valid() {
					return
				}
				timeStatic.SetText(str)
			})
		}
	}()

	var msg win32.MSG
	for {
		r := win32.GetMessageW(&msg, 0, 0, 0)
		if r == -1 {
			panic(r)
		}
		if r == 0 {
			break
		}
		if msg.Hwnd == 0 {
			continue // Messages not associated with a window cannot be dispatched
		}

		win32.TranslateMessage(&msg)
		win32.DispatchMessageW(&msg)
	}
	ticker.Stop()
	app.Destroy()
	os.Exit(int(msg.WParam))
}
