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
	customTip *bool
	data      sys.NOTIFYICONDATAW
}

// Apply applies the change.
func (m *Modifier) Apply() error {
	return sys.Shell_NotifyIconW(sys.NIM_MODIFY, &m.data)
}

func (m *Modifier) SetIcon(icon win32.HICON) *Modifier {
	m.data.Flags |= sys.NIF_ICON
	if !*m.customTip {
		m.data.Flags |= sys.NIF_SHOWTIP
	}
	m.data.Icon = icon
	return m
}

// SetTip sets the tooltip string.
// if tip is "", [NIN_POPUPOPEN] event is sent to allow showing application-drawn tooltip.
func (m *Modifier) SetTip(tip string) *Modifier {
	if tip == "" {
		*m.customTip = true
		m.data.Flags &= ^(sys.NIF_TIP | sys.NIF_SHOWTIP)
		return m
	}
	*m.customTip = false
	var buf []win32.WCHAR
	m.data.Flags |= (sys.NIF_TIP | sys.NIF_SHOWTIP)
	win32util.CString(tip, &buf)
	win32util.CopyCString(m.data.Tip[:], buf)
	return m
}

// SetNotify sets the ballon notification.
func (m *Modifier) SetNotify(spec *NotifySpec) *Modifier {
	if !*m.customTip {
		m.data.Flags |= sys.NIF_SHOWTIP
	}
	setNotify(&m.data, spec)
	return m
}

// NotifyIcon is an icon in taskbar's status area.
type NotifyIcon struct {
	w         win32.HWND
	id        win32.UINT
	customTip bool
}

// New adds an icon to taskbar's status area.
// An NotifyIcon is identified by a window and an ID.
func New(w *window.Window, id win32.WORD, spec *Spec) (*NotifyIcon, error) {
	if id == 0 {
		panic("id must > 0") // For [ParseCallback] HACK.
	}
	data := newData(spec, w.HWND(), id)
	data.Version = sys.NOTIFYICON_VERSION_4
	if err := sys.Shell_NotifyIconW(sys.NIM_ADD, data); err != nil {
		return nil, err
	}
	if err := sys.Shell_NotifyIconW(sys.NIM_SETVERSION, data); err != nil {
		return nil, err
	}
	if spec.OnEvent != nil {
		k := w.AddMsgListener(CallbackMessage, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
			spec.OnEvent(ParseCallback(wParam, lParam))
		})
		w.AddMsgListener(win32.WM_DESTROY, func(hwnd win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) {
			k.Remove()
		})
	}
	return &NotifyIcon{w.HWND(), win32.UINT(id), data.Flags&sys.NIF_SHOWTIP == 0}, nil
}

// StartModify returns a [Modifier] which can be used to modify the notify icon identified by w and id.
func (icon *NotifyIcon) StartModify() *Modifier {
	var ret Modifier
	ret.data.Size = win32.DWORD(unsafe.Sizeof(ret.data))
	ret.data.ID = icon.id
	ret.data.Wnd = icon.w
	ret.customTip = &icon.customTip
	return &ret
}

func (icon *NotifyIcon) ID() win32.UINT {
	return win32.UINT(icon.id)
}

// Delete removes the icon.
func (icon *NotifyIcon) Delete() error {
	return sys.Shell_NotifyIconW(sys.NIM_DELETE,
		&sys.NOTIFYICONDATAW{
			Size: win32.DWORD(unsafe.Sizeof(sys.NOTIFYICONDATAW{})),
			ID:   icon.id,
			Wnd:  icon.w,
		})
}

// SetFocus returns focus to the taskbar notification area.
// Notification area icons should use this method when they have completed their UI operation.
// For example, if the icon displays a shortcut menu, but the user presses ESC to cancel it, this method
// to return focus to the notification area.
func (icon *NotifyIcon) SetFocus() error {
	return sys.Shell_NotifyIconW(sys.NIM_SETFOCUS,
		&sys.NOTIFYICONDATAW{
			Size: win32.DWORD(unsafe.Sizeof(sys.NOTIFYICONDATAW{})),
			ID:   icon.id,
			Wnd:  icon.w,
		})
}

func newData(spec *Spec, hwnd win32.HWND, id win32.WORD) *sys.NOTIFYICONDATAW {
	var ret sys.NOTIFYICONDATAW
	ret.Size = win32.DWORD(unsafe.Sizeof(ret))
	ret.Flags = sys.NIF_MESSAGE
	ret.Wnd = hwnd
	ret.ID = win32.UINT(id)
	ret.CallbackMessage = CallbackMessage
	if spec.Icon != 0 {
		ret.Flags |= sys.NIF_ICON
		ret.Icon = spec.Icon
	}
	if spec.Tip != "" {
		var buf []win32.WCHAR
		ret.Flags |= (sys.NIF_TIP | sys.NIF_SHOWTIP)
		win32util.CString(spec.Tip, &buf)
		win32util.CopyCString(ret.Tip[:], buf)
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
		win32util.CopyCString(n.Info[:], buf)
		win32util.CString(spec.Title, &buf)
		win32util.CopyCString(n.InfoTitle[:], buf)
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
