package gw

import (
	"sync"
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const doubleBufferedTestClassName = "github.com/mkch/gw#DoubleBufferedTest"

var ensureDoubleBufferedTestClassRegistered = sync.OnceValue(func() error {
	_, err := win32util.RegisterClass(&win32util.WndClass{
		ClassName: doubleBufferedTestClassName,
		WndProc:   win32.DefWindowProcW,
	})
	return err
})

type doubleBufferedTestWindow struct {
	BaseWindowImpl
	onPaint func(evt *events.PaintEvent)
}

func (w *doubleBufferedTestWindow) CreateHandle() (win32.HWND, error) {
	if err := ensureDoubleBufferedTestClassRegistered(); err != nil {
		return 0, err
	}
	return win32util.CreateWindow(&win32util.Wnd{
		ClassName: doubleBufferedTestClassName,
		Style:     win32.WS_OVERLAPPEDWINDOW,
		Width:     300,
		Height:    200,
	})
}

func (w *doubleBufferedTestWindow) OnPaint(evt *events.PaintEvent) {
	if w.onPaint != nil {
		w.onPaint(evt)
		return
	}
	w.BaseWindowImpl.OnPaint(evt)
}

func newDoubleBufferedTestWindow(t *testing.T, app *app.App) *doubleBufferedTestWindow {
	t.Helper()
	w := chkerr.Must(Init(&doubleBufferedTestWindow{}))
	w.SetOnDestroyListener(func() { app.Quit(0) })
	return w
}

func TestSetDoubleBufferedToggleAndIdempotent(t *testing.T) {
	Run(func(app *app.App) {
		w := newDoubleBufferedTestWindow(t, app)

		if err := w.SetDoubleBuffered(true); err != nil {
			t.Fatalf("SetDoubleBuffered(true) failed: %v", err)
		}
		if w.paintBuffer == nil {
			t.Fatal("paintBuffer should be created after SetDoubleBuffered(true)")
		}
		firstBuffer := w.paintBuffer

		if err := w.SetDoubleBuffered(true); err != nil {
			t.Fatalf("second SetDoubleBuffered(true) failed: %v", err)
		}
		if w.paintBuffer != firstBuffer {
			t.Fatal("SetDoubleBuffered(true) should be idempotent when buffer already exists")
		}

		if err := w.SetDoubleBuffered(false); err != nil {
			t.Fatalf("SetDoubleBuffered(false) failed: %v", err)
		}
		if w.paintBuffer != nil {
			t.Fatal("paintBuffer should be nil after SetDoubleBuffered(false)")
		}

		if err := w.SetDoubleBuffered(false); err != nil {
			t.Fatalf("second SetDoubleBuffered(false) failed: %v", err)
		}

		if err := w.Destroy(); err != nil {
			t.Fatalf("Destroy failed: %v", err)
		}
	}, nil)
}

func TestSetDoubleBufferedPaintRoute(t *testing.T) {
	phase := 0

	Run(func(app *app.App) {
		w := newDoubleBufferedTestWindow(t, app)

		w.onPaint = func(evt *events.PaintEvent) {
			dc, err := evt.Begin()
			if err != nil {
				t.Fatalf("evt.Begin failed: %v", err)
			}

			switch phase {
			case 0:
				if _, ok := dc.(*events.WindowPaintDC); !ok {
					t.Fatalf("expected WindowPaintDC when double buffering is off, got %T", dc)
				}
				if err := w.SetDoubleBuffered(true); err != nil {
					t.Fatalf("SetDoubleBuffered(true) failed: %v", err)
				}
				phase = 1
			case 1:
				if _, ok := dc.(*events.BufferedPaintDC); !ok {
					t.Fatalf("expected BufferedPaintDC when double buffering is on, got %T", dc)
				}
			default:
				t.Fatalf("unexpected paint phase: %d", phase)
			}
		}

		if _, err := win32.SendMessageW(w.HWND(), win32.WM_PAINT, 0, 0); err != nil {
			t.Fatalf("SendMessageW WM_PAINT (phase 0) failed: %v", err)
		}
		if _, err := win32.SendMessageW(w.HWND(), win32.WM_PAINT, 0, 0); err != nil {
			t.Fatalf("SendMessageW WM_PAINT (phase 1) failed: %v", err)
		}
		if phase != 1 {
			t.Fatalf("expected phase 1 after two paints, got %d", phase)
		}

		if err := w.Destroy(); err != nil {
			t.Fatalf("Destroy failed: %v", err)
		}
	}, nil)
}

func TestSetDoubleBufferedOnWMSizeResizesBuffer(t *testing.T) {
	Run(func(app *app.App) {
		w := newDoubleBufferedTestWindow(t, app)

		if err := w.SetDoubleBuffered(true); err != nil {
			t.Fatalf("SetDoubleBuffered(true) failed: %v", err)
		}
		if w.paintBuffer == nil {
			t.Fatal("paintBuffer should not be nil")
		}

		oldW := w.paintBuffer.Width()
		oldH := w.paintBuffer.Height()
		targetW := oldW + 120
		targetH := oldH + 90

		sz := win32util.EventSize(0)
		sz.SetWidth(int16(targetW))
		sz.SetHeight(int16(targetH))

		if _, err := win32.SendMessageW(w.HWND(), win32.WM_SIZE, win32.WPARAM(win32.SIZE_RESTORED), win32.LPARAM(sz)); err != nil {
			t.Fatalf("SendMessageW WM_SIZE failed: %v", err)
		}

		if w.paintBuffer.Width() < targetW || w.paintBuffer.Height() < targetH {
			t.Fatalf("buffer size not updated, got (%d, %d), want at least (%d, %d)", w.paintBuffer.Width(), w.paintBuffer.Height(), targetW, targetH)
		}

		if err := w.Destroy(); err != nil {
			t.Fatalf("Destroy failed: %v", err)
		}
	}, nil)
}

func TestSetDoubleBufferedOnDisplayChangeRecreatesBuffer(t *testing.T) {
	Run(func(app *app.App) {
		w := newDoubleBufferedTestWindow(t, app)

		if err := w.SetDoubleBuffered(true); err != nil {
			t.Fatalf("SetDoubleBuffered(true) failed: %v", err)
		}
		if w.paintBuffer == nil {
			t.Fatal("paintBuffer should not be nil")
		}
		oldBuffer := w.paintBuffer

		if _, err := win32.SendMessageW(w.HWND(), win32.WM_DISPLAYCHANGE, 0, 0); err != nil {
			t.Fatalf("SendMessageW WM_DISPLAYCHANGE failed: %v", err)
		}
		if w.paintBuffer == nil {
			t.Fatal("paintBuffer should not be nil after WM_DISPLAYCHANGE")
		}
		if w.paintBuffer == oldBuffer {
			t.Fatal("paintBuffer should be recreated after WM_DISPLAYCHANGE")
		}

		if err := w.Destroy(); err != nil {
			t.Fatalf("Destroy failed: %v", err)
		}
	}, nil)
}

func TestDestroyCleansPaintBuffer(t *testing.T) {
	Run(func(app *app.App) {
		w := newDoubleBufferedTestWindow(t, app)

		if err := w.SetDoubleBuffered(true); err != nil {
			t.Fatalf("SetDoubleBuffered(true) failed: %v", err)
		}
		if w.paintBuffer == nil {
			t.Fatal("paintBuffer should not be nil")
		}

		if err := w.Destroy(); err != nil {
			t.Fatalf("Destroy failed: %v", err)
		}
		if w.paintBuffer != nil {
			t.Fatal("paintBuffer should be nil after Destroy")
		}
	}, nil)
}
