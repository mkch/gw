package fontdlg

import (
	"runtime"
	"unsafe"

	"github.com/mkch/gw/combobox"
	"github.com/mkch/gw/paint/font"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"golang.org/x/sys/windows"
)

type Limit struct {
	Min, Max win32.INT
}

type DefaultHookProc func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.UINT_PTR

type HookProc func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, cf *win32.CHOOSEFONTW, def DefaultHookProc) win32.UINT_PTR

type Spec struct {
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

	// If not nil, it is used as the hook procedure instead of the default one.
	HookProc HookProc
}

type FontChosen struct {
	Font      *font.LogFont
	Type      win32.CHOOSE_FONT_TYPE
	PointSize win32.INT
	Color     win32.COLORREF
}

type chooseFontCustomData struct {
	dpi      win32.UINT
	onApply  func(*FontChosen)
	hookProc HookProc
}

var chooseFontProp = win32util.NewWindowProp[win32.CHOOSEFONTW]("github.com/mkch/gw#ChooseFontProp")

const WM_CHOOSEFONT_GETLOGFONT = (win32.WM_USER + 1)

// defaultChooseFontHookProc is the default
func defaultChooseFontHookProc(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.UINT_PTR {
	switch message {
	case win32.WM_INITDIALOG:
		cf := (*win32.CHOOSEFONTW)(unsafe.Add(nil, lParam))
		chooseFontProp.Set(hwnd, cf)
	case win32.WM_NCDESTROY:
		chooseFontProp.Set(hwnd, nil)
	case win32.WM_COMMAND:
		id := win32.LOWORD(wParam)
		if id == 0x402 { // #include <dlgs.h> psh3
			cf := *chooseFontProp.Get(hwnd)
			cf.LogFont = &win32.LOGFONTW{}
			win32.SendMessageW(hwnd, WM_CHOOSEFONT_GETLOGFONT, 0, win32.LPARAM(uintptr(unsafe.Pointer(cf.LogFont))))

			if cf.Flags&win32.CF_EFFECTS != 0 {
				color, err := effectsColor(hwnd)
				if err == nil {
					cf.Color = color
				}
			}

			data := (*chooseFontCustomData)(unsafe.Pointer(uintptr(unsafe.Pointer(nil)) + uintptr(cf.CustomData)))
			data.onApply(newFontChosen(&cf, data.dpi))
		}
	}
	return 0
}

var hookProc = windows.NewCallback(
	func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.UINT_PTR {
		if message == win32.WM_INITDIALOG {
			// Must call the default hook proc at WM_INITDIALOG to set the CHOOSEFONTW pointer to the window property
			ret := defaultChooseFontHookProc(hwnd, message, wParam, lParam)
			cf := (*win32.CHOOSEFONTW)(unsafe.Add(nil, lParam))
			data := (*chooseFontCustomData)(unsafe.Pointer(uintptr(unsafe.Pointer(nil)) + uintptr(cf.CustomData)))
			if data.hookProc != nil {
				return data.hookProc(hwnd, message, wParam, lParam, cf, func(win32.HWND, win32.UINT, win32.WPARAM, win32.LPARAM) win32.UINT_PTR { return 0 })
			}
			return ret
		}
		cf := chooseFontProp.Get(hwnd)
		if cf == nil {
			return 0 // Cannot call the user hook proc without CHOOSEFONTW pointer.
		}
		data := (*chooseFontCustomData)(unsafe.Pointer(uintptr(unsafe.Pointer(nil)) + uintptr(cf.CustomData)))
		if data.hookProc != nil {
			return data.hookProc(hwnd, message, wParam, lParam, cf, defaultChooseFontHookProc)
		}
		return defaultChooseFontHookProc(hwnd, message, wParam, lParam)
	},
)

// effectsColor retrieves the color from the combo box in the Effects group.
func effectsColor(dlg win32.HWND) (win32.COLORREF, error) {
	combo, err := win32.GetDlgItem(dlg, 0x473) // #include <dlgs.h> cmb4
	if err != nil {
		return 0, err
	}
	i, err := win32.SendMessageW(combo, combobox.CB_GETCURSEL, 0, 0)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, nil
	}
	color, err := win32.SendMessageW(combo, combobox.CB_GETITEMDATA, win32.WPARAM(i), 0)
	if err != nil {
		return 0, err
	}
	return win32.COLORREF(color), nil
}

// ChooseFont displays a Font dialog.
// If the user cancels or closes the Font dialog box, it returns nil, nil.
// Nil spec means default setting.
func ChooseFont(spec *Spec) (*FontChosen, error) {
	// ChooseFont does not work well under PER_MONITOR_AWARE or PER_MONITOR_AWARE_V2.
	if oldDpiCtx, err := win32.SetThreadDpiAwarenessContext(win32.DPI_AWARENESS_CONTEXT_SYSTEM_AWARE); err != nil {
		return nil, err
	} else {
		defer win32.SetThreadDpiAwarenessContext(oldDpiCtx)
	}
	dpi := win32.GetDpiForSystem()
	if spec == nil {
		spec = &Spec{}
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
	} else {
		cf.LogFont = &win32.LOGFONTW{}
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
	}
	if spec.HookProc != nil {
		cf.Flags |= win32.CF_ENABLEHOOK
		cf.Hook = hookProc
	}

	p := &chooseFontCustomData{
		dpi:      dpi,
		onApply:  spec.OnApply,
		hookProc: spec.HookProc,
	}
	var pinner runtime.Pinner
	pinner.Pin(p)
	defer pinner.Unpin()
	cf.CustomData = win32.LPARAM(uintptr(unsafe.Pointer(p)))

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
