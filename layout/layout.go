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
	"github.com/mkch/gw/win32"
)

// Build constructs an Element tree from the specified Widget tree.
//
// All windowed widgets in the tree must form a valid Win32 parent-child
// hierarchy. Widgets that do not own a window (HWND == 0) are treated as
// logical containers and do not affect parent-child validation.
//
// If an invalid parent-child relationship is detected, Build returns a
// [*WrongParentError] identifying the location of the offending widget.
func Build(root Widget) (tree Element, err error) {
	var parent win32.HWND
	return build(root, &parent, nil)
}

// build recursively constructs an Element subtree rooted at root.
//
// The parent parameter points to the expected Win32 parent window for the
// current branch of the tree. Logical containers (HWND == 0) inherit and
// propagate this expectation unchanged. If no expected parent has been
// established yet, the first windowed widget encountered determines it for
// all subsequently visited descendants until another windowed widget is
// reached.
//
// If a widget has a window, its actual parent window is validated against
// the expected parent. That widget's window then becomes the expected parent
// for its own children.
//
// The indices slice records the widget's position within the tree and is used
// to construct a *WrongParentError when validation fails.
func build(root Widget, parent *win32.HWND, indices []int) (tree Element, err error) {
	tree, err = root.CreateElement()
	if err != nil {
		return
	}
	hwnd := tree.Widget().HWND()
	if *parent != 0 {
		if hwnd != 0 {
			var realParent win32.HWND
			if realParent, err = win32.GetParent(hwnd); err != nil {
				return
			}
			if realParent != *parent {
				return nil, &WrongParentError{Indices: indices, Widget: root}
			}
		}
	} else if hwnd != 0 {
		var realParent win32.HWND
		if realParent, err = win32.GetParent(hwnd); err != nil {
			return
		}
		*parent = realParent
	}

	var parentForChildren *win32.HWND
	if hwnd != 0 {
		parentForChildren = &hwnd
	} else {
		parentForChildren = parent
	}
	for i, child := range root.ChildWidgets() {
		var childTree Element
		childTree, err = build(child, parentForChildren, append(indices, i))
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
