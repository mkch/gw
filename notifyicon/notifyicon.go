package notifyicon

import (
	"unsafe"

	"github.com/mkch/gw/internal/appmsg"
	"github.com/mkch/gw/notifyicon/sys"
	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/win32util"
	"github.com/mkch/gw/window"
)

// CallbackMessage is the message sent to the window when mouse or keyboard interaction
// occurs in the bounding rectangle of the icon.
// Process this message manually or use [Spec].OnEvent.
const CallbackMessage = appmsg.NOTIFY_ICON_CALLBACK

// Event is the interaction event occurs in the bounding rectangle of the icon.
type Event win32.WORD

const (
	NIN_SELECT           Event = (win32.WM_USER + 0) // Notify icon is selected by mouse or keyboard.
	NINF_KEY             Event = 0x1
	NIN_KEYSELECT        Event = (NIN_SELECT | NINF_KEY) // Notify icon is selected by or keyboard.
	NIN_BALLOONSHOW      Event = (win32.WM_USER + 2)
	NIN_BALLOONHIDE      Event = (win32.WM_USER + 3)
	NIN_BALLOONTIMEOUT   Event = (win32.WM_USER + 4)
	NIN_BALLOONUSERCLICK Event = (win32.WM_USER + 5)
	NIN_POPUPOPEN        Event = (win32.WM_USER + 6) // Tooltip is needed.
	NIN_POPUPCLOSE       Event = (win32.WM_USER + 7)
	MIN_CONTEXTMENU      Event = win32.WM_CONTEXTMENU // Action to trigger the shortcut menu(i.e. right click).
)

// ParseCallback parses the parameters of CallbackMessage.
// id is the identifier of notify icon. event is the interaction event occurred. eventX/Y is the event coordinates.
func ParseCallback(wParam win32.WPARAM, lParam win32.LPARAM) (id win32.WORD, event Event, eventX win32.SHORT, eventY win32.SHORT) {
	event = Event(win32.LOWORD(uintptr(lParam)))
	id = win32.HIWORD(uintptr(lParam))
	if id == 0 && wParam > 0 && wParam <= 0xFFFF { // HACK: possible Shell_NotifyIconW BUG?
		id = win32.LOWORD(wParam)
	}
	if event == NIN_POPUPOPEN || event == NIN_SELECT || event == NIN_KEYSELECT ||
		(event >= win32.WM_MOUSEFIRST && event <= win32.WM_MOUSELAST) {
		eventX = win32.SHORT(win32.LOWORD(wParam))
		eventY = win32.SHORT(win32.HIWORD(wParam))
	}
	return
}

const (
	BI_NONE    win32.HICON = win32.HICON(sys.NIIF_NONE)
	BI_INFO    win32.HICON = win32.HICON(sys.NIIF_INFO)
	BI_WARNING win32.HICON = win32.HICON(sys.NIIF_WARNING)
	BI_ERROR   win32.HICON = win32.HICON(sys.NIIF_ERROR)
)

type NotifySpec struct {
	// Icon of notification. BI_xxx or a real HICON
	// If Title is empty string, the icon is not shown.
	Icon win32.HICON
	// Title of notification.
	Title string
	// Message of notification.
	Message string
	// Do not play the notifications sound.
	NoSound bool
	// Use (SM_CXICON x SM_CYICON) instead of (SM_CXSMICON x SM_CYSMICON) icon dimension.
	LargeIcon bool
	// If the notification cannot be displayed immediately, discard it.
	RealTime bool
	// Do not display the balloon notification if the current user is in "quiet time"
	RespectQuietTime bool
}

type Spec struct {
	// Icon in system tray.
	Icon win32.HICON
	// Tool tip string. if Tip is "", [NIN_POPUPOPEN] event is sent to allow showing application-drawn tooltip.
	Tip string
	// The ballon notification.
	Notify *NotifySpec
	// Optional callback to receive events.
	OnEvent func(id win32.WORD, event Event, eventX win32.SHORT, eventY win32.SHORT)
}

// Modifier is used to modify notify icon.
// A typical usage is a
//
// modifier.SetXXX(...).SetXXX(...).Apply()
//
// See [StartModify].
type Modifier struct {
	data sys.NOTIFYICONDATAW
}

// Apply applies the change.
func (m *Modifier) Apply() error {
	return sys.Shell_NotifyIconW(sys.NIM_MODIFY, &m.data)
}

func (m *Modifier) SetIcon(icon win32.HICON) *Modifier {
	m.data.Flags |= sys.NIF_ICON
	m.data.Icon = icon
	return m
}

// SetTip sets the tooltip string. Overrides custom tooltip.
func (m *Modifier) SetStringTip(tip string) *Modifier {
	var buf []win32.WCHAR
	m.data.Flags |= (sys.NIF_TIP | sys.NIF_SHOWTIP)
	win32util.CString(tip, &buf)
	copy(m.data.Tip[:len(m.data.Tip)-1], buf)
	return m
}

func (m *Modifier) SetNotify(spec *NotifySpec) *Modifier {
	setNotify(&m.data, spec)
	return m
}

// Add adds a notify icon to taskbar's status area.
// id can be any number to identify the icon.
func Add(w *window.Window, id win32.WORD, spec *Spec) error {
	if id == 0 {
		panic("id must > 0") // For [ParseCallback] HACK.
	}
	data := newData(spec)
	data.ID = win32.UINT(id)
	data.Wnd = w.HWND()
	data.Version = sys.NOTIFYICON_VERSION_4
	if err := sys.Shell_NotifyIconW(sys.NIM_ADD, data); err != nil {
		return err
	}
	if err := sys.Shell_NotifyIconW(sys.NIM_SETVERSION, data); err != nil {
		return err
	}
	if spec.OnEvent != nil {
		k := w.AddMsgListener(CallbackMessage, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
			spec.OnEvent(ParseCallback(wParam, lParam))
		})
		w.AddMsgListener(win32.WM_DESTROY, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
			k.Remove()
		})
	}
	return nil
}

// StartModify returns a [Modifier] which can be used to modify the notify icon
// identified by w and id.
// if customTip is true, sending [NIN_POPUPOPEN] instead of showing tooltip string.
func StartModify(w *window.Window, id win32.WORD, customTip bool) *Modifier {
	var ret Modifier
	ret.data.Size = win32.DWORD(unsafe.Sizeof(ret.data))
	ret.data.ID = win32.UINT(id)
	ret.data.Wnd = w.HWND()
	if !customTip {
		ret.data.Flags = sys.NIF_SHOWTIP
	}
	return &ret
}

func Delete(w *window.Window, id win32.WORD) error {
	return sys.Shell_NotifyIconW(sys.NIM_DELETE,
		&sys.NOTIFYICONDATAW{
			Size: win32.DWORD(unsafe.Sizeof(sys.NOTIFYICONDATAW{})),
			ID:   win32.UINT(id),
			Wnd:  w.HWND(),
		})
}

func SetFocus(w *window.Window, id win32.WORD) error {
	return sys.Shell_NotifyIconW(sys.NIM_SETFOCUS,
		&sys.NOTIFYICONDATAW{
			Size: win32.DWORD(unsafe.Sizeof(sys.NOTIFYICONDATAW{})),
			ID:   win32.UINT(id),
			Wnd:  w.HWND(),
		})
}

func newData(spec *Spec) *sys.NOTIFYICONDATAW {
	var ret sys.NOTIFYICONDATAW
	ret.Size = win32.DWORD(unsafe.Sizeof(ret))
	ret.Flags = sys.NIF_MESSAGE
	ret.CallbackMessage = CallbackMessage
	if spec.Icon != 0 {
		ret.Flags |= sys.NIF_ICON
		ret.Icon = spec.Icon
	}
	if spec.Tip != "" {
		var buf []win32.WCHAR
		ret.Flags |= (sys.NIF_TIP | sys.NIF_SHOWTIP)
		win32util.CString(spec.Tip, &buf)
		copy(ret.Tip[:len(ret.Tip)-1], buf)
	}
	setNotify(&ret, spec.Notify)
	return &ret
}

func setNotify(n *sys.NOTIFYICONDATAW, spec *NotifySpec) {
	if spec == nil {
		n.Flags &= ^sys.NIF_INFO
	} else {
		var buf []win32.WCHAR // shared string buffer.
		n.Flags |= sys.NIF_INFO
		if spec.Icon <= win32.HICON(sys.NIIF_ERROR) {
			n.InfoFlags |= sys.InfoFlag(spec.Icon)
		} else {
			n.InfoFlags |= sys.NIIF_USER
			n.BalloonIcon = spec.Icon
		}
		win32util.CString(spec.Message, &buf)
		copy(n.Info[:len(n.Info)-1], buf)
		win32util.CString(spec.Title, &buf)
		copy(n.InfoTitle[:len(n.InfoTitle)-1], buf)
		if spec.NoSound {
			n.InfoFlags |= sys.NIIF_NOSOUND
		}
		if spec.LargeIcon {
			n.InfoFlags |= sys.NIIF_LARGE_ICON
		}
		if spec.RealTime {
			n.Flags |= sys.NIF_REALTIME
		}
		if spec.RespectQuietTime {
			n.InfoFlags |= sys.NIIF_RESPECT_QUIET_TIME
		}
	}
}
