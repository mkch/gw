package win32util_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const testClassName = "gw_win32util_setprop_cls"
const testPropName = "gw_win32util_setprop_prop"

func testWndProc(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	return win32.DefWindowProcW(hwnd, message, wParam, lParam)
}

func TestMain(m *testing.M) {
	if _, err := win32util.RegisterClass(&win32util.WndClass{
		ClassName: testClassName,
		WndProc:   testWndProc,
	}); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func cString(str string) []win32.WCHAR {
	var buf []win32.WCHAR
	win32util.CString(str, &buf)
	return buf
}

func mustCreateTestWindow(t *testing.T) win32.HWND {
	t.Helper()

	hwnd, err := win32util.CreateWindow(&win32util.Wnd{
		ClassName: testClassName,
		Style:     win32.WS_OVERLAPPEDWINDOW,
		X:         win32.CW_USEDEFAULT,
		Y:         win32.CW_USEDEFAULT,
		Width:     160,
		Height:    120,
	})
	if err != nil {
		t.Fatal(err)
	}
	if hwnd == 0 {
		t.Fatal("CreateWindow returned null hwnd")
	}

	t.Cleanup(func() {
		if hwnd != 0 {
			_ = win32.DestroyWindow(hwnd)
		}
	})

	return hwnd
}

func mustSetIntPropFromLocalVar(t *testing.T, hwnd win32.HWND, prop *win32util.WindowProp[int], value int) {
	t.Helper()
	local := value
	if err := prop.Set(hwnd, &local); err != nil {
		t.Fatal(err)
	}
}

func getIntPropAfterGC(hwnd win32.HWND, prop *win32util.WindowProp[int]) *int {
	runtime.GC()
	runtime.GC()
	runtime.GC()
	return prop.Get(hwnd)
}

func TestSetWindowProp_BasicSetGet(t *testing.T) {
	hwnd := mustCreateTestWindow(t)

	prop := win32util.NewWindowProp[int](testPropName)
	value := 42
	if err := prop.Set(hwnd, &value); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = prop.Set(hwnd, nil)
	})

	got := prop.Get(hwnd)
	if got == nil {
		t.Fatal("WindowProp returned nil")
	}
	if *got != value {
		t.Fatalf("got %d, want %d", *got, value)
	}
}

func TestSetWindowProp_ReplaceValue(t *testing.T) {
	hwnd := mustCreateTestWindow(t)
	prop := win32util.NewWindowProp[int](testPropName)

	v1 := 100
	if err := prop.Set(hwnd, &v1); err != nil {
		t.Fatal(err)
	}
	v2 := 200
	if err := prop.Set(hwnd, &v2); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = prop.Set(hwnd, nil)
	})

	got := prop.Get(hwnd)
	if got == nil {
		t.Fatal("WindowProp returned nil")
	}
	if *got != v2 {
		t.Fatalf("got %d, want %d", *got, v2)
	}
}

func TestSetWindowProp_RemoveWithNil(t *testing.T) {
	hwnd := mustCreateTestWindow(t)

	value := 7
	prop := win32util.NewWindowProp[int](testPropName)
	if err := prop.Set(hwnd, &value); err != nil {
		t.Fatal(err)
	}
	if err := prop.Set(hwnd, nil); err != nil {
		t.Fatal(err)
	}

	got := prop.Get(hwnd)
	if got != nil {
		t.Fatalf("got %v, want nil", *got)
	}
}

func TestWindowProp_NotFoundReturnsNil(t *testing.T) {
	hwnd := mustCreateTestWindow(t)
	prop := win32util.NewWindowProp[int](testPropName)

	got := prop.Get(hwnd)
	if got != nil {
		t.Fatalf("got %v, want nil", *got)
	}
}

func TestInvalidHandleBehavior(t *testing.T) {
	prop := win32util.NewWindowProp[int](testPropName)
	value := 9

	invalidHwnd := win32.HWND(0)
	if err := prop.Set(invalidHwnd, &value); err == nil {
		t.Fatal("SetWindowProp should return error for invalid hwnd")
	}
	if got := prop.Get(invalidHwnd); got != nil {
		t.Fatalf("WindowProp returned %v, want nil", *got)
	}
}

func TestSetWindowProp_GCKeepsDataAlive(t *testing.T) {
	hwnd := mustCreateTestWindow(t)
	prop := win32util.NewWindowProp[int](testPropName)
	t.Cleanup(func() {
		_ = prop.Set(hwnd, nil)
	})

	mustSetIntPropFromLocalVar(t, hwnd, prop, 13579)
	got := getIntPropAfterGC(hwnd, prop)
	if got == nil {
		t.Fatal("WindowProp returned nil after runtime.GC")
	}
	if *got != 13579 {
		t.Fatalf("got %d, want %d", *got, 13579)
	}
}
