package layout

import (
	"iter"
	"slices"

	"github.com/mkch/gw/win32"
)

// Widget is a lightweight description of a layout element.
type Widget interface {
	HWND() win32.HWND
	CreateElement() (Element, error)
	ChildWidgets() iter.Seq2[int, Widget]
}

// Element is a layout element that performs the actual layout.
// It is created from a Widget.
type Element interface {
	Widget() Widget
	Parent() Element
	Children() iter.Seq2[int, Element]
	Measure(Constraints) (Size, error)
	Arrange(Point) error

	setParentField(parent Element)
	appendChildToSlice(child Element)
}

// BaseElement provides common fields and methods for all elements.
// It is a build block for concrete [Element] implementations.
type BaseElement struct {
	widget   Widget
	parent   Element
	children []Element
}

func (e *BaseElement) Widget() Widget {
	return e.widget
}

func (e *BaseElement) Parent() Element {
	return e.parent
}

func (e *BaseElement) Children() iter.Seq2[int, Element] {
	return slices.All(e.children)
}

func (e *BaseElement) setParentField(parent Element) {
	e.parent = parent
}

func (e *BaseElement) appendChildToSlice(child Element) {
	e.children = append(e.children, child)
}

func element_AddChild(parent Element, child Element) {
	child.setParentField(parent)
	parent.appendChildToSlice(child)
}
