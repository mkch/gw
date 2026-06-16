package window

import (
	"errors"
	"sync"

	"github.com/mkch/gg"
	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/internal"
	"github.com/mkch/gw/internal/app"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"golang.org/x/sys/windows"
)

const defClassName = "github.com/mkch/gw#Window"

var defClassSpec = win32util.WndClass{
	WndProc:    win32.DefWindowProcW,
	Background: win32.HBRUSH(win32.COLOR_WINDOW + 1),
	Cursor:     gg.Must(win32.LoadImageW_uintptr[win32.HCURSOR](0, uintptr(win32.OCR_NORMAL), win32.IMAGE_CURSOR, 0, 0, win32.LR_DEFAULTSIZE|win32.LR_SHARED)),
}

// classRegistered is a function that returns two functions: setRegistered and registered.
// setRegistered marks a class name as registered, and registered checks if a class name is marked as registered.
// The returned functions are concurrency-safe.
var classRegistered = sync.OnceValues(func() (func(string), func(string) bool) {
	var lock sync.RWMutex
	var registeredClasses = make(gg.Set[string])
	return func(s string) {
			lock.Lock()
			defer lock.Unlock()
			registeredClasses.Add(s)
		}, func(s string) bool {
			lock.RLock()
			defer lock.RUnlock()
			return registeredClasses.Contains(s)
		}
})

type Spec struct {
	Parent    win32.HWND
	ClassName string // Custom class name. If empty, the default class will be used.
	Text      string
	Style     win32.WindowStyle
	ExStyle   win32.WindowExStyle
	X         metrics.Dimension
	Y         metrics.Dimension
	Width     metrics.Dimension
	Height    metrics.Dimension
	Menu      *menu.Menu
	Instance  win32.HINSTANCE // 0 for this module.

	OnCreate  func()
	OnClose   func() bool // Return true to allow closing, false to prevent.
	OnDestroy func()
}

// TopLevel is an interface for top-level windows, which can be closed by user.
type TopLevel interface {
	gw.BaseWindow
	Close() bool
	SetOnCloseListener(listener func() bool)
	SetMenu(menu *menu.Menu) error
}

// Ensure [*Window] implements [TopLevel].
// Should have no performance overhead.
var _ TopLevel = (*Window)(nil)

type Window struct {
	gw.BaseWindowImpl
	// Spec is used to create the window and is set to nil after creation.
	Spec *Spec

	menu               *menu.Menu
	menuAccel          []menu.ItemAccel // Accelerator table of the window menu.
	popupMenuAccel     []menu.ItemAccel // Accelerator table of the popup menu(context menu).
	accelKeyTable      win32.HACCEL
	accelToMenuItemMap map[win32.WORD]*menu.Item // ID of accelerator to the corresponding menu item. Used in WM_COMMAND  handlers to find the menu item of an accelerator command.

	onCreateHandler  func()
	onCloseHandler   func() bool
	onDestroyHandler func()
}

func (w *Window) Close() bool {
	ret, _ := win32.SendMessageW(w.HWND(), win32.WM_CLOSE, 0, 0)
	// If the WM_CLOSE handler allows, it returns 0. Otherwise -1.
	return ret == 0
}

// SetOnCloseListener sets a listener for the close event of the popup window.
// Listener returns true to allow closing, false to prevent.
func (w *Window) SetOnCloseListener(listener func() bool) {
	w.onCloseHandler = listener
}

// OnClose() is called when the window receives WM_CLOSE message.
// It returns true to allow closing, false to prevent.
func (w *Window) OnClose() bool {
	return w.callOnCloseListener()
}

func (w *Window) callOnCloseListener() bool {
	if w.onCloseHandler != nil {
		return w.onCloseHandler()
	}
	return true
}

func (w *Window) rebuildAccelTable() error {
	if w.accelKeyTable != 0 {
		if err := win32.DestroyAcceleratorTable(w.accelKeyTable); err != nil {
			return err
		}
		w.accelKeyTable = 0
	}
	clear(w.accelToMenuItemMap)

	var table []win32.ACCEL
	id := internal.MinMenuItemID

	// Fill table and w.accelToMenuItemMap from t.
	// A unique ID is assigned to each accelerator.
	processMenuItemAccelTable := func(t []menu.ItemAccel) {
		for _, accel := range t {
			if id > internal.MaxMenuItemID {
				panic("out of menu item IDs")
			}
			accel.Accel.Cmd = win32.WORD(id)
			table = append(table, accel.Accel)
			if w.accelToMenuItemMap == nil {
				w.accelToMenuItemMap = make(map[win32.WORD]*menu.Item)
			}
			w.accelToMenuItemMap[accel.Accel.Cmd] = accel.Item
			id++
			if id == win32.IDTIMEOUT {
				id++
			}
		}
	}
	processMenuItemAccelTable(w.menuAccel)
	processMenuItemAccelTable(w.popupMenuAccel)

	if len(table) > 0 {
		h, err := win32.CreateAcceleratorTableW(win32.AlignSlice(table))
		if err != nil {
			return err
		}
		w.accelKeyTable = h
	}
	return nil
}

// setMsgPreTranslator sets a MsgProc to process a message sent to this window
// before TranslateMessage is called in the message loop.
// If p returns true, no further processing will be performed.
// A nil p removes the pre-translator.
func (w *Window) setMsgPreTranslator(p func(msg *win32.MSG) bool) {
	if p == nil {
		app.RemoveMsgPreTranslator(app.ThreadLocalApp(), w.HWND())
	} else {
		app.AddMsgPreTranslator(app.ThreadLocalApp(), w.HWND(), p)
	}
}

func (w *Window) SetMenu(menu *menu.Menu) error {
	if menu == w.menu {
		return nil
	}
	var hMenu win32.HMENU
	if menu != nil {
		// In case the keyboard layout has been changed since the menu is created.
		menu.OnKeyboardLayoutChange(win32.GetKeyboardLayout(0))
		hMenu = menu.HMENU()
	}
	if err := win32.SetMenu(w.HWND(), hMenu); err != nil {
		return err
	}
	if w.menu != nil {
		w.menu.OnAccelKeyChanged = nil
	}
	hasOldMenu := w.menu != nil
	w.menu = menu
	if w.menu != nil {
		var err error
		if w.menuAccel, err = w.menu.AccelKeyTable(); err != nil {
			return err
		}
		w.menu.OnAccelKeyChanged = func() (err error) {
			if w.menuAccel, err = w.menu.AccelKeyTable(); err != nil {
				return
			}
			if err = w.rebuildAccelTable(); err != nil {
				return
			}
			return
		}
		w.setMsgPreTranslator(w.preTranslateMessage)
	} else {
		if hasOldMenu {
			w.setMsgPreTranslator(nil)
		}
		w.menuAccel = nil
	}

	return w.rebuildAccelTable()
}

func (w *Window) preTranslateMessage(p *win32.MSG) bool {
	if w.accelKeyTable == 0 {
		return false
	}
	hwnd := w.HWND()
	ok, err := win32.TranslateAcceleratorW(hwnd, w.accelKeyTable, p)
	// TranslateAcceleratorW can send WM_COMMAND where the handler may destroy the window,
	// causing TranslateAcceleratorW returning ERROR_INVALID_WINDOW_HANDLE.
	if err != nil && !errors.Is(err, windows.ERROR_INVALID_WINDOW_HANDLE) {
		panic(err)
	}
	return ok
}

func (w *Window) CreateHandle() (win32.HWND, error) {
	if w.Spec == nil {
		w.Spec = &Spec{}
	}
	className := w.Spec.ClassName
	if className == "" {
		// Use default class name.
		className = defClassName
	}
	setRegistered, registered := classRegistered()
	if !registered(className) {
		cls := new(defClassSpec)
		cls.ClassName = className
		chkerr.Must(win32util.RegisterClass(cls))
		setRegistered(className)
	}

	var visible bool
	var style = w.Spec.Style
	var showCmd win32.SHOW_WINDOW_CMD = -1
	visible = w.Spec.Style&win32.WS_VISIBLE != 0
	style = w.Spec.Style &^ win32.WS_VISIBLE
	var x, y, cx, cy win32.INT
	var useDefPos, useDefSize bool
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-createwindowexw
	// If an overlapped window is created with the WS_VISIBLE style bit set and the x parameter is set to CW_USEDEFAULT,
	// then the y parameter determines how the window is shown.
	if w.Spec.X == gw.CW_USEDEFAULT {
		useDefPos = true
		x = win32.CW_USEDEFAULT
		if visible {
			// If an overlapped window is created with the WS_VISIBLE style bit set and the x parameter is set to CW_USEDEFAULT,
			// then the y parameter determines how the window is shown.
			var y win32.INT
			if w.Spec.Y == nil {
				y = 0
			} else {
				y = w.Spec.Y.(metrics.Px).Value() // Must be a Px.
			}
			showCmd = gg.If(y == win32.CW_USEDEFAULT, win32.SW_SHOW, win32.SHOW_WINDOW_CMD(y))
		}
	}

	// For overlapped windows, if width is CW_USEDEFAULT, the system selects a default width and height for the window
	// and ignores the height.
	if w.Spec.Width == gw.CW_USEDEFAULT {
		useDefSize = true
		cx = win32.CW_USEDEFAULT
	}

	hwnd, err := win32util.CreateWindow((&win32util.Wnd{
		WndParent:  w.Spec.Parent,
		ClassName:  className,
		WindowName: w.Spec.Text,
		Style:      style,
		ExStyle:    w.Spec.ExStyle,
		X:          x,
		Y:          y,
		Width:      cx,
		Height:     cy,
		Instance:   w.Spec.Instance,
	}))
	if err != nil {
		return 0, err
	}

	// Use the window's own DPI after creating it
	dpi := gg.Must(win32.GetDpiForWindow(hwnd))
	var swpFlags win32.UINT = win32.SWP_NOACTIVATE | win32.SWP_NOZORDER
	if useDefPos {
		swpFlags |= win32.SWP_NOMOVE
	} else {
		x = metrics.ToPx(w.Spec.X, dpi).Value()
		y = metrics.ToPx(w.Spec.Y, dpi).Value()
	}

	if useDefSize {
		swpFlags |= win32.SWP_NOSIZE
	} else {
		cx = metrics.ToPx(w.Spec.Width, dpi).Value()
		cy = metrics.ToPx(w.Spec.Height, dpi).Value()
	}

	if !useDefPos || !useDefSize {
		win32.SetWindowPos(hwnd, 0, x, y, cx, cy, swpFlags)
	}
	if showCmd != -1 {
		win32.ShowWindow(hwnd, showCmd)
	}

	if w.Spec.Menu != nil {
		if err := w.SetMenu(w.Spec.Menu); err != nil {
			return 0, err
		}
	}
	return hwnd, nil
}

func (w *Window) WndProc(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	switch message {
	case win32.WM_COMMAND:
		if lParam == 0 {
			// Menu or accelerator command
			// Because all the menus are notified by position,
			// only accelerator commands goes here.
			if item, ok := w.accelToMenuItemMap[win32.LOWORD(wParam)]; ok {
				item.OnClick()
			}
		}
	case win32.WM_CLOSE:
		oc, ok := gw.LookupWindow(w.HWND()).(interface{ OnClose() bool })
		if ok && !oc.OnClose() {
			return -1 // prevent the window from closing. Used by [Window.Close].
		}
	case win32.WM_DESTROY:
		w.setMsgPreTranslator(nil)
		if w.menu != nil {
			// Although the menu is automatically destroyed by the system when the window is destroyed,
			// the cleanup for menuMap is not called then.

			// Don't do this after WM_DESTROY because SetMenu(0) causes WM_SIZE and related messages
			// to be sent which may cause problems in event handlers if the window handle is already destroyed.
			win32.SetMenu(w.HWND(), 0)
			// Manually destroy the menu to avoid resource leak in menuMap.
			w.menu.Destroy()
		}
		if w.accelKeyTable != 0 {
			win32.DestroyAcceleratorTable(w.accelKeyTable)
			w.accelKeyTable = 0
		}
	case win32.WM_INPUTLANGCHANGE:
		if w.menu != nil {
			w.menu.OnKeyboardLayoutChange(win32.HKL(lParam))
		}
	}
	return w.BaseWindowImpl.WndProc(hwnd, message, wParam, lParam)
}

func (w *Window) OnInit() error {
	if err := w.BaseWindowImpl.OnInit(); err != nil {
		return err
	}

	defer func() { w.Spec = nil }()
	if w.Spec.OnCreate != nil {
		w.Spec.OnCreate()
	}
	if w.Spec.OnClose != nil {
		w.SetOnCloseListener(w.Spec.OnClose)
	}
	if w.Spec.OnDestroy != nil {
		w.SetOnDestroyListener(w.Spec.OnDestroy)
	}
	return nil
}

func New(spec *Spec) (*Window, error) {
	if spec != nil && spec.Style&win32.WS_CHILD != 0 {
		return nil, errors.New("WS_CHILD is not allowed")
	}
	return gw.Init(&Window{Spec: spec})
}
