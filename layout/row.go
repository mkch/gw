package layout

import (
	"iter"
	"slices"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
)

// Row is a layout that arranges its children in a row.
type Row struct {
	// Hwnd is the window to layout.
	// It can be 0, in which case Row only arranges its children without resizing the window.
	Hwnd win32.HWND
	// MainAxisAlign is the alignment of children along the main axis(X axis).
	// Only useful when MainAxisSize is AxisSizeMax.
	MainAxisAlign AxisAlignment
	// CrossAxisAlign is the alignment of children along the cross axis(Y axis).
	CrossAxisAlign AxisAlignment
	// MainAxisSize determines how to calculate the width of Row.
	MainAxisSize AxisSize
	// Children are the children to layout.
	Children []Widget
}

func (r *Row) HWND() win32.HWND {
	return r.Hwnd
}

func (r *Row) ChildWidgets() iter.Seq2[int, Widget] {
	return slices.All(r.Children)
}

func (r *Row) CreateElement() (Element, error) {
	return &colRowElement{
		BaseElement: BaseElement{
			widget: r,
		},
		MainSize:       func(s *Size) *metrics.Dip { return &s.Width },
		CrossSize:      func(s *Size) *metrics.Dip { return &s.Height },
		MaxMainSize:    func(c *Constraints) *metrics.Dip { return &c.MaxWidth },
		MinMainSize:    func(c *Constraints) *metrics.Dip { return &c.MinWidth },
		MaxCrossSize:   func(c *Constraints) *metrics.Dip { return &c.MaxHeight },
		MinCrossSize:   func(c *Constraints) *metrics.Dip { return &c.MinHeight },
		MainCoord:      func(p *Point) *metrics.Dip { return &p.X },
		CrossCoord:     func(p *Point) *metrics.Dip { return &p.Y },
		MainAxisSize:   func() AxisSize { return r.MainAxisSize },
		MainAxisAlign:  func() AxisAlignment { return r.MainAxisAlign },
		CrossAxisAlign: func() AxisAlignment { return r.CrossAxisAlign },
	}, nil
}
