// Package gwapp implements application initialization and message loop
// that can be used in any goroutine.
package gwapp

import (
	"math"
	"runtime"
	"sync"

	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
	"golang.org/x/sys/windows"
)

type GwApp struct {
	uiThreadId win32.DWORD
	postMap    safeMap
}

// New creates a GwApp and do application initialization.
func New() *GwApp {
	runtime.LockOSThread()
	return &GwApp{
		uiThreadId: win32.DWORD(windows.GetCurrentThreadId()),
		postMap:    safeMap{ObjectMap: objectmap.New[func()](1, math.MaxUint)},
	}
}

// Run runs the message loop.
func (app *GwApp) Run() int {
	defer runtime.UnlockOSThread()
	var msg win32.MSG
	for {
		r := win32.GetMessageW(&msg, 0, 0, 0)
		if r == -1 {
			panic(r)
		}
		if r == 0 {
			return int(msg.WParam)
		}
		if msg.Hwnd == 0 && msg.Message == appmsg.POST {
			app.postMap.Value(objectmap.Handle(msg.WParam))()
		}

		if !window.PreTranslateMessage(&msg) {
			win32.TranslateMessage(&msg)
		}
		win32.DispatchMessageW(&msg)
	}
}

// Post put f into the UI message queue, f will run in the UI thread ASAP.
func (app *GwApp) Post(f func()) error {
	var h objectmap.Handle
	h = app.postMap.Add(func() {
		f()
		app.postMap.Remove(h)
	})
	return win32.PostThreadMessageW(app.uiThreadId, appmsg.POST, win32.WPARAM(h), 0)
}

// Quit calls win32.PostQuitMessage which tells the message loop to exit.
// The exit code will be the return value of Run.
func (app *GwApp) Quit(exitCode int) {
	win32.PostQuitMessage(exitCode)
}

type safeMap struct {
	*objectmap.ObjectMap[func()]
	l sync.RWMutex
}

func (m *safeMap) Add(f func()) objectmap.Handle {
	m.l.Lock()
	defer m.l.Unlock()
	return m.ObjectMap.Add(f)
}

func (m *safeMap) Value(h objectmap.Handle) func() {
	m.l.RLock()
	defer m.l.RUnlock()
	f, _ := m.ObjectMap.Value(h)
	return f
}

func (m *safeMap) Remove(h objectmap.Handle) {
	m.l.Lock()
	defer m.l.Unlock()
	m.ObjectMap.Remove(h)
}
