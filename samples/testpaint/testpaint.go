package main

import (
	"strconv"

	"github.com/mkch/gg"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/paint/font"
	"github.com/mkch/gw/paint/pen"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

func main() {

	bkWin := gg.Must(window.New(&window.Spec{
		Text:    "Full screen",
		Style:   win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
		X:       metrics.Px(win32.CW_USEDEFAULT),
		Y:       metrics.Px(win32.INT(win32.SW_SHOWMAXIMIZED)),
		Width:   metrics.Px(win32.CW_USEDEFAULT),
		OnClose: func() { app.Quit(0) },
	}))

	linePen := gg.Must(pen.NewCosmetic(win32.PS_SOLID, win32.RGB(255, 0, 0)))
	defer linePen.Release()
	dpi := gg.Must(bkWin.DPI())
	lsf := font.SysDefault().LOGFONTW()
	lsf.Height = lsf.Height * 2 / 3
	textFont := gg.Must(font.New(font.NewLogFont(lsf, font.SysDefault().DPI()), dpi))
	defer textFont.Release()

	gridDpi := dpi
	ctxMenu := menu.New(true)
	var defDpiMenuItem *menu.Item
	var curDpiMenuItem *menu.Item
	defDpiMenuItem = gg.Must(ctxMenu.InsertItem(-1, &menu.ItemSpec{
		Title:   "Default DPI",
		Checked: false,
		OnClick: func() {
			gridDpi = win32.USER_DEFAULT_SCREEN_DPI
			bkWin.InvalidateRect(nil, true)
			defDpiMenuItem.SetChecked(true)
			curDpiMenuItem.SetChecked(false)
		},
	}))
	curDpiMenuItem = gg.Must(ctxMenu.InsertItem(-1, &menu.ItemSpec{
		Title:   "Current DPI",
		Checked: true,
		OnClick: func() {
			gridDpi = dpi
			bkWin.InvalidateRect(nil, true)
			defDpiMenuItem.SetChecked(false)
			curDpiMenuItem.SetChecked(true)
		},
	}))

	bkWin.OnRButtonDown = func(opt window.MouseClickOpt, x, y int) {
		bkWin.TrackPopupMenu(ctxMenu, nil)
	}

	const gridSize = win32.INT(50)

	bkWin.AddMsgListener(win32.WM_DPICHANGED, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
		gg.MustOK(textFont.ChangeDPI(gg.Must(bkWin.DPI())))
		bkWin.InvalidateRect(nil, true)
	})

	bkWin.SetPaintCallback(func(dc *paint.PaintDC, prev func(*paint.PaintDC)) {
		rcClient := gg.Must(bkWin.GetClientRect())
		rcClient.Right = metrics.DPIConv(rcClient.Right, dpi, gridDpi)
		rcClient.Bottom = metrics.DPIConv(rcClient.Bottom, dpi, gridDpi)

		defer gg.Must(paint.SelectObject(dc.HDC(), linePen.HPEN())).Restore()
		defer gg.Must(paint.SelectObject(dc.HDC(), textFont.HFONT())).Restore()

		var charBuf []win32.WCHAR
		for x := win32.INT(rcClient.Left) + gridSize; x <= win32.INT(rcClient.Right); x += gridSize {
			drawX := metrics.DPIConv(x, gridDpi, dpi)
			gg.MustOK(win32.MoveToEx(dc.HDC(), drawX, 0, nil))
			gg.MustOK(win32.LineTo(dc.HDC(), drawX, win32.INT(metrics.DPIConv(rcClient.Bottom, gridDpi, dpi))))
			win32util.CString(strconv.Itoa(int(x)), &charBuf)
			rect := win32.RECT{}
			gg.Must(win32.DrawTextExW(dc.HDC(), &charBuf[0], -1, &rect, win32.DT_CALCRECT, nil))
			width := rect.Width()
			rect.Left = win32.LONG(drawX) - width/2
			rect.Right = rect.Left + width
			gg.Must(win32.DrawTextExW(dc.HDC(), &charBuf[0], -1, &rect, win32.DT_CENTER, nil))
		}
		for y := win32.INT(rcClient.Top) + gridSize; y <= win32.INT(rcClient.Bottom); y += gridSize {
			drawY := metrics.DPIConv(y, gridDpi, dpi)
			gg.MustOK(win32.MoveToEx(dc.HDC(), 0, drawY, nil))
			gg.MustOK(win32.LineTo(dc.HDC(), win32.INT(metrics.DPIConv(rcClient.Right, gridDpi, dpi)), drawY))
			win32util.CString(strconv.Itoa(int(y)), &charBuf)
			rect := win32.RECT{}
			gg.Must(win32.DrawTextExW(dc.HDC(), &charBuf[0], -1, &rect, win32.DT_CALCRECT, nil))
			height := rect.Height()
			rect.Top = win32.LONG(drawY) - height/2
			rect.Bottom = rect.Top + height
			gg.Must(win32.DrawTextExW(dc.HDC(), &charBuf[0], -1, &rect, win32.DT_CENTER, nil))
		}
	})

	gg.Must(win32util.MessageBox(bkWin.HWND(), "Use context menu to change display", "Hint", win32.MB_ICONINFORMATION))

	textFontForWin1 := textFont.Clone()
	defer textFontForWin1.Release()

	win1 := gg.Must(window.New(&window.Spec{
		WndParent: bkWin.HWND(),
		Text:      "500 X 500",
		Style:     win32.WS_POPUP | win32.WS_CAPTION | win32.WS_VISIBLE,
		X:         metrics.Px(win32.CW_USEDEFAULT),
		Y:         metrics.Px(win32.INT(win32.SW_SHOWNORMAL)),
		Width:     metrics.Dip(500),
		Height:    metrics.Dip(500),
	}))
	win1.OnLButtonDown = func(opt window.MouseClickOpt, x, y int) {
		gg.Must(win32.SendMessageW(win1.HWND(), win32.WM_NCLBUTTONDOWN, win32.HTCAPTION, 0))
	}

	var win1CharBuf []win32.WCHAR
	win32util.CString("字体测试 ABCDEFG 50 100 150", &win1CharBuf)
	win1.SetPaintCallback(func(dc *paint.PaintDC, prev func(*paint.PaintDC)) {
		defer gg.Must(paint.SelectObject(dc.HDC(), textFontForWin1.HFONT())).Restore()
		rect := gg.Must(win1.GetClientRect())
		gg.Must(win32.DrawTextExW(dc.HDC(), &win1CharBuf[0], -1, rect, win32.DT_CENTER|win32.DT_SINGLELINE|win32.DT_VCENTER, nil))
	})

	win1.AddMsgListener(win32.WM_DPICHANGED, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
		gg.MustOK(textFontForWin1.ChangeDPI(gg.Must(win1.DPI())))
		win1.InvalidateRect(nil, true)
	})

	app.Run()
}
