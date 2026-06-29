package win32util

import (
	"errors"
	"unicode/utf16"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/win32"
	"golang.org/x/sys/windows"
)

func CString(str string, p *[]win32.WCHAR) {
	*p = (*p)[:0]
	s := utf16.Encode([]rune(str))
	if len(s) != 0 {
		a := unsafe.Slice((*win32.WCHAR)(unsafe.Pointer(&s[0])), len(s))
		*p = append(*p, a...)
	}
	*p = append(*p, 0) // 0 terminated.
}

// CStrLen returns the length of a null terminated C string, not including the terminating 0.
// Parameter str is a pointer to the C string, and bufSize is the maximum buffer size(in unit of win32.WCHAR),
// including the terminating 0.
// If str is nil and bufSize is zero, CStrLen returns 0.
// If there's no 0 terminator within the first bufSize elements, a run-time panic occurs.
// if bufSize is negative, or if str is nil and bufSize is not zero, a run-time panic occurs
func CStrLen(str *win32.WCHAR, bufSize int) int {
	if str == nil && bufSize == 0 {
		return 0
	}
	if str == nil && bufSize != 0 {
		panic("str is nil but bufSize is not zero")
	}

	for i := range bufSize {
		c := *(*win32.WCHAR)(unsafe.Add(unsafe.Pointer(str), uintptr(i)*unsafe.Sizeof(win32.WCHAR(0))))
		if c == 0 {
			return i
		}
	}
	panic("no 0 terminator in c string")
}

// GoString converts a null terminated C string to go string.
// size is the buffer length of C string, including the terminating null.
//
// If size is -1, the buffer size will be calculated using CStrLen(p, 0xFFF)+1.
func GoString(p *win32.WCHAR, size int) string {
	if size == -1 {
		size = CStrLen(p, 0xFFF) + 1
	}
	s := unsafe.Slice((*uint16)(unsafe.Pointer(p)), size-1) // size-1 to exclude the terminating null.
	return string(utf16.Decode(s))
}

// CopyCString copies null terminated C string src to dest.
// Returns the count of win32.WCHAR copied, including the terminating null.
func CopyCString(dest, src []win32.WCHAR) (charCopied int) {
	charCopied = copy(dest[:len(dest)-1], src)
	if charCopied < len(src) {
		dest[charCopied] = 0
		charCopied++
	}
	return
}

type WndClass struct {
	ClassName  string
	WndProc    win32.WndProc
	Instance   win32.HINSTANCE // 0 for this module.
	Style      win32.CLASS_STYLE
	ClsExtra   win32.INT
	WndExtra   win32.INT
	Icon       win32.HICON
	Cursor     win32.HCURSOR
	Background win32.HBRUSH
	MenuName   string
	IconSm     win32.HICON
}

func RegisterClass(cls *WndClass) (win32.ATOM, error) {
	var classNameBuf []win32.WCHAR
	CString(cls.ClassName, &classNameBuf)
	var menuNamePtr *win32.WCHAR
	if len(cls.MenuName) > 0 {
		var buf []win32.WCHAR
		CString(cls.MenuName, &buf)
		menuNamePtr = &buf[0]
	}

	var wndClass = win32.WNDCLASSEXW{
		Size:       win32.UINT(unsafe.Sizeof(win32.WNDCLASSEXW{})),
		ClassName:  &classNameBuf[0],
		WndProc:    windows.NewCallback(cls.WndProc),
		Style:      cls.Style,
		ClsExtra:   cls.ClsExtra,
		WndExtra:   cls.WndExtra,
		Instance:   cls.Instance,
		Icon:       cls.Icon,
		Cursor:     cls.Cursor,
		Background: cls.Background,
		MenuName:   menuNamePtr,
		IconSm:     cls.IconSm,
	}
	if wndClass.Instance == 0 {
		instance, _ := win32.GetModuleHandleW[win32.HINSTANCE](nil)
		wndClass.Instance = instance
	}

	return win32.RegisterClassExW(&wndClass)
}

func UnregisterClass(className string, instance win32.HINSTANCE) error {
	var classNameBuf []win32.WCHAR
	CString(className, &classNameBuf)
	return win32.UnregisterClassW(&classNameBuf[0], instance)
}

type Wnd struct {
	ClassName  string
	WindowName string
	Style      win32.WindowStyle
	ExStyle    win32.WindowExStyle
	X          win32.INT
	Y          win32.INT
	Width      win32.INT
	Height     win32.INT
	WndParent  win32.HWND
	Menu       win32.HMENU
	Instance   win32.HINSTANCE // 0 for this module.
	Param      win32.UINT_PTR
}

func CreateWindow(spec *Wnd) (win32.HWND, error) {
	instance := spec.Instance
	if instance == 0 {
		instance, _ = win32.GetModuleHandleW[win32.HINSTANCE](nil)
	}
	var classNameBuf []win32.WCHAR
	CString(spec.ClassName, &classNameBuf)
	var windowNamePtr *win32.WCHAR
	if len(spec.WindowName) > 0 {
		var windowNameBuf []win32.WCHAR
		CString(spec.WindowName, &windowNameBuf)
		windowNamePtr = &windowNameBuf[0]
	}

	return win32.CreateWindowExW(spec.ExStyle, &classNameBuf[0], windowNamePtr, spec.Style,
		spec.X, spec.Y, spec.Width, spec.Height,
		spec.WndParent, spec.Menu, instance, spec.Param)
}

func GetWindowText(hwnd win32.HWND) (string, error) {
	l, err := win32.GetWindowTextLengthW(hwnd)
	if err != nil {
		return "", err
	}
	buf := make([]win32.WCHAR, l+1)
	n, err := win32.GetWindowTextW(hwnd, &buf[0], len(buf))
	if n == 0 && err != nil {
		return "", err
	}
	return GoString(&buf[0], n+1), nil
}

func SetWindowText(hwnd win32.HWND, str string) error {
	var buf []win32.WCHAR
	CString(str, &buf)
	return win32.SetWindowTextW(hwnd, &buf[0])
}

// EmptyDialogTemplate allocates an empty dialog template.
// x, y, cx, cy are in pixel format in screen coordinates.
func EmptyDialogTemplate(style win32.DWORD, exStyle win32.DWORD, x win32.SHORT, y win32.SHORT, cx win32.SHORT, cy win32.SHORT) *win32.DLGTEMPLATE {
	base := win32.GetDialogBaseUnits()
	xBase, yBase := win32.LOWORD(uintptr(base)), win32.HIWORD(uintptr(base))
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getdialogbaseunits
	x = win32.SHORT(win32.MulDiv(win32.INT(x), 4, win32.INT(xBase)))
	y = win32.SHORT(win32.MulDiv(win32.INT(y), 4, win32.INT(yBase)))
	cx = win32.SHORT(win32.MulDiv(win32.INT(cx), 4, win32.INT(xBase)))
	cy = win32.SHORT(win32.MulDiv(win32.INT(cy), 4, win32.INT(yBase)))

	type template struct {
		win32.DLGTEMPLATE
		_, _, _ win32.WORD
	}

	return (*win32.DLGTEMPLATE)(unsafe.Pointer(&template{
		DLGTEMPLATE: win32.DLGTEMPLATE{
			Style:   style,
			ExStyle: exStyle,
			X:       x,
			Y:       y,
			CX:      cx,
			CY:      cy,
		},
	}))
}

func ClientToScreen(hwnd win32.HWND, rect *win32.RECT) error {
	if err := win32.ClientToScreen(hwnd, rect.TopLeft()); err != nil {
		return err
	}
	if err := win32.ClientToScreen(hwnd, rect.BottomRight()); err != nil {
		return err
	}
	return nil
}

func ScreenToClient(hwnd win32.HWND, rect *win32.RECT) error {
	if err := win32.ScreenToClient(hwnd, rect.TopLeft()); err != nil {
		return err
	}
	if err := win32.ScreenToClient(hwnd, rect.BottomRight()); err != nil {
		return err
	}
	return nil
}

func MessageBoxEx(owner win32.HWND, text string, caption string, typ win32.MESSAGE_BOX_TYPE, langID win32.WORD) (int, error) {
	var textBuf []win32.WCHAR
	CString(text, &textBuf)
	var captionBuf []win32.WCHAR
	CString(caption, &captionBuf)
	r, err := win32.MessageBoxExW(owner, &textBuf[0], &captionBuf[0], typ, langID)
	return int(r), err
}

func MessageBox(owner win32.HWND, text string, caption string, typ win32.MESSAGE_BOX_TYPE) (int, error) {
	return MessageBoxEx(owner, text, caption, typ, 0)
}

// CreatePen simulates CreatePen API using ExtCreatePen.
// Style can be one of PS_SOLID, PS_DASH, PS_DOT, PS_DASHDOT, PS_DASHDOTDOT, PS_NULL, PS_INSIDEFRAME.
func CreatePen(style win32.PEN_STYLE, width win32.DWORD, color win32.COLORREF) (win32.HPEN, error) {
	brush := win32.LOGBRUSH{
		Style: win32.BS_SOLID,
		Color: color,
	}
	return win32.ExtCreatePen(style|win32.PS_GEOMETRIC, width, &brush, nil)
}

type ModifyStyleSpec struct {
	Add    win32.WindowStyle
	Remove win32.WindowStyle
}

// ModifyWindowStyle modifies the window style of the specified window by removing and adding styles.
func ModifyWindowStyle(hwnd win32.HWND, spec ModifyStyleSpec) error {
	style, err := win32.GetWindowLongPtrW(hwnd, win32.GWL_STYLE)
	if err != nil {
		return err
	}
	style &^= win32.LONG_PTR(spec.Remove)
	style |= win32.LONG_PTR(spec.Add)
	_, err = win32.SetWindowLongPtrW(hwnd, win32.GWL_STYLE, style)
	return err
}

type ModifyExStyleSpec struct {
	Add    win32.WindowExStyle
	Remove win32.WindowExStyle
}

// ModifyWindowExStyle modifies the extended window style of the specified window by removing and adding styles.
func ModifyWindowExStyle(hwnd win32.HWND, spec ModifyExStyleSpec) error {
	exStyle, err := win32.GetWindowLongPtrW(hwnd, win32.GWL_EXSTYLE)
	if err != nil {
		return err
	}
	exStyle &^= win32.LONG_PTR(spec.Remove)
	exStyle |= win32.LONG_PTR(spec.Add)
	_, err = win32.SetWindowLongPtrW(hwnd, win32.GWL_EXSTYLE, exStyle)
	return err
}

// CopyWindowClass registers a new window class by copying an existing class specified by src, and changing the class name to newClassName.
func CopyWindowClass(src win32.ATOM, newClassName string) (win32.ATOM, error) {
	// Retrieve the default class info
	var cls = win32.WNDCLASSEXW{Size: win32.UINT(unsafe.Sizeof(win32.WNDCLASSEXW{}))}
	if err := win32.GetClassInfoExWAtom(gg.Must(win32.GetModuleHandleW[win32.HINSTANCE](nil)), src, &cls); err != nil {
		return 0, err
	}
	// Change the class name
	var classNameBuf []win32.WCHAR
	CString(newClassName, &classNameBuf)
	cls.ClassName = &classNameBuf[0]
	// Register the new class
	return win32.RegisterClassExW(&cls)
}

// KeyName returns the localized string representation of vk.
// If doNotCareLeftRight is true, left and right modifier keys will not be distinguished.
// If keyboardLayout is not 0, the name will be obtained according to the specified keyboard layout;
// otherwise, the layout of current thread will be used.
func KeyName(vk win32.VKCode, doNotCareLeftRight bool, keyboardLayout win32.HKL) (string, error) {
	if keyboardLayout == 0 {
		keyboardLayout = win32.GetKeyboardLayout(0)
	}
	// Get complete scan code (with E0/E1 extended key flag)

	// Here are the special cases for keys with non-standard scan code mappings.
	// For most keys, MapVirtualKeyW with MAPVK_VK_TO_VSC_EX can be used to get the complete scan code.
	var fullScanCode win32.UINT
	switch vk {
	case win32.VK_PAUSE:
		// VK_PAUSE has a compound scan code (E1-1D-45); MapVirtualKeyW returns
		// 0xE11D, which causes GetKeyNameTextW to return "Right Ctrl" when the
		// extended bit is set. The correct lParam to obtain "Pause" is scan
		// code 0x45 without the extended bit.
		fullScanCode = 0x45
	case win32.VK_NUMLOCK:
		// MapVirtualKeyW omits the E0 prefix for VK_NUMLOCK; supply it manually.
		fullScanCode = 0xE045
	case win32.VK_SNAPSHOT:
		// MapVirtualKeyW maps VK_SNAPSHOT to scan code 0x54 (Sys Req), causing
		// GetKeyNameTextW to return "Sys Req". Print Screen is scan code E0-37.
		fullScanCode = 0xE037
	default:
		fullScanCode = win32.MapVirtualKeyExW(win32.UINT(vk), win32.MAPVK_VK_TO_VSC_EX, keyboardLayout)
		if fullScanCode == 0 {
			return "", errors.New("invalid virtual key code")
		}
		switch vk {
		case win32.VK_PRIOR, win32.VK_NEXT, win32.VK_END, win32.VK_HOME, win32.VK_LEFT, win32.VK_UP, win32.VK_RIGHT, win32.VK_DOWN, win32.VK_INSERT, win32.VK_DELETE:
			fullScanCode |= 0xE000 // These keys have the E0 extended key flag.

		}
	}
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeynametextw
	// 16-23	The scan code. The value depends on the OEM.
	lParam := (fullScanCode & 0xFF) << 16
	highByte := (fullScanCode >> 8) & 0xFF
	// MapVirtualKeyW(MAPVK_VK_TO_VSC_EX):
	// If the scan code is an extended scan code, the high byte of the returned value will contain either 0xe0 or 0xe1
	if highByte == 0xE0 || highByte == 0xE1 {
		// GetKeyNameTextW:
		// 24	Indicates whether the key is an extended key
		lParam |= 1 << 24
	}
	if doNotCareLeftRight {
		// GetKeyNameTextW:
		// 25	"Do not care" bit. ... should not distinguish between left and right CTRL and SHIFT keys, for example.
		lParam |= 1 << 25
	}

	var nameBuf [32]win32.WCHAR
	n, err := win32.GetKeyNameTextW(win32.LONG(lParam), &nameBuf[0], len(nameBuf))
	if err != nil {
		return "", err
	} else if n <= 0 {
		return "", errors.New("win32.GetKeyNameTextW failed")
	} else if n == win32.INT(len(nameBuf)) {
		panic("buffer too small")
	}
	return GoString(&nameBuf[0], int(n+1)), nil
}

// SetClientSize modifies the size of the window to ensure a desired client size.
func SetClientSize(hwnd win32.HWND, width, height win32.INT) error {
	style, err := win32.GetWindowLongPtrW(hwnd, win32.GWL_STYLE)
	if err != nil {
		return err
	}
	exStyle, err := win32.GetWindowLongPtrW(hwnd, win32.GWL_EXSTYLE)
	if err != nil {
		return err
	}
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-adjustwindowrectex
	exStyle &^= win32.LONG_PTR(win32.WS_OVERLAPPED)

	var rect = win32.RECT{Right: win32.LONG(width), Bottom: win32.LONG(height)}
	// Calculate the estimated window size based on the desired client area size and the window styles.
	menu, err := win32.GetMenu(hwnd)
	if err != nil {
		return err
	}
	if err = win32.AdjustWindowRectEx(&rect, win32.WindowStyle(style), menu != 0, win32.WindowExStyle(exStyle)); err != nil {
		return err
	}
	// Set the window to the estimated size.
	if err = win32.SetWindowPos(hwnd, 0, 0, 0, win32.INT(rect.Right-rect.Left), win32.INT(rect.Bottom-rect.Top), win32.SWP_NOMOVE|win32.SWP_NOZORDER); err != nil {
		return err
	}

	// Get the actual client area size after resizing.
	var clientRect win32.RECT
	if err = win32.GetClientRect(hwnd, &clientRect); err != nil {
		return err
	}
	// Do any correction if the actual client area size is different from the desired size.
	// Usually caused by menu bar wrapping.
	var diffX = clientRect.Width() - win32.LONG(width)
	var diffY = clientRect.Height() - win32.LONG(height)
	if diffX != 0 || diffY != 0 {
		// Get the new window size.
		err = win32.GetWindowRect(hwnd, &rect)
		if err != nil {
			return err
		}
		// Calculate the fixed window size.
		fixedWidth := win32.INT(rect.Right - rect.Left - diffX)
		fixedHeight := win32.INT(rect.Bottom - rect.Top - diffY)
		// Set the window to the fixed size.
		if err = win32.SetWindowPos(hwnd, 0, 0, 0, fixedWidth, fixedHeight, win32.SWP_NOMOVE|win32.SWP_NOZORDER); err != nil {
			return err
		}
	}
	return nil
}
