package static_test

import (
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

type MyStatic struct {
	static.Static
}

func TestSimple(t *testing.T) {
	ret := gw.Run(func(app *app.App) {
		parent := chkerr.Must(window.New(&window.Spec{
			OnDestroy: func() { app.Quit(1) },
		}))
		defer parent.Close()

		ctrl := chkerr.Must(static.New(&static.Spec{
			Parent: parent,
			Text:   "Hello",
			X:      metrics.Px(10),
			Y:      metrics.Px(11),
			Width:  metrics.Px(100),
			Height: metrics.Px(50),
			Style:  static.SS_CENTER,
		}))

		hwnd := ctrl.HWND()
		if hwnd == 0 {
			t.Fatal("Static HWND is 0")
		}

		clientRect, err := ctrl.GetClientRect()
		if err != nil {
			t.Fatalf("GetClientRect failed: %v", err)
		}
		if *clientRect != (win32.RECT{Right: 100, Bottom: 50}) {
			t.Fatalf("Unexpected client rect: %+v", *clientRect)
		}

		text, err := ctrl.Text()
		if err != nil {
			t.Fatalf("Get text failed: %v", err)
		}
		if text != "Hello" {
			t.Fatalf("Unexpected text: %q", text)
		}

		if got := *static.BackgroundColor(ctrl); got != win32.COLORREF(win32.GetSysColor(win32.COLOR_WINDOW)) {
			t.Fatalf("Unexpected default background color: %v", got)
		}

		newColor := win32.RGB(1, 2, 3)
		if err := ctrl.SetBackgroundColor(newColor); err != nil {
			t.Fatalf("SetBackgroundColor failed: %v", err)
		}
		if got := *static.BackgroundColor(ctrl); got != newColor {
			t.Fatalf("Unexpected background color after SetBackgroundColor: %v", got)
		}

		if err := ctrl.SetText("World"); err != nil {
			t.Fatalf("SetText failed: %v", err)
		}
		text, err = ctrl.Text()
		if err != nil {
			t.Fatalf("Get text after SetText failed: %v", err)
		}
		if text != "World" {
			t.Fatalf("Unexpected text after SetText: %q", text)
		}

	}, nil)

	if ret != 1 {
		t.Fatalf("Unexpected app return value: %d", ret)
	}
}

func TestWrapper(t *testing.T) {
	ret := gw.Run(func(app *app.App) {
		parent := chkerr.Must(window.New(&window.Spec{
			OnDestroy: func() { app.Quit(1) },
		}))
		defer parent.Close()

		ctrl := chkerr.Must(gw.Init(&MyStatic{Static: static.Static{Spec: &static.Spec{
			Parent: parent,
			Text:   "Hello",
			X:      metrics.Px(10),
			Y:      metrics.Px(11),
			Width:  metrics.Px(100),
			Height: metrics.Px(50),
			Style:  static.SS_CENTER,
		}}}))

		hwnd := ctrl.HWND()
		if hwnd == 0 {
			t.Fatal("Static HWND is 0")
		}

		if myStatic, ok := gw.LookupWindow(hwnd).(*MyStatic); !ok || myStatic != ctrl {
			t.Fatal("LookupWindow did not return the correct static instance")
		}

		clientRect, err := ctrl.GetClientRect()
		if err != nil {
			t.Fatalf("GetClientRect failed: %v", err)
		}
		if *clientRect != (win32.RECT{Right: 100, Bottom: 50}) {
			t.Fatalf("Unexpected client rect: %+v", *clientRect)
		}

		text, err := ctrl.Text()
		if err != nil {
			t.Fatalf("Get text failed: %v", err)
		}
		if text != "Hello" {
			t.Fatalf("Unexpected text: %q", text)
		}

		if got := *static.BackgroundColor(&ctrl.Static); got != win32.COLORREF(win32.GetSysColor(win32.COLOR_WINDOW)) {
			t.Fatalf("Unexpected default background color: %v", got)
		}

		newColor := win32.RGB(10, 20, 30)
		if err := ctrl.SetBackgroundColor(newColor); err != nil {
			t.Fatalf("SetBackgroundColor failed: %v", err)
		}
		if got := *static.BackgroundColor(&ctrl.Static); got != newColor {
			t.Fatalf("Unexpected background color after SetBackgroundColor: %v", got)
		}

		if err := ctrl.SetText("World"); err != nil {
			t.Fatalf("SetText failed: %v", err)
		}
		text, err = ctrl.Text()
		if err != nil {
			t.Fatalf("Get text after SetText failed: %v", err)
		}
		if text != "World" {
			t.Fatalf("Unexpected text after SetText: %q", text)
		}

	}, nil)

	if ret != 1 {
		t.Fatalf("Unexpected app return value: %d", ret)
	}
}
