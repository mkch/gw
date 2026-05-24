package main

import (
	"fmt"

	"github.com/mkch/gg"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/panel"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -ico main.ico -manifest manifest.xml
//go:generate rsrc -arch 386 -ico main.ico -manifest manifest.xml

func main() {
	win := gg.Must(window.New(&window.Spec{
		Text:  "Default Class Window",
		Style: win32.WS_OVERLAPPEDWINDOW,
		X:     metrics.Px(win32.CW_USEDEFAULT),
		Width: metrics.Dip(500), Height: metrics.Dip(300),
		OnDestroy: func() { app.Quit(0) },
	}))

	panel1 := gg.Must(panel.New(win.HWND(), &panel.Spec{
		X: metrics.Dip(20), Y: metrics.Dip(20),
		Width: metrics.Dip(200), Height: metrics.Dip(100),
	}))

	panel1.SetBackgroundColor(win32.COLORREF(0x00FF0000))

	win.Show(win32.SW_SHOW)

	const windowClassName = "my window class name"
	const panelClassName = "my panel class name"

	win2 := gg.Must(window.New(&window.Spec{
		ClassName: windowClassName,
		Text:      fmt.Sprintf("Window Class Name: %q", windowClassName),
		Style:     win32.WS_OVERLAPPEDWINDOW,
		X:         metrics.Px(win32.CW_USEDEFAULT),
		Width:     metrics.Dip(400), Height: metrics.Dip(200),
		OnDestroy: func() { app.Quit(0) },
	}))

	var nameBuffer [256]win32.WCHAR
	n := gg.Must(win32.GetClassNameW(win2.HWND(), &nameBuffer[0], len(nameBuffer)))
	if win32util.GoString(&nameBuffer[0], n+1) != windowClassName {
		panic("unexpected class name")
	}

	panel2 := gg.Must(panel.New(win2.HWND(), &panel.Spec{
		ClassName: panelClassName,
		X:         metrics.Dip(20), Y: metrics.Dip(20),
		Width: metrics.Dip(200), Height: metrics.Dip(100),
	}))

	panel2.SetBackgroundColor(win32.COLORREF(0x0000FF00))

	n = gg.Must(win32.GetClassNameW(panel2.HWND(), &nameBuffer[0], len(nameBuffer)))
	if win32util.GoString(&nameBuffer[0], n+1) != panelClassName {
		panic("unexpected class name")
	}

	win2.Show(win32.SW_SHOW)

	app.Run()
}
