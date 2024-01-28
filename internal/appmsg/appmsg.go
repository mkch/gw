package appmsg

import "github.com/mkch/gw/win32"

const (
	REFLECT_COMMAND = win32.WM_APP + 0xFFFF + iota
	POST
)
