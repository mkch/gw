package win32

import (
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/win32/sysutil"
	"golang.org/x/sys/windows"
)

var lzUser32 = windows.NewLazySystemDLL("user32.dll")
var lzKernel32 = windows.NewLazySystemDLL("kernel32")
var lzGdi32 = windows.NewLazySystemDLL("gdi32.dll")

var lzGetMessageW = lzUser32.NewProc("GetMessageW")

type Point struct {
	X, Y LONG
}

type MSG struct {
	Hwnd    HWND
	Message UINT
	WParam  WPARAM
	LParam  LPARAM
	Time    DWORD
	Point   Point
}

func GetMessageW(msg *MSG, hwnd HWND, msgFilterMin UINT, msgFilterMax UINT) BOOL {
	return sysutil.As[BOOL](lzGetMessageW.Call(uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(msgFilterMin), uintptr(msgFilterMax)))
}

type PeekMessageFlag UINT

const (
	PM_NOREMOVE = PeekMessageFlag(0x0000)
	PM_REMOVE   = PeekMessageFlag(0x0001)
	PM_NOYIELD  = PeekMessageFlag(0x0002)
)

var lzPeekMessageW = lzUser32.NewProc("PeekMessageW")

func PeekMessageW(msg *MSG, hwnd HWND, msgFilterMin UINT, msgFilterMax UINT, flags PeekMessageFlag) BOOL {
	return sysutil.As[BOOL](lzPeekMessageW.Call(uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(msgFilterMin), uintptr(msgFilterMax), uintptr(flags)))
}

var lzTranslateMessage = lzUser32.NewProc("TranslateMessage")

func TranslateMessage(msg *MSG) bool {
	return sysutil.AsBool(lzTranslateMessage.Call(uintptr(unsafe.Pointer(msg))))
}

var lzDispatchMessageW = lzUser32.NewProc("DispatchMessageW")

func DispatchMessageW(msg *MSG) LRESULT {
	return sysutil.As[LRESULT](lzDispatchMessageW.Call(uintptr(unsafe.Pointer(msg))))
}

var lzPostQuitMessage = lzUser32.NewProc("PostQuitMessage")

func PostQuitMessage(code int) {
	lzPostQuitMessage.Call(uintptr(code))
}

type CLASS_STYLE UINT

type WNDCLASSEXW struct {
	Size       UINT
	Style      CLASS_STYLE
	WndProc    uintptr
	ClsExtra   INT
	WndExtra   INT
	Instance   HINSTANCE
	Icon       HICON
	Cursor     HCURSOR
	Background HBRUSH
	MenuName   *WCHAR
	ClassName  *WCHAR
	IconSm     HICON
}

const (
	CS_BYTEALIGNCLIENT CLASS_STYLE = 0x1000
	CS_BYTEALIGNWINDOW CLASS_STYLE = 0x2000
	CS_CLASSDC         CLASS_STYLE = 0x0040
	CS_DBLCLKS         CLASS_STYLE = 0x0008
	CS_DROPSHADOW      CLASS_STYLE = 0x0002000
	CS_GLOBALCLASS     CLASS_STYLE = 0x4000
	CS_HREDRAW         CLASS_STYLE = 0x0002
	CS_NOCLOSE         CLASS_STYLE = 0x0200
	CS_OWNDC           CLASS_STYLE = 0x0020
	CS_PARENTDC        CLASS_STYLE = 0x0080
	CS_SAVEBITS        CLASS_STYLE = 0x0800
	CS_VREDRAW         CLASS_STYLE = 0x0001
)

var lzRegisterClassExW = lzUser32.NewProc("RegisterClassExW")

func RegisterClassExW(cls *WNDCLASSEXW) (ATOM, error) {
	return sysutil.MustNotZero[ATOM](lzRegisterClassExW.Call(uintptr(unsafe.Pointer(cls))))
}

type WINDOW_EX_STYLE DWORD

const (
	WS_EX_ACCEPTFILES         WINDOW_EX_STYLE = 0x00000010
	WS_EX_APPWINDOW           WINDOW_EX_STYLE = 0x00040000
	WS_EX_CLIENTEDGE          WINDOW_EX_STYLE = 0x00000200
	WS_EX_COMPOSITED          WINDOW_EX_STYLE = 0x02000000
	WS_EX_CONTEXTHELP         WINDOW_EX_STYLE = 0x00000400
	WS_EX_CONTROLPARENT       WINDOW_EX_STYLE = 0x00010000
	WS_EX_DLGMODALFRAME       WINDOW_EX_STYLE = 0x00000001
	WS_EX_LAYERED             WINDOW_EX_STYLE = 0x00080000
	WS_EX_LAYOUTRTL           WINDOW_EX_STYLE = 0x00400000
	WS_EX_LEFT                WINDOW_EX_STYLE = 0x00000000
	WS_EX_LEFTSCROLLBAR       WINDOW_EX_STYLE = 0x00004000
	WS_EX_LTRREADING          WINDOW_EX_STYLE = 0x00000000
	WS_EX_MDICHILD            WINDOW_EX_STYLE = 0x00000040
	WS_EX_NOACTIVATE          WINDOW_EX_STYLE = 0x08000000
	WS_EX_NOINHERITLAYOUT     WINDOW_EX_STYLE = 0x00100000
	WS_EX_NOPARENTNOTIFY      WINDOW_EX_STYLE = 0x00000004
	WS_EX_NOREDIRECTIONBITMAP WINDOW_EX_STYLE = 0x00200000
	WS_EX_OVERLAPPEDWINDOW    WINDOW_EX_STYLE = WS_EX_WINDOWEDGE | WS_EX_CLIENTEDGE
	WS_EX_PALETTEWINDOW       WINDOW_EX_STYLE = WS_EX_WINDOWEDGE | WS_EX_TOOLWINDOW | WS_EX_TOPMOST
	WS_EX_RIGHT               WINDOW_EX_STYLE = 0x00001000
	WS_EX_RIGHTSCROLLBAR      WINDOW_EX_STYLE = 0x00000000
	WS_EX_RTLREADING          WINDOW_EX_STYLE = 0x00002000
	WS_EX_STATICEDGE          WINDOW_EX_STYLE = 0x00020000
	WS_EX_TOOLWINDOW          WINDOW_EX_STYLE = 0x00000080
	WS_EX_TOPMOST             WINDOW_EX_STYLE = 0x00000008
	WS_EX_TRANSPARENT         WINDOW_EX_STYLE = 0x00000020
	WS_EX_WINDOWEDGE          WINDOW_EX_STYLE = 0x00000100
)

type WINDOW_STYLE DWORD

const (
	WS_BORDER           WINDOW_STYLE = 0x00800000
	WS_CAPTION          WINDOW_STYLE = 0x00C00000
	WS_CHILD            WINDOW_STYLE = 0x40000000
	WS_CHILDWINDOW      WINDOW_STYLE = 0x40000000
	WS_CLIPCHILDREN     WINDOW_STYLE = 0x02000000
	WS_CLIPSIBLINGS     WINDOW_STYLE = 0x04000000
	WS_DISABLED         WINDOW_STYLE = 0x08000000
	WS_DLGFRAME         WINDOW_STYLE = 0x00400000
	WS_GROUP            WINDOW_STYLE = 0x00020000
	WS_HSCROLL          WINDOW_STYLE = 0x00100000
	WS_ICONIC           WINDOW_STYLE = 0x20000000
	WS_MAXIMIZE         WINDOW_STYLE = 0x01000000
	WS_MAXIMIZEBOX      WINDOW_STYLE = 0x00010000
	WS_MINIMIZE         WINDOW_STYLE = 0x20000000
	WS_MINIMIZEBOX      WINDOW_STYLE = 0x00020000
	WS_OVERLAPPED       WINDOW_STYLE = 0x00000000
	WS_OVERLAPPEDWINDOW WINDOW_STYLE = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_POPUP            WINDOW_STYLE = 0x80000000
	WS_POPUPWINDOW      WINDOW_STYLE = WS_POPUP | WS_BORDER | WS_SYSMENU
	WS_SIZEBOX          WINDOW_STYLE = 0x0004000
	WS_SYSMENU          WINDOW_STYLE = 0x00080000
	WS_TABSTOP          WINDOW_STYLE = 0x00010000
	WS_THICKFRAME       WINDOW_STYLE = 0x00040000
	WS_TILED            WINDOW_STYLE = 0x00000000
	WS_TILEDWINDOW      WINDOW_STYLE = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_VISIBLE          WINDOW_STYLE = 0x10000000
	WS_VSCROLL          WINDOW_STYLE = 0x00200000
)

const CW_USEDEFAULT INT = -2147483648 //0x80000000

var lzCreateWindowExW = lzUser32.NewProc("CreateWindowExW")

func CreateWindowExW(
	exStyle WINDOW_EX_STYLE,
	className *WCHAR,
	windowName *WCHAR,
	style WINDOW_STYLE,
	x INT, y INT, width INT, height INT,
	wndParent HWND,
	menu HMENU,
	instance HINSTANCE,
	param UINT_PTR,
) (HWND, error) {
	return sysutil.MustNotZero[HWND](lzCreateWindowExW.Call(uintptr(exStyle), uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)), uintptr(style),
		uintptr(x), uintptr(y), uintptr(width), uintptr(height),
		uintptr(wndParent), uintptr(menu), uintptr(instance), uintptr(param)))
}

var lzDestroyWindow = lzUser32.NewProc("DestroyWindow")

func DestroyWindow(hwnd HWND) error {
	return sysutil.MustTrue(lzDestroyWindow.Call(uintptr(hwnd)))
}

var lzDefWindowProcW = lzUser32.NewProc("DefWindowProcW")

func DefWindowProcW(hwnd HWND, message UINT, wParam WPARAM, lParam LPARAM) LRESULT {
	return sysutil.As[LRESULT](lzDefWindowProcW.Call(uintptr(hwnd), uintptr(message), uintptr(wParam), uintptr(lParam)))
}

var lzGetModuleHandleW = lzKernel32.NewProc("GetModuleHandleW")

func GetModuleHandleW[H HMODULE | HINSTANCE](moduleName *WCHAR) (H, error) {
	return sysutil.MustNotZero[H](lzGetModuleHandleW.Call(uintptr(unsafe.Pointer(moduleName))))
}

const (
	COLOR_3DDKSHADOW              = 21
	COLOR_3DFACE                  = 15
	COLOR_3DHIGHLIGHT             = 20
	COLOR_3DHILIGHT               = 20
	COLOR_3DLIGHT                 = 22
	COLOR_3DSHADOW                = 16
	COLOR_ACTIVEBORDER            = 10
	COLOR_ACTIVECAPTION           = 2
	COLOR_APPWORKSPACE            = 12
	COLOR_BACKGROUND              = 1
	COLOR_BTNFACE                 = 15
	COLOR_BTNHIGHLIGHT            = 20
	COLOR_BTNHILIGHT              = 20
	COLOR_BTNSHADOW               = 16
	COLOR_BTNTEX                  = 18
	COLOR_CAPTIONTEXT             = 9
	COLOR_DESKTOP                 = 1
	COLOR_GRADIENTACTIVECAPTION   = 27
	COLOR_GRADIENTINACTIVECAPTION = 28
	COLOR_GRAYTEXT                = 17
	COLOR_HIGHLIGHT               = 13
	COLOR_HIGHLIGHTTEXT           = 14
	COLOR_HOTLIGHT                = 26
	COLOR_INACTIVEBORDER          = 11
	COLOR_INACTIVECAPTION         = 3
	COLOR_INACTIVECAPTIONTEXT     = 19
	COLOR_INFOBK                  = 24
	COLOR_INFOTEXT                = 23
	COLOR_MENU                    = 4
	COLOR_MENUHILIGHT             = 29
	COLOR_MENUBAR                 = 30
	COLOR_MENUTEXT                = 7
	COLOR_SCROLLBAR               = 0
	COLOR_WINDOW                  = 5
	COLOR_WINDOWFRAME             = 6
	COLOR_WINDOWTEXT              = 8
)

var lzGetSysColor = lzUser32.NewProc("GetSysColor")

func GetSysColor(index int) DWORD {
	return sysutil.As[DWORD](lzGetSysColor.Call(uintptr(index)))
}

type SHOW_WINDOW_CMD INT

const (
	SW_HIDE            SHOW_WINDOW_CMD = 0
	SW_SHOWNORMAL      SHOW_WINDOW_CMD = 1
	SW_NORMAL          SHOW_WINDOW_CMD = 1
	SW_SHOWMINIMIZED   SHOW_WINDOW_CMD = 2
	SW_SHOWMAXIMIZED   SHOW_WINDOW_CMD = 3
	SW_MAXIMIZE        SHOW_WINDOW_CMD = 3
	SW_SHOWNOACTIVATE  SHOW_WINDOW_CMD = 4
	SW_SHOW            SHOW_WINDOW_CMD = 5
	SW_MINIMIZE        SHOW_WINDOW_CMD = 6
	SW_SHOWMINNOACTIVE SHOW_WINDOW_CMD = 7
	SW_SHOWNA          SHOW_WINDOW_CMD = 8
	SW_RESTORE         SHOW_WINDOW_CMD = 9
	SW_SHOWDEFAULT     SHOW_WINDOW_CMD = 10
	SW_FORCEMINIMIZE   SHOW_WINDOW_CMD = 11
)

var lzShowWindow = lzUser32.NewProc("ShowWindow")

func ShowWindow(hwnd HWND, cmdShow SHOW_WINDOW_CMD) error {
	return sysutil.MustTrue(lzShowWindow.Call(uintptr(hwnd), uintptr(cmdShow)))
}

type WndProc = func(hwnd HWND, message UINT, wParam WPARAM, lParam LPARAM) LRESULT

const (
	GWL_EXSTYLE     = -20
	GWLP_HINSTANCE  = -6
	GWLP_HWNDPARENT = -8
	GWLP_ID         = -12
	GWL_STYLE       = -16
	GWLP_USERDATA   = -21
	GWLP_WNDPROC    = -4
	DWLP_DLGPROC    = DWLP_MSGRESULT + unsafe.Sizeof(LRESULT(0))
	DWLP_MSGRESULT  = 0
	DWLP_USER       = DWLP_DLGPROC + unsafe.Sizeof(UINT_PTR(0))
)

func SetWindowLongPtrW(hwnd HWND, index int, newLong LONG_PTR) (LONG_PTR, error) {
	r, _, err := lzSetWindowLongPtrW.Call(uintptr(hwnd), uintptr(index), uintptr(newLong))
	if r != 0 || sysutil.IsNoError(err) {
		err = nil
	}
	return LONG_PTR(r), err
}

func GetWindowLongPtrW(hwnd HWND, index int) (LONG_PTR, error) {
	r, _, err := lzGetWindowLongPtrW.Call(uintptr(hwnd), uintptr(index))
	if r != 0 || sysutil.IsNoError(err) {
		err = nil
	}
	return LONG_PTR(r), err
}

var lzSetPropW = lzUser32.NewProc("SetPropW")

func SetPropW(hwnd HWND, key *WCHAR, data HANDLE) error {
	return sysutil.MustTrue(lzSetPropW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(key)), uintptr(data)))
}

var lzRemovePropW = lzUser32.NewProc("RemovePropW")

func RemovePropW(hwnd HWND, key *WCHAR) (HANDLE, error) {
	return sysutil.MustNotZero[HANDLE](lzRemovePropW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(key))))
}

var lzGetPropW = lzUser32.NewProc("GetPropW")

func GetPropW(hwnd HWND, key *WCHAR) (HANDLE, error) {
	r, _, err := lzGetPropW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(key)))
	if r != 0 || sysutil.IsNoError(err) {
		err = nil
	}
	return HANDLE(r), err
}

var lzCallWindowProcW = lzUser32.NewProc("CallWindowProcW")

func CallWindowProcW(proc uintptr, hwnd HWND, msg UINT, wParam WPARAM, lParam LPARAM) LRESULT {
	return sysutil.As[LRESULT](lzCallWindowProcW.Call(uintptr(proc), uintptr(hwnd), uintptr(msg), uintptr(wParam), uintptr(lParam)))
}

var lzSendMessageW = lzUser32.NewProc("SendMessageW")

func SendMessageW(hwnd HWND, message UINT, wParam WPARAM, lParam LPARAM) (LRESULT, error) {
	r, _, err := lzSendMessageW.Call(uintptr(hwnd), uintptr(message), uintptr(wParam), uintptr(lParam))
	if sysutil.IsNoError(err) {
		err = nil
	}
	return LRESULT(r), err
}

var lzPostMessageW = lzUser32.NewProc("PostMessageW")

func PostMessageW(hwnd HWND, message UINT, wParam WPARAM, lParam LPARAM) error {
	return sysutil.MustTrue(lzPostMessageW.Call(uintptr(hwnd), uintptr(message), uintptr(wParam), uintptr(lParam)))
}

var lzPostThreadMessageW = lzUser32.NewProc("PostThreadMessageW")

func PostThreadMessageW(threadId DWORD, msg UINT, wParam WPARAM, lParam LPARAM) error {
	return sysutil.MustTrue(lzPostThreadMessageW.Call(uintptr(threadId), uintptr(msg), uintptr(wParam), uintptr(lParam)))
}

const (
	IMAGE_BITMAP = 0
	IMAGE_ICON   = 1
	IMAGE_CURSOR = 2
)

const (
	LR_DEFAULTCOLOR     = 0x00000000
	LR_MONOCHROME       = 0x00000001
	LR_COLOR            = 0x00000002
	LR_COPYRETURNORG    = 0x00000004
	LR_COPYDELETEORG    = 0x00000008
	LR_LOADFROMFILE     = 0x00000010
	LR_LOADTRANSPARENT  = 0x00000020
	LR_DEFAULTSIZE      = 0x00000040
	LR_VGACOLOR         = 0x00000080
	LR_LOADMAP3DCOLORS  = 0x00001000
	LR_CREATEDIBSECTION = 0x00002000
	LR_COPYFROMRESOURCE = 0x00004000
	LR_SHARED           = 0x00008000
)

const (
	IDC_ARROW       = 32512
	IDC_IBEAM       = 32513
	IDC_WAIT        = 32514
	IDC_CROSS       = 32515
	IDC_UPARROW     = 32516
	IDC_SIZENWSE    = 32642
	IDC_SIZENESW    = 32643
	IDC_SIZEWE      = 32644
	IDC_SIZENS      = 32645
	IDC_SIZEALL     = 32646
	IDC_NO          = 32648
	IDC_HAND        = 32649
	IDC_APPSTARTING = 32650
	IDC_HELP        = 32651
	IDC_PIN         = 32671
	IDC_PERSON      = 32672
)

const (
	OCR_NORMAL      = 32512
	OCR_IBEAM       = 32513
	OCR_WAIT        = 32514
	OCR_CROSS       = 32515
	OCR_UP          = 32516
	OCR_SIZENWSE    = 32642
	OCR_SIZENESW    = 32643
	OCR_SIZEWE      = 32644
	OCR_SIZENS      = 32645
	OCR_SIZEALL     = 32646
	OCR_NO          = 32648
	OCR_HAND        = 32649
	OCR_APPSTARTING = 32650
)

var lzLoadImageW = lzUser32.NewProc("LoadImageW")

func LoadImageW[H HBITMAP | HCURSOR | HICON](instance HINSTANCE, name *WCHAR, imageType UINT, cx INT, cy INT, flag UINT) (H, error) {
	return sysutil.MustNotZero[H](lzLoadImageW.Call(uintptr(instance), uintptr(unsafe.Pointer(name)), uintptr(imageType), uintptr(cx), uintptr(cy), uintptr(flag)))
}

func LoadImageW_uintptr[H HBITMAP | HCURSOR | HICON](instance HINSTANCE, name uintptr, imageType UINT, cx INT, cy INT, flag UINT) (H, error) {
	return sysutil.MustNotZero[H](lzLoadImageW.Call(uintptr(instance), name, uintptr(imageType), uintptr(cx), uintptr(cy), uintptr(flag)))
}

var lzSetWindowTextW = lzUser32.NewProc("SetWindowTextW")

func SetWindowTextW(hwnd HWND, str *WCHAR) error {
	return sysutil.MustTrue(lzSetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(str))))
}

var lzGetWindowTextLengthW = lzUser32.NewProc("GetWindowTextLengthW")

func GetWindowTextLengthW(hwnd HWND) (int, error) {
	r, _, err := lzGetWindowTextLengthW.Call(uintptr(hwnd))
	if r != 0 || sysutil.IsNoError(err) {
		err = nil
	}
	return int(r), err
}

var lzGetWindowTextW = lzUser32.NewProc("GetWindowTextW")

func GetWindowTextW(hwnd HWND, buffer *WCHAR, maxCount int) (int, error) {
	r, _, err := lzGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(buffer)), uintptr(maxCount))
	if r != 0 || sysutil.IsNoError(err) {
		err = nil
	}
	return int(r), err
}

var lzCreateMenu = lzUser32.NewProc("CreateMenu")

func CreateMenu() (HMENU, error) {
	return sysutil.MustNotZero[HMENU](lzCreateMenu.Call())
}

var lzCreatePopupMenu = lzUser32.NewProc("CreatePopupMenu")

func CreatePopupMenu() (HMENU, error) {
	return sysutil.MustNotZero[HMENU](lzCreatePopupMenu.Call())
}

var lzDestroyMenu = lzUser32.NewProc("DestroyMenu")

func DestroyMenu(menu HMENU) error {
	return sysutil.MustTrue(lzDestroyMenu.Call(uintptr(menu)))
}

const (
	MF_BYCOMMAND  = 0x0000000
	MF_BYPOSITION = 0x00000400
)

var lzDeleteMenu = lzUser32.NewProc("DeleteMenu")

func DeleteMenu(menu HMENU, pos UINT, flags UINT) error {
	return sysutil.MustTrue(lzDeleteMenu.Call(uintptr(menu), uintptr(pos), uintptr(flags)))
}

var lzRemoveMenu = lzUser32.NewProc("RemoveMenu")

func RemoveMenu(menu HMENU, pos UINT, flags UINT) error {
	return sysutil.MustTrue(lzRemoveMenu.Call(uintptr(menu), uintptr(pos), uintptr(flags)))
}

const (
	MIIM_BITMAP     = 0x00000080
	MIIM_CHECKMARKS = 0x00000008
	MIIM_DATA       = 0x00000020
	MIIM_FTYPE      = 0x00000100
	MIIM_ID         = 0x00000002
	MIIM_STATE      = 0x00000001
	MIIM_STRING     = 0x00000040
	MIIM_SUBMENU    = 0x00000004
	MIIM_TYPE       = 0x00000010
)

const (
	MFT_BITMAP       = 0x00000004
	MFT_MENUBARBREAK = 0x00000020
	MFT_MENUBREAK    = 0x00000040
	MFT_OWNERDRAW    = 0x00000100
	MFT_RADIOCHECK   = 0x00000200
	MFT_RIGHTJUSTIFY = 0x00004000
	MFT_RIGHTORDER   = 0x00002000
	MFT_SEPARATOR    = 0x00000800
	MFT_STRING       = 0x00000000
)

const (
	MFS_CHECKED   = 0x00000008
	MFS_DEFAULT   = 0x00001000
	MFS_DISABLED  = 0x00000003
	MFS_ENABLED   = 0x00000000
	MFS_GRAYED    = 0x00000003
	MFS_HILITE    = 0x00000080
	MFS_UNCHECKED = 0x00000000
	MFS_UNHILITE  = 0x00000000
)

type MENUITEMINFOW struct {
	Size            UINT
	Mask            UINT
	Type            UINT
	State           UINT
	ID              UINT
	SubMenu         HMENU
	CheckedBitmap   HBITMAP
	UncheckedBitmap HBITMAP
	ItemData        ULONG_PTR
	TypeData        *WCHAR // If TypeData needs to be a non pointer, a new struct and a new version of InsertMenuItemW must be defined instead of conversion to pointer.
	Cch             UINT
	ItemBitmap      HBITMAP
}

var lzInsertMenuItemW = lzUser32.NewProc("InsertMenuItemW")

func InsertMenuItemW(menu HMENU, item UINT, byPos bool, mii *MENUITEMINFOW) error {
	var byPosInt BOOL = 0
	if byPos {
		byPosInt = 1
	}
	return sysutil.MustTrue(lzInsertMenuItemW.Call(uintptr(menu), uintptr(item), uintptr(byPosInt), uintptr(unsafe.Pointer(mii))))
}

var lzGetMenuItemCount = lzUser32.NewProc("GetMenuItemCount")

func GetMenuItemCount(menu HMENU) (INT, error) {
	return sysutil.MustNotNegativeOne[INT](lzGetMenuItemCount.Call(uintptr(menu)))
}

var lzGetMenuItemInfoW = lzUser32.NewProc("GetMenuItemInfoW")

func GetMenuItemInfoW(menu HMENU, item UINT, byPos bool, mii *MENUITEMINFOW) error {
	return sysutil.MustTrue(lzGetMenuItemInfoW.Call(uintptr(menu), uintptr(item), uintptr(gg.If(byPos, 1, 0)), uintptr(unsafe.Pointer(mii))))
}

var lzSetMenuItemInfoW = lzUser32.NewProc("SetMenuItemInfoW")

func SetMenuItemInfoW(menu HMENU, item UINT, byPos bool, mmi *MENUITEMINFOW) error {
	return sysutil.MustTrue(lzSetMenuItemInfoW.Call(uintptr(menu), uintptr(item), gg.If[uintptr](byPos, 1, 0), uintptr(unsafe.Pointer(mmi))))
}

var lzSetMenu = lzUser32.NewProc("SetMenu")

func SetMenu(hwnd HWND, menu HMENU) error {
	return sysutil.MustTrue(lzSetMenu.Call(uintptr(hwnd), uintptr(menu)))
}

func HIWORD[T ~uintptr](l T) WORD {
	return WORD((l >> 16) & 0xFFFF)
}

func LOWORD[T ~uintptr](l T) WORD {
	return WORD(l & 0xFFFF)
}

func MAKEWORD[T ~byte](a, b T) WORD {
	return WORD(a) | WORD(b)<<8
}

func MAKELONG[T ~uint16](a, b T) LONG {
	return LONG(uint32(a) | uint32(b)<<16)
}

func LOBYTE[T ~uint16](w T) BYTE {
	return BYTE(w & 0xff)
}

func HIBYTE[T ~uint16](w T) BYTE {
	return BYTE((w >> 8) & 0xff)
}

type ACCEL_FVIRT BYTE

const (
	FALT      ACCEL_FVIRT = 0x10
	FCONTROL  ACCEL_FVIRT = 0x08
	FNOINVERT ACCEL_FVIRT = 0x02
	FSHIFT    ACCEL_FVIRT = 0x04
	FVIRTKEY  ACCEL_FVIRT = 1
)

type ACCEL struct {
	Virt ACCEL_FVIRT
	Key  WORD
	Cmd  WORD
}

// alignSlice makes &s[0] word-aligned.
// If len(s) == 0, s is unchanged.
func alignSlice[T any](s *[]T) {
	const WordSize = unsafe.Sizeof(uintptr(0))
	if len(*s) == 0 || uintptr(unsafe.Pointer(&(*s)[0]))%WordSize == 0 {
		return
	}
	// p = alloc(slice_data_size+WordSize)
	p := make([]byte, len(*s)*int(unsafe.Sizeof((*s)[0]))+int(WordSize))
	// aligned = p+WordSize-(p%WordSize)
	aligned := unsafe.Slice((*T)(unsafe.Pointer(uintptr(unsafe.Pointer(&p[0]))+WordSize-uintptr(unsafe.Pointer(&p[0]))%WordSize)), len(*s))
	copy(aligned, *s)
	*s = aligned
}

var lzCreateAcceleratorTableW = lzUser32.NewProc("CreateAcceleratorTableW")

func CreateAcceleratorTableW(accel []ACCEL) (HACCEL, error) {
	// For some reason, &accel[0] may not be aligned.
	alignSlice(&accel)
	r, r2, err := lzCreateAcceleratorTableW.Call(uintptr(unsafe.Pointer(&accel[0])), uintptr(len(accel)))
	return sysutil.MustNotZero[HACCEL](r, r2, err)
}

var lzDestroyAcceleratorTable = lzUser32.NewProc("DestroyAcceleratorTable")

func DestroyAcceleratorTable(table HACCEL) error {
	return sysutil.MustTrue(lzDestroyAcceleratorTable.Call(uintptr(table)))
}

var lzTranslateAcceleratorW = lzUser32.NewProc("TranslateAcceleratorW")

func TranslateAcceleratorW(hwnd HWND, accTable HACCEL, msg *MSG) (bool, error) {
	r, _, err := lzTranslateAcceleratorW.Call(uintptr(hwnd), uintptr(accTable), uintptr(unsafe.Pointer(msg)))
	if sysutil.IsNoError(err) {
		err = nil
	}
	return r != 0, err
}

var lzGetActiveWindow = lzUser32.NewProc("GetActiveWindow")

func GetActiveWindow() HWND {
	r, _, _ := lzGetActiveWindow.Call()
	return HWND(r)
}

type TRACK_POPUP_MENU_FLAG UINT

const (
	TPM_CENTERALIGN     TRACK_POPUP_MENU_FLAG = 0x0004
	TPM_LEFTALIGN       TRACK_POPUP_MENU_FLAG = 0x0000
	TPM_RIGHTALIGN      TRACK_POPUP_MENU_FLAG = 0x0008
	TPM_BOTTOMALIGN     TRACK_POPUP_MENU_FLAG = 0x0020
	TPM_TOPALIGN        TRACK_POPUP_MENU_FLAG = 0x0000
	TPM_VCENTERALIGN    TRACK_POPUP_MENU_FLAG = 0x001
	TPM_NONOTIFY        TRACK_POPUP_MENU_FLAG = 0x0080
	TPM_RETURNCMD       TRACK_POPUP_MENU_FLAG = 0x0100
	TPM_LEFTBUTTON      TRACK_POPUP_MENU_FLAG = 0x0000
	TPM_RIGHTBUTTON     TRACK_POPUP_MENU_FLAG = 0x0002
	TPM_RECURSE         TRACK_POPUP_MENU_FLAG = 0x0001
	TPM_HORNEGANIMATION TRACK_POPUP_MENU_FLAG = 0x0800
	TPM_HORPOSANIMATION TRACK_POPUP_MENU_FLAG = 0x0400
	TPM_NOANIMATION     TRACK_POPUP_MENU_FLAG = 0x4000
	TPM_VERNEGANIMATION TRACK_POPUP_MENU_FLAG = 0x2000
	TPM_VERPOSANIMATION TRACK_POPUP_MENU_FLAG = 0x1000
	TPM_HORIZONTAL      TRACK_POPUP_MENU_FLAG = 0x0000
	TPM_VERTICAL        TRACK_POPUP_MENU_FLAG = 0x0040
)

type POINT struct {
	X, Y LONG
}

type RECT struct {
	Left, Top, Right, Bottom LONG
}

func (rect *RECT) Width() LONG {
	return rect.Right - rect.Left
}

func (rect *RECT) Height() LONG {
	return rect.Bottom - rect.Top
}

func (rect *RECT) TopLeft() *POINT {
	return (*POINT)(unsafe.Pointer(&rect.Left))
}

func (rect *RECT) BottomRight() *POINT {
	return (*POINT)(unsafe.Pointer(&rect.Right))
}

type TPMPARAMS struct {
	Size    UINT
	Exclude RECT
}

var lzTrackPopupMenuEx = lzUser32.NewProc("TrackPopupMenuEx")

func TrackPopupMenuEx(menu HMENU, flags TRACK_POPUP_MENU_FLAG, x INT, y INT, hwnd HWND, params *TPMPARAMS) (int, error) {
	r, _, err := lzTrackPopupMenuEx.Call(uintptr(menu), uintptr(flags), uintptr(x), uintptr(y), uintptr(hwnd), uintptr(unsafe.Pointer(params)))
	if !sysutil.IsNoError(err) {
		return 0, err
	}
	return int(r), nil
}

var lzGetCursorPos = lzUser32.NewProc("GetCursorPos")

const (
	IDOK       = 1
	IDCANCEL   = 2
	IDABORT    = 3
	IDRETRY    = 4
	IDIGNORE   = 5
	IDYES      = 6
	IDNO       = 7
	IDCLOSE    = 8
	IDHELP     = 9
	IDTRYAGAIN = 10
	IDCONTINUE = 11
	IDTIMEOUT  = 32000
)

const (
	BN_CLICKED       = 0
	BN_PAINT         = 1
	BN_HILITE        = 2
	BN_UNHILITE      = 3
	BN_DISABLE       = 4
	BN_DOUBLECLICKED = 5
	BN_PUSHED        = BN_HILITE
	BN_UNPUSHED      = BN_UNHILITE
	BN_DBLCLK        = BN_DOUBLECLICKED
	BN_SETFOCUS      = 6
	BN_KILLFOCUS     = 7
)

func GetCursorPos() (*POINT, error) {
	var pos POINT
	r, _, err := lzGetCursorPos.Call(uintptr(unsafe.Pointer(&pos)))
	if r == 0 {
		return nil, err
	}
	return &pos, nil
}

type DLGTEMPLATE struct {
	Style        DWORD
	ExStyle      DWORD
	ItemCount    WORD
	X, Y, CX, CY SHORT
}

var lzDialogBoxIndirectParamW = lzUser32.NewProc("DialogBoxIndirectParamW")

func DialogBoxIndirectParamW(instance HINSTANCE, template *DLGTEMPLATE, parent HWND, dialogFunc uintptr, param LPARAM) (UINT_PTR, error) {
	r, _, err := lzDialogBoxIndirectParamW.Call(uintptr(instance), uintptr(unsafe.Pointer(template)), uintptr(parent), dialogFunc, uintptr(param))
	if sysutil.IsNoError(err) {
		err = nil
	}
	return UINT_PTR(r), err
}

var lzGetDialogBaseUnits = lzUser32.NewProc("GetDialogBaseUnits")

func GetDialogBaseUnits() LONG {
	r, _, _ := lzGetDialogBaseUnits.Call()
	return LONG(r)
}

var lzMulDiv = lzKernel32.NewProc("MulDiv")

func MulDiv(number INT, numerator INT, denominator INT) INT {
	r, _, _ := lzMulDiv.Call(uintptr(number), uintptr(numerator), uintptr(denominator))
	return INT(r)
}

var lzEndDialog = lzUser32.NewProc("EndDialog")

func EndDialog(hwnd HWND, result INT_PTR) error {
	return sysutil.MustTrue(lzEndDialog.Call(uintptr(hwnd), uintptr(result)))
}

var lzGetParent = lzUser32.NewProc("GetParent")

func GetParent(hwnd HWND) (HWND, error) {
	return sysutil.MustNoError[HWND](lzGetParent.Call(uintptr(hwnd)))
}

type GET_ANCESTOR_FLAG UINT

const (
	GA_PARENT    GET_ANCESTOR_FLAG = 1
	GA_ROOT      GET_ANCESTOR_FLAG = 2
	GA_ROOTOWNER GET_ANCESTOR_FLAG = 3
)

var lzGetAncestor = lzUser32.NewProc("GetAncestor")

func GetAncestor(hwnd HWND, flags GET_ANCESTOR_FLAG) (HWND, error) {
	return sysutil.MustNoError[HWND](lzGetAncestor.Call(uintptr(hwnd), uintptr(flags)))
}

type GET_WINDOW_CMD UINT

const (
	GW_CHILD        GET_WINDOW_CMD = 5
	GW_ENABLEDPOPUP GET_WINDOW_CMD = 6
	GW_HWNDFIRST    GET_WINDOW_CMD = 0
	GW_HWNDLAST     GET_WINDOW_CMD = 1
	GW_HWNDNEXT     GET_WINDOW_CMD = 2
	GW_HWNDPREV     GET_WINDOW_CMD = 3
	GW_OWNER        GET_WINDOW_CMD = 4
)

var lzGetWindow = lzUser32.NewProc("GetWindow")

func GetWindow(hwnd HWND, cmd GET_WINDOW_CMD) (HWND, error) {
	return sysutil.MustNoError[HWND](lzGetWindow.Call(uintptr(hwnd), uintptr(cmd)))
}

var lzGetDlgItem = lzUser32.NewProc("GetDlgItem")

func GetDlgItem(hwnd HWND, id INT) (HWND, error) {
	return sysutil.MustNotZero[HWND](lzGetDlgItem.Call(uintptr(hwnd), uintptr(id)))
}

const (
	BS_PUSHBUTTON      = 0x00000000
	BS_DEFPUSHBUTTON   = 0x00000001
	BS_CHECKBOX        = 0x00000002
	BS_AUTOCHECKBOX    = 0x00000003
	BS_RADIOBUTTON     = 0x00000004
	BS_3STATE          = 0x00000005
	BS_AUTO3STATE      = 0x00000006
	BS_GROUPBOX        = 0x00000007
	BS_USERBUTTON      = 0x00000008
	BS_AUTORADIOBUTTON = 0x00000009
	BS_PUSHBOX         = 0x0000000A
	BS_OWNERDRAW       = 0x0000000B
	BS_TYPEMASK        = 0x0000000F
	BS_LEFTTEXT        = 0x00000020
	BS_TEXT            = 0x00000000
	BS_ICON            = 0x00000040
	BS_BITMAP          = 0x00000080
	BS_LEFT            = 0x00000100
	BS_RIGHT           = 0x00000200
	BS_CENTER          = 0x00000300
	BS_TOP             = 0x00000400
	BS_BOTTOM          = 0x00000800
	BS_VCENTER         = 0x00000C00
	BS_PUSHLIKE        = 0x00001000
	BS_MULTILINE       = 0x00002000
	BS_NOTIFY          = 0x00004000
	BS_FLAT            = 0x00008000
	BS_RIGHTBUTTON     = BS_LEFTTEXT
)

const (
	RT_ICON       = 3
	RT_GROUP_ICON = RT_ICON + 11
	RT_MANIFEST   = 24
)

type HUPDATE HANDLE

var lzBeginUpdateResourceW = lzKernel32.NewProc("BeginUpdateResourceW")

func BeginUpdateResourceW(fileName *WCHAR, deleteExisting bool) (HUPDATE, error) {
	return sysutil.MustNotZero[HUPDATE](lzBeginUpdateResourceW.Call(uintptr(unsafe.Pointer(fileName)), gg.If[uintptr](deleteExisting, 1, 0)))
}

var lzEndUpdateResourceW = lzKernel32.NewProc("EndUpdateResourceW")

func EndUpdateResourceW(update HUPDATE, discard bool) error {
	r, _, err := lzEndUpdateResourceW.Call(uintptr(update), gg.If[uintptr](discard, 1, 0))
	if r != 0 && sysutil.IsNoError(err) {
		err = nil
	}
	return err
}

var lzUpdateResourceW = lzKernel32.NewProc("UpdateResourceW")

func UpdateResourceW(update HUPDATE, resType uintptr, name uintptr, lang WORD, data unsafe.Pointer, dataSize DWORD) error {
	r, _, err := lzUpdateResourceW.Call(uintptr(update), resType, name, uintptr(lang), uintptr(data), uintptr(dataSize))
	if r != 0 && sysutil.IsNoError(err) {
		err = nil
	}
	return err
}

var lzEnumResourceNamesExW = lzKernel32.NewProc("EnumResourceNamesExW")

func EnumResourceNamesExW(module HMODULE, typ uintptr, enumProc uintptr, lParam uintptr, flags DWORD, langId DWORD) error {
	return sysutil.MustTrue(lzEnumResourceNamesExW.Call(uintptr(module), typ, enumProc, lParam, uintptr(flags), uintptr(langId)))
}

var lzFindResourceW = lzKernel32.NewProc("FindResourceW")

func FindResourceW(module HMODULE, name uintptr, typ uintptr) (HRSRC, error) {
	return sysutil.MustNotZero[HRSRC](lzFindResourceW.Call(uintptr(module), name, typ))
}

var lzLoadResource = lzKernel32.NewProc("LoadResource")

func LoadResource(module HMODULE, res HRSRC) (HGLOBAL, error) {
	return sysutil.MustNotZero[HGLOBAL](lzLoadResource.Call(uintptr(module), uintptr(res)))
}

var lzLockResource = lzKernel32.NewProc("LockResource")

func LockResource(res HGLOBAL) (PVOID, error) {
	return sysutil.MustNotZero[PVOID](lzLockResource.Call(uintptr(res)))
}

type PAINTSTRUCT struct {
	HDC     HDC
	Erase   BOOL
	RcPaint RECT
	_       BOOL
	_       BOOL
	_       [32]byte
}

var lzBeginPaint = lzUser32.NewProc("BeginPaint")

func BeginPaint(hwnd HWND, p *PAINTSTRUCT) (HDC, error) {
	return sysutil.MustNotZero[HDC](lzBeginPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(p))))
}

var lzEndPaint = lzUser32.NewProc("EndPaint")

func EndPaint(hwnd HWND, p *PAINTSTRUCT) error {
	return sysutil.MustTrue(lzEndPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(p))))
}

var lzDrawFocusRect = lzUser32.NewProc("DrawFocusRect")

func DrawFocusRect(hdc HDC, rect *RECT) error {
	return sysutil.MustTrue(lzDrawFocusRect.Call(uintptr(hdc), uintptr(unsafe.Pointer(rect))))
}

var lzCreateCompatibleDC = lzGdi32.NewProc("CreateCompatibleDC")

func CreateCompatibleDC(hdc HDC) (HDC, error) {
	return sysutil.MustNotZero[HDC](lzCreateCompatibleDC.Call(uintptr(hdc)))
}

var lzGetDC = lzUser32.NewProc("GetDC")

func GetDC(hwnd HWND) (HDC, error) {
	return sysutil.MustNotZero[HDC](lzGetDC.Call(uintptr(hwnd)))
}

var lzReleaseDC = lzUser32.NewProc("ReleaseDC")

func ReleaseDC(hwnd HWND, hdc HDC) bool {
	return sysutil.AsBool(lzReleaseDC.Call(uintptr(hwnd), uintptr(hdc)))
}

var lzCreateCompatibleBitmap = lzGdi32.NewProc("CreateCompatibleBitmap")

func CreateCompatibleBitmap(hdc HDC, cx INT, cy INT) (HBITMAP, error) {
	return sysutil.MustNotZero[HBITMAP](lzCreateCompatibleBitmap.Call(uintptr(hdc), uintptr(cx), uintptr(cy)))
}

const (
	SRCCOPY     = 0x00CC0020
	SRCPAINT    = 0x00EE0086
	SRCAND      = 0x008800C6
	SRCINVERT   = 0x00660046
	SRCERASE    = 0x00440328
	NOTSRCCOPY  = 0x00330008
	NOTSRCERASE = 0x001100A6
	MERGECOPY   = 0x00C000CA
	MERGEPAINT  = 0x00BB0226
	PATCOPY     = 0x00F00021
	PATPAINT    = 0x00FB0A09
	PATINVERT   = 0x005A0049
	DSTINVERT   = 0x00550009
	BLACKNESS   = 0x00000042
	WHITENESS   = 0x00FF0062
)

var lzBitBlt = lzGdi32.NewProc("BitBlt")

func BitBlt(hdc HDC, x int, y int, cx int, cy int, srcDC HDC, srcX int, srcY int, op DWORD) error {
	return sysutil.MustTrue(lzBitBlt.Call(uintptr(hdc), uintptr(x), uintptr(y), uintptr(cx), uintptr(cy), uintptr(srcDC), uintptr(srcX), uintptr(srcY), uintptr(op)))
}

var lzDeleteObject = lzGdi32.NewProc("DeleteObject")

func DeleteObject[H HGDIOBJ](h H) error {
	return sysutil.MustTrue(lzDeleteObject.Call(uintptr(h)))
}

const (
	GDI_ERROR = 0xFFFFFFFF
)

var lzSelectObject = lzGdi32.NewProc("SelectObject")

func SelectObject[H HGDIOBJ](hdc HDC, obj H) (H, error) {
	h, _, err := lzSelectObject.Call(uintptr(hdc), uintptr(obj))
	if h == 0 || h == GDI_ERROR {
		return 0, err
	}
	return H(h), nil
}

var lzRectangle = lzGdi32.NewProc("Rectangle")

func Rectangle(hdc HDC, left int, top int, right int, bottom int) error {
	return sysutil.MustTrue(lzRectangle.Call(uintptr(hdc), uintptr(left), uintptr(top), uintptr(right), uintptr(bottom)))
}

var lzFillRect = lzUser32.NewProc("FillRect")

func FillRect(hdc HDC, rect *RECT, brush HBRUSH) error {
	return sysutil.MustTrue(lzFillRect.Call(uintptr(hdc), uintptr(unsafe.Pointer(rect)), uintptr(brush)))
}

const (
	LF_FACESIZE = 32
)

type LOGFONTW struct {
	Height         LONG
	Width          LONG
	Escapement     LONG
	Orientation    LONG
	Weight         LONG
	Italic         BYTE
	Underline      BYTE
	StrikeOut      BYTE
	CharSet        BYTE
	OutPrecision   BYTE
	ClipPrecision  BYTE
	Quality        BYTE
	PitchAndFamily BYTE
	FaceName       [LF_FACESIZE]WCHAR
}

type LOGPEN struct {
	Style PEN_STYLE
	Width LONG
	_     LONG
	Color COLORREF
}

type NONCLIENTMETRICSW struct {
	Size              UINT
	BorderWidth       INT
	ScrollWidth       INT
	ScrollHeight      INT
	CaptionWidth      INT
	CaptionHeight     INT
	CaptionFont       LOGFONTW
	SmCaptionWidth    INT
	SmCaptionHeight   INT
	SmCaptionFont     LOGFONTW
	MenuWidth         INT
	MenuHeight        INT
	MenuFont          LOGFONTW
	StatusFont        LOGFONTW
	MessageFont       LOGFONTW
	PaddedBorderWidth INT
}

const (
	SPI_GETNONCLIENTMETRICS = 0x0029
)

var lzSystemParametersInfoW = lzUser32.NewProc("SystemParametersInfoW")

func SystemParametersInfoW(action UINT, param UINT, p PVOID, winIni UINT) error {
	return sysutil.MustTrue(lzSystemParametersInfoW.Call(uintptr(action), uintptr(param), uintptr(p), uintptr(winIni)))
}

var lzSystemParametersInfoForDpi = lzUser32.NewProc("SystemParametersInfoForDpi")

func SystemParametersInfoForDpi(action UINT, param UINT, p PVOID, winIni UINT, dpi UINT) error {
	return sysutil.MustTrue(lzSystemParametersInfoForDpi.Call(uintptr(action), uintptr(param), uintptr(unsafe.Pointer(p)), uintptr(winIni), uintptr(dpi)))
}

var lzCreateFontIndirectW = lzGdi32.NewProc("CreateFontIndirectW")

func CreateFontIndirectW(f *LOGFONTW) (HFONT, error) {
	return sysutil.MustNotZero[HFONT](lzCreateFontIndirectW.Call(uintptr(unsafe.Pointer(f))))
}

var lzCreatePenIndirect = lzGdi32.NewProc("CreatePenIndirect")

func CreatePenIndirect(p *LOGPEN) (HPEN, error) {
	return sysutil.MustNotZero[HPEN](lzCreatePenIndirect.Call(uintptr(unsafe.Pointer(p))))
}

const (
	CLR_INVALID COLORREF = 0xFFFFFFFF
)

var lzSetTextColor = lzGdi32.NewProc("SetTextColor")

func SetTextColor(hdc HDC, color COLORREF) (COLORREF, error) {
	r, _, err := lzSetTextColor.Call(uintptr(hdc), uintptr(color))
	if r == uintptr(CLR_INVALID) {
		return COLORREF(r), err
	}
	return COLORREF(r), nil
}

const (
	HWND_BOTTOM    HWND = 1
	HWND_NOTOPMOST HWND = 2
	HWND_TOP       HWND = 0
	HWND_TOPMOST   HWND = ^HWND(0)
)

const (
	SWP_ASYNCWINDOWPOS = 0x4000
	SWP_DEFERERASE     = 0x2000
	SWP_DRAWFRAME      = 0x0020
	SWP_FRAMECHANGED   = 0x0020
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOACTIVATE     = 0x0010
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOMOVE         = 0x0002
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOREDRAW       = 0x0008
	SWP_NOREPOSITION   = 0x0200
	SWP_NOSENDCHANGING = 0x0400
	SWP_NOSIZE         = 0x0001
	SWP_NOZORDER       = 0x0004
	SWP_SHOWWINDOW     = 0x0040
)

var lzEnableWindow = lzUser32.NewProc("EnableWindow")

func EnableWindow(hwnd HWND, enable bool) bool {
	return sysutil.AsBool(lzEnableWindow.Call(uintptr(hwnd), gg.If[uintptr](enable, 1, 0)))
}

var lzIsWindowEnabled = lzUser32.NewProc("IsWindowEnabled")

func IsWindowEnabled(hwnd HWND) bool {
	return sysutil.AsBool(lzIsWindowEnabled.Call(uintptr(hwnd)))
}

var lzSetWindowPos = lzUser32.NewProc("SetWindowPos")

func SetWindowPos(hwnd HWND, hwndInsertAfter HWND, x INT, y INT, cx INT, cy INT, flags UINT) error {
	return sysutil.MustTrue(lzSetWindowPos.Call(uintptr(hwnd), uintptr(hwndInsertAfter), uintptr(x), uintptr(y), uintptr(cx), uintptr(cy), uintptr(flags)))
}

var lzGetDpiForWindow = lzUser32.NewProc("GetDpiForWindow")

func GetDpiForWindow(hwnd HWND) (UINT, error) {
	return sysutil.MustNoError[UINT](lzGetDpiForWindow.Call(uintptr(hwnd)))
}

var lzGetDpiForSystem = lzUser32.NewProc("GetDpiForSystem")

func GetDpiForSystem() UINT {
	r, _, _ := lzGetDpiForSystem.Call()
	return UINT(r)
}

const USER_DEFAULT_SCREEN_DPI = 96

var lzGetDesktopWindow = lzUser32.NewProc("GetDesktopWindow")

func GetDesktopWindow() HWND {
	r, _, _ := lzGetDesktopWindow.Call()
	return HWND(r)
}

var lzGetWindowRect = lzUser32.NewProc("GetWindowRect")

func GetWindowRect(hwnd HWND, rect *RECT) error {
	return sysutil.MustTrue(lzGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(rect))))
}

var lzScreenToClient = lzUser32.NewProc("ScreenToClient")

func ScreenToClient(hwnd HWND, pt *POINT) error {
	return sysutil.MustTrue(lzScreenToClient.Call(uintptr(hwnd), uintptr(unsafe.Pointer(pt))))
}

var lzClientToScreen = lzUser32.NewProc("ClientToScreen")

func ClientToScreen(hwnd HWND, pt *POINT) error {
	return sysutil.MustTrue(lzClientToScreen.Call(uintptr(hwnd), uintptr(unsafe.Pointer(pt))))
}

var lzGetClientRect = lzUser32.NewProc("GetClientRect")

func GetClientRect(hwnd HWND, rect *RECT) error {
	return sysutil.MustTrue(lzGetClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(rect))))
}

var lzLineTo = lzGdi32.NewProc("LineTo")

func LineTo(hdc HDC, x INT, y INT) error {
	return sysutil.MustTrue(lzLineTo.Call(uintptr(hdc), uintptr(x), uintptr(y)))
}

var lzMoveToEx = lzGdi32.NewProc("MoveToEx")

func MoveToEx(hdc HDC, x INT, y INT, prev *POINT) error {
	return sysutil.MustTrue(lzMoveToEx.Call(uintptr(hdc), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(prev))))
}

type MESSAGE_BOX_TYPE UINT

const (
	MB_OK                        MESSAGE_BOX_TYPE = 0x00000000
	MB_OKCANCEL                  MESSAGE_BOX_TYPE = 0x00000001
	MB_ABORTRETRYIGNORE          MESSAGE_BOX_TYPE = 0x00000002
	MB_YESNOCANCEL               MESSAGE_BOX_TYPE = 0x00000003
	MB_YESNO                     MESSAGE_BOX_TYPE = 0x00000004
	MB_RETRYCANCEL               MESSAGE_BOX_TYPE = 0x00000005
	MB_CANCELTRYCONTINUE         MESSAGE_BOX_TYPE = 0x00000006
	MB_ICONHAND                  MESSAGE_BOX_TYPE = 0x00000010
	MB_ICONQUESTION              MESSAGE_BOX_TYPE = 0x00000020
	MB_ICONEXCLAMATION           MESSAGE_BOX_TYPE = 0x00000030
	MB_ICONASTERISK              MESSAGE_BOX_TYPE = 0x00000040
	MB_USERICON                  MESSAGE_BOX_TYPE = 0x00000080
	MB_ICONWARNING               MESSAGE_BOX_TYPE = MB_ICONEXCLAMATION
	MB_ICONERROR                 MESSAGE_BOX_TYPE = MB_ICONHAND
	MB_ICONINFORMATION           MESSAGE_BOX_TYPE = MB_ICONASTERISK
	MB_ICONSTOP                  MESSAGE_BOX_TYPE = MB_ICONHAND
	MB_DEFBUTTON1                MESSAGE_BOX_TYPE = 0x00000000
	MB_DEFBUTTON2                MESSAGE_BOX_TYPE = 0x00000100
	MB_DEFBUTTON3                MESSAGE_BOX_TYPE = 0x00000200
	MB_DEFBUTTON4                MESSAGE_BOX_TYPE = 0x00000300
	MB_APPLMODAL                 MESSAGE_BOX_TYPE = 0x00000000
	MB_SYSTEMMODAL               MESSAGE_BOX_TYPE = 0x00001000
	MB_TASKMODAL                 MESSAGE_BOX_TYPE = 0x00002000
	MB_HELP                      MESSAGE_BOX_TYPE = 0x00004000
	MB_NOFOCUS                   MESSAGE_BOX_TYPE = 0x00008000
	MB_SETFOREGROUND             MESSAGE_BOX_TYPE = 0x00010000
	MB_DEFAULT_DESKTOP_ONLY      MESSAGE_BOX_TYPE = 0x00020000
	MB_TOPMOST                   MESSAGE_BOX_TYPE = 0x00040000
	MB_RIGHT                     MESSAGE_BOX_TYPE = 0x00080000
	MB_RTLREADING                MESSAGE_BOX_TYPE = 0x00100000
	MB_SERVICE_NOTIFICATION      MESSAGE_BOX_TYPE = 0x00200000
	MB_SERVICE_NOTIFICATION_NT3X MESSAGE_BOX_TYPE = 0x00040000
	MB_TYPEMASK                  MESSAGE_BOX_TYPE = 0x0000000F
	MB_ICONMASK                  MESSAGE_BOX_TYPE = 0x000000F0
	MB_DEFMASK                   MESSAGE_BOX_TYPE = 0x00000F00
	MB_MODEMASK                  MESSAGE_BOX_TYPE = 0x00003000
	MB_MISCMASK                  MESSAGE_BOX_TYPE = 0x0000C000
)

var lzMessageBoxExW = lzUser32.NewProc("MessageBoxExW")

func MessageBoxExW(owner HWND, text *WCHAR, caption *WCHAR, typ MESSAGE_BOX_TYPE, lang WORD) (INT, error) {
	return sysutil.MustNotZero[INT](lzMessageBoxExW.Call(uintptr(owner), uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(caption)), uintptr(typ), uintptr(lang)))
}

type PEN_STYLE DWORD

const (
	PS_SOLID       PEN_STYLE = 0
	PS_DASH        PEN_STYLE = 1 /* -------  */
	PS_DOT         PEN_STYLE = 2 /* .......  */
	PS_DASHDOT     PEN_STYLE = 3 /* _._._._  */
	PS_DASHDOTDOT  PEN_STYLE = 4 /* _.._.._  */
	PS_NULL        PEN_STYLE = 5
	PS_INSIDEFRAME PEN_STYLE = 6
	PS_USERSTYLE   PEN_STYLE = 7
	PS_ALTERNATE   PEN_STYLE = 8
	PS_STYLE_MASK  PEN_STYLE = 0x0000000F

	PS_ENDCAP_ROUND  PEN_STYLE = 0x00000000
	PS_ENDCAP_SQUARE PEN_STYLE = 0x00000100
	PS_ENDCAP_FLAT   PEN_STYLE = 0x00000200
	PS_ENDCAP_MASK   PEN_STYLE = 0x00000F00

	PS_JOIN_ROUND PEN_STYLE = 0x00000000
	PS_JOIN_BEVEL PEN_STYLE = 0x00001000
	PS_JOIN_MITER PEN_STYLE = 0x00002000
	PS_JOIN_MASK  PEN_STYLE = 0x0000F000

	PS_COSMETIC  PEN_STYLE = 0x00000000
	PS_GEOMETRIC PEN_STYLE = 0x00010000
	PS_TYPE_MASK PEN_STYLE = 0x000F0000
)

var lzExtCreatePen = lzGdi32.NewProc("ExtCreatePen")

func ExtCreatePen(style PEN_STYLE, width DWORD, brush *LOGBRUSH, userStyles []DWORD) (HPEN, error) {
	var (
		r1, r2 uintptr
	)
	var err error
	if len(userStyles) == 0 {
		r1, r2, err = lzExtCreatePen.Call(uintptr(style), uintptr(width), uintptr(unsafe.Pointer(brush)), 0, 0)
	} else {
		r1, r2, err = lzExtCreatePen.Call(uintptr(style), uintptr(width), uintptr(unsafe.Pointer(brush)), uintptr(len(userStyles)), uintptr(unsafe.Pointer(&userStyles[0])))
	}

	return sysutil.MustNotZero[HPEN](r1, r2, err)
}

type BRUSH_STYLE UINT

const (
	BS_SOLID         BRUSH_STYLE = 0
	BS_NULL          BRUSH_STYLE = 1
	BS_HOLLOW        BRUSH_STYLE = BS_NULL
	BS_HATCHED       BRUSH_STYLE = 2
	BS_PATTERN       BRUSH_STYLE = 3
	BS_INDEXED       BRUSH_STYLE = 4
	BS_DIBPATTERN    BRUSH_STYLE = 5
	BS_DIBPATTERNPT  BRUSH_STYLE = 6
	BS_PATTERN8X8    BRUSH_STYLE = 7
	BS_DIBPATTERN8X8 BRUSH_STYLE = 8
	BS_MONOPATTERN   BRUSH_STYLE = 9
)

type HATCH_STYLE ULONG_PTR

const (
	HS_HORIZONTAL HATCH_STYLE = 0 /* ----- */
	HS_VERTICAL   HATCH_STYLE = 1 /* ||||| */
	HS_FDIAGONAL  HATCH_STYLE = 2 /* \\\\\ */
	HS_BDIAGONAL  HATCH_STYLE = 3 /* ///// */
	HS_CROSS      HATCH_STYLE = 4 /* +++++ */
	HS_DIAGCROSS  HATCH_STYLE = 5 /* xxxxx */
)

type LOGBRUSH struct {
	Style BRUSH_STYLE
	Color COLORREF
	Hatch HATCH_STYLE
}

var lzCreateBrushIndirect = lzGdi32.NewProc("CreateBrushIndirect")

func CreateBrushIndirect(p *LOGBRUSH) (HBRUSH, error) {
	return sysutil.MustNotZero[HBRUSH](lzCreateBrushIndirect.Call(uintptr(unsafe.Pointer(p))))
}

type EXTLOGPEN struct {
	PenStyle   PEN_STYLE
	Width      DWORD
	BrushStyle BRUSH_STYLE
	Color      COLORREF
	Hatch      HATCH_STYLE
	NumEntries DWORD
	StyleEntry []DWORD
}

type BK_MODE INT

const (
	TRANSPARENT BK_MODE = 1
	OPAQUE      BK_MODE = 2
)

var lzSetBkMode = lzGdi32.NewProc("SetBkMode")

func SetBkMode(hdc HDC, mode BK_MODE) (BK_MODE, error) {
	return sysutil.MustNotZero[BK_MODE](lzSetBkMode.Call(uintptr(hdc), uintptr(mode)))
}

type DRAWTEXTPARAMS struct {
	Size        UINT
	TabLength   INT
	LeftMargin  INT
	RightMargin INT
	LengthDrawn UINT
}

type DRAW_TEXT_FORMAT UINT

const (
	DT_TOP             DRAW_TEXT_FORMAT = 0x00000000
	DT_LEFT            DRAW_TEXT_FORMAT = 0x00000000
	DT_CENTER          DRAW_TEXT_FORMAT = 0x00000001
	DT_RIGHT           DRAW_TEXT_FORMAT = 0x00000002
	DT_VCENTER         DRAW_TEXT_FORMAT = 0x00000004
	DT_BOTTOM          DRAW_TEXT_FORMAT = 0x00000008
	DT_WORDBREAK       DRAW_TEXT_FORMAT = 0x00000010
	DT_SINGLELINE      DRAW_TEXT_FORMAT = 0x00000020
	DT_EXPANDTABS      DRAW_TEXT_FORMAT = 0x00000040
	DT_TABSTOP         DRAW_TEXT_FORMAT = 0x00000080
	DT_NOCLIP          DRAW_TEXT_FORMAT = 0x00000100
	DT_EXTERNALLEADING DRAW_TEXT_FORMAT = 0x00000200
	DT_CALCRECT        DRAW_TEXT_FORMAT = 0x00000400
	DT_NOPREFIX        DRAW_TEXT_FORMAT = 0x00000800
	DT_INTERNAL        DRAW_TEXT_FORMAT = 0x00001000

	DT_EDITCONTROL          DRAW_TEXT_FORMAT = 0x00002000
	DT_PATH_ELLIPSIS        DRAW_TEXT_FORMAT = 0x00004000
	DT_END_ELLIPSIS         DRAW_TEXT_FORMAT = 0x00008000
	DT_MODIFYSTRING         DRAW_TEXT_FORMAT = 0x00010000
	DT_RTLREADING           DRAW_TEXT_FORMAT = 0x00020000
	DT_WORD_ELLIPSIS        DRAW_TEXT_FORMAT = 0x00040000
	DT_NOFULLWIDTHCHARBREAK DRAW_TEXT_FORMAT = 0x00080000
	DT_HIDEPREFIX           DRAW_TEXT_FORMAT = 0x00100000
	DT_PREFIXONLY           DRAW_TEXT_FORMAT = 0x00200000
)

var lzDrawTextExW = lzUser32.NewProc("DrawTextExW")

func DrawTextExW(hdc HDC, text *WCHAR, cchText INT, rect *RECT, format DRAW_TEXT_FORMAT, param *DRAWTEXTPARAMS) (INT, error) {
	return sysutil.MustNotZero[INT](lzDrawTextExW.Call(uintptr(hdc), uintptr(unsafe.Pointer(text)), uintptr(cchText), uintptr(unsafe.Pointer(rect)), uintptr(format), uintptr(unsafe.Pointer(param))))
}

var lzInvalidateRect = lzUser32.NewProc("InvalidateRect")

func InvalidateRect(hwnd HWND, rect *RECT, erase bool) error {
	return sysutil.MustTrue(lzInvalidateRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(rect)), gg.If[uintptr](erase, 1, 0)))
}

var lzDeleteDC = lzGdi32.NewProc("DeleteDC")

func DeleteDC(hdc HDC) error {
	return sysutil.MustTrue(lzDeleteDC.Call(uintptr(hdc)))
}

type DPI_AWARENESS_CONTEXT HANDLE

const (
	DPI_AWARENESS_CONTEXT_UNAWARE              DPI_AWARENESS_CONTEXT = ^DPI_AWARENESS_CONTEXT(0)     // -1
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         DPI_AWARENESS_CONTEXT = ^DPI_AWARENESS_CONTEXT(0) - 1 // -2
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    DPI_AWARENESS_CONTEXT = ^DPI_AWARENESS_CONTEXT(0) - 2 // -3
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 DPI_AWARENESS_CONTEXT = ^DPI_AWARENESS_CONTEXT(0) - 3 // -4
	DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED    DPI_AWARENESS_CONTEXT = ^DPI_AWARENESS_CONTEXT(0) - 4 //-5
)

var lzSetThreadDpiAwarenessContext = lzUser32.NewProc("SetThreadDpiAwarenessContext")

func SetThreadDpiAwarenessContext(ctx DPI_AWARENESS_CONTEXT) (DPI_AWARENESS_CONTEXT, error) {
	return sysutil.MustNotZero[DPI_AWARENESS_CONTEXT](lzSetThreadDpiAwarenessContext.Call(uintptr(ctx)))
}

var lzAreDpiAwarenessContextsEqual = lzUser32.NewProc("AreDpiAwarenessContextsEqual")

func AreDpiAwarenessContextsEqual(ctx1, ctx2 DPI_AWARENESS_CONTEXT) bool {
	r, _, _ := lzAreDpiAwarenessContextsEqual.Call(uintptr(ctx1), uintptr(ctx2))
	return r != 0
}

var lzSetFocus = lzUser32.NewProc("SetFocus")

func SetFocus(hwnd HWND) HWND {
	return sysutil.As[HWND](lzSetFocus.Call(uintptr(hwnd)))
}

var lzGetModuleFileNameW = lzKernel32.NewProc("GetModuleFileNameW")

func GetModuleFileNameW(h HMODULE, buf []WCHAR) DWORD {
	var p *WCHAR
	if len(buf) > 0 {
		p = &buf[0]
	}
	return sysutil.As[DWORD](lzGetModuleFileNameW.Call(uintptr(h), uintptr(unsafe.Pointer(p)), uintptr(len(buf))))
}

type UUID struct {
	unused1 ULONG
	unused2 USHORT
	unused3 USHORT
	unused4 [8]UCHAR
	// Don't make these fields blanks(_).
	// Blank fields are not considered when comparing equality.
}

type GUID = UUID

var lzLoadIconW = lzUser32.NewProc("LoadIconW")

const (
	IDI_APPLICATION uintptr = 32512
	IDI_HAND        uintptr = 32513
	IDI_QUESTION    uintptr = 32514
	IDI_EXCLAMATION uintptr = 32515
	IDI_ASTERISK    uintptr = 32516
	IDI_WINLOGO     uintptr = 32517
	IDI_SHIELD      uintptr = 32518
	IDI_WARNING     uintptr = IDI_EXCLAMATION
	IDI_ERROR       uintptr = IDI_HAND
	IDI_INFORMATION uintptr = IDI_ASTERISK
)

func LoadIconW(instance HINSTANCE, name *WCHAR) (HICON, error) {
	return sysutil.MustNotZero[HICON](lzLoadIconW.Call(uintptr(instance), uintptr(unsafe.Pointer(name))))
}

var lzSetForegroundWindow = lzUser32.NewProc("SetForegroundWindow")

func SetForegroundWindow(hwnd HWND) bool {
	return sysutil.AsBool(lzSetForegroundWindow.Call(uintptr(hwnd)))
}

var lzGetSystemMetrics = lzUser32.NewProc("GetSystemMetrics")

func GetSystemMetrics(index SystemMetricsIndex) INT {
	return sysutil.As[INT](lzGetSystemMetrics.Call(uintptr(index)))
}

type SystemMetricsIndex INT

const (
	SM_CXSCREEN                    SystemMetricsIndex = 0
	SM_CYSCREEN                    SystemMetricsIndex = 1
	SM_CXVSCROLL                   SystemMetricsIndex = 2
	SM_CYHSCROLL                   SystemMetricsIndex = 3
	SM_CYCAPTION                   SystemMetricsIndex = 4
	SM_CXBORDER                    SystemMetricsIndex = 5
	SM_CYBORDER                    SystemMetricsIndex = 6
	SM_CXDLGFRAME                  SystemMetricsIndex = 7
	SM_CYDLGFRAME                  SystemMetricsIndex = 8
	SM_CYVTHUMB                    SystemMetricsIndex = 9
	SM_CXHTHUMB                    SystemMetricsIndex = 10
	SM_CXICON                      SystemMetricsIndex = 11
	SM_CYICON                      SystemMetricsIndex = 12
	SM_CXCURSOR                    SystemMetricsIndex = 13
	SM_CYCURSOR                    SystemMetricsIndex = 14
	SM_CYMENU                      SystemMetricsIndex = 15
	SM_CXFULLSCREEN                SystemMetricsIndex = 16
	SM_CYFULLSCREEN                SystemMetricsIndex = 17
	SM_CYKANJIWINDOW               SystemMetricsIndex = 18
	SM_MOUSEPRESENT                SystemMetricsIndex = 19
	SM_CYVSCROLL                   SystemMetricsIndex = 20
	SM_CXHSCROLL                   SystemMetricsIndex = 21
	SM_DEBUG                       SystemMetricsIndex = 22
	SM_SWAPBUTTON                  SystemMetricsIndex = 23
	SM_RESERVED1                   SystemMetricsIndex = 24
	SM_RESERVED2                   SystemMetricsIndex = 25
	SM_RESERVED3                   SystemMetricsIndex = 26
	SM_RESERVED4                   SystemMetricsIndex = 27
	SM_CXMIN                       SystemMetricsIndex = 28
	SM_CYMIN                       SystemMetricsIndex = 29
	SM_CXSIZE                      SystemMetricsIndex = 30
	SM_CYSIZE                      SystemMetricsIndex = 31
	SM_CXFRAME                     SystemMetricsIndex = 32
	SM_CYFRAME                     SystemMetricsIndex = 33
	SM_CXMINTRACK                  SystemMetricsIndex = 34
	SM_CYMINTRACK                  SystemMetricsIndex = 35
	SM_CXDOUBLECLK                 SystemMetricsIndex = 36
	SM_CYDOUBLECLK                 SystemMetricsIndex = 37
	SM_CXICONSPACING               SystemMetricsIndex = 38
	SM_CYICONSPACING               SystemMetricsIndex = 39
	SM_MENUDROPALIGNMENT           SystemMetricsIndex = 40
	SM_PENWINDOWS                  SystemMetricsIndex = 41
	SM_DBCSENABLED                 SystemMetricsIndex = 42
	SM_CMOUSEBUTTONS               SystemMetricsIndex = 43
	SM_CXFIXEDFRAME                SystemMetricsIndex = SM_CXDLGFRAME
	SM_CYFIXEDFRAME                SystemMetricsIndex = SM_CYDLGFRAME
	SM_CXSIZEFRAME                 SystemMetricsIndex = SM_CXFRAME
	SM_CYSIZEFRAME                 SystemMetricsIndex = SM_CYFRAME
	SM_SECURE                      SystemMetricsIndex = 44
	SM_CXEDGE                      SystemMetricsIndex = 45
	SM_CYEDGE                      SystemMetricsIndex = 46
	SM_CXMINSPACING                SystemMetricsIndex = 47
	SM_CYMINSPACING                SystemMetricsIndex = 48
	SM_CXSMICON                    SystemMetricsIndex = 49
	SM_CYSMICON                    SystemMetricsIndex = 50
	SM_CYSMCAPTION                 SystemMetricsIndex = 51
	SM_CXSMSIZE                    SystemMetricsIndex = 52
	SM_CYSMSIZE                    SystemMetricsIndex = 53
	SM_CXMENUSIZE                  SystemMetricsIndex = 54
	SM_CYMENUSIZE                  SystemMetricsIndex = 55
	SM_ARRANGE                     SystemMetricsIndex = 56
	SM_CXMINIMIZED                 SystemMetricsIndex = 57
	SM_CYMINIMIZED                 SystemMetricsIndex = 58
	SM_CXMAXTRACK                  SystemMetricsIndex = 59
	SM_CYMAXTRACK                  SystemMetricsIndex = 60
	SM_CXMAXIMIZED                 SystemMetricsIndex = 61
	SM_CYMAXIMIZED                 SystemMetricsIndex = 62
	SM_NETWORK                     SystemMetricsIndex = 63
	SM_CLEANBOOT                   SystemMetricsIndex = 67
	SM_CXDRAG                      SystemMetricsIndex = 68
	SM_CYDRAG                      SystemMetricsIndex = 69
	SM_SHOWSOUNDS                  SystemMetricsIndex = 70
	SM_CXMENUCHECK                 SystemMetricsIndex = 71
	SM_CYMENUCHECK                 SystemMetricsIndex = 72
	SM_SLOWMACHINE                 SystemMetricsIndex = 73
	SM_MIDEASTENABLED              SystemMetricsIndex = 74
	SM_MOUSEWHEELPRESENT           SystemMetricsIndex = 75
	SM_XVIRTUALSCREEN              SystemMetricsIndex = 76
	SM_YVIRTUALSCREEN              SystemMetricsIndex = 77
	SM_CXVIRTUALSCREEN             SystemMetricsIndex = 78
	SM_CYVIRTUALSCREEN             SystemMetricsIndex = 79
	SM_CMONITORS                   SystemMetricsIndex = 80
	SM_SAMEDISPLAYFORMAT           SystemMetricsIndex = 81
	SM_IMMENABLED                  SystemMetricsIndex = 82
	SM_CXFOCUSBORDER               SystemMetricsIndex = 83
	SM_CYFOCUSBORDER               SystemMetricsIndex = 84
	SM_TABLETPC                    SystemMetricsIndex = 86
	SM_MEDIACENTER                 SystemMetricsIndex = 87
	SM_STARTER                     SystemMetricsIndex = 88
	SM_SERVERR2                    SystemMetricsIndex = 89
	SM_MOUSEHORIZONTALWHEELPRESENT SystemMetricsIndex = 91
	SM_CXPADDEDBORDER              SystemMetricsIndex = 92
	SM_DIGITIZER                   SystemMetricsIndex = 94
	SM_MAXIMUMTOUCHES              SystemMetricsIndex = 95
	SM_REMOTESESSION               SystemMetricsIndex = 0x1000
	SM_SHUTTINGDOWN                SystemMetricsIndex = 0x2000
	SM_REMOTECONTROL               SystemMetricsIndex = 0x2001
	SM_CARETBLINKINGENABLED        SystemMetricsIndex = 0x2002
	SM_CONVERTIBLESLATEMODE        SystemMetricsIndex = 0x2003
	SM_SYSTEMDOCKED                SystemMetricsIndex = 0x2004
)

var lzGetStockObject = lzGdi32.NewProc("GetStockObject")

// GetStockObject retrieves a handle to one of the stock pens, brushes, fonts, or
// palettes.
//
// Returns 0 if it fails, no additional error information is available.
func GetStockObject[H HGDIOBJ](object StockObjectType) H {
	h, _, _ := lzGetStockObject.Call(uintptr(object))
	return H(h)
}

type StockObjectType int

const (
	WHITE_BRUSH         = StockObjectType(0)
	LTGRAY_BRUSH        = StockObjectType(1)
	GRAY_BRUSH          = StockObjectType(2)
	DKGRAY_BRUSH        = StockObjectType(3)
	BLACK_BRUSH         = StockObjectType(4)
	NULL_BRUSH          = StockObjectType(5)
	HOLLOW_BRUSH        = NULL_BRUSH
	WHITE_PEN           = StockObjectType(6)
	BLACK_PEN           = StockObjectType(7)
	NULL_PEN            = StockObjectType(8)
	OEM_FIXED_FONT      = StockObjectType(10)
	ANSI_FIXED_FONT     = StockObjectType(11)
	ANSI_VAR_FONT       = StockObjectType(12)
	SYSTEM_FONT         = StockObjectType(13)
	DEVICE_DEFAULT_FONT = StockObjectType(14)
	DEFAULT_PALETTE     = StockObjectType(15)
	SYSTEM_FIXED_FONT   = StockObjectType(16)
	DEFAULT_GUI_FONT    = StockObjectType(17)
	DC_BRUSH            = StockObjectType(18)
	DC_PEN              = StockObjectType(19)
	STOCK_LAST          = StockObjectType(19)
)

var lzSetWindowsHookExW = lzUser32.NewProc("SetWindowsHookExW")

func SetWindowsHookExW(idHook HookID, lpfn uintptr, hMod HINSTANCE, dwThreadId DWORD) (HHOOK, error) {
	return sysutil.MustNotZero[HHOOK](lzSetWindowsHookExW.Call(uintptr(idHook), lpfn, uintptr(hMod), uintptr(dwThreadId)))
}

type HookID INT

const (
	WH_MSGFILTER       HookID = -1
	WH_JOURNALRECORD   HookID = 0
	WH_JOURNALPLAYBACK HookID = 1
	WH_KEYBOARD        HookID = 2
	WH_GETMESSAGE      HookID = 3
	WH_CALLWNDPROC     HookID = 4
	WH_CBT             HookID = 5
	WH_SYSMSGFILTER    HookID = 6
	WH_MOUSE           HookID = 7
	WH_HARDWARE        HookID = 8
	WH_DEBUG           HookID = 9
	WH_SHELL           HookID = 10
	WH_FOREGROUNDIDLE  HookID = 11
	WH_CALLWNDPROCRET  HookID = 12
	WH_KEYBOARD_LL     HookID = 13
	WH_MOUSE_LL        HookID = 14
)

var lzUnhookWindowsHookEx = lzUser32.NewProc("UnhookWindowsHookEx")

func UnhookWindowsHookEx(hhk HHOOK) error {
	return sysutil.MustTrue(lzUnhookWindowsHookEx.Call(uintptr(hhk)))
}

var lzCallNextHookEx = lzUser32.NewProc("CallNextHookEx")

func CallNextHookEx(hhk HHOOK, nCode HookCode, wParam WPARAM, lParam LPARAM) LRESULT {
	return sysutil.As[LRESULT](lzCallNextHookEx.Call(uintptr(hhk), uintptr(nCode), uintptr(wParam), uintptr(lParam)))
}

type HookCode INT

const (
	HC_ACTION      HookCode = 0
	HC_GETNEXT     HookCode = 1
	HC_SKIP        HookCode = 2
	HC_NOREMOVE    HookCode = 3
	HC_NOREM       HookCode = HC_NOREMOVE
	HC_SYSMODALON  HookCode = 4
	HC_SYSMODALOFF HookCode = 5
)

var lzSetViewportOrgEx = lzGdi32.NewProc("SetViewportOrgEx")

func SetViewportOrgEx(hdc HDC, x INT, y INT, prev *POINT) error {
	return sysutil.MustTrue(lzSetViewportOrgEx.Call(uintptr(hdc), uintptr(x), uintptr(y), uintptr(unsafe.Pointer(prev))))
}

var lzSetLayeredWindowAttributes = lzUser32.NewProc("SetLayeredWindowAttributes")

type LayeredWindowFlag DWORD

const (
	LWA_ALPHA    LayeredWindowFlag = 0x00000002
	LWA_COLORKEY LayeredWindowFlag = 0x00000001
)

func SetLayeredWindowAttributes(hwnd HWND, crKey COLORREF, bAlpha BYTE, dwFlags LayeredWindowFlag) error {
	return sysutil.MustTrue(lzSetLayeredWindowAttributes.Call(uintptr(hwnd), uintptr(crKey), uintptr(bAlpha), uintptr(dwFlags)))
}
