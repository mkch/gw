// Package metrics implements conversion between
// physical pixels(Px) and device-independent pixels(Dip).
package metrics

import (
	"math"

	"github.com/mkch/gw/win32"
	"golang.org/x/exp/constraints"
)

// DPIConv converts a value from old DPI to new DPI.
func DPIConv[T constraints.Integer | constraints.Float, DPI constraints.Integer](oldValue T, oldDPI, newDPI DPI) (newValue T) {
	return oldValue * T(newDPI) / T(oldDPI)
}

// FromDefaultDPI convert value from USER_DEFAULT_SCREEN_DPI(96) to a new DPI.
func FromDefaultDPI[T constraints.Integer | constraints.Float, DPI constraints.Integer](value T, dpi DPI) T {
	return DPIConv(value, win32.USER_DEFAULT_SCREEN_DPI, dpi)
}

// ToDefaultDPI convert value from a DPI to USER_DEFAULT_SCREEN_DPI(96).
func ToDefaultDPI[T constraints.Integer | constraints.Float, DPI constraints.Integer](value T, dpi DPI) T {
	return DPIConv(value, dpi, win32.USER_DEFAULT_SCREEN_DPI)
}

// Dimension represents a dimension that can be converted between Px and Dip.
type Dimension interface {
	// Px converts the dimension to physical pixels using the given DPI.
	Px(dpi win32.UINT) Px
	// Dip converts the dimension to device-independent pixels using the given DPI.
	Dip(dpi win32.UINT) Dip
}

// Px is a [Dimension] that represents physical pixels.
type Px win32.INT

func (p Px) Value() win32.INT {
	return win32.INT(p)
}

func (p Px) Px(dpi win32.UINT) Px {
	return p
}

func (p Px) Dip(dpi win32.UINT) Dip {
	return ToDefaultDPI(Dip(p), dpi)
}

// Dip is a [Dimension] that represents device-independent pixels.
type Dip float64

func (d Dip) Value() float64 {
	return float64(d)
}

func (d Dip) Px(dpi win32.UINT) Px {
	return Px(round(float64(FromDefaultDPI(d, dpi))))
}

func (d Dip) Dip(dpi win32.UINT) Dip {
	return d
}

// ToPx converts a Dimension to Px using the given DPI.
// If d is nil, it returns 0.
func ToPx(d Dimension, dpi win32.UINT) Px {
	if d == nil {
		return 0
	}
	return d.Px(dpi)
}

// ToDip converts a Dimension to Dip using the given DPI.
// If d is nil, it returns 0.
func ToDip(d Dimension, dpi win32.UINT) Dip {
	if d == nil {
		return 0
	}
	return d.Dip(dpi)
}

// round converts a float64 to int.
// -1 is returned if f is NaN, Inf or out of int bounds.
func round(f float64) int32 {
	f = math.Round(f)

	if math.IsNaN(f) || math.IsInf(f, 0) || f < float64(math.MinInt32) || f > float64(math.MaxInt32) {
		return -1
	}
	return int32(f)
}
