package main

import (
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -ico main.ico
//go:generate rsrc -arch 386 -ico main.ico

func main() {
	win, _ := window.New(&window.Spec{
		Text:  "Hello, Go!",
		Style: win32.WS_OVERLAPPEDWINDOW,
		X:     win32.CW_USEDEFAULT,
		Width: 500, Height: 300,
		OnClose: func() { app.Quit(0) },
	})
	button.New(win.HWND(), &button.Spec{
		Text:  "Hello",
		Style: win32.WS_VISIBLE,
		X:     200, Y: 120,
		Width: 100, Height: 60,
		OnClick: func() {
			win32util.MessageBox(win.HWND(),
				"Hello GUI!", "Button clicked",
				win32.MB_ICONINFORMATION)
		},
	})
	win.Show(win32.SW_SHOW)
	app.Run()
}
