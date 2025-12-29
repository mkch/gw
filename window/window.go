package window

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const defClassName = "github.com/mkch/gw/wnd_class"

var defClassRegistered = false

type Spec struct {
	ClassName string
	Text      string
	Style     win32.WINDOW_STYLE
	ExStyle   win32.WINDOW_EX_STYLE
	X         metrics.Dimension
	Y         metrics.Dimension
	Width     metrics.Dimension
	Height    metrics.Dimension
	WndParent win32.HWND
	Menu      *menu.Menu
	Instance  win32.HINSTANCE // 0 for this module.
	OnCreate  func()
	OnClose   func() bool // Return true to allow closing, false to prevent.
	OnDestroy func()
}

type Window struct {
	WindowBase
	OnCreate  func()
	OnClose   func() bool
	OnDestroy func()
}

func New(spec *Spec) (*Window, error) {
	if !defClassRegistered {
		gg.Must(win32util.RegisterClass(&win32util.WndClass{
			ClassName: defClassName,
			WndProc: func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
				return win32.DefWindowProcW(hwnd, message, wParam, lParam)
			},
			Background: win32.HBRUSH(win32.COLOR_WINDOW + 1),
			Cursor:     gg.Must(win32.LoadImageW_uintptr[win32.HCURSOR](0, uintptr(win32.OCR_NORMAL), win32.IMAGE_CURSOR, 0, 0, win32.LR_DEFAULTSIZE|win32.LR_SHARED)),
		}))
		defClassRegistered = true
	}
	if spec.ClassName == "" {
		copy := *spec
		copy.ClassName = defClassName
		spec = &copy
	}

	var visible bool
	var style = spec.Style
	var showCmd win32.SHOW_WINDOW_CMD = -1
	visible = spec.Style&win32.WS_VISIBLE != 0
	style = spec.Style &^ win32.WS_VISIBLE
	var x, y, cx, cy win32.INT
	var useDefPos, useDefSize bool
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-createwindowexw
	// If an overlapped window is created with the WS_VISIBLE style bit set and the x parameter is set to CW_USEDEFAULT,
	// then the y parameter determines how the window is shown.
	if spec.X.Unit == metrics.PX && spec.X.Value == win32.CW_USEDEFAULT {
		useDefPos = true
		x = win32.CW_USEDEFAULT
		if visible {
			// If an overlapped window is created with the WS_VISIBLE style bit set and the x parameter is set to CW_USEDEFAULT,
			// then the y parameter determines how the window is shown.
			showCmd = gg.If(spec.Y.Value == win32.CW_USEDEFAULT, win32.SW_SHOW, win32.SHOW_WINDOW_CMD(spec.Y.Value))
		}
	}

	// For overlapped windows, if width is CW_USEDEFAULT, the system selects a default width and height for the window
	// and ignores the height.
	if spec.Width.Unit == metrics.PX && spec.Width.Value == win32.CW_USEDEFAULT {
		useDefSize = true
		cx = win32.CW_USEDEFAULT
	}

	hwnd, err := win32util.CreateWindow((&win32util.Wnd{
		ClassName:  spec.ClassName,
		WindowName: spec.Text,
		Style:      style,
		ExStyle:    spec.ExStyle,
		X:          x,
		Y:          y,
		Width:      cx,
		Height:     cy,
		WndParent:  spec.WndParent,
		Instance:   spec.Instance,
	}))
	if err != nil {
		return nil, err
	}

	// Use the window's own DPI after creating it
	dpi := gg.Must(win32.GetDpiForWindow(hwnd))
	var swpFlags win32.UINT = win32.SWP_NOACTIVATE | win32.SWP_NOZORDER
	if useDefPos {
		swpFlags |= win32.SWP_NOMOVE
	} else {
		x = spec.X.Px(dpi)
		y = spec.Y.Px(dpi)
	}

	if useDefSize {
		swpFlags |= win32.SWP_NOSIZE
	} else {
		cx = spec.Width.Px(dpi)
		cy = spec.Height.Px(dpi)
	}

	if !useDefPos || !useDefSize {
		win32.SetWindowPos(hwnd, 0, x, y, cx, cy, swpFlags)
	}
	if showCmd != -1 {
		win32.ShowWindow(hwnd, showCmd)
	}

	win := &Window{OnCreate: spec.OnCreate, OnClose: spec.OnClose, OnDestroy: spec.OnDestroy}
	if err := Attach(hwnd, &win.WindowBase); err != nil {
		return nil, err
	}
	if win.OnCreate != nil {
		win.OnCreate()
	}

	win.SetWndProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevWndProc win32.WndProc) win32.LRESULT {
		switch message {
		case win32.WM_CLOSE:
			if win.OnClose != nil {
				if !win.OnClose() {
					return 0 // prevent closing
				}
			}
		case win32.WM_DESTROY:
			if win.OnDestroy != nil {
				win.OnDestroy()
			}
		}
		return prevWndProc(hwnd, message, wParam, lParam)
	})
	if spec.Menu != nil {
		win.SetMenu(spec.Menu)
	}
	return win, nil
}

func (w *Window) SetMenu(menu *menu.Menu) error {
	return w.setMenu(menu)
}

func (w *Window) preTranslateMessage(p *win32.MSG) bool {
	if w.accelKeyTable == 0 {
		return false
	}
	ok, err := win32.TranslateAcceleratorW(w.hwnd, w.accelKeyTable, p)
	if err != nil {
		panic(err)
	}
	return ok
}
