package main

import (
	"os"

	"github.com/mkch/gg"
	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/dialog/fontdlg"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/paint/font"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

type MainWindow struct {
	window.Window
	logFont   *font.LogFont
	textFont  *font.Font
	textBuf   []win32.WCHAR
	dpi       win32.UINT
	textColor win32.COLORREF
}

func (w *MainWindow) TextFont() *font.LogFont {
	return w.logFont
}

func (w *MainWindow) TextColor() win32.COLORREF {
	return w.textColor
}

func (w *MainWindow) SetTextFont(f *fontdlg.FontChosen) {
	w.logFont = f.Font
	w.textFont.Release()
	w.textFont = gg.Must(font.New(w.logFont, w.dpi))
	w.textColor = f.Color
	w.InvalidateRect(nil, true)
}

func (w *MainWindow) OnInit() error {
	if err := w.Window.OnInit(); err != nil {
		return err
	}
	w.dpi = gg.Must(w.DPI())
	lf := font.SysDefault()
	w.textFont = gg.Must(font.New(lf, w.dpi))
	win32util.CString("微软中文软件 Test font", &w.textBuf)

	w.AddMsgListener(win32.WM_SIZE, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
		w.InvalidateRect(nil, true)
	})

	w.AddMsgListener(win32.WM_DPICHANGED, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
		w.dpi = gg.Must(w.DPI())
		gg.MustOK(w.textFont.ChangeDPI(w.dpi))
		w.InvalidateRect(nil, true)
	})
	return nil
}

func (w *MainWindow) OnDestroy() {
	w.textFont.Release()
	w.Window.OnDestroy()
}

func (w *MainWindow) OnPaint(evt *events.PaintEvent) {
	dc := gg.Must(evt.Begin())
	defer gg.Must(paint.SelectObject(dc.HDC(), w.textFont.HFONT())).Restore()
	gg.Must(win32.SetTextColor(dc.HDC(), w.textColor))
	rcClient, _ := w.GetClientRect()
	win32.DrawTextExW(dc.HDC(), &w.textBuf[0], -1, rcClient, win32.DT_CENTER|win32.DT_VCENTER|win32.DT_SINGLELINE, nil)
}

func NewMainWindow(spec *window.Spec) (*MainWindow, error) {
	return gw.Init(&MainWindow{Window: window.Window{Spec: spec}})
}

func ui(app *app.App) {

	win := chkerr.Must(NewMainWindow(&window.Spec{
		Text:      "Test font",
		Style:     win32.WS_OVERLAPPEDWINDOW,
		X:         metrics.Px(win32.CW_USEDEFAULT),
		Width:     metrics.Dip(500),
		Height:    metrics.Dip(300),
		OnDestroy: func() { app.Quit(0) },
	}))

	fontMenu := menu.New(false)
	fontMenu.InsertItem(-1, &menu.ItemSpec{
		Title: "Choose &font",
		OnClick: func() {
			r, err := fontdlg.ChooseFont(&fontdlg.Spec{
				Owner:   win.HWND(),
				Flags:   win32.CF_EFFECTS,
				Color:   new(win.TextColor()),
				LogFont: win.TextFont(),
				OnApply: func(curFont *fontdlg.FontChosen) {
					win.SetTextFont(curFont)
				},
				HookProc: func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, cf *win32.CHOOSEFONTW, def fontdlg.DefaultHookProc) win32.UINT_PTR {
					if message == win32.WM_INITDIALOG {
						c := chkerr.Must(win32.GetDlgItem(hwnd, fontdlg.ID_CHARSET_STATIC))
						var buf []win32.WCHAR
						win32util.CString("~字符集~:", &buf)
						chkerr.MustOK(win32.SetWindowTextW(c, &buf[0]))
					}
					return def(hwnd, message, wParam, lParam)
				},
			})
			if err != nil {
				panic(err)
			}
			if r != nil {
				win.SetTextFont(r)
			}
		},
	})
	mainMenu := menu.New(false)
	mainMenu.InsertItem(-1, &menu.ItemSpec{
		Title:   "&Font",
		Submenu: fontMenu,
	})

	win.SetMenu(mainMenu)

	win.Show(win32.SW_SHOW)
}

func main() {
	ret := gw.Run(ui, nil)
	os.Exit(ret)
}
