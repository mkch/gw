package template_test

import (
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/dialog/template"
	"github.com/mkch/gw/paint"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"golang.org/x/sys/windows"
)

func TestNewTemplate(t *testing.T) {

	const (
		BUTTON_ID win32.DWORD = 1000 + iota
		EDIT_ID
		STATIC_ID
		LISTBOX_ID
		SCROLLBAR_ID
		COMBOBOX_ID
		CUSTOM_CONTROL_ID
	)

	gw.Run(func(app *app.App) {

		chkerr.Must(win32util.RegisterClass(&win32util.WndClass{
			ClassName: "CustomControl",
			WndProc: func(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
				switch msg {
				case win32.WM_PAINT:
					dc := chkerr.Must(paint.NewPaintDC(hwnd))
					defer dc.EndPaint()
					win32.FillRect(dc.HDC(), dc.Rect(), win32.HBRUSH(win32.COLOR_HIGHLIGHT+1))
					return 0
				}
				return win32.DefWindowProcW(hwnd, msg, wParam, lParam)
			},
		}))

		dlgTpl, err := template.New(&template.Dialog{
			Style: win32.WS_OVERLAPPEDWINDOW,
			CX:    200,
			CY:    300,
			Title: "Hello, world!",
		})
		if err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		dlgTpl.Add(
			&template.Control{
				Type:  template.Button,
				Style: win32.WS_CHILD | win32.WS_VISIBLE | win32.BS_DEFPUSHBUTTON,
				Title: "OK",
				X:     10,
				Y:     10,
				CX:    80,
				CY:    30,
				ID:    BUTTON_ID,
			},
			&template.Control{
				Type:  template.Edit,
				Style: win32.WS_CHILD | win32.WS_VISIBLE | win32.WS_BORDER,
				Title: "OK",
				X:     10,
				Y:     50,
				CX:    80,
				CY:    30,
				ID:    EDIT_ID,
			},
			&template.Control{
				Type:  template.Static,
				Style: win32.WS_CHILD | win32.WS_VISIBLE | win32.WS_BORDER,
				Title: "OK",
				X:     10,
				Y:     90,
				CX:    80,
				CY:    30,
				ID:    STATIC_ID,
			},
			&template.Control{
				Type:  template.ListBox,
				Style: win32.WS_CHILD | win32.WS_VISIBLE | win32.WS_BORDER,
				Title: "OK",
				X:     10,
				Y:     130,
				CX:    80,
				CY:    30,
				ID:    LISTBOX_ID,
			},
			&template.Control{
				Type:  template.ScrollBar,
				Style: win32.WS_CHILD | win32.WS_VISIBLE | win32.WS_BORDER,
				Title: "OK",
				X:     10,
				Y:     170,
				CX:    80,
				CY:    30,
				ID:    SCROLLBAR_ID,
			},
			&template.Control{
				Type:  template.ComboBox,
				Style: win32.WS_CHILD | win32.WS_VISIBLE,
				Title: "OK",
				X:     10,
				Y:     210,
				CX:    80,
				CY:    30,
				ID:    COMBOBOX_ID,
			},
			&template.Control{
				Type:  template.ControlType("CustomControl"),
				Style: win32.WS_CHILD | win32.WS_VISIBLE,
				X:     10,
				Y:     250,
				CX:    80,
				CY:    30,
				ID:    CUSTOM_CONTROL_ID,
			},
		)
		r, err := win32.DialogBoxIndirectParamW(0, dlgTpl.Data(), 0,
			windows.NewCallback(func(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
				switch msg {
				case win32.WM_INITDIALOG:
					var nameBuffer [256]win32.WCHAR
					win32.GetClassNameW(chkerr.Must(win32.GetDlgItem(hwnd, win32.INT(BUTTON_ID))), &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "Button" {
						t.Errorf("Unexpected class name for button control: got %q, want %q", className, "Button")
					}
					win32.GetClassNameW(chkerr.Must(win32.GetDlgItem(hwnd, win32.INT(EDIT_ID))), &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "Edit" {
						t.Errorf("Unexpected class name for edit control: got %q, want %q", className, "Edit")
					}
					win32.GetClassNameW(chkerr.Must(win32.GetDlgItem(hwnd, win32.INT(STATIC_ID))), &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "Static" {
						t.Errorf("Unexpected class name for static control: got %q, want %q", className, "Static")
					}
					win32.GetClassNameW(chkerr.Must(win32.GetDlgItem(hwnd, win32.INT(LISTBOX_ID))), &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "ListBox" {
						t.Errorf("Unexpected class name for listbox control: got %q, want %q", className, "ListBox")
					}
					win32.GetClassNameW(chkerr.Must(win32.GetDlgItem(hwnd, win32.INT(SCROLLBAR_ID))), &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "ScrollBar" {
						t.Errorf("Unexpected class name for scrollbar control: got %q, want %q", className, "ScrollBar")
					}
					win32.GetClassNameW(chkerr.Must(win32.GetDlgItem(hwnd, win32.INT(COMBOBOX_ID))), &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "ComboBox" {
						t.Errorf("Unexpected class name for combobox control: got %q, want %q", className, "ComboBox")
					}
					win32.GetClassNameW(chkerr.Must(win32.GetDlgItem(hwnd, win32.INT(CUSTOM_CONTROL_ID))), &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "CustomControl" {
						t.Errorf("Unexpected class name for custom control: got %q, want %q", className, "CustomControl")
					}
					win32.EndDialog(hwnd, 100)
					return 1
				case win32.WM_DESTROY:
					app.Quit(0)
					return 0
				case win32.WM_COMMAND:
					id := win32.LOWORD(wParam)
					if id == win32.IDOK || id == win32.IDCANCEL {
						win32.EndDialog(hwnd, 100)
						return 1
					}
				}
				return 0
			}), 0)
		if err != nil {
			t.Fatalf("DialogBoxIndirectParamW failed: %v", err)
		}
		if r != 100 {
			t.Fatalf("Unexpected dialog result: got %d, want 100", r)
		}

	}, func(app *app.App) {
		chkerr.MustOK(win32util.UnregisterClass("CustomControl", 0))
	})
}

func TestNewTemplateCustomClass(t *testing.T) {

	gw.Run(func(app *app.App) {

		chkerr.Must(win32util.RegisterClass(&win32util.WndClass{
			ClassName:  "CustomWindow",
			WndExtra:   win32.DLGWINDOWEXTRA,
			WndProc:    win32.DefDlgProcW,
			Background: win32.HBRUSH(win32.COLOR_WINDOW + 1),
			Cursor:     chkerr.Must(win32.LoadImageW_uintptr[win32.HCURSOR](0, uintptr(win32.OCR_NORMAL), win32.IMAGE_CURSOR, 0, 0, win32.LR_DEFAULTSIZE|win32.LR_SHARED)),
		}))

		dlgTpl, err := template.New(&template.Dialog{
			Class: "CustomWindow",
			Style: win32.WS_OVERLAPPEDWINDOW,
			CX:    200,
			CY:    300,
			Title: "Hello, world!",
		})
		if err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		const BUTTON_ID = 1000

		dlgTpl.Add(
			&template.Control{
				Type:  template.Button,
				Style: win32.WS_CHILD | win32.WS_VISIBLE | win32.BS_DEFPUSHBUTTON,
				Title: "OK",
				X:     10,
				Y:     10,
				CX:    80,
				CY:    30,
				ID:    BUTTON_ID,
			},
		)
		r, err := win32.DialogBoxIndirectParamW(0, dlgTpl.Data(), 0,
			windows.NewCallback(func(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
				switch msg {
				case win32.WM_INITDIALOG:
					var nameBuffer [256]win32.WCHAR
					win32.GetClassNameW(hwnd, &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "CustomWindow" {
						t.Errorf("Unexpected class name for dialog: got %q, want %q", className, "CustomWindow")
					}
					win32.GetClassNameW(chkerr.Must(win32.GetDlgItem(hwnd, win32.INT(BUTTON_ID))), &nameBuffer[0], len(nameBuffer))
					if className := win32util.GoString(&nameBuffer[0], -1); className != "Button" {
						t.Errorf("Unexpected class name for button control: got %q, want %q", className, "Button")
					}
					win32.EndDialog(hwnd, 100)
					return 1
				case win32.WM_DESTROY:
					app.Quit(0)
					return 0
				case win32.WM_COMMAND:
					id := win32.LOWORD(wParam)
					if id == win32.IDOK || id == win32.IDCANCEL {
						win32.EndDialog(hwnd, 100)
						return 1
					}
				}
				return 0
			}), 0)
		if err != nil {
			t.Fatalf("DialogBoxIndirectParamW failed: %v", err)
		}
		if r != 100 {
			t.Fatalf("Unexpected dialog result: got %d, want 100", r)
		}

	}, func(app *app.App) {
		chkerr.MustOK(win32util.UnregisterClass("CustomWindow", 0))
	})
}
