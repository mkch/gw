package dialog

import "github.com/mkch/gw/win32"

const (
	DS_ABSALIGN      win32.WindowStyle = 0x01
	DS_SYSMODAL      win32.WindowStyle = 0x02
	DS_LOCALEDIT     win32.WindowStyle = 0x20  /* 16-bit: Edit items get Local storage. 32-bit and up: meaningless. */
	DS_SETFONT       win32.WindowStyle = 0x40  /* User specified font for Dlg controls */
	DS_MODALFRAME    win32.WindowStyle = 0x80  /* Can be combined with WS_CAPTION  */
	DS_NOIDLEMSG     win32.WindowStyle = 0x100 /* WM_ENTERIDLE message will not be sent */
	DS_SETFOREGROUND win32.WindowStyle = 0x200
	DS_3DLOOK        win32.WindowStyle = 0x0004
	DS_FIXEDSYS      win32.WindowStyle = 0x0008
	DS_NOFAILCREATE  win32.WindowStyle = 0x0010
	DS_CONTROL       win32.WindowStyle = 0x0400
	DS_CENTER        win32.WindowStyle = 0x0800
	DS_CENTERMOUSE   win32.WindowStyle = 0x1000
	DS_CONTEXTHELP   win32.WindowStyle = 0x2000
	DS_SHELLFONT     win32.WindowStyle = (DS_SETFONT | DS_FIXEDSYS)
)
