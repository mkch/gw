package dialog

import (
	"errors"
	"math"

	"github.com/mkch/gg"
	"github.com/mkch/gw/button"
	"github.com/mkch/gw/internal"
	"github.com/mkch/gw/internal/objectmap"
	"github.com/mkch/gw/menu"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
	"golang.org/x/sys/windows"
)

// HWND -> Dialog
var dialogMap = make(map[win32.HWND]*Dialog)

// Dialog return code -> return value
var retMap = objectmap.New[any](1, math.MaxUint)

// LPARAM of WM_INITDIALOG
var dialogParamMap = objectmap.New[*Dialog](1, math.MaxUint)

type Dialog struct {
	window.Window
	OnCreate    func(dlg *Dialog)
	OnOK        func(dlg *Dialog) bool
	OnCancel    func(dlg *Dialog) bool
	initSpec    *Spec
	dlgProc     DlgProc
	prevDlgProc Proc
}

type Spec struct {
	Text      string
	Style     win32.DWORD
	ExStyle   win32.DWORD
	X         int
	Y         int
	Width     int
	Height    int
	WndParent win32.HWND
	Menu      *menu.Menu
	Instance  win32.HINSTANCE // 0 for this module.
	OnCreate  func(dlg *Dialog)
	OnOK      func(dlg *Dialog) bool
	OnCancel  func(dlg *Dialog) bool
}

var dlgProc = windows.NewCallback(func(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.UINT_PTR {
	switch msg {
	case win32.WM_INITDIALOG:
		// Find the *Dialog set in lParam
		dialog, _ := dialogParamMap.Value(objectmap.Handle(lParam))
		if err := window.Attach(hwnd, &dialog.WindowBase); err != nil {
			panic(err)
		}
		// Put dialog in dialogMap for following messages to retrieve.
		dialogMap[hwnd] = dialog
		dialog.SetText(dialog.initSpec.Text)
		dialog.SetMenu(dialog.initSpec.Menu)
		return dialog.callDlgProc(hwnd, msg, wParam, lParam)
	case win32.WM_NCDESTROY:
		r := dialogMap[hwnd].callDlgProc(hwnd, msg, wParam, lParam)
		delete(dialogMap, hwnd)
		return r
	default:
		return dialogMap[hwnd].callDlgProc(hwnd, msg, wParam, lParam)
	}
})

type Proc func(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) bool
type DlgProc func(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevDlgProc Proc) bool

func (d *Dialog) callDlgProc(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.UINT_PTR {
	return gg.If[win32.UINT_PTR](d.dlgProc(hwnd, msg, wParam, lParam, d.prevDlgProc), 1, 0)
}

func (d *Dialog) SetDlgProc(dlgProc DlgProc) {
	if dlgProc == nil {
		panic(errors.New("nil DlgProc"))
	}
	oldProc, oldPrevProc := d.dlgProc, d.prevDlgProc
	d.prevDlgProc = func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) bool {
		return oldProc(hwnd, message, wParam, lParam, oldPrevProc)
	}
	d.dlgProc = dlgProc
}

// End ends the dialog and set the result value.
func (d *Dialog) End(result any) error {
	return win32.EndDialog(d.HWND(), win32.INT_PTR(retMap.Add(result)))
}

// Reposition repositions a top-level dialog box so that it fits within the desktop area.
func (d *Dialog) Reposition() error {
	_, err := win32.SendMessageW(d.HWND(), win32.DM_REPOSITION, 0, 0)
	return err
}

const (
	defaultButtonID = internal.MinMenuItemID - 1 - iota
)

var ErrNotChild = errors.New("no a child")

// SetDefault sets the default button of this dialog.
// Panic if btn is not a child of d.
func (d *Dialog) SetDefault(btn *button.Button) error {
	if parent, err := win32.GetAncestor(btn.HWND(), win32.GA_PARENT); err != nil {
		return err
	} else if parent != d.HWND() {
		return ErrNotChild
	}

	if oldDef, _ := win32.GetDlgItem(d.HWND(), defaultButtonID); oldDef != 0 {
		if _, err := win32.SetWindowLongPtrW(oldDef, win32.GWLP_ID, 0); err != nil {
			return err
		}
		if oldStyle, err := win32.GetWindowLongPtrW(oldDef, win32.GWL_STYLE); err != nil {
			return err
		} else {
			// Use BM_SETSTYLE instead of SetWindowLongPtr(GWL_STYLE) because it has a redraw option.
			win32.SendMessageW(oldDef, win32.BM_SETSTYLE, win32.WPARAM(oldStyle & ^win32.BS_DEFPUSHBUTTON), 1)
		}
	}

	var id win32.LONG_PTR
	var err error
	if id, err = win32.GetWindowLongPtrW(btn.HWND(), win32.GWLP_ID); err != nil {
		return err
	} else if id == 0 { // btn is a OK,Cancel button etc.
		if _, err := win32.SetWindowLongPtrW(btn.HWND(), win32.GWLP_ID, defaultButtonID); err != nil {
			return err
		}
		id = defaultButtonID
	}

	if oldStyle, err := win32.GetWindowLongPtrW(btn.HWND(), win32.GWL_STYLE); err != nil {
		return err
	} else {
		win32.SendMessageW(btn.HWND(), win32.BM_SETSTYLE, win32.WPARAM(oldStyle|win32.BS_DEFPUSHBUTTON), 1)
	}

	win32.SendMessageW(d.HWND(), win32.DM_SETDEFID, win32.WPARAM(id), 0)

	return nil
}

func (d *Dialog) setButtonID(btn *button.Button, id int) error {
	if parent, err := win32.GetAncestor(btn.HWND(), win32.GA_PARENT); err != nil {
		return err
	} else if parent != d.HWND() {
		return ErrNotChild
	}

	if oldOK, _ := win32.GetDlgItem(d.HWND(), win32.INT(id)); oldOK != 0 {
		if _, err := win32.SetWindowLongPtrW(oldOK, win32.GWLP_ID, 0); err != nil {
			return err
		}
	}

	if _, err := win32.SetWindowLongPtrW(btn.HWND(), win32.GWLP_ID, win32.LONG_PTR(id)); err != nil {
		return err
	}

	return nil
}

// SetOK sets the OK button of this dialog.
// Panic if btn is not a child of d.
func (d *Dialog) SetOK(btn *button.Button) error {
	return d.setButtonID(btn, win32.IDOK)
}

// SetCancel sets the Cancel button of this dialog.
// Panic if btn is not a child of d.
func (d *Dialog) SetCancel(btn *button.Button) error {
	return d.setButtonID(btn, win32.IDCANCEL)
}

// Modal shows a modal dialog box.
// The ret is the return value set by Dialog.End(), nil if Dialog.End() is not called.
func Modal(spec *Spec) (ret any, err error) {
	tpl := win32util.EmptyDialogTemplate(spec.Style, spec.ExStyle, win32.SHORT(spec.X), win32.SHORT(spec.Y), win32.SHORT(spec.Width), win32.SHORT(spec.Height))
	instance := spec.Instance
	if instance == 0 {
		instance, _ = win32.GetModuleHandleW[win32.HINSTANCE](nil)
	}
	dialog := &Dialog{
		OnCreate: spec.OnCreate,
		OnOK:     spec.OnOK,
		OnCancel: spec.OnCancel,
		initSpec: spec,
	}
	dialog.prevDlgProc = func(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) bool {
		switch msg {
		case win32.WM_INITDIALOG:
			if dialog.OnCreate != nil {
				dialog.OnCreate(dialog)
			}
		case win32.WM_COMMAND:
			if win32.HIWORD(wParam) == win32.BN_CLICKED {
				switch win32.LOWORD(wParam) {
				case win32.IDOK:
					if dialog.OnOK == nil || !dialog.OnOK(dialog) {
						win32.EndDialog(hwnd, 0)
					}
				case win32.IDCANCEL:
					if dialog.OnCancel == nil || !dialog.OnCancel(dialog) {
						win32.EndDialog(hwnd, 0)
					}
				}
			}
		}
		return false
	}
	dialog.dlgProc = func(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM, prevDlgProc Proc) bool {
		return prevDlgProc(hwnd, msg, wParam, lParam)
	}
	param := dialogParamMap.Add(dialog)
	defer dialogParamMap.Remove(param)
	r, err := win32.DialogBoxIndirectParamW(instance, tpl, spec.WndParent, dlgProc, win32.LPARAM(param))
	if err != nil {
		return nil, err
	}
	ret, _ = retMap.Value(objectmap.Handle(r))
	retMap.Remove(objectmap.Handle(r))
	return
}
