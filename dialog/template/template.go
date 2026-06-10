package template

import (
	"encoding/binary"
	"errors"
	"unsafe"

	"github.com/mkch/gw/dialog"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
)

// https://learn.microsoft.com/en-us/windows/win32/dlgbox/dlgtemplateex

// Template represents a dialog box template in memory.
// It can be constructed using [New] and modified using [Append].
type Template struct {
	buf []byte
}

// Dialog represents the parameters for creating a dialog box template.
type Dialog struct {
	HelpID       win32.DWORD
	ExStyle      win32.WindowExStyle
	Style        win32.WindowStyle
	X, Y, CX, CY int16  // In dialog box units.
	Class        string // Window class of the dialog box. If empty, the predefined dialog box class is used.
	Title        string
}

// ControlType specifies the type of a control in a dialog box template.
type ControlType string

const (
	Button    ControlType = "Button"
	Edit      ControlType = "Edit"
	Static    ControlType = "Static"
	ListBox   ControlType = "ListBox"
	ScrollBar ControlType = "ScrollBar"
	ComboBox  ControlType = "ComboBox"
)

// predefined returns the predefined control type ordinal for the given ControlType, or 0 if it is not a predefined type.
func (c ControlType) predefined() uint16 {
	switch c {
	case Button:
		return _Button
	case Edit:
		return _Edit
	case Static:
		return _Static
	case ListBox:
		return _ListBox
	case ScrollBar:
		return _ScrollBar
	case ComboBox:
		return _ComboBox
	default:
		return 0
	}
}

const (
	_Button = 0x0080 + iota
	_Edit
	_Static
	_ListBox
	_ScrollBar
	_ComboBox
)

// Control represents a control to be added to a dialog box template.
type Control struct {
	HelpID       win32.DWORD
	ExStyle      win32.WindowExStyle
	Style        win32.WindowStyle
	X, Y, CX, CY int16 // In dialog box units.
	ID           win32.DWORD
	Type         ControlType
	Title        string
}

var endianness = binary.LittleEndian

// Add adds one or more controls to the dialog box template.
// Like the builtin append function, it returns the modified template.
func (t *Template) Add(controls ...*Control) {
	for _, ctrl := range controls {
		t.addControl(ctrl)
	}
}

// addControl add a control to the dialog box template.
func (t *Template) addControl(control *Control) {
	// Each DLGITEMTEMPLATEEX structure must be aligned on a DWORD boundary
	if r := len(t.buf) % 4; r != 0 {
		t.buf = append(t.buf, make([]byte, 4-r)...)
	}
	t.buf = endianness.AppendUint32(t.buf, uint32(control.HelpID))
	t.buf = endianness.AppendUint32(t.buf, uint32(control.ExStyle))
	t.buf = endianness.AppendUint32(t.buf, uint32(control.Style))
	t.buf = endianness.AppendUint16(t.buf, uint16(control.X))
	t.buf = endianness.AppendUint16(t.buf, uint16(control.Y))
	t.buf = endianness.AppendUint16(t.buf, uint16(control.CX))
	t.buf = endianness.AppendUint16(t.buf, uint16(control.CY))
	t.buf = endianness.AppendUint32(t.buf, uint32(control.ID))
	// Type
	predefined := control.Type.predefined()
	if predefined != 0 {
		t.buf = endianness.AppendUint16(t.buf, 0xFFFF)
		t.buf = endianness.AppendUint16(t.buf, uint16(predefined))
	} else {
		var buf []win32.WCHAR
		win32util.CString(string(control.Type), &buf)
		for _, r := range buf {
			t.buf = endianness.AppendUint16(t.buf, uint16(r))
		}
	}
	// Title
	if control.Title != "" {
		var strBuf []win32.WCHAR
		win32util.CString(control.Title, &strBuf)
		for _, r := range strBuf {
			t.buf = endianness.AppendUint16(t.buf, uint16(r))
		}
	} else {
		t.buf = endianness.AppendUint16(t.buf, 0) // no title
	}
	// extraCount
	t.buf = append(t.buf, 0, 0)

	// Increase the item count in the template header.
	endianness.PutUint16(t.buf[16:], endianness.Uint16(t.buf[16:])+1)
}

// Data returns a pointer to the dialog box template data that can be passed to Win32 API functions.
func (t Template) Data() *win32.DLGTEMPLATE {
	return (*win32.DLGTEMPLATE)(unsafe.Pointer(&win32.AlignSlice(t.buf)[0]))
}

// New creates a new dialog box template from the given Dialog specification.
func New(spec *Dialog) (*Template, error) {
	if spec == nil {
		spec = &Dialog{}
	}
	if spec.Style&(dialog.DS_SETFONT|dialog.DS_SHELLFONT) != 0 {
		return nil, errors.New("DS_SETFONT and DS_SHELLFONT styles are not supported")
	}
	var buf []byte
	buf = endianness.AppendUint16(buf, 1)      // version
	buf = endianness.AppendUint16(buf, 0xFFFF) // signature
	buf = endianness.AppendUint32(buf, 0)      // help ID
	buf = endianness.AppendUint32(buf, uint32(spec.ExStyle))
	buf = endianness.AppendUint32(buf, uint32(spec.Style))
	buf = endianness.AppendUint16(buf, 0) // number of items
	buf = endianness.AppendUint16(buf, uint16(spec.X))
	buf = endianness.AppendUint16(buf, uint16(spec.Y))
	buf = endianness.AppendUint16(buf, uint16(spec.CX))
	buf = endianness.AppendUint16(buf, uint16(spec.CY))
	buf = endianness.AppendUint16(buf, 0) // menu
	var strBuf []win32.WCHAR
	if spec.Class != "" {
		win32util.CString(spec.Class, &strBuf)
		for _, r := range strBuf {
			buf = endianness.AppendUint16(buf, uint16(r)) // window class
		}
	} else {
		buf = endianness.AppendUint16(buf, 0) // predefined dialog box class
	}
	if spec.Title != "" {
		win32util.CString(spec.Title, &strBuf)
		for _, r := range strBuf {
			buf = endianness.AppendUint16(buf, uint16(r)) // title
		}
		if len(buf)%2 != 0 {
			buf = append(buf, 0) // padding: aligned on WORD boundaries.
		}
	} else {
		buf = endianness.AppendUint16(buf, 0) // no title
	}

	return &Template{buf}, nil
}
