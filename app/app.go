// Package app initializes application and message loop in main goroutine.
// Running gw in main goroutine involves importing this package and calling [Run] in main function.
// This package creates a [GwApp] in initialization and the functions in this package operate on that [GwApp].
// See [gwapp] for details of application creation in other scenarios.
package app

import "github.com/mkch/gw/app/gwapp"

var app *gwapp.GwApp = gwapp.New()

// Run runs the message loop.
func Run() int {
	return app.Run()
}

// Post put f into the UI message queue, f will run in the UI thread ASAP.
func Post(f func()) error {
	return app.Post(f)
}

// Quit calls win32.PostQuitMessage which tells the message loop to exit.
// The exit code will be the return value of Run.
func Quit(exitCode int) {
	app.Quit(exitCode)
}

// AddMessageRetListener adds a listener that is called after a message is processed in any window procedure.
// The returned key can be used to remove the listener by calling [RemoveMessageRetListener].
func AddMessageRetListener(listener gwapp.MessageRetListener) (key gwapp.MessageRetListenerKey) {
	return app.AddMessageRetListener(listener)
}

// RemoveMessageRetListener removes the listener added by AddMessageRetListener.
func RemoveMessageRetListener(key gwapp.MessageRetListenerKey) {
	app.RemoveMessageRetListener(key)
}
