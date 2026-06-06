package win32util_test

import (
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

func Test_SetClientSize_Child(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:  win32.WS_POPUPWINDOW,
			Width:  metrics.Px(300),
			Height: metrics.Px(300),
		}))

		btn := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
			Width:  metrics.Px(100),
			Height: metrics.Px(100),
		}))

		chkerr.MustOK(win32util.SetClientSize(btn.HWND(), 100, 100))

		clientRect := chkerr.Must(btn.GetClientRect())
		if clientRect.Width() != 100 || clientRect.Height() != 100 {
			t.Errorf("unexpected client size: %v", clientRect)
		}

		chkerr.MustOK(win32util.SetClientSize(btn.HWND(), 110, 120))

		clientRect = chkerr.Must(btn.GetClientRect())
		if clientRect.Width() != 110 || clientRect.Height() != 120 {
			t.Errorf("unexpected client size: %v", clientRect)
		}

		app.Quit(0)

	}, func(app *app.App) { app.DestroyAllWindows() })
}

func Test_SetClientSize_Popup(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:  win32.WS_POPUPWINDOW,
			Width:  metrics.Px(300),
			Height: metrics.Px(300),
		}))

		chkerr.MustOK(win32util.SetClientSize(win.HWND(), 100, 100))

		clientRect := chkerr.Must(win.GetClientRect())
		if clientRect.Width() != 100 || clientRect.Height() != 100 {
			t.Errorf("unexpected client size: %v", clientRect)
		}

		chkerr.MustOK(win32util.SetClientSize(win.HWND(), 110, 120))

		clientRect = chkerr.Must(win.GetClientRect())
		if clientRect.Width() != 110 || clientRect.Height() != 120 {
			t.Errorf("unexpected client size: %v", clientRect)
		}

		app.Quit(0)

	}, func(app *app.App) { app.DestroyAllWindows() })
}

func Test_SetClientSize_PopupWithMenu(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Style:  win32.WS_CAPTION | win32.WS_POPUPWINDOW,
			Width:  metrics.Px(300),
			Height: metrics.Px(300),
		}))

		menuBar := menu.New(false)
		menuBar.InsertItem(-1, &menu.ItemSpec{Title: "111111111111"})
		menuBar.InsertItem(-1, &menu.ItemSpec{Title: "222222222222"})
		menuBar.InsertItem(-1, &menu.ItemSpec{Title: "3333333333333"})
		menuBar.InsertItem(-1, &menu.ItemSpec{Title: "44444444444444"})
		menuBar.InsertItem(-1, &menu.ItemSpec{Title: "55555555555555"})

		chkerr.MustOK(win.SetMenu(menuBar))

		chkerr.MustOK(win32util.SetClientSize(win.HWND(), 100, 100))

		clientRect := chkerr.Must(win.GetClientRect())
		if clientRect.Width() != 100 || clientRect.Height() != 100 {
			t.Errorf("unexpected client size: %v", clientRect)
		}

		chkerr.MustOK(win32util.SetClientSize(win.HWND(), 110, 120))

		clientRect = chkerr.Must(win.GetClientRect())
		if clientRect.Width() != 110 || clientRect.Height() != 120 {
			t.Errorf("unexpected client size: %v", clientRect)
		}

		app.Quit(0)

	}, func(app *app.App) { app.DestroyAllWindows() })
}
