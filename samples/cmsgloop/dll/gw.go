package main

// #include <windows.h>
import "C"

import (
	"time"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw/app/gwapp"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -ico main.ico -manifest manifest.xml
//go:generate rsrc -arch 386 -ico main.ico -manifest manifest.xml

var app *gwapp.BareApp
var ticker *time.Ticker

//export Show
func Show(parent C.HWND) {
	const TickerDuration = time.Millisecond * 100
	ticker = time.NewTicker(TickerDuration)

	app = gwapp.NewBare(nil)

	mainWindow := chkerr.Must(window.New(&window.Spec{
		WndParent: win32.HWND(unsafe.Pointer(parent)),
		Text:      "GW window",
		Style:     win32.WS_OVERLAPPEDWINDOW,
		X:         metrics.Px(win32.CW_USEDEFAULT),
		Width:     metrics.Dip(500),
		Height:    metrics.Dip(300),
	}))

	fileMenu := menu.New(true)
	fileMenu.InsertItem(-1, &menu.ItemSpec{
		Title:    "&Close",
		AccelKey: menu.AccelKey{Mod: menu.ModCtrl, VKeyCode: 'Q'},
		OnClick:  func() { mainWindow.Destroy() },
	})

	mainMenu := menu.New(false)
	mainMenu.InsertItem(-1, &menu.ItemSpec{
		Title:   "&File",
		Submenu: fileMenu,
	})

	mainWindow.SetMenu(mainMenu)

	timeStatic := gg.Must(static.New(mainWindow.HWND(), &static.Spec{
		Text:  "Time",
		Style: win32.WS_VISIBLE | static.SS_CENTER | static.SS_CENTERIMAGE,
		X:     metrics.Dip(200), Y: metrics.Dip(30),
		Width: metrics.Dip(100), Height: metrics.Dip(60),
	}))

	mainWindow.Show(win32.SW_SHOW)

	go func() {
		for t := range ticker.C {
			str := t.Local().Format("15:04:05")
			// Run SetText() in UI goroutine.
			app.Post(func() { timeStatic.SetText(str) })
		}
	}()
}

//export Cleanup
func Cleanup() {
	ticker.Stop()
}

func main() {}
