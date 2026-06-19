package listbox_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/listbox"
	"github.com/mkch/gw/listbox/listcombo"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		if count, err := list.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != 0 {
			t.Fatalf("expected item count 0, got %d", count)
		}

		if i, err := list.AppendItemString("Item 1"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		} else if i != 0 {
			t.Fatalf("expected index 0, got %d", i)
		}

		if i, err := list.AppendItemString("Item 2"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		} else if i != 1 {
			t.Fatalf("expected index 1, got %d", i)
		}

		if i, err := list.AppendItemString("Item 3"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		} else if i != 2 {
			t.Fatalf("expected index 2, got %d", i)
		}

		if str, err := list.ItemString(0); err != nil {
			t.Fatalf("failed to get item string: %v", err)
		} else if str != "Item 1" {
			t.Fatalf("expected 'Item 1', got '%s'", str)
		}

		if str, err := list.ItemString(1); err != nil {
			t.Fatalf("failed to get item string: %v", err)
		} else if str != "Item 2" {
			t.Fatalf("expected 'Item 2', got '%s'", str)
		}

		if str, err := list.ItemString(2); err != nil {
			t.Fatalf("failed to get item string: %v", err)
		} else if str != "Item 3" {
			t.Fatalf("expected 'Item 3', got '%s'", str)
		}

		if count, err := list.ItemCount(); err != nil {
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | (listbox.LBS_STANDARD &^ listbox.LBS_SORT),
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		if _, err := list.AppendItemString("Item 1"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}
		if _, err := list.AppendItemString("Item 3"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}

		if count, err := list.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != 2 {
			t.Fatalf("expected item count 2, got %d", count)
		}

		if err := list.InsertItemString(1, "Item 2"); err != nil {
			t.Fatalf("failed to insert item: %v", err)
		}

		expectedItems := []string{"Item 1", "Item 2", "Item 3"}
		for i, expected := range expectedItems {
			if str, err := list.ItemString(i); err != nil {
				t.Fatalf("failed to get item string at index %d: %v", i, err)
			} else if str != expected {
				t.Fatalf("expected '%s' at index %d, got '%s'", expected, i, str)
			}
		}

		if count, err := list.ItemCount(); err != nil {
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | (listbox.LBS_STANDARD &^ listbox.LBS_SORT),
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		if _, err := list.AppendItemString("Item 1"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}
		if _, err := list.AppendItemString("Item 4"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}

		if count, err := list.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != 2 {
			t.Fatalf("expected item count 2, got %d", count)
		}

		insertItems := []string{"Item 2", "Item 3"}
		for i, item := range insertItems {
			if err := list.InsertItemString(1+i, item); err != nil {
				t.Fatalf("failed to insert item: %v", err)
			}
		}

		expectedItems := []string{"Item 1", "Item 2", "Item 3", "Item 4"}
		for i, expected := range expectedItems {
			if str, err := list.ItemString(i); err != nil {
				t.Fatalf("failed to get item string at index %d: %v", i, err)
			} else if str != expected {
				t.Fatalf("expected '%s' at index %d, got '%s'", expected, i, str)
			}
		}

		if count, err := list.ItemCount(); err != nil {
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | (listbox.LBS_STANDARD &^ listbox.LBS_SORT),
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
			if _, err := list.AppendItemString(item); err != nil {
				t.Fatalf("failed to append item: %v", err)
			}
		}

		if count, err := list.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != len(items) {
			t.Fatalf("expected item count %d, got %d", len(items), count)
		}

		if err := list.DeleteItem(1); err != nil {
			t.Fatalf("failed to delete item: %v", err)
		}

		expectedAfterDelete := []string{"Item 1", "Item 3", "Item 4"}
		for i, expected := range expectedAfterDelete {
			if str, err := list.ItemString(i); err != nil {
				t.Fatalf("failed to get item string at index %d: %v", i, err)
			} else if str != expected {
				t.Fatalf("expected '%s' at index %d, got '%s'", expected, i, str)
			}
		}

		if _, err := list.ItemString(3); err == nil {
			t.Fatalf("expected error when getting removed item at index 3")
		}

		if count, err := list.ItemCount(); err != nil {
			t.Fatalf("failed to get item count: %v", err)
		} else if count != len(expectedAfterDelete) {
			t.Fatalf("expected item count %d, got %d", len(expectedAfterDelete), count)
		}

		if err := list.DeleteAllItems(); err != nil {
			t.Fatalf("failed to delete all items: %v", err)
		}

		if _, err := list.ItemString(0); err == nil {
			t.Fatalf("expected error when getting item from empty list")
		}

		if count, err := list.ItemCount(); err != nil {
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | (listbox.LBS_STANDARD &^ listbox.LBS_SORT),
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
			if _, err := list.AppendItemString(item); err != nil {
				t.Fatalf("failed to append item: %v", err)
			}
		}

		if data, err := list.ItemData(0); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != uintptr(0) {
			t.Fatalf("expected item data uintptr(0) at index 0 before any SetItemData, got %v", data)
		}

		if err := list.SetItemData(0, "abc"); err != nil {
			t.Fatalf("failed to set item data: %v", err)
		}
		if err := list.SetItemData(1, 200); err != nil {
			t.Fatalf("failed to set item data: %v", err)
		}

		if data, err := list.ItemData(0); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != "abc" {
			t.Fatalf(`expected item data "abc" at index 0, got %v`, data)
		}

		if data, err := list.ItemData(1); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != 200 {
			t.Fatalf(`expected item data 200 at index 1, got %v`, data)
		}

		if data, err := list.ItemData(2); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != uintptr(0) {
			t.Fatalf("expected item data uintptr(0) at index 2 (never set), got %v", data)
		}

		if err := list.SetItemData(0, nil); err != nil {
			t.Fatalf("failed to set item data to nil: %v", err)
		}

		if data, err := list.ItemData(0); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != uintptr(0) {
			t.Fatalf("expected item data uintptr(0) at index 0 after SetItemData(0, nil), got %v", data)
		}

		// Verify that SetItemData does not affect item strings.
		for i, expected := range items {
			if str, err := list.ItemString(i); err != nil {
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | (listbox.LBS_STANDARD &^ listbox.LBS_SORT),
			X:      metrics.Dip(10),
			Y:      metrics.Dip(10),
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		if _, err := list.AppendItemString("Item 1"); err != nil {
			t.Fatalf("failed to append item: %v", err)
		}

		// Set external item data via raw Win32 message.
		if _, err := win32.SendMessageW(list.HWND(), listbox.LB_SETITEMDATA,
			win32.WPARAM(0), win32.LPARAM(uintptr(1))); err != nil {
			t.Fatalf("failed to set external item data: %v", err)
		}

		if data, err := list.ItemData(0); err != nil {
			t.Fatalf("failed to get item data: %v", err)
		} else if data != uintptr(1) {
			t.Fatalf("expected item data uintptr(1) at index 0, got %v", data)
		}

		if err := list.SetItemData(0, "test"); !errors.Is(err, listcombo.ErrExternalItemData) {
			t.Fatalf("expected ErrExternalItemData from SetItemData, got %v", err)
		}

		// After external data is still there, verify SetItemData(nil) also fails.
		if err := list.SetItemData(0, nil); !errors.Is(err, listcombo.ErrExternalItemData) {
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

		unsorted, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | (listbox.LBS_STANDARD &^ listbox.LBS_SORT),
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

		sorted, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
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

		unsorted, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | (listbox.LBS_STANDARD &^ listbox.LBS_SORT),
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

		var list *listbox.ListBox
		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_OWNERDRAWFIXED,
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

		list.AppendItem("C")
		list.AppendItem("A")
		list.AppendItem("B")

		win32.SendMessageW(list.HWND(), listbox.LB_ADDSTRING, 0, 100)

		data, err := list.ItemData(0)
		if err != nil {
			t.Fatalf("failed to get item data: %v", err)
		}
		if data != uintptr(100) {
			t.Fatalf(`expected item data uintptr(100) at index 0 due to OnCompareItem, got %#v`, data)
		}

		data, err = list.ItemData(1)
		if err != nil {
			t.Fatalf("failed to get item data: %v", err)
		}
		if data != "A" {
			t.Fatalf(`expected item data "A" at index 1 due to OnCompareItem, got %#v`, data)
		}

		data, err = list.ItemData(2)
		if err != nil {
			t.Fatalf("failed to get item data: %v", err)
		}
		if data != "B" {
			t.Fatalf(`expected item data "B" at index 2 due to OnCompareItem, got %#v`, data)
		}

		data, err = list.ItemData(3)
		if err != nil {
			t.Fatalf("failed to get item data: %v", err)
		}
		if data != "C" {
			t.Fatalf(`expected item data "C" at index 2 due to OnCompareItem, got %#v`, data)
		}

		_, err = list.SelectedItems()
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected LB_ERR from SelectedItems on owner-drawn ListBox without LBS_MULTIPLESEL, got %v", err)
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Item 1")
		list.AppendItemString("Item 2")
		list.AppendItemString("Item 3")

		i, err := list.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if i != -1 {
			t.Fatalf("expected no selection (-1), got %d", i)
		}

		err = list.SetCurSelected(1)
		if err != nil {
			t.Fatalf("failed to set current selection: %v", err)
		}
		i, err = list.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if i != 1 {
			t.Fatalf("expected selection index 1, got %d", i)
		}

		err = list.SetCurSelected(-1)
		if err != nil {
			t.Fatalf("failed to set current selection: %v", err)
		}
		i, err = list.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if i != -1 {
			t.Fatalf("expected no selection (-1), got %d", i)
		}

		err = list.SetSelected(0, true)
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected LB_ERR from SetSelected on single-select ListBox, got %v", err)
		}

		win.Close()

	}, nil)
}

func TestMultiSel(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_MULTIPLESEL,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Item 1")
		list.AppendItemString("Item 2")
		list.AppendItemString("Item 3")

		i, err := list.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if i != 0 { // 0 for multi-select ListBox
			t.Fatalf("expected no selection(0), got %d", i)
		}

		s, err := list.SelectedItems()
		if err != nil {
			t.Fatalf("failed to get selected items: %v", err)
		}
		if len(s) != 0 {
			t.Fatalf("expected no selected items, got %d", len(s))
		}

		err = list.SetCurSelected(1)
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected LB_ERR from SetCurSelected on multi-select ListBox, got %v", err)
		}

		err = list.SetSelected(0, true)
		if err != nil {
			t.Fatalf("failed to set item 0 selected: %v", err)
		}
		selected, err := list.Selected(0)
		if err != nil {
			t.Fatalf("failed to get selected item 0: %v", err)
		}
		if !selected {
			t.Fatalf("expected item 0 to be selected")
		}

		s, err = list.SelectedItems()
		if err != nil {
			t.Fatalf("failed to get selected items: %v", err)
		}
		if !slices.Equal(s, []int32{0}) {
			t.Fatalf("expected selected items [0], got %v", s)
		}

		err = list.SetCurSelected(1)
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected LB_ERR from SetCurSelected on multi-select ListBox, got %v", err)
		}

		err = list.SetSelected(1, true)
		if err != nil {
			t.Fatalf("failed to set item 1 selected: %v", err)
		}
		selected, err = list.Selected(1)
		if err != nil {
			t.Fatalf("failed to get selected item 1: %v", err)
		}
		if !selected {
			t.Fatalf("expected item 1 to be selected")
		}

		s, err = list.SelectedItems()
		if err != nil {
			t.Fatalf("failed to get selected items: %v", err)
		}
		if !slices.Equal(s, []int32{0, 1}) {
			t.Fatalf("expected selected items [0, 1], got %v", s)
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = list.SetHorizontalExtent(200)
		if err != nil {
			t.Fatalf("failed to set horizontal extent, got: %v", err)
		}
		ext, err := list.HorizontalExtent()
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | win32.WS_HSCROLL,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = list.SetHorizontalExtent(200)
		if err != nil {
			t.Fatalf("failed to set horizontal extent: %v", err)
		}
		ext, err := list.HorizontalExtent()
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = list.SetItemHeight(0, 30)
		if err != nil {
			t.Fatalf("failed to set item height: %v", err)
		}

		h, err := list.ItemHeight(0)
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_OWNERDRAWVARIABLE | listbox.LBS_HASSTRINGS,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = list.SetItemHeight(0, 30)
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected ErrFailed when set item height 0 on a empty list box, got %v", err)
		}

		list.AppendItemString("item1")

		err = list.SetItemHeight(0, 30)
		if err != nil {
			t.Fatalf("failed to set item height: %v", err)
		}

		h, err := list.ItemHeight(0)
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_OWNERDRAWVARIABLE | listbox.LBS_HASSTRINGS,
			OnMeasureItem: func(index int, itemData any) (width, height int) {
				return 100, 50
			},
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("item1")

		h, err := list.ItemHeight(0)
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		enUS := win32.MAKELCID(
			win32.MAKELANGID(win32.LANG_ENGLISH, win32.SUBLANG_ENGLISH_US),
			win32.SORT_DEFAULT,
		)

		old, err := list.SetLocale(enUS)
		if err != nil {
			t.Fatalf("failed to set locale: %v", err)
		}
		if old == 0 {
			t.Fatalf("expected non-zero old locale, got %d", old)
		}

		old, err = list.SetLocale(old)
		if err != nil {
			t.Fatalf("failed to set locale: %v", err)
		}
		if old != enUS {
			t.Fatalf("expected old locale 0x0409, got %d", old)
		}

		win.Close()

	}, nil)
}

func TestItemRect(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		rect, err := list.ItemRect(0)
		if err != nil {
			t.Fatalf("failed to get item rect: %v", err)
		}
		if rect == nil {
			t.Fatalf("expected non-nil rect, got nil")
		}

		win.Close()

	}, nil)
}

func TestColumns(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_NOINTEGRALHEIGHT | listbox.LBS_MULTICOLUMN | win32.WS_HSCROLL,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.SetItemHeight(0, 40)
		list.SetColumnWidth(100)

		list.AppendItemString("Item 1")
		list.AppendItemString("Item 2")
		list.AppendItemString("Item 3")

		n, err := list.ItemsPerColumn()
		if err != nil {
			t.Fatalf("failed to get items per column: %v", err)
		}
		if n != 2 {
			t.Fatalf("expected 2 items per column, got %d", n)
		}

		win.Close()

	}, nil)
}

func TestAnchorCart(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Item 1")
		list.AppendItemString("Item 2")
		list.AppendItemString("Item 3")

		a, err := list.AnchorIndex()
		if err != nil {
			t.Fatalf("failed to get anchor index: %v", err)
		}
		if a != -1 {
			t.Fatalf("expected anchor index -1, got %d", a)
		}
		err = list.SetAnchorIndex(1)
		if err != nil {
			t.Fatalf("failed to set anchor index: %v", err)
		}
		a, err = list.AnchorIndex()
		if err != nil {
			t.Fatalf("failed to get anchor index: %v", err)
		}
		if a != 1 {
			t.Fatalf("expected anchor index 1, got %d", a)
		}

		c, err := list.CaretIndex()
		if err != nil {
			t.Fatalf("failed to get caret index: %v", err)
		}
		if c != 0 {
			t.Fatalf("expected caret index 0, got %d", c)
		}
		err = list.SetCaretIndex(2)
		if err != nil {
			t.Fatalf("failed to set caret index: %v", err)
		}
		c, err = list.CaretIndex()
		if err != nil {
			t.Fatalf("failed to get caret index: %v", err)
		}
		if c != 2 {
			t.Fatalf("expected caret index 2, got %d", c)
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_NOINTEGRALHEIGHT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Item 1")
		list.AppendItemString("Item 2")
		list.AppendItemString("Item 3")
		list.AppendItemString("Item 4")
		list.AppendItemString("Item 5")
		list.AppendItemString("Item 6")
		list.AppendItemString("Item 7")

		list.SetItemHeight(0, 30)

		v, err := list.FirstVisible()
		if err != nil {
			t.Fatalf("failed to get first visible item: %v", err)
		}
		if v != 0 {
			t.Fatalf("expected first visible item index 0, got %d", v)
		}

		err = list.EnsureVisible(3)
		if err != nil {
			t.Fatalf("failed to ensure item 3 visible: %v", err)
		}
		v, err = list.FirstVisible()
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_NOINTEGRALHEIGHT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = list.InitStorage(10, 20)
		if err != nil {
			t.Fatalf("failed to init storage: %v", err)
		}

		win.Close()

	}, nil)
}

func TestItemFromPoint(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_NOINTEGRALHEIGHT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Item 1")
		list.AppendItemString("Item 2")
		list.AppendItemString("Item 3")

		list.SetItemHeight(0, 30)

		index, err := list.ItemFromPoint(10, 10)
		if err != nil {
			t.Fatalf("failed to get item from point: %v", err)
		}
		if index != 0 {
			t.Fatalf("expected item index 0 from point (10, 10), got %d", index)
		}

		index, err = list.ItemFromPoint(10, 40)
		if err != nil {
			t.Fatalf("failed to get item from point: %v", err)
		}
		if index != 1 {
			t.Fatalf("expected item index 1 from point (10, 40), got %d", index)
		}

		index, err = list.ItemFromPoint(500, 200)
		if err != nil {
			t.Fatalf("failed to get item from point: %v", err)
		}
		if index != 2 {
			t.Fatalf("expected item index 2 from point (500, 200), got %d", index)
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

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_NOINTEGRALHEIGHT,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Apple")
		list.AppendItemString("Banana")
		list.AppendItemString("Cherry")

		err = list.SelectItemString("Banana")
		if err != nil {
			t.Fatalf("failed to select item string: %v", err)
		}
		sel, err := list.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if sel != 1 {
			t.Fatalf("expected current selection index 1 for 'Banana', got %d", sel)
		}

		err = list.SelectItemStringIndex(2, "Apple")
		if err != nil {
			t.Fatalf("failed to select item string at index 2: %v", err)
		}
		sel, err = list.CurSelected()
		if err != nil {
			t.Fatalf("failed to get current selection: %v", err)
		}
		if sel != 0 {
			t.Fatalf("expected current selection index 0 for 'Apple' after selecting 'Apple' at index 2, got %d", sel)
		}

		err = list.SelectItemString("xxxx")
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected ErrFailed when selecting non-existent item string, got %v", err)
		}

		win.Close()

	}, nil)
}

func TestSelectRange(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_MULTIPLESEL,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Apple")
		list.AppendItemString("Banana")
		list.AppendItemString("Cherry")

		err = list.SelectRange(true, 0, 1)
		if err != nil {
			t.Fatalf("failed to select range: %v", err)
		}
		sel, err := list.SelectedItems()
		if err != nil {
			t.Fatalf("failed to get selected items: %v", err)
		}
		if !slices.Equal(sel, []int32{0, 1}) {
			t.Fatalf("expected selected items [0, 1] for range selection, got %v", sel)
		}

		err = list.SelectRange(true, 1, 2)
		if err != nil {
			t.Fatalf("failed to select range: %v", err)
		}
		sel, err = list.SelectedItems()
		if err != nil {
			t.Fatalf("failed to get selected items: %v", err)
		}
		if !slices.Equal(sel, []int32{0, 1, 2}) {
			t.Fatalf("expected selected items [0, 1, 2] for range selection, got %v", sel)
		}

		err = list.SelectRange(false, 0, 2)
		if err != nil {
			t.Fatalf("failed to deselect range: %v", err)
		}
		sel, err = list.SelectedItems()
		if err != nil {
			t.Fatalf("failed to get selected items: %v", err)
		}
		if n := len(sel); n != 0 {
			t.Fatalf("expected no selected items after deselecting range, got %d", n)
		}

		win.Close()

	}, nil)
}

func TestSetCount(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		list, err := listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = list.SetCount(5)
		if !errors.Is(err, listcombo.ErrFailed) {
			t.Fatalf("expected ErrFailed when setting count on ListBox without LBS_NODATA, got %v", err)
		}
		list.Destroy()

		list, err = listbox.New(&listbox.Spec{
			Parent: win,
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_OWNERDRAWFIXED | listbox.LBS_NODATA,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		err = list.SetCount(5)
		if err != nil {
			t.Fatalf("failed to set count on ListBox with LBS_NODATA: %v", err)
		}
		n, err := list.ItemCount()
		if err != nil {
			t.Fatalf("failed to get item count: %v", err)
		}
		if n != 5 {
			t.Fatalf("expected item count 5 after SetCount, got %d", n)
		}

		win.Close()
	}, nil)
}

func TestOnItemDoubleClick(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		var itemDbClickCalled, setFocusCalled, killFocusCalled bool
		list, err := listbox.New(&listbox.Spec{
			Parent:            win,
			Width:             metrics.Px(400),
			Height:            metrics.Px(100),
			Style:             win32.WS_VISIBLE | listbox.LBS_STANDARD,
			OnItemDoubleClick: func() { itemDbClickCalled = true },
			OnSetFocus:        func() { setFocusCalled = true },
			OnKillFocus:       func() { killFocusCalled = true },
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list2, err := listbox.New(&listbox.Spec{
			Parent: win,
			Y:      metrics.Px(105),
			Width:  metrics.Px(400),
			Height: metrics.Px(100),
			Style:  win32.WS_VISIBLE | listbox.LBS_STANDARD,
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Apple")
		list2.AppendItemString("Banana")

		rc, err := list.ItemRect(0)
		win32.SendMessageW(list.HWND(), win32.WM_LBUTTONDBLCLK, 0, win32.LPARAM(win32.MAKELONG(uint16(rc.Left), uint16(rc.Top))))

		if !itemDbClickCalled {
			t.Fatalf("expected OnItemDoubleClick listener to be called, but it was not")
		}

		if !setFocusCalled {
			t.Fatalf("expected OnSetFocus listener to be called, but it was not")
		}

		if killFocusCalled {
			t.Fatal("expected OnKillFocus listener to not be called, but it was")
		}

		rc, err = list2.ItemRect(0)
		win32.SendMessageW(list2.HWND(), win32.WM_LBUTTONDOWN, 0, win32.LPARAM(win32.MAKELONG(uint16(rc.Left+1), uint16(rc.Top+1))))

		if !killFocusCalled {
			t.Fatal("expected OnKillFocus listener to be called, but it was not")
		}

		win.Close()
	}, nil)
}

func TestSelChange(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:     win32.WS_OVERLAPPEDWINDOW,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		win.Show(win32.SW_SHOW)

		var selChangeCalled bool
		list, err := listbox.New(&listbox.Spec{
			Parent:      win,
			Width:       metrics.Px(400),
			Height:      metrics.Px(100),
			Style:       win32.WS_VISIBLE | win32.WS_BORDER | listbox.LBS_STANDARD,
			OnSelChange: func() { selChangeCalled = true },
		})
		if err != nil {
			t.Fatalf("failed to create ListBox: %v", err)
		}

		list.AppendItemString("Apple")
		list.AppendItemString("Banana")
		rc1, err := list.ItemRect(0)
		rc2, err := list.ItemRect(1)

		win32.SendMessageW(list.HWND(), win32.WM_LBUTTONDOWN, 0, win32.LPARAM(win32.MAKELONG(uint16(rc1.Left), uint16(rc1.Top))))
		win32.SendMessageW(list.HWND(), win32.WM_LBUTTONUP, 0, win32.LPARAM(win32.MAKELONG(uint16(rc1.Left), uint16(rc1.Top))))

		if !selChangeCalled {
			t.Fatalf("expected OnSelChange listener to be called, but it was not")
		}

		selChangeCalled = false
		win32.SendMessageW(list.HWND(), win32.WM_LBUTTONDOWN, 0, win32.LPARAM(win32.MAKELONG(uint16(rc2.Left), uint16(rc2.Top+1))))
		win32.SendMessageW(list.HWND(), win32.WM_LBUTTONUP, 0, win32.LPARAM(win32.MAKELONG(uint16(rc2.Left), uint16(rc2.Top+1))))

		if !selChangeCalled {
			t.Fatalf("expected OnSelChange listener to be called, but it was not")
		}

		win.Close()

	}, nil)
}
