package app

import (
	_ "unsafe" // for go:linkname

	"github.com/mkch/gw/app"
	"github.com/mkch/gw/internal/msghandler"
	"github.com/mkch/gw/win32"
)

// ThreadLocalApp returns the app instance associated with the current thread.
//
// If a nil-pointer returned here, there may be windows that were not properly destroyed
// by a previous app instance.
// These stale windows can receive messages when the new app instance calls
// win32.PeekMessageW before the TLS value is initialized in [BaseApp.Init].
//
// Ref: https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-peekmessagew
// As per MSDN, PeekMessageW "dispatches incoming nonqueued messages..."
//
// To avoid this issue in tests, always invoke app.DestroyAllWindows()
// before calling app.Quit().
//
//go:linkname ThreadLocalApp
func ThreadLocalApp() *app.BaseApp

//go:linkname AddMsgPreTranslator
func AddMsgPreTranslator(app *app.BaseApp, hwnd win32.HWND, translator func(msg *win32.MSG) bool)

//go:linkname RemoveMsgPreTranslator
func RemoveMsgPreTranslator(app *app.BaseApp, hwnd win32.HWND)

//go:linkname CallMsgRetListeners
func CallMsgRetListeners(app *app.BaseApp, hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT)

//go:linkname MsgHandlers
func MsgHandlers(app *app.BaseApp) map[win32.UINT]*msghandler.Chain
