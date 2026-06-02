package events

import (
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

// KeyEvent is the event of keyboard messages (WM_KEYDOWN, WM_KEYUP, WM_SYSKEYDOWN, WM_SYSKEYUP).
type KeyEvent struct {
	VKCode win32.WPARAM
	State  win32util.KeyMessageLParam
}

func NewKeyEvent(wParam win32.WPARAM, lParam win32.LPARAM) KeyEvent {
	return KeyEvent{
		VKCode: wParam,
		State:  win32util.KeyMessageLParam(lParam),
	}
}
