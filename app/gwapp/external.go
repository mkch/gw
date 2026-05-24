package gwapp

import (
	"math"
	"sync"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
	"golang.org/x/sys/windows"
)

type ExternalApp struct {
	uiThreadId    win32.DWORD
	postMap       safeMap
	threadMsgHook win32.HHOOK
}

// NewExternal creates a [ExternalApp] for an external message loop.
// The external initialization code must call NewExternal in the main thread that runs the message loop.
// If hookMsgLoop is true, a thread message hook will be installed to process messages with app.PreTranslateMessage
// and call app.Destroy when WM_QUIT is received.
// If hookMsgLoop is false, it is the responsibility of the external code to call app.PreTranslateMessage
// in the message loop and call app.Destroy before exiting.
func NewExternal(hookMsgLoop bool) (app *ExternalApp) {
	app = &ExternalApp{
		uiThreadId: win32.DWORD(windows.GetCurrentThreadId()),
		postMap:    safeMap{ObjectMap: objectmap.New[func()](1, math.MaxUint)},
	}

	// Prepare postMap
	// See https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-postthreadmessagew#remarks

	// Initialize message queue
	win32.PeekMessageW(&win32.MSG{}, 0, 0, 0, win32.PM_NOREMOVE)
	// Install thread message hook
	var proc func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT
	if hookMsgLoop {
		proc = func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
			if code >= 0 && win32.PeekMessageFlag(wParam) == win32.PM_REMOVE {
				msg := (*win32.MSG)(unsafe.Add(nil, lParam))
				switch msg.Message {
				case win32.WM_QUIT:
					app.Destroy() // Destroy the app when the message loop is about to exit.
				case appmsg.POST:
					// Handle posted functions
					app.postMap.Value(objectmap.Handle(msg.WParam))()
					msg.Message = win32.WM_NULL // Stop WNDPROC processing
					return 0                    // Stop other hook processing
				default:
					if app.PreTranslateMessage(msg) {
						msg.Message = win32.WM_NULL // Stop WNDPROC processing
						return 0                    // Stop other hook processing
					}
				}
			}
			return win32.CallNextHookEx(app.threadMsgHook, code, wParam, lParam)
		}
	} else {
		proc = func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
			if code >= 0 && win32.PeekMessageFlag(wParam) == win32.PM_REMOVE {
				if msg := (*win32.MSG)(unsafe.Add(nil, lParam)); msg.Message == appmsg.POST {
					// Handle posted functions
					app.postMap.Value(objectmap.Handle(msg.WParam))()
					msg.Message = win32.WM_NULL // Stop WNDPROC processing
					return 0                    // Stop other hook processing
				}
			}
			return win32.CallNextHookEx(app.threadMsgHook, code, wParam, lParam)
		}
	}
	app.threadMsgHook = gg.Must(win32.SetWindowsHookExW(win32.WH_GETMESSAGE, windows.NewCallback(proc), 0, app.uiThreadId))

	return app
}

// Destroy destroys the app for an external message loop.
// This function must be called in the external message loop before exiting.
// After calling this function, the app should not be used anymore.
func (app *ExternalApp) Destroy() {
	gg.MustOK(win32.UnhookWindowsHookEx(app.threadMsgHook))
}

// PreTranslateMessage should be called in the external message loop before TranslateMessage.
// If it returns true, the message must not be passed to TranslateMessage and DispatchMessage.
func (app *ExternalApp) PreTranslateMessage(msg *win32.MSG) bool {
	return window.PreTranslateMessage(msg)
}

// Post put f into the UI message queue, f will run in the UI thread ASAP.
func (app *ExternalApp) Post(f func()) error {
	var h objectmap.Handle
	h = app.postMap.Add(func() {
		f()
		app.postMap.Remove(h)
	})
	return win32.PostThreadMessageW(app.uiThreadId, appmsg.POST, win32.WPARAM(h), 0)
}

// Quit calls win32.PostQuitMessage which tells the message loop to exit.
// The exit code will be the return value of Run.
func (app *ExternalApp) Quit(exitCode int) {
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
