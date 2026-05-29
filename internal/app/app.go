package app

import (
	"unsafe"
	_ "unsafe" // for go:linkname

	"github.com/mkch/gw/app"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/win32"
)

//go:linkname ThreadLocalApp
func ThreadLocalApp() *app.BareApp

//go:linkname AddMsgPreTranslator
func AddMsgPreTranslator(app *app.BareApp, hwnd win32.HWND, translator func(msg *win32.MSG) bool)

//go:linkname RemoveMsgPreTranslator
func RemoveMsgPreTranslator(app *app.BareApp, hwnd win32.HWND)

//go:linkname MenuMap
func MenuMap(app *app.BareApp) map[win32.HMENU]unsafe.Pointer

//go:linkname MenuItemMap
func MenuItemMap(app *app.BareApp) *objectmap.ObjectMap[unsafe.Pointer]

//go:linkname CallMsgRetListeners
func CallMsgRetListeners(app *app.BareApp, hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT)
