package comctl32_v6_test

import (
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/combobox"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

//go:generate rsrc -arch amd64 -manifest manifest.xml
//go:generate rsrc -arch 386 -manifest manifest.xml

func TestCueBanner(t *testing.T) {
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

		cue, err := combo.CueBanner()
		if err != nil {
			t.Fatalf("failed to get cue banner: %v", err)
		}
		if cue != "" {
			t.Fatalf("expected empty cue banner, got %q", cue)
		}

		err = combo.SetCueBanner("Enter text here")
		if err != nil {
			t.Fatalf("failed to set cue banner: %v", err)
		}
		cue, err = combo.CueBanner()
		if err != nil {
			t.Fatalf("failed to get cue banner: %v", err)
		}
		if cue != "Enter text here" {
			t.Fatalf("expected cue banner to be %q, got %q", "Enter text here", cue)
		}

		win.Close()

	}, nil)
}

func TestMiniVisible(t *testing.T) {
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

		err = combo.SetMinVisible(5)
		if err != nil {
			t.Fatalf("failed to set minimum visible items: %v", err)
		}

		n, err := combo.MinVisible()
		if err != nil {
			t.Fatalf("failed to get minimum visible items: %v", err)
		}
		if n != 5 {
			t.Fatalf("expected minimum visible items to be 5, got %d", n)
		}

		win.Close()

	}, nil)
}
