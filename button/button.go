package button

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/control"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

type Button struct {
	control.Control
	OnClick func()
}

func (b *Button) SetWindowText(str string) error {
	return win32util.SetWindowText(b.HWND(), str)
}

func (b *Button) GetWindowText() (string, error) {
	return win32util.GetWindowText(b.HWND())
}

type Spec struct {
	Text    string
	OnClick func()
	X       metrics.Dimension
	Y       metrics.Dimension
	Width   metrics.Dimension
	Height  metrics.Dimension
	Style   win32.WINDOW_STYLE
	ExStyle win32.WINDOW_EX_STYLE
}

func New(parent win32.HWND, spec *Spec) (*Button, error) {
	dpi := gg.Must(win32.GetDpiForWindow(parent))
	hwnd, err := win32util.CreateWindow(&win32util.Wnd{
		ClassName:  "BUTTON",
		WndParent:  parent,
		WindowName: spec.Text,
		X:          spec.X.Px(dpi),
		Y:          spec.Y.Px(dpi),
		Width:      spec.Width.Px(dpi),
		Height:     spec.Height.Px(dpi),
		Style:      spec.Style | win32.WS_CHILD,
		ExStyle:    spec.ExStyle,
	})
	if err != nil {
		return nil, err
	}
	var button = Button{OnClick: spec.OnClick}
	if err := control.Attach(hwnd, &button.Control); err != nil {
		return nil, err
	}
	button.SetWndProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prev win32.WndProc) win32.LRESULT {
		switch message {
		case appmsg.REFLECT_COMMAND:
			if button.OnClick != nil {
				button.OnClick()
			}
		}
		return prev(hwnd, message, wParam, lParam)
	})
	return &button, nil
}
