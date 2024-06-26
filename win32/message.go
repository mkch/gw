package win32

const (
	WM_CREATE                  = 0x0001
	WM_DESTROY                 = 0x0002
	WM_SIZE                    = 0x0005
	WM_ACTIVATE                = 0x0006
	WM_SETFOCUS                = 0x0007
	WM_KILLFOCUS               = 0x0008
	WM_ENABLE                  = 0x000A
	WM_SETREDRAW               = 0x000B
	WM_SETTEXT                 = 0x000C
	WM_GETTEXT                 = 0x000D
	WM_GETTEXTLENGTH           = 0x000E
	WM_PAINT                   = 0x000F
	WM_CLOSE                   = 0x0010
	WM_QUIT                    = 0x0012
	WM_CONTEXTMENU             = 0x007B
	WM_COMMAND                 = 0x0111
	WM_NCDESTROY               = 0x0082
	WM_MOUSEFIRST              = 0x0200
	WM_MOUSEMOVE               = 0x0200
	WM_LBUTTONDOWN             = 0x0201
	WM_LBUTTONUP               = 0x0202
	WM_LBUTTONDBLCLK           = 0x0203
	WM_RBUTTONDOWN             = 0x0204
	WM_RBUTTONUP               = 0x0205
	WM_RBUTTONDBLCLK           = 0x0206
	WM_MBUTTONDOWN             = 0x0207
	WM_MBUTTONUP               = 0x0208
	WM_MBUTTONDBLCLK           = 0x0209
	WM_MOUSEWHEEL              = 0x020A
	WM_XBUTTONDOWN             = 0x020B
	WM_XBUTTONUP               = 0x020C
	WM_XBUTTONDBLCLK           = 0x020D
	WM_MOUSEHWHEEL             = 0x020E
	WM_MOUSELAST               = 0x020E
	WM_INITDIALOG              = 0x0110
	WM_USER                    = 0x0400
	WM_SETFONT                 = 0x0030
	WM_DPICHANGED              = 0x02E0
	WM_DPICHANGED_BEFOREPARENT = 0x02E2
	WM_DPICHANGED_AFTERPARENT  = 0x02E3
	WM_GETDPISCALEDSIZE        = 0x02E4
	WM_SIZING                  = 0x0214
	WM_NCLBUTTONDOWN           = 0x00A1
	WM_NCHITTEST               = 0x0084

	WM_APP = 0x8000

	DM_GETDEFID   = WM_USER + 0
	DM_SETDEFID   = WM_USER + 1
	DM_REPOSITION = WM_USER + 2

	BM_GETCHECK     = 0x00F0
	BM_SETCHECK     = 0x00F1
	BM_GETSTATE     = 0x00F2
	BM_SETSTATE     = 0x00F3
	BM_SETSTYLE     = 0x00F4
	BM_CLICK        = 0x00F5
	BM_GETIMAGE     = 0x00F6
	BM_SETIMAGE     = 0x00F7
	BM_SETDONTCLICK = 0x00F8
)

const (
	HTERROR       = -2
	HTTRANSPARENT = -1
	HTNOWHERE     = 0
	HTCLIENT      = 1
	HTCAPTION     = 2
	HTSYSMENU     = 3
	HTGROWBOX     = 4
	HTSIZE        = HTGROWBOX
	HTMENU        = 5
	HTHSCROLL     = 6
	HTVSCROLL     = 7
	HTMINBUTTON   = 8
	HTMAXBUTTON   = 9
	HTLEFT        = 10
	HTRIGHT       = 11
	HTTOP         = 12
	HTTOPLEFT     = 13
	HTTOPRIGHT    = 14
	HTBOTTOM      = 15
	HTBOTTOMLEFT  = 16
	HTBOTTOMRIGHT = 17
	HTBORDER      = 18
	HTREDUCE      = HTMINBUTTON
	HTZOOM        = HTMAXBUTTON
	HTSIZEFIRST   = HTLEFT
	HTSIZELAST    = HTBOTTOMRIGHT
	HTOBJECT      = 19
	HTCLOSE       = 20
	HTHELP        = 21
)

const (
	MK_CONTROL  = 0x0008
	MK_LBUTTON  = 0x0001
	MK_MBUTTON  = 0x0010
	MK_RBUTTON  = 0x0002
	MK_SHIFT    = 0x0004
	MK_XBUTTON1 = 0x0020
	MK_XBUTTON2 = 0x0040
)
