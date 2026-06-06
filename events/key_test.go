package events_test

import (
	"testing"

	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

// keyEventWindow captures keyboard events for testing.
type keyEventWindow struct {
	gw.BaseWindowImpl
	keyDownEvent    *events.KeyEvent
	keyUpEvent      *events.KeyEvent
	sysKeyDownEvent *events.KeyEvent
	sysKeyUpEvent   *events.KeyEvent
}

func (w *keyEventWindow) OnKeyDown(event *events.KeyEvent) {
	w.keyDownEvent = event
}

func (w *keyEventWindow) OnKeyUp(event *events.KeyEvent) {
	w.keyUpEvent = event
}

func (w *keyEventWindow) OnSysKeyDown(event *events.KeyEvent) {
	w.sysKeyDownEvent = event
}

func (w *keyEventWindow) OnSysKeyUp(event *events.KeyEvent) {
	w.sysKeyUpEvent = event
}

// TestNewKeyEvent verifies that NewKeyEvent correctly populates VKCode and State from wParam/lParam.
func TestNewKeyEvent(t *testing.T) {
	var lp win32util.KeyMessageLParam
	lp.SetRepeatCount(3)
	lp.SetScanCode(0x1E)
	lp.SetExtended(true)
	lp.SetAltDown(false)
	lp.SetPreviousDown(true)
	lp.SetKeyUp(false)

	e := events.NewKeyEvent(win32.WPARAM('A'), win32.LPARAM(lp))

	if e.VKCode != win32.WPARAM('A') {
		t.Fatalf("VKCode: got 0x%X, want 0x%X ('A')", e.VKCode, win32.WPARAM('A'))
	}
	if e.State.RepeatCount() != 3 {
		t.Fatalf("RepeatCount: got %d, want 3", e.State.RepeatCount())
	}
	if e.State.ScanCode() != 0x1E {
		t.Fatalf("ScanCode: got 0x%02X, want 0x1E", e.State.ScanCode())
	}
	if !e.State.Extended() {
		t.Fatal("Extended: got false, want true")
	}
	if e.State.AltDown() {
		t.Fatal("AltDown: got true, want false")
	}
	if !e.State.PreviousDown() {
		t.Fatal("PreviousDown: got false, want true")
	}
	if e.State.KeyUp() {
		t.Fatal("KeyUp: got true, want false")
	}
}

// TestKeyEventDispatch verifies that key messages posted to a window are correctly
// dispatched to the overridden On...() methods.
func TestKeyEventDispatch(t *testing.T) {
	if err := ensureTestClassRegistered(); err != nil {
		t.Fatal(err)
	}

	var wnd keyEventWindow

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

		wnd.SetOnDestroyListener(func() {
			a.Quit(0)
		})

		// Post WM_KEYDOWN for key 'A'.
		var keyDownLParam win32util.KeyMessageLParam
		keyDownLParam.SetRepeatCount(1)
		keyDownLParam.SetScanCode(0x1E)
		win32.PostMessageW(hwnd, win32.WM_KEYDOWN, win32.WPARAM('A'), win32.LPARAM(keyDownLParam))

		// Post WM_KEYUP for key 'A'.
		var keyUpLParam win32util.KeyMessageLParam
		keyUpLParam.SetRepeatCount(1)
		keyUpLParam.SetScanCode(0x1E)
		keyUpLParam.SetKeyUp(true)
		win32.PostMessageW(hwnd, win32.WM_KEYUP, win32.WPARAM('A'), win32.LPARAM(keyUpLParam))

		// Post WM_SYSKEYDOWN for VK_MENU (Alt key).
		var sysKeyDownLParam win32util.KeyMessageLParam
		sysKeyDownLParam.SetRepeatCount(1)
		sysKeyDownLParam.SetScanCode(0x38)
		sysKeyDownLParam.SetAltDown(true)
		win32.PostMessageW(hwnd, win32.WM_SYSKEYDOWN, win32.WPARAM(win32.VK_MENU), win32.LPARAM(sysKeyDownLParam))

		// Post WM_SYSKEYUP for VK_MENU (Alt key).
		var sysKeyUpLParam win32util.KeyMessageLParam
		sysKeyUpLParam.SetRepeatCount(1)
		sysKeyUpLParam.SetScanCode(0x38)
		sysKeyUpLParam.SetAltDown(true)
		sysKeyUpLParam.SetKeyUp(true)
		win32.PostMessageW(hwnd, win32.WM_SYSKEYUP, win32.WPARAM(win32.VK_MENU), win32.LPARAM(sysKeyUpLParam))

		// Post WM_CLOSE to trigger window destruction and quit the message loop.
		win32.PostMessageW(hwnd, win32.WM_CLOSE, 0, 0)
	}, func(app *app.App) { app.DestroyAllWindows() })

	// Verify WM_KEYDOWN was dispatched to OnKeyDown with correct event data.
	if wnd.keyDownEvent == nil {
		t.Fatal("OnKeyDown was not called")
	}
	if wnd.keyDownEvent.VKCode != win32.WPARAM('A') {
		t.Fatalf("OnKeyDown VKCode: got 0x%X, want 0x41 ('A')", wnd.keyDownEvent.VKCode)
	}
	if wnd.keyDownEvent.State.ScanCode() != 0x1E {
		t.Fatalf("OnKeyDown ScanCode: got 0x%02X, want 0x1E", wnd.keyDownEvent.State.ScanCode())
	}
	if wnd.keyDownEvent.State.KeyUp() {
		t.Fatal("OnKeyDown State.KeyUp: got true, want false")
	}

	// Verify WM_KEYUP was dispatched to OnKeyUp with correct event data.
	if wnd.keyUpEvent == nil {
		t.Fatal("OnKeyUp was not called")
	}
	if wnd.keyUpEvent.VKCode != win32.WPARAM('A') {
		t.Fatalf("OnKeyUp VKCode: got 0x%X, want 0x41 ('A')", wnd.keyUpEvent.VKCode)
	}
	if !wnd.keyUpEvent.State.KeyUp() {
		t.Fatal("OnKeyUp State.KeyUp: got false, want true")
	}

	// Verify WM_SYSKEYDOWN was dispatched to OnSysKeyDown with correct event data.
	if wnd.sysKeyDownEvent == nil {
		t.Fatal("OnSysKeyDown was not called")
	}
	if wnd.sysKeyDownEvent.VKCode != win32.WPARAM(win32.VK_MENU) {
		t.Fatalf("OnSysKeyDown VKCode: got 0x%X, want 0x%X (VK_MENU)", wnd.sysKeyDownEvent.VKCode, win32.WPARAM(win32.VK_MENU))
	}
	if !wnd.sysKeyDownEvent.State.AltDown() {
		t.Fatal("OnSysKeyDown State.AltDown: got false, want true")
	}
	if wnd.sysKeyDownEvent.State.KeyUp() {
		t.Fatal("OnSysKeyDown State.KeyUp: got true, want false")
	}

	// Verify WM_SYSKEYUP was dispatched to OnSysKeyUp with correct event data.
	if wnd.sysKeyUpEvent == nil {
		t.Fatal("OnSysKeyUp was not called")
	}
	if wnd.sysKeyUpEvent.VKCode != win32.WPARAM(win32.VK_MENU) {
		t.Fatalf("OnSysKeyUp VKCode: got 0x%X, want 0x%X (VK_MENU)", wnd.sysKeyUpEvent.VKCode, win32.WPARAM(win32.VK_MENU))
	}
	if !wnd.sysKeyUpEvent.State.KeyUp() {
		t.Fatal("OnSysKeyUp State.KeyUp: got false, want true")
	}
}
