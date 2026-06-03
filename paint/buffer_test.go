package paint

import (
	"errors"
	"testing"

	"github.com/mkch/gw/win32"
)

func newScreenBuffer(t *testing.T, width, height int) *Buffer {
	t.Helper()

	buf, err := NewBuffer(ClientDC(0), width, height)
	if err != nil {
		t.Skipf("NewBuffer failed, skip: %v", err)
	}

	t.Cleanup(func() {
		if err := buf.Destroy(); err != nil {
			t.Fatalf("Destroy failed: %v", err)
		}
	})

	return buf
}

func TestBufferResizeNoopForSmallShrink(t *testing.T) {
	buf := newScreenBuffer(t, 100, 80)
	oldHDC := buf.hdc
	oldBitmap := buf.bitmap
	oldOldBitmap := buf.oldBitmap

	if err := buf.Resize(90, 70); err != nil {
		t.Fatalf("Resize failed: %v", err)
	}
	if buf.Width() != 90 || buf.Height() != 70 {
		t.Fatalf("unexpected logical size after noop resize: got (%d,%d), want (90,70)", buf.Width(), buf.Height())
	}
	if buf.hdc != oldHDC || buf.bitmap != oldBitmap || buf.oldBitmap != oldOldBitmap {
		t.Fatal("noop resize should not recreate GDI resources")
	}
}

func TestNewBufferBasic(t *testing.T) {
	buf, err := NewBuffer(ClientDC(0), 64, 48)
	if err != nil {
		t.Skipf("NewBuffer failed, skip: %v", err)
	}
	defer func() {
		if err := buf.Destroy(); err != nil {
			t.Fatalf("Destroy failed: %v", err)
		}
	}()

	if buf.Width() != 64 || buf.Height() != 48 {
		t.Fatalf("unexpected buffer size: got (%d,%d), want (64,48)", buf.Width(), buf.Height())
	}
	if buf.HDC() == 0 {
		t.Fatal("buffer HDC should not be zero")
	}
}

func TestNewBufferWithProviderError(t *testing.T) {
	expectedErr := errors.New("provider error")
	provider := DCProviderFunc(func(func(dc win32.HDC) error) error {
		return expectedErr
	})

	buf, err := NewBuffer(provider, 10, 10)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("NewBuffer error: got %v, want %v", err, expectedErr)
	}
	if buf != nil {
		t.Fatal("buffer should be nil when provider fails")
	}
}

func TestBufferResizeGrow(t *testing.T) {
	buf := newScreenBuffer(t, 64, 48)

	if err := buf.Resize(200, 160); err != nil {
		t.Fatalf("Resize failed: %v", err)
	}
	if buf.Width() != 200 || buf.Height() != 160 {
		t.Fatalf("unexpected size after grow resize: got (%d,%d), want (200,160)", buf.Width(), buf.Height())
	}
	if buf.capWidth != 220 || buf.capHeight != 176 {
		t.Fatalf("unexpected bitmap size after grow resize: got (%d,%d), want (220,176)", buf.capWidth, buf.capHeight)
	}
}

func TestBufferResizeWhenOneDimensionGrows(t *testing.T) {
	buf := newScreenBuffer(t, 100, 80)

	if err := buf.Resize(100, 120); err != nil {
		t.Fatalf("Resize failed: %v", err)
	}
	if buf.Width() != 100 || buf.Height() != 120 {
		t.Fatalf("unexpected size after one-dimension grow: got (%d,%d), want (100,120)", buf.Width(), buf.Height())
	}
	if buf.capWidth != 100 || buf.capHeight != 132 {
		t.Fatalf("unexpected bitmap size after one-dimension grow: got (%d,%d), want (100,132)", buf.capWidth, buf.capHeight)
	}
}

func TestBufferResizeWithinHeadroomDoesNotRecreate(t *testing.T) {
	buf := newScreenBuffer(t, 100, 80)

	if err := buf.Resize(200, 160); err != nil {
		t.Fatalf("first Resize failed: %v", err)
	}
	oldHDC := buf.hdc
	oldBitmap := buf.bitmap
	oldOldBitmap := buf.oldBitmap

	if err := buf.Resize(210, 170); err != nil {
		t.Fatalf("second Resize failed: %v", err)
	}
	if buf.Width() != 210 || buf.Height() != 170 {
		t.Fatalf("unexpected logical size after second resize: got (%d,%d), want (210,170)", buf.Width(), buf.Height())
	}
	if buf.hdc != oldHDC || buf.bitmap != oldBitmap || buf.oldBitmap != oldOldBitmap {
		t.Fatal("resize within headroom should not recreate GDI resources")
	}
}

func TestBufferResizeShrinkToQuarterBoundary(t *testing.T) {
	buf := newScreenBuffer(t, 200, 160)

	// Quarter boundary should trigger real resize:
	// width*4 > oldWidth is false when width == oldWidth/4.
	if err := buf.Resize(50, 40); err != nil {
		t.Fatalf("Resize failed: %v", err)
	}
	if buf.Width() != 50 || buf.Height() != 40 {
		t.Fatalf("unexpected size after shrink resize: got (%d,%d), want (50,40)", buf.Width(), buf.Height())
	}
	if buf.capWidth != 50 || buf.capHeight != 40 {
		t.Fatalf("unexpected bitmap size after shrink resize: got (%d,%d), want (50,40)", buf.capWidth, buf.capHeight)
	}
}

func TestBufferResizeContinuousShrinkEventuallyRecreates(t *testing.T) {
	buf := newScreenBuffer(t, 1000, 800)
	oldHDC := buf.hdc
	oldBitmap := buf.bitmap
	oldOldBitmap := buf.oldBitmap

	if err := buf.Resize(400, 320); err != nil {
		t.Fatalf("first Resize failed: %v", err)
	}
	if buf.hdc != oldHDC || buf.bitmap != oldBitmap || buf.oldBitmap != oldOldBitmap {
		t.Fatal("first shrink above quarter should not recreate GDI resources")
	}

	if err := buf.Resize(200, 160); err != nil {
		t.Fatalf("second Resize failed: %v", err)
	}
	if buf.Width() != 200 || buf.Height() != 160 {
		t.Fatalf("unexpected logical size after continuous shrink: got (%d,%d), want (200,160)", buf.Width(), buf.Height())
	}
	if buf.capWidth != 200 || buf.capHeight != 160 {
		t.Fatalf("unexpected bitmap size after continuous shrink: got (%d,%d), want (200,160)", buf.capWidth, buf.capHeight)
	}
}

func TestBufferDestroy(t *testing.T) {
	buf, err := NewBuffer(ClientDC(0), 32, 24)
	if err != nil {
		t.Skipf("NewBuffer failed, skip: %v", err)
	}

	if err := buf.Destroy(); err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}
}
