package layout

import "github.com/mkch/gw/metrics"

// colRowElement is the element of Col and Row.
type colRowElement struct {
	BaseElement

	// MainSize returns the size of the main axis.
	MainSize func(*Size) *metrics.Dip
	// CrossSize returns the size of the cross axis.
	CrossSize func(*Size) *metrics.Dip
	// MaxMainSize returns the maximum size of the main axis.
	MaxMainSize func(*Constraints) *metrics.Dip
	// MinMainSize returns the minimum size of the main axis.
	MinMainSize func(*Constraints) *metrics.Dip
	// MaxCrossSize returns the maximum size of the cross axis.
	MaxCrossSize func(*Constraints) *metrics.Dip
	// MinCrossSize returns the minimum size of the cross axis.
	MinCrossSize func(*Constraints) *metrics.Dip
	// MainCoord returns the coordinate of the main axis.
	MainCoord func(*Point) *metrics.Dip
	// CrossCoord returns the coordinate of the cross axis.
	CrossCoord func(*Point) *metrics.Dip

	// MainAxisSize returns the size of the main axis.
	MainAxisSize func() AxisSize
	// MainAxisAlign returns the alignment of the main axis.
	MainAxisAlign func() AxisAlignment
	// CrossAxisAlign returns the alignment of the cross axis.
	CrossAxisAlign func() AxisAlignment

	layoutSize  Size
	itemSizes   []Size
	itemOffsets []Point
}

func (e *colRowElement) Measure(cst Constraints) (size Size, err error) {
	var maxCross metrics.Dip
	var totalMain metrics.Dip
	var availMain = *e.MaxMainSize(&cst)
	for _, item := range e.children {
		var itemSize Size
		var childCst Constraints
		*e.MinMainSize(&childCst) = *e.MinMainSize(&cst) //
		*e.MaxMainSize(&childCst) = availMain
		*e.MinCrossSize(&childCst) = *e.MinCrossSize(&cst)
		*e.MaxCrossSize(&childCst) = *e.MaxCrossSize(&cst)
		itemSize, err = item.Measure(childCst)
		if err != nil {
			return
		}
		checkOverflow(cst, itemSize)
		availMain -= *e.MainSize(&itemSize)

		if itemSize.Height > maxCross {
			maxCross = *e.CrossSize(&itemSize)
		}
		totalMain += *e.MainSize(&itemSize)
		e.itemSizes = append(e.itemSizes, itemSize)
	}

	size.Height = maxCross
	if e.MainAxisSize() == AxisSizeMin {
		*e.MainSize(&size) = totalMain
	} else {
		*e.MainSize(&size) = *e.MaxMainSize(&cst)
	}

	e.layoutSize = size

	for i, itemSize := range e.itemSizes {
		var offset Point
		switch e.CrossAxisAlign() {
		case AlignStart:
			// NOP
		case AlignCenter:
			*e.CrossCoord(&offset) = (maxCross - *e.CrossSize(&itemSize)) / 2
		case AlignEnd:
			*e.CrossCoord(&offset) = maxCross - *e.CrossSize(&itemSize)
		}
		if i == 0 {
			switch e.MainAxisAlign() {
			case AlignStart:
				// NOP
			case AlignCenter:
				*e.MainCoord(&offset) = (*e.MainSize(&size) - totalMain) / 2
			case AlignEnd:
				*e.MainCoord(&offset) = *e.MainSize(&size) - totalMain
			}
		} else {
			*e.MainCoord(&offset) = *e.MainSize(&e.itemSizes[i-1]) + *e.MainCoord(&e.itemOffsets[i-1])
		}
		e.itemOffsets = append(e.itemOffsets, offset)
	}
	return
}

func (e *colRowElement) Arrange(pt Point) (err error) {
	defer func() {
		e.layoutSize = Size{}
		e.itemSizes = e.itemSizes[:0]
		e.itemOffsets = e.itemOffsets[:0]
	}()
	hwnd := e.Widget().HWND()
	if hwnd != 0 {
		if err = positionWindow(hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height); err != nil {
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
