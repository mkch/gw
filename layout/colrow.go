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
	e.itemSizes = make([]Size, len(e.children))
	e.itemOffsets = make([]Point, len(e.children))

	err = MeasureChildren(cst, e.children,
		func(index int, cst Constraints) (size Size, err error) {
			var childCst Constraints
			*e.MinMainSize(&childCst) = *e.MinMainSize(&cst)
			*e.MaxMainSize(&childCst) = availMain
			*e.MinCrossSize(&childCst) = *e.MinCrossSize(&cst)
			*e.MaxCrossSize(&childCst) = *e.MaxCrossSize(&cst)
			size, err = e.children[index].Measure(childCst)
			if err != nil {
				return
			}
			checkOverflow(cst, size)
			availMain -= *e.MainSize(&size)

			if *e.CrossSize(&size) > maxCross {
				maxCross = *e.CrossSize(&size)
			}
			totalMain += *e.MainSize(&size)
			e.itemSizes[index] = size
			return
		},
		func(cst Constraints, used Size) Constraints {
			// Subtract the used space from the main axis of constraints.
			*e.MinMainSize(&cst) = max(0, *e.MinMainSize(&cst)-*e.MainSize(&used))
			*e.MaxMainSize(&cst) = max(0, *e.MaxMainSize(&cst)-*e.MainSize(&used))
			return cst
		},
		func(size1, size2 Size) (result Size) {
			// Add the main axises of size1 and size2.
			*e.MainSize(&size1) += *e.MainSize(&size2)
			return size1
		},
		func(cst Constraints, used Size) (remain metrics.Dip) {
			// Calculate the remaining space of the main axis.
			remain = max(0, *e.MaxMainSize(&cst)-*e.MainSize(&used))
			return
		},
		func(allocatedSize metrics.Dip) (result Constraints) {
			result = cst
			// Set the allocated size to the main axis of constraints.
			*e.MinMainSize(&result) = allocatedSize
			*e.MaxMainSize(&result) = allocatedSize
			return
		},
		func(index int, size Size) {
			// Callback for each measured Expanded child with non-zero size.
			availMain -= *e.MainSize(&size)
			if *e.CrossSize(&size) > maxCross {
				maxCross = *e.CrossSize(&size)
			}
			totalMain += *e.MainSize(&size)
			e.itemSizes[index] = size
		},
	)

	*e.CrossSize(&size) = maxCross
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
		e.itemOffsets[i] = offset
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
