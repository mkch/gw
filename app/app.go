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
type App BareApp

// New creates a App and do application initialization.
// The returned app should be destroyed by calling [App.Destroy] after use.
// No gw operation should be performed after the app is destroyed.
func New() *App {
	runtime.LockOSThread()
	return (*App)(newBare(false))
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

		if (*BareApp)(app).preTranslateMessage(&msg) {
			continue
		}
		win32.TranslateMessage(&msg)
		win32.DispatchMessageW(&msg)
	}
}

// AddMessageRetListener adds a listener that is called after a message is processed in any window procedure.
// The returned key can be used to remove the listener by calling [RemoveMessageRetListener].
func (app *App) AddMessageRetListener(listener MessageRetListener) (key MessageRetListenerKey) {
	return (*BareApp)(app).AddMessageRetListener(listener)
}

// RemoveMessageRetListener removes the listener added by AddMessageRetListener.
func (app *App) RemoveMessageRetListener(key MessageRetListenerKey) {
	(*BareApp)(app).RemoveMessageRetListener(key)
}

// Post put f into the UI message queue, f will run in the UI thread ASAP.
func (app *App) Post(f func()) error {
	return (*BareApp)(app).Post(f)
}

// Quit calls win32.PostQuitMessage which tells the message loop to exit.
// The exit code will be the return value of [App.Run].
func (app *App) Quit(exitCode int) {
	(*BareApp)(app).Quit(exitCode)
}

// Destroy destroys the app for an external message loop.
// This function must be called in the external message loop before exiting.
// After calling this function, the app should not be used anymore and no gw
// operation should be performed after that.
func (app *App) Destroy() {
	(*BareApp)(app).Destroy()
	runtime.UnlockOSThread()
}
