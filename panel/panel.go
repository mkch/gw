package panel

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/control"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/paint/brush"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const className = "github.com/mkch/gw/panel_class"

var classRegistered = false

type Spec struct {
	X      metrics.Dimension
	Y      metrics.Dimension
	Width  metrics.Dimension
	Height metrics.Dimension
}

type Panel struct {
	control.Control
	backgroundColor win32.COLORREF
	backgroundBrush *brush.Brush
}

func New(parent win32.HWND, spec *Spec) (*Panel, error) {
	if !classRegistered {
		gg.Must(win32util.RegisterClass(&win32util.WndClass{
			ClassName: className,
			WndProc: func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
				return win32.DefWindowProcW(hwnd, message, wParam, lParam)
			},
			Cursor: gg.Must(win32.LoadImageW_uintptr[win32.HCURSOR](0, uintptr(win32.OCR_NORMAL), win32.IMAGE_CURSOR, 0, 0, win32.LR_DEFAULTSIZE|win32.LR_SHARED)),
		}))
		classRegistered = true
	}
	dpi := gg.Must(win32.GetDpiForWindow(parent))
	hwnd, err := win32util.CreateWindow((&win32util.Wnd{
		ClassName: className,
		Style:     win32.WS_CHILD | win32.WS_VISIBLE,
		X:         spec.X.Px(dpi),
		Y:         spec.Y.Px(dpi),
		Width:     spec.Width.Px(dpi),
		Height:    spec.Height.Px(dpi),
		WndParent: parent,
	}))
	if err != nil {
		return nil, err
	}
	var panel = &Panel{}
	if err := control.Attach(hwnd, &panel.Control); err != nil {
		return nil, err
	}

	if err := panel.SetBackgroundColor(win32.COLORREF(win32.GetSysColor(win32.COLOR_WINDOW))); err != nil {
		return nil, err
	}

	panel.SetWndProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevWndProc win32.WndProc) win32.LRESULT {
		switch message {
		case win32.WM_NCDESTROY:
			if panel.backgroundBrush != nil {
				panel.backgroundBrush.Release()
			}
		case win32.WM_PAINT:
			dc := gg.Must(paint.NewPaintDC(hwnd))
			defer dc.EndPaint()
			win32.FillRect(dc.HDC(), dc.Rect(), panel.backgroundBrush.HBRUSH())
		}
		return prevWndProc(hwnd, message, wParam, lParam)
	})

	return panel, nil
}

func (p *Panel) BackgroundColor() win32.COLORREF {
	return p.backgroundColor
}

func (p *Panel) SetBackgroundColor(color win32.COLORREF) (err error) {
	if p.backgroundBrush != nil {
		p.backgroundBrush.Release()
	}
	p.backgroundColor = color
	p.backgroundBrush, err = brush.New(&win32.LOGBRUSH{
		Style: win32.BS_SOLID,
		Color: p.backgroundColor,
	})
	if err != nil {
		return err
	}

	return win32.InvalidateRect(p.Control.HWND(), nil, true)
}
