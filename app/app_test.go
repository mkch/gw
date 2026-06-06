package app_test

import (
	"testing"

	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/app"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/window"
)

func Test_App_DestroyAllWindows(t *testing.T) {
	gw.Run(func(app *app.App) {
		var win1Destroyed bool
		var win2Destroyed bool
		var winMsgDestroyed bool
		chkerr.Must(window.New(&window.Spec{
			OnDestroy: func() { win1Destroyed = true },
		}))
		chkerr.Must(window.New(&window.Spec{
			OnDestroy: func() { win2Destroyed = true },
		}))
		chkerr.Must(window.New(&window.Spec{
			Parent:    win32.HWND_MESSAGE,
			OnDestroy: func() { winMsgDestroyed = true },
		}))

		n := app.DestroyAllWindows()
		if n != 3 {
			t.Errorf("expected 3 windows destroyed, got %d", n)
		}
		if !win1Destroyed {
			t.Error("win1 was not destroyed")
		}
		if !win2Destroyed {
			t.Error("win2 was not destroyed")
		}
		if !winMsgDestroyed {
			t.Error("message-only window was not destroyed")
		}

		app.Quit(0)
	}, nil)
}
