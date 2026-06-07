package layout

import (
	"iter"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/iter2"
)

// Expanded is a widget that expands to fill the available space in the parent container.
// Expanded passes a tight constraint to its child widget to occupy all the allocated space.
// If more than one Expanded widget are present in a parent, the available space is divided
// among them according to their Flex factor.
type Expanded struct {
	Hwnd win32.HWND
	// The Flex factor to use for this Expanded widget.
	// Flex must be a non-negative.
	Flex float64
	// Child is the child to layout. Must be non-nil.
	Child Widget
}

func (e *Expanded) HWND() win32.HWND {
	return e.Hwnd
}

func (e *Expanded) ChildWidgets() iter.Seq2[int, Widget] {
	return iter2.Just2(struct {
		K int
		V Widget
	}{0, e.Child})
}

func (e *Expanded) CreateElement() (Element, error) {
	if e.Child == nil {
		return nil, &NoChildError{Layout: "Expanded"}
	}
	if e.Flex < 0 {
		return nil, &InvalidFlexError{Flex: e.Flex}
	}
	return &expandedElement{
		BaseElement: BaseElement{
			widget: e,
		},
	}, nil
}

type expandedElement struct {
	BaseElement
	layoutSize Size
}

func (e *expandedElement) Measure(cst Constraints) (size Size, err error) {
	size = cst.MaxSize()
	e.layoutSize = size
	child := e.children[0]
	childSize, err := child.Measure(cst.TightMax())
	if err != nil {
		return Size{}, err
	}
	checkOverflow(cst, childSize)
	return
}

func (e *expandedElement) Arrange(pt Point) (err error) {
	defer func() {
		e.layoutSize = Size{}
	}()
	w := e.Widget().(*Expanded)
	child := e.children[0]
	if w.Hwnd != 0 {
		if err = positionWindow(w.Hwnd, pt.X, pt.Y, e.layoutSize.Width, e.layoutSize.Height); err != nil {
			return err
		}
		return child.Arrange(Point{0, 0})
	}
	return child.Arrange(pt)
}

// MeasureChildren measures the children of a layout that supports [Expanded] widgets.
// The function first measures all non-Expanded children to determine the remaining space for Expanded children.
// Then it measures the [Expanded] children based on their Flex factors and the remaining space.
// The parameters are as follows:
//
// - cst: The constraints for the layout.
//
// - children: The child elements to measure.
//
// - measureChild: A function to measure a child element with given constraints.
//
// - subtractConstraints: A function to subtract the size of a measured child from the remaining constraints.
//
// - addSize: A function to add the size of a measured child to the total size.
//
// - remainSpace: A function to calculate the remaining space after measuring non-Expanded children.
//
// - buildConstraints: A function to build constraints for an Expanded child based on the allocated size.
//
// - onMeasureExpanded: A callback function that is called with the size of each measured Expanded child whose size is not zero.
func MeasureChildren(cst Constraints,
	children []Element,
	measureChild func(index int, cst Constraints) (Size, error),
	subtractConstraints func(Constraints, Size) Constraints,
	addSize func(Size, Size) Size,
	remainSpace func(cst Constraints, used Size) metrics.Dip,
	buildConstraints func(allocatedSize metrics.Dip) Constraints,
	onMeasureExpanded func(index int, size Size),
) (err error) {
	var expandedElements []*expandedElement
	var expandedIndexes []int
	var totalExpandedFlex float64
	var childCst = cst
	var nonExpSize Size
	for i, child := range children {
		// Collect all Expanded elements to measure them later.
		if expanded, ok := child.Widget().(*Expanded); ok {
			expandedElements = append(expandedElements, child.(*expandedElement))
			expandedIndexes = append(expandedIndexes, i)
			totalExpandedFlex += expanded.Flex
			continue
		}
		// Measure non-Expanded children immediately.
		var childSize Size
		if childSize, err = measureChild(i, childCst); err != nil {
			return
		}
		childCst = subtractConstraints(childCst, childSize)
		nonExpSize = addSize(nonExpSize, childSize)
	}

	// No Expanded children, we're done.
	if len(expandedElements) == 0 {
		return
	}

	// Only one Expanded child, give it all the remaining space.
	if len(expandedElements) == 1 {
		var availableSpace = remainSpace(cst, nonExpSize)
		var childCst = buildConstraints(availableSpace)
		var childSize Size
		if childSize, err = expandedElements[0].Measure(childCst); err != nil {
			return
		}
		checkOverflow(childCst, childSize)
		onMeasureExpanded(expandedIndexes[0], childSize)
		return
	}

	// If total Flex is zero, all Expanded children get zero space.
	if totalExpandedFlex == 0 {
		var zeroCst Constraints
		for i, e := range expandedElements {
			if _, err = e.Measure(zeroCst); err != nil {
				return
			}
			checkOverflow(zeroCst, zeroCst.MaxSize())
			onMeasureExpanded(expandedIndexes[i], zeroCst.MaxSize())
		}
		return
	}

	totalFlexInv := 1 / totalExpandedFlex

	var remain = remainSpace(cst, nonExpSize)
	for i, e := range expandedElements {
		var size metrics.Dip
		if i == len(expandedElements)-1 {
			// Give all remaining space to the last Expanded child to avoid rounding errors.
			size = remain
		} else if remain > 0 {
			if flex := e.Widget().(*Expanded).Flex; flex > 0 {
				size = metrics.Dip(float64(remain) * flex * totalFlexInv)
			}
			remain -= size
		}
		var childCst = buildConstraints(size)
		var childSize Size
		if childSize, err = e.Measure(childCst); err != nil {
			return
		}
		checkOverflow(childCst, childSize)
		onMeasureExpanded(expandedIndexes[i], childSize)
	}

	return
}
