package sys

import (
	"unsafe"

	"github.com/mkch/gw/win32"
	"github.com/mkch/gw/win32/sysutil"
	"golang.org/x/sys/windows"
)

var lzShell32 = windows.NewLazyDLL("shell32.dll")

const NOTIFYICON_VERSION_4 win32.UINT = 4

type NOTIFYICONDATAW struct {
	Size            win32.DWORD
	Wnd             win32.HWND
	ID              win32.UINT
	Flags           Flag
	CallbackMessage win32.UINT
	Icon            win32.HICON
	Tip             [128]win32.WCHAR
	State           State
	StateMask       State
	Info            [256]win32.WCHAR
	Version         win32.UINT
	InfoTitle       [64]win32.WCHAR
	InfoFlags       InfoFlag
	GuidItem        win32.GUID
	BalloonIcon     win32.HICON
}

type Message win32.DWORD

const (
	NIM_ADD        Message = 0x00000000
	NIM_MODIFY     Message = 0x00000001
	NIM_DELETE     Message = 0x00000002
	NIM_SETFOCUS   Message = 0x00000003
	NIM_SETVERSION Message = 0x00000004
)

type Flag win32.UINT

const (
	NIF_MESSAGE  Flag = 0x00000001
	NIF_ICON     Flag = 0x00000002
	NIF_TIP      Flag = 0x00000004
	NIF_STATE    Flag = 0x00000008
	NIF_INFO     Flag = 0x00000010
	NIF_GUID     Flag = 0x00000020
	NIF_REALTIME Flag = 0x00000040
	NIF_SHOWTIP  Flag = 0x00000080
)

type State win32.DWORD

const (
	NIS_HIDDEN     State = 0x00000001
	NIS_SHAREDICON State = 0x00000002
)

type InfoFlag win32.DWORD

const (
	NIIF_NONE               InfoFlag = 0x00000000
	NIIF_INFO               InfoFlag = 0x00000001
	NIIF_WARNING            InfoFlag = 0x00000002
	NIIF_ERROR              InfoFlag = 0x00000003
	NIIF_USER               InfoFlag = 0x00000004
	NIIF_ICON_MASK          InfoFlag = 0x0000000F
	NIIF_NOSOUND            InfoFlag = 0x00000010
	NIIF_LARGE_ICON         InfoFlag = 0x00000020
	NIIF_RESPECT_QUIET_TIME InfoFlag = 0x00000080
)

var lzShell_NotifyIconW = lzShell32.NewProc("Shell_NotifyIconW")

func Shell_NotifyIconW(message Message, data *NOTIFYICONDATAW) error {
	return sysutil.MustTrue(lzShell_NotifyIconW.Call(uintptr(message), uintptr(unsafe.Pointer(data))))
}
