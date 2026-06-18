package app

import (
	"math"
	"runtime"
	"sync"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/internal/msghandler"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/win32"
	"golang.org/x/sys/windows"
)

// tlsApp is the TLS index for storing *BaseApp in the UI thread.
var tlsApp = gg.Must(win32.TlsAlloc())

//go:linkname threadLocalApp github.com/mkch/gw/internal/app.ThreadLocalApp
func threadLocalApp() *BaseApp {
	return (*BaseApp)(gg.Must(win32.TlsGetValue(tlsApp)))
}

//go:linkname app_AddMsgPreTranslator github.com/mkch/gw/internal/app.AddMsgPreTranslator
func app_AddMsgPreTranslator(app *BaseApp, hwnd win32.HWND, translator func(msg *win32.MSG) bool) {
	app.msgPreTranslators[hwnd] = translator
}

//go:linkname app_RemoveMsgPreTranslator github.com/mkch/gw/internal/app.RemoveMsgPreTranslator
func app_RemoveMsgPreTranslator(app *BaseApp, hwnd win32.HWND) {
	delete(app.msgPreTranslators, hwnd)
}

//go:linkname app_callMsgRetListeners github.com/mkch/gw/internal/app.CallMsgRetListeners
func app_callMsgRetListeners(app *BaseApp, hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT) {
	app.callMsgRetListeners(hwnd, message, wParam, lParam, result)
}

//go:linkname app_MsgHandlers github.com/mkch/gw/internal/app.MsgHandlers
func app_MsgHandlers(app *BaseApp) map[win32.UINT]*msghandler.Chain {
	return app.msgHandlers
}

// MessageRetListener is a function type used by [BaseApp.AddMessageRetListener] and [App.AddMessageRetListener].
type MessageRetListener func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT)

// MessageRetListenerKey is the key for a listener added by [BaseApp.AddMessageRetListener]  and [App.AddMessageRetListener].
// It can be used to remove the listener by calling [BaseApp.RemoveMessageRetListener] and [App.RemoveMessageRetListener].
type MessageRetListenerKey struct{ p *MessageRetListener }

// BaseApp provides basic functionalities for managing the gw application.
// Usually [App] is sufficient for most use cases.
type BaseApp struct {
	uiThreadId        win32.DWORD
	postMap           safeMap
	getMsgHook        win32.HHOOK
	msgPreTranslators map[win32.HWND]func(msg *win32.MSG) bool
	msgRetListeners   map[MessageRetListenerKey]MessageRetListener
	msgHandlers       map[win32.UINT]*msghandler.Chain
	pinner            runtime.Pinner
}

// NewBase creates a [BaseApp] that do not manage the message loop.
// When integrating with external code that already has a message loop, call this function to create a
// [BaseApp] and then operate gw in that goroutine.
// Any other gw functionalities must be used in the same goroutine as the one that calls this function.
// The external initialization code must call this function in the main thread that runs the message loop.
// The returned app should be destroyed by calling [BaseApp.Destroy] after the message loop exits and no gw
// operation should be performed after that.
func NewBase() (app *BaseApp) {
	return (&BaseApp{}).init(true)
}

func (b *BaseApp) init(hookGetMsg bool) *BaseApp {
	if threadLocalApp() != nil {
		panic("app already exists in this thread")
	}
	b.uiThreadId = win32.DWORD(windows.GetCurrentThreadId())
	b.postMap = safeMap{ObjectMap: objectmap.New[func()](1, math.MaxUint)}
	b.msgPreTranslators = make(map[win32.HWND]func(msg *win32.MSG) bool)
	b.msgRetListeners = make(map[MessageRetListenerKey]MessageRetListener)
	b.msgHandlers = make(map[win32.UINT]*msghandler.Chain)

	// Prepare postMap
	// See https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-postthreadmessagew#remarks

	// Initialize message queue

	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-peekmessagew
	// PeekMessageW "dispatches incoming nonqueued messages..."
	//
	// If undestroyed windows from a previous instance still exist, they may receive
	// messages when the new instance calls win32.PeekMessageW before the TLS value is set.
	// This will cause a nil-pointer panic in the window procedure.
	win32.PeekMessageW(&win32.MSG{}, 0, 0, 0, win32.PM_NOREMOVE)

	// Install WH_GETMESSAGE hook
	var msgHookProc func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT
	if hookGetMsg {
		msgHookProc = func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
			if code >= 0 && win32.PeekMessageFlag(wParam) == win32.PM_REMOVE {
				msg := (*win32.MSG)(unsafe.Add(nil, lParam))
				switch msg.Message {
				case appmsg.POST:
					// Handle posted functions
					b.postMap.Value(objectmap.Handle(msg.WParam))()
					msg.Message = win32.WM_NULL // Stop WNDPROC processing
					return 0                    // Stop other hook processing
				default:
					if b.preTranslateMessage(msg) {
						msg.Message = win32.WM_NULL // Stop WNDPROC processing
						return 0                    // Stop other hook processing
					}
				}
			}
			return win32.CallNextHookEx(b.getMsgHook, code, wParam, lParam)
		}
	} else {
		msgHookProc = func(code win32.HookCode, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
			if code >= 0 && win32.PeekMessageFlag(wParam) == win32.PM_REMOVE {
				if msg := (*win32.MSG)(unsafe.Add(nil, lParam)); msg.Message == appmsg.POST {
					// Handle posted functions
					b.postMap.Value(objectmap.Handle(msg.WParam))()
					msg.Message = win32.WM_NULL // Stop WNDPROC processing
					return 0                    // Stop other hook processing
				}
			}
			return win32.CallNextHookEx(b.getMsgHook, code, wParam, lParam)
		}
	}
	b.getMsgHook = gg.Must(win32.SetWindowsHookExW(win32.WH_GETMESSAGE, windows.NewCallback(msgHookProc), 0, b.uiThreadId))

	b.pinner.Pin(b)
	gg.MustOK(win32.TlsSetValue(tlsApp, win32.PVOID(unsafe.Pointer(b))))
	return b
}

// AddMessageRetListener adds a listener that is called after a message is processed in any window procedure.
// The returned key can be used to remove the listener by calling [RemoveMessageRetListener].
func (b *BaseApp) AddMessageRetListener(listener MessageRetListener) (key MessageRetListenerKey) {
	key = MessageRetListenerKey{p: &listener}
	b.msgRetListeners[key] = listener
	return
}

// RemoveMessageRetListener removes the listener added by AddMessageRetListener.
func (b *BaseApp) RemoveMessageRetListener(key MessageRetListenerKey) {
	delete(b.msgRetListeners, key)
}

// Destroy destroys the app for an external message loop.
// This function must be called in the external message loop before exiting.
// After calling this function, the app should not be used anymore and no gw
// function should be called after that.
func (b *BaseApp) Destroy() {
	gg.MustOK(win32.UnhookWindowsHookEx(b.getMsgHook))
	gg.MustOK(win32.TlsSetValue(tlsApp, nil))
	b.pinner.Unpin()
	*b = BaseApp{} // Clear all fields
}

// DestroyAllWindows destroys all non-child windows associated with the application's
// UI thread (including message-only windows).
//
// It returns the number of windows successfully destroyed.
//
// This method is particularly useful in unit tests to ensure a clean state after
// each test case, preventing stale windows from interfering with subsequent tests.
func (b *BaseApp) DestroyAllWindows() (nDestroyed int) {
	// Enumerate all non-child windows associated with the UI thread and destroy them.
	win32.EnumThreadWindows(b.uiThreadId, func(hwnd win32.HWND) bool {
		win32.DestroyWindow(hwnd)
		nDestroyed++
		return true
	})
	// Find all message-only windows and destroy them.
	var hwnd win32.HWND
	for {
		hwnd = win32.FindWindowExW(win32.HWND_MESSAGE, hwnd, nil, nil)
		if hwnd == 0 {
			break
		}
		win32.DestroyWindow(hwnd)
		nDestroyed++
	}
	return
}

// preTranslateMessage should be called in the external message loop before TranslateMessage.
// If it returns true, the message must not be passed to TranslateMessage and DispatchMessage.
func (b *BaseApp) preTranslateMessage(msg *win32.MSG) (processed bool) {
	for _, translator := range b.msgPreTranslators {
		if translator(msg) {
			return true
		}
	}
	if p := b.msgPreTranslators[win32.GetActiveWindow()]; p != nil {
		return p(msg)
	}
	return false
}

// callMsgRetListeners calls all message return listeners with the given parameters.
// It should be called in the WNDPROC after processing a message and before returning the result.
func (b *BaseApp) callMsgRetListeners(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, result win32.LRESULT) {
	// If a nil-pointer panic occurs here, there may be windows that were not properly destroyed
	// by a previous app instance.
	// See wndProc variable and [BaseApp.Init] for details.
	for _, listener := range b.msgRetListeners {
		listener(hwnd, message, wParam, lParam, result)
	}
}

// Post put f into the UI message queue, f will run in the UI thread ASAP.
func (b *BaseApp) Post(f func()) error {
	var h objectmap.Handle
	h = b.postMap.Add(func() {
		f()
		b.postMap.Remove(h)
	})
	return win32.PostThreadMessageW(b.uiThreadId, appmsg.POST, win32.WPARAM(h), 0)
}

// Quit calls win32.PostQuitMessage which tells the message loop to exit.
// The exit code will be the return value of [win32.GetMessage].
func (b *BaseApp) Quit(exitCode int) {
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
