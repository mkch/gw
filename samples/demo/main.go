package main

import (
	"time"

	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/static"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -ico main.ico -manifest manifest.xml
//go:generate rsrc -arch 386 -ico main.ico -manifest manifest.xml

const TickerDuration = time.Millisecond * 100

// CreateMenu creates the window menu
func createMenu(ticker *time.Ticker, tickerStopped *bool) *menu.Menu {
	fileMenu := menu.New(true)
	var tickerMenuItem *menu.Item
	tickerMenuItem, _ = fileMenu.InsertItem(-1, &menu.ItemSpec{
		Title:    "&Stop Ticker",
		AccelKey: menu.AccelKey{Mod: menu.ModCtrl, VKeyCode: 'T'},
		OnClick: func() {
			if *tickerStopped {
				ticker.Reset(TickerDuration)
				tickerMenuItem.SetTitle("&Stop Ticker")
			} else {
				ticker.Stop()
				tickerMenuItem.SetTitle("&Start Ticker")
			}
			*tickerStopped = !*tickerStopped
		},
	})
	m := menu.New(false)
	m.InsertItem(-1, &menu.ItemSpec{Title: "&File", Submenu: fileMenu})
	return m
}

func main() {
	ticker := time.NewTicker(TickerDuration)
	tickerStopped := false

	win, _ := window.New(&window.Spec{
		Text:  "Hello, Go!",
		Style: win32.WS_OVERLAPPEDWINDOW,
		X:     win32.CW_USEDEFAULT,
		Width: 500, Height: 300,
		OnClose: func() { app.Quit(0) },
	})
	win.SetMenu(createMenu(ticker, &tickerStopped))

	timeStatic, _ := static.New(win.HWND(), &static.Spec{
		Text:        "Time",
		Style:       win32.WS_VISIBLE | static.SS_CENTER | static.SS_CENTERIMAGE,
		InParentDPI: true,
		X:           200, Y: 30,
		Width: 100, Height: 60,
	})

	button.New(win.HWND(), &button.Spec{
		Text:  "Hello",
		Style: win32.WS_VISIBLE,
		X:     200, Y: 120,
		Width: 100, Height: 60,
		InParentDPI: true,
		OnClick: func() {
			win32util.MessageBox(win.HWND(),
				"Hello GUI!", "Button clicked",
				win32.MB_ICONINFORMATION)
		},
	})

	win.Show(win32.SW_SHOW)

	go func() {
		for {
			str := (<-ticker.C).Local().Format("15:04:05")
			// Run SetText() in UI goroutine.
			app.Post(func() { timeStatic.SetText(str) })
		}
	}()
	app.Run()
}
