package layout

import (
	"iter"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
)

// Padding is a layout that adds padding around its child.
type Padding struct {
	// Hwnd is the window to layout.
	// It can be 0, in which case Padding only sizes its child without resizing the window.
	Hwnd win32.HWND
	// The padding on each side. Must be non-negative.
	Left, Top, Right, Bottom metrics.Dip
	// Child is the child to layout. Must be non-nil.
	Child Widget
}

func (p *Padding) HWND() win32.HWND {
	return p.Hwnd
}

func (p *Padding) ChildWidgets() iter.Seq2[int, Widget] {
	return func(yield func(int, Widget) bool) {
		yield(0, p.Child)
	}
}

func (p *Padding) CreateElement() (Element, error) {
	return &paddingElement{
		BaseElement: BaseElement{
			widget: p,
		},
	}, nil
}

type paddingElement struct {
	BaseElement
	layoutSize Size
}

func (e *paddingElement) Measure(cst Constraints) (size Size, err error) {
	w := e.Widget().(*Padding)
	child := e.children[0]

	cst.MinWidth = max(0, cst.MinWidth-w.Left-w.Right)
	cst.MinHeight = max(0, cst.MinHeight-w.Top-w.Bottom)
	cst.MaxWidth = max(0, cst.MaxWidth-w.Left-w.Right)
	cst.MaxHeight = max(0, cst.MaxHeight-w.Top-w.Bottom)

	var childSize Size
	if childSize, err = child.Measure(cst); err != nil {
		return
	}
	checkOverflow(cst, childSize)

	size.Width = childSize.Width + w.Left + w.Right
	size.Height = childSize.Height + w.Top + w.Bottom
	e.layoutSize = size
	return
}

func (e *paddingElement) Arrange(pt Point) (err error) {
	w := e.Widget().(*Padding)
	child := e.children[0]
	if w.Hwnd != 0 {
		if err = positionWindow(w.Hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height); err != nil {
			return
		}
		return child.Arrange(Point{
			X: w.Left,
			Y: w.Top,
		})
	}
	return child.Arrange(Point{
		X: pt.X + w.Left,
		Y: pt.Y + w.Top,
	})
}
