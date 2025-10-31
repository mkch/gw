// Package metrics implements conversion between
// physical pixels(PX) and device-independent pixels(DIP).
package metrics

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/win32"
)

// DPIConv converts a value from old DPI to new DPI.
func DPIConv[T ~int32 | ~uint32](oldValue T, oldDPI, newDPI win32.UINT) (newValue T) {
	return T(win32.MulDiv(win32.INT(oldValue), win32.INT(newDPI), win32.INT(oldDPI)))
}

// FromDefaultDPI convert value from USER_DEFAULT_SCREEN_DPI(96) to a new DPI.
func FromDefaultDPI[T ~int32 | ~uint32](value T, dpi win32.UINT) T {
	return DPIConv(value, win32.USER_DEFAULT_SCREEN_DPI, dpi)
}

// Unit is a measurement uit.
type Unit uint8

const (
	// Physical Pixel.
	PX Unit = iota
	// Device-Independent Pixel.
	// When the DPI is equal to 96, 1 DIP is equal to 1 PX.
	DIP
)

// Dimension represents a dimension with value and unit.
// For example, Dimension{1, PX} means 1 px, and Dimension{2, DIP} means 2 dips.
type Dimension struct {
	Value win32.INT
	Unit  Unit
}

// Px converts the dimension to physical pixels in the given DPI.
func (dim Dimension) Px(dpi win32.UINT) win32.INT {
	if dim.Unit == PX {
		return dim.Value
	}
	return FromDefaultDPI(dim.Value, dpi)
}

// WindowPx converts the dimension to physical pixels in the DPI of the given window.
func (dim Dimension) WidowPx(hwnd win32.HWND) win32.INT {
	return dim.Px(gg.Must(win32.GetDpiForWindow(hwnd)))
}

// Px creates a Dimension in physical pixels.
func Px(n win32.INT) Dimension {
	return Dimension{n, PX}
}

// Dip creates a Dimension in device-independent pixels.
func Dip(n win32.INT) Dimension {
	return Dimension{n, DIP}
}
