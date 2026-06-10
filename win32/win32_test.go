package win32_test

import (
	"slices"
	"testing"
	"unsafe"

	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"golang.org/x/sys/windows"
)

func TestMessageLoop(t *testing.T) {
	var msg win32.MSG

	const exitCode = 100
	for {
		win32.PostQuitMessage(exitCode)

		r := win32.GetMessageW(&msg, 0, 0, 0)
		if r == 0 {
			break
		}

		win32.TranslateMessage(&msg)

		win32.DispatchMessageW(&msg)
	}

	if msg.Message != win32.WM_QUIT {
		t.Fatal(msg)
	}
	if msg.WParam != exitCode {
		t.Fatal(msg.WParam)
	}
}

func TestCreateWindow(t *testing.T) {
	module, err := win32.GetModuleHandleW[win32.HINSTANCE](nil)
	if module == 0 {
		t.Fatal(err)
	}
	var className []win32.WCHAR
	win32util.CString("cls1", &className)
	var cls = win32.WNDCLASSEXW{
		Size:      win32.UINT(unsafe.Sizeof(win32.WNDCLASSEXW{})),
		WndProc:   windows.NewCallback(testCreateWindowWndproc),
		Instance:  module,
		ClassName: &className[0],
	}
	if _, err = win32.RegisterClassExW(&cls); err != nil {
		t.Fatal(err)
	}
	var windowName []win32.WCHAR
	win32util.CString("A window", &windowName)
	hwnd, err := win32.CreateWindowExW(0, &className[0], &windowName[0], win32.WS_OVERLAPPEDWINDOW, win32.CW_USEDEFAULT, win32.CW_USEDEFAULT, 300, 200, 0, 0, 0, 0)
	if hwnd == 0 {
		t.Fatal(err)
	}
}

func testCreateWindowWndproc(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	return win32.DefWindowProcW(hwnd, message, wParam, lParam)
}

func TestAlignSlice(t *testing.T) {
	empty := []byte{}
	empty = win32.AlignSlice(empty)
	if len(empty) != 0 {
		t.Errorf("empty slice should remain empty, got %v", empty)
	}

	buf := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	for i := range len(buf) - 1 {
		s := buf[i:]
		orig := slices.Clone(s)
		s = win32.AlignSlice(s)
		if uintptr(unsafe.Pointer(&s[0]))%unsafe.Sizeof(uintptr(0)) != 0 {
			t.Errorf("slice starting at index %d is not aligned: %p", i, &s[0])
		}
		if !slices.Equal(s, orig) {
			t.Errorf("aligning slice starting at index %d changed its contents: got %v, want %v", i, s, orig)
		}
	}

	type S struct {
		a int
		b byte
	}

	buf2 := []S{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}, {11, 12}, {13, 14}, {15, 16}}
	for i := range len(buf2) - 1 {
		s := buf2[i:]
		orig := slices.Clone(s)
		s = win32.AlignSlice(s)
		if uintptr(unsafe.Pointer(&s[0]))%unsafe.Sizeof(uintptr(0)) != 0 {
			t.Errorf("slice starting at index %d is not aligned: %p", i, &s[0])
		}
		if !slices.Equal(s, orig) {
			t.Errorf("aligning slice starting at index %d changed its contents: got %v, want %v", i, s, orig)
		}
	}
}
