package layout

import (
	"iter"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
)

// Padding is a layout that adds padding around its child.
type Padding struct {
	// The padding on each side. Must be non-negative.
	Left, Top, Right, Bottom metrics.Dip
	// Child is the child to layout.
	Child Widget
}

func (p *Padding) HWND() win32.HWND {
	return 0
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

	cst.MinWidth -= w.Left + w.Right
	cst.MinHeight -= w.Top + w.Bottom
	cst.MaxWidth -= w.Left + w.Right
	cst.MaxHeight -= w.Top + w.Bottom
	cst.MaxHeight = max(0, cst.MaxHeight)

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

func (e *paddingElement) Arrange(pt Point) error {
	w := e.Widget().(*Padding)
	child := e.children[0]
	return child.Arrange(Point{
		X: pt.X + w.Left,
		Y: pt.Y + w.Top,
	})
}
