package events

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

// MouseClickEvent is the event for mouse click messages, including WM_LBUTTONDOWN, WM_LBUTTONUP, WM_RBUTTONDOWN, and WM_RBUTTONUP.
type MouseClickEvent struct {
	Opt MouseClickOpt
	X   int
	Y   int
}
