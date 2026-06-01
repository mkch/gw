package main

import (
	"os"
	"strconv"

	"github.com/mkch/gg"
	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/events"
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

const gridSize = win32.INT(50)

type MainWindow struct {
	window.Window
	linePen  *pen.Pen
	textFont *font.Font
	dpi      win32.UINT
	gridDpi  win32.UINT
}

func (w *MainWindow) OnInit() {
	w.Window.OnInit()
	w.linePen = gg.Must(pen.NewCosmetic(win32.PS_SOLID, win32.RGB(255, 0, 0)))

	w.dpi = gg.Must(w.DPI())
	w.gridDpi = w.dpi
	lsf := font.SysDefault().LOGFONTW()
	lsf.Height = lsf.Height * 2 / 3
	w.textFont = gg.Must(font.New(font.NewLogFont(lsf, font.SysDefault().DPI()), w.dpi))
}

func (w *MainWindow) OnDestroy() {
	w.linePen.Release()
	w.textFont.Release()
	w.Window.OnDestroy()
}

func (w *MainWindow) SetGridDpi(dpi win32.UINT) {
	w.gridDpi = dpi
	w.InvalidateRect(nil, true)
}

func (w *MainWindow) OnPaint() {
	dc := chkerr.Must(paint.NewPaintDC(w.HWND()))
	defer func() { chkerr.MustOK(dc.EndPaint()) }()

	rcClient := gg.Must(w.GetClientRect())
	rcClient.Right = metrics.DPIConv(rcClient.Right, w.dpi, w.gridDpi)
	rcClient.Bottom = metrics.DPIConv(rcClient.Bottom, w.dpi, w.gridDpi)

	defer gg.Must(paint.SelectObject(dc.HDC(), w.linePen.HPEN())).Restore()
	defer gg.Must(paint.SelectObject(dc.HDC(), w.textFont.HFONT())).Restore()

	var charBuf []win32.WCHAR
	for x := win32.INT(rcClient.Left) + gridSize; x <= win32.INT(rcClient.Right); x += gridSize {
		drawX := metrics.DPIConv(x, w.gridDpi, w.dpi)
		gg.MustOK(win32.MoveToEx(dc.HDC(), drawX, 0, nil))
		gg.MustOK(win32.LineTo(dc.HDC(), drawX, win32.INT(metrics.DPIConv(rcClient.Bottom, w.gridDpi, w.dpi))))
		win32util.CString(strconv.Itoa(int(x)), &charBuf)
		rect := win32.RECT{}
		gg.Must(win32.DrawTextExW(dc.HDC(), &charBuf[0], -1, &rect, win32.DT_CALCRECT, nil))
		width := rect.Width()
		rect.Left = win32.LONG(drawX) - width/2
		rect.Right = rect.Left + width
		gg.Must(win32.DrawTextExW(dc.HDC(), &charBuf[0], -1, &rect, win32.DT_CENTER, nil))
	}
	for y := win32.INT(rcClient.Top) + gridSize; y <= win32.INT(rcClient.Bottom); y += gridSize {
		drawY := metrics.DPIConv(y, w.gridDpi, w.dpi)
		gg.MustOK(win32.MoveToEx(dc.HDC(), 0, drawY, nil))
		gg.MustOK(win32.LineTo(dc.HDC(), win32.INT(metrics.DPIConv(rcClient.Right, w.gridDpi, w.dpi)), drawY))
		win32util.CString(strconv.Itoa(int(y)), &charBuf)
		rect := win32.RECT{}
		gg.Must(win32.DrawTextExW(dc.HDC(), &charBuf[0], -1, &rect, win32.DT_CALCRECT, nil))
		height := rect.Height()
		rect.Top = win32.LONG(drawY) - height/2
		rect.Bottom = rect.Top + height
		gg.Must(win32.DrawTextExW(dc.HDC(), &charBuf[0], -1, &rect, win32.DT_CENTER, nil))
	}
}

func (w *MainWindow) CloneTextFont() *font.Font {
	return w.textFont.Clone()
}

func NewMainWindow(spec *window.Spec) *MainWindow {
	return gw.Init(&MainWindow{Window: window.Window{Spec: spec}})
}

type Window1 struct {
	window.Window
	textFont *font.Font
	charBuf  []win32.WCHAR
}

func (w *Window1) OnInit() {
	w.Window.OnInit()
	win32util.CString("500 X 500", &w.charBuf)
	w.AddMsgListener(win32.WM_DPICHANGED, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
		gg.MustOK(w.textFont.ChangeDPI(gg.Must(w.DPI())))
		w.InvalidateRect(nil, true)
	})
}

func (w *Window1) OnPaint() {
	dc := chkerr.Must(paint.NewPaintDC(w.HWND()))
	defer func() { chkerr.MustOK(dc.EndPaint()) }()
	defer gg.Must(paint.SelectObject(dc.HDC(), w.textFont.HFONT())).Restore()
	rect := gg.Must(w.GetClientRect())
	gg.Must(win32.DrawTextExW(dc.HDC(), &w.charBuf[0], -1, rect, win32.DT_CENTER|win32.DT_SINGLELINE|win32.DT_VCENTER, nil))
}

func (w *Window1) OnDestroy() {
	w.Window.OnDestroy()
	w.textFont.Release()
}

func NewWindow1(spec *window.Spec, textFont *font.Font) *Window1 {
	return gw.Init(&Window1{Window: window.Window{Spec: spec}, textFont: textFont})
}

func main() {
	os.Exit(gw.Run(ui, nil))
}

func ui(app *app.App) {
	var textFont *font.Font

	bkWin := NewMainWindow(&window.Spec{
		Text:      "Full screen",
		Style:     win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
		X:         metrics.Px(win32.CW_USEDEFAULT),
		Y:         metrics.Px(win32.INT(win32.SW_SHOWMAXIMIZED)),
		Width:     metrics.Px(win32.CW_USEDEFAULT),
		OnDestroy: func() { app.Quit(0) },
	})

	ctxMenu := menu.New(true)
	var defDpiMenuItem *menu.Item
	var curDpiMenuItem *menu.Item
	defDpiMenuItem = gg.Must(ctxMenu.InsertItem(-1, &menu.ItemSpec{
		Title:   "Default DPI",
		Checked: false,
		OnClick: func() {
			bkWin.SetGridDpi(win32.USER_DEFAULT_SCREEN_DPI)
			defDpiMenuItem.SetChecked(true)
			curDpiMenuItem.SetChecked(false)
		},
	}))
	curDpiMenuItem = gg.Must(ctxMenu.InsertItem(-1, &menu.ItemSpec{
		Title:   "Current DPI",
		Checked: true,
		OnClick: func() {
			bkWin.SetGridDpi(gg.Must(bkWin.DPI()))
			defDpiMenuItem.SetChecked(false)
			curDpiMenuItem.SetChecked(true)
		},
	}))

	bkWin.SetOnRButtonDownListener(func(event events.MouseClickEvent) {
		gg.MustOK(bkWin.TrackPopupMenu(ctxMenu, nil))
	})

	bkWin.AddMsgListener(win32.WM_DPICHANGED, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
		gg.MustOK(textFont.ChangeDPI(gg.Must(bkWin.DPI())))
		bkWin.InvalidateRect(nil, true)
	})

	gg.Must(win32util.MessageBox(bkWin.HWND(), "Use context menu to change display", "Hint", win32.MB_ICONINFORMATION))

	win1 := window.New(&window.Spec{
		Parent: bkWin.HWND(),
		Text:   "500 X 500",
		Style:  win32.WS_POPUP | win32.WS_CAPTION | win32.WS_VISIBLE,
		X:      metrics.Px(win32.CW_USEDEFAULT),
		Y:      metrics.Px(win32.INT(win32.SW_SHOWNORMAL)),
		Width:  metrics.Dip(500),
		Height: metrics.Dip(500),
	})
	win1.SetOnLButtonDownListener(func(event events.MouseClickEvent) {
		gg.Must(win32.SendMessageW(win1.HWND(), win32.WM_NCLBUTTONDOWN, win32.HTCAPTION, 0))
	})
}
