package edit

import (
	"github.com/mkch/gg"
	"github.com/mkch/gg/errortrace/chkerr"
	"github.com/mkch/gw"
	"github.com/mkch/gw/control"
	"github.com/mkch/gw/metrics"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

const (
	ES_LEFT        win32.WINDOW_STYLE = 0x0000
	ES_CENTER      win32.WINDOW_STYLE = 0x0001
	ES_RIGHT       win32.WINDOW_STYLE = 0x0002
	ES_MULTILINE   win32.WINDOW_STYLE = 0x0004
	ES_UPPERCASE   win32.WINDOW_STYLE = 0x0008
	ES_LOWERCASE   win32.WINDOW_STYLE = 0x0010
	ES_PASSWORD    win32.WINDOW_STYLE = 0x0020
	ES_AUTOVSCROLL win32.WINDOW_STYLE = 0x0040
	ES_AUTOHSCROLL win32.WINDOW_STYLE = 0x0080
	ES_NOHIDESEL   win32.WINDOW_STYLE = 0x0100
	ES_OEMCONVERT  win32.WINDOW_STYLE = 0x0400
	ES_READONLY    win32.WINDOW_STYLE = 0x0800
	ES_WANTRETURN  win32.WINDOW_STYLE = 0x1000
	ES_NUMBER      win32.WINDOW_STYLE = 0x2000
)

const (
	EM_GETSEL              win32.UINT = 0x00B0
	EM_SETSEL              win32.UINT = 0x00B1
	EM_GETRECT             win32.UINT = 0x00B2
	EM_SETRECT             win32.UINT = 0x00B3
	EM_SETRECTNP           win32.UINT = 0x00B4
	EM_SCROLL              win32.UINT = 0x00B5
	EM_LINESCROLL          win32.UINT = 0x00B6
	EM_SCROLLCARET         win32.UINT = 0x00B7
	EM_GETMODIFY           win32.UINT = 0x00B8
	EM_SETMODIFY           win32.UINT = 0x00B9
	EM_GETLINECOUNT        win32.UINT = 0x00BA
	EM_LINEINDEX           win32.UINT = 0x00BB
	EM_SETHANDLE           win32.UINT = 0x00BC
	EM_GETHANDLE           win32.UINT = 0x00BD
	EM_GETTHUMB            win32.UINT = 0x00BE
	EM_LINELENGTH          win32.UINT = 0x00C1
	EM_REPLACESEL          win32.UINT = 0x00C2
	EM_GETLINE             win32.UINT = 0x00C4
	EM_LIMITTEXT           win32.UINT = 0x00C5
	EM_CANUNDO             win32.UINT = 0x00C6
	EM_UNDO                win32.UINT = 0x00C7
	EM_FMTLINES            win32.UINT = 0x00C8
	EM_LINEFROMCHAR        win32.UINT = 0x00C9
	EM_SETTABSTOPS         win32.UINT = 0x00CB
	EM_SETPASSWORDCHAR     win32.UINT = 0x00CC
	EM_EMPTYUNDOBUFFER     win32.UINT = 0x00CD
	EM_GETFIRSTVISIBLELINE win32.UINT = 0x00CE
	EM_SETREADONLY         win32.UINT = 0x00CF
	EM_SETWORDBREAKPROC    win32.UINT = 0x00D0
	EM_GETWORDBREAKPROC    win32.UINT = 0x00D1
	EM_GETPASSWORDCHAR     win32.UINT = 0x00D2
	EM_SETMARGINS          win32.UINT = 0x00D3
	EM_GETMARGINS          win32.UINT = 0x00D4
	EM_SETLIMITTEXT        win32.UINT = EM_LIMITTEXT
	EM_GETLIMITTEXT        win32.UINT = 0x00D5
	EM_POSFROMCHAR         win32.UINT = 0x00D6
	EM_CHARFROMPOS         win32.UINT = 0x00D7
	EM_SETIMESTATUS        win32.UINT = 0x00D8
	EM_GETIMESTATUS        win32.UINT = 0x00D9
	EM_ENABLEFEATURE       win32.UINT = 0x00DA
)

type Edit struct {
	control.Control
	// Spec is used to create the window and is cleared after creation.
	Spec *Spec
}

type Spec struct {
	Parent  gw.WindowParent
	Text    string
	X       metrics.Dimension
	Y       metrics.Dimension
	Width   metrics.Dimension
	Height  metrics.Dimension
	Style   win32.WINDOW_STYLE
	ExStyle win32.WINDOW_EX_STYLE
}

func (e *Edit) OnInit() {
	defer func() { e.Spec = nil }()
	e.Control.OnInit()
}

func (e *Edit) CreateHandle() win32.HWND {
	dpi := gg.Must(win32.GetDpiForWindow(e.Spec.Parent.HWND()))
	return chkerr.Must(win32util.CreateWindow(&win32util.Wnd{
		ClassName:  "EDIT",
		WndParent:  e.Spec.Parent.HWND(),
		WindowName: e.Spec.Text,
		X:          e.Spec.X.Px(dpi),
		Y:          e.Spec.Y.Px(dpi),
		Width:      e.Spec.Width.Px(dpi),
		Height:     e.Spec.Height.Px(dpi),
		Style:      e.Spec.Style | win32.WS_CHILD,
		ExStyle:    e.Spec.ExStyle,
	}))
}

// New creates a new Edit control with the specified specification.
func New(spec *Spec) (edit *Edit) {
	edit = &Edit{Spec: spec}
	gw.Init(edit)
	return
}
