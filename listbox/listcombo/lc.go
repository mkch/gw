// package listcombo provides common functionality for listbox and combobox controls.
package listcombo

import (
	"errors"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw"
	"github.com/mkch/gw/control"
	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/util"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const (
	LC_OKAY     = 0
	LC_ERR      = -1
	LC_ERRSPACE = -2
)

// ErrItemStringNotSupported is returned when trying to use item strings with an owner-drawn listbox or combobox that does not have the HASSTRINGS style.
var ErrItemStringNotSupported = errors.New("owner-drawn listbox or combobox without HASSTRINGS style cannot use item string")

// ErrExternalItemData is returned when trying to use item data that was not set by gw itself, which may cause memory safety issues.
var ErrExternalItemData = errors.New("external item data")

// ErrFindStringNotSupported is returned when trying to find an item by string in an owner-drawn listbox or combobox without the HASSTRINGS and SORT styles.
var ErrFindStringNotSupported = errors.New("finding item by string is not supported for owner-drawn listbox or combobox without HASSTRINGS and SORT style")

// ResultError represents an error result from a listbox or combobox operation.
type ResultError int

func (e ResultError) Error() string {
	if e == LC_OKAY {
		return "operation succeeded"
	}
	switch e {
	case LC_ERR:
		return "operation failed"
	case LC_ERRSPACE:
		return "insufficient space to store the new item"
	default:
		return "unknown error"
	}
}

// ErrFailed is returned when a listbox or combobox operation failed.
const ErrFailed = ResultError(LC_ERR)

// ErrSpace is returned when a listbox or combobox operation failed due to insufficient space to store the new item.
const ErrSpace = ResultError(LC_ERRSPACE)

type Config struct {
	MsgAddString           win32.UINT
	MsgGetItemText         win32.UINT
	MsgGetItemTextLen      win32.UINT
	MsgInsertString        win32.UINT
	MsgDeleteString        win32.UINT
	MsgResetContent        win32.UINT
	MsgGetItemCount        win32.UINT
	MsgGetItemData         win32.UINT
	MsgSetItemData         win32.UINT
	MsgGetCurSel           win32.UINT
	MsgSetCurSel           win32.UINT
	MsgFindString          win32.UINT
	MsgFindStringExact     win32.UINT
	MsgSetHorizontalExtent win32.UINT
	MsgGetHorizontalExtent win32.UINT
	MsgSetItemHeight       win32.UINT
	MsgGetItemHeight       win32.UINT
	MsgGetLocale           win32.UINT
	MsgSetLocale           win32.UINT
	MsgGetTopIndex         win32.UINT
	MsgSetTopIndex         win32.UINT
	MsgInitStorage         win32.UINT
	MsgSelectString        win32.UINT

	StyleOwnerDrawFixed    win32.WindowStyle
	StyleOwnerDrawVariable win32.WindowStyle
	StyleHasStrings        win32.WindowStyle
	StyleSort              win32.WindowStyle

	OnCompareItem func(item1, item2 any, locale win32.DWORD) int
	OnMeasureItem func(index int, itemData any) (width, height int)
	OnDrawItem    func(info *DrawItemInfo)
}

// DrawItemInfo is the information for drawing an item in an owner-drawn listbox or combobox.
// See [ListCombo.OnDrawItem] for details.
type DrawItemInfo struct {
	// Index of the item to draw or -1 if the control is not an item.
	Index      int
	ItemAction win32.OwnerDrawAction
	ItemState  win32.OwnerDrawState
	HwndItem   win32.HWND
	HDC        win32.HDC
	RcItem     win32.RECT
	// ItemData is the item data of the item to draw. See [ListCombo.ItemData] for details.
	ItemData any
}

// ListCombo represents the common functionality for listbox and combobox controls.
type ListCombo struct {
	control.Control
	config         Config
	pinnedItemData gg.Set[*ItemData]
	// Whether the listbox or combobox is an owner-drawn control.
	ownerDraw bool
	// whether the listbox or combobox is an owner-drawn control without the HASSTRINGS style,
	// which cannot use item strings.
	noItemStrings bool
	//Whether SORT style is set.
	sorted bool
}

func (l *ListCombo) OnInit(spec Config) (err error) {
	err = l.Control.OnInit()
	if err != nil {
		return
	}

	l.config = spec

	hwnd := l.HWND()

	ls, err := win32.GetWindowLongPtrW(hwnd, win32.GWL_STYLE)
	if err != nil {
		return
	}
	style := win32.WindowStyle(ls)

	ownerDraw := style&l.config.StyleOwnerDrawFixed != 0 || style&l.config.StyleOwnerDrawVariable != 0
	l.ownerDraw = ownerDraw
	l.noItemStrings = ownerDraw && style&l.config.StyleHasStrings == 0
	l.sorted = style&l.config.StyleSort != 0

	if ownerDraw {
		var p win32.HWND
		p, err = win32.GetAncestor(hwnd, win32.GA_PARENT)
		if err != nil {
			return
		}
		if err = gw.LookupWindow(p).AssignID(hwnd); err != nil {
			return
		}
	}

	return
}

func (l *ListCombo) OnDestroy() {
	for p := range l.pinnedItemData {
		p.Unpin()
	}
	if l.ownerDraw {
		hwnd := l.HWND()
		p, err := win32.GetAncestor(hwnd, win32.GA_PARENT)
		if err != nil {
			panic(err)
		}
		gw.LookupWindow(p).RemoveID(hwnd)
	}
	l.Control.OnDestroy()
}

func (l *ListCombo) WndProc(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) (result win32.LRESULT) {
	switch msg {
	case appmsg.REFLECT_COMPAREITEM:
		cmp := (*win32.COMPAREITEMSTRUCT)(unsafe.Add(nil, lParam))
		if c, ok := gw.LookupWindow(l.HWND()).(interface {
			OnCompareItem(item1, item2 any, locale win32.DWORD) int
		}); ok {
			return win32.LRESULT(c.OnCompareItem(l.ItemDataFromRaw(uintptr(cmp.ItemData1)), l.ItemDataFromRaw(uintptr(cmp.ItemData2)), cmp.LocaleId))
		}
	case appmsg.REFLECT_MEASUREITEM:
		ms := (*win32.MEASUREITEMSTRUCT)(unsafe.Add(nil, lParam))
		if m, ok := gw.LookupWindow(l.HWND()).(interface {
			OnMeasureItem(index int, itemData any) (width, height int)
		}); ok {
			width, height := m.OnMeasureItem(int(ms.ItemID), l.ItemDataFromRaw(uintptr(ms.ItemData)))
			ms.ItemWidth = win32.UINT(width)
			ms.ItemHeight = win32.UINT(height)
			return 1
		}
	case appmsg.REFLECT_DRAWITEM:
		ds := (*win32.DRAWITEMSTRUCT)(unsafe.Add(nil, lParam))
		if d, ok := gw.LookupWindow(l.HWND()).(interface{ OnDrawItem(info *DrawItemInfo) }); ok {
			d.OnDrawItem(&DrawItemInfo{
				Index:      int(ds.ItemID),
				ItemAction: ds.ItemAction,
				ItemState:  ds.ItemState,
				HwndItem:   ds.HwndItem,
				HDC:        ds.HDC,
				RcItem:     ds.RcItem,
				ItemData:   l.ItemDataFromRaw(uintptr(ds.ItemData)),
			})
			return 1
		}
	}
	return l.Control.WndProc(hwnd, msg, wParam, lParam)
}

// OnCompareItem is called when an new item is added to or FindItemString method is called on an
// owner-drawn control with LBS_SORT style but without LBS_HASSTRINGS style.
// The parameter item1 and item2 are the item data of the two items to compare. See [ListCombo.ItemData] for details.
// The parameter locale is the current locale id of the list box or combo box.
// It should return -1 if item1 is less than item2, 0 if they are equal, and 1 if item1 is greater than item2.
// The default implementation calls the item comparator(See [ListCombo.SetItemComparator]) if it is not nil, otherwise it returns 0.
func (l *ListCombo) OnCompareItem(item1, item2 any, locale win32.DWORD) int {
	if l.config.OnCompareItem == nil {
		return 0
	}
	return l.config.OnCompareItem(item1, item2, locale)
}

// OnMeasureItem is called when the control needs to know the size of an item in an owner-drawn control.
// The parameter index is the index of the item to measure or 0 if the control is ownerdraw-fixed.
// The parameter itemData is the item data of the item to measure. See [ListCombo.ItemData] for details.
// It should return the width and height of the item in pixels.
// The default implementation calls the OnMeasureItem of the [Config] if it is not nil, otherwise it returns 0, 0.
func (l *ListCombo) OnMeasureItem(index int, itemData any) (width, height int) {
	if l.config.OnMeasureItem == nil {
		return 0, 0
	}
	return l.config.OnMeasureItem(index, itemData)
}

// OnDrawItem is called when the control needs to draw an item in an owner-drawn control.
// The default implementation calls the OnDrawItem of the [Config] if it is not nil.
func (l *ListCombo) OnDrawItem(info *DrawItemInfo) {
	if l.config.OnDrawItem == nil {
		return
	}
	l.config.OnDrawItem(info)
}

// ItemData is the data associated with an item in a listbox or combobox.
type ItemData = util.DataPinner[any]

func (l *ListCombo) checkCanUseItemString() error {
	if l.noItemStrings {
		return ErrItemStringNotSupported
	}
	return nil
}

// AppendItem adds an item with the specified data to a listbox or combobox and returns any error that occurred.
// If the listbox or combobox is not an owner-drawn control without the HASSTRINGS style, the newly added item
// will have an empty string.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-addstring for details.
func (l *ListCombo) AppendItem(data any) error {
	if !l.noItemStrings {
		i, err := l.addString("") // Must add an item before setting item data
		if err != nil {
			return err
		}
		return l.SetItemData(i, data)
	}
	var d = &ItemData{Data: &data}
	d.Pin()
	if l.pinnedItemData == nil {
		l.pinnedItemData = make(gg.Set[*ItemData])
	}
	l.pinnedItemData.Add(d)
	return SendMessageOkay(l.HWND(), l.config.MsgAddString, 0, unsafe.Pointer(d))
}

// AppendItemString adds a string to a listbox or combobox and returns the index of the new item, or an error if the operation failed.
//
// // See https://learn.microsoft.com/en-us/windows/win32/controls/lb-addstring for details.
func (l *ListCombo) AppendItemString(s string) (int, error) {
	if err := l.checkCanUseItemString(); err != nil {
		return 0, err
	}
	return l.addString(s)
}

func (l *ListCombo) addString(s string) (int, error) {
	var buf []win32.WCHAR
	win32util.CString(s, &buf)
	i, err := SendMessageRet[int](l.HWND(), l.config.MsgAddString, 0, unsafe.Pointer(&buf[0]))
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, ResultError(i)
	}
	return int(i), nil
}

// GetItemString retrieves the string of the item at the specified index in a listbox or combobox.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-gettext for details.
func (l *ListCombo) GetItemString(index int) (string, error) {
	if err := l.checkCanUseItemString(); err != nil {
		return "", err
	}
	textLen, err := SendMessageRetNoError[int](l.HWND(), l.config.MsgGetItemTextLen, index, 0)
	if err != nil {
		return "", err
	}
	buf := make([]win32.WCHAR, textLen+1)
	textLen, err = SendMessageRetNoError[int](l.HWND(), l.config.MsgGetItemText, index, unsafe.Pointer(&buf[0]))
	if err != nil {
		return "", err
	}
	return win32util.GoString(&buf[0], int(textLen)+1), nil
}

func (l *ListCombo) insertItemString(index int, s string) (int, error) {
	var buf []win32.WCHAR
	win32util.CString(s, &buf)
	i, err := SendMessageRet[int](l.HWND(), l.config.MsgInsertString, index, unsafe.Pointer(&buf[0]))
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, ResultError(i)
	}
	return int(i), nil
}

// InsertItemString inserts a string at the specified index in a listbox or combobox and returns
// any error that occurred.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-insertstring for details.
func (l *ListCombo) InsertItemString(index int, s string) error {
	if err := l.checkCanUseItemString(); err != nil {
		return err
	}
	_, err := l.insertItemString(index, s)
	return err
}

// DeleteItem deletes the item at the specified index in a listbox or combobox and returns any error that occurred.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-deletestring for details.
func (l *ListCombo) DeleteItem(index int) error {
	// Delete item data
	_, data, err := l.getItemData(index)
	if err != nil {
		return err
	}
	if data != nil {
		data.Unpin()
		l.pinnedItemData.Delete(data)
	}
	// Delete the item
	return SendMessageNoError(l.HWND(), l.config.MsgDeleteString, index, 0)
}

// DeleteAllItems removes all items from a listbox or combobox and returns any error that occurred.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-resetcontent for details.
func (l *ListCombo) DeleteAllItems() error {
	// Delete all item data
	itemCount, err := l.ItemCount()
	if err != nil {
		return err
	}
	for i := range itemCount {
		_, data, err := l.getItemData(i)
		if err != nil {
			return err
		}
		if data == nil {
			continue
		}
		data.Unpin()
		l.pinnedItemData.Delete(data)
	}
	// Delete all items
	if _, err = win32.SendMessageW(l.HWND(), l.config.MsgResetContent, 0, 0); err != nil {
		return err
	}
	return nil
}

// ItemCount returns the number of items in the listbox or an error if the operation failed.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getcount for details.
func (l *ListCombo) ItemCount() (int, error) {
	i, err := SendMessageRet[int](l.HWND(), l.config.MsgGetItemCount, 0, 0)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, ResultError(i)
	}
	return i, nil
}

// getItemData returns the item data of the item at the specified index in a listbox or combobox, or an error if the operation failed.
// The returned uintptr is the raw data retrieved from Win32 API, the returned *ItemData is the item data associated with the item if
// the data is set by gw itself.
// If the item does not have item data, it returns (0, nil, nil).
// If the item data is not set by gw itself, it returns the item data as uintptr.
func (l *ListCombo) getItemData(index int) (uintptr, *ItemData, error) {
	i, err := SendMessageRetNoError[uintptr](l.HWND(), l.config.MsgGetItemData, index, 0)
	if err != nil {
		return 0, nil, err
	}
	if i == 0 {
		return 0, nil, nil
	}
	p := unsafe.Add(nil, i)
	if !l.pinnedItemData.Contains((*ItemData)(p)) {
		return uintptr(i), nil, ErrExternalItemData
	}
	return i, (*ItemData)(p), nil
}

// SetItemData sets the item data of the item at the specified index in a listbox or combobox and returns any error that occurred.
// If data is nil, it removes the item data of the item.
// If the item already has item data that was not set by gw itself, it returns [ErrExternalItemData].
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setitemdata for details.
func (l *ListCombo) SetItemData(index int, data any) error {
	_, p, err := l.getItemData(index)
	if err != nil {
		return err
	}
	if p != nil {
		if data == nil {
			// Remove item data
			p.Unpin()
			l.pinnedItemData.Delete(p)
			return SendMessageNoError(l.HWND(), l.config.MsgSetItemData, index, 0)
		} else {
			// Update the data of the existing itemData
			p.Data = &data
		}
		return nil
	}

	if data == nil {
		return nil
	}

	p = &ItemData{Data: &data}
	p.Pin()

	if l.pinnedItemData == nil {
		l.pinnedItemData = make(gg.Set[*ItemData])
	}
	l.pinnedItemData.Add(p)

	return SendMessageNoError(l.HWND(), l.config.MsgSetItemData, index, unsafe.Pointer(p))

}

// ItemData retrieves the item data of the item at the specified index in a listbox or combobox, or an error if the operation failed.
// If the item data is not set by gw itself, it returns the data as a uintptr.
// If the item does not have item data, it returns (uintptr(0), nil).
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getitemdata for details.
func (l *ListCombo) ItemData(index int) (any, error) {
	i, p, err := l.getItemData(index)
	if err == ErrExternalItemData {
		return uintptr(i), nil
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return uintptr(0), nil
	}
	return *p.Data, nil
}

// ItemDataFromRaw converts the raw item data retrieved from Win32 API to the actual item data.
// If the raw data is 0, it returns nil.
// If the raw data is not 0 but does not correspond to any item data set by gw itself, it returns the raw data as uintptr.
func (l *ListCombo) ItemDataFromRaw(raw uintptr) any {
	if raw == 0 {
		return uintptr(0)
	}
	p := unsafe.Add(nil, raw)
	if !l.pinnedItemData.Contains((*ItemData)(p)) {
		return uintptr(raw)
	}
	return *(*ItemData)(p).Data
}

// CurSelected retrieves the index of the currently selected item in a listbox or combobox,
// or an error if the operation failed.
// If no item is selected, it returns (-1, nil).
// If applied to multiple-selection list box or combobox, it returns the index of the item that
// has the focus rectangle. If no items are selected, the returned index is 0.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getcursel for details.
func (l *ListCombo) CurSelected() (int, error) {
	i, err := SendMessageRet[int](l.HWND(), l.config.MsgGetCurSel, 0, 0)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// SetCurSelected sets the selection to the item at the specified index in a listbox or combobox and returns any error that occurred.
// If index is -1, it removes the selection from all items.
// If applied to multiple-selection list box or combobox, it returns [ErrFailed].
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setcursel for details.
func (l *ListCombo) SetCurSelected(index int) error {
	i, err := SendMessageRet[win32.LRESULT](l.HWND(), l.config.MsgSetCurSel, index, 0)
	if err != nil {
		return err
	}
	// If the wParam parameter is -1, the return value is LB_ERR even though no error occurred.
	if i == LC_ERR && index != -1 {
		return ResultError(i)
	}
	return nil
}

// FindItemStringIndex finds the first item in a listbox or combobox that begins with the specified string and
// returns its index, or an error if the operation failed. If no such item found, it returns -1, nil.
//
// The prefix is case-insensitive.
//
// The search begins with the item after the startIndex and continues to the end of the list,
// and then wraps around to the beginning of the list until it reaches the startIndex.
// If startIndex is -1, the search starts from the beginning of the list.
//
// If the listbox or combobox is an owner-drawn control without the HASSTRINGS style,
// what this method does depends on whether SORT style is used.
// If SORT style is not set, it returns [ErrFindStringNotSupported], otherwise it
// calls CompareItem method to determine which item matches the specified string.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-findstring for details.
func (l *ListCombo) FindItemStringIndex(startIndex int, prefix string) (int, error) {
	if l.noItemStrings && !l.sorted {
		return 0, ErrFindStringNotSupported
	}
	var buf []win32.WCHAR
	win32util.CString(prefix, &buf)
	return SendMessageRet[int](l.HWND(), l.config.MsgFindString, win32.WPARAM(startIndex), win32.LPARAM(uintptr(unsafe.Pointer(&buf[0]))))
}

// FindItemString calls FindItemStringIndex(-1, item).
func (l *ListCombo) FindItemString(item string) (int, error) {
	return l.FindItemStringIndex(-1, item)
}

// FindItemStringExactIndex finds the first item in a listbox or combobox that exactly matches the specified string and
// returns its index, or an error if the operation failed. If no such item found, it returns -1, nil.
//
// The search begins with the item after the startIndex and continues to the end of the list,
// and then wraps around to the beginning of the list until it reaches the startIndex.
// If startIndex is -1, the search starts from the beginning of the list.
//
// If the listbox or combobox is an owner-drawn control without the HASSTRINGS style,
// what this method does depends on whether SORT style is used.
// If SORT style is not set, it returns [ErrFindStringNotSupported], otherwise it
// calls CompareItem method to determine which item matches the specified string.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-findstringexact for details.
func (l *ListCombo) FindItemStringExactIndex(startIndex int, s string) (int, error) {
	if l.noItemStrings && !l.sorted {
		return 0, ErrFindStringNotSupported
	}
	var buf []win32.WCHAR
	win32util.CString(s, &buf)

	return SendMessageRet[int](l.HWND(), l.config.MsgFindStringExact, win32.WPARAM(startIndex), win32.LPARAM(uintptr(unsafe.Pointer(&buf[0]))))
}

// FindItemStringExact calls FindItemStringExactIndex(-1, item).
func (l *ListCombo) FindItemStringExact(item string) (int, error) {
	return l.FindItemStringExactIndex(-1, item)
}

// SetSelectedItems sets the width, in pixels, by which a list box can be scrolled horizontally (the scrollable width).
// If the width of the list box is smaller than this value, the horizontal scroll bar horizontally scrolls items in
// the list box. If the width of the list box is equal to or greater than this value, the horizontal scroll bar is hidden.
// To respond to SetHorizontalExtent, the list box must have the WS_HSCROLL style.
// This method has no effect on a multiple-column list box.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-sethorizontalextent for details.
func (l *ListCombo) SetHorizontalExtent(extent int) error {
	return SendMessageNoError(l.HWND(), l.config.MsgSetHorizontalExtent, extent, 0)
}

// HorizontalExtent returns the value sets by [ListCombo.SetHorizontalExtent] or 0 if SetHorizontalExtent is not called.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-gethorizontalextent for details.
func (l *ListCombo) HorizontalExtent() (int, error) {
	return SendMessageRet[int](l.HWND(), l.config.MsgGetHorizontalExtent, 0, 0)
}

// SetItemHeight Sets the height, in pixels, of items in a list box.
// If the list box has the LBS_OWNERDRAWVARIABLE style, this message sets the height of the item specified by itemIndex parameter.
// Otherwise, this message sets the height of all items in the list box.
//
// If the list box has the LBS_OWNERDRAWVARIABLE style, the itemIndex specifies the zero-based index of the item in the list box;
// otherwise set it to 0.
//
// The height parameter specifies the height, in pixels, of the item. The maximum height is 255 pixels.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setitemheight for details.
func (l *ListCombo) SetItemHeight(itemIndex int, height int) error {
	return SendMessageNoError(l.HWND(), l.config.MsgSetItemHeight, itemIndex, height)
}

// ItemHeight sets the height of items in a list box.
// The return value is the height of the item specified by the itemIndex parameter if the list box has the LBS_OWNERDRAWVARIABLE style;
// otherwise, it is the height of each item in the list box.
//
// The itemIndex parameter is used only if the list box has the LBS_OWNERDRAWVARIABLE style; otherwise, it must be zero.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-getitemheight for details.
func (l *ListCombo) ItemHeight(itemIndex int) (int, error) {
	return SendMessageRetNoError[int](l.HWND(), l.config.MsgGetItemHeight, itemIndex, 0)
}

// SetLocale sets the current locale of the list box or combo box.
// The locale can be used to determine the correct sorting order of displayed text.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setlocale for details.
func (l *ListCombo) SetLocale(locale win32.LCID) (oldLocale win32.LCID, err error) {
	return SendMessageRetNoError[win32.LCID](l.HWND(), l.config.MsgSetLocale, win32.WPARAM(locale), 0)
}

// Locale returns the current locale of the list box.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-setlocale for details.
func (l *ListCombo) Locale() (win32.LCID, error) {
	return SendMessageRetNoError[win32.LCID](l.HWND(), l.config.MsgGetLocale, 0, 0)
}

// FirstVisible returns the index of the first visible item in a list box.
// Initially the item with index 0 is at the top of the list box, but if the list box contents have been scrolled another item may be at the top.
// The first visible item in a multiple-column list box is the top-left item.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-gettopindex for details.
func (l *ListCombo) FirstVisible() (int, error) {
	return SendMessageRet[int](l.HWND(), l.config.MsgGetTopIndex, 0, 0)
}

// EnsureVisible scrolls the list box contents so that either the specified item appears at the top of the list box or the maximum scroll range has been reached.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-settopindex for details.
func (l *ListCombo) EnsureVisible(index int) error {
	return SendMessageNoError(l.HWND(), l.config.MsgSetTopIndex, index, 0)
}

// InitStorage pre-allocates memory for storing list items.
// Parameter itemCount is the number of items to pre-allocate for, and avgStringLen is the average length of the item strings in characters.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-initstorage for details.
func (l *ListCombo) InitStorage(itemCount, avgStringLen int) error {
	i, err := SendMessageRet[win32.LRESULT](l.HWND(), l.config.MsgInitStorage, itemCount, itemCount*(avgStringLen+1)*int(unsafe.Sizeof(win32.WCHAR(0))))
	if err != nil {
		return err
	}
	if i < 0 {
		return ResultError(i)
	}
	return nil
}

// SelectStringIndex searches a list box for an item that has the specified prefix. If a matching item is found, the item is selected.
// If the search is successful, the return value is the index of the selected item. If the search is unsuccessful, the return value is [ErrFailed]
// and the current selection is not changed.
//
// See [ListCombo.FindItemStringIndex] for details on how the search is performed and the limitations of this method for owner-drawn controls.
// The prefix is case-insensitive.
//
// See https://learn.microsoft.com/en-us/windows/win32/controls/lb-selectstring for details.
func (l *ListCombo) SelectItemStringIndex(startIndex int, prefix string) error {
	if l.noItemStrings && !l.sorted {
		return ErrFindStringNotSupported
	}
	var buf []win32.WCHAR
	win32util.CString(prefix, &buf)
	_, err := SendMessageRetNoError[int](l.HWND(), l.config.MsgSelectString, win32.WPARAM(startIndex), win32.LPARAM(uintptr(unsafe.Pointer(&buf[0]))))
	return err
}

// SelectItemString calls SelectItemStringIndex(-1, prefix).
func (l *ListCombo) SelectItemString(prefix string) error {
	return l.SelectItemStringIndex(-1, prefix)
}
