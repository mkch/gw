package window

import (
	"errors"
	"sync"

	"github.com/mkch/gg"
	"github.com/mkch/gw/internal"
	"github.com/mkch/gw/internal/app"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"golang.org/x/sys/windows"
)

const defClassName = "github.com/mkch/gw#Window"

var defClassSpec = win32util.WndClass{
	WndProc: func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
		return win32.DefWindowProcW(hwnd, message, wParam, lParam)
	},
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
	ClassName string // Custom class name. If empty, the default class will be used.
	Text      string
	Style     win32.WINDOW_STYLE
	ExStyle   win32.WINDOW_EX_STYLE
	X         metrics.Dimension
	Y         metrics.Dimension
	Width     metrics.Dimension
	Height    metrics.Dimension
	WndParent win32.HWND
	Menu      *menu.Menu
	Instance  win32.HINSTANCE // 0 for this module.
	OnCreate  func()
	OnClose   func() bool // Return true to allow closing, false to prevent.
	OnDestroy func()
}

type Window struct {
	WindowBase
	OnCreate  func()
	OnClose   func() bool
	OnDestroy func()

	menu *menu.Menu

	menuAccel          []menu.ItemAccel // Accelerator table of the window menu.
	popupMenuAccel     []menu.ItemAccel // Accelerator table of the popup menu(context menu).
	accelKeyTable      win32.HACCEL
	accelToMenuItemMap map[win32.WORD]*menu.Item // ID of accelerator to the corresponding menu item. Used in WM_COMMAND  handlers to find the menu item of an accelerator command.
}

func New(spec *Spec) (*Window, error) {
	className := spec.ClassName
	if className == "" {
		// Use default class name.
		className = defClassName
	}
	setRegistered, registered := classRegistered()
	if !registered(className) {
		cls := new(defClassSpec)
		cls.ClassName = className
		if _, err := win32util.RegisterClass(cls); err != nil {
			return nil, err
		}
		setRegistered(className)
	}

	var visible bool
	var style = spec.Style
	var showCmd win32.SHOW_WINDOW_CMD = -1
	visible = spec.Style&win32.WS_VISIBLE != 0
	style = spec.Style &^ win32.WS_VISIBLE
	var x, y, cx, cy win32.INT
	var useDefPos, useDefSize bool
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-createwindowexw
	// If an overlapped window is created with the WS_VISIBLE style bit set and the x parameter is set to CW_USEDEFAULT,
	// then the y parameter determines how the window is shown.
	if spec.X.Unit == metrics.PX && spec.X.Value == win32.CW_USEDEFAULT {
		useDefPos = true
		x = win32.CW_USEDEFAULT
		if visible {
			// If an overlapped window is created with the WS_VISIBLE style bit set and the x parameter is set to CW_USEDEFAULT,
			// then the y parameter determines how the window is shown.
			showCmd = gg.If(spec.Y.Value == win32.CW_USEDEFAULT, win32.SW_SHOW, win32.SHOW_WINDOW_CMD(spec.Y.Value))
		}
	}

	// For overlapped windows, if width is CW_USEDEFAULT, the system selects a default width and height for the window
	// and ignores the height.
	if spec.Width.Unit == metrics.PX && spec.Width.Value == win32.CW_USEDEFAULT {
		useDefSize = true
		cx = win32.CW_USEDEFAULT
	}

	hwnd, err := win32util.CreateWindow((&win32util.Wnd{
		ClassName:  className,
		WindowName: spec.Text,
		Style:      style,
		ExStyle:    spec.ExStyle,
		X:          x,
		Y:          y,
		Width:      cx,
		Height:     cy,
		WndParent:  spec.WndParent,
		Instance:   spec.Instance,
	}))
	if err != nil {
		return nil, err
	}

	// Use the window's own DPI after creating it
	dpi := gg.Must(win32.GetDpiForWindow(hwnd))
	var swpFlags win32.UINT = win32.SWP_NOACTIVATE | win32.SWP_NOZORDER
	if useDefPos {
		swpFlags |= win32.SWP_NOMOVE
	} else {
		x = spec.X.Px(dpi)
		y = spec.Y.Px(dpi)
	}

	if useDefSize {
		swpFlags |= win32.SWP_NOSIZE
	} else {
		cx = spec.Width.Px(dpi)
		cy = spec.Height.Px(dpi)
	}

	if !useDefPos || !useDefSize {
		win32.SetWindowPos(hwnd, 0, x, y, cx, cy, swpFlags)
	}
	if showCmd != -1 {
		win32.ShowWindow(hwnd, showCmd)
	}

	win := &Window{OnCreate: spec.OnCreate, OnClose: spec.OnClose, OnDestroy: spec.OnDestroy}
	if err := Attach(hwnd, &win.WindowBase); err != nil {
		return nil, err
	}
	if win.OnCreate != nil {
		win.OnCreate()
	}

	win.SetWndProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevWndProc win32.WndProc) win32.LRESULT {
		switch message {
		case win32.WM_COMMAND:
			if lParam != 0 { // Control command
				win32.SendMessageW(win32.HWND(lParam), appmsg.REFLECT_COMMAND, wParam, lParam)
			} else { // Menu or accelerator command
				// Because all the menus are notified by position,
				// only accelerator commands goes here.
				if item, ok := win.accelToMenuItemMap[win32.LOWORD(wParam)]; ok {
					item.OnClick()
				}
			}
			return 0
		case win32.WM_MENUCOMMAND:
			menu.OnWmMenuCommand(wParam, lParam)
			return 0
		case win32.WM_CLOSE:
			if win.OnClose != nil {
				if !win.OnClose() {
					return 0 // prevent closing
				}
			}
		case win32.WM_DESTROY:
			win.setMsgPreTranslator(nil)
			if win.menu != nil {
				// Although the menu is automatically destroyed by the system when the window is destroyed,
				// the cleanup for menuMap is not called then.

				// Don't do this after WM_DESTROY because SetMenu(0) causes WM_SIZE and related messages
				// to be sent which may cause problems in event handlers if the window handle is already destroyed.
				win32.SetMenu(win.hwnd, 0)
				// Manually destroy the menu to avoid resource leak in menuMap.
				win.menu.Destroy()
			}
			if win.accelKeyTable != 0 {
				win32.DestroyAcceleratorTable(win.accelKeyTable)
				win.accelKeyTable = 0
			}

			if win.OnDestroy != nil {
				win.OnDestroy()
			}
		case win32.WM_INPUTLANGCHANGE:
			if win.menu != nil {
				win.menu.RefreshDisplayTitle()
			}
		}
		return prevWndProc(hwnd, message, wParam, lParam)
	})
	if spec.Menu != nil {
		win.SetMenu(spec.Menu)
	}
	return win, nil
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
		h, err := win32.CreateAcceleratorTableW(table)
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
func (w *Window) setMsgPreTranslator(p msgProc) {
	if p == nil {
		app.RemoveMsgPreTranslator(app.ThreadLocalApp(), w.hwnd)
	} else {
		app.AddMsgPreTranslator(app.ThreadLocalApp(), w.hwnd, p)
	}
}

func (w *Window) SetMenu(menu *menu.Menu) error {
	return w.setMenu(menu)
}

func (w *Window) preTranslateMessage(p *win32.MSG) bool {
	if w.accelKeyTable == 0 {
		return false
	}
	hwnd := w.hwnd
	ok, err := win32.TranslateAcceleratorW(hwnd, w.accelKeyTable, p)
	// TranslateAcceleratorW can send WM_COMMAND where the handler may destroy the window,
	// causing TranslateAcceleratorW returning ERROR_INVALID_WINDOW_HANDLE.
	if err != nil && !errors.Is(err, windows.ERROR_INVALID_WINDOW_HANDLE) {
		panic(err)
	}
	return ok
}
