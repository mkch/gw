package combobox_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/combobox"
	"github.com/mkch/gw/listbox/listcombo"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

func TestAppendItemString(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != 0 {
			t.Fatalf("expected item count 0, got %d", count)
		}

		if i, err := combo.AppendItemString("Item 1"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		} else if i != 0 {
			t.Fatalf("expected index 0, got %d", i)
		}

		if i, err := combo.AppendItemString("Item 2"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		} else if i != 1 {
			t.Fatalf("expected index 1, got %d", i)
		}

		if i, err := combo.AppendItemString("Item 3"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		} else if i != 2 {
			t.Fatalf("expected index 2, got %d", i)
		}

		if str, err := combo.GetItemString(0); err != nil {
			t.Fatalf("failed to get item string: %v", err)
		} else if str != "Item 1" {
			t.Fatalf("expected 'Item 1', got '%s'", str)
		}

		if str, err := combo.GetItemString(1); err != nil {
			t.Fatalf("failed to get item string: %v", err)
		} else if str != "Item 2" {
			t.Fatalf("expected 'Item 2', got '%s'", str)
		}

		if str, err := combo.GetItemString(2); err != nil {
			t.Fatalf("failed to get item string: %v", err)
		} else if str != "Item 3" {
			t.Fatalf("expected 'Item 3', got '%s'", str)
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != 3 {
			t.Fatalf("expected item count 3, got %d", count)
		}

		win.Close()

	}, nil)
}

func TestInsertItemString(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		if _, err := combo.AppendItemString("Item 1"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}
		if _, err := combo.AppendItemString("Item 3"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != 2 {
			t.Fatalf("expected item count 2, got %d", count)
		}

		if err := combo.InsertItemString(1, "Item 2"); err != nil {
			t.Fatalf("failed to insert item: %v", err)
		}

		expectedItems := []string{"Item 1", "Item 2", "Item 3"}
		for i, expected := range expectedItems {
			if str, err := combo.GetItemString(i); err != nil {
				t.Fatalf("failed to get item string at index %d: %v", i, err)
			} else if str != expected {
				t.Fatalf("expected '%s' at index %d, got '%s'", expected, i, str)
			}
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != len(expectedItems) {
			t.Fatalf("expected item count %d, got %d", len(expectedItems), count)
		}

		win.Close()

	}, nil)
}

func TestInsertItemStrings(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		if _, err := combo.AppendItemString("Item 1"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}
		if _, err := combo.AppendItemString("Item 4"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != 2 {
			t.Fatalf("expected item count 2, got %d", count)
		}

		insertItems := []string{"Item 2", "Item 3"}
		for i, item := range insertItems {
			if err := combo.InsertItemString(1+i, item); err != nil {
				t.Fatalf("failed to insert item: %v", err)
			}
		}

		expectedItems := []string{"Item 1", "Item 2", "Item 3", "Item 4"}
		for i, expected := range expectedItems {
			if str, err := combo.GetItemString(i); err != nil {
				t.Fatalf("failed to get item string at index %d: %v", i, err)
			} else if str != expected {
				t.Fatalf("expected '%s' at index %d, got '%s'", expected, i, str)
			}
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != len(expectedItems) {
			t.Fatalf("expected item count %d, got %d", len(expectedItems), count)
		}

		win.Close()

	}, nil)
}

func TestDeleteItems(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		items := []string{"Item 1", "Item 2", "Item 3", "Item 4"}
		for _, item := range items {
			if _, err := combo.AppendItemString(item); err != nil {
				t.Fatalf("failed to append item: %v", err)
			}
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != len(items) {
			t.Fatalf("expected item count %d, got %d", len(items), count)
		}

		if err := combo.DeleteItem(1); err != nil {
			t.Fatalf("failed to delete item: %v", err)
		}

		expectedAfterDelete := []string{"Item 1", "Item 3", "Item 4"}
		for i, expected := range expectedAfterDelete {
			if str, err := combo.GetItemString(i); err != nil {
				t.Fatalf("failed to get item string at index %d: %v", i, err)
			} else if str != expected {
				t.Fatalf("expected '%s' at index %d, got '%s'", expected, i, str)
			}
		}

		if _, err := combo.GetItemString(3); err == nil {
			t.Fatalf("expected error when getting removed item at index 3")
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != len(expectedAfterDelete) {
			t.Fatalf("expected item count %d, got %d", len(expectedAfterDelete), count)
		}

		if err := combo.DeleteAllItems(); err != nil {
			t.Fatalf("failed to delete all items: %v", err)
		}

		if _, err := combo.GetItemString(0); err == nil {
			t.Fatalf("expected error when getting item from empty combo")
		}

		if count, err := combo.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != 0 {
			t.Fatalf("expected item count 0, got %d", count)
		}

		win.Close()

	}, nil)
}

func TestItemData(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		items := []string{"Item 1", "Item 2", "Item 3"}
		for _, item := range items {
			if _, err := combo.AppendItemString(item); err != nil {
				t.Fatalf("failed to append item: %v", err)
			}
		}

		if data, err := combo.ItemData(0); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != uintptr(0) {
			t.Fatalf("expected item data uintptr(0) at index 0 before any SetItemData, got %v", data)
		}

		if err := combo.SetItemData(0, "abc"); err != nil {
			t.Fatalf("failed to set item data: %v", err)
		}
		if err := combo.SetItemData(1, 200); err != nil {
			t.Fatalf("failed to set item data: %v", err)
		}

		if data, err := combo.ItemData(0); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != "abc" {
			t.Fatalf(`expected item data "abc" at index 0, got %v`, data)
		}

		if data, err := combo.ItemData(1); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != 200 {
			t.Fatalf(`expected item data 200 at index 1, got %v`, data)
		}

		if data, err := combo.ItemData(2); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != uintptr(0) {
			t.Fatalf("expected item data uintptr(0) at index 2 (never set), got %v", data)
		}

		if err := combo.SetItemData(0, nil); err != nil {
			t.Fatalf("failed to set item data to nil: %v", err)
		}

		if data, err := combo.ItemData(0); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != uintptr(0) {
			t.Fatalf("expected item data uintptr(0) at index 0 after SetItemData(0, nil), got %v", data)
		}

		// Verify that SetItemData does not affect item strings.
		for i, expected := range items {
			if str, err := combo.GetItemString(i); err != nil {
				t.Fatalf("failed to get item string at index %d: %v", i, err)
			} else if str != expected {
				t.Fatalf("expected '%s' at index %d, got '%s'", expected, i, str)
			}
		}

		win.Close()

	}, nil)
}

func TestExternalItemData(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		if _, err := combo.AppendItemString("Item 1"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}

		// Set external item data via raw Win32 message.
		if _, err := win32.SendMessageW(combo.HWND(), combobox.CB_SETITEMDATA,
			win32.WPARAM(0), win32.LPARAM(uintptr(1))); err != nil {
			t.Fatalf("failed to set external item data: %v", err)
		}

		if data, err := combo.ItemData(0); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != uintptr(1) {
			t.Fatalf("expected item data uintptr(1) at index 0, got %v", data)
		}

		if err := combo.SetItemData(0, "test"); !errors.Is(err, listcombo.ErrExternalItemData) {
			t.Fatalf("expected ErrExternalItemData from SetItemData, got %v", err)
		}

		// After external data is still there, verify SetItemData(nil) also fails.
		if err := combo.SetItemData(0, nil); !errors.Is(err, listcombo.ErrExternalItemData) {
			t.Fatalf("expected ErrExternalItemData from SetItemData(0, nil), got %v", err)
		}

		win.Close()

	}, nil)
}

func TestFindItemString(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		unsorted, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		unsorted.AppendItemString("C2")
		unsorted.AppendItemString("A1")
		unsorted.AppendItemString("b")

		i, err := unsorted.FindItemString("A")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != 1 {
			t.Fatalf("expected index 1 for item 'A', got %d", i)
		}

		i, err = unsorted.FindItemString("b")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != 2 {
			t.Fatalf("expected index 2 for item 'b', got %d", i)
		}

		i, err = unsorted.FindItemString("C")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != 0 {
			t.Fatalf("expected index 0 for item 'C', got %d", i)
		}

		sorted, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE | combobox.CBS_SORT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		sorted.AppendItemString("C")
		sorted.AppendItemString("A")
		sorted.AppendItemString("b")

		i, err = sorted.FindItemString("A")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != 0 {
			t.Fatalf("expected index 0 for item 'A', got %d", i)
		}

		i, err = sorted.FindItemString("b")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != 1 {
			t.Fatalf("expected index 1 for item 'b', got %d", i)
		}

		i, err = sorted.FindItemString("C")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != 2 {
			t.Fatalf("expected index 2 for item 'C', got %d", i)
		}

		win.Close()

	}, nil)
}

func TestFindItemStringExact(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		unsorted, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		unsorted.AppendItemString("C2")
		unsorted.AppendItemString("A1")
		unsorted.AppendItemString("b")

		i, err := unsorted.FindItemStringExact("A")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != -1 {
			t.Fatalf("expected index -1 for item 'A' (no exact match), got %d", i)
		}

		i, err = unsorted.FindItemString("b")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != 2 {
			t.Fatalf("expected index 2 for item 'b', got %d", i)
		}

		i, err = unsorted.FindItemString("C2")
		if err != nil {
			t.Fatalf("failed to find item string: %v", err)
		}
		if i != 0 {
			t.Fatalf("expected index 0 for item 'C2', got %d", i)
		}

		win.Close()

	}, nil)
}

func TestOnCompareItem(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		var combo *combobox.ComboBox
		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE | combobox.CBS_SORT | combobox.CBS_OWNERDRAWFIXED,
			OnCompareItem: func(item1, item2 any, locale win32.DWORD) int {
				str1, isstr1 := item1.(string)
				str2, isstr2 := item2.(string)
				if !isstr1 {
					return -1 // Non-string is less than string
				}
				if !isstr2 {
					return 1 // String is greater than non-string
				}
				return strings.Compare(str1, str2)
			},
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		combo.AppendItem("C")
		combo.AppendItem("A")
		combo.AppendItem("B")

		win32.SendMessageW(combo.HWND(), combobox.CB_ADDSTRING, 0, 100)

		data, err := combo.ItemData(0)
		if err != nil {
			t.Fatalf("failed to get item data: %v", err)
		}
		if data != uintptr(100) {
			t.Fatalf(`expected item data uintptr(100) at index 0 due to OnCompareItem, got %#v`, data)
		}

		data, err = combo.ItemData(1)
		if err != nil {
			t.Fatalf("failed to get item data: %v", err)
		}
		if data != "A" {
			t.Fatalf(`expected item data "A" at index 1 due to OnCompareItem, got %#v`, data)
		}

		data, err = combo.ItemData(2)
		if err != nil {
			t.Fatalf("failed to get item data: %v", err)
		}
		if data != "B" {
			t.Fatalf(`expected item data "B" at index 2 due to OnCompareItem, got %#v`, data)
		}

		data, err = combo.ItemData(3)
		if err != nil {
			t.Fatalf("failed to get item data: %v", err)
		}
		if data != "C" {
			t.Fatalf(`expected item data "C" at index 2 due to OnCompareItem, got %#v`, data)
		}

		win.Close()
	}, nil)
}

func TestSingleSel(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		combo.AppendItemString("Item 1")
		combo.AppendItemString("Item 2")
		combo.AppendItemString("Item 3")

		i, err := combo.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if i != -1 {
			t.Fatalf("expected no selection (-1), got %d", i)
		}

		err = combo.SetCurSelected(1)
		if err != nil {
			t.Fatalf("failed to set current selection: %v", err)
		}
		i, err = combo.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if i != 1 {
			t.Fatalf("expected selection index 1, got %d", i)
		}

		err = combo.SetCurSelected(-1)
		if err != nil {
			t.Fatalf("failed to set current selection: %v", err)
		}
		i, err = combo.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if i != -1 {
			t.Fatalf("expected no selection (-1), got %d", i)
		}

		win.Close()

	}, nil)
}

func TestHorizontalExtent(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = combo.SetHorizontalExtent(200)
		if err != nil {
			t.Fatalf("failed to set horizontal extent, got: %v", err)
		}
		ext, err := combo.HorizontalExtent()
		if err != nil {
			t.Fatalf("failed to get horizontal extent, got %v", err)
		}
		if ext != 200 {
			t.Fatalf("expected horizontal extent 200 from ListBox without WS_HSCROLL, got %d", ext)
		}

		win.Close()

	}, nil)
}

func TestHorizontalExtentHScroll(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE | win32.WS_HSCROLL,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = combo.SetHorizontalExtent(200)
		if err != nil {
			t.Fatalf("failed to set horizontal extent: %v", err)
		}
		ext, err := combo.HorizontalExtent()
		if err != nil {
			t.Fatalf("failed to get horizontal extent: %v", err)
		}
		if ext != 200 {
			t.Fatalf("expected horizontal extent 200 from ListBox with WS_HSCROLL, got %d", ext)
		}

		win.Close()

	}, nil)
}

func TestItemHeight(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = combo.SetItemHeight(0, 30)
		if err != nil {
			t.Fatalf("failed to set item height: %v", err)
		}

		h, err := combo.ItemHeight(0)
		if err != nil {
			t.Fatalf("failed to get item height: %v", err)
		}
		if h != 30 {
			t.Fatalf("expected item height 30, got %d", h)
		}

		win.Close()

	}, nil)
}

func TestItemHeightOwnerdrawVar(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE | combobox.CBS_OWNERDRAWVARIABLE | combobox.CBS_HASSTRINGS,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = combo.SetItemHeight(0, 30)
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected ErrFailed when set item height 0 on a empty combo box, got %v", err)
		}

		combo.AppendItemString("item1")

		err = combo.SetItemHeight(0, 30)
		if err != nil {
			t.Fatalf("failed to set item height: %v", err)
		}

		h, err := combo.ItemHeight(0)
		if err != nil {
			t.Fatalf("failed to get item height: %v", err)
		}
		if h != 30 {
			t.Fatalf("expected item height 30, got %d", h)
		}

		win.Close()

	}, nil)
}

func TestOnMeasureItem(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE | combobox.CBS_OWNERDRAWVARIABLE | combobox.CBS_HASSTRINGS,
			OnMeasureItem: func(index int, itemData any) (width, height int) {
				return 100, 50
			},
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		combo.AppendItemString("item1")

		h, err := combo.ItemHeight(0)
		if err != nil {
			t.Fatalf("failed to get item height: %v", err)
		}
		if h != 50 {
			t.Fatalf("expected item height 50, got %d", h)
		}

		win.Close()

	}, nil)
}

func TestLocale(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		enUS := win32.MAKELCID(
			win32.MAKELANGID(win32.LANG_ENGLISH, win32.SUBLANG_ENGLISH_US),
			win32.SORT_DEFAULT,
		)

		old, err := combo.SetLocale(enUS)
		if err != nil {
			t.Fatalf("failed to set locale: %v", err)
		}
		if old == 0 {
			t.Fatalf("expected non-zero old locale, got %d", old)
		}

		old, err = combo.SetLocale(old)
		if err != nil {
			t.Fatalf("failed to set locale: %v", err)
		}
		if old != enUS {
			t.Fatalf("expected old locale 0x0409, got %d", old)
		}

		win.Close()

	}, nil)
}

func TestVisibleItem(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE | combobox.CBS_NOINTEGRALHEIGHT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		combo.AppendItemString("Item 1")
		combo.AppendItemString("Item 2")
		combo.AppendItemString("Item 3")
		combo.AppendItemString("Item 4")
		combo.AppendItemString("Item 5")
		combo.AppendItemString("Item 6")
		combo.AppendItemString("Item 7")

		combo.SetItemHeight(0, 30)

		v, err := combo.FirstVisible()
		if err != nil {
			t.Fatalf("failed to get first visible item: %v", err)
		}
		if v != 0 {
			t.Fatalf("expected first visible item index 0, got %d", v)
		}

		err = combo.EnsureVisible(3)
		if err != nil {
			t.Fatalf("failed to ensure item 3 visible: %v", err)
		}
		v, err = combo.FirstVisible()
		if err != nil {
			t.Fatalf("failed to get first visible item: %v", err)
		}
		if v != 3 {
			t.Fatalf("expected first visible item index 2 after EnsureVisible(3), got %d", v)
		}

		win.Close()

	}, nil)
}

func TestInitStorage(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE | combobox.CBS_NOINTEGRALHEIGHT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = combo.InitStorage(10, 20)
		if err != nil {
			t.Fatalf("failed to init storage: %v", err)
		}

		win.Close()

	}, nil)
}

func TestSelectItemString(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_SIMPLE | combobox.CBS_NOINTEGRALHEIGHT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		combo.AppendItemString("Apple")
		combo.AppendItemString("Banana")
		combo.AppendItemString("Cherry")

		err = combo.SelectItemString("Banana")
		if err != nil {
			t.Fatalf("failed to select item string: %v", err)
		}
		sel, err := combo.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if sel != 1 {
			t.Fatalf("expected current selection index 1 for 'Banana', got %d", sel)
		}

		err = combo.SelectItemStringIndex(2, "Apple")
		if err != nil {
			t.Fatalf("failed to select item string at index 2: %v", err)
		}
		sel, err = combo.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if sel != 0 {
			t.Fatalf("expected current selection index 0 for 'Apple' after selecting 'Apple' at index 2, got %d", sel)
		}

		err = combo.SelectItemString("xxxx")
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected ErrFailed when selecting non-existent item string, got %v", err)
		}

		win.Close()

	}, nil)
}

func TestComboboxInfo(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_DROPDOWN | combobox.CBS_NOINTEGRALHEIGHT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		info, err := combo.ComboBoxInfo()
		if err != nil {
			t.Fatalf("failed to get combo box info: %v", err)
		}
		if info.HwndCombo != combo.HWND() {
			t.Fatalf("expected ComboBox field to be %v, got %v", combo.HWND(), info.HwndCombo)
		}

		win.Close()

	}, nil)
}

func TestDroppedCtrlRect(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_DROPDOWN,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		rect, err := combo.DroppedCtrlRect()
		if err != nil {
			t.Fatalf("failed to get dropped control rect: %v", err)
		}
		if rect == nil {
			t.Fatalf("expected non-nil dropped control rect")
		}

		win.Close()

	}, nil)
}

func TestDroppedDown(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_DROPDOWN,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		dropped, err := combo.DroppedDown()
		if err != nil {
			t.Fatalf("failed to get dropped down state: %v", err)
		}
		if dropped {
			t.Fatalf("expected dropped down state to be false initially")
		}

		err = combo.ShowDropDown(true)
		if err != nil {
			t.Fatalf("failed to show drop down: %v", err)
		}

		dropped, err = combo.DroppedDown()
		if err != nil {
			t.Fatalf("failed to get dropped down state: %v", err)
		}
		if !dropped {
			t.Fatalf("expected dropped down state to be true after ShowDropDown(true)")
		}

		err = combo.ShowDropDown(false)
		if err != nil {
			t.Fatalf("failed to hide drop down: %v", err)
		}
		dropped, err = combo.DroppedDown()
		if err != nil {
			t.Fatalf("failed to get dropped down state: %v", err)
		}
		if dropped {
			t.Fatalf("expected dropped down state to be false after ShowDropDown(false)")
		}

		win.Close()

	}, nil)
}

func TestMinimumDroppedWidth(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_DROPDOWN,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		width, err := combo.MinDroppedWidth()
		if err != nil {
			t.Fatalf("failed to get minimum dropped width: %v", err)
		}
		if width != 400 {
			t.Fatalf("expected minimum dropped width 400 (initial width), got %d", width)
		}

		err = combo.SetMinDroppedWidth(500)
		if err != nil {
			t.Fatalf("failed to set minimum dropped width: %v", err)
		}
		width, err = combo.MinDroppedWidth()
		if err != nil {
			t.Fatalf("failed to get minimum dropped width: %v", err)
		}
		if width != 500 {
			t.Fatalf("expected minimum dropped width 500, got %d", width)
		}

		win.Close()

	}, nil)
}

func TestEditSel(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_DROPDOWN,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		start, end, err := combo.EditSel()
		if err != nil {
			t.Fatalf("failed to get edit selection: %v", err)
		}
		if start != 0 || end != 0 {
			t.Fatalf("expected edit selection start=0, end=0 initially, got start=%d, end=%d", start, end)
		}

		info, err := combo.ComboBoxInfo()
		if err != nil {
			t.Fatalf("failed to get combo box info: %v", err)
		}
		win32util.SetWindowText(info.HwndItem, "Hello, World!") // Set text in the edit control
		err = combo.SetEditSel(7, 12)                           // Select "World"
		if err != nil {
			t.Fatalf("failed to set edit selection: %v", err)
		}
		start, end, err = combo.EditSel()
		if err != nil {
			t.Fatalf("failed to get edit selection: %v", err)
		}
		if start != 7 || end != 12 {
			t.Fatalf("expected edit selection start=7, end=12 after SetEditSel, got start=%d, end=%d", start, end)
		}

		err = combo.SetEditSel(0, -1) // Select all
		if err != nil {
			t.Fatalf("failed to set edit selection: %v", err)
		}
		start, end, err = combo.EditSel()
		if err != nil {
			t.Fatalf("failed to get edit selection: %v", err)
		}
		if start != 0 || end != 13 {
			t.Fatalf("expected edit selection start=0, end=13 after SetEditSel(0, -1), got start=%d, end=%d", start, end)
		}

		win.Close()

	}, nil)
}

func TestExtendedUI(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_DROPDOWN,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		extended, err := combo.ExtendedUI()
		if err != nil {
			t.Fatalf("failed to get extended UI state: %v", err)
		}
		if extended {
			t.Fatalf("expected extended UI state to be false initially")
		}

		err = combo.SetExtendedUI(true)
		if err != nil {
			t.Fatalf("failed to set extended UI state: %v", err)
		}

		extended, err = combo.ExtendedUI()
		if err != nil {
			t.Fatalf("failed to get extended UI state: %v", err)
		}
		if !extended {
			t.Fatalf("expected extended UI state to be true after SetExtendedUI(true)")
		}

		win.Close()

	}, nil)
}

func TestLimitText(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		combo, err := combobox.New(&combobox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | combobox.CBS_DROPDOWN,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = combo.SetTextLimit(10)
		if err != nil {
			t.Fatalf("failed to set text limit: %v", err)
		}

		win.Close()

	}, nil)
}
