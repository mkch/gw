// Package gwapp implements application initialization and message loop
// that can be used in any goroutine.
package gwapp

import (
	"runtime"

	"github.com/mkch/gw/win32"
)

type app BareApp // Avoid exposing that a [GwApp] is actually an [BaseApp].

// GwApp is the main application struct of gw.
// Any other gw functionalities must be used in the same goroutine as the one that creates the GwApp.
// Usually functions in the [app] package are sufficient to run a gw app instead of using this struct directly.
// However, if there is a need to run gw in a goroutine other than main, a GwApp must be created by calling [New]
// before operating gw there.
type GwApp app

// New creates a GwApp and do application initialization.
func New() *GwApp {
	runtime.LockOSThread()
	return (*GwApp)(newBase(false))
}

// Run runs the message loop.
func (app *GwApp) Run() int {
	defer func() {
		(*BareApp)(app).destroy()
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

		if (*BareApp)(app).preTranslateMessage(&msg) {
			continue
		}
		win32.TranslateMessage(&msg)
		win32.DispatchMessageW(&msg)
	}
}

// AddMessageRetListener adds a listener that is called after a message is processed in any window procedure.
// The returned key can be used to remove the listener by calling [RemoveMessageRetListener].
func (app *GwApp) AddMessageRetListener(listener MessageRetListener) (key MessageRetListenerKey) {
	key = MessageRetListenerKey{p: &listener}
	app.msgRetListeners[key] = listener
	return
}

// RemoveMessageRetListener removes the listener added by AddMessageRetListener.
func (app *GwApp) RemoveMessageRetListener(key MessageRetListenerKey) {
	delete(app.msgRetListeners, key)
}

// Post put f into the UI message queue, f will run in the UI thread ASAP.
func (app *GwApp) Post(f func()) error {
	return (*BareApp)(app).Post(f)
}

// Quit calls win32.PostQuitMessage which tells the message loop to exit.
// The exit code will be the return value of Run.
func (app *GwApp) Quit(exitCode int) {
	(*BareApp)(app).Quit(exitCode)
}
