package app

import "github.com/mkch/gw/win32"

// https://learn.microsoft.com/en-us/windows/win32/winmsg/wm-app
// Message numbers in the third range (0x8000 through 0xBFFF) are available for applications to use as private messages.
// Messages in this range do not conflict with system messages.

const (
	msgPost = win32.WM_APP + iota
	MSG_NOTIFY_ICON_CALLBACK
	MSG_REFLECT_COMMAND
	MSG_REFLECT_CTLCOLORSTATIC
	MSG_REFLECT_COMPAREITEM
	MSG_REFLECT_MEASUREITEM
	MSG_REFLECT_DRAWITEM

	MSG_LAST_REFLECT
)
