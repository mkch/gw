// package msghandler provides types shared by app and gw packages for message handling.
package msghandler

import (
	"github.com/mkch/go-events/events"
	"github.com/mkch/gw/win32"
)

// Arg is the argument type of [msghandler.Chain] handlers.
type Arg struct {
	Hwnd    win32.HWND
	Message win32.UINT
	WParam  win32.WPARAM
	LParam  win32.LPARAM
}

// HandlerKey is the key type for message handlers in a [Chain].
type HandlerKey = events.HandlerKey[*Arg, win32.LRESULT]

// Chain is a [events.Chain] of message handlers.
type Chain = events.Chain[*Arg, win32.LRESULT]
