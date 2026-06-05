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
			Width:     metrics.Px(500),
			Height:    metrics.Px(300),
			OnDestroy: func() { app.Quit(0) },
		}))

		btn := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Px(300),
			Height: metrics.Px(100),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		center := &layout.Center{
			Item: &layout.Intrinsic{
				Hwnd: btn.HWND(),
			},
		}

		chkerr.MustOK(layout.PerformWindow(center, win.HWND()))

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
			Width:     metrics.Px(600),
			Height:    metrics.Px(500),
			OnDestroy: func() { app.Quit(0) },
		}))

		btn1 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Px(300),
			Height: metrics.Px(100),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		btn2 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Px(200),
			Height: metrics.Px(200),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		column := &layout.Column{
			Items: []layout.Layout{
				&layout.Intrinsic{
					Hwnd: btn1.HWND(),
				},
				&layout.Intrinsic{
					Hwnd: btn2.HWND(),
				},
			},
		}
		chkerr.MustOK(layout.PerformWindow(column, win.HWND()))

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

		chkerr.MustOK(layout.PerformWindow(column, win.HWND()))

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

		chkerr.MustOK(layout.PerformWindow(column, win.HWND()))

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
			Width:     metrics.Px(600),
			Height:    metrics.Px(500),
			OnDestroy: func() { app.Quit(0) },
		}))

		btn1 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Px(300),
			Height: metrics.Px(100),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		btn2 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Px(200),
			Height: metrics.Px(200),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		row := &layout.Row{
			Items: []layout.Layout{
				&layout.Intrinsic{
					Hwnd: btn1.HWND(),
				},
				&layout.Intrinsic{
					Hwnd: btn2.HWND(),
				},
			},
		}

		chkerr.MustOK(layout.PerformWindow(row, win.HWND()))

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

		chkerr.MustOK(layout.PerformWindow(row, win.HWND()))

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

		chkerr.MustOK(layout.PerformWindow(row, win.HWND()))

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

		chkerr.MustOK(layout.PerformWindow(row, win.HWND()))

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

func TestPadding(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Text:   "Test Padding Layout",
			Style:  win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
			X:      gw.CW_USEDEFAULT,
			Width:  metrics.Px(600),
			Height: metrics.Px(500),
		}))

		btn := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE,
			Text:   "Button",
			Width:  metrics.Px(300),
			Height: metrics.Px(100),
		}))

		padding := &layout.Padding{
			Left: metrics.Dip(5),
			Top:  metrics.Dip(10),
			Item: &layout.Intrinsic{
				Hwnd: btn.HWND(),
			},
		}

		chkerr.MustOK(layout.PerformWindow(padding, win.HWND()))

		dpi := chkerr.Must(win32.GetDpiForWindow(win.HWND()))

		dip5 := int(metrics.Dip(5).Px(dpi).Value())
		dip10 := int(metrics.Dip(10).Px(dpi).Value())
		x, y, w, h := chkerr.Must4(btn.Dimensions())
		if x != dip5 {
			t.Errorf("x=%v, want %v", x, dip5)
		}
		if y != dip10 {
			t.Errorf("y=%v, want %v", y, dip10)
		}
		if w != 300 {
			t.Errorf("w=%v, want %v", w, 300)
		}
		if h != 100 {
			t.Errorf("h=%v, want %v", h, 100)
		}

		app.Quit(0)

	}, nil)
}
