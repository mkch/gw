package events

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/win32"
)

// PaintDC is the interface for drawing context of WM_PAINT message.
type PaintDC interface {
	HDC() win32.HDC
	EraseBackground() bool
	Rect() *win32.RECT
	endPaint() error
}

// paintDC is the drawing context for WM_PAINT message.
// It implements [PaintDC].
type paintDC paint.PaintDC

func (p *paintDC) HDC() win32.HDC {
	return (*paint.PaintDC)(p).HDC()
}

func (p *paintDC) EraseBackground() bool {
	return (*paint.PaintDC)(p).EraseBackground()
}

func (p *paintDC) Rect() *win32.RECT {
	return (*paint.PaintDC)(p).Rect()
}

// endPaint is a proxy for [paint.PaintDC.EndPaint] but is unexported.
func (p *paintDC) endPaint() error {
	return (*paint.PaintDC)(p).EndPaint()
}

// newPaintDC creates a new [paintDC] for the given window handle.
func newPaintDC(w win32.HWND) (*paintDC, error) {
	dc, err := paint.NewPaintDC(w)
	if err != nil {
		return nil, err
	}
	return (*paintDC)(dc), nil
}

// bufferedPaintDC is a double-buffered drawing context.
// It implements [PaintDC].
type bufferedPaintDC struct {
	*paintDC
	buffer *paint.Buffer
}

// HDC returns the HDC of the buffer.
func (b *bufferedPaintDC) HDC() win32.HDC {
	return b.buffer.HDC()
}

func (b *bufferedPaintDC) endPaint() (err error) {
	defer gg.CollectError(b.paintDC.endPaint, &err)
	defer gg.CollectError(b.buffer.Destroy, &err)
	err = win32.BitBlt(b.paintDC.HDC(), 0, 0, b.buffer.Width(), b.buffer.Height(), b.buffer.HDC(), 0, 0, win32.SRCCOPY)
	return
}

// newBufferedPaintDC creates a new [bufferedPaintDC] for the given window handle.
// The size of the buffer is the same as the client area of w.
func newBufferedPaintDC(w win32.HWND) (*bufferedPaintDC, error) {
	var clientRect win32.RECT
	err := win32.GetClientRect(w, &clientRect)
	if err != nil {
		return nil, err
	}
	dc, err := newPaintDC(w)
	if err != nil {
		return nil, err
	}

	buffer, err := paint.NewBuffer(dc.HDC(), int(clientRect.Width()), int(clientRect.Height()))
	if err != nil {
		dc.endPaint()
		return nil, err
	}

	return &bufferedPaintDC{
		paintDC: dc,
		buffer:  buffer,
	}, nil
}

// PaintEvent is the event for WM_PAINT message.
type PaintEvent struct {
	hwnd win32.HWND
	// Save wParam and lParam of WM_PAINT message, because
	// "For some common controls, the default WM_PAINT message processing checks the wParam parameter ..."
	// https://learn.microsoft.com/en-us/windows/win32/gdi/wm-paint
	wParam win32.WPARAM
	lParam win32.LPARAM

	defProcCalled bool
	dc            PaintDC
	buffered      bool
}

// Begin returns the drawing context for WM_PAINT message.
// Multiple calls to this method will return the same drawing context for a single valid PaintEvent.
// If the default window procedure is called before this method, it panics.
func (e *PaintEvent) Begin() (PaintDC, error) {
	if e.defProcCalled {
		panic("Begin called after default window procedure called")
	}
	if e.dc == nil {
		var err error
		if e.buffered {
			e.dc, err = newBufferedPaintDC(e.hwnd)
		} else {
			e.dc, err = newPaintDC(e.hwnd)
		}
		if err != nil {
			return nil, err
		}
	}
	return e.dc, nil
}

// end ends the painting.
// If the default window procedure is called before this method or [PaintEvent.Begin] is not called,
// it does nothing but returns nil.
func (e *PaintEvent) end() error {
	if e.dc != nil && !e.defProcCalled {
		err := e.dc.endPaint()
		e.dc = nil
		return err
	}
	return nil
}

// callDefProc calls the proc and returns its return value if this method and [PaintEvent.Begin] are both not called,
// otherwise it does nothing but return 0. This method is used to call the default window procedure for WM_PAINT message.
func (e *PaintEvent) CallDefProc(proc func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT) win32.LRESULT {
	if e.defProcCalled || e.dc != nil {
		return 0
	}
	return proc(e.hwnd, win32.WM_PAINT, e.wParam, e.lParam)
}

// NewPaintEvent creates a new PaintEvent.
// The returned release function should be called to end the painting.
func NewPaintEvent(hwnd win32.HWND, wParam win32.WPARAM, lParam win32.LPARAM, buffered bool) (evt *PaintEvent, release func() error) {
	evt = &PaintEvent{
		hwnd:     hwnd,
		wParam:   wParam,
		lParam:   lParam,
		buffered: buffered,
	}
	release = evt.end
	return
}
