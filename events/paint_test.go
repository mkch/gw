package events

import (
	"testing"

	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/win32"
)

type fakePaintDC struct {
	rect     win32.RECT
	endCalls int
	endErr   error
}

func (f *fakePaintDC) HDC() win32.HDC {
	return 1
}

func (f *fakePaintDC) EraseBackground() bool {
	return false
}

func (f *fakePaintDC) Rect() *win32.RECT {
	return &f.rect
}

func (f *fakePaintDC) endPaint() error {
	f.endCalls++
	return f.endErr
}

func TestNewPaintEvent(t *testing.T) {
	buf := &paint.Buffer{}
	evt, release := NewPaintEvent(100, 200, 300, buf)
	t.Cleanup(func() {
		if err := release(); err != nil {
			t.Fatalf("release failed: %v", err)
		}
	})

	if evt.hwnd != 100 {
		t.Fatalf("hwnd: got %v, want 100", evt.hwnd)
	}
	if evt.wParam != 200 {
		t.Fatalf("wParam: got %v, want 200", evt.wParam)
	}
	if evt.lParam != 300 {
		t.Fatalf("lParam: got %v, want 300", evt.lParam)
	}
	if evt.extBuffer != buf {
		t.Fatal("extBuffer mismatch")
	}
}

func TestPaintEventCallDefProc(t *testing.T) {
	evt, release := NewPaintEvent(11, 22, 33, nil)
	t.Cleanup(func() {
		if err := release(); err != nil {
			t.Fatalf("release failed: %v", err)
		}
	})

	called := 0
	ret := evt.CallDefProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
		called++
		if hwnd != 11 || message != win32.WM_PAINT || wParam != 22 || lParam != 33 {
			t.Fatalf("CallDefProc forwarded unexpected args: hwnd=%v message=%v wParam=%v lParam=%v", hwnd, message, wParam, lParam)
		}
		return 77
	})
	if ret != 77 {
		t.Fatalf("first CallDefProc return: got %v, want 77", ret)
	}
	if called != 1 {
		t.Fatalf("proc call count after first call: got %d, want 1", called)
	}

	ret = evt.CallDefProc(func(win32.HWND, win32.UINT, win32.WPARAM, win32.LPARAM) win32.LRESULT {
		called++
		return 99
	})
	if ret != 0 {
		t.Fatalf("second CallDefProc return: got %v, want 0", ret)
	}
	if called != 1 {
		t.Fatalf("proc call count after second call: got %d, want 1", called)
	}

	panicked := false
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		_, _ = evt.Begin()
	}()
	if !panicked {
		t.Fatal("Begin should panic after CallDefProc")
	}
}

func TestPaintEventCallDefProcAfterBegin(t *testing.T) {
	fake := &fakePaintDC{}
	evt := &PaintEvent{dc: fake}

	if _, err := evt.Begin(); err != nil {
		t.Fatalf("Begin failed: %v", err)
	}

	called := 0
	ret := evt.CallDefProc(func(win32.HWND, win32.UINT, win32.WPARAM, win32.LPARAM) win32.LRESULT {
		called++
		return 88
	})

	if ret != 0 {
		t.Fatalf("CallDefProc return after Begin: got %v, want 0", ret)
	}
	if called != 0 {
		t.Fatalf("proc should not be called after Begin, got call count %d", called)
	}
}

func TestPaintEventBeginReturnsExistingDC(t *testing.T) {
	fake := &fakePaintDC{}
	evt := &PaintEvent{dc: fake}

	dc, err := evt.Begin()
	if err != nil {
		t.Fatalf("Begin failed: %v", err)
	}
	if dc != fake {
		t.Fatal("Begin should return the existing dc")
	}
}

func TestPaintEventEnd(t *testing.T) {
	fake := &fakePaintDC{}
	evt := &PaintEvent{dc: fake}

	if err := evt.end(); err != nil {
		t.Fatalf("end failed: %v", err)
	}
	if fake.endCalls != 1 {
		t.Fatalf("endPaint call count: got %d, want 1", fake.endCalls)
	}
	if evt.dc != nil {
		t.Fatal("dc should be cleared after end")
	}
}
