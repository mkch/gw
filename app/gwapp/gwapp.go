// Package gwapp implements application initialization and message loop
// that can be used in any goroutine.
package gwapp

import (
	"math"
	"runtime"
	"sync"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
	"golang.org/x/sys/windows"
)

type GwApp struct {
	uiThreadId    win32.DWORD
	postMap       safeMap
	threadMsgHook win32.HHOOK
}

// New creates a GwApp and do application initialization.
func New() *GwApp {
	runtime.LockOSThread()

	app := &GwApp{
		uiThreadId: win32.DWORD(windows.GetCurrentThreadId()),
		postMap:    safeMap{ObjectMap: objectmap.New[func()](1, math.MaxUint)},
	}

	// Prepare postMap
	// See https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-postthreadmessagew#remarks

	// Initialize message queue
	win32.PeekMessageW(&win32.MSG{}, 0, 0, 0, win32.PM_NOREMOVE)
	// Install thread message hook
	proc := windows.NewCallback(func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
		if code >= 0 && win32.PeekMessageFlag(wParam) == win32.PM_REMOVE {
			if msg := (*win32.MSG)(unsafe.Add(nil, lParam)); msg.Message == appmsg.POST {
				// Handle posted functions
				app.postMap.Value(objectmap.Handle(msg.WParam))()
				msg.Message = win32.WM_NULL // Stop WNDPROC processing
				return 0                    // Stop other hook processing
			}
		}
		return win32.CallNextHookEx(app.threadMsgHook, code, wParam, lParam)
	})
	app.threadMsgHook = gg.Must(win32.SetWindowsHookExW(win32.WH_GETMESSAGE, proc, 0, app.uiThreadId))

	return app
}

// Run runs the message loop.
func (app *GwApp) Run() int {
	defer func() {
		gg.MustOK(win32.UnhookWindowsHookEx(app.threadMsgHook))
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

// Len returns the number of elements in the map.
// For debugging use only.
func (m *safeMap) Len() int {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.ObjectMap.Len()
}
