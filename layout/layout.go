// Package layout provides a DPI-aware layout engine for arranging HWND-backed
// controls in a layout tree.
//
// Layouts are defined as a tree of Widget, which are lightweight descriptions of the layout.
// They are converted to an Element tree, which performs the actual layout.
//
// Elements follow a two-phase contract: Measure computes desired size
// under Constraints, and Arrange places the layout and its children at a Point.
// Perform runs this flow for a tree root.
package layout

import (
	"errors"

	"github.com/mkch/gw/win32"
)

// Build builds the element tree for the given layout tree.
func Build(root Widget) (tree Element, err error) {
	tree, err = root.CreateElement()
	if err != nil {
		return
	}

	for _, child := range root.ChildWidgets() {
		var childTree Element
		childTree, err = Build(child)
		if err != nil {
			return
		}
		element_AddChild(tree, childTree)
	}

	return
}

// PerformWindow calls [Perform] with the given layout and the client size of the given window.
func PerformWindow(root Element, hwnd win32.HWND) (err error) {
	var size Size
	size, err = ClientSize(hwnd)
	if err != nil {
		return
	}
	return Perform(root, size)
}

// ErrWrongRoot is returned by [Perform] and [PerformWindow] if the root layout has an associated window.
var ErrWrongRoot = errors.New("root layout must not have an associated window")

// Perform performs the layout by measuring the layout with the given size constraint and arranging the layout at (0, 0).
// The root must not have an associated window.
func Perform(root Element, size Size) (err error) {
	if root.Widget().HWND() != 0 {
		return ErrWrongRoot
	}
	root.Measure(Constraints{
		MaxWidth:  size.Width,
		MaxHeight: size.Height,
	})
	return root.Arrange(Point{X: 0, Y: 0})
}
