package layout

import (
	"errors"
	"fmt"
)

// ErrWrongRoot is returned by [Perform] and [PerformWindow] if the root layout has an associated window.
var ErrWrongRoot = errors.New("root layout must not have an associated window")

// WrongParentError is returned by [Build] if a Widget's actual Win32 parent window
// does not match the parent window expected by the widget tree structure.
//
// Indices identifies the path from the root widget to the offending widget.
// Each element is the child index within its parent's ChildWidgets slice.
type WrongParentError struct {
	Indices []int
	Widget  Widget
}

func (e *WrongParentError) Error() string {
	return fmt.Sprintf("widget of type %T at %v has an unexpected native parent window", e.Widget, e.Indices)
}

// NoChildError is returned by a layout that requires a child but the child is nil.
type NoChildError struct {
	Layout string
}

func (e *NoChildError) Error() string {
	return fmt.Sprintf("%s layout must have a child", e.Layout)
}

// InvalidSizeError is returned by a layout if the size factor is invalid.
type InvalidFactorError struct {
	Name  string
	Value float64
}

func (e *InvalidFactorError) Error() string {
	return fmt.Sprintf("invalid %s %v: must be 0 or >= 1", e.Name, e.Value)
}

// NoHwndError is returned by a layout which requires a non-zero Hwnd but the Hwnd is zero.
type NoHwndError struct {
	Layout string
}

func (e *NoHwndError) Error() string {
	return fmt.Sprintf("%s layout must have a non-zero Hwnd", e.Layout)
}

// InvalidFlexError is returned when a Flex factor is invalid.
type InvalidFlexError struct {
	Flex float64
}

func (e *InvalidFlexError) Error() string {
	return fmt.Sprintf("invalid Flex %v: must be non-negative", e.Flex)
}
