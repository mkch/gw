package win32util_test

import (
	"slices"
	"testing"

	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

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
