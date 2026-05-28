package gwapp

import (
	"math"
	"sync"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/internal"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/win32"
	"golang.org/x/sys/windows"
)

// tlsApp is the TLS index for storing *BaseApp in the UI thread.
var tlsApp = gg.Must(win32.TlsAlloc())

//go:linkname threadLocalApp github.com/mkch/gw/internal/app.ThreadLocalApp
func threadLocalApp() *BareApp {
	return (*BareApp)(win32.PVOID(unsafe.Pointer(gg.Must(win32.TlsGetValue(tlsApp)))))
}

//go:linkname app_AddMsgPreTranslator github.com/mkch/gw/internal/app.AddMsgPreTranslator
func app_AddMsgPreTranslator(app *BareApp, hwnd win32.HWND, translator func(msg *win32.MSG) bool) {
	app.msgPreTranslators[hwnd] = translator
}

//go:linkname app_RemoveMsgPreTranslator github.com/mkch/gw/internal/app.RemoveMsgPreTranslator
func app_RemoveMsgPreTranslator(app *BareApp, hwnd win32.HWND) {
	delete(app.msgPreTranslators, hwnd)
}

//go:linkname app_MenuMap github.com/mkch/gw/internal/app.MenuMap
func app_MenuMap(app *BareApp) map[win32.HMENU]unsafe.Pointer {
	return app.menuMap
}

//go:linkname app_MenuItemMap github.com/mkch/gw/internal/app.MenuItemMap
func app_MenuItemMap(app *BareApp) *objectmap.ObjectMap[unsafe.Pointer] {
	return app.menuItemMap
}

//go:linkname app_callMsgRetListeners github.com/mkch/gw/internal/app.CallMsgRetListeners
func app_callMsgRetListeners(app *BareApp, hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT) {
	app.callMsgRetListeners(hwnd, message, wParam, lParam, result)
}

// MessageRetListener is a function type used by [BareApp.AddMessageRetListener] and [GwApp.AddMessageRetListener].
type MessageRetListener func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT)

// MessageRetListenerKey is the key for a listener added by [BareApp.AddMessageRetListener]  and [GwApp.AddMessageRetListener].
// It can be used to remove the listener by calling [BareApp.RemoveMessageRetListener] and [GwApp.RemoveMessageRetListener].
type MessageRetListenerKey struct{ p *MessageRetListener }

// BareApp is the application for external message loop.
// When integrating with external code that already has a message loop, call [NewBare] to create a
// [BareApp] and then operate gw in that goroutine.
// Any other gw functionalities must be used in the same goroutine as the one that creates the BareApp.
// See [NewBare] for details.
type BareApp struct {
	uiThreadId        win32.DWORD
	postMap           safeMap
	getMsgHook        win32.HHOOK
	msgPreTranslators map[win32.HWND]func(msg *win32.MSG) bool
	msgRetListeners   map[MessageRetListenerKey]MessageRetListener
	menuMap           map[win32.HMENU]unsafe.Pointer       // unsafe.Pointer is [*github.com/mkch/gw/menu.Menu]
	menuItemMap       *objectmap.ObjectMap[unsafe.Pointer] // unsafe.Pointer is [*github.com/mkch/gw/menu.MenuItem]
	cleanup           func()                               // called before the app is destroyed, can be nil.
}

// NewBare creates a [BareApp] that do not manage the message loop.
// The cleanup function will be called when the application is about to exit, it can be nil if no cleanup is needed.
// The external initialization code must call NewBare in the main thread that runs the message loop.
func NewBare(cleanup func()) (app *BareApp) {
	return newBare(true, cleanup)
}

func newBare(hookGetMsg bool, cleanup func()) (app *BareApp) {
	if threadLocalApp() != nil {
		panic("app already exists in this thread")
	}
	app = &BareApp{
		uiThreadId:        win32.DWORD(windows.GetCurrentThreadId()),
		postMap:           safeMap{ObjectMap: objectmap.New[func()](1, math.MaxUint)},
		msgPreTranslators: make(map[win32.HWND]func(msg *win32.MSG) bool),
		msgRetListeners:   make(map[MessageRetListenerKey]MessageRetListener),
		menuMap:           make(map[win32.HMENU]unsafe.Pointer),
		menuItemMap:       objectmap.New[unsafe.Pointer](internal.MinMenuItemID, internal.MaxMenuItemID),
		cleanup:           cleanup,
	}

	// Prepare postMap
	// See https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-postthreadmessagew#remarks

	// Initialize message queue
	win32.PeekMessageW(&win32.MSG{}, 0, 0, 0, win32.PM_NOREMOVE)

	// Install WH_GETMESSAGE hook
	var msgHookProc func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT
	if hookGetMsg {
		msgHookProc = func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
			if code >= 0 && win32.PeekMessageFlag(wParam) == win32.PM_REMOVE {
				msg := (*win32.MSG)(unsafe.Add(nil, lParam))
				switch msg.Message {
				case win32.WM_QUIT:
					if app.cleanup != nil {
						app.cleanup()
					}
					app.destroy() // Destroy the app when the message loop is about to exit.
				case appmsg.POST:
					// Handle posted functions
					app.postMap.Value(objectmap.Handle(msg.WParam))()
					msg.Message = win32.WM_NULL // Stop WNDPROC processing
					return 0                    // Stop other hook processing
				default:
					if app.preTranslateMessage(msg) {
						msg.Message = win32.WM_NULL // Stop WNDPROC processing
						return 0                    // Stop other hook processing
					}
				}
			}
			return win32.CallNextHookEx(app.getMsgHook, code, wParam, lParam)
		}
	} else {
		msgHookProc = func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
			if code >= 0 && win32.PeekMessageFlag(wParam) == win32.PM_REMOVE {
				if msg := (*win32.MSG)(unsafe.Add(nil, lParam)); msg.Message == appmsg.POST {
					// Handle posted functions
					app.postMap.Value(objectmap.Handle(msg.WParam))()
					msg.Message = win32.WM_NULL // Stop WNDPROC processing
					return 0                    // Stop other hook processing
				}
			}
			return win32.CallNextHookEx(app.getMsgHook, code, wParam, lParam)
		}
	}
	app.getMsgHook = gg.Must(win32.SetWindowsHookExW(win32.WH_GETMESSAGE, windows.NewCallback(msgHookProc), 0, app.uiThreadId))

	gg.MustOK(win32.TlsSetValue(tlsApp, win32.PVOID(unsafe.Pointer(app))))
	return
}

// AddMessageRetListener adds a listener that is called after a message is processed in any window procedure.
// The returned key can be used to remove the listener by calling [RemoveMessageRetListener].
func (app *BareApp) AddMessageRetListener(listener MessageRetListener) (key MessageRetListenerKey) {
	key = MessageRetListenerKey{p: &listener}
	app.msgRetListeners[key] = listener
	return
}

// RemoveMessageRetListener removes the listener added by AddMessageRetListener.
func (app *BareApp) RemoveMessageRetListener(key MessageRetListenerKey) {
	delete(app.msgRetListeners, key)
}

// destroy destroys the app for an external message loop.
// This function must be called in the external message loop before exiting.
// After calling this function, the app should not be used anymore.
func (app *BareApp) destroy() {
	gg.MustOK(win32.UnhookWindowsHookEx(app.getMsgHook))
	gg.MustOK(win32.TlsSetValue(tlsApp, nil))
}

// preTranslateMessage should be called in the external message loop before TranslateMessage.
// If it returns true, the message must not be passed to TranslateMessage and DispatchMessage.
func (app *BareApp) preTranslateMessage(msg *win32.MSG) (processed bool) {
	for _, translator := range app.msgPreTranslators {
		if translator(msg) {
			return true
		}
	}
	if p := app.msgPreTranslators[win32.GetActiveWindow()]; p != nil {
		return p(msg)
	}
	return false
}

// callMsgRetListeners calls all message return listeners with the given parameters.
// It should be called in the WNDPROC after processing a message and before returning the result.
func (app *BareApp) callMsgRetListeners(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT) {
	for _, listener := range app.msgRetListeners {
		listener(hwnd, message, wParam, lParam, result)
	}
}

// Post put f into the UI message queue, f will run in the UI thread ASAP.
func (app *BareApp) Post(f func()) error {
	var h objectmap.Handle
	h = app.postMap.Add(func() {
		f()
		app.postMap.Remove(h)
	})
	return win32.PostThreadMessageW(app.uiThreadId, appmsg.POST, win32.WPARAM(h), 0)
}

// Quit calls win32.PostQuitMessage which tells the message loop to exit.
// The exit code will be the return value of Run.
func (app *BareApp) Quit(exitCode int) {
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
