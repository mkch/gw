package layout

import (
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
)

type Layout interface {
	Measure(Constraints) (Size, error)
	Arrange(Point) error
	Children() []Layout
	HWND() win32.HWND
}

func Perform(l Layout, size *Size) (err error) {
	if size == nil {
		var clientRect win32.RECT
		if err = win32.GetClientRect(l.HWND(), &clientRect); err != nil {
			return err
		}
		var dpi win32.UINT
		dpi, err = win32.GetDpiForWindow(l.HWND())
		if err != nil {
			return
		}
		size = &Size{
			Width:  PxToDP(clientRect.Right-clientRect.Left, dpi),
			Height: PxToDP(clientRect.Bottom-clientRect.Top, dpi),
		}
	}
	l.Measure(Constraints{
		MaxWidth:  size.Width,
		MaxHeight: size.Height,
	})
	return l.Arrange(Point{X: 0, Y: 0})
}

type Window struct {
	Hwnd       win32.HWND
	layoutSize Size
}

func (w *Window) HWND() win32.HWND {
	return w.Hwnd
}

func (ww *Window) Children() []Layout {
	return nil
}

func (w *Window) Measure(cst Constraints) (size Size, err error) {
	var winSize win32.RECT
	if err = win32.GetWindowRect(win32.HWND(w.Hwnd), &winSize); err != nil {
		return
	}
	dpi, err := win32.GetDpiForWindow(win32.HWND(w.Hwnd))
	size = cst.Clamp(Size{
		Width:  PxToDP(winSize.Right-winSize.Left, dpi),
		Height: PxToDP(winSize.Bottom-winSize.Top, dpi),
	})
	w.layoutSize = size
	return
}

func (w *Window) Arrange(pt Point) (err error) {
	return setWindowPos(w.Hwnd, pt.X, pt.Y, w.layoutSize.Width, w.layoutSize.Height)
}

type Center struct {
	Hwnd win32.HWND
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
		if err := setWindowPos(c.Hwnd, pt.X, pt.Y, c.layoutSize.Width, c.layoutSize.Height); err != nil {
			return err
		}
	}
	return c.Item.Arrange(Point{
		X: pt.X + c.childOffset.X,
		Y: pt.Y + c.childOffset.Y,
	})
}

type AxisAlignment int

const (
	AlignStart AxisAlignment = iota
	AlignCenter
	AlignEnd
)

type AxisSize int

const (
	AxisSizeMax AxisSize = iota
	AxisSizeMin
)

type Column struct {
	Hwnd           win32.HWND
	MainAxisAlign  AxisAlignment
	CrossAxisAlign AxisAlignment
	MainAxisSize   AxisSize
	Items          []Layout

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
		if err = setWindowPos(c.Hwnd, pt.X, pt.Y, c.layoutSize.Width, c.layoutSize.Height); err != nil {
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

type Row struct {
	Hwnd           win32.HWND
	MainAxisAlign  AxisAlignment
	CrossAxisAlign AxisAlignment
	MainAxisSize   AxisSize
	Items          []Layout

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
		if err = setWindowPos(r.Hwnd, pt.X, pt.Y, r.layoutSize.Width, r.layoutSize.Height); err != nil {
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
