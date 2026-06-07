package layout_test

import (
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/layout"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/panel"
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
			Child: &layout.Intrinsic{
				Hwnd: btn.HWND(),
			},
		}

		tree := chkerr.Must(layout.Build(center))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

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

	}, func(app *app.App) { app.DestroyAllWindows() })

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
			Width:  metrics.Px(250),
			Height: metrics.Px(100),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		btn2 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Px(150),
			Height: metrics.Px(300),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		column := &layout.Column{
			Children: []layout.Widget{
				&layout.Intrinsic{
					Hwnd: btn1.HWND(),
				},
				&layout.Intrinsic{
					Hwnd: btn2.HWND(),
				},
			},
		}

		tree := chkerr.Must(layout.Build(column))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		x1, y1, w1, h1 := chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 250 {
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
		if w2 != 150 {
			t.Errorf(" w2=%v", w2)
		}
		if h2 != 300 {
			t.Errorf("h2=%v", h2)
		}

		column.CrossAxisAlign = layout.AlignCenter

		tree = chkerr.Must(layout.Build(column))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 250 {
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
		if w2 != 150 {
			t.Errorf(" w2=%v", w2)
		}
		if h2 != 300 {
			t.Errorf("h2=%v", h2)
		}

		column.CrossAxisAlign = layout.AlignEnd

		tree = chkerr.Must(layout.Build(column))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 250 {
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
		if w2 != 150 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 300 {
			t.Errorf("h2=%v", h2)
		}

		win.Close()

	}, func(app *app.App) { app.DestroyAllWindows() })

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
			Width:  metrics.Px(280),
			Height: metrics.Px(220),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		btn2 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button",
			Width:  metrics.Px(180),
			Height: metrics.Px(100),
			Style:  win32.WS_CHILD | win32.WS_VISIBLE,
		}))

		row := &layout.Row{
			Children: []layout.Widget{
				&layout.Intrinsic{
					Hwnd: btn1.HWND(),
				},
				&layout.Intrinsic{
					Hwnd: btn2.HWND(),
				},
			},
		}

		tree := chkerr.Must(layout.Build(row))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		x1, y1, w1, h1 := chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 280 {
			t.Errorf("w1=%v", w1)
		}
		if h1 != 220 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 := chkerr.Must4(btn2.Dimensions())
		if x2 != 280 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 0 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 180 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 100 {
			t.Errorf("h2=%v", h2)
		}

		row.CrossAxisAlign = layout.AlignCenter

		tree = chkerr.Must(layout.Build(row))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 280 {
			t.Errorf("w1=%v", w1)
		}
		if h1 != 220 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != 280 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 60 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 180 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 100 {
			t.Errorf("h2=%v", h2)
		}

		row.CrossAxisAlign = layout.AlignEnd

		tree = chkerr.Must(layout.Build(row))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("x1=%v", x1)
		}
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 280 {
			t.Errorf("w1=%v", w1)
		}
		if h1 != 220 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != 280 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 120 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 180 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 100 {
			t.Errorf("h2=%v", h2)
		}

		row.MainAxisAlign = layout.AlignCenter

		tree = chkerr.Must(layout.Build(row))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		x1, y1, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if y1 != 0 {
			t.Errorf("y1=%v", y1)
		}
		if w1 != 280 {
			t.Errorf("w1=%v", w1)
		}
		if h1 != 220 {
			t.Errorf("h1=%v", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != x1+w1 {
			t.Errorf("x2=%v", x2)
		}
		if y2 != 120 {
			t.Errorf("y2=%v", y2)
		}
		if w2 != 180 {
			t.Errorf("w2=%v", w2)
		}
		if h2 != 100 {
			t.Errorf("h2=%v", h2)
		}

		win.Close()

	}, func(app *app.App) { app.DestroyAllWindows() })

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
			Child: &layout.Intrinsic{
				Hwnd: btn.HWND(),
			},
		}

		tree := chkerr.Must(layout.Build(padding))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

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

	}, func(app *app.App) { app.DestroyAllWindows() })
}

func TestPaddingWithHwnd(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Text:  "Test Padding Layout",
			Style: win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
			X:     gw.CW_USEDEFAULT,
		}))
		win32util.SetClientSize(win.HWND(), 600, 500)

		pnl := chkerr.Must(panel.New(&panel.Spec{
			Parent: win,
		}))
		pnl.SetBackgroundColor(win32.COLORREF(0xFF0000))

		btn := chkerr.Must(button.New(&button.Spec{
			Parent: pnl,
			Style:  win32.WS_VISIBLE,
			Text:   "Button",
			Width:  metrics.Px(200),
			Height: metrics.Px(100),
		}))

		center := &layout.Center{
			Child: &layout.Padding{
				Hwnd: pnl.HWND(),
				Left: metrics.Dip(6),
				Top:  metrics.Dip(10),
				Child: &layout.Intrinsic{
					Hwnd: btn.HWND(),
				},
			},
		}

		tree := chkerr.Must(layout.Build(center))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		dpi := chkerr.Must(win32.GetDpiForWindow(win.HWND()))

		dip6 := win32.LONG(metrics.Dip(6).Px(dpi).Value())
		dip10 := win32.LONG(metrics.Dip(10).Px(dpi).Value())
		btnRect := chkerr.Must(btn.GetWindowRect())
		chkerr.MustOK(win32util.ScreenToClient(win.HWND(), btnRect))
		x, y, w, h := btnRect.Left, btnRect.Top, btnRect.Width(), btnRect.Height()
		if desiredX := (600-200-dip6)/2 + dip6; x != desiredX {
			t.Errorf("x=%v, want %v", x, desiredX)
		}
		if desiredY := (500-100-dip10)/2 + dip10; y != desiredY {
			t.Errorf("y=%v, want %v", y, desiredY)
		}
		if w != 200 {
			t.Errorf("w=%v, want %v", w, 200)
		}
		if h != 100 {
			t.Errorf("h=%v, want %v", h, 100)
		}

		app.Quit(0)

	}, func(app *app.App) { app.DestroyAllWindows() })
}

func TestSized(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Text:   "Test Sized Layout",
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

		sized := &layout.Sized{
			Child:  &layout.Intrinsic{Hwnd: btn.HWND()},
			Width:  metrics.Dip(50),
			Height: metrics.Dip(100),
		}

		tree := chkerr.Must(layout.Build(sized))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		dpi := chkerr.Must(win32.GetDpiForWindow(win.HWND()))

		dip50 := int(metrics.Dip(50).Px(dpi).Value())
		dip100 := int(metrics.Dip(100).Px(dpi).Value())
		x, y, w, h := chkerr.Must4(btn.Dimensions())
		if x != 0 {
			t.Errorf("x=%v, want %v", x, 0)
		}
		if y != 0 {
			t.Errorf("y=%v, want %v", y, 0)
		}
		if w != dip50 {
			t.Errorf("w=%v, want %v", w, dip50)
		}
		if h != dip100 {
			t.Errorf("h=%v, want %v", h, dip100)
		}

		app.Quit(0)
	}, func(app *app.App) { app.DestroyAllWindows() })
}

func TestSizedHwnd(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Text:   "Test Sized Layout",
			Style:  win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
			X:      gw.CW_USEDEFAULT,
			Width:  metrics.Px(600),
			Height: metrics.Px(500),
		}))

		btn := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Style:  win32.WS_VISIBLE,
			Text:   "Button",
		}))

		center := &layout.Center{
			Child: &layout.Sized{
				Hwnd:   btn.HWND(),
				Width:  metrics.Dip(50),
				Height: metrics.Dip(100),
			},
		}

		tree := chkerr.Must(layout.Build(center))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		dpi := chkerr.Must(win32.GetDpiForWindow(win.HWND()))

		dip50 := int(metrics.Dip(50).Px(dpi).Value())
		dip100 := int(metrics.Dip(100).Px(dpi).Value())
		_, _, w, h := chkerr.Must4(btn.Dimensions())
		if w != dip50 {
			t.Errorf("w=%v, want %v", w, dip50)
		}
		if h != dip100 {
			t.Errorf("h=%v, want %v", h, dip100)
		}

		app.Quit(0)

	}, func(app *app.App) { app.DestroyAllWindows() })
}

func TestExpanded(t *testing.T) {
	gw.Run(func(app *app.App) {
		win := chkerr.Must(window.New(&window.Spec{
			Text:      "Test Expanded Layout",
			Style:     win32.WS_OVERLAPPEDWINDOW | win32.WS_VISIBLE,
			X:         gw.CW_USEDEFAULT,
			Width:     metrics.Px(600),
			Height:    metrics.Px(500),
			OnDestroy: func() { app.Quit(0) },
		}))

		// btn1 has a fixed intrinsic size used as the non-Expanded sibling.
		btn1 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button1",
			Width:  metrics.Px(250),
			Height: metrics.Px(100),
		}))

		// btn2 and btn3 have no initial size; they will be sized by Expanded.
		btn2 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button2",
		}))

		btn3 := chkerr.Must(button.New(&button.Spec{
			Parent: win,
			Text:   "Button3",
		}))

		clientRect := chkerr.Must(win.GetClientRect())
		clientW := int(clientRect.Width())
		clientH := int(clientRect.Height())

		// Case 1: Column with one fixed Intrinsic child followed by one Expanded child.
		// The Expanded child must fill the remaining height and the full client width.
		column := &layout.Column{
			Children: []layout.Widget{
				&layout.Intrinsic{Hwnd: btn1.HWND()},
				&layout.Expanded{
					Flex:  1,
					Child: &layout.Intrinsic{Hwnd: btn2.HWND()},
				},
			},
		}
		tree := chkerr.Must(layout.Build(column))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		x1, y1, w1, h1 := chkerr.Must4(btn1.Dimensions())
		if x1 != 0 {
			t.Errorf("case1: x1=%v, want 0", x1)
		}
		if y1 != 0 {
			t.Errorf("case1: y1=%v, want 0", y1)
		}
		if w1 != 250 {
			t.Errorf("case1: w1=%v, want 250", w1)
		}
		if h1 != 100 {
			t.Errorf("case1: h1=%v, want 100", h1)
		}

		x2, y2, w2, h2 := chkerr.Must4(btn2.Dimensions())
		if x2 != 0 {
			t.Errorf("case1: x2=%v, want 0", x2)
		}
		if y2 != 100 {
			t.Errorf("case1: y2=%v, want 100", y2)
		}
		if w2 != clientW {
			t.Errorf("case1: w2=%v, want %v (clientW)", w2, clientW)
		}
		if h2 != clientH-100 {
			t.Errorf("case1: h2=%v, want %v (clientH-100)", h2, clientH-100)
		}

		// Case 2: Column with two Expanded children with Flex 1:2.
		// Their heights must sum to clientH, and the second child must be taller than the first.
		column2 := &layout.Column{
			Children: []layout.Widget{
				&layout.Expanded{
					Flex:  1,
					Child: &layout.Intrinsic{Hwnd: btn2.HWND()},
				},
				&layout.Expanded{
					Flex:  2,
					Child: &layout.Intrinsic{Hwnd: btn3.HWND()},
				},
			},
		}
		tree = chkerr.Must(layout.Build(column2))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		_, y2, _, h2 = chkerr.Must4(btn2.Dimensions())
		_, y3, _, h3 := chkerr.Must4(btn3.Dimensions())
		if y2 != 0 {
			t.Errorf("case2: y2=%v, want 0", y2)
		}
		if y3 != h2 {
			t.Errorf("case2: y3=%v, want %v (h2)", y3, h2)
		}
		if h2+h3 != clientH {
			t.Errorf("case2: h2+h3=%v, want %v (clientH)", h2+h3, clientH)
		}
		if h2 >= h3 {
			t.Errorf("case2: expected h2 < h3 for Flex 1:2, got h2=%v h3=%v", h2, h3)
		}

		// Case 3: Row with one fixed Intrinsic child followed by one Expanded child.
		// The Expanded child must fill the remaining width and the full client height.
		row := &layout.Row{
			Children: []layout.Widget{
				&layout.Intrinsic{Hwnd: btn1.HWND()},
				&layout.Expanded{
					Flex:  1,
					Child: &layout.Intrinsic{Hwnd: btn2.HWND()},
				},
			},
		}
		tree = chkerr.Must(layout.Build(row))
		chkerr.MustOK(layout.PerformWindow(tree, win.HWND()))

		_, _, w1, h1 = chkerr.Must4(btn1.Dimensions())
		if w1 != 250 {
			t.Errorf("case3: w1=%v, want 250", w1)
		}
		if h1 != 100 {
			t.Errorf("case3: h1=%v, want 100", h1)
		}

		x2, y2, w2, h2 = chkerr.Must4(btn2.Dimensions())
		if x2 != 250 {
			t.Errorf("case3: x2=%v, want 250", x2)
		}
		if y2 != 0 {
			t.Errorf("case3: y2=%v, want 0", y2)
		}
		if w2 != clientW-250 {
			t.Errorf("case3: w2=%v, want %v (clientW-250)", w2, clientW-250)
		}
		if h2 != clientH {
			t.Errorf("case3: h2=%v, want %v (clientH)", h2, clientH)
		}

		win.Close()
	}, func(app *app.App) { app.DestroyAllWindows() })
}
