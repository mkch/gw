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

	var rect win32.RECT
	if err := win32.GetWindowRect(w.Hwnd, &rect); err != nil {
		return nil, err
	}
	dpi, err := win32.GetDpiForWindow(w.Hwnd)
	if err != nil {
		return nil, err
	}

	return &intrinsicElement{
		BaseElement: BaseElement{
			widget: w,
		},
		intrinsicSize: Size{
			Width:  metrics.Px(rect.Width()).Dip(dpi),
			Height: metrics.Px(rect.Height()).Dip(dpi),
		},
	}, nil
}

type intrinsicElement struct {
	BaseElement

	intrinsicSize Size
	layoutSize    Size
}

func (e *intrinsicElement) Measure(cst Constraints) (size Size, err error) {
	size = cst.Clamp(e.intrinsicSize)
	e.layoutSize = size
	return
}

func (e *intrinsicElement) Arrange(pt Point) (err error) {
	w := e.Widget().(*Intrinsic)
	return positionWindow(w.Hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height)
}
