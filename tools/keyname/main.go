package main

import (
	"fmt"
	"os"

	"github.com/mkch/gg"
	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

type MainWindow struct {
	window.Window
	label *static.Static
}

func (w *MainWindow) OnInit() (err error) {
	err = w.Window.OnInit()
	if err != nil {
		return err
	}
	w.label, err = static.New(&static.Spec{
		Parent: w,
		Style:  win32.WS_VISIBLE | win32.WS_BORDER | static.SS_CENTER,
	})
	return err
}

func (w *MainWindow) WndProc(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	switch message {
	case win32.WM_SIZE:
		dpi := gg.Must(win32.GetDpiForWindow(hwnd))
		winWidth := win32.INT(win32.GET_X_LPARAM(lParam))
		padding := win32.INT(metrics.Dip(5).Px(dpi))
		h := win32.INT(metrics.Dip(100).Px(dpi))
		win32.SetWindowPos(w.label.HWND(), 0, padding, padding, winWidth-padding*2, h, win32.SWP_NOZORDER)
	}
	return w.Window.WndProc(hwnd, message, wParam, lParam)
}

func (w *MainWindow) OnKeyDown(event *events.KeyEvent) {
	w.Window.OnKeyDown(event)

	var buf [32]win32.WCHAR
	_, err := win32.GetKeyNameTextW(win32.LONG(event.State), &buf[0], len(buf))
	if err != nil {
		w.label.SetText("Error: " + err.Error())
		return
	}
	keyName := win32util.GoString(&buf[0], win32util.CStrLen(&buf[0], len(buf))+1)
	w.label.SetText(fmt.Sprintf("Virtual Key: 0x%02X('%s')\n Key: %s\n lParam: %s", event.VKCode, string(rune(event.VKCode)), keyName, event.State))
}

func NewMainWindow(spec *window.Spec) (*MainWindow, error) {
	return gw.Init(&MainWindow{Window: window.Window{Spec: spec}})
}

func ui(app *app.App) {
	win := chkerr.Must(NewMainWindow(&window.Spec{
		Text:      "Key Name",
		Style:     win32.WS_OVERLAPPEDWINDOW,
		X:         gw.CW_USEDEFAULT,
		Width:     metrics.Dip(500),
		Height:    metrics.Dip(300),
		OnDestroy: func() { app.Quit(0) },
	}))

	win.Show(win32.SW_SHOW)
}

func main() {
	ret := gw.Run(ui, nil)
	os.Exit(ret)
}
