package gw_test

import (
	"testing"

	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

func TestBaseWindowImpl_AssignID(t *testing.T) {
	gw.Run(func(app *app.App) {
		win1, err := window.New(&window.Spec{})
		if err != nil {
			t.Fatalf("failed to create window: %v", err)
		}
		defer win1.Close()

		win2, err := window.New(&window.Spec{})
		if err != nil {
			t.Fatalf("failed to create window: %v", err)
		}
		defer win2.Close()

		btn1, err := button.New(&button.Spec{Parent: win1})
		if err != nil {
			t.Fatalf("failed to create button: %v", err)
		}

		btn11, err := button.New(&button.Spec{Parent: win1})
		if err != nil {
			t.Fatalf("failed to create button: %v", err)
		}

		btn2, err := button.New(&button.Spec{Parent: win2})

		err = win1.AssignID(btn1.HWND())
		if err != nil {
			t.Fatalf("failed to assign ID to btn1: %v", err)
		}
		id1, err := win32.GetWindowLongPtrW(btn1.HWND(), win32.GWLP_ID)
		if err != nil {
			t.Fatalf("failed to get ID of btn1: %v", err)
		}
		if id1 == 0 {
			t.Fatalf("ID of btn1 is 0")
		}

		err = win1.AssignID(btn1.HWND())
		if err == nil {
			t.Fatalf("expected error when assigning ID to btn1 again, but got nil")
		}

		err = win1.AssignID(btn11.HWND())
		if err != nil {
			t.Fatalf("failed to assign ID to btn11: %v", err)
		}
		id11, err := win32.GetWindowLongPtrW(btn11.HWND(), win32.GWLP_ID)
		if err != nil {
			t.Fatalf("failed to get ID of btn11: %v", err)
		}
		if id11 == 0 {
			t.Fatalf("ID of btn11 is 0")
		}
		if id11 == id1 {
			t.Fatalf("ID of btn11 is the same as ID of btn1")
		}

		err = win1.AssignID(btn2.HWND())
		if err == nil {
			t.Fatalf("expected error when assigning ID to btn2, but got nil")
		}

		err = win1.RemoveID(btn1.HWND())
		if err != nil {
			t.Fatalf("failed to remove ID of btn1: %v", err)
		}
		id1, err = win32.GetWindowLongPtrW(btn1.HWND(), win32.GWLP_ID)
		if err != nil {
			t.Fatalf("failed to get ID of btn1 after removing ID: %v", err)
		}
		if id1 != 0 {
			t.Fatalf("ID of btn1 is not 0 after removing ID")
		}

		err = win1.RemoveID(btn11.HWND())
		if err != nil {
			t.Fatalf("failed to remove ID of btn1: %v", err)
		}
		id11, err = win32.GetWindowLongPtrW(btn1.HWND(), win32.GWLP_ID)
		if err != nil {
			t.Fatalf("failed to get ID of btn1 after removing ID: %v", err)
		}
		if id11 != 0 {
			t.Fatalf("ID of btn1 is not 0 after removing ID")
		}

		app.Quit(0)

	}, nil)
}
