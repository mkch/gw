package layout

import (
	"iter"

	"github.com/mkch/gg"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/iter2"
)

type Sized struct {
	// Hwnd is the window to layout.
	// It can be 0, in which case Sized only sizes its child without resizing the window.
	Hwnd win32.HWND
	// Child is the child to layout.
	// Can be nil, in which case Sized just takes the specified size.
	Child  Widget
	Width  metrics.Dip
	Height metrics.Dip
}

func (s *Sized) HWND() win32.HWND {
	return s.Hwnd
}

func (s *Sized) ChildWidgets() iter.Seq2[int, Widget] {
	return gg.If(s.Child == nil,
		iter2.Empty2[int, Widget],
		func(yield func(int, Widget) bool) { yield(0, s.Child) },
	)
}

func (s *Sized) CreateElement() (Element, error) {
	return &sizedElement{
		BaseElement: BaseElement{
			widget: s,
		},
	}, nil
}

type sizedElement struct {
	BaseElement
	layoutSize Size
}

func (e *sizedElement) Measure(cst Constraints) (size Size, err error) {
	w := e.Widget().(*Sized)
	size = cst.Clamp(Size{Width: w.Width, Height: w.Height})
	e.layoutSize = size
	if len(e.children) == 0 {
		return
	}
	child := e.children[0]

	childCst := Constraints{
		MinWidth:  size.Width,
		MaxWidth:  size.Width,
		MinHeight: size.Height,
		MaxHeight: size.Height,
	}
	childSize, err := child.Measure(childCst)
	if err != nil {
		return
	}
	checkOverflow(childCst, childSize)
	return
}

func (e *sizedElement) Arrange(pt Point) (err error) {
	w := e.Widget().(*Sized)
	arrangeChild := gg.If(len(e.children) == 0,
		func(Point) error { return nil },
		func(pt Point) error { return e.children[0].Arrange(pt) },
	)
	if w.Hwnd != 0 {
		if err = positionWindow(w.Hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height); err != nil {
			return
		}
		return arrangeChild(Point{X: 0, Y: 0})
	}

	return arrangeChild(pt)
}
