package layout

import (
	"iter"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/iter2"
)

// Intrinsic is a layout that measures itself to the size of its associated window.
// It is useful to adding a HWND to the layout tree.
type Intrinsic struct {
	// Hwnd is the window to layout.
	// Must be non-zero.
	Hwnd win32.HWND
}

func (w *Intrinsic) HWND() win32.HWND {
	return w.Hwnd
}

func (w *Intrinsic) ChildWidgets() iter.Seq2[int, Widget] {
	return iter2.Empty2[int, Widget]
}

func (w *Intrinsic) CreateElement() (Element, error) {
	if w.Hwnd == 0 {
		return nil, &NoHwndError{Layout: "Intrinsic"}
	}
	return &intrinsicElement{
		BaseElement: BaseElement{
			widget: w,
		},
	}, nil
}

type intrinsicElement struct {
	BaseElement
	layoutSize Size
}

func (e *intrinsicElement) Measure(cst Constraints) (size Size, err error) {
	w := e.Widget().(*Intrinsic)

	var winSize win32.RECT
	if err = win32.GetWindowRect(w.Hwnd, &winSize); err != nil {
		return
	}
	dpi, err := win32.GetDpiForWindow(w.Hwnd)
	size = cst.Clamp(Size{
		Width:  metrics.Px(winSize.Right - winSize.Left).Dip(dpi),
		Height: metrics.Px(winSize.Bottom - winSize.Top).Dip(dpi),
	})
	e.layoutSize = size
	return
}

func (e *intrinsicElement) Arrange(pt Point) (err error) {
	w := e.Widget().(*Intrinsic)
	return positionWindow(w.Hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height)
}
