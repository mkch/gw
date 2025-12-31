package paint

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/win32"
)

// DC is a readonly wrapper of win32.HDC.
type DC struct {
	hdc win32.HDC
}

func (dc DC) HDC() win32.HDC {
	return dc.hdc
}

func SelectObject[T win32.HGDIOBJ](hdc win32.HDC, object T) (struct{ Restore func() }, error) {
	if oldObject, err := win32.SelectObject(hdc, object); err != nil {
		return struct{ Restore func() }{}, err
	} else {
		return struct{ Restore func() }{func() {
			gg.Must(win32.SelectObject(hdc, oldObject))
		}}, nil
	}
}

type PaintDC struct {
	DC
	hwnd win32.HWND
	p    win32.PAINTSTRUCT
}

func (dc *PaintDC) EraseBackground() bool {
	return dc.p.Erase != 0
}

func (dc *PaintDC) Rect() *win32.RECT {
	return &dc.p.RcPaint
}

func (dc *PaintDC) EndPaint() error {
	return win32.EndPaint(dc.hwnd, &dc.p)
}

func NewPaintDC(hwnd win32.HWND) (dc *PaintDC, err error) {
	dc = &PaintDC{}
	var hdc win32.HDC
	if hdc, err = win32.BeginPaint(hwnd, &dc.p); err != nil {
		return
	}
	dc.DC.hdc = hdc
	dc.hwnd = hwnd
	return
}

type PaintData struct {
	DC    win32.HDC
	Erase bool
	Rect  win32.RECT
}
