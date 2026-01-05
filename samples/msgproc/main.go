package main

import (
	"fmt"

	"github.com/mkch/gg"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

func main() {
	win, _ := window.New(&window.Spec{
		Text:  "Hello, Go!",
		Style: win32.WS_OVERLAPPEDWINDOW,
		X:     metrics.Px(win32.CW_USEDEFAULT),
		Width: metrics.Dip(500), Height: metrics.Dip(300),
		OnDestroy: func() { app.Quit(0) },
	})

	label := gg.Must(static.New(win.HWND(), &static.Spec{
		Style: win32.WS_VISIBLE | static.SS_CENTER | static.SS_CENTERIMAGE | static.SS_SUNKEN | static.SS_ENDELLIPSIS,
		X:     metrics.Dip(20), Y: metrics.Dip(20),
		Width: metrics.Dip(400), Height: metrics.Dip(30),
		Text: "Label",
	}))

	button.New(win.HWND(), &button.Spec{
		Text:  "Hello",
		Style: win32.WS_VISIBLE,
		X:     metrics.Dip(200), Y: metrics.Dip(120),
		Width: metrics.Dip(100), Height: metrics.Dip(60),
		OnClick: func() {
			win32util.MessageBox(win.HWND(),
				"Hello GUI!", "Button clicked",
				win32.MB_ICONINFORMATION)
		},
	})

	win.Show(win32.SW_SHOW)

	app.SetMessageDispatcher(func(msg *win32.MSG, prevProc func(msg *win32.MSG) win32.LRESULT) win32.LRESULT {
		switch msg.Message {
		case win32.WM_MOUSEMOVE:
			label.SetText(fmt.Sprintf("Mouse moved: HWND=0x%0x, X=%v, Y=%v", msg.Hwnd, win32.GET_X_LPARAM(msg.LParam), win32.GET_Y_LPARAM(msg.LParam)))
		}
		return prevProc(msg)
	})

	app.SetMessageDispatcher(func(msg *win32.MSG, prevProc func(msg *win32.MSG) win32.LRESULT) (r win32.LRESULT) {
		r = prevProc(msg)
		switch msg.Message {
		case win32.WM_MOUSEMOVE:
			fmt.Println("Mouse moved2:", msg.Hwnd, win32.GET_X_LPARAM(msg.LParam), win32.GET_Y_LPARAM(msg.LParam))
		}
		return
	})

	app.Run()
}
