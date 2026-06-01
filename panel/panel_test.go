package panel_test

import (
	"testing"

	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/panel"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

type MyPanel struct {
	panel.Panel
}

func TestSimple(t *testing.T) {
	ret := gw.Run(func(app *app.App) {
		parent := window.New(&window.Spec{
			OnDestroy: func() { app.Quit(1) },
		})
		defer parent.Close()

		ctrl := panel.New(&panel.Spec{
			Parent: parent,
			X:      metrics.Px(10),
			Y:      metrics.Px(11),
			Width:  metrics.Px(100),
			Height: metrics.Px(50),
		})

		hwnd := ctrl.HWND()
		if hwnd == 0 {
			t.Fatal("Panel HWND is 0")
		}

		clientRect, err := ctrl.GetClientRect()
		if err != nil {
			t.Fatalf("GetClientRect failed: %v", err)
		}
		if *clientRect != (win32.RECT{Right: 100, Bottom: 50}) {
			t.Fatalf("Unexpected client rect: %+v", *clientRect)
		}

		if got := ctrl.BackgroundColor(); got != win32.COLORREF(win32.GetSysColor(win32.COLOR_WINDOW)) {
			t.Fatalf("Unexpected default background color: %v", got)
		}

		newColor := win32.RGB(1, 2, 3)
		if err := ctrl.SetBackgroundColor(newColor); err != nil {
			t.Fatalf("SetBackgroundColor failed: %v", err)
		}
		if got := ctrl.BackgroundColor(); got != newColor {
			t.Fatalf("Unexpected background color after SetBackgroundColor: %v", got)
		}

	}, nil)

	if ret != 1 {
		t.Fatalf("Unexpected app return value: %d", ret)
	}
}

func TestWrapper(t *testing.T) {
	ret := gw.Run(func(app *app.App) {
		parent := window.New(&window.Spec{
			OnDestroy: func() { app.Quit(1) },
		})
		defer parent.Close()

		const panelClassName = "my panel class name"
		ctrl := gw.Init(&MyPanel{Panel: panel.Panel{Spec: &panel.Spec{
			Parent:    parent,
			ClassName: panelClassName,
			X:         metrics.Px(10),
			Y:         metrics.Px(11),
			Width:     metrics.Px(100),
			Height:    metrics.Px(50),
		}}})

		hwnd := ctrl.HWND()
		if hwnd == 0 {
			t.Fatal("Panel HWND is 0")
		}

		if myPanel, ok := gw.LookupWindow(hwnd).(*MyPanel); !ok || myPanel != ctrl {
			t.Fatal("LookupWindow did not return the correct panel instance")
		}

		var nameBuffer [256]win32.WCHAR
		n, err := win32.GetClassNameW(ctrl.HWND(), &nameBuffer[0], len(nameBuffer))
		if err != nil {
			t.Fatalf("GetClassNameW failed: %v", err)
		}
		if win32util.GoString(&nameBuffer[0], n+1) != panelClassName {
			t.Fatalf("Unexpected class name: %q", win32util.GoString(&nameBuffer[0], n+1))
		}

		clientRect, err := ctrl.GetClientRect()
		if err != nil {
			t.Fatalf("GetClientRect failed: %v", err)
		}
		if *clientRect != (win32.RECT{Right: 100, Bottom: 50}) {
			t.Fatalf("Unexpected client rect: %+v", *clientRect)
		}

		if got := ctrl.BackgroundColor(); got != win32.COLORREF(win32.GetSysColor(win32.COLOR_WINDOW)) {
			t.Fatalf("Unexpected default background color: %v", got)
		}

		newColor := win32.RGB(10, 20, 30)
		if err := ctrl.SetBackgroundColor(newColor); err != nil {
			t.Fatalf("SetBackgroundColor failed: %v", err)
		}
		if got := ctrl.BackgroundColor(); got != newColor {
			t.Fatalf("Unexpected background color after SetBackgroundColor: %v", got)
		}

	}, nil)

	if ret != 1 {
		t.Fatalf("Unexpected app return value: %d", ret)
	}
}
