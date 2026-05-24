// Package gwapp implements application initialization and message loop
// that can be used in any goroutine.
package gwapp

import (
	"runtime"

	"github.com/mkch/gw/win32"
)

type app ExternalApp // Avoid exposing that a [GwApp] is actually an [ExternalApp].

type GwApp app

// New creates a GwApp and do application initialization.
func New() *GwApp {
	runtime.LockOSThread()
	return (*GwApp)(NewExternal(false))
}

// Run runs the message loop.
func (app *GwApp) Run() int {
	defer func() {
		(*ExternalApp)(app).Destroy()
		runtime.UnlockOSThread()
	}()
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

		if (*ExternalApp)(app).PreTranslateMessage(&msg) {
			continue
		}
		win32.TranslateMessage(&msg)
		win32.DispatchMessageW(&msg)
	}
}

// Post put f into the UI message queue, f will run in the UI thread ASAP.
func (app *GwApp) Post(f func()) error {
	return (*ExternalApp)(app).Post(f)
}

// Quit calls win32.PostQuitMessage which tells the message loop to exit.
// The exit code will be the return value of Run.
func (app *GwApp) Quit(exitCode int) {
	(*ExternalApp)(app).Quit(exitCode)
}
