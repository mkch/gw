package dialog

import (
	"unsafe"

	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/paint/font"
	"github.com/mkch/gw/win32"
	"golang.org/x/sys/windows"
)

type Limit struct {
	Min, Max win32.INT
}

type ChooseFontSpec struct {
	Owner   win32.HWND
	LogFont *font.LogFont // Initial value or nil.
	/* The following flags are automatically removed(added):
	CF_INITTOLOGFONTSTRUCT
	CF_APPLY
	CF_EFFECTS
	CF_ENABLEHOOK
	CF_ENABLETEMPLATE
	CF_ENABLETEMPLATEHANDLE
	CF_LIMITSIZE
	CF_USESTYLE
	*/
	Flags win32.CHOOSE_FONT_FLAG
	// If not nil, font effects(color, strikeout, underline) are enabled.
	// Strikeout and underline are specified in LogFont.
	Color          *win32.COLORREF
	PointSizeLimit *Limit                    // In point. Nil for none.
	OnApply        func(curFont *FontChosen) // If not nil, an Apply button is displayed, and OnApply is called if it is pressed.
}

type FontChosen struct {
	Font      *font.LogFont
	Type      win32.CHOOSE_FONT_TYPE
	PointSize win32.INT
	Color     win32.COLORREF
}

type chooseFontCustomData struct {
	dpi     win32.UINT
	onApply func(*FontChosen)
}

var chooseFontCustomDataMap = objectmap.New[*chooseFontCustomData](0, 0xFF)

type chooseFontData struct {
	*chooseFontCustomData
	*win32.CHOOSEFONTW
}

var chooseFontHwndMap = make(map[win32.HWND]*chooseFontData)

const WM_CHOOSEFONT_GETLOGFONT = (win32.WM_USER + 1)

var hookProc = windows.NewCallback(
	func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.UINT_PTR {
		switch message {
		case win32.WM_INITDIALOG:
			cf := (*win32.CHOOSEFONTW)(unsafe.Pointer(uintptr(unsafe.Pointer(nil)) + uintptr(lParam)))
			if customData, ok := chooseFontCustomDataMap.Value(objectmap.Handle(cf.CustomData)); ok {
				chooseFontHwndMap[hwnd] = &chooseFontData{chooseFontCustomData: customData, CHOOSEFONTW: cf}
			}
		case win32.WM_NCDESTROY:
			delete(chooseFontHwndMap, hwnd)
		case win32.WM_COMMAND:
			id := win32.LOWORD(wParam)
			if id == 1026 { // What is the const name for 1026??
				data := chooseFontHwndMap[hwnd]
				cf := *data.CHOOSEFONTW
				cf.LogFont = &win32.LOGFONTW{}
				win32.SendMessageW(hwnd, WM_CHOOSEFONT_GETLOGFONT, 0, win32.LPARAM(uintptr(unsafe.Pointer(cf.LogFont))))
				data.onApply(newFontChosen(&cf, data.dpi))
			}
		}
		return 0
	})

// ChooseFont displays a Font dialog.
// If the user cancels or closes the Font dialog box, it returns nil, nil.
// Nil spec means default setting.
func ChooseFont(spec *ChooseFontSpec) (*FontChosen, error) {
	// ChooseFont does not work well under PER_MONITOR_AWARE or PER_MONITOR_AWARE_V2.
	if oldDpiCtx, err := win32.SetThreadDpiAwarenessContext(win32.DPI_AWARENESS_CONTEXT_SYSTEM_AWARE); err != nil {
		return nil, err
	} else {
		defer win32.SetThreadDpiAwarenessContext(oldDpiCtx)
	}
	dpi := win32.GetDpiForSystem()
	if spec == nil {
		spec = &ChooseFontSpec{}
	}
	var cf = win32.CHOOSEFONTW{
		StructSize: win32.DWORD(unsafe.Sizeof(win32.CHOOSEFONTW{})),
		Owner:      spec.Owner,
		Flags:      spec.Flags,
	}
	cf.Flags &= ^(win32.CF_INITTOLOGFONTSTRUCT | win32.CF_APPLY | win32.CF_EFFECTS |
		win32.CF_ENABLEHOOK | win32.CF_ENABLETEMPLATE | win32.CF_ENABLETEMPLATEHANDLE |
		win32.CF_LIMITSIZE | win32.CF_USESTYLE)
	if spec.LogFont != nil {
		cf.LogFont = spec.LogFont.ForDPI(dpi)
		cf.Flags |= win32.CF_INITTOLOGFONTSTRUCT
	}
	if spec.Color != nil {
		cf.Color = *spec.Color
		cf.Flags |= win32.CF_EFFECTS
	}
	if spec.PointSizeLimit != nil {
		cf.SizeMin = spec.PointSizeLimit.Min
		cf.SizeMax = spec.PointSizeLimit.Max
		cf.Flags |= win32.CF_LIMITSIZE
	}
	if spec.OnApply != nil {
		cf.Flags |= (win32.CF_APPLY | win32.CF_ENABLEHOOK)
		cf.Hook = hookProc
		h := chooseFontCustomDataMap.Add(&chooseFontCustomData{dpi: dpi, onApply: spec.OnApply})
		defer chooseFontCustomDataMap.Remove(h)
		cf.CustomData = win32.LPARAM(h)
	}

	if ok, err := win32.ChooseFontW(&cf); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	} else {
		return newFontChosen(&cf, dpi), nil
	}
}

func newFontChosen(cf *win32.CHOOSEFONTW, dpi win32.UINT) *FontChosen {
	return &FontChosen{
		Font:      font.NewLogFont(cf.LogFont, dpi),
		Type:      cf.FontType,
		PointSize: cf.PointSize,
		Color:     cf.Color,
	}
}
