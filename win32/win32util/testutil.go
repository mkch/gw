package win32util

import "github.com/mkch/gw/win32"

// SimulateControlCommand simulates a control command by sending a WM_COMMAND message to the parent window of the control.
// wParam is the WPARAM of the WM_COMMAND message, which typically contains the control ID and notification code.
func SimulateControlCommand(parent, control win32.HWND, wParam win32.WPARAM) {
	win32.SendMessageW(parent, win32.WM_COMMAND, wParam, win32.LPARAM(control))
}
