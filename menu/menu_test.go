package menu_test

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

// getItemMII retrieves MENUITEMINFOW for the item at index in menu h.
// The pattern mirrors the read calls used throughout menu.go.
func getItemMII(t *testing.T, h win32.HMENU, index int, mask win32.UINT) win32.MENUITEMINFOW {
	t.Helper()
	mii := win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: mask,
	}
	if err := win32.GetMenuItemInfoW(h, win32.UINT(index), true, &mii); err != nil {
		t.Fatal(err)
	}
	return mii
}

// getItemTitle reads the display title (including accelerator tab) of the item at
// index directly from the win32 menu h. The two-pass approach mirrors Item.DisplayTitle().
func getItemTitle(t *testing.T, h win32.HMENU, index int) string {
	t.Helper()
	// First pass: get string length via MIIM_TYPE (populates Cch).
	mii := win32.MENUITEMINFOW{
		Size: win32.UINT(unsafe.Sizeof(win32.MENUITEMINFOW{})),
		Mask: win32.MIIM_TYPE,
	}
	if err := win32.GetMenuItemInfoW(h, win32.UINT(index), true, &mii); err != nil {
		t.Fatal(err)
	}
	buf := make([]win32.WCHAR, mii.Cch+1)
	// Second pass: read the actual string.
	// mii.Cch must be updated to the buffer size; Windows copies at most Cch-1 chars.
	mii.Mask = win32.MIIM_STRING
	mii.TypeData = &buf[0]
	mii.Cch = win32.UINT(len(buf))
	if err := win32.GetMenuItemInfoW(h, win32.UINT(index), true, &mii); err != nil {
		t.Fatal(err)
	}
	return win32util.GoString(&buf[0], len(buf))
}

// ── Phase 1: Creation and Destruction ────────────────────────────────────────

func TestNew_Popup(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	if m.HMENU() == 0 {
		t.Fatal("HMENU() should not be 0")
	}
	if !m.Popup() {
		t.Fatal("Popup() should be true")
	}
	count, err := m.ItemCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("ItemCount() = %d, want 0", count)
	}
}

func TestNew_Regular(t *testing.T) {
	m := menu.New(false)
	defer m.Destroy()
	if m.HMENU() == 0 {
		t.Fatal("HMENU() should not be 0")
	}
	if m.Popup() {
		t.Fatal("Popup() should be false")
	}
	count, err := m.ItemCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("ItemCount() = %d, want 0", count)
	}
}

func TestDestroy(t *testing.T) {
	m := menu.New(true)
	if _, err := m.InsertItem(-1, &menu.ItemSpec{Title: "A"}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.InsertItem(-1, &menu.ItemSpec{Title: "B"}); err != nil {
		t.Fatal(err)
	}
	if err := m.Destroy(); err != nil {
		t.Fatalf("Destroy() error: %v", err)
	}
	if m.HMENU() != 0 {
		t.Fatal("HMENU() should be 0 after Destroy()")
	}
}

// ── Phase 2: InsertItem and ItemCount ────────────────────────────────────────

func TestInsertItem_Append(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	for _, title := range []string{"A", "B", "C"} {
		if _, err := m.InsertItem(-1, &menu.ItemSpec{Title: title}); err != nil {
			t.Fatal(err)
		}
	}
	count, err := m.ItemCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("ItemCount() = %d, want 3", count)
	}
	// Cross-verify with win32.
	win32Count, err := win32.GetMenuItemCount(m.HMENU())
	if err != nil {
		t.Fatal(err)
	}
	if int(win32Count) != 3 {
		t.Fatalf("win32 GetMenuItemCount = %d, want 3", win32Count)
	}
}

func TestInsertItem_AtIndex(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	if _, err := m.InsertItem(-1, &menu.ItemSpec{Title: "B"}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.InsertItem(0, &menu.ItemSpec{Title: "A"}); err != nil {
		t.Fatal(err)
	}
	titleAt0 := getItemTitle(t, m.HMENU(), 0)
	if !strings.HasPrefix(titleAt0, "A") {
		t.Fatalf("index 0 title = %q, want prefix %q", titleAt0, "A")
	}
	titleAt1 := getItemTitle(t, m.HMENU(), 1)
	if !strings.HasPrefix(titleAt1, "B") {
		t.Fatalf("index 1 title = %q, want prefix %q", titleAt1, "B")
	}
}

func TestInsertSeparator(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	if _, err := m.InsertSeparator(-1); err != nil {
		t.Fatal(err)
	}
	count, err := m.ItemCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("ItemCount() = %d, want 1 after InsertSeparator", count)
	}
	mii := getItemMII(t, m.HMENU(), 0, win32.MIIM_FTYPE)
	if mii.Type&win32.MFT_SEPARATOR == 0 {
		t.Fatal("MFT_SEPARATOR should be set for separator item")
	}
}

// ── Phase 3: Item properties ──────────────────────────────────────────────────

func TestItem_Title(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "Hello"})
	if err != nil {
		t.Fatal(err)
	}
	if item.Title() != "Hello" {
		t.Fatalf("Title() = %q, want %q", item.Title(), "Hello")
	}
	display, err := item.DisplayTitle()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(display, "Hello") {
		t.Fatalf("DisplayTitle() = %q, want prefix %q", display, "Hello")
	}
	win32Title := getItemTitle(t, m.HMENU(), 0)
	if !strings.HasPrefix(win32Title, "Hello") {
		t.Fatalf("win32 title = %q, want prefix %q", win32Title, "Hello")
	}
}

func TestItem_SetTitle(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "Old"})
	if err != nil {
		t.Fatal(err)
	}
	if err := item.SetTitle("New"); err != nil {
		t.Fatal(err)
	}
	if item.Title() != "New" {
		t.Fatalf("Title() = %q, want %q", item.Title(), "New")
	}
	win32Title := getItemTitle(t, m.HMENU(), 0)
	if !strings.HasPrefix(win32Title, "New") {
		t.Fatalf("win32 title after SetTitle = %q, want prefix %q", win32Title, "New")
	}
}

func TestItem_Checked_Init(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "X", Checked: true})
	if err != nil {
		t.Fatal(err)
	}
	checked, err := item.Checked()
	if err != nil {
		t.Fatal(err)
	}
	if !checked {
		t.Fatal("Checked() should be true when inserted with Checked:true")
	}
	mii := getItemMII(t, m.HMENU(), 0, win32.MIIM_STATE)
	if mii.State&win32.MFS_CHECKED == 0 {
		t.Fatal("win32 MFS_CHECKED should be set")
	}
}

func TestItem_SetChecked(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "X"})
	if err != nil {
		t.Fatal(err)
	}

	if err := item.SetChecked(true); err != nil {
		t.Fatal(err)
	}
	checked, err := item.Checked()
	if err != nil {
		t.Fatal(err)
	}
	if !checked {
		t.Fatal("Checked() should be true after SetChecked(true)")
	}
	mii := getItemMII(t, m.HMENU(), 0, win32.MIIM_STATE)
	if mii.State&win32.MFS_CHECKED == 0 {
		t.Fatal("win32 MFS_CHECKED should be set after SetChecked(true)")
	}

	if err := item.SetChecked(false); err != nil {
		t.Fatal(err)
	}
	checked, err = item.Checked()
	if err != nil {
		t.Fatal(err)
	}
	if checked {
		t.Fatal("Checked() should be false after SetChecked(false)")
	}
	mii = getItemMII(t, m.HMENU(), 0, win32.MIIM_STATE)
	if mii.State&win32.MFS_CHECKED != 0 {
		t.Fatal("win32 MFS_CHECKED should be cleared after SetChecked(false)")
	}
}

func TestItem_Disabled_Init(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "X", Disabled: true})
	if err != nil {
		t.Fatal(err)
	}
	disabled, err := item.Disabled()
	if err != nil {
		t.Fatal(err)
	}
	if !disabled {
		t.Fatal("Disabled() should be true when inserted with Disabled:true")
	}
	mii := getItemMII(t, m.HMENU(), 0, win32.MIIM_STATE)
	if mii.State&win32.MFS_DISABLED == 0 {
		t.Fatal("win32 MFS_DISABLED should be set")
	}
}

func TestItem_SetDisabled(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "X"})
	if err != nil {
		t.Fatal(err)
	}

	if err := item.SetDisabled(true); err != nil {
		t.Fatal(err)
	}
	disabled, err := item.Disabled()
	if err != nil {
		t.Fatal(err)
	}
	if !disabled {
		t.Fatal("Disabled() should be true after SetDisabled(true)")
	}
	mii := getItemMII(t, m.HMENU(), 0, win32.MIIM_STATE)
	if mii.State&win32.MFS_DISABLED == 0 {
		t.Fatal("win32 MFS_DISABLED should be set after SetDisabled(true)")
	}

	if err := item.SetDisabled(false); err != nil {
		t.Fatal(err)
	}
	disabled, err = item.Disabled()
	if err != nil {
		t.Fatal(err)
	}
	if disabled {
		t.Fatal("Disabled() should be false after SetDisabled(false)")
	}
	mii = getItemMII(t, m.HMENU(), 0, win32.MIIM_STATE)
	if mii.State&win32.MFS_DISABLED != 0 {
		t.Fatal("win32 MFS_DISABLED should be cleared after SetDisabled(false)")
	}
}

func TestItem_SetSeparator(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "Hello"})
	if err != nil {
		t.Fatal(err)
	}

	if err := item.SetSeparator(true); err != nil {
		t.Fatal(err)
	}
	if !item.Separator() {
		t.Fatal("Separator() should be true after SetSeparator(true)")
	}
	if disabled, err := item.Disabled(); err != nil {
		t.Fatal(err)
	} else if !disabled {
		t.Fatal("Disabled() should be true for string item converted from separator")
	}
	mii := getItemMII(t, m.HMENU(), 0, win32.MIIM_FTYPE)
	if mii.Type&win32.MFT_SEPARATOR == 0 {
		t.Fatal("win32 MFT_SEPARATOR should be set after SetSeparator(true)")
	}

	if err := item.SetSeparator(false); err != nil {
		t.Fatal(err)
	}
	if item.Separator() {
		t.Fatal("Separator() should be false after SetSeparator(false)")
	}
	mii = getItemMII(t, m.HMENU(), 0, win32.MIIM_FTYPE)
	if mii.Type&win32.MFT_SEPARATOR != 0 {
		t.Fatal("win32 MFT_SEPARATOR should be cleared after SetSeparator(false)")
	}
	win32Title := getItemTitle(t, m.HMENU(), 0)
	if !strings.HasPrefix(win32Title, "Hello") {
		t.Fatalf("title after SetSeparator(false) = %q, want prefix %q", win32Title, "Hello")
	}
}

func TestItem_SetSeparator_RestoresTitle(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "Hello"})
	if err != nil {
		t.Fatal(err)
	}
	if err := item.SetSeparator(true); err != nil {
		t.Fatal(err)
	}
	if err := item.SetSeparator(false); err != nil {
		t.Fatal(err)
	}
	win32Title := getItemTitle(t, m.HMENU(), 0)
	if !strings.HasPrefix(win32Title, "Hello") {
		t.Fatalf("title after SetSeparator(true/false) = %q, want prefix %q", win32Title, "Hello")
	}
}

// ── Phase 4: Accelerator keys ─────────────────────────────────────────────────

func TestItem_AccelKey(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{
		Title:    "Save",
		AccelKey: menu.AccelKey{Mod: menu.ModCtrl, VKeyCode: 'S'},
	})
	if err != nil {
		t.Fatal(err)
	}

	k := item.AccelKey()
	if k.VKeyCode != 'S' {
		t.Fatalf("AccelKey().VKeyCode = %v, want 'S'", k.VKeyCode)
	}
	if k.Mod&menu.ModCtrl == 0 {
		t.Fatal("AccelKey().Mod should include ModCtrl")
	}

	win32Title := getItemTitle(t, m.HMENU(), 0)
	if !strings.Contains(win32Title, "\tCtrl+S") {
		t.Fatalf("win32 title = %q, want it to contain %q", win32Title, "\tCtrl+S")
	}

	table, err := m.AccelKeyTable()
	if err != nil {
		t.Fatal(err)
	}
	if len(table) != 1 {
		t.Fatalf("AccelKeyTable() len = %d, want 1", len(table))
	}
	if table[0].Item != item {
		t.Fatal("AccelKeyTable()[0].Item mismatch")
	}
}

func TestItem_DisplayTitle(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{
		Title:    "Save",
		AccelKey: menu.AccelKey{Mod: menu.ModCtrl, VKeyCode: 'S'},
	})
	if err != nil {
		t.Fatal(err)
	}
	display, err := item.DisplayTitle()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(display, "\tCtrl+S") {
		t.Fatalf("DisplayTitle() = %q, want it to contain %q", display, "\tCtrl+S")
	}
}

// ── Phase 5: Submenus ─────────────────────────────────────────────────────────

func TestItem_Submenu_Init(t *testing.T) {
	parent := menu.New(false)
	defer parent.Destroy() // also destroys child submenu
	child := menu.New(true)

	item, err := parent.InsertItem(-1, &menu.ItemSpec{Title: "Sub", Submenu: child})
	if err != nil {
		t.Fatal(err)
	}
	childHMenu := child.HMENU()

	// Verify via win32.
	mii := getItemMII(t, parent.HMENU(), 0, win32.MIIM_SUBMENU)
	if mii.SubMenu != childHMenu {
		t.Fatalf("win32 SubMenu = %v, want %v", mii.SubMenu, childHMenu)
	}

	// Verify via Go API.
	sub, err := item.Submenu()
	if err != nil {
		t.Fatal(err)
	}
	if sub != child {
		t.Fatal("Submenu() should return child")
	}
}

func TestItem_SetSubmenu(t *testing.T) {
	parent := menu.New(false)
	defer parent.Destroy()

	item, err := parent.InsertItem(-1, &menu.ItemSpec{Title: "Sub"})
	if err != nil {
		t.Fatal(err)
	}

	// No submenu initially.
	sub, err := item.Submenu()
	if err != nil {
		t.Fatal(err)
	}
	if sub != nil {
		t.Fatal("Submenu() should be nil initially")
	}

	// Attach child.
	child := menu.New(true)
	childHMenu := child.HMENU()
	if err := item.SetSubmenu(child); err != nil {
		t.Fatal(err)
	}
	mii := getItemMII(t, parent.HMENU(), 0, win32.MIIM_SUBMENU)
	if mii.SubMenu != childHMenu {
		t.Fatalf("win32 SubMenu = %v, want %v", mii.SubMenu, childHMenu)
	}
	sub, err = item.Submenu()
	if err != nil {
		t.Fatal(err)
	}
	if sub != child {
		t.Fatal("Submenu() should return child after SetSubmenu")
	}

	// Detach child; child is no longer owned by parent.
	if err := item.SetSubmenu(nil); err != nil {
		t.Fatal(err)
	}
	mii = getItemMII(t, parent.HMENU(), 0, win32.MIIM_SUBMENU)
	if mii.SubMenu != 0 {
		t.Fatalf("win32 SubMenu = %v, want 0 after SetSubmenu(nil)", mii.SubMenu)
	}
	// Manually destroy the detached child.
	if err := child.Destroy(); err != nil {
		t.Fatal(err)
	}
}

// ── Phase 6: Deletion ─────────────────────────────────────────────────────────

func TestDeleteItemIndex(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	for _, title := range []string{"A", "B", "C"} {
		if _, err := m.InsertItem(-1, &menu.ItemSpec{Title: title}); err != nil {
			t.Fatal(err)
		}
	}

	if err := m.DeleteItemIndex(1); err != nil {
		t.Fatal(err)
	}

	count, err := m.ItemCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("ItemCount() = %d, want 2 after deleting index 1", count)
	}
	win32Count, err := win32.GetMenuItemCount(m.HMENU())
	if err != nil {
		t.Fatal(err)
	}
	if int(win32Count) != 2 {
		t.Fatalf("win32 GetMenuItemCount = %d, want 2", win32Count)
	}
	titleAt0 := getItemTitle(t, m.HMENU(), 0)
	if !strings.HasPrefix(titleAt0, "A") {
		t.Fatalf("index 0 title = %q, want prefix %q", titleAt0, "A")
	}
	titleAt1 := getItemTitle(t, m.HMENU(), 1)
	if !strings.HasPrefix(titleAt1, "C") {
		t.Fatalf("index 1 title = %q, want prefix %q", titleAt1, "C")
	}
}

func TestDeleteItem(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	item, err := m.InsertItem(-1, &menu.ItemSpec{Title: "X"})
	if err != nil {
		t.Fatal(err)
	}
	if err := m.DeleteItem(item); err != nil {
		t.Fatal(err)
	}
	count, err := m.ItemCount()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("ItemCount() = %d, want 0 after DeleteItem", count)
	}
}

// ── Phase 7: Message handling ─────────────────────────────────────────────────

func TestOnWmMenuCommand(t *testing.T) {
	m := menu.New(true)
	defer m.Destroy()
	called := false
	if _, err := m.InsertItem(-1, &menu.ItemSpec{
		Title:   "Click",
		OnClick: func() { called = true },
	}); err != nil {
		t.Fatal(err)
	}
	result := menu.OnWmMenuCommand(0, win32.LPARAM(m.HMENU()))
	if !result {
		t.Fatal("OnWmMenuCommand should return true when item has OnClick")
	}
	if !called {
		t.Fatal("OnClick callback should have been invoked")
	}
}
