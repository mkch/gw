// Package menu manipulates windows menus.
// A Menu uses certain system resource which must be released
// by Menu.Destroy(). A window menu, the one set by Window.SetMenu(),
// will be destroyed when the windows is closed.
package menu

import (
	"errors"
	"strings"
	"unicode"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/internal"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

var itemMap = objectmap.New[*Item](internal.MinMenuItemID, internal.MaxMenuItemID)

// OnWmCommand handles menu commands.
// Called by the default WndProc of window.
func OnWmCommand(id win32.WORD) bool {
	if item, ok := itemMap.Value(objectmap.Handle(id)); ok {
		return item.CallOnClick()
	}
	return false
}

var menuMap = make(map[win32.HMENU]*Menu)

type Menu struct {
	// OnAccelKeyChanged is called when the accelerator key of any item
	// in this menu and submenus changed.
	OnAccelKeyChanged func() error
	h                 win32.HMENU
	parent            *Item
	popup             bool
}

// AccelKeyTable returns all accelerator keys in this menu and its submenus.
// The order of ACCEL is unspecified.
func (m *Menu) AccelKeyTable() ([]win32.ACCEL, error) {
	count, err := m.ItemCount()
	if err != nil {
		return nil, err
	}
	var table []win32.ACCEL
	for i := 0; i < count; i++ {
		item, err := m.Item(i)
		if err != nil {
			return nil, err
		}
		k := item.accelKey
		if k != (AccelKey{}) {
			table = append(table, win32.ACCEL{
				Virt: win32.ACCEL_FVIRT(k.Mod),
				Key:  k.VKeyCode,
				Cmd:  item.id,
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

func New(popup bool) *Menu {
	r := &Menu{
		h:      gg.If(popup, gg.Must(win32.CreatePopupMenu()), gg.Must(win32.CreateMenu())),
		parent: nil,
		popup:  popup}
	menuMap[r.h] = r
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

func (m *Menu) Item(i int) (*Item, error) {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_ID,
	}
	if err := win32.GetMenuItemInfoW(m.h, win32.UINT(i), true, &mii); err != nil {
		return nil, err
	}
	if item, ok := itemMap.Value(objectmap.Handle(mii.ID)); !ok {
		panic("no this item")
	} else {
		return item, nil
	}
}

// DeleteItem deletes an Item and it's submenu if any.
// Use Item.SetSubmenu(nil) before deleting if the submenu
// is intended to be used later.
func (m *Menu) DeleteItem(item *Item) error {
	if item.Menu() != m {
		return errors.New("invalid item")
	}
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
	if err := win32.RemoveMenu(m.h, win32.UINT(item.ID()), win32.MF_BYCOMMAND); err != nil {
		return err
	}

	itemMap.Remove(objectmap.Handle(item.ID()))
	item.invalidate()
	return nil
}

func (m *Menu) DeleteItemIndex(index int) error {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_ID,
	}
	if err := win32.GetMenuItemInfoW(m.h, win32.UINT(index), true, &mii); err != nil {
		return err
	}
	item, _ := itemMap.Value(objectmap.Handle(mii.ID))
	return m.DeleteItem(item)
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

	if err := win32.DestroyMenu(m.h); err != nil {
		return err
	}
	delete(menuMap, m.h)
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

	var item = &Item{OnClick: spec.OnClick, title: spec.Title, menu: m}
	item.id = win32.WORD(itemMap.Add(item))
	for item.id == win32.IDTIMEOUT {
		itemMap.Remove(objectmap.Handle(item.id))
		item.id = win32.WORD(itemMap.Add(item))
	}

	if err = win32.InsertMenuItemW(m.h, win32.UINT(indexBefore), true, &win32.MENUITEMINFOW{
		Size:     win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask:     win32.MIIM_ID | win32.MIIM_STATE | win32.MIIM_FTYPE | win32.MIIM_STRING | win32.MIIM_SUBMENU,
		Type:     win32.UINT(gg.If(spec.Separator, win32.MFT_SEPARATOR, 0)),
		State:    win32.UINT(gg.If(spec.Checked, win32.MFS_CHECKED, 0) | gg.If(spec.Disabled, win32.MFS_DISABLED, 0)),
		ID:       win32.UINT(item.ID()),
		TypeData: &titleBuf[0],
		SubMenu:  hSubmenu,
	}); err != nil {
		return nil, err
	}
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
		modKeys = append(modKeys, "Shift")
	}
	if a&ModCtrl == ModCtrl {
		modKeys = append(modKeys, "Ctrl")
	}
	if a&ModAlt == ModAlt {
		modKeys = append(modKeys, "Alt")
	}
	return strings.Join(modKeys, "+")
}

const (
	ModAlt   AccelMod = AccelMod(win32.FALT)
	ModCtrl  AccelMod = AccelMod(win32.FCONTROL)
	ModShift AccelMod = AccelMod(win32.FSHIFT)
)

// AccelKey is an accelerator key
type AccelKey struct {
	Mod      AccelMod
	VKeyCode win32.WORD // Virtual key code. 'A', 'b' for example.
}

func (k AccelKey) String() string {
	var buf = make([]string, 0, 2)
	if mod := k.Mod.String(); mod != "" {
		buf = append(buf, mod)
	}
	buf = append(buf, string(unicode.ToUpper(rune(k.VKeyCode))))
	return strings.Join(buf, "+")
}

type Item struct {
	OnClick  func()
	menu     *Menu
	id       win32.WORD
	accelKey AccelKey
	title    string //title without accelerator key
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

func (item *Item) ID() win32.WORD {
	return item.id
}

func (item *Item) invalidate() {
	item.menu = nil
	item.id = 0
}

func (item *Item) Separator() (bool, error) {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_FTYPE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return false, err
	}
	return mii.Type&win32.MFT_SEPARATOR != 0, nil
}

// SetSeparator sets whether item is a separator.
// If SetSeparator(false) is called on a separator item, the
// item is changed to a disabled string item.
func (item *Item) SetSeparator(sep bool) error {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_FTYPE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return err
	}
	var buf []win32.WCHAR
	win32util.CString(itemDisplayTitle(item.title, item.accelKey), &buf)
	if sep {
		mii.Type |= win32.MFT_SEPARATOR
	} else {
		mii.Type &= ^win32.UINT(win32.MFT_SEPARATOR)
	}
	return win32.SetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii)
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
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_TYPE, // Retrieve Cch.
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return "", err
	}
	var buf = make([]win32.WCHAR, mii.Cch+1)
	mii.Mask = win32.MIIM_STRING
	mii.TypeData = &buf[0]
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return "", err
	}
	return win32util.GoString(&buf[0], len(buf)), nil
}

func itemDisplayTitle(title string, accelKey AccelKey) string {
	return strings.Join([]string{title, accelKey.String()}, "\t")
}

func (item *Item) SetTitle(title string) error {
	item.title = title
	displayTitle := itemDisplayTitle(item.title, item.accelKey)
	var buf []win32.WCHAR
	win32util.CString(displayTitle, &buf)
	return win32.SetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &win32.MENUITEMINFOW{
		Size:     win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask:     win32.MIIM_STRING,
		TypeData: &buf[0],
	})
}

func (item *Item) Checked() (bool, error) {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_STATE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return false, err
	}
	return mii.State&win32.MFS_CHECKED != 0, nil
}

func (item *Item) SetChecked(checked bool) error {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_STATE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return err
	}
	if checked {
		mii.State |= win32.MFS_CHECKED
	} else {
		mii.State &= ^win32.UINT(win32.MFS_CHECKED)
	}
	return win32.SetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii)
}

func (item *Item) Disabled() (bool, error) {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_STATE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return false, err
	}
	return mii.State&win32.MFS_DISABLED != 0, nil
}

// SetDisabled sets the disabled state of item.
// Has no effect on separators.
func (item *Item) SetDisabled(disabled bool) error {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_STATE,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return err
	}
	if disabled {
		mii.State |= win32.MFS_DISABLED
	} else {
		mii.State &= ^win32.UINT(win32.MFS_DISABLED)
	}
	return win32.SetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii)
}

func (item *Item) Submenu() (*Menu, error) {
	var mii = win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_SUBMENU,
	}
	if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &mii); err != nil {
		return nil, err
	}
	return menuMap[mii.SubMenu], nil
}

func (item *Item) SetSubmenu(menu *Menu) error {
	oldSubmenu, err := item.Submenu()
	if err != nil {
		return err
	}
	if oldSubmenu != nil {
		// Remove the item and insert a new one without the submenu.
		// SetMenuItemInfoW will destroy the old submenu if it is used
		// to replace the submenu.
		count, err := item.menu.ItemCount()
		if err != nil {
			return err
		}
		var index = -1
		var titleLen win32.UINT
		for i := 0; i < count; i++ {
			var mii = win32.MENUITEMINFOW{
				Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
				Mask: win32.MIIM_ID | win32.MIIM_TYPE, // MIIM_TYPE to retrieve Cch
			}
			if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(i), true, &mii); err != nil {
				return err
			}
			if mii.ID == win32.UINT(item.id) {
				index = i
				titleLen = mii.Cch
				break
			}
		}
		if index == -1 {
			panic("no such item") // A item must be in its parent.
		}
		var strBuf = make([]win32.WCHAR, titleLen+1)
		var mii = win32.MENUITEMINFOW{
			Size:     win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
			Mask:     win32.MIIM_BITMAP | win32.MIIM_CHECKMARKS | win32.MIIM_DATA | win32.MIIM_FTYPE | win32.MIIM_ID | win32.MIIM_STATE | win32.MIIM_STRING, // everything except win32.MIIM_SUBMENU
			TypeData: &strBuf[0],
			Cch:      win32.UINT(len(strBuf)),
		}
		if err := win32.GetMenuItemInfoW(item.menu.h, win32.UINT(index), true, &mii); err != nil {
			return err
		}
		if err := win32.RemoveMenu(item.menu.h, win32.UINT(index), win32.MF_BYPOSITION); err != nil {
			return err
		}
		mii.Cch = titleLen
		if err := win32.InsertMenuItemW(item.menu.h, win32.UINT(index), true, &mii); err != nil {
			return err
		}
	}
	if menu != nil {
		if err := win32.SetMenuItemInfoW(item.menu.h, win32.UINT(item.id), false, &win32.MENUITEMINFOW{
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
