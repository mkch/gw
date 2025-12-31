package paint

import "github.com/mkch/gw/win32"

type Buffer struct {
	DC
	realDC        win32.HDC
	bitmap        win32.HBITMAP
	oldBitmap     win32.HBITMAP
	width, height int
}

func NewBuffer(dc win32.HDC, width, height int) (*Buffer, error) {
	memDC, err := win32.CreateCompatibleDC(dc)
	if err != nil {
		return nil, err
	}
	bitmap, err := win32.CreateCompatibleBitmap(dc, win32.INT(width), win32.INT(height))
	if err != nil {
		win32.DeleteDC(memDC)
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
		realDC:    dc,
		bitmap:    bitmap,
		oldBitmap: oldBitmap,
		width:     width, height: height}, nil
}

// Resize resizes the buffer to hold at least the given width and height.
// The content of buffer will be lost after resizing.
func (buf *Buffer) Resize(width, height int) error {
	if width <= buf.width && height <= buf.height && // The new sizes are smaller than current buffer
		width*4 > buf.width && height*4 > buf.height { // but not smaller than quarter size
		return nil // No need to resize
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
	memDC, err := win32.CreateCompatibleDC(buf.realDC)
	if err != nil {
		return err
	}
	bitmap, err := win32.CreateCompatibleBitmap(buf.realDC, win32.INT(width), win32.INT(height))
	if err != nil {
		win32.DeleteDC(memDC)
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
	return nil
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
