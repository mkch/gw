package events

import (
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

// MouseClickEvent is the event for mouse click messages.
type MouseClickEvent struct {
	Opt win32util.MouseClickOpt
	Pt  win32util.MouseLocation
}

func NewMouseClickEvent(wParam win32.WPARAM, lParam win32.LPARAM) MouseClickEvent {
	return MouseClickEvent{
		Opt: win32util.MouseClickOpt(wParam),
		Pt:  win32util.MouseLocation(lParam),
	}
}
