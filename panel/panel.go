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

const className = "github.com/mkch/gw/panel"

var classAtom win32.ATOM

var customClassNames gg.Set[string]

type Spec struct {
	ClassName string // Custom class name. If empty, the default class will be used.
	X         metrics.Dimension
	Y         metrics.Dimension
	Width     metrics.Dimension
	Height    metrics.Dimension
	ExStyle   win32.WINDOW_EX_STYLE
}

type Panel struct {
	control.Control
	backgroundColor win32.COLORREF
	backgroundBrush *brush.Brush
}

func New(parent win32.HWND, spec *Spec) (*Panel, error) {
	if classAtom == 0 {
		classAtom = gg.Must(win32util.RegisterClass(&win32util.WndClass{
			ClassName: className,
			WndProc: func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
				return win32.DefWindowProcW(hwnd, message, wParam, lParam)
			},
			Cursor: gg.Must(win32.LoadImageW_uintptr[win32.HCURSOR](0, uintptr(win32.OCR_NORMAL), win32.IMAGE_CURSOR, 0, 0, win32.LR_DEFAULTSIZE|win32.LR_SHARED)),
		}))
	}
	if spec.ClassName == "" { // Use default class name.
		spec = new(*spec)
		spec.ClassName = className
	} else { // Use custom class name.
		if customClassNames == nil {
			customClassNames = make(gg.Set[string])
		}
		if !customClassNames.Contains(spec.ClassName) {
			if _, err := win32util.CopyWindowClass(classAtom, spec.ClassName); err != nil {
				return nil, err
			}
			customClassNames.Add(spec.ClassName)
		}
	}

	dpi := gg.Must(win32.GetDpiForWindow(parent))
	hwnd, err := win32util.CreateWindow((&win32util.Wnd{
		ClassName: spec.ClassName,
		Style:     win32.WS_CHILD | win32.WS_VISIBLE,
		ExStyle:   spec.ExStyle,
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
		}
		return prevWndProc(hwnd, message, wParam, lParam)
	})
	panel.AddPaintCallback(func(paintData *paint.PaintData, prev func(*paint.PaintData)) {
		prev(paintData)
		win32.FillRect(paintData.DC, &paintData.Rect, panel.backgroundBrush.HBRUSH())
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
