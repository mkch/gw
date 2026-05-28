package app

import (
	"unsafe"
	_ "unsafe" // for go:linkname

	"github.com/mkch/gw/app/gwapp"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/win32"
)

//go:linkname ThreadLocalApp
func ThreadLocalApp() *gwapp.BareApp

//go:linkname AddMsgPreTranslator
func AddMsgPreTranslator(app *gwapp.BareApp, hwnd win32.HWND, translator func(msg *win32.MSG) bool)

//go:linkname RemoveMsgPreTranslator
func RemoveMsgPreTranslator(app *gwapp.BareApp, hwnd win32.HWND)

//go:linkname MenuMap
func MenuMap(app *gwapp.BareApp) map[win32.HMENU]unsafe.Pointer

//go:linkname MenuItemMap
func MenuItemMap(app *gwapp.BareApp) *objectmap.ObjectMap[unsafe.Pointer]
