package layout

import (
	"fmt"
	"math"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

// round converts a float64 to int.
// -1 is returned if f is NaN, Inf or out of int bounds.
func round(f float64) int32 {
	f = math.Round(f)

	if math.IsNaN(f) || math.IsInf(f, 0) || f < float64(math.MinInt32) || f > float64(math.MaxInt32) {
		return -1
	}
	return int32(f)
}

type DIP win32.INT

// Infinity represents an infinite size(unbounded) constraint.
const Infinity = math.MaxInt32

type Size struct {
	Width  DIP
	Height DIP
}

func (s *Size) String() string {
	return fmt.Sprintf("{Width: %v, Height: %v}", s.Width, s.Height)
}

type Point struct {
	X, Y DIP
}

func (pt Point) String() string {
	return fmt.Sprintf("{X: %v, Y: %v}", pt.X, pt.Y)
}

type Rect struct {
	Left, Top, Right, Bottom DIP
}

func NewRect(pt Point, size Size) Rect {
	return Rect{
		Left:   pt.X,
		Top:    pt.Y,
		Right:  pt.X + size.Width,
		Bottom: pt.Y + size.Height,
	}
}

func (r *Rect) Size() Size {
	return Size{
		Width:  r.Width(),
		Height: r.Height(),
	}
}

func (r *Rect) Width() DIP {
	return r.Right - r.Left
}

func (r *Rect) Height() DIP {
	return r.Bottom - r.Top
}

func (r *Rect) TopLeft() Point {
	return Point{X: r.Left, Y: r.Top}
}

func (r *Rect) BottomRight() Point {
	return Point{X: r.Right, Y: r.Bottom}
}

// Constraints represents layout constraints.
type Constraints struct {
	MinWidth  DIP
	MinHeight DIP
	MaxWidth  DIP
	MaxHeight DIP
}

func (c *Constraints) String() string {
	return fmt.Sprintf("{MinWidth: %v, MinHeight: %v, MaxWidth: %v, MaxHeight: %v}",
		c.MinWidth, c.MinHeight, c.MaxWidth, c.MaxHeight)
}

// TightWidth returns true if the constraint has a finite and equal min and max width.
func (c *Constraints) TightWidth() bool {
	return c.MinWidth == c.MaxWidth && c.MinWidth != Infinity
}

// TightHeight returns true if the constraint has a finite and equal min and max height.
func (c *Constraints) TightHeight() bool {
	return c.MinHeight == c.MaxHeight && c.MinHeight != Infinity
}

// UnboundWidth returns true if no constraint is imposed on width.
func (c *Constraints) UnboundWidth() bool {
	return c.MaxWidth == Infinity
}

// UnboundHeight returns true if no constraint is imposed on height.
func (c *Constraints) UnboundHeight() bool {
	return c.MaxHeight == Infinity
}

// MinSize returns the minimum size allowed by the constraints.
func (c *Constraints) MinSize() Size {
	return Size{Width: c.MinWidth, Height: c.MinHeight}
}

// MaxSize returns the maximum size allowed by the constraints.
func (c *Constraints) MaxSize() Size {
	return Size{Width: c.MaxWidth, Height: c.MaxHeight}
}

// Clamp clamps the given size between the constraints.
func (c *Constraints) Clamp(size Size) Size {
	return Size{
		Width:  clamp(size.Width, c.MinWidth, c.MaxWidth),
		Height: clamp(size.Height, c.MinHeight, c.MaxHeight),
	}
}

// ClampHeight clamps height between the constraints.
func (c *Constraints) ClampWidth(width DIP) DIP {
	return clamp(width, c.MinWidth, c.MaxWidth)
}

// ClampHeight clamps height between the constraints.
func (c *Constraints) ClampHeight(height DIP) DIP {
	return clamp(height, c.MinHeight, c.MaxHeight)
}

// TightMax returns a new Constraints with tight max constraints.
func (c *Constraints) TightMax() Constraints {
	return Constraints{
		MinWidth:  c.MaxWidth,
		MinHeight: c.MaxHeight,
		MaxWidth:  c.MaxWidth,
		MaxHeight: c.MaxHeight,
	}
}

// TightMin returns a new Constraints with tight min constraints.
func (c *Constraints) TightMin() Constraints {
	return Constraints{
		MinWidth:  c.MinWidth,
		MinHeight: c.MinHeight,
		MaxWidth:  c.MinWidth,
		MaxHeight: c.MinHeight,
	}
}

// clamp clamps value between min and max.
func clamp(value, minBound, maxBound DIP) DIP {
	return min(max(value, minBound), maxBound)
}

func PxToDP[T ~int | ~int32 | ~uint32](px T, dpi win32.UINT) DIP {
	return DIP(metrics.DPIConv(px, dpi, win32.USER_DEFAULT_SCREEN_DPI))
}

// ClientSize returns the size of the client area of the given window in DPs.
func ClientSize(hwnd win32.HWND) (size Size, err error) {
	var clientRect win32.RECT
	if err = win32.GetClientRect(hwnd, &clientRect); err != nil {
		return
	}
	if err = win32util.ScreenToClient(hwnd, &clientRect); err != nil {
		return
	}
	dpi, err := win32.GetDpiForWindow(hwnd)
	if err != nil {
		return
	}
	size = Size{
		Width:  PxToDP(clientRect.Right-clientRect.Left, dpi),
		Height: PxToDP(clientRect.Bottom-clientRect.Top, dpi),
	}
	return
}

func setWindowPos(hwnd win32.HWND, x, y, width, height DIP) (err error) {
	var dpi win32.UINT
	dpi, err = win32.GetDpiForWindow(hwnd)
	if err != nil {
		return err
	}
	err = win32.SetWindowPos(hwnd, 0,
		metrics.Dip(win32.INT(x)).Px(dpi),
		metrics.Dip(win32.INT(y)).Px(dpi),
		metrics.Dip(win32.INT(width)).Px(dpi),
		metrics.Dip(win32.INT(height)).Px(dpi), win32.SWP_NOZORDER)
	if err != nil {
		return
	}
	return
}

func checkOverflow(cst Constraints, size Size) {
	if size.Width < cst.MinWidth || size.Width > cst.MaxWidth ||
		size.Height < cst.MinHeight || size.Height > cst.MaxHeight {
		panic(fmt.Sprintf("Size %v is out of constraints %v", size, cst))
	}
}
