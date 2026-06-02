package button_test

import (
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

type MyButton struct {
	button.Button
}

func TestSimple(t *testing.T) {
	ret := gw.Run(func(app *app.App) {
		var onClickCalled bool
		parent := chkerr.Must(window.New(&window.Spec{
			OnDestroy: func() { app.Quit(1) },
		}))
		defer parent.Close()

		btn := chkerr.Must(button.New(&button.Spec{
			Parent:  parent,
			Text:    "Click me",
			X:       metrics.Px(10),
			Y:       metrics.Px(11),
			Width:   metrics.Px(100),
			Height:  metrics.Px(50),
			Style:   win32.BS_LEFT,
			OnClick: func() { onClickCalled = true },
		}))
		hwnd := btn.HWND()
		if hwnd == 0 {
			t.Fatal("Button HWND is 0")
		}

		clientRect, err := btn.GetClientRect()
		if err != nil {
			t.Fatalf("GetClientRect failed: %v", err)
		}
		if *clientRect != (win32.RECT{Right: 100, Bottom: 50}) {
			t.Fatalf("Unexpected client rect: %+v", *clientRect)
		}

		win32util.SimulateControlCommand(parent.HWND(), hwnd, 0)
		if !onClickCalled {
			t.Fatal("Button OnClick was not called")
		}

		onClickCalled = false

		btn.SetOnClickListener(nil)
		win32util.SimulateControlCommand(parent.HWND(), hwnd, 0)
		if onClickCalled {
			t.Fatal("Button OnClick was called after being set to nil")
		}

		btn.SetOnClickListener(func() { onClickCalled = true })
		win32util.SimulateControlCommand(parent.HWND(), hwnd, 0)
		if !onClickCalled {
			t.Fatal("Button OnClick was not called after being set again")
		}

	}, nil)

	if ret != 1 {
		t.Fatalf("Unexpected app return value: %d", ret)
	}
}

func TestWrapper(t *testing.T) {
	ret := gw.Run(func(app *app.App) {
		var onClickCalled bool
		parent := chkerr.Must(window.New(&window.Spec{
			OnDestroy: func() { app.Quit(1) },
		}))
		defer parent.Close()

		btn := chkerr.Must(gw.Init(&MyButton{Button: button.Button{Spec: &button.Spec{
			Parent:  parent,
			Text:    "Click me",
			X:       metrics.Px(10),
			Y:       metrics.Px(11),
			Width:   metrics.Px(100),
			Height:  metrics.Px(50),
			Style:   win32.BS_LEFT,
			OnClick: func() { onClickCalled = true },
		}}}))

		hwnd := btn.HWND()
		if hwnd == 0 {
			t.Fatal("Button HWND is 0")
		}

		if myBtn, ok := gw.LookupWindow(hwnd).(*MyButton); !ok || myBtn != btn {
			t.Fatal("LookupWindow did not return the correct button instance")
		}

		clientRect, err := btn.GetClientRect()
		if err != nil {
			t.Fatalf("GetClientRect failed: %v", err)
		}
		if *clientRect != (win32.RECT{Right: 100, Bottom: 50}) {
			t.Fatalf("Unexpected client rect: %+v", *clientRect)
		}

		win32util.SimulateControlCommand(parent.HWND(), hwnd, 0)
		if !onClickCalled {
			t.Fatal("Button OnClick was not called")
		}

		onClickCalled = false

		btn.SetOnClickListener(nil)
		win32util.SimulateControlCommand(parent.HWND(), hwnd, 0)
		if onClickCalled {
			t.Fatal("Button OnClick was called after being set to nil")
		}

		btn.SetOnClickListener(func() { onClickCalled = true })
		win32util.SimulateControlCommand(parent.HWND(), hwnd, 0)
		if !onClickCalled {
			t.Fatal("Button OnClick was not called after being set again")
		}

	}, nil)

	if ret != 1 {
		t.Fatalf("Unexpected app return value: %d", ret)
	}
}
