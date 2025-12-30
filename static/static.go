package static

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/control"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/paint/brush"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const (
	SS_LEFT            win32.WINDOW_STYLE = 0x00000000
	SS_CENTER          win32.WINDOW_STYLE = 0x00000001
	SS_RIGHT           win32.WINDOW_STYLE = 0x00000002
	SS_ICON            win32.WINDOW_STYLE = 0x00000003
	SS_BLACKRECT       win32.WINDOW_STYLE = 0x00000004
	SS_GRAYRECT        win32.WINDOW_STYLE = 0x00000005
	SS_WHITERECT       win32.WINDOW_STYLE = 0x00000006
	SS_BLACKFRAME      win32.WINDOW_STYLE = 0x00000007
	SS_GRAYFRAME       win32.WINDOW_STYLE = 0x00000008
	SS_WHITEFRAME      win32.WINDOW_STYLE = 0x00000009
	SS_USERITEM        win32.WINDOW_STYLE = 0x0000000A
	SS_SIMPLE          win32.WINDOW_STYLE = 0x0000000B
	SS_LEFTNOWORDWRAP  win32.WINDOW_STYLE = 0x0000000C
	SS_OWNERDRAW       win32.WINDOW_STYLE = 0x0000000D
	SS_BITMAP          win32.WINDOW_STYLE = 0x0000000E
	SS_ENHMETAFILE     win32.WINDOW_STYLE = 0x0000000F
	SS_ETCHEDHORZ      win32.WINDOW_STYLE = 0x00000010
	SS_ETCHEDVERT      win32.WINDOW_STYLE = 0x00000011
	SS_ETCHEDFRAME     win32.WINDOW_STYLE = 0x00000012
	SS_TYPEMASK        win32.WINDOW_STYLE = 0x0000001F
	SS_REALSIZECONTROL win32.WINDOW_STYLE = 0x00000040
	SS_NOPREFIX        win32.WINDOW_STYLE = 0x00000080
	SS_NOTIFY          win32.WINDOW_STYLE = 0x00000100
	SS_CENTERIMAGE     win32.WINDOW_STYLE = 0x00000200
	SS_RIGHTJUST       win32.WINDOW_STYLE = 0x00000400
	SS_REALSIZEIMAGE   win32.WINDOW_STYLE = 0x00000800
	SS_SUNKEN          win32.WINDOW_STYLE = 0x00001000
	SS_EDITCONTROL     win32.WINDOW_STYLE = 0x00002000
	SS_ENDELLIPSIS     win32.WINDOW_STYLE = 0x00004000
	SS_PATHELLIPSIS    win32.WINDOW_STYLE = 0x00008000
	SS_WORDELLIPSIS    win32.WINDOW_STYLE = 0x0000C000
	SS_ELLIPSISMASK    win32.WINDOW_STYLE = 0x0000C000
)

type Static struct {
	control.Control
	backgroundColor win32.COLORREF
	backgroundBrush *brush.Brush // can't bi nil
}

// BackgroundColor returns the background color of the Static control.
func BackgroundColor(static *Static) *win32.COLORREF {
	return &static.backgroundColor
}

// SetBackgroundColor sets the background color of the Static control.
func (static *Static) SetBackgroundColor(color win32.COLORREF) error {
	if static.backgroundColor == color {
		return nil
	}

	if static.backgroundBrush != nil {
		static.backgroundBrush.Release()
	}

	static.backgroundColor = color
	var err error
	if static.backgroundBrush, err = createBackgroundBrush(color); err != nil {
		return err
	}
	return win32.InvalidateRect(static.HWND(), nil, true)

}

// createBackgroundBrush is a helper function that creates a solid brush with the specified color.
func createBackgroundBrush(color win32.COLORREF) (*brush.Brush, error) {
	return brush.New(&win32.LOGBRUSH{Style: win32.BS_SOLID, Color: color})
}

type Spec struct {
	Text    string
	X       metrics.Dimension
	Y       metrics.Dimension
	Width   metrics.Dimension
	Height  metrics.Dimension
	Style   win32.WINDOW_STYLE
	ExStyle win32.WINDOW_EX_STYLE
}

func New(parent win32.HWND, spec *Spec) (*Static, error) {
	dpi := gg.Must(win32.GetDpiForWindow(parent))
	hwnd, err := win32util.CreateWindow(&win32util.Wnd{
		ClassName:  "STATIC",
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
	var static Static
	if err := control.Attach(hwnd, &static.Control); err != nil {
		return nil, err
	}

	static.SetBackgroundColor(win32.COLORREF(win32.GetSysColor(win32.COLOR_WINDOW)))
	static.SetWndProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevWndProc win32.WndProc) win32.LRESULT {
		switch message {
		case win32.WM_DESTROY:
			static.backgroundBrush.Release()
		case appmsg.REFLECT_CTLCOLORSTATIC:
			win32.SetBkMode(win32.HDC(wParam), win32.TRANSPARENT) // The *text* background is transparent.
			return win32.LRESULT(static.backgroundBrush.HBRUSH())
		}
		return prevWndProc(hwnd, message, wParam, lParam)
	})
	return &static, nil
}
