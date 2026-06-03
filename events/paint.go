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

// WindowPaintDC is the drawing context for WM_PAINT message.
// It implements [PaintDC].
type WindowPaintDC paint.PaintDC

func (p *WindowPaintDC) HWND() win32.HWND {
	return (*paint.PaintDC)(p).HWND()
}

func (p *WindowPaintDC) HDC() win32.HDC {
	return (*paint.PaintDC)(p).HDC()
}

func (p *WindowPaintDC) EraseBackground() bool {
	return (*paint.PaintDC)(p).EraseBackground()
}

func (p *WindowPaintDC) Rect() *win32.RECT {
	return (*paint.PaintDC)(p).Rect()
}

// endPaint is a proxy for [paint.PaintDC.EndPaint] but is unexported.
func (p *WindowPaintDC) endPaint() error {
	return (*paint.PaintDC)(p).EndPaint()
}

// NewWindowPaintDC creates a new [WindowPaintDC] for the given window handle.
func NewWindowPaintDC(w win32.HWND) (*WindowPaintDC, error) {
	dc, err := paint.NewPaintDC(w)
	if err != nil {
		return nil, err
	}
	return (*WindowPaintDC)(dc), nil
}

// BufferedPaintDC is a double-buffered drawing context.
// It implements [PaintDC].
type BufferedPaintDC struct {
	*WindowPaintDC
	buffer         *paint.Buffer
	externalBuffer bool // buffer is created outside of this struct, so it should not be destroyed by this struct
}

// HDC returns the HDC of the buffer.
func (b *BufferedPaintDC) HDC() win32.HDC {
	return b.buffer.HDC()
}

func (b *BufferedPaintDC) endPaint() (err error) {
	var clientRect win32.RECT
	err = win32.GetClientRect(b.WindowPaintDC.HWND(), &clientRect)
	if err != nil {
		return err
	}
	w := int(clientRect.Width())
	h := int(clientRect.Height())
	defer gg.CollectError(b.WindowPaintDC.endPaint, &err)
	if !b.externalBuffer {
		defer gg.CollectError(b.buffer.Destroy, &err)
	}
	err = win32.BitBlt(b.WindowPaintDC.HDC(), 0, 0, w, h, b.buffer.HDC(), 0, 0, win32.SRCCOPY)
	return
}

// NewBufferedPaintDC creates a new [BufferedPaintDC] for the given window handle.
// The size of the buffer is the same as the client area of w.
func NewBufferedPaintDC(w win32.HWND) (*BufferedPaintDC, error) {
	var clientRect win32.RECT
	err := win32.GetClientRect(w, &clientRect)
	if err != nil {
		return nil, err
	}
	dc, err := NewWindowPaintDC(w)
	if err != nil {
		return nil, err
	}

	dcProvider := paint.DCProviderFunc(func(f func(dc win32.HDC) error) error {
		return f(dc.HDC())
	})
	buffer, err := paint.NewBuffer(dcProvider, int(clientRect.Width()), int(clientRect.Height()))
	if err != nil {
		dc.endPaint()
		return nil, err
	}

	return &BufferedPaintDC{
		WindowPaintDC: dc,
		buffer:        buffer,
	}, nil
}

// NewBufferedPaintDCWithBuffer creates a new [BufferedPaintDC] for the given window handle and buffer.
// The buffer must be valid for the lifetime of the returned [bufferedPaintDC].
func NewBufferedPaintDCWithBuffer(w win32.HWND, buffer *paint.Buffer) (*BufferedPaintDC, error) {
	dc, err := NewWindowPaintDC(w)
	if err != nil {
		return nil, err
	}
	return &BufferedPaintDC{
		WindowPaintDC:  dc,
		buffer:         buffer,
		externalBuffer: true,
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
	extBuffer     *paint.Buffer // external buffer for double buffering
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
		if e.extBuffer != nil {
			e.dc, err = NewBufferedPaintDCWithBuffer(e.hwnd, e.extBuffer)
		} else {
			e.dc, err = NewWindowPaintDC(e.hwnd)
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
	e.defProcCalled = true
	return proc(e.hwnd, win32.WM_PAINT, e.wParam, e.lParam)
}

// NewPaintEvent creates a new PaintEvent.
// If buffer is not nil, it will be used for double buffering, otherwise the painting will be done directly on the window.
// If the buffer is not nil, it must be valid for the lifetime of the returned PaintEvent.
// The returned release function should be called to end the painting.
func NewPaintEvent(hwnd win32.HWND, wParam win32.WPARAM, lParam win32.LPARAM, buffer *paint.Buffer) (evt *PaintEvent, release func() error) {
	evt = &PaintEvent{
		hwnd:      hwnd,
		wParam:    wParam,
		lParam:    lParam,
		extBuffer: buffer,
	}
	release = evt.end
	return
}
