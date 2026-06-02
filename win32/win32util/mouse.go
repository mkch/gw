package win32util

import "github.com/mkch/gw/win32"

type MouseClickOpt win32.WPARAM

func (m MouseClickOpt) Control() bool {
	return m&win32.MK_CONTROL != 0
}

func (m *MouseClickOpt) set(opt MouseClickOpt, on bool) {
	if on {
		*m |= opt
	} else {
		*m &^= opt
	}
}

func (m *MouseClickOpt) SetControl(down bool) {
	m.set(win32.MK_CONTROL, down)
}

func (m MouseClickOpt) LButton() bool {
	return m&win32.MK_LBUTTON != 0
}

func (m *MouseClickOpt) SetLButton(down bool) {
	m.set(win32.MK_LBUTTON, down)
}

func (m MouseClickOpt) MButton() bool {
	return m&win32.MK_MBUTTON != 0
}

func (m *MouseClickOpt) SetMButton(down bool) {
	m.set(win32.MK_MBUTTON, down)
}

func (m MouseClickOpt) RButton() bool {
	return m&win32.MK_RBUTTON != 0
}

func (m *MouseClickOpt) SetRButton(down bool) {
	m.set(win32.MK_RBUTTON, down)
}

func (m MouseClickOpt) Shift() bool {
	return m&win32.MK_SHIFT != 0
}

func (m *MouseClickOpt) SetShift(down bool) {
	m.set(win32.MK_SHIFT, down)
}

func (m MouseClickOpt) XButton1() bool {
	return m&win32.MK_XBUTTON1 != 0
}

func (m *MouseClickOpt) SetXButton1(down bool) {
	m.set(win32.MK_XBUTTON1, down)
}

func (m MouseClickOpt) XButton2() bool {
	return m&win32.MK_XBUTTON2 != 0
}

func (m *MouseClickOpt) SetXButton2(down bool) {
	m.set(win32.MK_XBUTTON2, down)
}

// EventPoint is the point parameter of events.
type EventPoint win32.LPARAM

func (p EventPoint) X() int16 {
	return win32.GET_X_LPARAM(win32.LPARAM(p))
}

func (p *EventPoint) SetX(x int16) {
	*p = EventPoint(win32.MAKELONG(uint16(x), uint16(p.Y())))
}

func (p EventPoint) Y() int16 {
	return win32.GET_Y_LPARAM(win32.LPARAM(p))
}

func (p *EventPoint) SetY(y int16) {
	*p = EventPoint(win32.MAKELONG(uint16(p.X()), uint16(y)))
}

// MouseClickEvent is the event for mouse click messages.
type MouseClickEvent struct {
	Opt MouseClickOpt
	Pt  EventPoint
}
