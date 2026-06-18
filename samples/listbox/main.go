package main

import (
	"fmt"
	"os"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/events"
	"github.com/mkch/gw/layout"
	"github.com/mkch/gw/listbox"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

func main() {
	os.Exit(gw.Run(ui, nil))
}

func ui(app *app.App) {

	win := chkerr.Must(window.New(&window.Spec{
		Text:   "Static demo",
		Style:  win32.WS_OVERLAPPEDWINDOW,
		X:      gw.CW_USEDEFAULT,
		Width:  metrics.Dip(500),
		Height: metrics.Dip(300),
		OnDestroy: func() {
			app.Quit(0)
		},
	}))

	list := chkerr.Must(listbox.New(&listbox.Spec{
		Parent:    win,
		Style:     win32.WS_VISIBLE | listbox.LBS_STANDARD | listbox.LBS_NOINTEGRALHEIGHT | listbox.LBS_USETABSTOPS,
		Width:     metrics.Dip(200),
		Height:    metrics.Dip(200),
		Draggable: true,
	}))

	list.SetTabStops([]int32{4 * 8, 4 * 16})

	list.AppendItemString("Item\t1\tApple tree")
	list.AppendItemString("Item\t2\tBanana tree")
	list.AppendItemString("Item\t3\tCherry tree")
	list.AppendItemString("Item\t4\tDate tree")
	list.AppendItemString("Item\t5\tElderberry tree")

	btnGetSel := chkerr.Must(button.New(&button.Spec{
		Parent: win,
		Text:   "Get selected index",
		Style:  win32.WS_VISIBLE | win32.WS_TABSTOP,
		Width:  metrics.Dip(150),
		Height: metrics.Dip(50),
		OnClick: func() {
			fmt.Println("CurSel:", chkerr.Must(list.CurSelected()))
			fmt.Println("AnchorIndex:", chkerr.Must(list.AnchorIndex()))
			fmt.Println("CaretIndex:", chkerr.Must(list.CaretIndex()))
			fmt.Println()
		},
	}))

	btnSetSel := chkerr.Must(button.New(&button.Spec{
		Parent: win,
		Text:   "Select last item",
		Style:  win32.WS_VISIBLE | win32.WS_TABSTOP,
		Width:  metrics.Dip(150),
		Height: metrics.Dip(50),
		OnClick: func() {
			fmt.Println(list.SetCurSelected(2))
		},
	}))

	root := layout.Padding{
		Left:   metrics.Dip(10),
		Top:    metrics.Dip(10),
		Right:  metrics.Dip(10),
		Bottom: metrics.Dip(10),
		Child: &layout.Row{
			Children: []layout.Widget{
				&layout.Expanded{
					Child: &layout.Intrinsic{Hwnd: list.HWND()},
				},
				&layout.Sized{Width: metrics.Dip(10)},
				&layout.Column{
					MainAxisSize:  layout.AxisSizeMax,
					MainAxisAlign: layout.AlignCenter,
					Children: []layout.Widget{
						&layout.Intrinsic{Hwnd: btnGetSel.HWND()},
						&layout.Sized{Height: metrics.Dip(10)},
						&layout.Intrinsic{Hwnd: btnSetSel.HWND()},
					},
				},
			},
		},
	}

	e := chkerr.Must(layout.Build(&root))

	layout.PerformWindow(e, win.HWND())
	win.SetOnSizeListener(func(event events.SizeEvent) {
		layout.PerformWindow(e, win.HWND())
	})

	win.Show(win32.SW_SHOW)
}
