package button

import (
	"github.com/mkch/gw"
	"github.com/mkch/gw/control"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

type Button struct {
	control.Control
	// Spec is used to create the window and is set to nil after creation.
	Spec *Spec

	onClickListener func()
}

func (b *Button) SetOnClickListener(listener func()) {
	b.onClickListener = listener
}

func (b *Button) callOnClickListener() {
	if b.onClickListener != nil {
		b.onClickListener()
	}
}

func (b *Button) SetWindowText(str string) error {
	return win32util.SetWindowText(b.HWND(), str)
}

func (b *Button) GetWindowText() (string, error) {
	return win32util.GetWindowText(b.HWND())
}

type Spec struct {
	Parent  gw.BaseWindow
	Text    string
	OnClick func()
	X       metrics.Dimension
	Y       metrics.Dimension
	Width   metrics.Dimension
	Height  metrics.Dimension
	Style   win32.WindowStyle
	ExStyle win32.WindowExStyle
}

func (b *Button) OnInit() error {
	defer func() { b.Spec = nil }()
	if err := b.Control.OnInit(); err != nil {
		return err
	}
	if b.Spec.OnClick != nil {
		b.SetOnClickListener(b.Spec.OnClick)
	}
	return nil
}

// OnClick is called when the button is clicked.
// The default implementation calls the on-click listener if it is set.
func (b *Button) OnClick() {
	b.callOnClickListener()
}

func (b *Button) WndProc(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	switch message {
	case appmsg.REFLECT_COMMAND:
		gw.LookupWindow(b.HWND()).(interface{ OnClick() }).OnClick()
		return 0
	}
	return b.Control.WndProc(hwnd, message, wParam, lParam)
}

func (b *Button) CreateHandle() (win32.HWND, error) {
	dpi, err := win32.GetDpiForWindow(b.Spec.Parent.HWND())
	if err != nil {
		return 0, err
	}
	return win32util.CreateWindow(&win32util.Wnd{
		ClassName:  "BUTTON",
		WndParent:  b.Spec.Parent.HWND(),
		WindowName: b.Spec.Text,
		X:          metrics.ToPx(b.Spec.X, dpi).Value(),
		Y:          metrics.ToPx(b.Spec.Y, dpi).Value(),
		Width:      metrics.ToPx(b.Spec.Width, dpi).Value(),
		Height:     metrics.ToPx(b.Spec.Height, dpi).Value(),
		Style:      b.Spec.Style | win32.WS_CHILD,
		ExStyle:    b.Spec.ExStyle,
	})
}

// New creates a new Button control with the specified specification.
func New(spec *Spec) (*Button, error) {
	return gw.Init(&Button{Spec: spec})
}
