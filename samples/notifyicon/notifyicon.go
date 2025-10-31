package main

import (
	"bytes"
	"unsafe"

	"github.com/mkch/gg"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/notifyicon"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml -ico main.ico
//go:generate rsrc -arch 386 -manifest manifest.xml -ico main.ico

func main() {
	win := gg.Must(window.New(&window.Spec{
		Text:    "Test Notify Icon",
		Style:   win32.WS_OVERLAPPEDWINDOW,
		X:       metrics.Px(win32.CW_USEDEFAULT),
		Width:   metrics.Dip(500),
		Height:  metrics.Dip(300),
		OnClose: func() { app.Quit(0) },
	}))

	tooltip := gg.Must(window.New(&window.Spec{
		Text:    "Tooltip",
		Style:   win32.WS_POPUP | win32.WS_BORDER | win32.WS_CAPTION,
		ExStyle: win32.WS_EX_TOOLWINDOW,
		Width:   metrics.Dip(200),
		Height:  metrics.Dip(100),
		OnClose: func() {},
	}))
	tooltip.SetWndProc(func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevWndProc win32.WndProc) win32.LRESULT {
		switch message {
		case win32.WM_CLOSE:
			tooltip.Show(win32.SW_HIDE)
			return 0
		}
		return prevWndProc(hwnd, message, wParam, lParam)
	})

	icon := gg.Must(win32.LoadIconW(gg.Must(win32.GetModuleHandleW[win32.HINSTANCE](nil)), (*win32.WCHAR)(unsafe.Add(unsafe.Pointer(nil), 2))))
	const iconID = 2
	var notifyIcon *notifyicon.NotifyIcon

	fileMenu := menu.New(true)
	fileMenu.InsertItem(-1, &menu.ItemSpec{
		Title: "Change Tip",
		OnClick: func() {
			b := bytes.NewBufferString("New Tip")
			for range 128 {
				b.WriteRune('.')
			}
			gg.MustOK(notifyIcon.StartModify().SetTip(b.String()).Apply())
		},
	})
	fileMenu.InsertItem(-1, &menu.ItemSpec{
		Title: "Custom Tip",
		OnClick: func() {
			gg.MustOK(notifyIcon.StartModify().SetTip("").Apply())

		},
	})
	fileMenu.InsertItem(-1, &menu.ItemSpec{
		Title: "Change Icon",
		OnClick: func() {
			icon := gg.Must(win32.LoadIconW(0, (*win32.WCHAR)(unsafe.Add(unsafe.Pointer(nil), win32.IDI_INFORMATION))))
			gg.MustOK(notifyIcon.StartModify().SetIcon(icon).Apply())
		},
	})
	fileMenu.InsertItem(-1, &menu.ItemSpec{
		Title: "Change Notification",
		OnClick: func() {
			gg.MustOK(notifyIcon.StartModify().
				SetNotify(&notifyicon.NotifySpec{
					Icon:    notifyicon.BI_ERROR,
					Title:   "New Title",
					Message: "New Message",
					NoSound: true,
				}).Apply())
		},
	})
	fileMenu.InsertItem(-1, &menu.ItemSpec{
		Title:   "E&xit",
		OnClick: func() { win32.SendMessageW(win.HWND(), win32.WM_CLOSE, 0, 0) },
	})

	notifyIcon = gg.Must(notifyicon.New(win, iconID, &notifyicon.Spec{
		Icon: icon,
		Tip:  "Notify icon tip",
		Notify: &notifyicon.NotifySpec{
			Icon:    notifyicon.BI_INFO,
			Title:   "Notification title",
			Message: "This is a notification",
		},
		OnEvent: func(id win32.WORD, event notifyicon.Event, eventX, eventY win32.SHORT) {
			if id != win32.WORD(notifyIcon.ID()) {
				panic("wrong id")
			}
			switch event {
			case notifyicon.MIN_CONTEXTMENU:
				win32.SetForegroundWindow(win.HWND())
				win.TrackPopupMenu(fileMenu, nil)
			case notifyicon.NIN_SELECT:
				win.Show(win32.SW_SHOWNORMAL)
				win32.SetForegroundWindow(win.HWND())
			case notifyicon.NIN_POPUPOPEN:
				var rect win32.RECT
				win32.GetWindowRect(tooltip.HWND(), &rect)
				win32.SetWindowPos(tooltip.HWND(), win32.HWND_TOPMOST,
					win32.INT(eventX)-win32.INT(rect.Width()), win32.INT(eventY)-win32.INT(win32.SHORT(rect.Height())),
					0, 0, win32.SWP_NOSIZE)
				win32.SetForegroundWindow(tooltip.HWND())
				tooltip.Show(win32.SW_SHOWNORMAL)
			case notifyicon.NIN_POPUPCLOSE:
				tooltip.Show(win32.SW_HIDE)
			}
		},
	}))

	app.Run()
}
