package layout

import (
	"iter"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
)

// Center is a layout that centers its child in the available space.
type Center struct {
	// Hwnd is the window to layout.
	// It can be 0, in which case Center only centers its child without resizing the window.
	Hwnd win32.HWND
	// Child is the child to layout. Must be non-nil.
	Child Widget
	// Width scaling factor. If not 0, the desired with of Center is calculated as
	// child's width multiplied by WidthFactor.
	// A 0 WidthFactor means to take all available width from parent.
	// A non-zero WidthFactor must be greater than 1, or it panics.
	WidthFactor float64
	// Height scaling factor. If not 0, the desired height of Center is calculated as
	// child's height multiplied by HeightFactor.
	// A 0 HeightFactor means to take all available height from parent.
	// A non-zero HeightFactor must be greater than 1, or it panics.
	HeightFactor float64

	layoutSize  Size
	childOffset Point
}

func (c *Center) HWND() win32.HWND {
	return c.Hwnd
}

func (c *Center) ChildWidgets() iter.Seq2[int, Widget] {
	return func(yield func(int, Widget) bool) {
		yield(0, c.Child)
	}
}

func (c *Center) CreateElement() (Element, error) {
	if c.Child == nil {
		return nil, &NoChildError{Layout: "Center"}
	}
	if c.WidthFactor != 0 && c.WidthFactor < 1 {
		return nil, &InvalidFactorError{Name: "WidthFactor", Value: c.WidthFactor}
	}

	if c.HeightFactor != 0 && c.HeightFactor < 1 {
		return nil, &InvalidFactorError{Name: "HeightFactor", Value: c.HeightFactor}
	}

	return &centerElement{
		BaseElement: BaseElement{
			widget: c,
		},
	}, nil
}

type centerElement struct {
	BaseElement
	layoutSize  Size
	childOffset Point
}

func (e *centerElement) Measure(cst Constraints) (size Size, err error) {
	w := e.Widget().(*Center)
	child := e.children[0]

	var childSize Size
	if childSize, err = child.Measure(cst); err != nil {
		return Size{}, err
	}

	if w.WidthFactor == 0 {
		size.Width = cst.MaxWidth
	} else if w.WidthFactor >= 1 {
		size.Width = metrics.Dip(float64(childSize.Width) * w.WidthFactor)
	} else {
		panic("WidthFactor must be 0 or >= 1")
	}
	if w.HeightFactor == 0 {
		size.Height = cst.MaxHeight
	} else if w.HeightFactor >= 1 {
		size.Height = metrics.Dip(float64(childSize.Height) * w.HeightFactor)
	} else {
		panic("HeightFactor must be 0 or >= 1")
	}
	e.childOffset = Point{
		X: (size.Width - childSize.Width) / 2,
		Y: (size.Height - childSize.Height) / 2,
	}

	e.layoutSize = size
	return size, nil
}

func (e *centerElement) Arrange(pt Point) error {
	defer func() {
		e.childOffset = Point{}
	}()
	w := e.Widget().(*Center)
	child := e.children[0]
	if w.Hwnd != 0 {
		if err := positionWindow(w.Hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height); err != nil {
			return err
		}
		return child.Arrange(Point{
			X: e.childOffset.X,
			Y: e.childOffset.Y,
		})
	}
	return child.Arrange(Point{
		X: pt.X + e.childOffset.X,
		Y: pt.Y + e.childOffset.Y,
	})
}
