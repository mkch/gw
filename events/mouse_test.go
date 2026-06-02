package events_test

import (
	"sync"
	"testing"

	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const keyTestClassName = "github.com/mkch/gw/events#KeyTest"

var (
	registerOnce sync.Once
	registerErr  error
)

func ensureTestClassRegistered() error {
	registerOnce.Do(func() {
		_, registerErr = win32util.RegisterClass(&win32util.WndClass{
			Style:     win32.CS_DBLCLKS, // Enable double-click messages for testing.
			ClassName: keyTestClassName,
			WndProc:   win32.DefWindowProcW,
		})
	})
	return registerErr
}

// mouseListenerWindow embeds BaseWindowImpl without overriding any On...() methods,
// so the default implementations in BaseWindowImpl call the registered listeners.
type mouseListenerWindow struct {
	gw.BaseWindowImpl
}

// mouseOverrideWindow embeds BaseWindowImpl and overrides all On...Button...() methods
// to capture the received events for verification.
type mouseOverrideWindow struct {
	gw.BaseWindowImpl
	onLButtonDown     *events.MouseClickEvent
	onLButtonUp       *events.MouseClickEvent
	onRButtonDown     *events.MouseClickEvent
	onRButtonUp       *events.MouseClickEvent
	onLButtonDblClick *events.MouseClickEvent
}

func (w *mouseOverrideWindow) OnLButtonDown(event *events.MouseClickEvent) {
	w.onLButtonDown = event
}

func (w *mouseOverrideWindow) OnLButtonUp(event *events.MouseClickEvent) {
	w.onLButtonUp = event
}

func (w *mouseOverrideWindow) OnRButtonDown(event *events.MouseClickEvent) {
	w.onRButtonDown = event
}

func (w *mouseOverrideWindow) OnRButtonUp(event *events.MouseClickEvent) {
	w.onRButtonUp = event
}

func (w *mouseOverrideWindow) OnLButtonDoubleClick(event *events.MouseClickEvent) {
	w.onLButtonDblClick = event
}

// simulateAllMouseEvents posts all five mouse messages to hwnd in order.
func simulateAllMouseEvents(hwnd win32.HWND) {
	var opt win32util.MouseClickOpt
	opt.SetLButton(true)
	// Left button click (down + up).
	win32util.SimulateLButtonClick(hwnd, opt, 0)

	// Right button down + up.
	opt.SetLButton(false)
	opt.SetRButton(true)
	win32util.SimulateRButtonDown(hwnd, 0, 0)
	win32util.SimulateRButtonUp(hwnd, 0, 0)

	// Left button double click.
	opt.SetLButton(true)
	opt.SetRButton(false)
	win32util.SimulateLButtonDoubleClick(hwnd, opt, 0)
}

// TestMouseListeners verifies that all Set...Listener callbacks are invoked when the
// corresponding mouse messages are received.
func TestMouseListeners(t *testing.T) {
	if err := ensureTestClassRegistered(); err != nil {
		t.Fatal(err)
	}

	var wnd mouseListenerWindow
	var lButtonDown, lButtonUp, rButtonDown, rButtonUp, lButtonDblClick *events.MouseClickEvent

	gw.Run(func(a *app.App) {
		hwnd, err := win32util.CreateWindow(&win32util.Wnd{
			ClassName: keyTestClassName,
			Style:     win32.WS_POPUP,
		})
		if err != nil {
			t.Fatal(err)
			return
		}
		if err := gw.Attach(hwnd, &wnd); err != nil {
			t.Fatal(err)
			return
		}

		wnd.SetOnDestroyListener(func() { a.Quit(0) })

		wnd.SetOnLButtonDownListener(func(e events.MouseClickEvent) { lButtonDown = &e })
		wnd.SetOnLButtonUpListener(func(e events.MouseClickEvent) { lButtonUp = &e })
		wnd.SetOnRButtonDownListener(func(e events.MouseClickEvent) { rButtonDown = &e })
		wnd.SetOnRButtonUpListener(func(e events.MouseClickEvent) { rButtonUp = &e })
		wnd.SetOnLButtonDoubleClickListener(func(e events.MouseClickEvent) { lButtonDblClick = &e })

		simulateAllMouseEvents(hwnd)
		win32.PostMessageW(hwnd, win32.WM_CLOSE, 0, 0)
	}, nil)

	if lButtonDown == nil {
		t.Fatal("LButtonDownListener was not called")
	}
	if lButtonUp == nil {
		t.Fatal("LButtonUpListener was not called")
	}
	if rButtonDown == nil {
		t.Fatal("RButtonDownListener was not called")
	}
	if rButtonUp == nil {
		t.Fatal("RButtonUpListener was not called")
	}
	if lButtonDblClick == nil {
		t.Fatal("LButtonDoubleClickListener was not called")
	}
}

// TestMouseListenerCoordinates verifies that mouse event coordinates and options are
// correctly propagated through the listener argument.
func TestMouseListenerCoordinates(t *testing.T) {
	if err := ensureTestClassRegistered(); err != nil {
		t.Fatal(err)
	}

	var wnd mouseListenerWindow
	var got *events.MouseClickEvent

	var loc win32util.MouseLocation
	loc.SetX(123)
	loc.SetY(456)
	var opts win32util.MouseClickOpt
	opts.SetShift(true)

	gw.Run(func(a *app.App) {
		hwnd, err := win32util.CreateWindow(&win32util.Wnd{
			ClassName: keyTestClassName,
			Style:     win32.WS_POPUP,
		})
		if err != nil {
			t.Fatal(err)
			return
		}
		if err := gw.Attach(hwnd, &wnd); err != nil {
			t.Fatal(err)
			return
		}

		wnd.SetOnDestroyListener(func() { a.Quit(0) })
		wnd.SetOnLButtonDownListener(func(e events.MouseClickEvent) { got = &e })

		win32util.SimulateLButtonDown(hwnd, opts, loc)
		win32.PostMessageW(hwnd, win32.WM_CLOSE, 0, 0)
	}, nil)

	if got == nil {
		t.Fatal("LButtonDownListener was not called")
	}
	if got.Pt.X() != 123 {
		t.Fatalf("Pt.X: got %d, want 123", got.Pt.X())
	}
	if got.Pt.Y() != 456 {
		t.Fatalf("Pt.Y: got %d, want 456", got.Pt.Y())
	}
	if !got.Opt.Shift() {
		t.Fatal("Opt.Shift: got false, want true")
	}
}

// TestMouseOnMethods verifies that overriding all On...Button...() methods receives events
// when the corresponding mouse messages are dispatched.
func TestMouseOnMethods(t *testing.T) {
	if err := ensureTestClassRegistered(); err != nil {
		t.Fatal(err)
	}

	var wnd mouseOverrideWindow

	gw.Run(func(a *app.App) {
		hwnd, err := win32util.CreateWindow(&win32util.Wnd{
			ClassName: keyTestClassName,
			Style:     win32.WS_POPUP,
		})
		if err != nil {
			t.Fatal(err)
			return
		}
		if err := gw.Attach(hwnd, &wnd); err != nil {
			t.Fatal(err)
			return
		}

		wnd.SetOnDestroyListener(func() { a.Quit(0) })

		simulateAllMouseEvents(hwnd)
		win32.PostMessageW(hwnd, win32.WM_CLOSE, 0, 0)
	}, nil)

	if wnd.onLButtonDown == nil {
		t.Fatal("OnLButtonDown was not called")
	}
	if wnd.onLButtonUp == nil {
		t.Fatal("OnLButtonUp was not called")
	}
	if wnd.onRButtonDown == nil {
		t.Fatal("OnRButtonDown was not called")
	}
	if wnd.onRButtonUp == nil {
		t.Fatal("OnRButtonUp was not called")
	}
	if wnd.onLButtonDblClick == nil {
		t.Fatal("OnLButtonDoubleClick was not called")
	}
}

// TestMouseListenerNil verifies that setting a listener to nil disables further callbacks.
func TestMouseListenerNil(t *testing.T) {
	if err := ensureTestClassRegistered(); err != nil {
		t.Fatal(err)
	}

	var wnd mouseListenerWindow
	var callCount int

	gw.Run(func(a *app.App) {
		hwnd, err := win32util.CreateWindow(&win32util.Wnd{
			ClassName: keyTestClassName,
			Style:     win32.WS_POPUP,
		})
		if err != nil {
			t.Fatal(err)
			return
		}
		if err := gw.Attach(hwnd, &wnd); err != nil {
			t.Fatal(err)
			return
		}

		wnd.SetOnDestroyListener(func() { a.Quit(0) })
		wnd.SetOnLButtonDownListener(func(e events.MouseClickEvent) {
			callCount++
			wnd.SetOnLButtonDownListener(nil) // Clear listener after first call.
		})

		var opt win32util.MouseClickOpt
		opt.SetLButton(true)
		win32util.SimulateLButtonDown(hwnd, opt, 0)
		win32util.SimulateLButtonUp(hwnd, opt, 0)
		win32.SendMessageW(hwnd, win32.WM_LBUTTONUP, 0, 0)
		win32util.SimulateLButtonClick(hwnd, 0, 0)

		win32.PostMessageW(hwnd, win32.WM_CLOSE, 0, 0)
	}, nil)

	if callCount != 1 {
		t.Fatalf("LButtonDown listener call count: got %d, want 1", callCount)
	}
}
