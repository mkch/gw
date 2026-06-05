package main

import (
	"os"

	"github.com/mkch/gg"
	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/dialog/fontdlg"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/layout"
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
	logFont  *font.LogFont
	textFont *font.Font
	textBuf  []win32.WCHAR
	dpi      win32.UINT
}

func (w *MainWindow) TextFont() *font.LogFont {
	return w.logFont
}

func (w *MainWindow) SetTextFont(f *fontdlg.FontChosen) {
	w.logFont = f.Font
	w.textFont.Release()
	w.textFont = gg.Must(font.New(w.logFont, w.dpi))
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
	gg.Must(win32.SetTextColor(dc.HDC(), win32.COLORREF(0x0000FF00))) // Green
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
		X:         gw.CW_USEDEFAULT,
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
				LogFont: win.TextFont(),
				HookProc: func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, cf *win32.CHOOSEFONTW, def fontdlg.DefaultHookProc) win32.UINT_PTR {
					if message == win32.WM_INITDIALOG {
						for child := range gw.Children(hwnd) {
							id := chkerr.Must(win32.GetWindowLongPtrW(child, win32.GWLP_ID))
							switch id {
							case fontdlg.ID_GLYPH_STATIC, fontdlg.ID_GLYPH_COMBO, fontdlg.ID_SIZE_STATIC, fontdlg.ID_SIZE_COMBO,
								fontdlg.ID_CHARSET_STATIC, fontdlg.ID_CHARSET_COMBO, fontdlg.ID_MEMO_STATIC,
								fontdlg.ID_LINK_MORE_FONTS:
								win32.ShowWindow(child, win32.SW_HIDE)
							}
						}

						center := &layout.Center{
							Item: &layout.Column{
								CrossAxisAlign: layout.AlignCenter,
								Items: []layout.Layout{
									&layout.Column{
										MainAxisSize: layout.AxisSizeMin,
										Items: []layout.Layout{
											&layout.Window{Hwnd: chkerr.Must(win32.GetDlgItem(hwnd, fontdlg.ID_FONT_STATIC))},
											&layout.Window{Hwnd: chkerr.Must(win32.GetDlgItem(hwnd, fontdlg.ID_FONT_COMBO))},
										},
									},
									&layout.Window{Hwnd: chkerr.Must(win32.GetDlgItem(hwnd, fontdlg.ID_SAMPLE_GROUPBOX))},
									&layout.Window{Hwnd: chkerr.Must(win32.GetDlgItem(hwnd, fontdlg.ID_SAMPLE_STATIC))},
									&layout.Row{
										MainAxisSize: layout.AxisSizeMin,
										Items: []layout.Layout{
											&layout.Window{Hwnd: chkerr.Must(win32.GetDlgItem(hwnd, fontdlg.ID_OK))},
											&layout.Window{Hwnd: chkerr.Must(win32.GetDlgItem(hwnd, fontdlg.ID_CANCEL))},
										},
									},
								},
							},
						}

						client := chkerr.Must(layout.ClientSize(hwnd))
						chkerr.MustOK(layout.Perform(center, &layout.Size{Width: client.Width, Height: client.Height}))
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
