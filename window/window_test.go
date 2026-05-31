package window_test

import (
	"testing"

	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

type MyWindow struct {
	window.Window
}

func (w *MyWindow) OnDestroy() {
	w.App().Quit(1)
	w.Window.OnDestroy()
}

func TestSimple(t *testing.T) {
	gw.Run(func(app *app.App) {
		w := window.New(&window.Spec{
			Text:      "Hello, World!",
			Style:     win32.WS_OVERLAPPEDWINDOW,
			OnDestroy: func() { app.Quit(0) },
		})

		hwnd := w.HWND()
		if hwnd == 0 {
			t.Fatal("Failed to create window")
		}

		if myWin, ok := gw.LookupWindow(hwnd).(*window.Window); !ok || myWin != w {
			t.Fatal("LookupWindow did not return the correct window instance")
		}

		text, err := w.Text()
		if err != nil {
			t.Fatalf("Failed to get window text: %v", err)
		}
		if text != "Hello, World!" {
			t.Fatalf("Expected window text to be 'Hello, World!', got '%s'", text)
		}

		style, err := win32.GetWindowLongPtrW(w.HWND(), win32.GWL_STYLE)
		if err != nil {
			t.Fatalf("Failed to get window style: %v", err)
		}
		if win32.WINDOW_STYLE(style)&win32.WS_OVERLAPPEDWINDOW == 0 {
			t.Fatalf("Expected window style to include WS_OVERLAPPEDWINDOW, got %x", style)
		}

		w.Close()

		var panicked bool
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			w.Text() // should panic because the window is closed
		}()
		if !panicked {
			t.Fatal("Expected panic when calling Text() on a closed window, but did not panic")
		}

	}, nil)
}

func TestWrapper(t *testing.T) {
	gw.Run(func(app *app.App) {
		w := &MyWindow{
			window.Window{
				Spec: &window.Spec{
					Text:      "Hello, World!",
					Style:     win32.WS_OVERLAPPEDWINDOW,
					OnDestroy: func() { app.Quit(0) },
				},
			},
		}
		gw.Init(w)

		hwnd := w.HWND()
		if hwnd == 0 {
			t.Fatal("Failed to create window")
		}

		if myWin, ok := gw.LookupWindow(hwnd).(*MyWindow); !ok || myWin != w {
			t.Fatal("LookupWindow did not return the correct window instance")
		}

		text, err := w.Text()
		if err != nil {
			t.Fatalf("Failed to get window text: %v", err)
		}
		if text != "Hello, World!" {
			t.Fatalf("Expected window text to be 'Hello, World!', got '%s'", text)
		}

		w.Close()

	}, nil)
}
