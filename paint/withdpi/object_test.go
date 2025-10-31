package withdpi

import (
	"testing"
	"unsafe"

	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
)

type LogFont = LogStruct[win32.LOGFONTW, win32.HFONT]
type Font = Object[win32.LOGFONTW, win32.HFONT]

func NewLogFont(l *win32.LOGFONTW, DPI win32.UINT) *LogFont {
	return NewLogStruct[win32.LOGFONTW, win32.HFONT](l, DPI, func(l *win32.LOGFONTW, oldDPI, newDPI win32.UINT) {
		l.Height = metrics.DPIConv(l.Height, oldDPI, newDPI)
	}, win32.CreateFontIndirectW)
}

func TestCache(t *testing.T) {
	const DPI = win32.USER_DEFAULT_SCREEN_DPI
	var metrics = win32.NONCLIENTMETRICSW{Size: win32.UINT(unsafe.Sizeof(win32.NONCLIENTMETRICSW{}))}
	if err := win32.SystemParametersInfoForDpi(win32.SPI_GETNONCLIENTMETRICS, win32.UINT(unsafe.Sizeof(metrics)), win32.PVOID(&metrics), 0, DPI); err != nil {
		panic(err)
	}

	lf := NewLogFont(&metrics.CaptionFont, 96)
	font168, err := New(lf, 168)
	if err != nil {
		t.Fatal(err)
	}
	if font168.l != lf || font168.dpi != 168 {
		t.Fatal()
	}
	if font168.Handle() == 0 {
		t.Fatal()
	}

	fontCopy := font168.Clone()
	if fontCopy.l != lf || fontCopy.dpi != 168 {
		t.Fatal()
	}
	if fontCopy.Handle() != font168.Handle() {
		t.Fatal()
	}

	if len(lf.cache) != 1 || lf.cache[168] == nil {
		t.Fatal()
	}

	err = fontCopy.ChangeDPI(96)
	if err != nil {
		t.Fatal(err)
	}
	if fontCopy.l != lf || fontCopy.dpi != 96 {
		t.Fatal()
	}

	err = fontCopy.ChangeDPI(192)
	if err != nil {
		t.Fatal(err)
	}
	if fontCopy.l != lf || fontCopy.dpi != 192 {
		t.Fatal()
	}

	if len(lf.cache) != 2 || lf.cache[168] == nil || lf.cache[192] == nil {
		t.Fatal()
	}

	fontCopy.Release()
	font168.Release()
	if len(lf.cache) != 0 {
		t.Fatal()
	}
}
