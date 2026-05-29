package win32util_test

import (
	"slices"
	"testing"

	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

func mustPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: expected panic but did not panic", name)
		}
	}()
	f()
}

func TestGoString(t *testing.T) {
	var buf []win32.WCHAR
	win32util.CString("abc", &buf)
	if !slices.Equal(buf, []win32.WCHAR{'a', 'b', 'c', 0}) {
		t.Fatal(buf)
	}
	if s := win32util.GoString(&buf[0], len(buf)); s != "abc" {
		t.Fatal(s)
	}

	win32util.CString("", &buf)
	if !slices.Equal(buf, []win32.WCHAR{0}) {
		t.Fatal(buf)
	}
	if s := win32util.GoString(&buf[0], len(buf)); s != "" {
		t.Fatal(s)
	}

	win32util.CString("中文abc", &buf)
	if s := win32util.GoString(&buf[0], len(buf)); s != "中文abc" {
		t.Fatal(s)
	}
}

func TestCStrLen(t *testing.T) {
	// Non-ASCII string with null terminator.
	var buf []win32.WCHAR
	win32util.CString("中文abc", &buf)
	if n := win32util.CStrLen(&buf[0], len(buf)); n != 5 {
		t.Fatalf("got %d, want %d", n, len(buf)-1)
	}

	// Normal string with null terminator.
	buf = []win32.WCHAR{'a', 'b', 'c', 0}
	if n := win32util.CStrLen(&buf[0], len(buf)); n != 3 {
		t.Fatalf("got %d, want 3", n)
	}

	// Empty string: only the null terminator.
	buf = []win32.WCHAR{0}
	if n := win32util.CStrLen(&buf[0], len(buf)); n != 0 {
		t.Fatalf("got %d, want 0", n)
	}

	// No null terminator within bufSize: must panic.
	mustPanic(t, "no null terminator within bufSize", func() {
		buf = []win32.WCHAR{'a', 'b', 'c'}
		win32util.CStrLen(&buf[0], len(buf))
	})

	// bufSize == 0 with non-nil str: no elements to search, must panic.
	mustPanic(t, "bufSize == 0 with non-nil str", func() {
		buf = []win32.WCHAR{'a', 'b', 'c', 0}
		win32util.CStrLen(&buf[0], 0)
	})

	// str == nil with bufSize == 0: returns 0.
	if n := win32util.CStrLen(nil, 0); n != 0 {
		t.Fatalf("got %d, want 0", n)
	}

	// Null terminator at the first position with bufSize == 1.
	buf = []win32.WCHAR{0, 'x'}
	if n := win32util.CStrLen(&buf[0], 1); n != 0 {
		t.Fatalf("got %d, want 0", n)
	}

	// Larger bufSize than the actual string length: should return the actual string length.
	buf = []win32.WCHAR{'a', 'b', 'c', 0}
	if n := win32util.CStrLen(&buf[0], 1024); n != 3 {
		t.Fatalf("got %d, want %d", n, 3)
	}
}

func TestCopyCString(t *testing.T) {
	dest := make([]win32.WCHAR, 3)
	src := []win32.WCHAR{'a', 'b', 0}
	if n := win32util.CopyCString(dest, src); n != 3 {
		t.Fatal(n)
	} else if !slices.Equal(dest, src) {
		t.Fatal(dest)
	}

	src = []win32.WCHAR{'a', 'b', 'c', 0}
	if n := win32util.CopyCString(dest, src); n != 3 {
		t.Fatal(n)
	} else if !slices.Equal(dest, []win32.WCHAR{'a', 'b', 0}) {
		t.Fatal(dest)
	}

	src = []win32.WCHAR{0}
	if n := win32util.CopyCString(dest, src); n != 1 {
		t.Fatal(n)
	} else if dest[0] != 0 {
		t.Fatal(dest)
	}
}
