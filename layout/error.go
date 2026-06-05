package layout

import "fmt"

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
