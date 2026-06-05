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
	return &columnElement{
		BaseElement: BaseElement{
			widget: c,
		},
	}, nil
}

type columnElement struct {
	BaseElement
	layoutSize  Size
	itemSizes   []Size
	itemOffsets []Point
}

func (e *columnElement) Measure(cst Constraints) (size Size, err error) {
	w := e.Widget().(*Column)

	var maxWidth metrics.Dip
	var totalHeight metrics.Dip
	var availHeight = cst.MaxHeight
	for _, child := range e.Children() {
		var itemSize Size
		itemSize, err = child.Measure(Constraints{
			MinWidth: cst.MinWidth, MaxWidth: cst.MaxWidth,
			MinHeight: cst.MinHeight, MaxHeight: availHeight})
		if err != nil {
			return
		}

		checkOverflow(cst, itemSize)
		availHeight -= itemSize.Height

		if itemSize.Width > maxWidth {
			maxWidth = itemSize.Width
		}
		totalHeight += itemSize.Height
		e.itemSizes = append(e.itemSizes, itemSize)
	}

	size.Width = maxWidth
	if w.MainAxisSize == AxisSizeMin {
		size.Height = totalHeight
	} else {
		size.Height = cst.MaxHeight
	}
	e.layoutSize = size

	for i, itemSize := range e.itemSizes {
		var offset Point
		switch w.CrossAxisAlign {
		case AlignStart:
			offset.X = 0
		case AlignCenter:
			offset.X = (maxWidth - itemSize.Width) / 2
		case AlignEnd:
			offset.X = maxWidth - itemSize.Width
		}
		if i == 0 {
			switch w.MainAxisAlign {
			case AlignStart:
				// NOP
			case AlignCenter:
				offset.Y = (size.Height - totalHeight) / 2
			case AlignEnd:
				offset.Y = size.Height - totalHeight
			}
		} else {
			offset.Y = e.itemSizes[i-1].Height + e.itemOffsets[i-1].Y
		}
		e.itemOffsets = append(e.itemOffsets, offset)
	}
	return
}

func (e *columnElement) Arrange(pt Point) (err error) {
	defer func() {
		e.layoutSize = Size{}
		e.itemSizes = e.itemSizes[:0]
		e.itemOffsets = e.itemOffsets[:0]
	}()
	w := e.Widget().(*Column)
	if w.Hwnd != 0 {
		if err = positionWindow(w.Hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height); err != nil {
			return err
		}
		for i, item := range e.children {
			if err = item.Arrange(Point{
				e.itemOffsets[i].X,
				e.itemOffsets[i].Y,
			}); err != nil {
				return
			}
		}
		return nil
	}
	for i, item := range e.children {
		if err = item.Arrange(Point{
			X: pt.X + e.itemOffsets[i].X,
			Y: pt.Y + e.itemOffsets[i].Y,
		}); err != nil {
			return
		}
	}
	return nil
}
