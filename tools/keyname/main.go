package main

import (
	"fmt"
	"os"

	"github.com/mkch/gg"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

func ui(app *app.App) {
	var label *static.Static
	win := gg.Must(window.New(&window.Spec{
		Text:      "Key Name",
		Style:     win32.WS_OVERLAPPEDWINDOW,
		X:         metrics.Px(win32.CW_USEDEFAULT),
		Width:     metrics.Dip(500),
		Height:    metrics.Dip(300),
		OnDestroy: func() { app.Quit(0) },
	}))
	win.SetWndProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevWndProc win32.WndProc) win32.LRESULT {
		switch message {
		case win32.WM_SIZE:
			dpi := gg.Must(win32.GetDpiForWindow(hwnd))
			winWidth := win32.INT(win32.GET_X_LPARAM(lParam))
			padding := metrics.Dip(5).Px(dpi)
			h := metrics.Dip(100).Px(dpi)
			win32.SetWindowPos(label.HWND(), 0, padding, padding, winWidth-padding*2, h, win32.SWP_NOZORDER)
		case win32.WM_KEYDOWN:
			var buf [32]win32.WCHAR
			_, err := win32.GetKeyNameTextW(win32.LONG(lParam), &buf[0], len(buf))
			if err != nil {
				label.SetText("Error: " + err.Error())
				break
			}
			keyName := win32util.GoString(&buf[0], win32util.CStrLen(&buf[0], len(buf))+1)
			label.SetText(fmt.Sprintf("Virtual Key: 0x%02X('%s')\n Key: %s\n lParam: %s", wParam, string(rune(wParam)), keyName, win32util.KeyMessageLParam(lParam)))
		}
		return prevWndProc(hwnd, message, wParam, lParam)
	})

	label = gg.Must(static.New(win.HWND(), &static.Spec{
		Style: win32.WS_VISIBLE | win32.WS_BORDER | static.SS_CENTER,
	}))

	win.Show(win32.SW_SHOW)
}

func main() {
	ret := gw.Run(ui, nil)
	os.Exit(ret)
}
