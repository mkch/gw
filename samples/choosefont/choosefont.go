package main

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/dialog"
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

func main() {
	win := gg.Must(window.New(&window.Spec{
		Text:      "Test font",
		Style:     win32.WS_OVERLAPPEDWINDOW,
		X:         metrics.Px(win32.CW_USEDEFAULT),
		Width:     metrics.Dip(500),
		Height:    metrics.Dip(300),
		OnDestroy: func() { app.Quit(0) },
	}))

	dpi := gg.Must(win.DPI())
	lf := font.SysDefault()
	textFont := gg.Must(font.New(lf, dpi))
	var textColor win32.COLORREF

	win.AddMsgListener(win32.WM_SIZE, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
		win.InvalidateRect(nil, true)
	})

	win.AddMsgListener(win32.WM_DPICHANGED, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
		dpi = gg.Must(win.DPI())
		gg.MustOK(textFont.ChangeDPI(dpi))
		win.InvalidateRect(nil, true)
	})

	updateFont := func() {
		textFont.Release()
		textFont = gg.Must(font.New(lf, dpi))
	}

	fontMenu := menu.New(false)
	fontMenu.InsertItem(-1, &menu.ItemSpec{
		Title: "Choose &font",
		OnClick: func() {
			r, err := dialog.ChooseFont(&dialog.ChooseFontSpec{
				Owner:   win.HWND(),
				Flags:   win32.CF_EFFECTS,
				Color:   &textColor,
				LogFont: lf,
				OnApply: func(curFont *dialog.FontChosen) {
					lf = curFont.Font
					textColor = curFont.Color
					updateFont()
					win.InvalidateRect(nil, true)
				},
			})
			if err != nil {
				panic(err)
			}
			if r != nil {
				lf = r.Font
				textColor = r.Color
				updateFont()
				win.InvalidateRect(nil, true)
			}
		},
	})
	mainMenu := menu.New(false)
	mainMenu.InsertItem(-1, &menu.ItemSpec{
		Title:   "&Font",
		Submenu: fontMenu,
	})

	win.SetMenu(mainMenu)

	const text = "微软中文软件 Test font"
	var textBuf []win32.WCHAR
	win32util.CString(text, &textBuf)
	win.AddPaintCallback(func(paintData *paint.PaintData, prev func(*paint.PaintData)) {
		defer gg.Must(paint.SelectObject(paintData.DC, textFont.HFONT())).Restore()
		gg.Must(win32.SetTextColor(paintData.DC, textColor))
		rcClient, _ := win.GetClientRect()
		win32.DrawTextExW(paintData.DC, &textBuf[0], -1, rcClient, win32.DT_CENTER|win32.DT_VCENTER|win32.DT_SINGLELINE, nil)
	})

	win.Show(win32.SW_SHOW)
	app.Run()
	textFont.Release()
}
