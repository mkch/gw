package paint

import "github.com/mkch/gw/win32"

type Buffer struct {
	DC
	bitmap    win32.HBITMAP
	oldBitmap win32.HBITMAP
}

func NewBuffer(dc win32.HDC, width, height win32.INT) (*Buffer, error) {
	dc, err := win32.CreateCompatibleDC(dc)
	if err != nil {
		return nil, err
	}
	bitmap, err := win32.CreateCompatibleBitmap(dc, width, height)
	if err != nil {
		win32.DeleteDC(dc)
		return nil, err
	}
	oldBitmap, err := win32.SelectObject(dc, bitmap)
	if err != nil {
		win32.DeleteDC(dc)
		win32.DeleteObject(bitmap)
		return nil, err
	}
	return &Buffer{DC: DC{dc}, bitmap: bitmap, oldBitmap: oldBitmap}, nil
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
