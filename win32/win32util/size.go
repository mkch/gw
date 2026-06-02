package win32util

import "github.com/mkch/gw/win32"

// EventSize is the size parameter of event.
type EventSize win32.LPARAM

func (p EventSize) Width() int16 {
	return win32.GET_X_LPARAM(win32.LPARAM(p))
}

func (p *EventSize) SetWidth(x int16) {
	*p = EventSize(win32.MAKELONG(uint16(x), uint16(p.Height())))
}

func (p EventSize) Height() int16 {
	return win32.GET_Y_LPARAM(win32.LPARAM(p))
}

func (p *EventSize) SetHeight(y int16) {
	*p = EventSize(win32.MAKELONG(uint16(p.Width()), uint16(y)))
}
