package edit_test

import (
	"testing"

	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/edit"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

type MyEdit struct {
	edit.Edit
}

func TestSimple(t *testing.T) {
	ret := gw.Run(func(app *app.App) {
		parent := window.New(&window.Spec{
			OnDestroy: func() { app.Quit(1) },
		})
		defer parent.Close()

		ctrl := edit.New(&edit.Spec{
			Parent: parent,
			Text:   "Hello",
			X:      metrics.Px(10),
			Y:      metrics.Px(11),
			Width:  metrics.Px(100),
			Height: metrics.Px(50),
			Style:  edit.ES_LEFT,
		})

		hwnd := ctrl.HWND()
		if hwnd == 0 {
			t.Fatal("Edit HWND is 0")
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
		parent := window.New(&window.Spec{
			OnDestroy: func() { app.Quit(1) },
		})
		defer parent.Close()

		ctrl := &MyEdit{Edit: edit.Edit{Spec: &edit.Spec{
			Parent: parent,
			Text:   "Hello",
			X:      metrics.Px(10),
			Y:      metrics.Px(11),
			Width:  metrics.Px(100),
			Height: metrics.Px(50),
			Style:  edit.ES_LEFT,
		}}}
		gw.Init(ctrl)

		hwnd := ctrl.HWND()
		if hwnd == 0 {
			t.Fatal("Edit HWND is 0")
		}

		if myEdit, ok := gw.LookupWindow(hwnd).(*MyEdit); !ok || myEdit != ctrl {
			t.Fatal("LookupWindow did not return the correct edit instance")
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
