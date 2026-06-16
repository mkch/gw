package listbox

import (
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw"
	"github.com/mkch/gw/listbox/listcombo"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

type NotifyCode win32.WORD

const (
	LBN_ERRSPACE  NotifyCode = 0xFFFE // -2
	LBN_SELCHANGE NotifyCode = 1
	LBN_DBLCLK    NotifyCode = 2
	LBN_SELCANCEL NotifyCode = 3
	LBN_SETFOCUS  NotifyCode = 4
	LBN_KILLFOCUS NotifyCode = 5
)

// Listbox messages
const (
	LB_ADDSTRING           = 0x0180
	LB_INSERTSTRING        = 0x0181
	LB_DELETESTRING        = 0x0182
	LB_SELITEMRANGEEX      = 0x0183
	LB_RESETCONTENT        = 0x0184
	LB_SETSEL              = 0x0185
	LB_SETCURSEL           = 0x0186
	LB_GETSEL              = 0x0187
	LB_GETCURSEL           = 0x0188
	LB_GETTEXT             = 0x0189
	LB_GETTEXTLEN          = 0x018A
	LB_GETCOUNT            = 0x018B
	LB_SELECTSTRING        = 0x018C
	LB_DIR                 = 0x018D
	LB_GETTOPINDEX         = 0x018E
	LB_FINDSTRING          = 0x018F
	LB_GETSELCOUNT         = 0x0190
	LB_GETSELITEMS         = 0x0191
	LB_SETTABSTOPS         = 0x0192
	LB_GETHORIZONTALEXTENT = 0x0193
	LB_SETHORIZONTALEXTENT = 0x0194
	LB_SETCOLUMNWIDTH      = 0x0195
	LB_ADDFILE             = 0x0196
	LB_SETTOPINDEX         = 0x0197
	LB_GETITEMRECT         = 0x0198
	LB_GETITEMDATA         = 0x0199
	LB_SETITEMDATA         = 0x019A
	LB_SELITEMRANGE        = 0x019B
	LB_SETANCHORINDEX      = 0x019C
	LB_GETANCHORINDEX      = 0x019D
	LB_SETCARETINDEX       = 0x019E
	LB_GETCARETINDEX       = 0x019F
	LB_SETITEMHEIGHT       = 0x01A0
	LB_GETITEMHEIGHT       = 0x01A1
	LB_FINDSTRINGEXACT     = 0x01A2
	LB_SETLOCALE           = 0x01A5
	LB_GETLOCALE           = 0x01A6
	LB_SETCOUNT            = 0x01A7
	LB_INITSTORAGE         = 0x01A8
	LB_ITEMFROMPOINT       = 0x01A9
	LB_MULTIPLEADDSTRING   = 0x01B1

	LB_GETLISTBOXINFO = 0x01B2
	LB_MSGMAX         = 0x01B3
)

// Listbox Styles
const (
	LBS_NOTIFY            win32.WindowStyle = 0x0001
	LBS_SORT              win32.WindowStyle = 0x0002
	LBS_NOREDRAW          win32.WindowStyle = 0x0004
	LBS_MULTIPLESEL       win32.WindowStyle = 0x0008
	LBS_OWNERDRAWFIXED    win32.WindowStyle = 0x0010
	LBS_OWNERDRAWVARIABLE win32.WindowStyle = 0x0020
	LBS_HASSTRINGS        win32.WindowStyle = 0x0040
	LBS_USETABSTOPS       win32.WindowStyle = 0x0080
	LBS_NOINTEGRALHEIGHT  win32.WindowStyle = 0x0100
	LBS_MULTICOLUMN       win32.WindowStyle = 0x0200
	LBS_WANTKEYBOARDINPUT win32.WindowStyle = 0x0400
	LBS_EXTENDEDSEL       win32.WindowStyle = 0x0800
	LBS_DISABLENOSCROLL   win32.WindowStyle = 0x1000
	LBS_NODATA            win32.WindowStyle = 0x2000
	LBS_NOSEL             win32.WindowStyle = 0x4000
	LBS_COMBOBOX          win32.WindowStyle = 0x8000
	LBS_STANDARD          win32.WindowStyle = (LBS_NOTIFY | LBS_SORT | win32.WS_VSCROLL | win32.WS_BORDER)
)

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

type ListBox struct {
	listcombo.ListCombo
	// Spec is used to create the window and is cleared after creation.
	Spec *Spec
}

func (l *ListBox) OnInit() error {
	defer func() { l.Spec = nil }()

	return l.ListCombo.OnInit(listcombo.Config{
		MsgAddString:           LB_ADDSTRING,
		MsgGetItemText:         LB_GETTEXT,
		MsgGetItemTextLen:      LB_GETTEXTLEN,
		MsgInsertString:        LB_INSERTSTRING,
		MsgDeleteString:        LB_DELETESTRING,
		MsgResetContent:        LB_RESETCONTENT,
		MsgGetItemCount:        LB_GETCOUNT,
		MsgGetItemData:         LB_GETITEMDATA,
		MsgSetItemData:         LB_SETITEMDATA,
		MsgGetCurSel:           LB_GETCURSEL,
		MsgSetCurSel:           LB_SETCURSEL,
		MsgFindString:          LB_FINDSTRING,
		MsgFindStringExact:     LB_FINDSTRINGEXACT,
		MsgSetHorizontalExtent: LB_SETHORIZONTALEXTENT,
		MsgGetHorizontalExtent: LB_GETHORIZONTALEXTENT,
		MsgSetItemHeight:       LB_SETITEMHEIGHT,
		MsgGetItemHeight:       LB_GETITEMHEIGHT,
		MsgGetLocale:           LB_GETLOCALE,
		MsgSetLocale:           LB_SETLOCALE,
		MsgGetTopIndex:         LB_GETTOPINDEX,
		MsgSetTopIndex:         LB_SETTOPINDEX,
		MsgInitStorage:         LB_INITSTORAGE,
		MsgSelectString:        LB_SELECTSTRING,

		StyleOwnerDrawFixed:    LBS_OWNERDRAWFIXED,
		StyleOwnerDrawVariable: LBS_OWNERDRAWVARIABLE,
		StyleHasStrings:        LBS_HASSTRINGS,
		StyleSort:              LBS_SORT,

		OnCompareItem: l.Spec.OnCompareItem,
		OnMeasureItem: l.Spec.OnMeasureItem,
		OnDrawItem:    l.Spec.OnDrawItem,
	})
}

func (l *ListBox) CreateHandle() (win32.HWND, error) {
	dpi, err := win32.GetDpiForWindow(l.Spec.Parent.HWND())
	if err != nil {
		return 0, err
	}
	return win32util.CreateWindow(&win32util.Wnd{
		ClassName: "ListBox",
		WndParent: l.Spec.Parent.HWND(),
		X:         metrics.ToPx(l.Spec.X, dpi).Value(),
		Y:         metrics.ToPx(l.Spec.Y, dpi).Value(),
		Width:     metrics.ToPx(l.Spec.Width, dpi).Value(),
		Height:    metrics.ToPx(l.Spec.Height, dpi).Value(),
		Style:     l.Spec.Style | win32.WS_CHILD,
		ExStyle:   l.Spec.ExStyle,
	})
}

// AnchorIndex returns the index of the anchor item that is, the item from which a multiple selection starts.
// A multiple selection spans all items from the anchor item to the caret item.
// If there is no anchor item, it returns -1, nil.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getanchorindex for details.
func (l *ListBox) AnchorIndex() (int, error) {
	return listcombo.SendMessageRet[int](l.HWND(), LB_GETANCHORINDEX, 0, 0)
}

// SetAnchorIndex sets the index of the anchor item in a multiple-selection list box.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setanchorindex for details.
func (l *ListBox) SetAnchorIndex(index int) error {
	return listcombo.SendMessageNoError(l.HWND(), LB_SETANCHORINDEX, index, 0)
}

// CaretIndex retrieves the index of the item that has the focus in a multiple-selection list box.
// The item may or may not be selected.
// If no item has the focus, it returns 0, nil.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getcaretindex for details.
func (l *ListBox) CaretIndex() (int, error) {
	return listcombo.SendMessageRet[int](l.HWND(), LB_GETCARETINDEX, 0, 0)
}

// SetCaretIndex sets the focus rectangle to the item at the specified index in a multiple-selection list box.
// If the item is not visible, it is scrolled into view.
//
// If this method is called on a single-selection list box that does not contain a selected item,
// the caret index is set to the item specified by the wParam parameter.
// If the single-selection list box does contain a selected item, the list box returns LB_ERR.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setcaretindex for details.
func (l *ListBox) SetCaretIndex(index int) error {
	return listcombo.SendMessageNoError(l.HWND(), LB_SETCARETINDEX, index, 0)
}

// Selected returns the selection state of an item.
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getsel for details.
func (l *ListBox) Selected(index int) (bool, error) {
	i, err := listcombo.SendMessageRetNoError[win32.LRESULT](l.HWND(), LB_GETSEL, index, 0)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

// SelectedItems returns the indexes of all selected items in a multiple-selection list box or combo box, or an error if the operation failed.
// If the list box or combo box is not a multiple-selection control, it returns [ErrFailed].
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getselcount and https://learn.microsoft.com/en-us/windows/win32/controls/lb-getselitems for details.
func (l *ListBox) SelectedItems() (indexes []int32, err error) {
	i, err := listcombo.SendMessageRetNoError[int](l.HWND(), LB_GETSELCOUNT, 0, 0)
	if err != nil {
		return
	}
	if i == 0 {
		return nil, nil
	}
	indexes = make([]int32, i)
	if err = listcombo.SendMessageNoError(l.HWND(), LB_GETSELITEMS, i, unsafe.Pointer(&indexes[0])); err != nil {
		return
	}
	return
}

// SetSelected selects or unselects an item in a multiple-selection list box or combo box and, if necessary, scrolls the item into view.
// The index parameter specifies the zero-based index of the item to select or unselect. If index is -1, all items are applied.
// If selected is true, the item or items are selected; if false, they are unselected.
//
// If the list box or combo box is not a multiple-selection control, it returns [ErrFailed].
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setsel for details.
func (l *ListBox) SetSelected(index int, selected bool) error {
	sel := gg.If(selected, 1, 0)
	return listcombo.SendMessageNoError(l.HWND(), LB_SETSEL, sel, index)
}

// ItemRect returns the dimensions of the rectangle that bounds a list box item as it is currently displayed in the list box.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getitemrect for details.
func (l *ListBox) ItemRect(itemIndex int) (*win32.RECT, error) {
	var rect win32.RECT
	if err := listcombo.SendMessageNoError(l.HWND(), LB_GETITEMRECT, itemIndex, unsafe.Pointer(&rect)); err != nil {
		return nil, err
	}
	return &rect, nil
}

//	ItemsPerColumn returns the number of items per column in a specified list box.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getlistboxinfo for details.
func (l *ListBox) ItemsPerColumn() (int, error) {
	return listcombo.SendMessageRet[int](l.HWND(), LB_GETLISTBOXINFO, 0, 0)
}

// SetColumnWidth sets the width, in pixels, of all columns in a multiple-column list box.
func (l *ListBox) SetColumnWidth(width int) error {
	return listcombo.SendMessageNoError(l.HWND(), LB_SETCOLUMNWIDTH, width, 0)
}

// ItemFromPoint returns the zero-based index of the item nearest the specified point in a list box.
// The point is specified in client coordinates of l.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-itemfrompoint for details.
func (l *ListBox) ItemFromPoint(x, y int) (index int, err error) {
	i, err := listcombo.SendMessageRet[win32.LRESULT](l.HWND(), LB_ITEMFROMPOINT, 0, win32.MAKELONG(uint16(x), uint16(y)))
	if err != nil {
		return 0, err
	}
	return int(win32.LOWORD(i)), nil
}

// SelectRange selects or deselects one or more consecutive items in a multiple-selection list box.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-selitemrange for details.
func (l *ListBox) SelectRange(selected bool, startIndex, endIndex int) error {
	sel := gg.If(selected, 1, 0)
	return listcombo.SendMessageNoError(l.HWND(), LB_SELITEMRANGE, sel, win32.MAKELONG(uint16(startIndex), uint16(endIndex)))
}

// SetCount sets the count of items in a list box created with the LBS_NODATA style and not created with the LBS_HASSTRINGS style.
// This method is supported only by list boxes created with the LBS_NODATA style and not created with the LBS_HASSTRINGS style.
// All other list boxes return [listcombo.ErrFailed].
//
// Note: A no-data list box must also have the LBS_OWNERDRAWFIXED style, but must not have the LBS_SORT or LBS_HASSTRINGS style.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setcount for details.
func (l *ListBox) SetCount(count int) error {
	return listcombo.SendMessageOkay(l.HWND(), LB_SETCOUNT, count, 0)
}

// SetTabStops sets the tab-stop positions in a list box.
// To respond to this method, the list box must have been created with the LBS_USETABSTOPS style.
//
// The elements of tabStops represent the number of quarters of the average character width for the font that is selected into the list box.
// For example, a tab stop of 4 is placed at 1.0 character units, and a tab stop of 6 is placed at 1.5 average character units.
// However, if the list box is part of a dialog box, the integers are in dialog template units.
//
// The tab stops must be sorted in ascending order; backward tabs are not allowed.
// If tabStops is empty, the default tab stop is two dialog template units.
// If tabStops has only one element, the list box will have tab stops separated by the distance specified by tapStops[0].
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-settabstops for details.
func (l *ListBox) SetTabStops(tabStops []int32) error {
	if len(tabStops) == 0 {
		return listcombo.SendMessageNoError(l.HWND(), LB_SETTABSTOPS, 0, 0)
	}
	tabStops = win32.AlignSlice(tabStops)
	return listcombo.SendMessageNoError(l.HWND(), LB_SETTABSTOPS, len(tabStops), unsafe.Pointer(&tabStops[0]))
}

// New creates a new ListBox control with the specified specification.
func New(spec *Spec) (*ListBox, error) {
	return gw.Init(&ListBox{Spec: spec})
}
