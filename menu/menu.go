// Package menu manipulates windows menus.
// A Menu uses certain system resource which must be released
// by Menu.Destroy(). A window menu, the one set by Window.SetMenu(),
// will be destroyed when the windows is closed.
package menu

import (
	"errors"
	"runtime"
	"slices"
	"strings"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

type menuPinner struct {
	runtime.Pinner
	*Menu
}

// setMenuData sets m as the menu data of h.
// If m is nil, the menu data of h will be set to nil.
func setMenuData(h win32.HMENU, m *Menu) error {
	var oldData win32.ULONG_PTR
	mi := &win32.MENUINFO{
		Size: win32.DWORD(unsafe.Sizeof(win32.MENUINFO{})),
		Mask: win32.MIM_MENUDATA,
	}
	if err := win32.GetMenuInfo(h, mi); err != nil {
		return err
	}
	oldData = mi.MenuData

	var data win32.ULONG_PTR
	if m != nil {
		if oldData != 0 {
			panic("menu already attached")
		}
		p := &menuPinner{Menu: m}
		p.Pin(p)
		data = win32.ULONG_PTR(unsafe.Pointer(p))
	} else if oldData != 0 {
		(*menuPinner)(unsafe.Add(nil, oldData)).Unpin()
	} else {
		return nil // No old data and no new data, do nothing.
	}
	return win32.SetMenuInfo(h, &win32.MENUINFO{
		Size:     win32.DWORD(unsafe.Sizeof(win32.MENUINFO{})),
		Mask:     win32.MIM_MENUDATA,
		MenuData: data,
	})
}

// lookupMenu looks up the Menu associated with h and returns it.
// Returns nil if no Menu is associated with h.
func lookupMenu(h win32.HMENU) *Menu {
	mi := &win32.MENUINFO{
		Size: win32.DWORD(unsafe.Sizeof(win32.MENUINFO{})),
		Mask: win32.MIM_MENUDATA,
	}
	if err := win32.GetMenuInfo(h, mi); err != nil {
		return nil
	}
	if mi.MenuData == 0 {
		return nil
	}
	return (*menuPinner)(unsafe.Add(nil, mi.MenuData)).Menu
}

func setMenuNotifyByPos(h win32.HMENU) error {
	return win32.SetMenuInfo(h, &win32.MENUINFO{
		Size:  win32.DWORD(unsafe.Sizeof(win32.MENUINFO{})),
		Mask:  win32.MIM_STYLE,
		Style: win32.MNS_NOTIFYBYPOS, // To receive WM_MENUCOMMAND.
	})
}

// OnWmMenuCommand handles menu commands.
// Called by the default WndProc of window when handling WM_MENUCOMMAND.
func OnWmMenuCommand(wParam win32.WPARAM, lParam win32.LPARAM) bool {
	menu := lookupMenu(win32.HMENU(lParam))
	if menu == nil {
		return false
	}
	return menu.Item(int(wParam)).CallOnClick()
}

type Menu struct {
	// OnAccelKeyChanged is called when the accelerator key of any item
	// in this menu and submenus changed.
	OnAccelKeyChanged func() error
	h                 win32.HMENU
	parent            *Item
	items             []*Item
	popup             bool
}

// ItemAccel represents an Item and its accelerator key.
type ItemAccel struct {
	Item  *Item
	Accel win32.ACCEL
}

// AccelKeyTable returns all accelerator keys in this menu and its submenus.
// The Cmd fields of Accel in returned ItemAccel are all 0. The caller should set Cmd fields
// before using the returned accelerator tables.
// The order of returned [ItemAccel] is unspecified.
func (m *Menu) AccelKeyTable() ([]ItemAccel, error) {
	count, err := m.ItemCount()
	if err != nil {
		return nil, err
	}
	var table []ItemAccel
	for i := 0; i < count; i++ {
		item := m.Item(i)
		k := item.accelKey
		if k != (AccelKey{}) {
			table = append(table, ItemAccel{
				Item: item,
				Accel: win32.ACCEL{
					Virt: win32.ACCEL_FVIRT(k.Mod),
					Key:  win32.WORD(k.VKeyCode),
				},
			})
		}

		submenu, err := item.Submenu()
		if err != nil {
			return nil, err
		}
		if submenu != nil {
			subTable, err := submenu.AccelKeyTable()
			if err != nil {
				return nil, err
			}
			table = append(table, subTable...)
		}
	}
	return table, nil
}

func (m *Menu) callAccelKeyChanged() error {
	if m.OnAccelKeyChanged != nil {
		if err := m.OnAccelKeyChanged(); err != nil {
			return err
		}
	}
	if m.parent != nil && m.parent.menu != nil {
		return m.parent.menu.callAccelKeyChanged()
	}
	return nil
}

// OnKeyboardLayoutChange handles the keyboard layout change event.
// It refreshes the display title of all items in this menu and its submenus
// if needed.
func (m *Menu) OnKeyboardLayoutChange(hkl win32.HKL) error {
	for _, item := range m.items {
		if item.hkl != hkl {
			item.hkl = hkl
			if item.AccelKey() != (AccelKey{}) {
				if err := item.SetTitle(item.title); err != nil {
					return err
				}
			}
		}
		if submenu, err := item.Submenu(); err != nil {
			return err
		} else if submenu != nil {
			if err := submenu.OnKeyboardLayoutChange(hkl); err != nil {
				return err
			}
		}
	}
	return nil
}

func New(popup bool) *Menu {
	r := &Menu{
		h:      gg.If(popup, gg.Must(win32.CreatePopupMenu()), gg.Must(win32.CreateMenu())),
		parent: nil,
		popup:  popup}
	if err := setMenuData(r.h, r); err != nil {
		panic(err)
	}
	if err := setMenuNotifyByPos(r.h); err != nil {
		panic(err)
	}
	return r
}

func (m *Menu) Popup() bool {
	return m.popup
}

func (m *Menu) HMENU() win32.HMENU {
	return m.h
}

func (m *Menu) ItemCount() (int, error) {
	if count, err := win32.GetMenuItemCount(m.h); err != nil {
		return 0, err
	} else {
		return int(count), nil
	}
}

func (m *Menu) Item(i int) *Item {
	return m.items[i]
}

// itemIndex returns the index of item in m.items.
// If item is not in m.items, a panic will occur.
func (m *Menu) itemIndex(item *Item) int {
	i := slices.Index(m.items, item)
	if i == -1 {
		panic("invalid menu item")
	}
	return i
}

func (m *Menu) deleteItem(item *Item, index int) error {
	if submenu, err := item.Submenu(); err != nil {
		return err
	} else if submenu != nil {
		if err := item.SetSubmenu(nil); err != nil {
			return err
		}
		if err := submenu.Destroy(); err != nil {
			return err
		}
	}
	if err := win32.RemoveMenu(m.h, win32.UINT(index), win32.MF_BYPOSITION); err != nil {
		return err
	}
	m.items = slices.Delete(m.items, index, index+1)
	return nil
}

// DeleteItem deletes an item and destroys its submenu if it has one.
// Use Item.SetSubmenu(nil) before deleting if the submenu
// is intended to be used later.
func (m *Menu) DeleteItem(item *Item) error {
	return m.deleteItem(item, m.itemIndex(item))
}

var ErrIndexOutOfRange = errors.New("index out of range")

// DeleteItemIndex deletes a menu item from the menu by index.
func (m *Menu) DeleteItemIndex(index int) error {
	if index < 0 || index >= len(m.items) {
		return ErrIndexOutOfRange
	}
	return m.deleteItem(m.items[index], index)
}

// Destroy destroys a Menu and releases all resources it uses.
// All its submenus(if any) will be destroyed recursively.
func (m *Menu) Destroy() error {
	if m.h == 0 {
		return nil
	}
	count, err := m.ItemCount()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		if err := m.DeleteItemIndex(0); err != nil {
			return err
		}
	}

	if err := setMenuData(m.h, nil); err != nil {
		return err
	}
	if err := win32.DestroyMenu(m.h); err != nil {
		return err
	}
	m.h = 0

	return nil
}

type ItemSpec struct {
	Separator bool
	Title     string
	Checked   bool
	Disabled  bool
	Submenu   *Menu
	AccelKey  AccelKey
	OnClick   func()
}

// InsertItem inserts an item before some item.
// If indexBefore is -1, the new item will be appended to the end of m.
func (m *Menu) InsertItem(indexBefore int, spec *ItemSpec) (*Item, error) {
	var err error
	if indexBefore == -1 {
		if indexBefore, err = m.ItemCount(); err != nil {
			return nil, err
		}
	}
	var titleBuf []win32.WCHAR
	win32util.CString(itemDisplayTitle(spec.Title, spec.AccelKey), &titleBuf)
	var hSubmenu win32.HMENU
	if spec.Submenu != nil {
		hSubmenu = spec.Submenu.h
	}

	if err = win32.InsertMenuItemW(m.h, win32.UINT(indexBefore), true, &win32.MENUITEMINFOW{
		Size:     win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask:     win32.MIIM_ID | win32.MIIM_STATE | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
		Type:     win32.UINT(gg.If(spec.Separator, win32.MFT_SEPARATOR, 0)),
		State:    win32.UINT(gg.If(spec.Checked, win32.MFS_CHECKED, 0) | gg.If(spec.Disabled, win32.MFS_DISABLED, 0)),
		TypeData: &titleBuf[0],
		SubMenu:  hSubmenu,
	}); err != nil {
		return nil, err
	}

	var item = &Item{separator: spec.Separator, OnClick: spec.OnClick, title: spec.Title, menu: m, hkl: win32.GetKeyboardLayout(0)}
	m.items = slices.Insert(m.items, indexBefore, item)
	if spec.Submenu != nil {
		spec.Submenu.parent = item
	}
	item.SetAccelKey(spec.AccelKey)
	return item, nil
}

func (m *Menu) InsertSeparator(indexBefore int) (*Item, error) {
	return m.InsertItem(indexBefore, &ItemSpec{Separator: true})
}

type AccelMod win32.ACCEL_FVIRT

func (a AccelMod) String() string {
	var modKeys = make([]string, 0, 3)
	if a&ModShift == ModShift {
		modKeys = append(modKeys, keyName(win32.VK_SHIFT, "Shift"))
	}
	if a&ModCtrl == ModCtrl {
		modKeys = append(modKeys, keyName(win32.VK_CONTROL, "Ctrl"))
	}
	if a&ModAlt == ModAlt {
		modKeys = append(modKeys, keyName(win32.VK_MENU, "Alt"))
	}
	return strings.Join(modKeys, "+")
}

// keyName returns the localized name of the key specified by vkCode.
// If the name cannot be retrieved, fallback is returned.
func keyName(vkCode win32.VKCode, fallback string) string {
	name, err := win32util.KeyName(vkCode, true, 0)
	if err != nil {
		return fallback
	}
	return name
}

const (
	ModAlt   AccelMod = AccelMod(win32.FALT)
	ModCtrl  AccelMod = AccelMod(win32.FCONTROL)
	ModShift AccelMod = AccelMod(win32.FSHIFT)
)

// AccelKey is an accelerator key
type AccelKey struct {
	Mod AccelMod
	// Virtual key code.
	// For alphanumeric keys, the value is the ASCII code of digit or uppercase letter: 'A', '0', etc.
	VKeyCode win32.VKCode
}

func (k AccelKey) String() string {
	var buf = make([]string, 0, 2)
	if mod := k.Mod.String(); mod != "" {
		buf = append(buf, mod)
	}
	if k.VKeyCode != 0 {
		buf = append(buf, keyName(k.VKeyCode, string(rune(k.VKeyCode))))
	}
	return strings.Join(buf, "+")
}

type Item struct {
	separator bool // Whether this item is a separator.
	OnClick   func()
	menu      *Menu
	accelKey  AccelKey
	title     string    //title without accelerator key
	hkl       win32.HKL // The input locale identifier used to display the accelerator key of this item.
}

func (item *Item) SetAccelKey(accel AccelKey) error {
	if accel != item.accelKey {
		item.accelKey = accel
		// http://stackoverflow.com/questions/23592079/why-does-createacceleratortable-not-work-without-fvirtkey
		// https://msdn.microsoft.com/en-us/library/windows/desktop/dd375731(v=vs.85).aspx
		item.accelKey.Mod |= AccelMod(win32.FVIRTKEY)
		if item.menu != nil {
			if err := item.menu.callAccelKeyChanged(); err != nil {
				return err
			}
		}
		// Update display title.
		return item.SetTitle(item.title)
	}
	return nil
}

func (item *Item) AccelKey() AccelKey {
	return item.accelKey
}

func (item *Item) Menu() *Menu {
	return item.menu
}

func (item *Item) Separator() bool {
	return item.separator
}

// SetSeparator sets whether item is a separator.
// If SetSeparator(false) is called on a separator item, the
// item is changed to a disabled string item.
func (item *Item) SetSeparator(sep bool) error {
	i := item.menu.itemIndex(item)
	mii := win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
	}
	if sep {
		mii.Mask = win32.MIIM_FTYPE
		mii.Type = win32.MFT_SEPARATOR
	} else {
		var buf []win32.WCHAR
		win32util.CString(itemDisplayTitle(item.title, item.accelKey), &buf)
		mii.Mask = win32.MIIM_FTYPE | win32.MIIM_STRING
		mii.TypeData = &buf[0]
	}
	if err := win32.SetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
		return err
	}
	item.separator = sep
	return nil
}

func (item *Item) CallOnClick() bool {
	if item.OnClick == nil {
		return false
	}
	item.OnClick()
	return true
}

func (item *Item) Title() string {
	return item.title
}

// Title with accelerator key.
func (item *Item) DisplayTitle() (string, error) {
	i := item.menu.itemIndex(item)
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_TYPE, // Retrieve Cch.
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
		return "", err
	}
	var buf = make([]win32.WCHAR, mii.Cch+1)
	mii.Mask = win32.MIIM_STRING
	mii.TypeData = &buf[0]
	mii.Cch++ // Include null terminator.
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
		return "", err
	}
	return win32util.GoString(&buf[0], len(buf)), nil
}

func itemDisplayTitle(title string, accelKey AccelKey) string {
	return strings.Join([]string{title, accelKey.String()}, "\t")
}

func (item *Item) SetTitle(title string) error {
	i := item.menu.itemIndex(item)
	item.title = title
	displayTitle := itemDisplayTitle(item.title, item.accelKey)
	var buf []win32.WCHAR
	win32util.CString(displayTitle, &buf)
	return win32.SetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &win32.MENUITEMINFOW{
		Size:     win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask:     win32.MIIM_STRING,
		TypeData: &buf[0],
	})
}

func (item *Item) Checked() (bool, error) {
	i := item.menu.itemIndex(item)
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_STATE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
		return false, err
	}
	return mii.State&win32.MFS_CHECKED != 0, nil
}

func (item *Item) SetChecked(checked bool) error {
	i := item.menu.itemIndex(item)
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_STATE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
		return err
	}
	if checked {
		mii.State |= win32.MFS_CHECKED
	} else {
		mii.State &= ^win32.UINT(win32.MFS_CHECKED)
	}
	return win32.SetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii)
}

func (item *Item) Disabled() (bool, error) {
	i := item.menu.itemIndex(item)
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_STATE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
		return false, err
	}
	return mii.State&win32.MFS_DISABLED != 0, nil
}

// SetDisabled sets the disabled state of item.
// Has no effect on separators.
func (item *Item) SetDisabled(disabled bool) error {
	i := item.menu.itemIndex(item)
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_STATE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
		return err
	}
	if disabled {
		mii.State |= win32.MFS_DISABLED
	} else {
		mii.State &= ^win32.UINT(win32.MFS_DISABLED)
	}
	return win32.SetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii)
}

func (item *Item) Submenu() (*Menu, error) {
	i := item.menu.itemIndex(item)
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_SUBMENU,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
		return nil, err
	}
	return lookupMenu(mii.SubMenu), nil
}

func (item *Item) SetSubmenu(menu *Menu) error {
	i := item.menu.itemIndex(item)
	oldSubmenu, err := item.Submenu()
	if err != nil {
		return err
	}
	if oldSubmenu != nil {
		// Remove the item and insert a new one without the submenu.
		// SetMenuItemInfoW will destroy the old submenu if it is used
		// to replace the submenu.

		// Get title string length
		var mii = win32.MENUITEMINFOW{
			Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
			Mask: win32.MIIM_TYPE, // MIIM_TYPE to retrieve Cch
		}
		if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
			return err
		}

		// Get old menu item info without submenu
		var strBuf = make([]win32.WCHAR, mii.Cch+1)
		mii = win32.MENUITEMINFOW{
			Size:     win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
			Mask:     win32.MIIM_BITMAP | win32.MIIM_CHECKMARKS | win32.MIIM_DATA | win32.MIIM_FTYPE | win32.MIIM_STATE | win32.MIIM_STRING, // everything except win32.MIIM_SUBMENU and win32.MIIM_ID
			TypeData: &strBuf[0],
			Cch:      win32.UINT(len(strBuf)),
		}
		if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
			return err
		}

		// Remove the item and insert a new one without the submenu.
		if err := win32.RemoveMenu(item.menu.h, win32.UINT(i), win32.MF_BYPOSITION); err != nil {
			return err
		}
		if err := win32.InsertMenuItemW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
			return err
		}
	}
	if menu != nil {
		if err := win32.SetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &win32.MENUITEMINFOW{
			Size:    win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
			Mask:    win32.MIIM_SUBMENU,
			SubMenu: menu.h,
		}); err != nil {
			return err
		}
	}

	if oldSubmenu != nil {
		oldSubmenu.parent = nil
	}
	if menu != nil {
		menu.parent = item
	}
	return nil
}
