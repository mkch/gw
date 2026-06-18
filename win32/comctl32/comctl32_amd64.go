package comctl32

import (
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/sysutil"
)

func LBItemFromPt(hwnd win32.HWND, pt win32.POINT, autoScroll bool) int {
	var autoScrollInt = gg.If(autoScroll, 1, 0)
	uintptrPt := *(*uintptr)(unsafe.Pointer(&pt))
	return int(sysutil.As[win32.INT](lzLBItemFromPt.Call(uintptr(hwnd), uintptr(uintptrPt), uintptr(autoScrollInt))))
}
