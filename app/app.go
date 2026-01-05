// Package app implements application initialization and message loop that
// can be used in main goroutine only.
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

type MsgProc = gwapp.MessageDispatcher

// SetMessageDispatcher sets a dispatcher for windows message dispatching.
// The default message dispatcher is [win32.DispatchMessageW].
func SetMessageDispatcher(msgProc MsgProc) {
	app.SetMessageDispatcher(msgProc)
}
