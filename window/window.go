package window

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const defClassName = "github.com/mkch/gw/wnd_class"

var defClassRegistered = false

type Spec struct {
	ClassName    string
	Text         string
	Style        win32.WINDOW_STYLE
	ExStyle      win32.WINDOW_EX_STYLE
	X            win32.INT
	Y            win32.INT
	Width        win32.INT
	Height       win32.INT
	InCurrentDPI bool // If true, X, Y, Width and Height are in current DPI, otherwise in USER_DEFAULT_SCREEN_DPI.
	WndParent    win32.HWND
	Menu         *menu.Menu
	Instance     win32.HINSTANCE // 0 for this module.
	OnCreate     func()
	OnClose      func()
}

type Window struct {
	WindowBase
	OnCreate func()
	OnClose  func()
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
	var y = spec.Y
	var showCmd win32.SHOW_WINDOW_CMD
	if !spec.InCurrentDPI {
		// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-createwindowexw
		// If an overlapped window is created with the WS_VISIBLE style bit set and the x parameter is set to CW_USEDEFAULT,
		// then the y parameter determines how the window is shown.
		visible = spec.Style&win32.WS_VISIBLE != 0
		style = spec.Style &^ win32.WS_VISIBLE
		showCmd = win32.SW_HIDE
		if visible {
			if spec.X == win32.CW_USEDEFAULT {
				showCmd = gg.If(spec.Y == win32.CW_USEDEFAULT, win32.SW_SHOW, win32.SHOW_WINDOW_CMD(spec.Y))
				y = 0
			}
		}
	}
	hwnd, err := win32util.CreateWindow((*win32util.Wnd)(&win32util.Wnd{
		ClassName:   spec.ClassName,
		WindowName:  spec.Text,
		Style:       style,
		X:           spec.X,
		Y:           y,
		Width:       spec.Width,
		Height:      spec.Height,
		WndParent:   spec.WndParent,
		InParentDPI: false, // DPI of itself will be used instead.
		Instance:    spec.Instance,
	}))
	if err != nil {
		return nil, err
	}

	if !spec.InCurrentDPI {
		// Use the window's own DPI after creating it
		dpi := gg.Must(win32.GetDpiForWindow(hwnd))
		var (
			x, y, cx, cy win32.INT
		)
		var swpFlags win32.UINT = win32.SWP_NOACTIVATE | win32.SWP_NOZORDER
		if spec.X == win32.CW_USEDEFAULT {
			swpFlags |= win32.SWP_NOMOVE
		} else {
			x = win32util.FromDefaultDPI(spec.X, dpi)
			y = win32util.FromDefaultDPI(spec.Y, dpi)
		}

		if spec.Width == win32.CW_USEDEFAULT {
			swpFlags |= win32.SWP_NOSIZE
		} else {
			cx = win32util.FromDefaultDPI(spec.Width, dpi)
			cy = win32util.FromDefaultDPI(spec.Height, dpi)
		}

		if swpFlags&win32.SWP_NOMOVE == 0 || swpFlags&win32.SWP_NOSIZE == 0 {
			win32.SetWindowPos(hwnd, 0, x, y, cx, cy, swpFlags)
		}
		if showCmd != win32.SW_HIDE {
			win32.ShowWindow(hwnd, showCmd)
		}
	}

	win := &Window{OnCreate: spec.OnCreate, OnClose: spec.OnClose}
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
				win.OnClose()
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
