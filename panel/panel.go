package panel

import (
	"sync"

	"github.com/mkch/gg"
	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/control"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/paint/brush"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const defClassName = "github.com/mkch/gw#Panel"

var defClassSpec = win32util.WndClass{
	WndProc: func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
		return win32.DefWindowProcW(hwnd, message, wParam, lParam)
	},
	Cursor: gg.Must(win32.LoadImageW_uintptr[win32.HCURSOR](0, uintptr(win32.OCR_NORMAL), win32.IMAGE_CURSOR, 0, 0, win32.LR_DEFAULTSIZE|win32.LR_SHARED)),
}

// classRegistered is a function that returns two functions: setRegistered and registered.
// setRegistered marks a class name as registered, and registered checks if a class name is marked as registered.
// The returned functions are concurrency-safe.
var classRegistered = sync.OnceValues(func() (func(string), func(string) bool) {
	var lock sync.RWMutex
	var registeredClasses = make(gg.Set[string])
	return func(s string) {
			lock.Lock()
			defer lock.Unlock()
			registeredClasses.Add(s)
		}, func(s string) bool {
			lock.RLock()
			defer lock.RUnlock()
			return registeredClasses.Contains(s)
		}
})

type Spec struct {
	Parent    gw.WindowParent
	ClassName string // Custom class name. If empty, the default class will be used.
	X         metrics.Dimension
	Y         metrics.Dimension
	Width     metrics.Dimension
	Height    metrics.Dimension
	ExStyle   win32.WINDOW_EX_STYLE
}

type Panel struct {
	control.Control
	// Spec is used to create the window and is cleared after creation.
	Spec            *Spec
	backgroundColor win32.COLORREF
	backgroundBrush *brush.Brush
}

func (p *Panel) OnInit() {
	defer func() { p.Spec = nil }()
	p.Control.OnInit()
	gg.MustOK(p.SetBackgroundColor(win32.COLORREF(win32.GetSysColor(win32.COLOR_WINDOW))))
}

func (p *Panel) WndProc(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	switch message {
	case win32.WM_NCDESTROY:
		if p.backgroundBrush != nil {
			p.backgroundBrush.Release()
		}
	}
	return p.Control.WndProc(hwnd, message, wParam, lParam)
}

func (p *Panel) OnPaint() {
	dc := gg.Must(paint.NewPaintDC(p.HWND()))
	defer func() { gg.MustOK(dc.EndPaint()) }()
	win32.FillRect(dc.HDC(), dc.Rect(), p.backgroundBrush.HBRUSH())
}

func (p *Panel) CreateHandle() win32.HWND {
	if p.Spec == nil {
		p.Spec = &Spec{}
	}

	className := p.Spec.ClassName
	if className == "" {
		// Use default class name.
		className = defClassName
	}
	setRegistered, registered := classRegistered()
	if !registered(className) {
		cls := new(defClassSpec)
		cls.ClassName = className
		chkerr.Must(win32util.RegisterClass(cls))
		setRegistered(className)
	}

	dpi := gg.Must(win32.GetDpiForWindow(p.Spec.Parent.HWND()))
	return chkerr.Must(win32util.CreateWindow(&win32util.Wnd{
		ClassName: className,
		Style:     win32.WS_CHILD | win32.WS_VISIBLE,
		ExStyle:   p.Spec.ExStyle,
		X:         p.Spec.X.Px(dpi),
		Y:         p.Spec.Y.Px(dpi),
		Width:     p.Spec.Width.Px(dpi),
		Height:    p.Spec.Height.Px(dpi),
		WndParent: p.Spec.Parent.HWND(),
	}))
}

func New(spec *Spec) (panel *Panel) {
	panel = &Panel{Spec: spec}
	gw.Init(panel)
	return
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

	return win32.InvalidateRect(p.HWND(), nil, true)
}
