package layout_test

import (
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/layout"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
	"golang.org/x/exp/constraints"
)

func abs[T constraints.Signed](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func TestCenter(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Text:      "Test Center Layout",
			Style:     win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(500),
			Height:    metrics.Dip(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		btn := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Dip(300),
			Height: metrics.Dip(100),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		center := &layout.Center{
			Hwnd: win.HWND(),
			Item: &layout.Window{
				Hwnd: btn.HWND(),
			},
		}

		layout.Perform(center, nil)

		winRect := chkerr.Must(win.GetClientRect())
		btnRect := chkerr.Must(btn.GetWindowRect())
		win32util.ScreenToClient(win.HWND(), btnRect)
		winCenterX := winRect.Width() / 2
		btnCenterX := btnRect.Left + btnRect.Width()/2
		if abs(winCenterX-btnCenterX) > 1 {
			t.Errorf("Button is not centered horizontally: winCenterX=%v, btnCenterX=%v", winCenterX, btnCenterX)
		}
		winCenterY := winRect.Height() / 2
		btnCenterY := btnRect.Top + btnRect.Height()/2
		if abs(winCenterY-btnCenterY) > 1 {
			t.Errorf("Button is not centered vertically: winCenterY=%v, btnCenterY=%v", winCenterY, btnCenterY)
		}

		win.Close()

	}, nil)

}

func TestColumn(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Text:      "Test Column Layout",
			Style:     win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(600),
			Height:    metrics.Dip(500),
			OnDestroy: func() { app.Quit(0) },
		}))

		btn1 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Dip(300),
			Height: metrics.Dip(100),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		btn2 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		column := &layout.Column{
			Items: []layout.Layout{
				&layout.Window{
					Hwnd: btn1.HWND(),
				},
				&layout.Window{
					Hwnd: btn2.HWND(),
				},
			},
		}

		clientSize := chkerr.Must(layout.ClientSize(win.HWND()))
		layout.Perform(column, &clientSize)

		x1, y1, w1, h1 := chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 300 {
			t.Errorf(" w1=%v", w1)
		}
		if h1 != 100 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 := chkerr.Must4(btn2.Dimensions())
		if x2 != 0 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 100 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 200 {
			t.Errorf(" w2=%v", w2)
		}
		if h2 != 200 {
			t.Errorf("h2=%v", h2)
		}

		column.CrossAxisAlign = layout.AlignCenter

		clientSize = chkerr.Must(layout.ClientSize(win.HWND()))
		layout.Perform(column, &clientSize)

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 300 {
			t.Errorf(" w1=%v", w1)
		}
		if h1 != 100 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != 50 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 100 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 200 {
			t.Errorf(" w2=%v", w2)
		}
		if h2 != 200 {
			t.Errorf("h2=%v", h2)
		}

		column.CrossAxisAlign = layout.AlignEnd

		clientSize = chkerr.Must(layout.ClientSize(win.HWND()))
		layout.Perform(column, &clientSize)

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 300 {
			t.Errorf(" w1=%v", w1)
		}
		if h1 != 100 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != 100 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 100 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 200 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 200 {
			t.Errorf("h2=%v", h2)
		}

		win.Close()

	}, nil)

}

func TestRow(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Text:      "Test Column Layout",
			Style:     win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Dip(600),
			Height:    metrics.Dip(500),
			OnDestroy: func() { app.Quit(0) },
		}))

		btn1 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Dip(300),
			Height: metrics.Dip(100),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		btn2 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Dip(200),
			Height: metrics.Dip(200),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		row := &layout.Row{
			Items: []layout.Layout{
				&layout.Window{
					Hwnd: btn1.HWND(),
				},
				&layout.Window{
					Hwnd: btn2.HWND(),
				},
			},
		}

		clientSize := chkerr.Must(layout.ClientSize(win.HWND()))
		layout.Perform(row, &clientSize)

		x1, y1, w1, h1 := chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 300 {
			t.Errorf("w1=%v", w1)
		}
		if h1 != 100 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 := chkerr.Must4(btn2.Dimensions())
		if x2 != 300 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 0 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 200 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 200 {
			t.Errorf("h2=%v", h2)
		}

		row.CrossAxisAlign = layout.AlignCenter

		clientSize = chkerr.Must(layout.ClientSize(win.HWND()))
		layout.Perform(row, &clientSize)

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 50 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 300 {
			t.Errorf("w1=%v", w1)
		}
		if h1 != 100 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != 300 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 0 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 200 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 200 {
			t.Errorf("h2=%v", h2)
		}

		row.CrossAxisAlign = layout.AlignEnd

		clientSize = chkerr.Must(layout.ClientSize(win.HWND()))
		layout.Perform(row, &clientSize)

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 100 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 300 {
			t.Errorf("w1=%v", w1)
		}
		if h1 != 100 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != 300 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 0 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 200 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 200 {
			t.Errorf("h2=%v", h2)
		}

		row.MainAxisAlign = layout.AlignCenter

		clientSize = chkerr.Must(layout.ClientSize(win.HWND()))
		layout.Perform(row, &clientSize)

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if y1 != 100 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 300 {
			t.Errorf("w1=%v", w1)
		}
		if h1 != 100 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != x1+w1 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 0 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 200 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 200 {
			t.Errorf("h2=%v", h2)
		}

		win.Close()

	}, nil)

}
