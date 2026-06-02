package win32util

import (
	"fmt"

	"github.com/mkch/gw/win32"
)

// KeyMessageLParam represents the lParam of keyboard messages (WM_KEYDOWN, WM_KEYUP, WM_SYSKEYDOWN, WM_SYSKEYUP).
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#keystroke-message-flags
type KeyMessageLParam win32.LPARAM

func (l KeyMessageLParam) String() string {
	return fmt.Sprintf("RepeatCount: %d, ScanCode: 0x%02X, Extended: %t, DialogMode: %t, MenuMode: %t, AltDown: %t, PreviousDown: %t, KeyUp: %t",
		l.RepeatCount(), l.ScanCode(), l.Extended(), l.DialogMode(), l.MenuMode(), l.AltDown(), l.PreviousDown(), l.KeyUp())
}

// RepeatCount returns the repeat count of the key message,
// which is the number of times the keystroke message should be sent while the user holds down the key.
func (l KeyMessageLParam) RepeatCount() uint16 {
	return uint16(l & 0xFFFF)
}

// SetRepeatCount sets the repeat count of the key message.
func (l *KeyMessageLParam) SetRepeatCount(count uint16) {
	*l = KeyMessageLParam((*l &^ 0xFFFF) | KeyMessageLParam(count))
}

// ScanCode returns the scan code of the key message.
// https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#scan-codes
func (l KeyMessageLParam) ScanCode() uint16 {
	return uint16((l >> 16) & 0xFF)
}

// SetScanCode sets the scan code of the key message.
func (l *KeyMessageLParam) SetScanCode(code uint16) {
	*l = KeyMessageLParam((*l &^ (0xFF << 16)) | (KeyMessageLParam(code) << 16))
}

// Extended returns the `Extended Key` flag in l.
// The `Extended Key` flag indicates whether the key is an extended key.
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#extended-key-flag
func (l KeyMessageLParam) Extended() bool {
	return l.bit(24)
}

// SetExtended sets the `Extended Key` flag in l.
func (l *KeyMessageLParam) SetExtended(extended bool) {
	l.setBit(24, extended)
}

// DialogMode returns whether a dialog box is active when the keystroke message was generated.
func (l KeyMessageLParam) DialogMode() bool {
	return l.bit(27)
}

// SetDialogMode sets whether a dialog box is active when the keystroke message was generated.
func (l *KeyMessageLParam) SetDialogMode(dialogMode bool) {
	l.setBit(27, dialogMode)
}

// MenuMode returns whether a menu is active when the keystroke message was generated.
func (l KeyMessageLParam) MenuMode() bool {
	return l.bit(28)
}

// SetMenuMode sets whether a menu is active when the keystroke message was generated.
func (l *KeyMessageLParam) SetMenuMode(menuMode bool) {
	l.setBit(28, menuMode)
}

// AltDown returns the `Context Code` flag in l.
// The `Context Code` flag indicates whether the Alt key was down when the keystroke message was generated.
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#context-code
func (l KeyMessageLParam) AltDown() bool {
	return l.bit(29)
}

// SetAltDown sets the `Context Code` flag in l.
func (l *KeyMessageLParam) SetAltDown(altDown bool) {
	l.setBit(29, altDown)
}

// PreviousDown returns the `Previous Key-State Flag` flag in l.
// The `Previous Key-State Flag` flag indicates whether the key was down before the message was sent.
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#previous-key-state-flag
func (l KeyMessageLParam) PreviousDown() bool {
	return l.bit(30)
}

// SetPreviousDown sets the `Previous Key-State Flag` flag in l.
func (l *KeyMessageLParam) SetPreviousDown(previousDown bool) {
	l.setBit(30, previousDown)
}

// KeyUp returns the `Transition-State Flag` flag in l.
// The `Transition-State Flag` flag indicates whether pressing a key generated the keystroke message.
// True for WM_KEYUP and WM_SYSKEYUP messages; false for WM_KEYDOWN and WM_SYSKEYDOWN messages;
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#transition-state-flag
func (l KeyMessageLParam) KeyUp() bool {
	return l.bit(31)
}

// SetKeyUp sets the 'KeyUp' flag in l.
func (l *KeyMessageLParam) SetKeyUp(keyUp bool) {
	l.setBit(31, keyUp)
}

// bit is a helper method to check if the i-th bit is set in l.
func (l KeyMessageLParam) bit(i int) bool {
	return (l & (1 << i)) != 0
}

// setBit is a helper method to set the i-th bit in l.
func (l *KeyMessageLParam) setBit(i int, value bool) {
	if value {
		*l |= KeyMessageLParam(1 << i)
	} else {
		*l &^= KeyMessageLParam(1 << i)
	}
}
