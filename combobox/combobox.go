package combobox

import (
	"structs"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw"
	"github.com/mkch/gw/listbox/listcombo"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const (
	CB_OKAY     = listcombo.LC_OKAY
	CB_ERR      = listcombo.LC_ERR
	CB_ERRSPACE = listcombo.LC_ERRSPACE
)

type NotifyCode win32.WORD

const (
	CBN_ERRSPACE     NotifyCode = 0xFFFF
	CBN_SELCHANGE    NotifyCode = 1
	CBN_DBLCLK       NotifyCode = 2
	CBN_SETFOCUS     NotifyCode = 3
	CBN_KILLFOCUS    NotifyCode = 4
	CBN_EDITCHANGE   NotifyCode = 5
	CBN_EDITUPDATE   NotifyCode = 6
	CBN_DROPDOWN     NotifyCode = 7
	CBN_CLOSEUP      NotifyCode = 8
	CBN_SELENDOK     NotifyCode = 9
	CBN_SELENDCANCEL NotifyCode = 10
)

const (
	CBS_SIMPLE            win32.WindowStyle = 0x0001
	CBS_DROPDOWN          win32.WindowStyle = 0x0002
	CBS_DROPDOWNLIST      win32.WindowStyle = 0x0003
	CBS_OWNERDRAWFIXED    win32.WindowStyle = 0x0010
	CBS_OWNERDRAWVARIABLE win32.WindowStyle = 0x0020
	CBS_AUTOHSCROLL       win32.WindowStyle = 0x0040
	CBS_OEMCONVERT        win32.WindowStyle = 0x0080
	CBS_SORT              win32.WindowStyle = 0x0100
	CBS_HASSTRINGS        win32.WindowStyle = 0x0200
	CBS_NOINTEGRALHEIGHT  win32.WindowStyle = 0x0400
	CBS_DISABLENOSCROLL   win32.WindowStyle = 0x0800
	CBS_UPPERCASE         win32.WindowStyle = 0x2000
	CBS_LOWERCASE         win32.WindowStyle = 0x4000
)

// Combo Box messages
const (
	CB_GETEDITSEL            = 0x0140
	CB_LIMITTEXT             = 0x0141
	CB_SETEDITSEL            = 0x0142
	CB_ADDSTRING             = 0x0143
	CB_DELETESTRING          = 0x0144
	CB_DIR                   = 0x0145
	CB_GETCOUNT              = 0x0146
	CB_GETCURSEL             = 0x0147
	CB_GETLBTEXT             = 0x0148
	CB_GETLBTEXTLEN          = 0x0149
	CB_INSERTSTRING          = 0x014A
	CB_RESETCONTENT          = 0x014B
	CB_FINDSTRING            = 0x014C
	CB_SELECTSTRING          = 0x014D
	CB_SETCURSEL             = 0x014E
	CB_SHOWDROPDOWN          = 0x014F
	CB_GETITEMDATA           = 0x0150
	CB_SETITEMDATA           = 0x0151
	CB_GETDROPPEDCONTROLRECT = 0x0152
	CB_SETITEMHEIGHT         = 0x0153
	CB_GETITEMHEIGHT         = 0x0154
	CB_SETEXTENDEDUI         = 0x0155
	CB_GETEXTENDEDUI         = 0x0156
	CB_GETDROPPEDSTATE       = 0x0157
	CB_FINDSTRINGEXACT       = 0x0158
	CB_SETLOCALE             = 0x0159
	CB_GETLOCALE             = 0x015A
	CB_GETTOPINDEX           = 0x015b
	CB_SETTOPINDEX           = 0x015c
	CB_GETHORIZONTALEXTENT   = 0x015d
	CB_SETHORIZONTALEXTENT   = 0x015e
	CB_GETDROPPEDWIDTH       = 0x015f
	CB_SETDROPPEDWIDTH       = 0x0160
	CB_INITSTORAGE           = 0x0161
	CB_MULTIPLEADDSTRING     = 0x0163
	CB_GETCOMBOBOXINFO       = 0x0164
	CB_MSGMAX                = 0x0165

	CBM_FIRST        = 0x1700 // Combobox control messages
	CB_SETMINVISIBLE = CBM_FIRST + 1
	CB_GETMINVISIBLE = CBM_FIRST + 2
	CB_SETCUEBANNER  = CBM_FIRST + 3
	CB_GETCUEBANNER  = CBM_FIRST + 4
)

type COMBOBOXINFO struct {
	_           structs.HostLayout
	Size        win32.DWORD
	RcItem      win32.RECT
	RcButton    win32.RECT
	StateButton win32.SystemState
	HwndCombo   win32.HWND
	HwndItem    win32.HWND
	HwndList    win32.HWND
}

type Spec struct {
	Parent              gw.BaseWindow
	Style               win32.WindowStyle
	ExStyle             win32.WindowExStyle
	X, Y, Width, Height metrics.Dimension

	// OnCompareItem is called in [ListBox.OnCompareItem] if not nil.
	OnCompareItem func(item1, item2 any, locale win32.DWORD) int
	// OnMeasureItem is called in [ListBox.OnMeasureItem] if not nil.
	OnMeasureItem func(index int, itemData any) (width, height int)
	// OnDrawItem is called in [ListBox.OnDrawItem] if not nil.
	OnDrawItem func(info *listcombo.DrawItemInfo)
}

type ComboBox struct {
	listcombo.ListCombo
	// Spec is used to create the window and is cleared after creation.
	Spec *Spec
}

func (l *ComboBox) WndProc(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	switch msg {
	case win32.WM_COMPAREITEM, win32.WM_MEASUREITEM, win32.WM_DRAWITEM:
		// A ComboBox acts as the parent window of its internal ListBox and receives these messages sent from the ListBox.
		// At this point, the CtlType and HwndItem fields within these messages are specific to the ListBox.
		// When these messages are passed to the ComboBox's DefWndProc, they are converted into messages targeting the ComboBox
		// itself and are forwarded to the ComboBox's parent window.
		// Upon receiving these messages, the parent window of ComboBox can perform proper message reflection based on CtlType and HwndItem.
		return l.ListCombo.DefWndProc(hwnd, msg, wParam, lParam)
	}
	return l.ListCombo.WndProc(hwnd, msg, wParam, lParam)
}

func (l *ComboBox) OnInit() error {
	defer func() { l.Spec = nil }()

	return l.ListCombo.OnInit(listcombo.Config{
		MsgAddString:           CB_ADDSTRING,
		MsgGetItemText:         CB_GETLBTEXT,
		MsgGetItemTextLen:      CB_GETLBTEXTLEN,
		MsgInsertString:        CB_INSERTSTRING,
		MsgDeleteString:        CB_DELETESTRING,
		MsgResetContent:        CB_RESETCONTENT,
		MsgGetItemCount:        CB_GETCOUNT,
		MsgGetItemData:         CB_GETITEMDATA,
		MsgSetItemData:         CB_SETITEMDATA,
		MsgGetCurSel:           CB_GETCURSEL,
		MsgSetCurSel:           CB_SETCURSEL,
		MsgFindString:          CB_FINDSTRING,
		MsgFindStringExact:     CB_FINDSTRINGEXACT,
		MsgSetHorizontalExtent: CB_SETHORIZONTALEXTENT,
		MsgGetHorizontalExtent: CB_GETHORIZONTALEXTENT,
		MsgSetItemHeight:       CB_SETITEMHEIGHT,
		MsgGetItemHeight:       CB_GETITEMHEIGHT,
		MsgGetLocale:           CB_GETLOCALE,
		MsgSetLocale:           CB_SETLOCALE,
		MsgGetTopIndex:         CB_GETTOPINDEX,
		MsgSetTopIndex:         CB_SETTOPINDEX,
		MsgInitStorage:         CB_INITSTORAGE,
		MsgSelectString:        CB_SELECTSTRING,

		StyleOwnerDrawFixed:    CBS_OWNERDRAWFIXED,
		StyleOwnerDrawVariable: CBS_OWNERDRAWVARIABLE,
		StyleHasStrings:        CBS_HASSTRINGS,
		StyleSort:              CBS_SORT,

		OnCompareItem: l.Spec.OnCompareItem,
		OnMeasureItem: l.Spec.OnMeasureItem,
		OnDrawItem:    l.Spec.OnDrawItem,
	})
}

func (l *ComboBox) CreateHandle() (win32.HWND, error) {
	dpi, err := win32.GetDpiForWindow(l.Spec.Parent.HWND())
	if err != nil {
		return 0, err
	}
	return win32util.CreateWindow(&win32util.Wnd{
		ClassName: "ComboBox",
		WndParent: l.Spec.Parent.HWND(),
		X:         metrics.ToPx(l.Spec.X, dpi).Value(),
		Y:         metrics.ToPx(l.Spec.Y, dpi).Value(),
		Width:     metrics.ToPx(l.Spec.Width, dpi).Value(),
		Height:    metrics.ToPx(l.Spec.Height, dpi).Value(),
		Style:     l.Spec.Style | win32.WS_CHILD,
		ExStyle:   l.Spec.ExStyle,
	})
}

// ComboBoxInfo returns the information about the specified combo box.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-getcomboboxinfo for details.
func (l *ComboBox) ComboBoxInfo() (*COMBOBOXINFO, error) {
	var info COMBOBOXINFO
	info.Size = win32.DWORD(uintptr(unsafe.Sizeof(info)))
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_GETCOMBOBOXINFO, 0, uintptr(unsafe.Pointer(&info)))
	if err != nil {
		return nil, err
	}
	if i == 0 {
		return nil, listcombo.ErrFailed
	}
	return &info, nil
}

// SetCueBanner sets the cue banner text that is displayed for the edit control of a combo box.
//
// Need comctl32.dll version 6 in the manifest.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-setcuebanner for details.
func (l *ComboBox) SetCueBanner(text string) error {
	var buf []win32.WCHAR
	win32util.CString(text, &buf)
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_SETCUEBANNER, 0, unsafe.Pointer(&buf[0]))
	if err != nil {
		return err
	}
	if i != 1 {
		return listcombo.ErrFailed
	}
	return nil
}

// CueBanner returns the cue banner text displayed in the edit control of a combo box.
//
// Need comctl32.dll version 6 in the manifest.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-setcuebanner for details.
func (l *ComboBox) CueBanner() (string, error) {
	var buf [256]win32.WCHAR
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_GETCUEBANNER, unsafe.Pointer(&buf[0]), len(buf))
	if err != nil {
		return "", err
	}
	if i == 0 {
		return "", nil
	}
	return win32util.GoString(&buf[0], win32util.CStrLen(&buf[0], len(buf))+1), nil
}

// DroppedCtrlRect retrieves the screen coordinates of a combo box in its dropped-down state.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-getdroppedcontrolrect for details.
func (l *ComboBox) DroppedCtrlRect() (*win32.RECT, error) {
	var rect win32.RECT
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_GETDROPPEDCONTROLRECT, 0, unsafe.Pointer(&rect))
	if err != nil {
		return nil, err
	}
	if i == 0 {
		return nil, listcombo.ErrFailed
	}
	return &rect, nil
}

// DroppedDown returns whether the combo box is in the dropped-down state.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-getdroppedstate for details.
func (l *ComboBox) DroppedDown() (bool, error) {
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_GETDROPPEDSTATE, 0, 0)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

// MinDroppedWidth retrieves the minimum allowable width, in pixels, of the list box of a combo box with the CBS_DROPDOWN or CBS_DROPDOWNLIST style.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-getdroppedwidth for details.
func (l *ComboBox) MinDroppedWidth() (int, error) {
	return listcombo.SendMessageRetNoError[int](l.HWND(), CB_GETDROPPEDWIDTH, 0, 0)
}

// SetMinDroppedWidth set the minimum allowable width, in pixels, of the list box of a combo box with the CBS_DROPDOWN or CBS_DROPDOWNLIST style.
//
// The width of the list box is either the minimum allowable width or
// the combo box width, whichever is larger.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-setdroppedwidth for details.
func (l *ComboBox) SetMinDroppedWidth(width int) error {
	return listcombo.SendMessageNoError(l.HWND(), CB_SETDROPPEDWIDTH, width, 0)
}

// SetEditSel select characters in the edit control of a combo box.
// If start is -1, the selection, if any, is removed.
// If end is -1, all text from the starting position to the last character in the edit control is selected.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-seteditsel for details.
func (l *ComboBox) SetEditSel(start, end int) error {
	return listcombo.SendMessageNoError(l.HWND(), CB_SETEDITSEL, 0, win32.MAKELONG(uint16(start), uint16(end)))
}

// EditSel returns the starting and ending character positions of the current selection in the edit control of a combo box.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-geteditsel for details.
func (l *ComboBox) EditSel() (start, end int, err error) {
	var sel win32.DWORD
	sel, err = listcombo.SendMessageRet[win32.DWORD](l.HWND(), CB_GETEDITSEL, 0, 0)
	if err != nil {
		return 0, 0, err
	}
	start = int(win32.LOWORD(sel))
	end = int(win32.HIWORD(sel))
	return start, end, nil
}

// ExtendedUI returns whether a combo box has the extended user interface rather than the default user interface.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-getextendedui for details.
func (l *ComboBox) ExtendedUI() (bool, error) {
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_GETEXTENDEDUI, 0, 0)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

// SetExtendedUI sets whether a combo box has the extended user interface rather than the default user interface.
//
// By default, the F4 key opens or closes the list and the DOWN ARROW changes the current selection.
// In a combo box with the extended user interface, the F4 key is disabled and pressing the DOWN ARROW key opens the drop-down list.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-setextendedui for details.
func (l *ComboBox) SetExtendedUI(extended bool) error {
	return listcombo.SendMessageNoError(l.HWND(), CB_SETEXTENDEDUI, gg.If(extended, 1, 0), 0)
}

// SetMinVisible sets  set the minimum number of visible items in the drop-down list of a combo box.
// When the number of items in the drop-down list is greater than the minimum, the combo box uses a scroll bar.
//
// Need comctl32.dll version 6 in the manifest.
//
// This method is ignored if the combo box control has style CBS_NOINTEGRALHEIGHT.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-setminvisible for details.
func (l *ComboBox) SetMinVisible(items int) error {
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_SETMINVISIBLE, items, 0)
	if err != nil {
		return err
	}
	if i == 0 {
		return listcombo.ErrFailed
	}
	return nil
}

// MinVisible returns the minimum number of visible items in the drop-down list of a combo box.
//
// Need comctl32.dll version 6 in the manifest.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-getminvisible for details.
func (l *ComboBox) MinVisible() (int, error) {
	return listcombo.SendMessageRet[int](l.HWND(), CB_GETMINVISIBLE, 0, 0)
}

// SetTextLimit limits the length of the text the user may type into the edit control of a combo box.
//
// If limit is zero, the text length is limited to 0x7FFFFFFE characters.
//
// If the combo box does not have the CBS_AUTOHSCROLL style, setting the text limit to be larger than the size of the edit control has no effect.
//
// The method limits only the text the user can enter. It has no effect on any text already in the edit control when this method is called,
// nor does it affect the length of the text copied to the edit control when a string in the list box is selected.
//
// The default limit to the text a user can enter in the edit control is 30,000 TCHARs.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-limittext for details.
func (l *ComboBox) SetTextLimit(limit int) error {
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_LIMITTEXT, limit, 0)
	if err != nil {
		return err
	}
	if i == 0 {
		return listcombo.ErrFailed
	}
	return nil
}

// ShowDropDown shows or hides the list box of a combo box that has the CBS_DROPDOWN or CBS_DROPDOWNLIST style.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/cb-showdropdown for details.
func (l *ComboBox) ShowDropDown(show bool) error {
	_, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), CB_SHOWDROPDOWN, gg.If(show, 1, 0), 0)
	return err
}

// New creates a new ListBox control with the specified specification.
func New(spec *Spec) (*ComboBox, error) {
	return gw.Init(&ComboBox{Spec: spec})
}
