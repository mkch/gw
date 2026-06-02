package events

import (
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

type SizeEvent struct {
	Type win32.SizingType
	Size win32util.EventSize
}

func NewSizeEvent(wParam win32.WPARAM, lParam win32.LPARAM) SizeEvent {
	return SizeEvent{
		Type: win32.SizingType(wParam),
		Size: win32util.EventSize(lParam),
	}
}
