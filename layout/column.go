package layout

import (
	"iter"
	"slices"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
)

// AxisAlignment represents alignment along an axis.
type AxisAlignment int

const (
	// AlignStart means to align to the start of the axis.
	AlignStart AxisAlignment = iota
	// AlignCenter means to align to the center of the axis.
	AlignCenter
	// AlignEnd means to align to the end of the axis.
	AlignEnd
)

// AxisSize represents how to determine the size of an axis.
type AxisSize int

const (
	// AxisSizeMax means to take the maximum available size along the axis.
	AxisSizeMax AxisSize = iota
	// AxisSizeMin means to take the minimum size required by the content along the axis.
	AxisSizeMin
)

// Column is a layout that arranges its children in a column.
type Column struct {
	// Hwnd is the window to layout.
	// It can be 0, in which case Column only arranges its children without resizing the window.
	Hwnd win32.HWND
	// MainAxisAlign is the alignment of children along the main axis(Y axis).
	// Only useful when MainAxisSize is AxisSizeMax.
	MainAxisAlign AxisAlignment
	// CrossAxisAlign is the alignment of children along the cross axis(X axis).
	CrossAxisAlign AxisAlignment
	// MainAxisSize determines how to calculate the height of Column.
	// If MainAxisSize is AxisSizeMax, the height of Column is the maximum available height from parent.
	// If MainAxisSize is AxisSizeMin, the height of Column is the total height of its children.
	MainAxisSize AxisSize
	// Children are the children to layout.
	Children []Widget
}

func (c *Column) HWND() win32.HWND {
	return c.Hwnd
}

func (c *Column) ChildWidgets() iter.Seq2[int, Widget] {
	return slices.All(c.Children)
}

func (c *Column) CreateElement() (Element, error) {
	return &colRowElement{
		BaseElement: BaseElement{
			widget: c,
		},
		MainSize:       func(s *Size) *metrics.Dip { return &s.Height },
		CrossSize:      func(s *Size) *metrics.Dip { return &s.Width },
		MaxMainSize:    func(c *Constraints) *metrics.Dip { return &c.MaxHeight },
		MinMainSize:    func(c *Constraints) *metrics.Dip { return &c.MinHeight },
		MaxCrossSize:   func(c *Constraints) *metrics.Dip { return &c.MaxWidth },
		MinCrossSize:   func(c *Constraints) *metrics.Dip { return &c.MinWidth },
		MainCoord:      func(p *Point) *metrics.Dip { return &p.Y },
		CrossCoord:     func(p *Point) *metrics.Dip { return &p.X },
		MainAxisSize:   func() AxisSize { return c.MainAxisSize },
		MainAxisAlign:  func() AxisAlignment { return c.MainAxisAlign },
		CrossAxisAlign: func() AxisAlignment { return c.CrossAxisAlign },
	}, nil
}
