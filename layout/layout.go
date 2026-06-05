// Package layout provides a DPI-aware layout engine for arranging HWND-backed
// controls in a layout tree.
//
// Layout objects follow a two-phase contract: Measure computes desired size
// under Constraints, and Arrange places the layout and its children at a Point.
// Perform runs this flow for a tree root.
// Geometry and constraints are expressed in DIP (device-independent pixels).
package layout

import (
	"errors"

	"github.com/mkch/gw/metrics"

	"github.com/mkch/gw/win32"
)

// Layout is the interface for layout objects.
type Layout interface {
	Measure(Constraints) (Size, error)
	Arrange(Point) error
	Children() []Layout
	HWND() win32.HWND
}

// PerformWindow calls [Perform] with the given layout and the client size of the given window.
func PerformWindow(root Layout, hwnd win32.HWND) (err error) {
	var size Size
	size, err = ClientSize(hwnd)
	if err != nil {
		return
	}
	return Perform(root, size)
}

// ErrWrongRoot is returned by Perform if the root layout has an associated window.
var ErrWrongRoot = errors.New("root layout must not have an associated window")

// Perform performs the layout: it measures the layout with the given size constraint and arranges the layout at (0, 0).
// The root must not have an associated window.
func Perform(root Layout, size Size) (err error) {
	if root.HWND() != 0 {
		return ErrWrongRoot
	}
	root.Measure(Constraints{
		MaxWidth:  size.Width,
		MaxHeight: size.Height,
	})
	return root.Arrange(Point{X: 0, Y: 0})
}

// Intrinsic is a layout that measures itself to the size of its associated window.
// It is useful to adding a HWND to the layout tree.
type Intrinsic struct {
	// Hwnd is the window to layout.
	// Must be non-zero.
	Hwnd win32.HWND

	layoutSize Size
}

func (w *Intrinsic) HWND() win32.HWND {
	return w.Hwnd
}

func (w *Intrinsic) Children() []Layout {
	return nil
}

func (w *Intrinsic) Measure(cst Constraints) (size Size, err error) {
	var winSize win32.RECT
	if err = win32.GetWindowRect(win32.HWND(w.Hwnd), &winSize); err != nil {
		return
	}
	dpi, err := win32.GetDpiForWindow(win32.HWND(w.Hwnd))
	size = cst.Clamp(Size{
		Width:  metrics.Px(winSize.Right - winSize.Left).Dip(dpi),
		Height: metrics.Px(winSize.Bottom - winSize.Top).Dip(dpi),
	})
	w.layoutSize = size
	return
}

func (w *Intrinsic) Arrange(pt Point) (err error) {
	return positionWindow(w.Hwnd, pt.X, pt.Y, w.layoutSize.Width, w.layoutSize.Height)
}

// Center is a layout that centers its child in the available space.
type Center struct {
	// Hwnd is the window to layout.
	// It can be 0, in which case Center only centers its child without resizing the window.
	Hwnd win32.HWND
	// Item is the child to layout. Must be non-nil.
	Item Layout
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

func (c *Center) Children() []Layout {
	return []Layout{c.Item}
}

func (c *Center) Measure(cst Constraints) (size Size, err error) {
	var childSize Size
	if childSize, err = c.Item.Measure(cst); err != nil {
		return Size{}, err
	}

	if c.WidthFactor == 0 {
		size.Width = cst.MaxWidth
	} else if c.WidthFactor >= 1 {
		size.Width = metrics.Dip(float64(childSize.Width) * c.WidthFactor)
	} else {
		panic("WidthFactor must be 0 or >= 1")
	}
	if c.HeightFactor == 0 {
		size.Height = cst.MaxHeight
	} else if c.HeightFactor >= 1 {
		size.Height = metrics.Dip(float64(childSize.Height) * c.HeightFactor)
	} else {
		panic("HeightFactor must be 0 or >= 1")
	}
	c.childOffset = Point{
		X: (size.Width - childSize.Width) / 2,
		Y: (size.Height - childSize.Height) / 2,
	}

	c.layoutSize = size
	return size, nil
}

func (c *Center) Arrange(pt Point) error {
	defer func() {
		c.childOffset = Point{}
	}()
	if c.Hwnd != 0 {
		if err := positionWindow(c.Hwnd, pt.X, pt.Y, c.layoutSize.Width, c.layoutSize.Height); err != nil {
			return err
		}
	}
	return c.Item.Arrange(Point{
		X: pt.X + c.childOffset.X,
		Y: pt.Y + c.childOffset.Y,
	})
}

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
	// Items are the children to layout.
	Items []Layout

	layoutSize  Size
	itemSizes   []Size
	itemOffsets []Point
}

func (c *Column) HWND() win32.HWND {
	return c.Hwnd
}

func (c *Column) Children() []Layout {
	return c.Items
}

func (c *Column) Measure(cst Constraints) (size Size, err error) {
	var maxWidth metrics.Dip
	var totalHeight metrics.Dip
	var availHeight = cst.MaxHeight
	for _, item := range c.Items {
		var itemSize Size
		itemSize, err = item.Measure(Constraints{
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
		c.itemSizes = append(c.itemSizes, itemSize)
	}

	size.Width = maxWidth
	if c.MainAxisSize == AxisSizeMin {
		size.Height = totalHeight
	} else {
		size.Height = cst.MaxHeight
	}
	c.layoutSize = size

	for i, itemSize := range c.itemSizes {
		var offset Point
		switch c.CrossAxisAlign {
		case AlignStart:
			offset.X = 0
		case AlignCenter:
			offset.X = (maxWidth - itemSize.Width) / 2
		case AlignEnd:
			offset.X = maxWidth - itemSize.Width
		}
		if i == 0 {
			switch c.MainAxisAlign {
			case AlignStart:
				// NOP
			case AlignCenter:
				offset.Y = (size.Height - totalHeight) / 2
			case AlignEnd:
				offset.Y = size.Height - totalHeight
			}
		} else {
			offset.Y = c.itemSizes[i-1].Height + c.itemOffsets[i-1].Y
		}
		c.itemOffsets = append(c.itemOffsets, offset)
	}
	return
}

func (c *Column) Arrange(pt Point) (err error) {
	defer func() {
		c.layoutSize = Size{}
		c.itemSizes = c.itemSizes[:0]
		c.itemOffsets = c.itemOffsets[:0]
	}()

	if c.Hwnd != 0 {
		if err = positionWindow(c.Hwnd, pt.X, pt.Y, c.layoutSize.Width, c.layoutSize.Height); err != nil {
			return err
		}
	}
	for i, item := range c.Items {
		if err = item.Arrange(Point{
			X: pt.X + c.itemOffsets[i].X,
			Y: pt.Y + c.itemOffsets[i].Y,
		}); err != nil {
			return
		}
	}
	return nil
}

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
	// Items are the children to layout.
	Items []Layout

	layoutSize  Size
	itemSizes   []Size
	itemOffsets []Point
}

func (r *Row) HWND() win32.HWND {
	return r.Hwnd
}

func (r *Row) Children() []Layout {
	return r.Items
}

func (r *Row) Measure(cst Constraints) (size Size, err error) {
	var maxHeight metrics.Dip
	var totalWidth metrics.Dip
	var availWidth = cst.MaxWidth
	for _, item := range r.Items {
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
		r.itemSizes = append(r.itemSizes, itemSize)
	}

	size.Height = maxHeight
	if r.MainAxisSize == AxisSizeMin {
		size.Width = totalWidth
	} else {
		size.Width = cst.MaxWidth
	}

	r.layoutSize = size

	for i, itemSize := range r.itemSizes {
		var offset Point
		switch r.CrossAxisAlign {
		case AlignStart:
			// NOP
		case AlignCenter:
			offset.Y = (maxHeight - itemSize.Height) / 2
		case AlignEnd:
			offset.Y = maxHeight - itemSize.Height
		}
		if i == 0 {
			switch r.MainAxisAlign {
			case AlignStart:
				// NOP
			case AlignCenter:
				offset.X = (size.Width - totalWidth) / 2
			case AlignEnd:
				offset.X = size.Width - totalWidth
			}
		} else {
			offset.X = r.itemSizes[i-1].Width + r.itemOffsets[i-1].X
		}
		r.itemOffsets = append(r.itemOffsets, offset)
	}
	return
}

func (r *Row) Arrange(pt Point) (err error) {
	defer func() {
		r.layoutSize = Size{}
		r.itemSizes = r.itemSizes[:0]
		r.itemOffsets = r.itemOffsets[:0]
	}()
	if r.Hwnd != 0 {
		if err = positionWindow(r.Hwnd, pt.X, pt.Y, r.layoutSize.Width, r.layoutSize.Height); err != nil {
			return err
		}
	}

	for i, item := range r.Items {
		if err = item.Arrange(Point{
			X: pt.X + r.itemOffsets[i].X,
			Y: pt.Y + r.itemOffsets[i].Y,
		}); err != nil {
			return
		}
	}
	return nil
}

// Padding is a layout that adds padding around its child.
type Padding struct {
	// The padding on each side. Must be non-negative.
	Left, Top, Right, Bottom metrics.Dip
	// Item is the child to layout.
	Item Layout

	layoutSize Size
}

func (p *Padding) HWND() win32.HWND {
	return 0
}

func (p *Padding) Children() []Layout {
	return []Layout{p.Item}
}

func (p *Padding) Measure(cst Constraints) (size Size, err error) {
	cst.MinWidth -= p.Left + p.Right
	cst.MinHeight -= p.Top + p.Bottom
	cst.MaxWidth -= p.Left + p.Right
	cst.MaxHeight -= p.Top + p.Bottom
	cst.MaxHeight = max(0, cst.MaxHeight)

	var childSize Size
	if childSize, err = p.Item.Measure(cst); err != nil {
		return
	}
	checkOverflow(cst, childSize)

	size.Width = childSize.Width + p.Left + p.Right
	size.Height = childSize.Height + p.Top + p.Bottom
	p.layoutSize = size
	return
}

func (p *Padding) Arrange(pt Point) error {
	return p.Item.Arrange(Point{
		X: pt.X + p.Left,
		Y: pt.Y + p.Top,
	})
}
