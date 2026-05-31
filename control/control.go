package control

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw"
	"github.com/mkch/gw/paint/font"
	"github.com/mkch/gw/win32"
)

type Control struct {
	gw.BaseWindowImpl
	font *font.Font
}

func (c *Control) OnInit() {
	c.font = gg.Must(font.New(font.SysDefault(), gg.Must(win32.GetDpiForWindow(c.HWND()))))
	c.applyFont()
}

func (c *Control) WndProc(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	switch message {
	case win32.WM_NCDESTROY:
		c.font.Release()
	case win32.WM_DPICHANGED_AFTERPARENT:
		dpi := gg.Must(win32.GetDpiForWindow(c.HWND()))
		gg.MustOK(c.font.ChangeDPI(dpi))
		c.applyFont()
	}
	return c.BaseWindowImpl.WndProc(hwnd, message, wParam, lParam)
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
