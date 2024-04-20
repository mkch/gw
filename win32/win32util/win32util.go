package win32util

import (
	"unicode/utf16"
	"unsafe"

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
	ClassName   string
	WindowName  string
	Style       win32.WINDOW_STYLE
	ExStyle     win32.WINDOW_EX_STYLE
	X           win32.INT
	Y           win32.INT
	Width       win32.INT
	Height      win32.INT
	WndParent   win32.HWND
	InParentDPI bool // Whether X, Y, Width and Height are in WndParent's DPI. USER_DEFAULT_SCREEN_DPI is used if false.
	Menu        win32.HMENU
	Instance    win32.HINSTANCE // 0 for this module.
	Param       win32.UINT_PTR
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
	x := spec.X
	y := spec.Y
	cx := spec.Width
	cy := spec.Height
	if spec.InParentDPI {
		// spec.WndParent can't be 0 if InCurrentDPI is true.
		dpi, err := win32.GetDpiForWindow(spec.WndParent)
		if err != nil {
			return 0, err
		}
		x = FromDefaultDPI(spec.X, dpi)
		y = FromDefaultDPI(spec.Y, dpi)
		cx = FromDefaultDPI(spec.Width, dpi)
		cy = FromDefaultDPI(spec.Height, dpi)
	}
	return win32.CreateWindowExW(spec.ExStyle, &classNameBuf[0], windowNamePtr, spec.Style,
		x, y, cx, cy,
		spec.WndParent, spec.Menu, instance, spec.Param)
}

func GetWindowText(hwnd win32.HWND) (string, error) {
	l, err := win32.GetWindowTextLengthW(hwnd)
	if err != nil {
		return "", err
	}
	buf := make([]win32.WCHAR, l)
	n, err := win32.GetWindowTextW(hwnd, &buf[0], l)
	if n == 0 && err != nil {
		return "", err
	}
	return GoString(&buf[0], l), nil
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

// DPIConv converts a value from old DPI to new DPI.
func DPIConv[T ~int32 | ~uint32](oldValue T, oldDPI, newDPI win32.UINT) (newValue T) {
	return T(win32.MulDiv(win32.INT(oldValue), win32.INT(newDPI), win32.INT(oldDPI)))
}

// FromDefaultDPI convert value from USER_DEFAULT_SCREEN_DPI(96) to a new DPI.
func FromDefaultDPI[T ~int32 | ~uint32](value T, dpi win32.UINT) T {
	return DPIConv(value, win32.USER_DEFAULT_SCREEN_DPI, dpi)
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
