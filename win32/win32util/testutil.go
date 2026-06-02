package win32util

import (
	"github.com/mkch/gw/win32"
)

// SimulateControlCommand simulates a control command by sending a WM_COMMAND message to the parent window of the control.
// wParam is the WPARAM of the WM_COMMAND message, which typically contains the control ID and notification code.
func SimulateControlCommand(parent, control win32.HWND, wParam win32.WPARAM) {
	win32.SendMessageW(parent, win32.WM_COMMAND, wParam, win32.LPARAM(control))
}

// SimulateLButtonDown simulates a mouse left button down event at the specified location with the given options.
func SimulateLButtonDown(hwnd win32.HWND, opts MouseClickOpt, loc EventPoint) {
	win32.PostMessageW(hwnd, win32.WM_LBUTTONDOWN, win32.WPARAM(opts), win32.LPARAM(loc))
}

// SimulateLButtonUp simulates a mouse left button up event at the specified location with the given options.
func SimulateLButtonUp(hwnd win32.HWND, opts MouseClickOpt, loc EventPoint) {
	win32.PostMessageW(hwnd, win32.WM_LBUTTONUP, win32.WPARAM(opts), win32.LPARAM(loc))
}

// SimulateLButtonClick simulates a mouse left button click at the specified location with the given options.
func SimulateLButtonClick(hwnd win32.HWND, opts MouseClickOpt, loc EventPoint) {
	SimulateLButtonDown(hwnd, opts, loc)
	SimulateLButtonUp(hwnd, opts, loc)
}

// SimulateLButtonDoubleClick simulates a mouse left button double click at the specified location with the given options.
func SimulateLButtonDoubleClick(hwnd win32.HWND, opts MouseClickOpt, loc EventPoint) {
	// https://learn.microsoft.com/en-us/windows/win32/inputdev/wm-lbuttondblclk
	// Double-clicking the left mouse button actually generates a sequence of four messages:
	// WM_LBUTTONDOWN, WM_LBUTTONUP, WM_LBUTTONDBLCLK, and WM_LBUTTONUP
	SimulateLButtonClick(hwnd, opts, loc)
	win32.PostMessageW(hwnd, win32.WM_LBUTTONDBLCLK, win32.WPARAM(opts), win32.LPARAM(loc))
	win32.PostMessageW(hwnd, win32.WM_LBUTTONUP, win32.WPARAM(opts), win32.LPARAM(loc))
}

// SimulateRButtonDown simulates a mouse right button down event at the specified location with the given options.
func SimulateRButtonDown(hwnd win32.HWND, opts MouseClickOpt, loc EventPoint) {
	win32.PostMessageW(hwnd, win32.WM_RBUTTONDOWN, win32.WPARAM(opts), win32.LPARAM(loc))
}

// SimulateRButtonUp simulates a mouse right button up event at the specified location with the given options.
func SimulateRButtonUp(hwnd win32.HWND, opts MouseClickOpt, loc EventPoint) {
	win32.PostMessageW(hwnd, win32.WM_RBUTTONUP, win32.WPARAM(opts), win32.LPARAM(loc))
}
