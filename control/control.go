package control

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/paint/font"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

type Control struct {
	window.WindowBase
	font *font.Font
}

func Attach(hwnd win32.HWND, control *Control) error {
	if err := window.Attach(hwnd, &control.WindowBase); err != nil {
		return err
	}
	if control.font == nil {
		control.font = gg.Must(font.New(font.SysDefault(), gg.Must(win32.GetDpiForWindow(hwnd))))
		control.applyFont()
	}
	control.SetWndProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevWndProc win32.WndProc) win32.LRESULT {
		switch message {
		case win32.WM_NCDESTROY:
			control.font.Release()
		case win32.WM_DPICHANGED_AFTERPARENT:
			dpi := gg.Must(win32.GetDpiForWindow(control.HWND()))
			gg.MustOK(control.font.ChangeDPI(dpi))
			control.applyFont()
		}
		return prevWndProc(hwnd, message, wParam, lParam)
	})
	return nil
}

func (ctrl *Control) applyFont() {
	win32.SendMessageW(ctrl.HWND(), win32.WM_SETFONT, win32.WPARAM(ctrl.font.HFONT()), 1)
}

// SetFont sets the font used by this control. System default font is used if font is nil.
func (ctrl *Control) SetFont(f *font.Font) {
	if ctrl.font != nil {
		ctrl.font.Release()
	}
	if f == nil {
		ctrl.font = gg.Must(font.New(font.SysDefault(), gg.Must(win32.GetDpiForWindow(ctrl.HWND()))))
	} else {
		ctrl.font = f.Clone()
	}
	ctrl.applyFont()
}
