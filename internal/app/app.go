package app

import (
	_ "unsafe" // for go:linkname

	"github.com/mkch/gw/app"
	"github.com/mkch/gw/win32"
)

//go:linkname ThreadLocalApp
func ThreadLocalApp() *app.BaseApp

//go:linkname AddMsgPreTranslator
func AddMsgPreTranslator(app *app.BaseApp, hwnd win32.HWND, translator func(msg *win32.MSG) bool)

//go:linkname RemoveMsgPreTranslator
func RemoveMsgPreTranslator(app *app.BaseApp, hwnd win32.HWND)

//go:linkname CallMsgRetListeners
func CallMsgRetListeners(app *app.BaseApp, hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT)
