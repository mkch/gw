package comctl32

import (
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/sysutil"
	"golang.org/x/sys/windows"
)

var lzComctl32 = windows.NewLazySystemDLL("comctl32.dll")

var lzDrawInsert = lzComctl32.NewProc("DrawInsert")

func DrawInsert(parent, listbox win32.HWND, itemIndex int) {
	lzDrawInsert.Call(uintptr(parent), uintptr(listbox), uintptr(itemIndex))
}

var lzMakeDragList = lzComctl32.NewProc("MakeDragList")

func MakeDragList(hwnd win32.HWND) error {
	return sysutil.MustTrue(lzMakeDragList.Call(uintptr(hwnd)))
}

var lzLBItemFromPt = lzComctl32.NewProc("LBItemFromPt")
