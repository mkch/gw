package comctl32

import (
	"github.com/mkch/gg"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/sysutil"
)

func LBItemFromPt(hwnd win32.HWND, pt win32.POINT, autoScroll bool) int {
	var autoScrollInt = gg.If(autoScroll, 1, 0)
	return sysutil.As[int](lzLBItemFromPt.Call(uintptr(hwnd), uintptr(pt.X), uintptr(pt.Y), uintptr(autoScrollInt)))
}
