package win32util

import (
	"errors"
	"fmt"
	"slices"
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
	s := unsafe.Slice((*win32.WCHAR)(unsafe.Pointer(str)), bufSize)
	i := slices.Index(s, 0)
	if i == -1 {
		panic("no 0 terminator in c string")
	}
	return i
}

// GoString converts a null terminated C string to go string.
// size is the buffer length of C string, includes terminating null.
func GoString(p *win32.WCHAR, size int) string {
	s := unsafe.Slice((*uint16)(unsafe.Pointer(p)), size-1) // size-1 to exclude the terminating null.
	return string(utf16.Decode(s))
}

// CopyCString copies null terminated C string src to dest.
// Returns the count of win32.WCHAR copied, includes terminating null.
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

type Wnd struct {
	ClassName  string
	WindowName string
	Style      win32.WINDOW_STYLE
	ExStyle    win32.WINDOW_EX_STYLE
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
	Add    win32.WINDOW_STYLE
	Remove win32.WINDOW_STYLE
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
	Add    win32.WINDOW_EX_STYLE
	Remove win32.WINDOW_EX_STYLE
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

// KeyMessageLParam represents the lParam of keyboard messages (WM_KEYDOWN, WM_KEYUP, WM_SYSKEYDOWN, WM_SYSKEYUP).
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#keystroke-message-flags
type KeyMessageLParam win32.LPARAM

func (l KeyMessageLParam) String() string {
	return fmt.Sprintf("RepeatCount: %d, ScanCode: 0x%02X, Extended: %t, DialogMode: %t, MenuMode: %t, AltDown: %t, PreviousDown: %t, KeyUp: %t",
		l.RepeatCount(), l.ScanCode(), l.Extended(), l.DialogMode(), l.MenuMode(), l.AltDown(), l.PreviousDown(), l.KeyUp())
}

// RepeatCount returns the repeat count of the key message,
// which is the number of times the keystroke message should be sent while the user holds down the key.
func (l KeyMessageLParam) RepeatCount() uint16 {
	return uint16(l & 0xFFFF)
}

// SetRepeatCount sets the repeat count of the key message.
func (l *KeyMessageLParam) SetRepeatCount(count uint16) {
	*l = KeyMessageLParam((*l &^ 0xFFFF) | KeyMessageLParam(count))
}

// ScanCode returns the scan code of the key message.
// https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#scan-codes
func (l KeyMessageLParam) ScanCode() uint16 {
	return uint16((l >> 16) & 0xFF)
}

// SetScanCode sets the scan code of the key message.
func (l *KeyMessageLParam) SetScanCode(code uint16) {
	*l = KeyMessageLParam((*l &^ (0xFF << 16)) | (KeyMessageLParam(code) << 16))
}

// Extended returns the `Extended Key` flag in l.
// The `Extended Key` flag indicates whether the key is an extended key.
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#extended-key-flag
func (l KeyMessageLParam) Extended() bool {
	return l.bit(24)
}

// SetExtended sets the `Extended Key` flag in l.
func (l *KeyMessageLParam) SetExtended(extended bool) {
	l.setBit(24, extended)
}

// DialogMode returns whether a dialog box is active when the keystroke message was generated.
func (l KeyMessageLParam) DialogMode() bool {
	return l.bit(27)
}

// SetDialogMode sets whether a dialog box is active when the keystroke message was generated.
func (l *KeyMessageLParam) SetDialogMode(dialogMode bool) {
	l.setBit(27, dialogMode)
}

// MenuMode returns whether a menu is active when the keystroke message was generated.
func (l KeyMessageLParam) MenuMode() bool {
	return l.bit(28)
}

// SetMenuMode sets whether a menu is active when the keystroke message was generated.
func (l *KeyMessageLParam) SetMenuMode(menuMode bool) {
	l.setBit(28, menuMode)
}

// AltDown returns the `Context Code` flag in l.
// The `Context Code` flag indicates whether the Alt key was down when the keystroke message was generated.
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#context-code
func (l KeyMessageLParam) AltDown() bool {
	return l.bit(29)
}

// SetAltDown sets the `Context Code` flag in l.
func (l *KeyMessageLParam) SetAltDown(altDown bool) {
	l.setBit(29, altDown)
}

// PreviousDown returns the `Previous Key-State Flag` flag in l.
// The `Previous Key-State Flag` flag indicates whether the key was down before the message was sent.
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#previous-key-state-flag
func (l KeyMessageLParam) PreviousDown() bool {
	return l.bit(30)
}

// SetPreviousDown sets the `Previous Key-State Flag` flag in l.
func (l *KeyMessageLParam) SetPreviousDown(previousDown bool) {
	l.setBit(30, previousDown)
}

// KeyUp returns the `Transition-State Flag` flag in l.
// The `Transition-State Flag` flag indicates whether pressing a key generated the keystroke message.
// True for WM_KEYUP and WM_SYSKEYUP messages; false for WM_KEYDOWN and WM_SYSKEYDOWN messages;
// See https://learn.microsoft.com/en-us/windows/win32/inputdev/about-keyboard-input#transition-state-flag
func (l KeyMessageLParam) KeyUp() bool {
	return l.bit(31)
}

// SetKeyUp sets the 'KeyUp' flag in l.
func (l *KeyMessageLParam) SetKeyUp(keyUp bool) {
	l.setBit(31, keyUp)
}

// bit is a helper method to check if the i-th bit is set in l.
func (l KeyMessageLParam) bit(i int) bool {
	return (l & (1 << i)) != 0
}

// setBit is a helper method to set the i-th bit in l.
func (l *KeyMessageLParam) setBit(i int, value bool) {
	if value {
		*l |= KeyMessageLParam(1 << i)
	} else {
		*l &^= KeyMessageLParam(1 << i)
	}
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
