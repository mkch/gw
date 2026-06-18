package appmsg

import "github.com/mkch/gw/win32"

// https://learn.microsoft.com/en-us/windows/win32/winmsg/wm-app
// Message numbers in the third range (0x8000 through 0xBFFF) are available for applications to use as private messages.
// Messages in this range do not conflict with system messages.

const (
	POST = win32.WM_APP + iota
	NOTIFY_ICON_CALLBACK
	REFLECT_COMMAND
	REFLECT_CTLCOLORSTATIC
	REFLECT_COMPAREITEM
	REFLECT_MEASUREITEM
	REFLECT_DRAWITEM

	LAST_REFLECT_MESSAGE
)
