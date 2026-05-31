// Package app implements application initialization and message loop
// that can be used in any goroutine.
package app

import (
	"runtime"

	"github.com/mkch/gw/win32"
)

// App is the main application struct of gw.
// Any other gw functionalities must be used in the same goroutine as the one that creates the App.
// Usually gw.Run is sufficient to run a gw app instead of using this type directly.
// However, if there is a need to run gw in a goroutine other than main, a App must be created by calling [New]
// before operating gw there.
type App struct {
	BaseApp
}

// New creates a App and do application initialization.
// The returned app should be destroyed by calling [App.Destroy] after use.
// No gw operation should be performed after the app is destroyed.
func New() *App {
	runtime.LockOSThread()
	var app App
	app.init(false)
	return &app
}

// Run runs the message loop.
// This function does not return until the message loop exits.
// To exit the message loop, call [App.Quit].
func (app *App) Run() int {
	var msg win32.MSG
	for {
		r := win32.GetMessageW(&msg, 0, 0, 0)
		if r == -1 {
			panic(r)
		}
		if r == 0 {
			return int(msg.WParam)
		}
		if msg.Hwnd == 0 {
			continue // Messages not associated with a window cannot be dispatched
		}

		if app.preTranslateMessage(&msg) {
			continue
		}
		win32.TranslateMessage(&msg)
		win32.DispatchMessageW(&msg)
	}
}

// Destroy destroys the app for an external message loop.
// This function must be called in the external message loop before exiting.
// After calling this function, the app should not be used anymore and no gw
// operation should be performed after that.
func (app *App) Destroy() {
	app.BaseApp.Destroy()
	runtime.UnlockOSThread()
}
