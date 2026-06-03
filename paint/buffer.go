package paint

import (
	"math"

	"github.com/mkch/gw/win32"
)

// Buffer is a double-buffered in-memory drawing context.
type Buffer struct {
	DC
	realDC              DCProvider
	bitmap              win32.HBITMAP
	oldBitmap           win32.HBITMAP
	width, height       int // The dimensions of the requested buffer size
	capWidth, capHeight int // The dimensions of the actual bitmap, capacity dimensions of the buffer.
}

// DCProvider is an interface that provides a DC for use in a callback function.
type DCProvider interface {
	// WithDC executes the given function with a DC.
	WithDC(func(dc win32.HDC) error) error
}

// DCProviderFunc is a function type that implements [DCProvider].
type DCProviderFunc func(func(dc win32.HDC) error) error

func (f DCProviderFunc) WithDC(fn func(dc win32.HDC) error) error {
	return f(fn)
}

// HDC is a [DCProvider] that provides a single HDC.
type HDC win32.HDC

func (h HDC) WithDC(f func(dc win32.HDC) error) error {
	return f(win32.HDC(h))
}

// ClientDC is a [DCProvider] that provides a DC from [win32.GetDC] and releases it after use.
type ClientDC win32.HWND

func (g ClientDC) WithDC(f func(dc win32.HDC) error) error {
	dc, err := win32.GetDC(win32.HWND(g))
	if err != nil {
		return err
	}
	defer win32.ReleaseDC(win32.HWND(g), dc)
	return f(dc)
}

// NewBuffer creates a [Buffer] with the given width and height.
// The realDC provider is used to create compatible DC and bitmap for the buffer.
func NewBuffer(realDC DCProvider, width, height int) (*Buffer, error) {
	var memDC win32.HDC
	var bitmap win32.HBITMAP
	if err := realDC.WithDC(func(dc win32.HDC) error {
		var err error
		if memDC, err = win32.CreateCompatibleDC(dc); err != nil {
			return err
		}
		if bitmap, err = win32.CreateCompatibleBitmap(dc, win32.INT(width), win32.INT(height)); err != nil {
			win32.DeleteDC(memDC)
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	oldBitmap, err := win32.SelectObject(memDC, bitmap)
	if err != nil {
		win32.DeleteDC(memDC)
		win32.DeleteObject(bitmap)
		return nil, err
	}
	return &Buffer{
		DC:        DC{memDC},
		realDC:    realDC,
		bitmap:    bitmap,
		oldBitmap: oldBitmap,
		width:     width,
		height:    height,
		capWidth:  width,
		capHeight: height,
	}, nil
}

// growByTenPercent returns a size that is 10% larger than the given size, rounded up.
func growByTenPercent(size int) int {
	if size <= 0 {
		return size
	}
	if size > (math.MaxInt-9)/11 {
		return math.MaxInt
	}
	return (size*11 + 9) / 10
}

// Resize resizes the buffer to hold at least the given width and height.
// The content of buffer will be lost after resizing.
func (buf *Buffer) Resize(width, height int) error {
	if width <= buf.capWidth && height <= buf.capHeight && // Requested size is within bitmap capacity
		width > buf.capWidth/4 && height > buf.capHeight/4 { // and not shrinking below quarter of current capacity
		buf.width = width
		buf.height = height
		return nil
	}
	capWidth := width
	capHeight := height
	if width > buf.capWidth {
		capWidth = growByTenPercent(width)
	}
	if height > buf.capHeight {
		capHeight = growByTenPercent(height)
	}

	// Select back the old bitmap and delete current bitmap and DC
	if _, err := win32.SelectObject(buf.hdc, buf.oldBitmap); err != nil {
		return err
	}
	// Delete current DC and bitmap
	if err := win32.DeleteDC(buf.hdc); err != nil {
		return err
	}
	if err := win32.DeleteObject(buf.bitmap); err != nil {
		return err
	}
	// Create new DC and bitmap
	var memDC win32.HDC
	var bitmap win32.HBITMAP
	if err := buf.realDC.WithDC(func(dc win32.HDC) error {
		var err error
		if memDC, err = win32.CreateCompatibleDC(dc); err != nil {
			return err
		}
		if bitmap, err = win32.CreateCompatibleBitmap(dc, win32.INT(capWidth), win32.INT(capHeight)); err != nil {
			win32.DeleteDC(memDC)
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	oldBitmap, err := win32.SelectObject(memDC, bitmap)
	if err != nil {
		win32.DeleteDC(memDC)
		win32.DeleteObject(bitmap)
		return err
	}
	// Update buffer fields
	buf.hdc = memDC
	buf.bitmap = bitmap
	buf.oldBitmap = oldBitmap
	buf.width = width
	buf.height = height
	buf.capWidth = capWidth
	buf.capHeight = capHeight
	return nil
}

// Width returns the width of the in-memory buffer.
func (buf *Buffer) Width() int {
	return buf.width
}

// Height returns the height of the in-memory buffer.
func (buf *Buffer) Height() int {
	return buf.height
}

func (buf *Buffer) Destroy() error {
	if _, err := win32.SelectObject(buf.hdc, buf.oldBitmap); err != nil {
		return err
	}
	if err := win32.DeleteDC(buf.DC.hdc); err != nil {
		return err
	}
	if err := win32.DeleteObject(buf.bitmap); err != nil {
		return err
	}
	return nil
}
