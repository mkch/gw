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
	return &rowElement{
		BaseElement: BaseElement{
			widget: r,
		},
	}, nil
}

type rowElement struct {
	BaseElement
	layoutSize  Size
	itemSizes   []Size
	itemOffsets []Point
}

func (e *rowElement) Measure(cst Constraints) (size Size, err error) {
	w := e.Widget().(*Row)

	var maxHeight metrics.Dip
	var totalWidth metrics.Dip
	var availWidth = cst.MaxWidth
	for _, item := range e.children {
		var itemSize Size
		itemSize, err = item.Measure(Constraints{
			MinWidth: cst.MinWidth, MaxWidth: availWidth,
			MinHeight: cst.MinHeight, MaxHeight: cst.MaxHeight})
		if err != nil {
			return
		}
		checkOverflow(cst, itemSize)
		availWidth -= itemSize.Width

		if itemSize.Height > maxHeight {
			maxHeight = itemSize.Height
		}
		totalWidth += itemSize.Width
		e.itemSizes = append(e.itemSizes, itemSize)
	}

	size.Height = maxHeight
	if w.MainAxisSize == AxisSizeMin {
		size.Width = totalWidth
	} else {
		size.Width = cst.MaxWidth
	}

	e.layoutSize = size

	for i, itemSize := range e.itemSizes {
		var offset Point
		switch w.CrossAxisAlign {
		case AlignStart:
			// NOP
		case AlignCenter:
			offset.Y = (maxHeight - itemSize.Height) / 2
		case AlignEnd:
			offset.Y = maxHeight - itemSize.Height
		}
		if i == 0 {
			switch w.MainAxisAlign {
			case AlignStart:
				// NOP
			case AlignCenter:
				offset.X = (size.Width - totalWidth) / 2
			case AlignEnd:
				offset.X = size.Width - totalWidth
			}
		} else {
			offset.X = e.itemSizes[i-1].Width + e.itemOffsets[i-1].X
		}
		e.itemOffsets = append(e.itemOffsets, offset)
	}
	return
}

func (e *rowElement) Arrange(pt Point) (err error) {
	defer func() {
		e.layoutSize = Size{}
		e.itemSizes = e.itemSizes[:0]
		e.itemOffsets = e.itemOffsets[:0]
	}()
	w := e.Widget().(*Row)
	if w.Hwnd != 0 {
		if err = positionWindow(w.Hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height); err != nil {
			return err
		}

		for i, child := range e.children {
			if err = child.Arrange(Point{
				X: e.itemOffsets[i].X,
				Y: e.itemOffsets[i].Y,
			}); err != nil {
				return
			}
		}
		return nil
	}

	for i, child := range e.children {
		if err = child.Arrange(Point{
			X: pt.X + e.itemOffsets[i].X,
			Y: pt.Y + e.itemOffsets[i].Y,
		}); err != nil {
			return
		}
	}
	return nil
}
